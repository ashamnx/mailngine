package queue

import "github.com/google/uuid"

// Task type constants for asynq job queue.
// Each task type maps to a handler registered in the worker.
//
// DEPLOY-SAFETY RULES — all task handlers MUST follow these:
//
//  1. No single task execution should exceed 30s in the happy path.
//     The worker's ShutdownTimeout (default 60s) gives in-flight tasks
//     time to finish during deploys. Tasks consistently exceeding this
//     will be force-killed.
//
//  2. Long operations (bulk sends, large exports) must be chunked with
//     a resumable cursor persisted to the database. If killed mid-work,
//     the task retries and resumes from the last committed cursor.
//
//  3. Each chunk must be committed to DB before proceeding to the next.
//     This ensures crash-safety — no work is lost on process termination.
//
//  4. The Valkey queue is durable. Only in-flight tasks (max Concurrency)
//     are affected by a deploy restart. Queued tasks are safe.
const (
	// Email tasks
	TaskSendEmail      = "email:send"
	TaskSendBatch      = "email:send_batch"
	TaskSendScheduled  = "email:send_scheduled"

	// Domain tasks
	TaskVerifyDNS      = "domain:verify_dns"
	TaskAutoCreateDNS  = "domain:auto_create_dns"

	// Webhook tasks
	TaskDispatchWebhook = "webhook:dispatch"

	// Analytics tasks
	TaskAggregateDaily  = "analytics:aggregate_daily"
	TaskAggregateMonthly = "analytics:aggregate_monthly"

	// Bounce/FBL tasks
	TaskProcessBounce   = "smtp:process_bounce"
	TaskProcessFBL      = "smtp:process_fbl"

	// IP warmup tasks
	TaskAdvanceWarmup   = "ip:advance_warmup"

	// Team tasks
	TaskSendInvite      = "team:send_invite"
)

// BatchSendPayload is the payload for the email:send_batch task.
//
// The handler MUST be deploy-safe: it chunks recipients into individual
// email:send tasks rather than processing all recipients inline.
//
// Flow:
//  1. Load batch record from DB (contains recipient list or query).
//  2. For each chunk of recipients (e.g., 100 at a time):
//     a. Create individual email records in DB.
//     b. Enqueue individual email:send tasks.
//     c. Update batch progress cursor in DB (commit before next chunk).
//  3. Mark batch as "enqueued".
//
// This ensures:
//   - Each individual email:send task takes <1s (deploy-safe).
//   - The batch task itself runs in bounded chunks (deploy-safe).
//   - If killed mid-batch, it resumes from the cursor (idempotent).
//   - Progress is trackable via the batch record.
type BatchSendPayload struct {
	BatchID uuid.UUID `json:"batch_id"`
}
