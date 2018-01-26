package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dnephin/cobra"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/storageos/go-api/soserror"
	"github.com/storageos/go-api/types/versions"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/commands"
	cliconfig "github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/cli/debug"
	cliflags "github.com/storageos/go-cli/cli/flags"
	"github.com/storageos/go-cli/pkg/term"
	"github.com/storageos/go-cli/version"
)

var shortDesc = `Converged storage for containers. 

By using this product, you are agreeing to the terms of the the StorageOS Ltd. End
User Subscription Agreement (EUSA) found at: https://storageos.com/legal/#eusa

WARNING: This is the beta version of StorageOS and should not be used in production.
To be notified about stable releases and latest features, sign up at https://my.storageos.com.
`

// Disable degugging (logging to stdout) until enabled.  In normal use we don't
// want logrus messages going to stdout/stderr.
func init() {
	debug.Disable()
}

func isCoreOS() (bool, error) {
	f, err := ioutil.ReadFile("/etc/lsb-release")
	if err != nil {
		return false, err
	}

	return strings.Contains(string(f), "DISTRIB_ID=CoreOS"), nil
}

func isInContainer() (bool, error) {
	f, err := ioutil.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false, err
	}

	// TODO: How reliable is this method of detection. Is there a better way?
	return strings.Contains(string(f), "docker"), nil
}

func verfyHostPlatform() error {
	// Detect native execution on coreOS
	// coreOS should not run user-land programs, and we will not work there (outside a container)
	if coreos, err := isCoreOS(); err == nil && coreos {

		// If we dont think we are in a container, fail and warn the user
		if inContainer, err := isInContainer(); err == nil && !inContainer {
			return errors.New("To use the StorageOS CLI on Container Linux, you need to run the storageos/cli image.")
		}
	}
	return nil
}

func newStorageOSCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opts := cliflags.NewClientOptions()
	var flags *pflag.FlagSet

	cmd := &cobra.Command{
		Use:              "storageos [OPTIONS] COMMAND [ARG...]",
		Short:            shortDesc,
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		Args:             noArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				showVersion()
				return nil
			}
			return storageosCli.ShowHelp(cmd, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// verify the host platform, exit immediately if known to be incompatible
			if err := verfyHostPlatform(); err != nil {
				return err
			}

			// flags must be the top-level command flags, not cmd.Flags()
			opts.Common.SetDefaultOptions(flags)
			preRun(opts)
			if err := storageosCli.Initialize(opts); err != nil {
				return err
			}
			return isSupported(cmd, storageosCli.Client().ClientVersion(), storageosCli.HasExperimental())
		},
	}
	cli.SetupRootCommand(cmd)

	flags = cmd.Flags()
	flags.BoolVarP(&opts.Version, "version", "v", false, "Print version information and quit")
	flags.StringVar(&opts.ConfigDir, "config", cliconfig.Dir(), "Location of client config files")
	opts.Common.InstallFlags(flags)

	setFlagErrorFunc(storageosCli, cmd, flags, opts)

	// setHelpFunc(storageosCli, cmd, flags, opts)

	cmd.SetOutput(storageosCli.Out())

	commands.AddCommands(cmd, storageosCli)

	setValidateArgs(storageosCli, cmd, flags, opts)

	return cmd
}

func setFlagErrorFunc(storageosCli *command.StorageOSCli, cmd *cobra.Command, flags *pflag.FlagSet, opts *cliflags.ClientOptions) {
	// When invoking `storageos volume --nonsense`, we need to make sure FlagErrorFunc return appropriate
	// output if the feature is not supported.
	// As above cli.SetupRootCommand(cmd) have already setup the FlagErrorFunc, we will add a pre-check before the FlagErrorFunc
	// is called.
	flagErrorFunc := cmd.FlagErrorFunc()
	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		initializeStorageOSCli(storageosCli, flags, opts)
		if err := isSupported(cmd, storageosCli.Client().ClientVersion(), storageosCli.HasExperimental()); err != nil {
			return err
		}
		return flagErrorFunc(cmd, err)
	})
}

