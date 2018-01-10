package mount

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func runFile(ctx context.Context, args ...string) (string, error) {
	const file = "/usr/bin/file"
	return runCmd(ctx, file, args...)
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
