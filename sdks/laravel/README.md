# Mailngine Laravel SDK

Official PHP/Laravel SDK for the [Mailngine](https://mailngine.com) email API.

## Requirements

- PHP 8.1+
- Guzzle 7.x
- Symfony Mailer 6.x or 7.x

## Installation

```bash
composer require mailngine/mailngine-laravel
```

## Quick Start

### Direct API Usage

```php
use Mailngine\Mailngine;

$mailngine = new Mailngine('your-api-key');

// Send an email
$email = $mailngine->emails()->send([
    'from' => 'you@example.com',
    'to' => ['recipient@example.com'],
    'subject' => 'Hello from Mailngine',
    'html' => '<h1>Welcome!</h1><p>Thanks for signing up.</p>',
    'text' => 'Welcome! Thanks for signing up.',
]);

echo $email['id']; // Email UUID
```

### Laravel Mail Integration

Add the Mailngine mailer to `config/mail.php`:

```php
'mailers' => [
    'mailngine' => [
        'transport' => 'mailngine',
        'key' => env('MAILNGINE_API_KEY'),
    ],
],
```

Register the transport in a service provider (e.g. `AppServiceProvider`):

```php
use Mailngine\Mailngine;
use Mailngine\MailngineTransport;
use Illuminate\Support\Facades\Mail;

public function boot(): void
{
    Mail::extend('mailngine', function (array $config) {
        $client = new Mailngine($config['key']);
        return new MailngineTransport($client);
    });
}
```

Set your API key in `.env`:

```
MAILNGINE_API_KEY=mn_your_api_key_here
```

Then send emails using Laravel's standard Mail facade:

```php
use Illuminate\Support\Facades\Mail;

Mail::mailer('mailngine')->to('user@example.com')->send(new WelcomeEmail());
```

## API Resources

### Emails

```php
// Send an email
$email = $mailngine->emails()->send([
    'from' => 'you@example.com',
    'to' => ['recipient@example.com'],
    'subject' => 'Order Confirmation',
    'html' => '<h1>Order #1234 confirmed</h1>',
    'tags' => ['order-confirmation'],
]);

// Send with a template
$email = $mailngine->emails()->send([
    'from' => 'you@example.com',
    'to' => ['recipient@example.com'],
    'template_id' => 'tmpl_abc123',
    'template_data' => ['name' => 'Alice', 'order_id' => '1234'],
]);

// Get an email by ID
$email = $mailngine->emails()->get('email-uuid');

// List emails (paginated)
$emails = $mailngine->emails()->list(['page' => 1, 'per_page' => 50]);
```

### Domains

```php
// Add a sending domain
$domain = $mailngine->domains()->create('example.com');

// List all domains
$domains = $mailngine->domains()->list();

// Get a domain
$domain = $mailngine->domains()->get('domain-uuid');

// Update tracking settings
$domain = $mailngine->domains()->update('domain-uuid', [
    'open_tracking' => true,
    'click_tracking' => true,
]);

// Verify DNS records
$records = $mailngine->domains()->verify('domain-uuid');

// Delete a domain
$mailngine->domains()->delete('domain-uuid');
```

### Webhooks

```php
// Create a webhook
$webhook = $mailngine->webhooks()->create([
    'url' => 'https://example.com/webhooks/mailngine',
    'events' => ['email.delivered', 'email.bounced', 'email.complained'],
]);

// List webhooks
$webhooks = $mailngine->webhooks()->list();

// Update a webhook
$webhook = $mailngine->webhooks()->update('webhook-uuid', [
    'url' => 'https://example.com/webhooks/mailngine',
    'events' => ['email.delivered'],
    'is_active' => true,
]);

// List delivery attempts
$deliveries = $mailngine->webhooks()->listDeliveries('webhook-uuid');

// Delete a webhook
$mailngine->webhooks()->delete('webhook-uuid');
```

### Templates

```php
// Create a template
$template = $mailngine->templates()->create([
    'name' => 'Welcome Email',
    'subject' => 'Welcome, {{name}}!',
    'html_body' => '<h1>Hello {{name}}</h1><p>Welcome to our service.</p>',
    'text_body' => 'Hello {{name}}, welcome to our service.',
    'variables' => ['name'],
]);

// List templates
$templates = $mailngine->templates()->list();

// Preview a rendered template
$preview = $mailngine->templates()->preview('template-uuid', [
    'name' => 'Alice',
]);

// Update a template
$template = $mailngine->templates()->update('template-uuid', [
    'name' => 'Welcome Email v2',
    'subject' => 'Welcome aboard, {{name}}!',
    'html_body' => '<h1>Welcome aboard, {{name}}!</h1>',
]);

// Delete a template
$mailngine->templates()->delete('template-uuid');
```

### API Keys

```php
// Create an API key (full key only returned once)
$key = $mailngine->apiKeys()->create([
    'name' => 'Production Sending Key',
    'permission' => 'send_only',
]);
echo $key['key']; // Store this securely

// List API keys
$keys = $mailngine->apiKeys()->list();

// Revoke an API key
$mailngine->apiKeys()->revoke('key-uuid');
```

## Error Handling

```php
use Mailngine\Exceptions\ApiException;
use Mailngine\Exceptions\MailngineException;

try {
    $email = $mailngine->emails()->send([...]);
} catch (ApiException $e) {
    // API returned an error response
    echo $e->statusCode;  // HTTP status code (e.g. 422)
    echo $e->errorCode;   // API error code (e.g. "domain_not_verified")
    echo $e->getMessage(); // Human-readable error message
} catch (MailngineException $e) {
    // Connection or unexpected error
    echo $e->getMessage();
}
```

## Configuration Options

```php
$mailngine = new Mailngine('your-api-key', [
    'base_url' => 'https://api.mailngine.com', // Custom API base URL
    'timeout' => 30,                            // Request timeout in seconds
]);
```

## License

MIT
