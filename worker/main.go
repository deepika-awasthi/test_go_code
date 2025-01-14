package main

import (
	"context"
	"log"

	otel "github.com/temporalio/samples-go/opentelemetry"
	// intercept "github.com/temporalio/samples-go/interceptor"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	// "go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
	// "go.temporal.io/sdk/workflow"
)

func main() {
	// Set up a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the OpenTelemetry tracer provider
	tp, err := otel.InitializeGlobalTracerProvider()
	if err != nil {
		log.Fatalln("Unable to create a global trace provider", err)
	}

	// Ensure tracer provider is shut down during cleanup
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Println("Error shutting down trace provider:", err)
		}
	}()

	// Set up tracing interceptor
	_, err = opentelemetry.NewTracingInterceptor(opentelemetry.TracerOptions{})
	if err != nil {
		log.Fatalln("Unable to create interceptor", err)
	}

	// Configure Temporal client options
	options := client.Options{
		// Interceptors: []interceptor.ClientInterceptor{
		// 	intercept.NewCustomOtelInterceptor(),
		// },
		// ContextPropagators: []workflow.ContextPropagator{
		// 	otel.NewContextPropagator(),
		// },
	}

	// Create Temporal client
	c, err := client.Dial(options)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// Create a worker for the "otel" task queue
	w := worker.New(c, "otel", worker.Options{
		// ContextPropagators: []workflow.ContextPropagator{
		// 	otel.NewContextPropagator(),
		// },
	})

	// Register workflows and activities
	w.RegisterWorkflow(otel.SampleScheduleWorkflow)
	w.RegisterActivity(otel.DoSomething)

	// Start the worker
	log.Println("Starting Temporal worker...")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Worker run failed", err)
	}
}
