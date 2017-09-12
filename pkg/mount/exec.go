package mount

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// OS utilities must be in path, which shouldn't be a problem as they need to
// run as root.
const (
	mount     = "/bin/mount"
	umount    = "/bin/umount"
	file      = "/usr/bin/file"
	mkfsExt2  = "/sbin/mkfs.ext2"
	mkfsExt3  = "/sbin/mkfs.ext3"
	mkfsExt4  = "/sbin/mkfs.ext4"
	mkfsXfs   = "/sbin/mkfs.xfs"
	mkfsBtrfs = "/bin/mkfs.btrfs"
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

func runCmd(ctx context.Context, cmd string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, cmd, args...)
	out, err := command.Output()

	if ctx.Err() == context.DeadlineExceeded {
		if log.StandardLogger().Level == log.DebugLevel {
			log.WithFields(log.Fields{
				"cmd":   cmd,
				"args":  args,
				"error": err,
			}).Error("fail to execute command, timeout exceeded")
		}
		return "", fmt.Errorf("timeout exceeded")
	}

	if err != nil {
		if log.StandardLogger().Level == log.DebugLevel {
			log.WithFields(log.Fields{
				"cmd":   cmd,
				"args":  args,
				"error": err,
			}).Error("fail to get output from command")
		}
		return strings.TrimSpace(string(out)), err
	}
	return strings.TrimSpace(string(out)), nil
}
