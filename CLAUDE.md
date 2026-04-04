# Hello Mail

Multi-tenant SaaS email platform (like Resend).

## Tech Stack
- **Backend**: Go 1.26+, chi router, pgx, sqlc, asynq, zerolog
- **Frontend**: Angular 19, Searlo Design System
- **MTA**: Postfix + Go milter integration
- **Database**: PostgreSQL (DigitalOcean managed)
- **Cache/Queue**: Valkey/Redis (DigitalOcean managed)
- **Storage**: DigitalOcean Spaces (S3-compatible)

## Project Structure
- `cmd/server/` - API server entry point
- `cmd/worker/` - Asynq background worker
- `cmd/migrate/` - Database migration runner
- `internal/` - Business logic (auth, email, domain, inbox, webhook, analytics, audit, billing, template, team, suppression, smtp)
- `internal/api/` - HTTP handlers, middleware, response helpers
- `internal/db/` - PostgreSQL connection, migrations, sqlc queries
- `internal/config/` - Environment-based configuration
- `internal/observability/` - Logging, metrics
- `internal/queue/` - Asynq task definitions
- `web/` - Angular frontend
- `sdks/` - Client SDKs (Go, Node.js, Laravel)
- `deploy/` - Nginx, systemd, Postfix configs

## Commands
- `make run` - Start API server
- `make worker` - Start background worker
- `make build` - Build binaries to bin/
- `make migrate-up` - Run DB migrations
- `make migrate-down` - Rollback one migration
- `make sqlc` - Regenerate sqlc code
- `make test` - Run tests with race detector
- `make lint` - Run golangci-lint
- `make tidy` - Tidy Go modules
- `make clean` - Remove build artifacts
- `make docker-build` - Build Docker image
- `make docker-up` - Start services via docker-compose
- `make docker-down` - Stop docker-compose services

## Key Patterns
- Auth context: `auth.OrgIDFromContext(ctx)`, `auth.UserIDFromContext(ctx)`
- Response envelope: `response.JSON(w, status, data)` produces `{"data": ...}`
- DB queries: sqlc-generated in `internal/db/sqlcdb/`
- Async jobs: asynq tasks defined in `internal/queue/`
- All tables scoped by `org_id` for multi-tenancy
- Config loaded from environment variables (see `.env.example`)

## API Routes
- `/health` - Health check
- `/metrics` - Prometheus metrics
- `/t/o/{id}` - Open tracking pixel
- `/t/c/{id}` - Click tracking redirect
- `/v1/auth/*` - Authentication (Google OAuth)
- `/v1/api-keys` - API key management
- `/v1/domains` - Domain management and verification
- `/v1/emails` - Email sending and listing
- `/v1/inbox/*` - Inbox threads, messages, labels
- `/v1/webhooks` - Webhook endpoints
- `/v1/templates` - Email templates
- `/v1/suppressions` - Suppression list
- `/v1/analytics/*` - Email analytics

## Environment Variables
See `.env.example` for all required configuration.
