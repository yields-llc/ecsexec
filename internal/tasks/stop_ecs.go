package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/yields-llc/ecsexec/internal/aws/clients"
)

type (
	// StopECS is a task to start ECS service.
	StopECS struct {
		// No fields are needed here.
	}
)

const (
	defaultWaitDurationForTasksStopped = 3 * time.Minute
)

// Run is a method to execute this task.
func (s *StopECS) Run(ctx context.Context) TaskResult {
	ctxValue := GetContextValue(ctx)
	if ctxValue.ServiceAlreadyStarted {
		return TaskResult{}
	}

	svc := ecs.NewFromConfig(ctxValue.AwsCfg)
	ecsClient := clients.NewECS(svc)

	fmt.Printf("%q service was started. it will be stopped now...\n", ctxValue.ServiceName)
	_, err := ecsClient.StopService(ctx, ctxValue.ClusterName, ctxValue.ServiceName)
	if err != nil {
		return TaskResult{Error: err}
	}
	err = ecsClient.WaitForTasksStopped(
		ctx,
		ctxValue.ClusterName,
		ctxValue.TaskARNs[0],
		defaultWaitDurationForTasksStopped,
	)
	if err != nil {
		return TaskResult{Error: err}
	}
	fmt.Printf("%q service has been stopped.\n", ctxValue.ServiceName)

	return TaskResult{}
}

// NewStopECS is creation of this task.
func NewStopECS() *StopECS {
	return &StopECS{}
}
