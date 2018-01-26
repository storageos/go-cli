package volume

import (
	"errors"
	"fmt"
	"strings"

	"context"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
)

type createOptions struct {
	name         string
	description  string
	size         int
	pool         string
	fsType       string
	namespace    string
	nodeSelector string
	labels       opts.ListOpts
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{
		labels: opts.NewListOpts(opts.ValidateEnv),
	}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] [VOLUME]",
		Short: "Create a volume",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var posarg string
			if len(args) > 0 {
				posarg = args[0]
			}

			var err error
			opt.namespace, opt.name, err = parseNamespaceVolume(opt.namespace, opt.name, posarg)
			if err != nil {
				return err
			}

			return runCreate(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.name, "name", "", "Volume name")
	flags.Lookup("name").Hidden = true
	flags.StringVarP(&opt.description, "description", "d", "", "Volume description")
	flags.IntVarP(&opt.size, "size", "s", 5, "Volume size in GB")
	flags.StringVarP(&opt.pool, "pool", "p", "default", "Volume capacity pool")
	flags.StringVarP(&opt.fsType, "fstype", "f", "", "Requested filesystem type")
	flags.StringVarP(&opt.namespace, "namespace", "n", "", `Volume namespace (default "default")`)
	flags.StringVar(&opt.nodeSelector, "nodeSelector", "", "Node selector")
	flags.Var(&opt.labels, "label", "Set metadata (key=value pairs) on the volume")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {
	client := storageosCli.Client()

	params := types.VolumeCreateOptions{
		Name:         opt.name,
		Description:  opt.description,
		Size:         opt.size,
		Pool:         opt.pool,
		FSType:       opt.fsType,
		Namespace:    opt.namespace,
		NodeSelector: opt.nodeSelector,
		Labels:       opts.ConvertKVStringsToMap(opt.labels.GetAll()),
		Context:      context.Background(),
	}

	vol, err := client.VolumeCreate(params)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "%s/%s\n", vol.Namespace, vol.Name)
	return nil
}

func parseNamespaceVolume(nsflag, vnflag, posarg string) (namespace string, volume string, err error) {
	switch {
	case posarg != "" && vnflag != "":
		return "", "", errors.New("Conflicting options: either specify --name or provide positional arg, not both\n")

	case posarg != "":
		split := strings.Split(posarg, "/")

		switch {
		case len(split) > 1 && nsflag != "":
			return "", "", errors.New("Conflicting options: either specify --namespace or use 'namespace/volumename' positional arg, not both\n")

		case len(split) > 1:
			return split[0], split[1], nil

		case nsflag != "":
			return nsflag, posarg, nil

		default:
			return "default", posarg, nil
		}

	case vnflag != "" && nsflag != "":
		return nsflag, vnflag, nil

	case vnflag != "":
		return "default", vnflag, nil

	default:
		return "", "", errors.New("Please provide a volume name\n")
	}
}
