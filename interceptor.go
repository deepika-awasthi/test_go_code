package opentelemetry

// // import (
// // 	"context"

// // 	"go.opentelemetry.io/otel/trace"
// // 	"go.temporal.io/sdk/converter"
// // 	"go.temporal.io/sdk/interceptor"
// // 	"go.temporal.io/sdk/log"
// // 	"go.temporal.io/sdk/workflow"
// // )


// // type workerInterceptor struct {
// // 	interceptor.WorkerInterceptorBase
// // 	options InterceptorOptions
// // }

// // type InterceptorOptions struct {
// // 	GetExtraLogTagsForWorkflow func(workflow.Context) []interface{}
// // 	GetExtraLogTagsForActivity func(context.Context) []interface{}
// // }

// // func NewWorkerInterceptor(options InterceptorOptions) interceptor.WorkerInterceptor {
// // 	return &workerInterceptor{options: options}
// // }

// // func (w *workerInterceptor) InterceptActivity(
// // 	ctx context.Context,
// // 	next interceptor.ActivityInboundInterceptor,
// // ) interceptor.ActivityInboundInterceptor {
// // 	i := &activityInboundInterceptor{root: w}
// // 	i.Next = next
// // 	return i
// // }

// // type activityInboundInterceptor struct {
// // 	interceptor.ActivityInboundInterceptorBase
// // 	root *workerInterceptor
// // }

// // func (a *activityInboundInterceptor) Init(outbound interceptor.ActivityOutboundInterceptor) error {
// // 	i := &activityOutboundInterceptor{root: a.root}
// // 	i.Next = outbound
// // 	return a.Next.Init(i)
// // }

// // type activityOutboundInterceptor struct {
// // 	interceptor.ActivityOutboundInterceptorBase
// // 	root *workerInterceptor
// // }

// // func (a *activityOutboundInterceptor) GetLogger(ctx context.Context) log.Logger {
// // 	logger := a.Next.GetLogger(ctx)
// // 	// Add extra tags if any
// // 	if a.root.options.GetExtraLogTagsForActivity != nil {
// // 		if extraTags := a.root.options.GetExtraLogTagsForActivity(ctx); len(extraTags) > 0 {
// // 			logger = log.With(logger, extraTags...)
// // 		}
// // 	}
// // 	return logger
// // }

// // func (w *workerInterceptor) InterceptWorkflow(
// // 	ctx workflow.Context,
// // 	next interceptor.WorkflowInboundInterceptor,
// // ) interceptor.WorkflowInboundInterceptor {
// // 	i := &workflowInboundInterceptor{root: w}
// // 	i.Next = next
// // 	return i
// // }

// // type workflowInboundInterceptor struct {
// // 	interceptor.WorkflowInboundInterceptorBase
// // 	root *workerInterceptor
// // }

// // func (w *workflowInboundInterceptor) Init(outbound interceptor.WorkflowOutboundInterceptor) error {
// // 	i := &workflowOutboundInterceptor{root: w.root}
// // 	i.Next = outbound
// // 	return w.Next.Init(i)
// // }

// // func (w *workflowInboundInterceptor) HandleSignal(ctx workflow.Context, input *interceptor.HandleSignalInput) error {
// // 	// Extract the signal payload from the input.Arg
// // 	var signalPayload SignalPayload
// // 	err := converter.GetDefaultDataConverter().FromPayloads(input.Arg, &signalPayload)
// // 	if err != nil {
// // 		w.root.options.GetExtraLogTagsForWorkflow = func(workflowCtx workflow.Context) []interface{} {
// // 			return []interface{}{"error", "failed to deserialize signal payload"}
// // 		}
// // 		return err
// // 	}

// // 	// Extract trace context from the signal payload
// // 	traceID, _ := trace.TraceIDFromHex(signalPayload.TraceID)
// // 	spanID, _ := trace.SpanIDFromHex(signalPayload.SpanID)

