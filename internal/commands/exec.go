package commands

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"github.com/yields-llc/ecsexec/internal/tasks"
)

type (
	// ExecCommand is the command to start ECS service.
	ExecCommand struct {
		profile       string
		clusterName   string
		serviceName   string
		containerName string
		command       string
	}
)

const (
	awsProfileEnvName = "AWS_PROFILE"
)

// Name is a name of this command.
func (p *ExecCommand) Name() string {
	return "exec"
}

// Synopsis is a synopsis of this command.
func (p *ExecCommand) Synopsis() string {
	return "this is a command to execute command by ECS exec"
}

// Usage is a usage of this command.
func (p *ExecCommand) Usage() string {
	return `exec -profile <AWS_PROFILE> \
    -cluster <ECS cluster name> \
    -service <ECS service name> \
    -container <container name> \
    -command <command>`
}

// SetFlags is definitions of this command arguments.
func (p *ExecCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.profile, "profile", "", "AWS profile")
	f.StringVar(&p.clusterName, "cluster", "", "ECS cluster name")
	f.StringVar(&p.serviceName, "service", "", "ECS service name")
	f.StringVar(&p.containerName, "container", "", "Container name")
	f.StringVar(&p.command, "command", "", "command")
}

// Execute is an implementations for this command.
func (p *ExecCommand) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if !p.validateFlags() {
		return subcommands.ExitFailure
	}

	p.setAWSProfile()

	ctx = tasks.ContextWithValue(ctx, &tasks.ContextValue{
		ClusterName: p.clusterName,
		ServiceName: p.serviceName,
		Container:   p.containerName,
		Command:     p.command,
	})

	runner := tasks.NewRunner([]tasks.Task{
		tasks.NewCheckAWSProfile(),
		tasks.NewStartECS(),
		tasks.NewExecuteCommand(),
		tasks.NewStopECS(),
	})

	exitCode := runner.Run(ctx)

	return exitCode
}

func (p *ExecCommand) validateFlags() bool {
	if p.clusterName == "" {
		fmt.Println("-cluster should not be empty")
		return false
	}
	if p.serviceName == "" {
		fmt.Println("-service should not be empty")
		return false
	}
	if p.containerName == "" {
		fmt.Println("-container should not be empty")
		return false
	}
	if p.command == "" {
		fmt.Println("-command should not be empty")
		return false
	}
	return true
}

func (p *ExecCommand) setAWSProfile() {
	if p.profile != "" {
		_ = os.Setenv(awsProfileEnvName, p.profile)
	}
}

// NewStartCommand is creation of this command.
func NewStartCommand() *ExecCommand {
	return &ExecCommand{}
}
