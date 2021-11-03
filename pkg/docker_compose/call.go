package docker_compose

import (
	"fmt"
	"os/exec"
	"github.com/Seb-C/watch-dock/pkg/context"
)

func CallOnce(ctx context.Context, args ...string) ([]byte, error) {
	allArgs := make([]string, 0, len(args)+len(ctx.DockerComposeArgs))
	allArgs = append(allArgs, ctx.DockerComposeArgs...)
	allArgs = append(allArgs, args...)

	cmd := exec.CommandContext(ctx, "docker-compose", allArgs...)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, isExitErr := err.(*exec.ExitError); isExitErr {
			return nil, fmt.Errorf("Error calling docker-compose: %w\n%v", exitErr, string(exitErr.Stderr))
		} else {
			return nil, err
		}
	}

	return output, nil
}
