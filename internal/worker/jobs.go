package worker

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

// TaskCreatedJob simulates sending an email when a task is assigned.
type TaskCreatedJob struct {
	TaskTitle  string
	AssignedTo uuid.UUID
}

// Process executes the job.
func (j *TaskCreatedJob) Process(ctx context.Context) error {
	// Simulate an expensive operation like generating an email template and sending via SMTP
	log.Printf("[Job] Starting email dispatch for new task: '%s'\n", j.TaskTitle)
	
	time.Sleep(2 * time.Second) // Simulate network delay
	
	if j.AssignedTo != uuid.Nil {
		log.Printf("[Job] Successfully sent 'Task Assigned' email to User ID: %s\n", j.AssignedTo)
	} else {
		log.Printf("[Job] Successfully sent 'Unassigned Task Created' notification to org admins\n")
	}

	return nil
}
