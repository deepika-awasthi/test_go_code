package opentelemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/workflow"
	// "go.temporal.io/sdk/activity"
	"github.com/google/uuid"
)

var tracer trace.Tracer

func init() {
	// Initialize the tracer with a name
	tracer = otel.Tracer("github.com/temporalio/samples-go/otel")
}

type res struct {
	Key   string
	Value string
}

// SampleScheduleWorkflow is the Temporal workflow
func SampleScheduleWorkflow(ctx workflow.Context, name string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Schedule workflow started.", "StartTime", workflow.Now(ctx))

	// Set up activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Second, // short timeout to fail quickly
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	tracking_id := uuid.New().String()

	// Add a span to the existing trace using the workflow_id
	ctx, span := startWorkflowSpan(ctx, "SampleScheduleWorkflow", tracking_id)
	defer span.End()

	var result interface{}
	if err := workflow.ExecuteActivity(ctx, DoSomething, tracking_id).Get(ctx, &result); err != nil {
		logger.Error("Workflow failed.", "Error", err)
		return err
	}

	logger.Info("Workflow completed.")
	return nil
}

// startWorkflowSpan starts a span for the workflow within the existing trace
func startWorkflowSpan(ctx workflow.Context, spanName string, trackingID string) (workflow.Context, trace.Span) {
	_, span := tracer.Start(context.Background(), spanName)
	span.SetAttributes(
		attribute.String("tracking", "check_" + trackingID),
	)
	return ctx, span
}

//activity
func DoSomething(ctx context.Context, trackingID string) (*res, error) {
	// Start a new span for the activity
	ctx, span := tracer.Start(ctx, "DoSomething")
	defer span.End()


	span.SetAttributes(
			attribute.String("tracking", "check_" + trackingID), // Add workflow_id to the activity span
		)
	return nil, nil
}
