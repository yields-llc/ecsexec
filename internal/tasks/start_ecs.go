package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/yields-llc/ecsexec/internal/aws/clients"
)

type (
	// StartECS is a task to start ECS service.
	StartECS struct {
		// No fields are needed here.
	}
)

const (
	defaultWaitDurationForServicesStable = 3 * time.Minute
)

// Run is a method to execute this task.
func (s *StartECS) Run(ctx context.Context) TaskResult {
	ctxValue := GetContextValue(ctx)
	svc := ecs.NewFromConfig(ctxValue.AwsCfg)
	ecsClient := clients.NewECS(svc)

	service, err := ecsClient.DescribeService(ctx, ctxValue.ClusterName, ctxValue.ServiceName)
	if err != nil {
		return TaskResult{Error: err}
	}

	if service.RunningCount == 0 {
		fmt.Printf("%q service was stopped. it will be started now...\n", ctxValue.ServiceName)
		service, err = ecsClient.StartService(ctx, ctxValue.ClusterName, ctxValue.ServiceName, 1)
		if err != nil {
			return TaskResult{Error: err}
		}
		err = ecsClient.WaitForServicesStable(
			ctx,
			ctxValue.ClusterName,
			ctxValue.ServiceName,
			defaultWaitDurationForServicesStable,
		)
		if err != nil {
			return TaskResult{Error: err}
		}

		fmt.Printf("%q service has been started.\n", ctxValue.ServiceName)
	} else {
		fmt.Printf("%q service was already started (running count = %d).\n", ctxValue.ServiceName, service.RunningCount)
		ctxValue.ServiceAlreadyStarted = true
	}

	taskARNs, err := ecsClient.ListTasks(ctx, ctxValue.ClusterName, ctxValue.ServiceName)
	if err != nil {
		return TaskResult{Error: err}
	}

	ctxValue.TaskARNs = taskARNs

	if service.RunningCount == 0 {
		fmt.Println("waiting for SSM agent running...")

		err = ecsClient.WaitForSSMAgentRunning(
			ctx,
			ctxValue.ClusterName,
			ctxValue.TaskARNs[0],
		)
		if err != nil {
			return TaskResult{Error: err}
		}

		fmt.Println("SSM agent is running")
	}

	return TaskResult{}
}

// NewStartECS is creation of this task.
func NewStartECS() *StartECS {
	return &StartECS{}
}
