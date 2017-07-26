package policy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"context"
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	//"github.com/storageos/go-cli/cli/opts"
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
	policies stringSlice
	stdin    bool
	args     []string
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{}

	cmd := &cobra.Command{
		Use: "create [jsonPolicy | jsonPolicyList]...",
		Short: `Create a new policy, Either provide the set of policy files, provide json input or write to stdin.
		E.g. "storageos policy create --policies='rules1.jsonl,rules2.jsonl'"`,
		Args: cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.args = args
			return runCreate(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.Var(&opt.policies, "policies", "Provide a new (comma seperated) list of policy files in json line format.")
	flags.BoolVar(&opt.stdin, "stdin", false, "Read policy input from stdin")
	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {
	switch {
	case opt.stdin:
		if len(opt.policies)+len(opt.args) > 0 {
			return fmt.Errorf("Please provide stdin or use other methods. (Not both)")
		}
		return runCreateFromStdin(storageosCli, opt)

	case len(opt.policies) > 0 && len(opt.args) > 0:
		return fmt.Errorf("Provide either a policy file, or a positional arg. (Not both)")

	case len(opt.policies) > 0:
		return runCreateFromFiles(storageosCli, opt)

	case len(opt.args) > 0:
		return runCreateFromArg(storageosCli, opt)

	default:
		return fmt.Errorf("Please provide input files, positional args or the stdin flag")
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

func runCreateFromArg(storageosCli *command.StorageOSCli, opt createOptions) error {
	type jsonList []*json.RawMessage
	type jsonObject *json.RawMessage

	var jsonlArgs [][]byte

	for i, policy := range opt.args {
		buf, err := jsonlParse([]byte(policy))
		if err != nil {
			return fmt.Errorf("error parsing positional arg %d: %s", i, err)
		}

		jsonlArgs = append(jsonlArgs, buf)
	}

	return sendJSONL(storageosCli, bytes.Join(jsonlArgs, []byte("\n")))
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
