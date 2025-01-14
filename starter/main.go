package main

import (
	"context"
	"log"
	"time"

	"github.com/pborman/uuid"
	otel "github.com/temporalio/samples-go/opentelemetry"
	// intercept "github.com/temporalio/samples-go/interceptor"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	// "go.temporal.io/sdk/workflow"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the OpenTelemetry tracer provider
	tp, err := otel.InitializeGlobalTracerProvider()
	if err != nil {
		log.Fatalln("Unable to create a global trace provider", err)
	}
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
		// Interceptors: []otel.ClientInterceptor{
		// 	otel.NewCustomOtelInterceptor(),
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

	// Create a schedule for the workflow
	scheduleID := "otel_schedule_" + uuid.New()
	// workflowID := "otel_workflow_" + uuid.New()
	uniqueTraceIdentifier := "trace_id_" + uuid.New()

	log.Println("Creating schedule", "ScheduleID", scheduleID)

	// Create a context with the unique propagation data
	// propagationValues := otel.Values{
	// 	Key:             "workflowKey",
	// 	Value:           uniqueTraceIdentifier,
	// 	// TraceIdentifier: uniqueTraceIdentifier,
	// }
	// propagationCtx := context.WithValue(ctx, otel.PropagateKey, propagationValues)


	// // Create the schedule without limiting actions
	// scheduleHandle, err := c.ScheduleClient().Create(propagationCtx, client.ScheduleOptions{
	// 	ID: scheduleID,
	// 	Spec: client.ScheduleSpec{
	// 		Intervals: []client.ScheduleIntervalSpec{
	// 			{
	// 				Every: 20 * time.Second, // Run every 2 minutes
	// 			},
	// 		},
	// 	},
	// 	Action: &client.ScheduleWorkflowAction{
	// 		ID:        "dynamic_workflow_" + uuid.New(),
	// 		Workflow:  otel.SampleScheduleWorkflow,
	// 		Args:      []interface{}{"Temporal"}, // Pass arguments to the workflow
	// 		TaskQueue: "otel",
	// 	},
	// })
	// if err != nil {
	// 	log.Fatalln("Unable to create schedule", err)
	// }

	scheduleHandle, err := c.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			Intervals: []client.ScheduleIntervalSpec{
				{
					Every: 20 * time.Second, // Run every 20 seconds
				},
			},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        uniqueTraceIdentifier, // Generate a unique Workflow ID
			Workflow:  otel.SampleScheduleWorkflow,
			Args:      []interface{}{"Temporal"},  
			TaskQueue: "otel",
		},
	})

	log.Println("Created schedule successfully", "ScheduleID", scheduleID, "workflowKey", uniqueTraceIdentifier)

	// Update the schedule to limit actions to 5
	err = scheduleHandle.Update(ctx, client.ScheduleUpdateOptions{
		DoUpdate: func(schedule client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
			// Limit the number of actions to 5
			schedule.Description.Schedule.State.LimitedActions = true
			schedule.Description.Schedule.State.RemainingActions = 4
			return &client.ScheduleUpdate{
				Schedule: &schedule.Description.Schedule,
			}, nil
		},
	})
	if err != nil {
		log.Fatalln("Unable to update schedule", err)
	}

	log.Println("Updated schedule to limit actions to 5", "ScheduleID", scheduleID)

	// Monitor the schedule
	for {
		// Check if the schedule is completing actions
		time.Sleep(10 * time.Second)
		description, err := scheduleHandle.Describe(ctx)
		if err != nil {
			log.Fatalln("Unable to describe schedule", err)
		}

		log.Println("Schedule Status:", "RemainingActions", description.Schedule.State.RemainingActions)
		if description.Schedule.State.RemainingActions == 0 {
			log.Println("Schedule completed all actions")
			break
		}
	}
}
