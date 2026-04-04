package queue

// Task type constants for asynq job queue.
// Each task type maps to a handler registered in the worker.
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
