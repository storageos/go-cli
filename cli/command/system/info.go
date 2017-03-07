package system

//
// import (
// 	"fmt"
// 	"sort"
// 	"strings"
// 	"time"
//
// 	"golang.org/x/net/context"
//
// 	"github.com/dnephin/cobra"
// 	"github.com/docker/docker/api/types/swarm"
// 	"github.com/docker/docker/pkg/ioutils"
// 	"github.com/docker/go-units"
// 	"github.com/storageos/go-api/types"
// 	"github.com/storageos/go-cli/cli"
// 	"github.com/storageos/go-cli/cli/command"
// 	"github.com/storageos/go-cli/cli/debug"
// 	"github.com/storageos/go-cli/pkg/templates"
// )
//
// type infoOptions struct {
// 	format string
// }
//
// // NewInfoCommand creates a new cobra.Command for `docker info`
// func NewInfoCommand(storageosCli *command.StorageOSCli) *cobra.Command {
// 	var opts infoOptions
//
// 	cmd := &cobra.Command{
// 		Use:   "info [OPTIONS]",
// 		Short: "Display system-wide information",
// 		Args:  cli.NoArgs,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return runInfo(storageosCli, &opts)
// 		},
// 	}
//
// 	flags := cmd.Flags()
//
// 	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")
//
// 	return cmd
// }
//
// func runInfo(storageosCli *command.StorageOSCli, opts *infoOptions) error {
// 	ctx := context.Background()
// 	info, err := storageosCli.Client().Info(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	if opts.format == "" {
// 		return prettyPrintInfo(storageosCli, info)
// 	}
// 	return formatInfo(storageosCli, info, opts.format)
// }
//
// func prettyPrintInfo(storageosCli *command.StorageOSCli, info types.Info) error {
// 	fmt.Fprintf(storageosCli.Out(), "Containers: %d\n", info.Containers)
// 	fmt.Fprintf(storageosCli.Out(), " Running: %d\n", info.ContainersRunning)
// 	fmt.Fprintf(storageosCli.Out(), " Paused: %d\n", info.ContainersPaused)
// 	fmt.Fprintf(storageosCli.Out(), " Stopped: %d\n", info.ContainersStopped)
// 	fmt.Fprintf(storageosCli.Out(), "Images: %d\n", info.Images)
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Server Version: %s\n", info.ServerVersion)
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Storage Driver: %s\n", info.Driver)
// 	if info.DriverStatus != nil {
// 		for _, pair := range info.DriverStatus {
// 			fmt.Fprintf(storageosCli.Out(), " %s: %s\n", pair[0], pair[1])
//
// 			// print a warning if devicemapper is using a loopback file
// 			if pair[0] == "Data loop file" {
// 				fmt.Fprintln(storageosCli.Err(), " WARNING: Usage of loopback devices is strongly discouraged for production use. Use `--storage-opt dm.thinpooldev` to specify a custom block storage device.")
// 			}
// 		}
//
// 	}
// 	if info.SystemStatus != nil {
// 		for _, pair := range info.SystemStatus {
// 			fmt.Fprintf(storageosCli.Out(), "%s: %s\n", pair[0], pair[1])
// 		}
// 	}
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Logging Driver: %s\n", info.LoggingDriver)
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Cgroup Driver: %s\n", info.CgroupDriver)
//
// 	fmt.Fprintf(storageosCli.Out(), "Plugins: \n")
// 	fmt.Fprintf(storageosCli.Out(), " Volume:")
// 	fmt.Fprintf(storageosCli.Out(), " %s", strings.Join(info.Plugins.Volume, " "))
// 	fmt.Fprintf(storageosCli.Out(), "\n")
// 	fmt.Fprintf(storageosCli.Out(), " Network:")
// 	fmt.Fprintf(storageosCli.Out(), " %s", strings.Join(info.Plugins.Network, " "))
// 	fmt.Fprintf(storageosCli.Out(), "\n")
//
// 	if len(info.Plugins.Authorization) != 0 {
// 		fmt.Fprintf(storageosCli.Out(), " Authorization:")
// 		fmt.Fprintf(storageosCli.Out(), " %s", strings.Join(info.Plugins.Authorization, " "))
// 		fmt.Fprintf(storageosCli.Out(), "\n")
// 	}
//
// 	fmt.Fprintf(storageosCli.Out(), "Swarm: %v\n", info.Swarm.LocalNodeState)
// 	if info.Swarm.LocalNodeState != swarm.LocalNodeStateInactive && info.Swarm.LocalNodeState != swarm.LocalNodeStateLocked {
// 		fmt.Fprintf(storageosCli.Out(), " NodeID: %s\n", info.Swarm.NodeID)
// 		if info.Swarm.Error != "" {
// 			fmt.Fprintf(storageosCli.Out(), " Error: %v\n", info.Swarm.Error)
// 		}
// 		fmt.Fprintf(storageosCli.Out(), " Is Manager: %v\n", info.Swarm.ControlAvailable)
// 		if info.Swarm.ControlAvailable && info.Swarm.Error == "" && info.Swarm.LocalNodeState != swarm.LocalNodeStateError {
// 			fmt.Fprintf(storageosCli.Out(), " ClusterID: %s\n", info.Swarm.Cluster.ID)
// 			fmt.Fprintf(storageosCli.Out(), " Managers: %d\n", info.Swarm.Managers)
// 			fmt.Fprintf(storageosCli.Out(), " Nodes: %d\n", info.Swarm.Nodes)
// 			fmt.Fprintf(storageosCli.Out(), " Orchestration:\n")
// 			taskHistoryRetentionLimit := int64(0)
// 			if info.Swarm.Cluster.Spec.Orchestration.TaskHistoryRetentionLimit != nil {
// 				taskHistoryRetentionLimit = *info.Swarm.Cluster.Spec.Orchestration.TaskHistoryRetentionLimit
// 			}
// 			fmt.Fprintf(storageosCli.Out(), "  Task History Retention Limit: %d\n", taskHistoryRetentionLimit)
// 			fmt.Fprintf(storageosCli.Out(), " Raft:\n")
// 			fmt.Fprintf(storageosCli.Out(), "  Snapshot Interval: %d\n", info.Swarm.Cluster.Spec.Raft.SnapshotInterval)
// 			if info.Swarm.Cluster.Spec.Raft.KeepOldSnapshots != nil {
// 				fmt.Fprintf(storageosCli.Out(), "  Number of Old Snapshots to Retain: %d\n", *info.Swarm.Cluster.Spec.Raft.KeepOldSnapshots)
// 			}
// 			fmt.Fprintf(storageosCli.Out(), "  Heartbeat Tick: %d\n", info.Swarm.Cluster.Spec.Raft.HeartbeatTick)
// 			fmt.Fprintf(storageosCli.Out(), "  Election Tick: %d\n", info.Swarm.Cluster.Spec.Raft.ElectionTick)
// 			fmt.Fprintf(storageosCli.Out(), " Dispatcher:\n")
// 			fmt.Fprintf(storageosCli.Out(), "  Heartbeat Period: %s\n", units.HumanDuration(time.Duration(info.Swarm.Cluster.Spec.Dispatcher.HeartbeatPeriod)))
// 			fmt.Fprintf(storageosCli.Out(), " CA Configuration:\n")
// 			fmt.Fprintf(storageosCli.Out(), "  Expiry Duration: %s\n", units.HumanDuration(info.Swarm.Cluster.Spec.CAConfig.NodeCertExpiry))
// 			if len(info.Swarm.Cluster.Spec.CAConfig.ExternalCAs) > 0 {
// 				fmt.Fprintf(storageosCli.Out(), "  External CAs:\n")
// 				for _, entry := range info.Swarm.Cluster.Spec.CAConfig.ExternalCAs {
// 					fmt.Fprintf(storageosCli.Out(), "    %s: %s\n", entry.Protocol, entry.URL)
// 				}
// 			}
// 		}
// 		fmt.Fprintf(storageosCli.Out(), " Node Address: %s\n", info.Swarm.NodeAddr)
// 		managers := []string{}
// 		for _, entry := range info.Swarm.RemoteManagers {
// 			managers = append(managers, entry.Addr)
// 		}
// 		if len(managers) > 0 {
// 			sort.Strings(managers)
// 			fmt.Fprintf(storageosCli.Out(), " Manager Addresses:\n")
// 			for _, entry := range managers {
// 				fmt.Fprintf(storageosCli.Out(), "  %s\n", entry)
// 			}
// 		}
// 	}
//
// 	if len(info.Runtimes) > 0 {
// 		fmt.Fprintf(storageosCli.Out(), "Runtimes:")
// 		for name := range info.Runtimes {
// 			fmt.Fprintf(storageosCli.Out(), " %s", name)
// 		}
// 		fmt.Fprint(storageosCli.Out(), "\n")
// 		fmt.Fprintf(storageosCli.Out(), "Default Runtime: %s\n", info.DefaultRuntime)
// 	}
//
// 	if info.OSType == "linux" {
// 		fmt.Fprintf(storageosCli.Out(), "Init Binary: %v\n", info.InitBinary)
//
// 		for _, ci := range []struct {
// 			Name   string
// 			Commit types.Commit
// 		}{
// 			{"containerd", info.ContainerdCommit},
// 			{"runc", info.RuncCommit},
// 			{"init", info.InitCommit},
// 		} {
// 			fmt.Fprintf(storageosCli.Out(), "%s version: %s", ci.Name, ci.Commit.ID)
// 			if ci.Commit.ID != ci.Commit.Expected {
// 				fmt.Fprintf(storageosCli.Out(), " (expected: %s)", ci.Commit.Expected)
// 			}
// 			fmt.Fprintf(storageosCli.Out(), "\n")
// 		}
// 		if len(info.SecurityOptions) != 0 {
// 			kvs, err := types.DecodeSecurityOptions(info.SecurityOptions)
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Fprintf(storageosCli.Out(), "Security Options:\n")
// 			for _, so := range kvs {
// 				fmt.Fprintf(storageosCli.Out(), " %s\n", so.Name)
// 				for _, o := range so.Options {
// 					switch o.Key {
// 					case "profile":
// 						if o.Value != "default" {
// 							fmt.Fprintf(storageosCli.Err(), "  WARNING: You're not using the default seccomp profile\n")
// 						}
// 						fmt.Fprintf(storageosCli.Out(), "  Profile: %s\n", o.Value)
// 					}
// 				}
// 			}
// 		}
// 	}
//
// 	// Isolation only has meaning on a Windows daemon.
// 	if info.OSType == "windows" {
// 		fmt.Fprintf(storageosCli.Out(), "Default Isolation: %v\n", info.Isolation)
// 	}
//
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Kernel Version: %s\n", info.KernelVersion)
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Operating System: %s\n", info.OperatingSystem)
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "OSType: %s\n", info.OSType)
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Architecture: %s\n", info.Architecture)
// 	fmt.Fprintf(storageosCli.Out(), "CPUs: %d\n", info.NCPU)
// 	fmt.Fprintf(storageosCli.Out(), "Total Memory: %s\n", units.BytesSize(float64(info.MemTotal)))
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Name: %s\n", info.Name)
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "ID: %s\n", info.ID)
// 	fmt.Fprintf(storageosCli.Out(), "Docker Root Dir: %s\n", info.DockerRootDir)
// 	fmt.Fprintf(storageosCli.Out(), "Debug Mode (client): %v\n", debug.IsEnabled())
// 	fmt.Fprintf(storageosCli.Out(), "Debug Mode (server): %v\n", info.Debug)
//
// 	if info.Debug {
// 		fmt.Fprintf(storageosCli.Out(), " File Descriptors: %d\n", info.NFd)
// 		fmt.Fprintf(storageosCli.Out(), " Goroutines: %d\n", info.NGoroutines)
// 		fmt.Fprintf(storageosCli.Out(), " System Time: %s\n", info.SystemTime)
// 		fmt.Fprintf(storageosCli.Out(), " EventsListeners: %d\n", info.NEventsListener)
// 	}
//
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Http Proxy: %s\n", info.HTTPProxy)
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "Https Proxy: %s\n", info.HTTPSProxy)
// 	ioutils.FprintfIfNotEmpty(storageosCli.Out(), "No Proxy: %s\n", info.NoProxy)
//
// 	if info.IndexServerAddress != "" {
// 		u := storageosCli.ConfigFile().AuthConfigs[info.IndexServerAddress].Username
// 		if len(u) > 0 {
// 			fmt.Fprintf(storageosCli.Out(), "Username: %v\n", u)
// 		}
// 		fmt.Fprintf(storageosCli.Out(), "Registry: %v\n", info.IndexServerAddress)
// 	}
//
// 	// Only output these warnings if the server does not support these features
// 	if info.OSType != "windows" {
// 		if !info.MemoryLimit {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: No memory limit support")
// 		}
// 		if !info.SwapLimit {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: No swap limit support")
// 		}
// 		if !info.KernelMemory {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: No kernel memory limit support")
// 		}
// 		if !info.OomKillDisable {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: No oom kill disable support")
// 		}
// 		if !info.CPUCfsQuota {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: No cpu cfs quota support")
// 		}
// 		if !info.CPUCfsPeriod {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: No cpu cfs period support")
// 		}
// 		if !info.CPUShares {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: No cpu shares support")
// 		}
// 		if !info.CPUSet {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: No cpuset support")
// 		}
// 		if !info.IPv4Forwarding {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: IPv4 forwarding is disabled")
// 		}
// 		if !info.BridgeNfIptables {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: bridge-nf-call-iptables is disabled")
// 		}
// 		if !info.BridgeNfIP6tables {
// 			fmt.Fprintln(storageosCli.Err(), "WARNING: bridge-nf-call-ip6tables is disabled")
// 		}
// 	}
//
// 	if info.Labels != nil {
// 		fmt.Fprintln(storageosCli.Out(), "Labels:")
// 		for _, attribute := range info.Labels {
// 			fmt.Fprintf(storageosCli.Out(), " %s\n", attribute)
// 		}
// 		// TODO: Engine labels with duplicate keys has been deprecated in 1.13 and will be error out
// 		// after 3 release cycles (1.16). For now, a WARNING will be generated. The following will
// 		// be removed eventually.
// 		labelMap := map[string]string{}
// 		for _, label := range info.Labels {
// 			stringSlice := strings.SplitN(label, "=", 2)
// 			if len(stringSlice) > 1 {
// 				// If there is a conflict we will throw out a warning
// 				if v, ok := labelMap[stringSlice[0]]; ok && v != stringSlice[1] {
// 					fmt.Fprintln(storageosCli.Err(), "WARNING: labels with duplicate keys and conflicting values have been deprecated")
// 					break
// 				}
// 				labelMap[stringSlice[0]] = stringSlice[1]
// 			}
// 		}
// 	}
//
// 	fmt.Fprintf(storageosCli.Out(), "Experimental: %v\n", info.ExperimentalBuild)
// 	if info.ClusterStore != "" {
// 		fmt.Fprintf(storageosCli.Out(), "Cluster Store: %s\n", info.ClusterStore)
// 	}
//
// 	if info.ClusterAdvertise != "" {
// 		fmt.Fprintf(storageosCli.Out(), "Cluster Advertise: %s\n", info.ClusterAdvertise)
// 	}
//
// 	if info.RegistryConfig != nil && (len(info.RegistryConfig.InsecureRegistryCIDRs) > 0 || len(info.RegistryConfig.IndexConfigs) > 0) {
// 		fmt.Fprintln(storageosCli.Out(), "Insecure Registries:")
// 		for _, registry := range info.RegistryConfig.IndexConfigs {
// 			if registry.Secure == false {
// 				fmt.Fprintf(storageosCli.Out(), " %s\n", registry.Name)
// 			}
// 		}
//
// 		for _, registry := range info.RegistryConfig.InsecureRegistryCIDRs {
// 			mask, _ := registry.Mask.Size()
// 			fmt.Fprintf(storageosCli.Out(), " %s/%d\n", registry.IP.String(), mask)
// 		}
// 	}
//
// 	if info.RegistryConfig != nil && len(info.RegistryConfig.Mirrors) > 0 {
// 		fmt.Fprintln(storageosCli.Out(), "Registry Mirrors:")
// 		for _, mirror := range info.RegistryConfig.Mirrors {
// 			fmt.Fprintf(storageosCli.Out(), " %s\n", mirror)
// 		}
// 	}
//
// 	fmt.Fprintf(storageosCli.Out(), "Live Restore Enabled: %v\n", info.LiveRestoreEnabled)
//
// 	return nil
// }
//
// func formatInfo(storageosCli *command.StorageOSCli, info types.Info, format string) error {
// 	tmpl, err := templates.Parse(format)
// 	if err != nil {
// 		return cli.StatusError{StatusCode: 64,
// 			Status: "Template parsing error: " + err.Error()}
// 	}
// 	err = tmpl.Execute(storageosCli.Out(), info)
// 	storageosCli.Out().Write([]byte{'\n'})
// 	return err
// }
