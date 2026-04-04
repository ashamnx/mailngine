-- Hello Mail Initial Schema
-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================
-- ORGANIZATIONS & USERS
-- ============================================================

CREATE TABLE organizations (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(100) NOT NULL UNIQUE,
    plan            VARCHAR(50) NOT NULL DEFAULT 'free',
    monthly_limit   INT NOT NULL DEFAULT 100,
    overage_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           VARCHAR(320) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    avatar_url      TEXT,
    google_id       VARCHAR(255) NOT NULL UNIQUE,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE org_members (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL DEFAULT 'member',
    invited_by      UUID REFERENCES users(id),
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, user_id)
);
CREATE INDEX idx_org_members_org ON org_members(org_id);
CREATE INDEX idx_org_members_user ON org_members(user_id);

-- ============================================================
-- DOMAINS
-- ============================================================

CREATE TABLE domains (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name            VARCHAR(253) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    region          VARCHAR(20) NOT NULL DEFAULT 'global',
    dkim_private_key TEXT,
    dkim_selector   VARCHAR(63) NOT NULL DEFAULT 'hm1',
    open_tracking   BOOLEAN NOT NULL DEFAULT TRUE,
    click_tracking  BOOLEAN NOT NULL DEFAULT TRUE,
    cloudflare_zone_id VARCHAR(64),
    cloudflare_api_token_enc TEXT,
    verified_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, name)
);
CREATE INDEX idx_domains_org ON domains(org_id);
CREATE INDEX idx_domains_name ON domains(name);

CREATE TABLE dns_records (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    domain_id       UUID NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    record_type     VARCHAR(10) NOT NULL,
    host            VARCHAR(253) NOT NULL,
    value           TEXT NOT NULL,
    purpose         VARCHAR(20) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    verified_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_dns_records_domain ON dns_records(domain_id);

-- ============================================================
-- API KEYS
-- ============================================================

CREATE TABLE api_keys (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    domain_id       UUID REFERENCES domains(id) ON DELETE SET NULL,
    name            VARCHAR(255) NOT NULL,
    prefix          VARCHAR(12) NOT NULL,
    key_hash        VARCHAR(128) NOT NULL,
    permission      VARCHAR(20) NOT NULL DEFAULT 'full',
    last_used_at    TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at      TIMESTAMPTZ
);
CREATE INDEX idx_api_keys_org ON api_keys(org_id);
CREATE INDEX idx_api_keys_prefix ON api_keys(prefix);
CREATE INDEX idx_api_keys_hash ON api_keys(key_hash);

-- ============================================================
-- TEMPLATES (must be before emails for FK)
-- ============================================================

CREATE TABLE templates (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    subject         TEXT NOT NULL,
    html_body       TEXT NOT NULL,
    text_body       TEXT,
    variables       JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, name)
);
CREATE INDEX idx_templates_org ON templates(org_id);

-- ============================================================
-- EMAILS (OUTBOUND)
-- ============================================================

CREATE TABLE emails (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id),
    domain_id       UUID NOT NULL REFERENCES domains(id),
    api_key_id      UUID REFERENCES api_keys(id),
    idempotency_key VARCHAR(255),
    from_address    VARCHAR(320) NOT NULL,
    from_name       VARCHAR(255),
    to_addresses    JSONB NOT NULL,
    cc_addresses    JSONB,
    bcc_addresses   JSONB,
    reply_to        VARCHAR(320),
    subject         TEXT NOT NULL,
    text_body_key   TEXT,
    html_body_key   TEXT,
    headers         JSONB,
    tags            JSONB,
    template_id     UUID REFERENCES templates(id),
    template_data   JSONB,
    status          VARCHAR(20) NOT NULL DEFAULT 'queued',
    scheduled_at    TIMESTAMPTZ,
    sent_at         TIMESTAMPTZ,
    delivered_at    TIMESTAMPTZ,
    message_id      VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, idempotency_key)
);
CREATE INDEX idx_emails_org ON emails(org_id);
CREATE INDEX idx_emails_domain ON emails(domain_id);
CREATE INDEX idx_emails_status ON emails(status);
CREATE INDEX idx_emails_created ON emails(created_at);
CREATE INDEX idx_emails_message_id ON emails(message_id);

