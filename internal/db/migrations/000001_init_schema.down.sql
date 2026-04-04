-- Drop triggers
DROP TRIGGER IF EXISTS trg_templates_updated_at ON templates;
DROP TRIGGER IF EXISTS trg_webhooks_updated_at ON webhooks;
DROP TRIGGER IF EXISTS trg_domains_updated_at ON domains;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP TRIGGER IF EXISTS trg_organizations_updated_at ON organizations;
DROP FUNCTION IF EXISTS update_updated_at();

-- Drop views
DROP VIEW IF EXISTS postfix_virtual_mailboxes;
DROP VIEW IF EXISTS postfix_virtual_domains;

-- Drop tables in reverse order of creation (respecting foreign keys)
DROP TABLE IF EXISTS ip_pools;
DROP TABLE IF EXISTS usage_monthly;
DROP TABLE IF EXISTS usage_daily;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS suppressions;
DROP TABLE IF EXISTS webhook_deliveries;
DROP TABLE IF EXISTS webhooks;
DROP TABLE IF EXISTS inbox_message_labels;
DROP TABLE IF EXISTS inbox_labels;
DROP TABLE IF EXISTS inbox_message_attachments;
DROP TABLE IF EXISTS inbox_messages;
DROP TABLE IF EXISTS inbox_threads;
DROP TABLE IF EXISTS email_events;
DROP TABLE IF EXISTS email_attachments;
DROP TABLE IF EXISTS emails;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS dns_records;
DROP TABLE IF EXISTS domains;
DROP TABLE IF EXISTS org_members;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;

-- Drop extensions
DROP EXTENSION IF EXISTS "pg_trgm";
DROP EXTENSION IF EXISTS "uuid-ossp";
