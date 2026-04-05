# Mailngine

A multi-tenant SaaS email platform for developers. Send, receive, and manage email at any scale with powerful APIs and a beautiful dashboard.

**Live:** [mailngine.com](https://mailngine.com) | **Dashboard:** [app.mailngine.com](https://app.mailngine.com)

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.22+, chi router, pgx, sqlc, asynq, zerolog |
| **Dashboard** | Angular 19, Searlo Design System |
| **Website** | Qwik (SSG, zero-hydration) |
| **MTA** | Postfix + Go milter integration |
| **Database** | PostgreSQL (DigitalOcean managed) |
| **Cache/Queue** | Valkey/Redis (DigitalOcean managed) |
| **Storage** | DigitalOcean Spaces (S3-compatible) |
| **Auth** | Google OAuth + JWT + API keys |
| **SDKs** | Go, Node.js, Laravel |

## Features

- **Email Sending** — REST API + SMTP relay, batch sending, scheduling, idempotency
- **Domain Management** — DKIM key generation, SPF/DMARC/MX records, DNS verification, Cloudflare auto-DNS
- **Inbox** — Gmail-like 3-panel UI, JWZ threading, labels, search, star/archive/trash
- **Webhooks** — HMAC-SHA256 signed delivery, exponential backoff retry, delivery tracking
- **Analytics** — Open/click tracking, timeseries dashboards, event breakdown
- **Templates** — Variable substitution with preview, reusable email templates
- **Team Management** — Invite members, role-based access (owner/admin/member/viewer)
- **Billing** — Plan tiers (Free/Starter/Pro/Enterprise), usage tracking, overage support
- **Audit Logs** — Every mutation logged with user, IP, metadata
- **Suppression** — Auto-suppress on bounce/complaint, Valkey-cached lookups
- **Security** — RBAC, rate limiting, security headers, max body size, DKIM encryption

## Project Structure

```
mailngine/
  cmd/
    server/           Go API server entry point
    worker/           Asynq background worker
    migrate/          Database migration CLI
  internal/
    api/              HTTP handlers, middleware, response helpers
    auth/             Google OAuth, JWT, API keys, RBAC
    email/            MIME builder, DKIM signing, SMTP delivery
    domain/           Domain CRUD, DNS verification, Cloudflare
    inbox/            Threading, receiver, search
    webhook/          HMAC dispatch, retry, delivery tracking
    analytics/        Open/click tracking, aggregation
    audit/            Fire-and-forget audit logging
    billing/          Plan definitions, usage tracking
    template/         Variable rendering, preview
    team/             Org management, invitations
    suppression/      Bounce/FBL processing, Valkey cache
    smtp/             Bounce/FBL processors
    queue/            Asynq task definitions
    db/               PostgreSQL pool, migrations, sqlc queries
    cache/            Valkey connection
    observability/    zerolog, Prometheus metrics
  pkg/mailngine/      Shared API types and errors
  web/                Angular 19 dashboard
  website/            Qwik marketing website (14 pages)
  sdks/
    go/               Go SDK
    node/             Node.js/TypeScript SDK
    laravel/          Laravel/PHP SDK
  deploy/
    production/       Nginx, systemd, .env configs
    postfix/          Postfix config templates
    deploy.sh         Automated deployment script
```

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 20+
- PostgreSQL (or DigitalOcean managed)
- Valkey/Redis (or DigitalOcean managed)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- [sqlc](https://sqlc.dev/) CLI

### Setup

```bash
# Clone
git clone https://github.com/mailngine/mailngine.git
cd mailngine

# Configure
cp .env.example .env
# Edit .env with your database, Valkey, and Google OAuth credentials

# Run migrations
make migrate-up

# Install frontend dependencies
cd web && npm install && cd ..
cd website && npm install && cd ..

# Start everything
make run          # API server (port 8080)
make worker       # Background worker (separate terminal)
cd web && ng serve  # Dashboard (port 4200)
cd website && npm start  # Marketing site (port 5173)
```

### Available Commands

```bash
make build          # Build Go binaries to bin/
make run            # Start API server
make worker         # Start background worker
make test           # Run tests
make lint           # Run linter
make migrate-up     # Run database migrations
make migrate-down   # Rollback last migration
make sqlc           # Regenerate sqlc code
make tidy           # go mod tidy
make docker-build   # Build Docker image
make docker-up      # Start with docker-compose
make docker-down    # Stop docker-compose
```

## API

Base URL: `https://app.mailngine.com/v1`

Auth: `Authorization: Bearer mn_live_...` (API key) or `Bearer <jwt>` (session)

### Endpoints

| Resource | Endpoints |
|----------|-----------|
| **Auth** | `GET /auth/google`, `GET /auth/google/callback`, `POST /auth/logout`, `GET /auth/me` |
| **Emails** | `POST /emails`, `GET /emails`, `GET /emails/{id}` |
| **Domains** | `POST /domains`, `GET /domains`, `GET /domains/{id}`, `PATCH /domains/{id}`, `DELETE /domains/{id}`, `POST /domains/{id}/verify`, `POST /domains/{id}/auto-dns` |
| **Webhooks** | `POST /webhooks`, `GET /webhooks`, `GET /webhooks/{id}`, `PATCH /webhooks/{id}`, `DELETE /webhooks/{id}`, `GET /webhooks/{id}/deliveries` |
| **Templates** | `POST /templates`, `GET /templates`, `GET /templates/{id}`, `PATCH /templates/{id}`, `DELETE /templates/{id}`, `POST /templates/{id}/preview` |
| **API Keys** | `POST /api-keys`, `GET /api-keys`, `DELETE /api-keys/{id}` |
| **Inbox** | `GET /inbox/threads`, `GET /inbox/threads/{id}`, `DELETE /inbox/threads/{id}`, `GET /inbox/messages/{id}`, `PATCH /inbox/messages/{id}`, `DELETE /inbox/messages/{id}`, labels CRUD, `GET /inbox/search` |
| **Suppressions** | `GET /suppressions`, `POST /suppressions`, `DELETE /suppressions/{id}` |
| **Analytics** | `GET /analytics/overview`, `GET /analytics/timeseries`, `GET /analytics/events` |
| **Billing** | `GET /billing/usage`, `GET /billing/usage/history`, `GET /billing/plan` |
| **Org/Team** | `GET /org`, `PATCH /org`, `GET /org/members`, `POST /org/members/invite`, `PATCH /org/members/{id}`, `DELETE /org/members/{id}` |
| **Audit** | `GET /audit-logs`, `GET /audit-logs/{id}` |
| **Tracking** | `GET /t/o/{id}` (open pixel), `GET /t/c/{id}` (click redirect) |
| **Health** | `GET /health`, `GET /metrics` |

## SDKs

### Go

```go
client := mailngine.New("mn_live_...")

email, err := client.Emails.Send(ctx, &mailngine.SendEmailParams{
    From:    "hello@example.com",
    To:      []string{"user@example.com"},
    Subject: "Welcome!",
    HTML:    "<h1>Hello from Go</h1>",
})
```

```bash
go get github.com/mailngine/mailngine-go
```

### Node.js

```typescript
import { Mailngine } from 'mailngine';

const client = new Mailngine('mn_live_...');

const email = await client.emails.send({
    from: 'hello@example.com',
    to: ['user@example.com'],
    subject: 'Welcome!',
    html: '<h1>Hello from Node.js</h1>',
});
```

```bash
npm install mailngine
```

### Laravel

```php
use Mailngine\Mailngine;

$client = new Mailngine('mn_live_...');

$email = $client->emails()->send([
    'from' => 'hello@example.com',
    'to' => ['user@example.com'],
    'subject' => 'Welcome!',
    'html' => '<h1>Hello from Laravel</h1>',
]);
```

Laravel Mail integration:
```php
// config/mail.php
'mailers' => [
    'mailngine' => [
        'transport' => 'mailngine',
        'key' => env('MAILNGINE_API_KEY'),
    ],
],

// Usage
Mail::mailer('mailngine')->to($user)->send(new WelcomeEmail());
```

```bash
composer require mailngine/mailngine-laravel
```

## Deployment

### Production (DigitalOcean VPS)

The project is deployed to a DigitalOcean VPS at `178.128.208.168`:

| Service | Domain | Port |
|---------|--------|------|
| Marketing Website | mailngine.com | Nginx static |
| Dashboard + API | app.mailngine.com | API on 8090, Nginx proxy |
| Background Worker | — | Systemd service |

### Deploy Commands

```bash
./deploy/deploy.sh full       # Build + deploy everything
./deploy/deploy.sh api        # Rebuild + deploy Go binaries only
./deploy/deploy.sh dashboard  # Rebuild + deploy Angular only
./deploy/deploy.sh website    # Rebuild + deploy Qwik site only
./deploy/deploy.sh nginx      # Update nginx configs only
```

### Infrastructure

- **SSL**: Let's Encrypt via certbot (auto-renewal)
- **Process Manager**: systemd (`mailngine-api`, `mailngine-worker`)
- **Reverse Proxy**: Nginx with security headers, gzip, rate limiting
- **CDN**: Cloudflare (DNS proxy)
- **Database**: DigitalOcean Managed PostgreSQL
- **Cache/Queue**: DigitalOcean Managed Valkey
- **Storage**: DigitalOcean Spaces (S3-compatible)

### Manual Server Setup

```bash
# On the VPS
mkdir -p /opt/mailngine/bin /opt/mailngine/dashboard /opt/mailngine/website /var/www/certbot

# Install certbot
apt install certbot python3-certbot-nginx

# Get SSL certificates
certbot certonly --webroot -w /var/www/certbot -d mailngine.com
certbot certonly --webroot -w /var/www/certbot -d app.mailngine.com
```

## Database

19 tables with UUID primary keys, multi-tenant via `org_id` on all tenant-scoped tables:

`organizations`, `users`, `org_members`, `domains`, `dns_records`, `api_keys`, `emails`, `email_attachments`, `email_events`, `inbox_threads`, `inbox_messages`, `inbox_labels`, `inbox_message_labels`, `inbox_message_attachments`, `webhooks`, `webhook_deliveries`, `suppressions`, `templates`, `audit_logs`, `usage_daily`, `usage_monthly`, `ip_pools`

Migrations: `internal/db/migrations/`

## Environment Variables

See [.env.example](.env.example) for all required configuration variables.

Key variables:
- `DATABASE_URL` — PostgreSQL connection string
- `VALKEY_URL` — Valkey/Redis connection string
- `JWT_SECRET` — Secret for JWT signing
- `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` — Google OAuth credentials
- `FRONTEND_URL` — Dashboard URL for CORS and OAuth redirects

## Incomplete / In Progress

### Domain Connect (Auto-DNS Configuration)

The codebase includes a Domain Connect implementation (`internal/domain/domainconnect.go`) that enables one-click DNS configuration via supported providers (Cloudflare, GoDaddy, IONOS, WordPress.com, NameSilo, Vercel). The backend discovers providers via `_domainconnect` TXT records, generates signed redirect URLs, and the frontend shows an "Auto-Configure DNS" button when supported.

**What works now:**
- Provider discovery via DNS TXT lookup
- Signed redirect URL generation (HMAC-SHA256)
- Frontend dynamically shows auto-configure card per provider
- CSRF state token validation
- Auto-verification on callback return

**What's needed to activate:**
1. Submit a DNS template PR to [Domain-Connect/Templates](https://github.com/Domain-Connect/Templates)
2. Register as a Service Provider with Cloudflare (email: domain-connect@cloudflare.com)
3. Generate a signing keypair and publish the public key as a DNS TXT record at `domainconnect.mailngine.com`
4. Add `DOMAIN_CONNECT_PRIVATE_KEY` to `.env`

See: [domainconnect.org](https://www.domainconnect.org/) | [Spec](https://github.com/Domain-Connect/spec)

### Payment Gateway

Billing service has plan definitions and usage tracking but no payment processing. Payment gateway integration (Stripe, LemonSqueezy, etc.) is planned.

### Postfix Milter

The Go milter server interface is defined but not yet connected to a running Postfix instance. Config templates exist at `deploy/postfix/`.

### DKIM Signing

DKIM signing interface exists (`internal/email/dkim.go`) but uses a passthrough placeholder. Full signing requires integrating `emersion/go-msgauth`.

## Testing

```bash
make test   # Run all tests

# Tests cover:
# - JWT generation/validation (including alg:none attack prevention)
# - API key generation, hashing, uniqueness
# - MIME message builder (text, HTML, multipart, headers)
# - Webhook HMAC signing/verification
# - DKIM keypair generation and PEM parsing
# - HTTP response envelope structure
```

## License

Private — All rights reserved.
