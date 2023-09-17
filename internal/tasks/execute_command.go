package tasks

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type (
	// ExecuteCommand is a task to execute command.
	ExecuteCommand struct {
		// No fields are needed here.
	}
)

// Run is a method to execute this task.
func (s *ExecuteCommand) Run(ctx context.Context) TaskResult {
	ctxValue := GetContextValue(ctx)
	command := fmt.Sprintf(
		`$(which aws) --profile %s ecs execute-command --cluster %s --task %s --container %s --interactive --command "%s"`,
		os.Getenv("AWS_PROFILE"),
		ctxValue.ClusterName,
		ctxValue.TaskARNs[0],
		ctxValue.Container,
		ctxValue.Command,
	)
	cmd := exec.Command(
		"bash", "-c", command,
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return TaskResult{Error: err}
	}

	return TaskResult{}
}

// NewExecuteCommand is creation of this task.
func NewExecuteCommand() *ExecuteCommand {
	return &ExecuteCommand{}
}
