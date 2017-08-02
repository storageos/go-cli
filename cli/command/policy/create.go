package policy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"context"
	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

type stringSlice []string

func (s *stringSlice) Type() string {
	return "string"
}

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(val string) error {
	*s = append(*s, strings.Split(val, ",")...)
	return nil
}

type createOptions struct {
	user      string
	group     string
	namespace string
	policies  stringSlice
	stdin     bool
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{}

	cmd := &cobra.Command{
		Use: "create [OPTIONS]",
		Short: `Create a new policy, Either provide the set of policy files, set with options or write to stdin.
		E.g. "storageos policy create --user awesomeUser --namespace testing"
		E.g. "storageos policy create --policies='rules1.jsonl,rules2.jsonl'"
		E.g. "echo '{"spec": {"group": "devs", "namespace": "develop"}}' | storageos policy create --stdin"`,
		Args: cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(cmd, storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.user, "user", "", "User field for a policy entry")
	flags.StringVar(&opt.group, "group", "", "Group field for a policy entry")
	flags.StringVar(&opt.namespace, "namespace", "", "Namespace field for a policy entry")
	flags.Var(&opt.policies, "policies", "Provide a new (comma seperated) list of policy files in json line format.")
	flags.BoolVar(&opt.stdin, "stdin", false, "Read policy input from stdin")
	return cmd
}

func runCreate(cmd *cobra.Command, storageosCli *command.StorageOSCli, opt createOptions) error {
	switch {
	case opt.stdin:
		if len(opt.policies) > 0 {
			return fmt.Errorf("Please provide stdin or use other methods. (Not both)")
		}
		return runCreateFromStdin(storageosCli, opt)

	case len(opt.policies) > 0 && (opt.user+opt.group+opt.namespace) != "":
		return fmt.Errorf("Provide either policy file(s), or use the user/group/namespace flags. (Not both)")

	case len(opt.policies) > 0:
		return runCreateFromFiles(storageosCli, opt)

	case (opt.user + opt.group + opt.namespace) != "":
		return runCreateFromFlags(storageosCli, opt)

	default:
		return fmt.Errorf(
			"Please provide input files, use the user/group/namespace flags or set the stdin flag\nSee '%s --help'.\n\nUsage:  %s\n\n%s",
			cmd.CommandPath(),
			cmd.UseLine(),
			cmd.Short,
		)
	}
}

// Attempts to parse JSONish data (multiple formats) to jsonl format
func jsonlParse(data []byte) ([]byte, error) {
	type jsonList []*json.RawMessage

	isJSONType := func(t interface{}, data []byte) bool {
		return json.Unmarshal(data, t) == nil
	}

	isJSONL := func(data []byte) bool {
		split := bytes.Split(data, []byte("\n"))
		obj := &json.RawMessage{}
		for _, elm := range split {
			if !isJSONType(obj, elm) {
				return false
			}
		}
		return true
	}

	if isJSONL(data) {
		return data, nil
	}

	list := jsonList{}
	obj := &json.RawMessage{}

	var concat []byte
	if isJSONType(&list, data) {
		for _, obj := range list {
			deref := *((*[]byte)(obj))
			concat = append(concat, deref...)
			concat = append(concat, '\n')
		}
		return concat, nil
	}

	if isJSONType(&obj, data) {
		deref := *((*[]byte)(obj))
		concat = append(concat, deref...)
		concat = append(concat, '\n')
		return concat, nil
	}

	return nil, fmt.Errorf("unknown data format")
}

func runCreateFromFiles(storageosCli *command.StorageOSCli, opt createOptions) error {
	var jsonlFiles [][]byte

	for _, file := range opt.policies {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read policy file (No. %s): %s", file, err)
		}

		jsonlBuf, err := jsonlParse(buf)
		if err != nil {
			return fmt.Errorf("failed to parse policy file (No. %s): %s", file, err)
		}

		jsonlFiles = append(jsonlFiles, jsonlBuf)
	}

	return sendJSONL(storageosCli, bytes.Join(jsonlFiles, []byte("\n")))
}

func runCreateFromFlags(storageosCli *command.StorageOSCli, opt createOptions) error {
	pol := types.Policy{}
	pol.Spec.User = opt.user
	pol.Spec.Group = opt.group
	pol.Spec.Namespace = opt.namespace

	data, err := json.Marshal(&pol)
	if err != nil {
		return err
	}

	return sendJSONL(storageosCli, append(data, '\n'))
}

func runCreateFromStdin(storageosCli *command.StorageOSCli, opt createOptions) error {
	buf, err := ioutil.ReadAll(storageosCli.In())
	if err != nil {
		return fmt.Errorf("failed to read stdin: %s", err)
	}

	jsonlBuf, err := jsonlParse(buf)
	if err != nil {
		return fmt.Errorf("failed to parse stdin: %s", err)
	}

	return sendJSONL(storageosCli, jsonlBuf)
}

func sendJSONL(storageosCli *command.StorageOSCli, jsonl []byte) error {
	return storageosCli.Client().PolicyCreate(jsonl, context.Background())
}
