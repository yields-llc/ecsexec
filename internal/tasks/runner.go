package tasks

import (
	"context"
	"fmt"

	"github.com/google/subcommands"
)

type (
	// Runnable is a interface to execute multiple tasks.
	Runnable interface {
		Run(ctx context.Context) subcommands.ExitStatus
	}

	// Runner is a struct to execute multiple tasks.
	Runner struct {
		tasks []Task
	}
)

// Run is a method to execute multiple tasks in sequence.
func (r *Runner) Run(ctx context.Context) subcommands.ExitStatus {
	for _, t := range r.tasks {
		result := t.Run(ctx)
		if result.Error != nil {
			return r.handleError(result.Error)
		}

		fmt.Println("")
	}

	return subcommands.ExitSuccess
}

func (r *Runner) handleError(err error) subcommands.ExitStatus {
	fmt.Printf("[error]: \n%v\n", err)

	return subcommands.ExitFailure
}

// NewRunner is creation of this runner.
func NewRunner(tasks []Task) *Runner {
	return &Runner{tasks: tasks}
}
