package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type (
	// ECS has ECS client and wrapped methods.
	ECS struct {
		client *ecs.Client
	}
)

// DescribeService is the method for describe ECS service.
func (e *ECS) DescribeService(ctx context.Context, clusterName string, serviceName string) (*types.Service, error) {
	services, err := e.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &clusterName,
		Services: []string{serviceName},
	})

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if len(services.Services) == 0 {
		return nil, fmt.Errorf("%s is not found", serviceName)
	}

	return &services.Services[0], nil
}

// StartService is the method for starting ECS service.
func (e *ECS) StartService(ctx context.Context, clusterName, serviceName string, taskCount int32) (*types.Service, error) {
	result, err := e.client.UpdateService(ctx, &ecs.UpdateServiceInput{
		Cluster:      &clusterName,
		Service:      &serviceName,
		DesiredCount: &taskCount,
	})

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result.Service, nil
}

// StopService is the method for stopping ECS service.
func (e *ECS) StopService(ctx context.Context, clusterName, serviceName string) (*types.Service, error) {
	taskCount := int32(0)
	result, err := e.client.UpdateService(ctx, &ecs.UpdateServiceInput{
		Cluster:      &clusterName,
		Service:      &serviceName,
		DesiredCount: &taskCount,
	})

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result.Service, nil
}

// WaitForServicesStable is the method for waiting until ECS service status is stable.
func (e *ECS) WaitForServicesStable(ctx context.Context, clusterName string, serviceName string, maxWait time.Duration) error {
	return ecs.NewServicesStableWaiter(e.client).Wait(
		ctx,
		&ecs.DescribeServicesInput{
			Cluster:  &clusterName,
			Services: []string{serviceName},
		},
		maxWait,
	)
}

// WaitForTasksStopped is the method for waiting until ECS tasks is stopped.
func (e *ECS) WaitForTasksStopped(ctx context.Context, clusterName string, taskARN string, maxWait time.Duration) error {
	return ecs.NewTasksStoppedWaiter(e.client).Wait(
		ctx,
		&ecs.DescribeTasksInput{
			Cluster: &clusterName,
			Tasks:   []string{taskARN},
		},
		maxWait,
	)
}

// WaitForSSMAgentRunning is the method for waiting until SSM agent is running.
func (e *ECS) WaitForSSMAgentRunning(ctx context.Context, clusterName string, taskARN string) error {
	timeout := 1 * time.Minute
	sleep := 1 * time.Second
	couldNotConfirmErr := fmt.Errorf("could not confirm that SSM agent is running")

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		task, err := e.DescribeTask(ctx, clusterName, taskARN)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		if len(task.Containers) == 0 {
			return couldNotConfirmErr
		}
		container := task.Containers[0]
		if len(container.ManagedAgents) == 0 {
			return couldNotConfirmErr
		}
		managedAgent := container.ManagedAgents[0]

		if *managedAgent.LastStatus == "RUNNING" {
			return nil
		}
		time.Sleep(sleep)
	}

	return couldNotConfirmErr
}

// ListTasks is the method for listing tasks of service.
func (e *ECS) ListTasks(ctx context.Context, clusterName string, serviceName string) ([]string, error) {
	tasks, err := e.client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:     &clusterName,
		ServiceName: &serviceName,
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if len(tasks.TaskArns) == 0 {
		return nil, fmt.Errorf("taskARNs not found")
	}

	return tasks.TaskArns, nil
}

// DescribeTask is the method for describing task.
func (e *ECS) DescribeTask(ctx context.Context, clusterName string, taskARN string) (*types.Task, error) {
	tasks, err := e.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &clusterName,
		Tasks:   []string{taskARN},
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if len(tasks.Tasks) == 0 {
		return nil, fmt.Errorf("tasks not found")
	}

	return &tasks.Tasks[0], nil
}

// ExecuteCommand is the method for executing command.
func (e *ECS) ExecuteCommand(ctx context.Context, clusterName string, taskID string, container string, command string) (*types.Session, error) {
	result, err := e.client.ExecuteCommand(ctx, &ecs.ExecuteCommandInput{
		Cluster:     &clusterName,
		Task:        &taskID,
		Container:   &container,
		Command:     &command,
		Interactive: true,
	})

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result.Session, nil
}

// NewECS is creation of this client.
func NewECS(client *ecs.Client) *ECS {
	return &ECS{client: client}
}
