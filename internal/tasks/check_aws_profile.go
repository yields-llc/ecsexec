package tasks

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type (
	// CheckAWSProfile is a task to check aws profile.
	CheckAWSProfile struct {
		// No fields are needed here.
	}
)

// Run is a method to execute this task.
func (c *CheckAWSProfile) Run(ctx context.Context) TaskResult {
	ctxValue := GetContextValue(ctx)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return TaskResult{Error: err}
	}

	svc := sts.NewFromConfig(cfg)
	_, err = svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return TaskResult{Error: err}
	}

	ctxValue.AwsCfg = cfg

	return TaskResult{}
}

// NewCheckAWSProfile is creation of this task.
func NewCheckAWSProfile() *CheckAWSProfile {
	return &CheckAWSProfile{}
}