// // 	// Continue the trace with the provided context
// // 	traceCtx := propagateTraceContext(traceID, spanID)

// // 	w.root.options.GetExtraLogTagsForWorkflow = func(workflowCtx workflow.Context) []interface{} {
// // 		return []interface{}{"trace_id", traceID.String()}
// // 	}

// // 	ctx = workflow.WithValue(ctx, "traceCtx", traceCtx)
// // 	return w.Next.HandleSignal(ctx, input)
// // }

// // type workflowOutboundInterceptor struct {
// // 	interceptor.WorkflowOutboundInterceptorBase
// // 	root *workerInterceptor
// // }

// // func (w *workflowOutboundInterceptor) GetLogger(ctx workflow.Context) log.Logger {
// // 	logger := w.Next.GetLogger(ctx)

// // 	if w.root.options.GetExtraLogTagsForWorkflow != nil {
// // 		if extraTags := w.root.options.GetExtraLogTagsForWorkflow(ctx); len(extraTags) > 0 {
// // 			logger = log.With(logger, extraTags...)
// // 		}
// // 	}
// // 	return logger
// // }

// // func propagateTraceContext(traceID trace.TraceID, spanID trace.SpanID) context.Context {
// // 	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
// // 		TraceID:    traceID,
// // 		SpanID:     spanID,
// // 		Remote:     true,
// // 		TraceFlags: trace.FlagsSampled,
// // 	})
// // 	return trace.ContextWithSpanContext(context.Background(), spanCtx)
// // }





// package opentelemetry

// import (
// 	"context"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/trace"
// 	"go.temporal.io/sdk/interceptor"
// 	// "go.temporal.io/sdk/workflow"
// )

// // CustomOtelInterceptor adds custom trace attributes for workflows and activities.
// type CustomOtelInterceptor struct {
// 	interceptor.ClientInterceptorBase
// 	tracer trace.Tracer
// }

// // NewCustomOtelInterceptor creates a new instance of the interceptor.
// func NewCustomOtelInterceptor() interceptor.ClientInterceptor {
// 	return &CustomOtelInterceptor{
// 		tracer: otel.Tracer("github.com/temporalio/samples-go/otel"),
// 	}
// }

// // ExecuteWorkflow adds trace attributes before starting a workflow.
// func (c *CustomOtelInterceptor) ExecuteWorkflow(
// 	ctx context.Context,
// 	in *interceptor.ClientExecuteWorkflowInput,
// ) (interceptor.ClientWorkflowRun, error) {
// 	// Extract trace identifier from context.
// 	values, ok := ctx.Value(PropagateKey).(Values)
// 	if ok {
// 		// Add trace attributes.
// 		_, span := c.tracer.Start(ctx, "ExecuteWorkflow", trace.WithAttributes(
// 			//attribute.String("trace_identifier", values.TraceIdentifier),
// 			attribute.String("key", values.Key),
// 			attribute.String("value", values.Value),
// 		))
// 		defer span.End()
// 	}

// 	// Continue the workflow execution.
// 	return c.ClientInterceptorBase.ExecuteWorkflow(ctx, in)
// }

// // ExecuteActivity adds trace attributes before starting an activity.
// func (c *CustomOtelInterceptor) ExecuteActivity(
// 	ctx context.Context,
// 	in *interceptor.ClientExecuteActivityInput,
// ) error {
// 	// Extract trace identifier from context.
// 	values, ok := ctx.Value(PropagateKey).(Values)
// 	if ok {
// 		// Add trace attributes.
// 		_, span := c.tracer.Start(ctx, "ExecuteActivity", trace.WithAttributes(
// 			//attribute.String("trace_identifier", values.TraceIdentifier),
// 			attribute.String("key", values.Key),
// 			attribute.String("value", values.Value),
// 		))
// 		defer span.End()
// 	}

// 	// Continue the activity execution.
// 	return c.ClientInterceptorBase.ExecuteActivity(ctx, in)
// }