func setHelpFunc(storageosCli *command.StorageOSCli, cmd *cobra.Command, flags *pflag.FlagSet, opts *cliflags.ClientOptions) {
	cmd.SetHelpFunc(func(ccmd *cobra.Command, args []string) {
		initializeStorageOSCli(storageosCli, flags, opts)
		fmt.Printf("VC: %s\n", storageosCli.Client().ClientVersion())
		fmt.Printf("HE: %t\n", storageosCli.HasExperimental())
		if err := isSupported(ccmd, storageosCli.Client().ClientVersion(), storageosCli.HasExperimental()); err != nil {
			fmt.Printf("ERRROR: %v\n", err)
			ccmd.Println(err)
			return
		}

		hideUnsupportedFeatures(ccmd, storageosCli.Client().ClientVersion(), storageosCli.HasExperimental())

		if err := ccmd.Help(); err != nil {
			ccmd.Println(err)
		}
	})
}

func setValidateArgs(storageosCli *command.StorageOSCli, cmd *cobra.Command, flags *pflag.FlagSet, opts *cliflags.ClientOptions) {
	// The Args is handled by ValidateArgs in cobra, which does not allows a pre-hook.
	// As a result, here we replace the existing Args validation func to a wrapper,
	// where the wrapper will check to see if the feature is supported or not.
	// The Args validation error will only be returned if the feature is supported.
	visitAll(cmd, func(ccmd *cobra.Command) {
		// if there is no tags for a command or any of its parent,
		// there is no need to wrap the Args validation.
		if !hasTags(ccmd) {
			return
		}

		if ccmd.Args == nil {
			return
		}

		cmdArgs := ccmd.Args
		ccmd.Args = func(cmd *cobra.Command, args []string) error {
			initializeStorageOSCli(storageosCli, flags, opts)
			if err := isSupported(cmd, storageosCli.Client().ClientVersion(), storageosCli.HasExperimental()); err != nil {
				return err
			}
			return cmdArgs(cmd, args)
		}
	})
}

func initializeStorageOSCli(storageosCli *command.StorageOSCli, flags *pflag.FlagSet, opts *cliflags.ClientOptions) {
	if storageosCli.Client() == nil { // when using --help, PersistentPreRun is not called, so initialization is needed.
		// flags must be the top-level command flags, not cmd.Flags()
		opts.Common.SetDefaultOptions(flags)
		preRun(opts)
		storageosCli.Initialize(opts)
	}
}

// visitAll will traverse all commands from the root.
// This is different from the VisitAll of cobra.Command where only parents
// are checked.
func visitAll(root *cobra.Command, fn func(*cobra.Command)) {
	for _, cmd := range root.Commands() {
		visitAll(cmd, fn)
	}
	fn(root)
}

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf("storageos: '%s' is not a valid command.\nSee 'storageos --help'", args[0])
}

func main() {
	// Set terminal emulation based on platform as required.
	stdin, stdout, stderr := term.StdStreams()
	logrus.SetOutput(stderr)

	storageosCli := command.NewStorageOSCli(stdin, stdout, stderr)
	cmd := newStorageOSCommand(storageosCli)

	if err := cmd.Execute(); err != nil {
		if customError, ok := err.(soserror.StorageOSError); ok {
			if msg := customError.String(); msg != "" {
				fmt.Fprintf(stderr, "error: %s\n", msg)
			}
			if cause := customError.Err(); cause != nil {
				fmt.Fprintf(stderr, "\ncaused by: %s\n", cause)
			}
			if help := customError.Help(); help != "" {
				fmt.Fprintf(stderr, "\n%s\n", help)
			}
			os.Exit(1)
		}

		if sterr, ok := err.(cli.StatusError); ok {
			if sterr.Status != "" {
				fmt.Fprintln(stderr, sterr.Status)
			}
			// StatusError should only be used for errors, and all errors should
			// have a non-zero exit status, so never exit with 0
			if sterr.StatusCode == 0 {
				os.Exit(1)
			}
			os.Exit(sterr.StatusCode)
		}
		fmt.Fprintln(stderr, err)
		os.Exit(1)
	}
}

