package policy

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

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

// jsonlValidate scans the byte array, limited by newline, and validates each
// line to be a valid json object.
func jsonlValidate(data []byte) error {
	var jsonObj map[string]interface{}
	var jsonArray []interface{}

	if len(data) == 0 {
		return errors.New("empty JSON line input")
	}

	dataReader := bytes.NewReader(data)
	scanner := bufio.NewScanner(dataReader)
	for scanner.Scan() {
		byteData := scanner.Bytes()
		if err := json.Unmarshal(byteData, &jsonObj); err != nil {
			// Invalidate if it's a json array object.
			if errArray := json.Unmarshal(byteData, &jsonArray); errArray == nil {
				return errors.New("expected a json object per line, got a json array")
			}
			return err
		}
	}

	return nil
}

func runCreateFromFiles(storageosCli *command.StorageOSCli, opt createOptions) error {
	var jsonlFiles [][]byte

	for _, file := range opt.policies {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read policy file (No. %s): %s", file, err)
		}

		if err := jsonlValidate(buf); err != nil {
			return fmt.Errorf("failed to parse policy file (No. %s): %s", file, err)
		}

		jsonlFiles = append(jsonlFiles, buf)
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

	if err := jsonlValidate(buf); err != nil {
		return fmt.Errorf("failed to parse stdin: %s", err)
	}

	return sendJSONL(storageosCli, buf)
}

func sendJSONL(storageosCli *command.StorageOSCli, jsonl []byte) error {
	jsonlReader := bytes.NewReader(jsonl)
	scanner := bufio.NewScanner(jsonlReader)
	for scanner.Scan() {
		if err := storageosCli.Client().PolicyCreate(scanner.Bytes(), context.Background()); err != nil {
			return err
		}
	}
	return nil
}
