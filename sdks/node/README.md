# mailngine

The official Node.js/TypeScript SDK for [Mailngine](https://mailngine.com).

## Requirements

- Node.js 18+ (uses the built-in `fetch` API)

## Installation

```bash
npm install mailngine
```

## Quick start

```typescript
import { Mailngine } from 'mailngine';

const client = new Mailngine('mn_live_...');

const email = await client.emails.send({
  from: 'hello@example.com',
  to: ['user@example.com'],
  subject: 'Hello!',
  html: '<h1>Welcome</h1>',
});

console.log(email.id);
```

## Resources

### Emails

```typescript
// Send an email
const email = await client.emails.send({
  from: 'hello@example.com',
  to: ['user@example.com'],
  subject: 'Welcome!',
  html: '<h1>Hello</h1>',
});

// Get an email by ID
const retrieved = await client.emails.get(email.id);

// List emails (paginated)
const list = await client.emails.list({ page: 1, perPage: 20 });
```

### Domains

```typescript
// Add a sending domain
const { domain, dns_records } = await client.domains.create({ name: 'example.com' });

// Verify DNS records
const records = await client.domains.verify(domain.id);

// Update tracking settings
await client.domains.update(domain.id, { open_tracking: true, click_tracking: true });

// List, get, delete
const domains = await client.domains.list();
const single = await client.domains.get(domain.id);
await client.domains.delete(domain.id);
```

### Templates

```typescript
// Create a template
const template = await client.templates.create({
  name: 'welcome',
  subject: 'Welcome, {{name}}!',
  html_body: '<h1>Hello {{name}}</h1>',
  variables: ['name'],
});

// Preview with data
const preview = await client.templates.preview(template.id, {
  data: { name: 'Alice' },
});

// Update, list, get, delete
await client.templates.update(template.id, { ...params });
const templates = await client.templates.list();
await client.templates.delete(template.id);
```

### Webhooks

```typescript
// Create a webhook
const webhook = await client.webhooks.create({
  url: 'https://example.com/webhook',
  events: ['email.delivered', 'email.bounced'],
});

// Update
await client.webhooks.update(webhook.id, {
  url: 'https://example.com/webhook',
  events: ['email.delivered'],
  is_active: true,
});

// List deliveries
const deliveries = await client.webhooks.listDeliveries(webhook.id);

// List, get, delete
const webhooks = await client.webhooks.list();
await client.webhooks.delete(webhook.id);
```

### API Keys

```typescript
// Create a key (full value is only returned once)
const key = await client.apiKeys.create({
  name: 'Production',
  permission: 'send_only',
});
console.log(key.key); // mn_live_...

// List keys
const keys = await client.apiKeys.list();

// Revoke a key
await client.apiKeys.revoke(key.id);
```

## Error handling

```typescript
import { Mailngine, MailngineError } from 'mailngine';

const client = new Mailngine('mn_live_...');

try {
  await client.emails.send({ /* ... */ });
} catch (err) {
  if (err instanceof MailngineError) {
    console.error(err.statusCode); // 422
    console.error(err.code);       // "domain_not_verified"
    console.error(err.message);    // "The sending domain has not been verified"
  }
}
```

## Configuration

```typescript
// Custom base URL (self-hosted or staging)
const client = new Mailngine('mn_live_...', {
  baseURL: 'https://api.staging.mailngine.com',
});
```

## License

MIT