func showVersion() {
	fmt.Printf("StorageOS version %s, build %s\n", version.Version, version.Revision)
	// TODO: better version
	// fmt.Printf("StorageOS API version %s\n", storageosCli.Client().ClientVersion())
}

func preRun(opts *cliflags.ClientOptions) {
	cliflags.SetLogLevel(opts.Common.LogLevel)

	if opts.ConfigDir != "" {
		cliconfig.SetDir(opts.ConfigDir)
	}

	if opts.Common.Debug {
		debug.Enable()
	}
}

func hideUnsupportedFeatures(cmd *cobra.Command, clientVersion string, hasExperimental bool) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// hide experimental flags
		if !hasExperimental {
			if _, ok := f.Annotations["experimental"]; ok {
				f.Hidden = true
			}
		}

		// hide flags not supported by the server
		if !isFlagSupported(f, clientVersion) {
			f.Hidden = true
		}

	})

	for _, subcmd := range cmd.Commands() {
		// hide experimental subcommands
		if !hasExperimental {
			if _, ok := subcmd.Tags["experimental"]; ok {
				subcmd.Hidden = true
			}
		}

		// hide subcommands not supported by the server
		if subcmdVersion, ok := subcmd.Tags["version"]; ok && versions.LessThan(clientVersion, subcmdVersion) {
			subcmd.Hidden = true
		}
	}
}

func isSupported(cmd *cobra.Command, clientVersion string, hasExperimental bool) error {
	// We check recursively so that, e.g., `storageos volume ls` will return the same output as `storageos volume`
	if !hasExperimental {
		for curr := cmd; curr != nil; curr = curr.Parent() {
			if _, ok := curr.Tags["experimental"]; ok {
				fmt.Print("e")
				return errors.New("only supported on a StorageOS with experimental features enabled")
			}
		}
	}

	if cmdVersion, ok := cmd.Tags["version"]; ok && versions.LessThan(clientVersion, cmdVersion) {
		fmt.Print("ERR: api version\n")
		return fmt.Errorf("requires API version %s, but the StorageOS API version is %s", cmdVersion, clientVersion)
	}

	// errs := []string{}

	// cmd.Flags().VisitAll(func(f *pflag.Flag) {
	// 	if f.Changed {
	// 		if !isFlagSupported(f, clientVersion) {
	// 			errs = append(errs, fmt.Sprintf("\"--%s\" requires API version %s, but the StorageOS API version is %s", f.Name, getFlagVersion(f), clientVersion))
	// 			return
	// 		}
	// 		if _, ok := f.Annotations["experimental"]; ok && !hasExperimental {
	// 			errs = append(errs, fmt.Sprintf("\"--%s\" is only supported on StorageOS with experimental features enabled", f.Name))
	// 		}
	// 	}
	// })
	// if len(errs) > 0 {
	// 	return errors.New(strings.Join(errs, "\n"))
	// }

	return nil
}

func getFlagVersion(f *pflag.Flag) string {
	if flagVersion, ok := f.Annotations["version"]; ok && len(flagVersion) == 1 {
		return flagVersion[0]
	}
	return ""
}

func isFlagSupported(f *pflag.Flag, clientVersion string) bool {
	if v := getFlagVersion(f); v != "" {
		return versions.GreaterThanOrEqualTo(clientVersion, v)
	}
	return true
}

// hasTags return true if any of the command's parents has tags
func hasTags(cmd *cobra.Command) bool {
	for curr := cmd; curr != nil; curr = curr.Parent() {
		if len(curr.Tags) > 0 {
			return true
		}
	}

	return false
}
