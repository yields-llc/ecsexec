package tasks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
)

type (
	// ExecuteCommand is a task to execute command.
	ExecuteCommand struct {
		// No fields are needed here.
	}
)

// Run is a method to execute this task.
func (s *ExecuteCommand) Run(ctx context.Context) TaskResult {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	ctxValue := GetContextValue(ctx)
	command := fmt.Sprintf(
		`$(which aws) --profile %s ecs execute-command --cluster %s --task %s --container %s --interactive --command "%s"`,
		os.Getenv("AWS_PROFILE"),
		ctxValue.ClusterName,
		ctxValue.TaskARNs[0],
		ctxValue.Container,
		ctxValue.Command,
	)
	cmd := exec.CommandContext(
		ctx, "bash", "-c", command,
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	err := cmd.Run()
	if errors.Is(err, context.Canceled) {
		return TaskResult{}
	} else if err != nil {
		return TaskResult{Error: err}
	}

	return TaskResult{}
}

// NewExecuteCommand is creation of this task.
func NewExecuteCommand() *ExecuteCommand {
	return &ExecuteCommand{}
}
