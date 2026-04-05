import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section, Card, Button, CodeBlock } from "~/components/ui";

const SDKS = [
  {
    name: "Go",
    install: "go get github.com/mailngine/mailngine-go",
    installLang: "bash",
    example: `package main

import (
  "context"
  mailngine "github.com/mailngine/mailngine-go"
)

func main() {
  client := mailngine.NewClient("mn_live_xxxxx")

  _, err := client.Emails.Send(context.Background(), &mailngine.SendEmailParams{
    From:    "hello@yourdomain.com",
    To:      []string{"user@example.com"},
    Subject: "Welcome!",
    HTML:    "<h1>Hello from Go</h1>",
  })
  if err != nil {
    panic(err)
  }
}`,
    exampleLang: "go",
    docsHref: "https://github.com/mailngine/mailngine-go",
  },
  {
    name: "Node.js",
    install: "npm install mailngine",
    installLang: "bash",
    example: `import { Mailngine } from 'mailngine';

const client = new Mailngine('mn_live_xxxxx');

await client.emails.send({
  from: 'hello@yourdomain.com',
  to: ['user@example.com'],
  subject: 'Welcome!',
  html: '<h1>Hello from Node.js</h1>',
});`,
    exampleLang: "typescript",
    docsHref: "https://github.com/mailngine/mailngine-node",
  },
  {
    name: "Laravel",
    install: "composer require mailngine/mailngine-laravel",
    installLang: "bash",
    example: `// config/mail.php — add to 'mailers' array:
'mailngine' => [
    'transport' => 'mailngine',
],

// .env
MAIL_MAILER=mailngine
MAILNGINE_API_KEY=mn_live_xxxxx

// Usage — standard Laravel Mail
Mail::to('user@example.com')
    ->send(new WelcomeEmail());`,
    exampleLang: "php",
    docsHref: "https://github.com/mailngine/mailngine-laravel",
  },
] as const;

const COMING_SOON = [
  { name: "Python" },
  { name: "Ruby" },
  { name: "Java" },
] as const;

export default component$(() => {
  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="integrations-hero">
          <h1 class="integrations-hero__title">Integrations</h1>
          <p class="integrations-hero__subtitle">
            Official SDKs and libraries to integrate Mailngine into your
            application. Get from zero to sending in minutes.
          </p>
        </div>
      </Section>

      <Section padding="md">
        <div class="sdk-grid">
          {SDKS.map((sdk) => (
            <Card key={sdk.name} class="sdk-card">
              <h2 class="sdk-card__name">{sdk.name}</h2>

              <h4 class="sdk-card__label">Install</h4>
              <CodeBlock code={sdk.install} language={sdk.installLang} />

              <h4 class="sdk-card__label">Send an email</h4>
              <CodeBlock code={sdk.example} language={sdk.exampleLang} />

              <div class="sdk-card__actions">
                <Button variant="primary" size="md" href={sdk.docsHref}>
                  Full Documentation
                </Button>
              </div>
            </Card>
          ))}
        </div>
      </Section>

      <Section background="surface" padding="md">
        <div class="coming-soon">
          <h2 class="coming-soon__title">More SDKs Coming Soon</h2>
          <p class="coming-soon__subtitle">
            We are actively working on official SDKs for more languages. Want
            to help? Contributions are welcome on GitHub.
          </p>
          <div class="coming-soon__grid">
            {COMING_SOON.map((lang) => (
              <Card key={lang.name} class="coming-soon__card">
                <span class="coming-soon__name">{lang.name}</span>
                <span class="coming-soon__badge">Coming Soon</span>
              </Card>
            ))}
          </div>
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .integrations-hero {
          text-align: center;
          max-width: 640px;
          margin: 0 auto;
        }
        .integrations-hero__title {
          font-size: var(--font-size-4xl);
          margin-bottom: var(--space-4);
        }
        .integrations-hero__subtitle {
          font-size: var(--font-size-lg);
          color: var(--color-text-secondary);
          line-height: 1.6;
        }

        /* --- SDK Grid --- */
        .sdk-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
          gap: var(--space-8);
        }
        @media (max-width: 480px) {
          .sdk-grid {
            grid-template-columns: 1fr;
          }
        }
        .sdk-card__name {
          font-size: var(--font-size-2xl);
          margin-bottom: var(--space-5);
        }
        .sdk-card__label {
          font-size: var(--font-size-xs);
          text-transform: uppercase;
          letter-spacing: 0.08em;
          color: var(--color-text-disabled);
          margin-bottom: var(--space-2);
          margin-top: var(--space-5);
        }
        .sdk-card__label:first-of-type {
          margin-top: 0;
        }
        .sdk-card .hm-code-block {
          margin-bottom: 0;
        }
        .sdk-card__actions {
          margin-top: var(--space-6);
        }

        /* --- Coming Soon --- */
        .coming-soon {
          text-align: center;
          max-width: 640px;
          margin: 0 auto;
        }
        .coming-soon__title {
          margin-bottom: var(--space-3);
        }
        .coming-soon__subtitle {
          font-size: var(--font-size-base);
          color: var(--color-text-secondary);
          margin-bottom: var(--space-8);
          line-height: 1.6;
        }
        .coming-soon__grid {
          display: flex;
          gap: var(--space-4);
          justify-content: center;
          flex-wrap: wrap;
        }
        .coming-soon__card {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: var(--space-2);
          min-width: 140px;
        }
        .coming-soon__name {
          font-family: var(--font-primary);
          font-size: var(--font-size-lg);
          font-weight: var(--font-weight-medium);
          color: var(--color-text-primary);
        }
        .coming-soon__badge {
          font-size: var(--font-size-xs);
          color: var(--color-text-disabled);
          background-color: var(--color-bg-surface-alt);
          padding: var(--space-1) var(--space-3);
          border-radius: var(--radius-pill);
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "Integrations | Mailngine",
  meta: [
    {
      name: "description",
      content:
        "Official Mailngine SDKs for Go, Node.js, and Laravel. Install, configure, and start sending emails in minutes.",
    },
  ],
};
