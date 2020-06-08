package update

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type errInvalidSizeArg struct {
	got string
}

func (e *errInvalidSizeArg) Error() string {
	return fmt.Sprintf("invalid size argument, got %s, use the following notations: '42 MiB', '42mib', or '44040192'", e.got)
}

func newErrInvalidSizeArg(got string) *errInvalidSizeArg {
	return &errInvalidSizeArg{
		got: got,
	}
}

type errInvalidArgNum struct {
	got  []string
	want int
}

func (e *errInvalidArgNum) Error() string {
	return fmt.Sprintf("invalid number of arguments, got %v, expected %d", e.got, e.want)
}

func newErrInvalidArgNum(got []string, want int) *errInvalidArgNum {
	return &errInvalidArgNum{
		got:  got,
		want: want,
	}
}

type errInvalidReplicaNum struct {
	got uint64
}

func (e *errInvalidReplicaNum) Error() string {
	return fmt.Sprintf("invalid number of replicas specified, got %d, cannot be more than %d", e.got, 5)
}

func newErrInvalidReplicaNum(got uint64) *errInvalidReplicaNum {
	return &errInvalidReplicaNum{
		got: got,
	}
}

// newVolumeUpdate configures the set of commands which are grouped by the
// "volume" noun.
func newVolumeUpdate(client Client, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "volume",
		Short: "Make changes to an existing volume",
	}

	command.AddCommand(
		newVolumeReplicas(os.Stdout, client, config),
		newVolumeDescription(os.Stdout, client, config),
		newVolumeLabels(os.Stdout, client, config),
		newVolumeSize(os.Stdout, client, config),
	)

	return command
}
