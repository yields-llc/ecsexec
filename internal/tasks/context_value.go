package tasks

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type (
	// ContextValue is struct for shared data in execution task.
	ContextValue struct {
		ClusterName string
		ServiceName string
		Container   string
		Command     string

		AwsCfg aws.Config

		ServiceAlreadyStarted bool

		TaskARNs []string
	}

	contextKey string
)

const (
	contextValueKey contextKey = "github.com/yields-llc/ecsexec:context-value"
)

// ContextWithValue is a method to set context value.
func ContextWithValue(ctx context.Context, ctxValue *ContextValue) context.Context {
	return context.WithValue(ctx, contextValueKey, ctxValue)
}

// GetContextValue is a method to get context value.
func GetContextValue(ctx context.Context) *ContextValue {
	ctxValue, ok := ctx.Value(contextValueKey).(*ContextValue)
	if !ok {
		panic("failed to get ctx.Value(\"value\") . please set context value with context.WithValue")
	}

	return ctxValue
}
