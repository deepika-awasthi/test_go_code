package opentelemetry

import (
	"context"

	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

type (
	contextKey struct{}
	propagator struct{}
	Values struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

// PropagateKey is the key used to store the value in the Context object
var PropagateKey = contextKey{}

// propagationKey is the key used by the propagator to pass values through the Temporal server headers
const propagationKey = "custom-header"

// DefaultKey and DefaultValue are used when no explicit value is set in the context
const (
	DefaultKey   = "myworkflow"
	DefaultValue = "default-value"
)

// NewContextPropagator initializes a new context propagator
func NewContextPropagator() workflow.ContextPropagator {
	return &propagator{}
}

func (s *propagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	// Fetch the value from the context or use the default value
	value := ctx.Value(PropagateKey)
	if value == nil {
		value = Values{
			Key:   DefaultKey,
			Value: DefaultValue,
		}
	}

	payload, err := converter.GetDefaultDataConverter().ToPayload(value)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)
	return nil
}

func (s *propagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	// Fetch the value from the workflow context or use the default value
	value := ctx.Value(PropagateKey) // Retrieve value using context key
	if value == nil {
		value = Values{
			Key:   DefaultKey,
			Value: DefaultValue,
		}
	}

	payload, err := converter.GetDefaultDataConverter().ToPayload(value)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)
	return nil
}

func (s *propagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	if value, ok := reader.Get(propagationKey); ok {
		var values Values
		if err := converter.GetDefaultDataConverter().FromPayload(value, &values); err != nil {
			return ctx, nil
		}
		ctx = context.WithValue(ctx, PropagateKey, values)
	}
	return ctx, nil
}

func (s *propagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	if value, ok := reader.Get(propagationKey); ok {
		var values Values
		if err := converter.GetDefaultDataConverter().FromPayload(value, &values); err != nil {
			return ctx, nil
		}
		// Attach the extracted value to the workflow context
		ctx = workflow.WithValue(ctx, PropagateKey, values)
	}
	return ctx, nil
}