CREATE TABLE email_attachments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email_id        UUID NOT NULL REFERENCES emails(id) ON DELETE CASCADE,
    filename        VARCHAR(255) NOT NULL,
    content_type    VARCHAR(127) NOT NULL,
    size_bytes      BIGINT NOT NULL,
    storage_key     TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_email_attachments_email ON email_attachments(email_id);

-- ============================================================
-- EMAIL EVENTS
-- ============================================================

CREATE TABLE email_events (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email_id        UUID NOT NULL REFERENCES emails(id),
    org_id          UUID NOT NULL REFERENCES organizations(id),
    event_type      VARCHAR(30) NOT NULL,
    recipient       VARCHAR(320) NOT NULL,
    metadata        JSONB,
    ip_address      INET,
    user_agent      TEXT,
    occurred_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_email_events_email ON email_events(email_id);
CREATE INDEX idx_email_events_org_type ON email_events(org_id, event_type);
CREATE INDEX idx_email_events_occurred ON email_events(occurred_at);

-- ============================================================
-- INBOX (INBOUND EMAILS)
-- ============================================================

CREATE TABLE inbox_threads (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    domain_id       UUID NOT NULL REFERENCES domains(id),
    subject         TEXT NOT NULL,
    participant_addresses JSONB NOT NULL,
    last_message_at TIMESTAMPTZ NOT NULL,
    message_count   INT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_inbox_threads_org ON inbox_threads(org_id);
CREATE INDEX idx_inbox_threads_domain ON inbox_threads(domain_id);
CREATE INDEX idx_inbox_threads_last_msg ON inbox_threads(last_message_at DESC);

CREATE TABLE inbox_messages (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    domain_id       UUID NOT NULL REFERENCES domains(id),
    thread_id       UUID REFERENCES inbox_threads(id),
    message_id_header VARCHAR(512),
    in_reply_to     VARCHAR(512),
    references_header TEXT,
    from_address    VARCHAR(320) NOT NULL,
    from_name       VARCHAR(255),
    to_addresses    JSONB NOT NULL,
    cc_addresses    JSONB,
    subject         TEXT NOT NULL,
    text_body_key   TEXT,
    html_body_key   TEXT,
    snippet         VARCHAR(255),
    is_read         BOOLEAN NOT NULL DEFAULT FALSE,
    is_starred      BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived     BOOLEAN NOT NULL DEFAULT FALSE,
    is_trashed      BOOLEAN NOT NULL DEFAULT FALSE,
    received_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_inbox_messages_org ON inbox_messages(org_id);
CREATE INDEX idx_inbox_messages_thread ON inbox_messages(thread_id);
CREATE INDEX idx_inbox_messages_domain ON inbox_messages(domain_id);
CREATE INDEX idx_inbox_messages_received ON inbox_messages(received_at DESC);
CREATE INDEX idx_inbox_messages_msgid ON inbox_messages(message_id_header);
CREATE INDEX idx_inbox_messages_subject_trgm ON inbox_messages USING gin (subject gin_trgm_ops);

CREATE TABLE inbox_message_attachments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id      UUID NOT NULL REFERENCES inbox_messages(id) ON DELETE CASCADE,
    filename        VARCHAR(255) NOT NULL,
    content_type    VARCHAR(127) NOT NULL,
    size_bytes      BIGINT NOT NULL,
    storage_key     TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_inbox_attachments_msg ON inbox_message_attachments(message_id);

CREATE TABLE inbox_labels (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    color           VARCHAR(7),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, name)
);

CREATE TABLE inbox_message_labels (
    message_id      UUID NOT NULL REFERENCES inbox_messages(id) ON DELETE CASCADE,
    label_id        UUID NOT NULL REFERENCES inbox_labels(id) ON DELETE CASCADE,
    PRIMARY KEY (message_id, label_id)
);

-- ============================================================
-- WEBHOOKS
-- ============================================================

CREATE TABLE webhooks (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    url             TEXT NOT NULL,
    events          JSONB NOT NULL,
    secret          VARCHAR(128) NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_webhooks_org ON webhooks(org_id);

CREATE TABLE webhook_deliveries (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    webhook_id      UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event_type      VARCHAR(30) NOT NULL,
    payload         JSONB NOT NULL,
    response_status INT,
    response_body   TEXT,
    attempt         INT NOT NULL DEFAULT 1,
    delivered_at    TIMESTAMPTZ,
    next_retry_at   TIMESTAMPTZ,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_webhook_deliveries_webhook ON webhook_deliveries(webhook_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);

-- ============================================================
-- SUPPRESSION LISTS
-- ============================================================

CREATE TABLE suppressions (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email           VARCHAR(320) NOT NULL,
    reason          VARCHAR(30) NOT NULL,
    metadata        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, email)
);
CREATE INDEX idx_suppressions_org_email ON suppressions(org_id, email);

-- ============================================================
-- AUDIT LOGS
-- ============================================================

CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id),
    user_id         UUID REFERENCES users(id),
    api_key_id      UUID REFERENCES api_keys(id),
    action          VARCHAR(100) NOT NULL,
    resource_type   VARCHAR(50) NOT NULL,
    resource_id     UUID,
    metadata        JSONB,
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_logs_org ON audit_logs(org_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- ============================================================
-- BILLING / USAGE TRACKING
-- ============================================================

CREATE TABLE usage_daily (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    date            DATE NOT NULL,
    emails_sent     INT NOT NULL DEFAULT 0,
    emails_delivered INT NOT NULL DEFAULT 0,
    emails_bounced  INT NOT NULL DEFAULT 0,
    emails_received INT NOT NULL DEFAULT 0,
    api_calls       INT NOT NULL DEFAULT 0,
    UNIQUE(org_id, date)
);
CREATE INDEX idx_usage_daily_org_date ON usage_daily(org_id, date);

CREATE TABLE usage_monthly (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    month           DATE NOT NULL,
    total_sent      INT NOT NULL DEFAULT 0,
    total_received  INT NOT NULL DEFAULT 0,
    total_api_calls INT NOT NULL DEFAULT 0,
    plan_limit      INT NOT NULL,
    overage_count   INT NOT NULL DEFAULT 0,
    UNIQUE(org_id, month)
);

-- ============================================================
-- IP POOLS & WARMUP
-- ============================================================

CREATE TABLE ip_pools (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(100) NOT NULL,
    ip_address      INET NOT NULL UNIQUE,
    is_dedicated    BOOLEAN NOT NULL DEFAULT FALSE,
    org_id          UUID REFERENCES organizations(id),
    warmup_phase    INT NOT NULL DEFAULT 0,
    warmup_started  TIMESTAMPTZ,
    daily_limit     INT NOT NULL DEFAULT 50,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- POSTFIX VIRTUAL DOMAIN LOOKUP VIEWS
-- ============================================================

CREATE VIEW postfix_virtual_domains AS
    SELECT name AS domain FROM domains WHERE status = 'verified';

CREATE VIEW postfix_virtual_mailboxes AS
    SELECT name AS domain FROM domains WHERE status = 'verified';

-- ============================================================
-- UPDATED_AT TRIGGER
-- ============================================================

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_organizations_updated_at
    BEFORE UPDATE ON organizations FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_domains_updated_at
    BEFORE UPDATE ON domains FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_webhooks_updated_at
    BEFORE UPDATE ON webhooks FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_templates_updated_at
    BEFORE UPDATE ON templates FOR EACH ROW EXECUTE FUNCTION update_updated_at();
