package nfs

import (
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/clierr"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"
)

type volumeNFSExportConfigsCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string
	volumeID  string
	exports   []volume.NFSExportConfig

	// useCAS determines whether the command makes the set replicas request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string
	useAsync   bool

	writer io.Writer
}

func (c *volumeNFSExportConfigsCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	var nsID id.Namespace

	if useIDs {
		nsID = id.Namespace(c.namespace)
	} else {
		ns, err := c.client.GetNamespaceByName(ctx, c.namespace)
		if err != nil {
			return err
		}
		nsID = ns.ID
	}

	var volID id.Volume

	if useIDs {
		volID = id.Volume(c.volumeID)
	} else {
		vol, err := c.client.GetVolumeByName(ctx, nsID, c.volumeID)
		if err != nil {
			return err
		}
		volID = vol.ID
	}

	params := &apiclient.UpdateNFSVolumeExportsRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	err = c.client.UpdateNFSVolumeExports(ctx, nsID, volID, c.exports, params)
	if err != nil {
		return err
	}

	// Display the "request submitted" message if it was async, instead of
	// the deletion confirmation below.
	if c.useAsync {
		return c.display.AsyncRequest(ctx, c.writer)
	}

	return c.display.UpdateNFSVolumeExports(ctx, c.writer, volID, output.NewNFSExportConfigs(c.exports))
}

func parseExportString(s string) (volume.NFSExportConfig, error) {
	s = strings.TrimSpace(s)

	if s == "" {
		return volume.NFSExportConfig{}, errEmptyExportString
	}

	ff := strings.Split(s, ",")

	if len(ff) != 4 {
		return volume.NFSExportConfig{}, newErrInvalidExportConfigArg(s)
	}

	idExp, err := strconv.ParseUint(ff[0], 10, 32)
	if err != nil {
		return volume.NFSExportConfig{}, errWrongExportID
	}

	acls, err := parseACLs(ff[3])
	if err != nil {
		return volume.NFSExportConfig{}, err
	}

	return volume.NFSExportConfig{
		ExportID:   uint(idExp),
		Path:       ff[1],
		PseudoPath: ff[2],
		ACLs:       acls,
	}, nil
}

func parseACLs(s string) ([]volume.NFSExportConfigACL, error) {
	s = strings.TrimSpace(s)
	ss := strings.Split(s, "+")

	acls := make([]volume.NFSExportConfigACL, 0)

	for _, a := range ss {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}

		fields := strings.Split(a, ";")
		if len(fields) != 6 {
			return nil, newErrInvalidExportConfigArg(a)
		}

		identityType, matcher, squash, accessLevel := fields[0], fields[1], fields[4], fields[5]

		uid, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return nil, errWrongACLUID
		}

		gid, err := strconv.ParseInt(fields[3], 10, 64)
		if err != nil {
			return nil, errWrongACLGID
		}

		switch identityType {
		case "cidr", "hostname", "netgroup":
			// ok
		default:
			return nil, errWrongIdentityType
		}

		switch squash {
		case "none", "root", "rootuid", "all":
			// ok
		default:
			return nil, errWrongSquash
		}

		switch accessLevel {
		case "ro", "rw":
			// ok
		default:
			return nil, errWrongSquashAccessLevel
		}

		acls = append(acls, volume.NFSExportConfigACL{
			Identity: volume.NFSExportConfigACLIdentity{
				IdentityType: identityType,
				Matcher:      matcher,
			},
			SquashConfig: volume.NFSExportConfigACLSquashConfig{
				GID:    gid,
				UID:    uid,
				Squash: squash,
			},
			AccessLevel: accessLevel,
		})
	}

	return acls, nil
}

func newSetExports(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeNFSExportConfigsCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "exports [volume name] [export...]",
		Short: "Updates volume's NFS export configs overwriting existing. Write no export to empty the list.",
		Example: `
$ storageos nfs exports my-volume-name '1,/path,/pseudo' -n my-namespace-name

SYNOPSIS:
 exportString ::= ID,PATH,PSEUDOPATH,[ ACL [+ACL]... ]
 ACL          ::= [cidr|hostname|netgroup];[MATCHER];[UID];[GID];[none|root|rootuid|all];[rw|ro]

EXAMPLES:
    "1,/dir,/other,cidr;10.0.0.0/16;1000;1001;root;rw"
    Expose volume directory /dir as the pseudo path /other
    Allow writes from 10.0.*.*
    Map UID & GID 0 to 1000 & 1001 respectively

      id=1
      path=/dir
      psuedopath=/other
      and only one ACL with
        accessLevel=rw
        type=cidr
        matcher=10.0.0.0/16
        squash=root
        uid=1000
        gid=1001

    "2,/path,/other,hostname;*.storageos.com;0;0;none;ro+cidr;10.0.0.0/8;3;4;all;rw"
    Expose volume directory /path as the pseudo path /other
    Allow reads from *.storageos.com
    Allow writes from 10.*.*.* mapping all ops to uid 3 and gid 4

      id=2
      path=/path
      psuedopath=/other
      and two ACLs with
        accessLevel=ro
        type=hostname
        matcher=*.storageos.com
        squash=none
        uid=0
        gid=0

        accessLevel=rw
        type=cidr
        matcher=10.0.0.0/8
        squash=all
        uid=3
        gid=4
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {

			if len(args) < 1 {
				return clierr.NewErrInvalidArgNum(args, 1, "storageos nfs exports [volume name] [export...]")
			}

			c.volumeID = args[0]

			c.exports = make([]volume.NFSExportConfig, 0, len(args)-1)

			for i := 1; i < len(args); i++ {
				parsed, err := parseExportString(args[i])
				if err != nil {
					return newErrInvalidExportConfigArg(args[1])
				}
				c.exports = append(c.exports, parsed)
			}

			return nil
		}),

		PreRunE: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
			c.display = SelectDisplayer(c.config)

			ns, err := c.config.Namespace()
			if err != nil {
				return err
			}

			if ns == "" {
				return clierr.ErrNoNamespaceSpecified
			}
			c.namespace = ns

			return nil

		}),

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.EnsureNamespaceSetWhenUseIDs(c.config),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	c.useCAS = flagutil.SupportCAS(cobraCommand.Flags(), &c.casVersion)
	flagutil.SupportAsync(cobraCommand.Flags(), &c.useAsync)

	return cobraCommand
}
