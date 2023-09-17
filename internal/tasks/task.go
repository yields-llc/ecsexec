package tasks

import (
	"context"
)

type (
	// TaskResult is a struct of task execution result.
	TaskResult struct {
		Error error
	}

	// Task declare interface of execution task.
	Task interface {
		Run(ctx context.Context) TaskResult
	}
)
