package mount

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// OS utilities must be in path, which shouldn't be a problem as they need to
// run as root.
const (
	mount    = "/bin/mount"
	umount   = "/bin/umount"
	file     = "/usr/bin/file"
	mkfsExt4 = "/sbin/mkfs"
)

func runMount(ctx context.Context, args ...string) (string, error) {
	return runCmd(ctx, mount, args...)
}

func runUmount(ctx context.Context, args ...string) (string, error) {
	return runCmd(ctx, umount, args...)
}

func runFile(ctx context.Context, args ...string) (string, error) {
	return runCmd(ctx, file, args...)
}

func runMkfsExt4(ctx context.Context, args ...string) (string, error) {
	return runCmd(ctx, mkfsExt4, args...)
}

func runCmd(ctx context.Context, cmd string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, cmd, args...)
	out, err := command.Output()

	if ctx.Err() == context.DeadlineExceeded {
		log.WithFields(log.Fields{
			"cmd":   cmd,
			"args":  args,
			"error": err,
		}).Error("fail to execute command, timeout exceeded")
		return "", fmt.Errorf("timeout exceeded")
	}

	if err != nil {
		log.WithFields(log.Fields{
			"cmd":   cmd,
			"args":  args,
			"error": err,
		}).Error("fail to get output from command")
		return strings.TrimSpace(string(out)), err
	}
	return strings.TrimSpace(string(out)), nil
}
