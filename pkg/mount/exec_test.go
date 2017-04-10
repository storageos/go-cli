package mount

import (
	"context"
	"time"

	"testing"
)

func TestRunCmdDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := runCmd(ctx, "ping", "-c 4", "-i 1", "8.8.8.8")

	// since ping will not stop - we should get an error
	if err.Error() != "timeout exceeded" {
		t.Errorf("expected error to be 'timeout exceeded' but got %s", err)
	}
}
