# Hello Mail Laravel SDK

Official PHP/Laravel SDK for the [Hello Mail](https://hellomail.dev) email API.

## Requirements

- PHP 8.1+
- Guzzle 7.x
- Symfony Mailer 6.x or 7.x

## Installation

```bash
composer require hellomail/hellomail-laravel
```

## Quick Start

### Direct API Usage

```php
use HelloMail\HelloMail;

$hellomail = new HelloMail('your-api-key');

// Send an email
$email = $hellomail->emails()->send([
    'from' => 'you@example.com',
    'to' => ['recipient@example.com'],
    'subject' => 'Hello from Hello Mail',
    'html' => '<h1>Welcome!</h1><p>Thanks for signing up.</p>',
    'text' => 'Welcome! Thanks for signing up.',
]);

echo $email['id']; // Email UUID
```

### Laravel Mail Integration

Add the Hello Mail mailer to `config/mail.php`:

```php
'mailers' => [
    'hellomail' => [
        'transport' => 'hellomail',
        'key' => env('HELLOMAIL_API_KEY'),
    ],
],
```

Register the transport in a service provider (e.g. `AppServiceProvider`):

```php
use HelloMail\HelloMail;
use HelloMail\HelloMailTransport;
use Illuminate\Support\Facades\Mail;

public function boot(): void
{
    Mail::extend('hellomail', function (array $config) {
        $client = new HelloMail($config['key']);
        return new HelloMailTransport($client);
    });
}
```

Set your API key in `.env`:

```
HELLOMAIL_API_KEY=hm_your_api_key_here
```

Then send emails using Laravel's standard Mail facade:

```php
use Illuminate\Support\Facades\Mail;

Mail::mailer('hellomail')->to('user@example.com')->send(new WelcomeEmail());
```

## API Resources

### Emails

```php
// Send an email
$email = $hellomail->emails()->send([
    'from' => 'you@example.com',
    'to' => ['recipient@example.com'],
    'subject' => 'Order Confirmation',
    'html' => '<h1>Order #1234 confirmed</h1>',
    'tags' => ['order-confirmation'],
]);

// Send with a template
$email = $hellomail->emails()->send([
    'from' => 'you@example.com',
    'to' => ['recipient@example.com'],
    'template_id' => 'tmpl_abc123',
    'template_data' => ['name' => 'Alice', 'order_id' => '1234'],
]);

// Get an email by ID
$email = $hellomail->emails()->get('email-uuid');

// List emails (paginated)
$emails = $hellomail->emails()->list(['page' => 1, 'per_page' => 50]);
```

### Domains

```php
// Add a sending domain
$domain = $hellomail->domains()->create('example.com');

// List all domains
$domains = $hellomail->domains()->list();

// Get a domain
$domain = $hellomail->domains()->get('domain-uuid');

// Update tracking settings
$domain = $hellomail->domains()->update('domain-uuid', [
    'open_tracking' => true,
    'click_tracking' => true,
]);

// Verify DNS records
$records = $hellomail->domains()->verify('domain-uuid');

// Delete a domain
$hellomail->domains()->delete('domain-uuid');
```

### Webhooks

```php
// Create a webhook
$webhook = $hellomail->webhooks()->create([
    'url' => 'https://example.com/webhooks/hellomail',
    'events' => ['email.delivered', 'email.bounced', 'email.complained'],
]);

// List webhooks
$webhooks = $hellomail->webhooks()->list();

// Update a webhook
$webhook = $hellomail->webhooks()->update('webhook-uuid', [
    'url' => 'https://example.com/webhooks/hellomail',
    'events' => ['email.delivered'],
    'is_active' => true,
]);

// List delivery attempts
$deliveries = $hellomail->webhooks()->listDeliveries('webhook-uuid');

// Delete a webhook
$hellomail->webhooks()->delete('webhook-uuid');
```

### Templates

```php
// Create a template
$template = $hellomail->templates()->create([
    'name' => 'Welcome Email',
    'subject' => 'Welcome, {{name}}!',
    'html_body' => '<h1>Hello {{name}}</h1><p>Welcome to our service.</p>',
    'text_body' => 'Hello {{name}}, welcome to our service.',
    'variables' => ['name'],
]);

// List templates
$templates = $hellomail->templates()->list();

// Preview a rendered template
$preview = $hellomail->templates()->preview('template-uuid', [
    'name' => 'Alice',
]);

// Update a template
$template = $hellomail->templates()->update('template-uuid', [
    'name' => 'Welcome Email v2',
    'subject' => 'Welcome aboard, {{name}}!',
    'html_body' => '<h1>Welcome aboard, {{name}}!</h1>',
]);

// Delete a template
$hellomail->templates()->delete('template-uuid');
```

### API Keys

```php
// Create an API key (full key only returned once)
$key = $hellomail->apiKeys()->create([
    'name' => 'Production Sending Key',
    'permission' => 'send_only',
]);
echo $key['key']; // Store this securely

// List API keys
$keys = $hellomail->apiKeys()->list();

// Revoke an API key
$hellomail->apiKeys()->revoke('key-uuid');
```

## Error Handling

```php
use HelloMail\Exceptions\ApiException;
use HelloMail\Exceptions\HelloMailException;

try {
    $email = $hellomail->emails()->send([...]);
} catch (ApiException $e) {
    // API returned an error response
    echo $e->statusCode;  // HTTP status code (e.g. 422)
    echo $e->errorCode;   // API error code (e.g. "domain_not_verified")
    echo $e->getMessage(); // Human-readable error message
} catch (HelloMailException $e) {
    // Connection or unexpected error
    echo $e->getMessage();
}
```

## Configuration Options

```php
$hellomail = new HelloMail('your-api-key', [
    'base_url' => 'https://api.hellomail.dev', // Custom API base URL
    'timeout' => 30,                           // Request timeout in seconds
]);
```

## License

MIT
