import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section, Card, Button, CodeBlock } from "~/components/ui";

const STEPS = [
  {
    number: 1,
    title: "Create your account",
    description:
      "Sign up with Google at app.mailngine.com. No credit card required to get started.",
  },
  {
    number: 2,
    title: "Add your domain",
    description:
      "Add and verify your sending domain. We support auto-DNS configuration via Cloudflare or manual DNS record setup.",
  },
  {
    number: 3,
    title: "Get your API key",
    description:
      "Generate an API key from the dashboard. Keep it secret and never expose it in client-side code.",
  },
  {
    number: 4,
    title: "Send your first email",
    description:
      "Use any of our SDKs or call the REST API directly to send your first email.",
  },
] as const;

const SDK_CARDS = [
  {
    name: "Go SDK",
    install: "go get github.com/mailngine/mailngine-go",
    href: "https://github.com/mailngine/mailngine-go",
    language: "bash",
  },
  {
    name: "Node.js SDK",
    install: "npm install mailngine",
    href: "https://github.com/mailngine/mailngine-node",
    language: "bash",
  },
  {
    name: "Laravel SDK",
    install: "composer require mailngine/mailngine-laravel",
    href: "https://github.com/mailngine/mailngine-laravel",
    language: "bash",
  },
] as const;

export default component$(() => {
  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="docs-hero">
          <h1 class="docs-hero__title">Getting Started</h1>
          <p class="docs-hero__subtitle">
            Start sending emails with Mailngine in under 5 minutes. Follow
            these four steps to go from sign-up to your first delivered email.
          </p>
        </div>
      </Section>

      <Section padding="md">
        <div class="docs-steps">
          {STEPS.map((step) => (
            <div key={step.number} class="docs-step">
              <div class="docs-step__number">{step.number}</div>
              <div class="docs-step__content">
                <h3 class="docs-step__title">{step.title}</h3>
                <p class="docs-step__description">{step.description}</p>
              </div>
            </div>
          ))}
        </div>
      </Section>

      <Section background="surface" padding="md">
        <h2 class="docs-section-title">Quick Links</h2>
        <div class="docs-quick-links">
          <Card elevated>
            <h3>API Reference</h3>
            <p>
              Full REST API documentation with request and response examples for
              every endpoint.
            </p>
            <Button variant="secondary" size="sm" href="/docs/api">
              View API Reference
            </Button>
          </Card>

          {SDK_CARDS.map((sdk) => (
            <Card key={sdk.name} elevated>
              <h3>{sdk.name}</h3>
              <CodeBlock code={sdk.install} language={sdk.language} />
              <div class="docs-sdk-card__actions">
                <Button variant="secondary" size="sm" href={sdk.href}>
                  View Docs
                </Button>
              </div>
            </Card>
          ))}
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .docs-hero {
          text-align: center;
          max-width: 640px;
          margin: 0 auto;
        }
        .docs-hero__title {
          font-size: var(--font-size-4xl);
          margin-bottom: var(--space-4);
        }
        .docs-hero__subtitle {
          font-size: var(--font-size-lg);
          color: var(--color-text-secondary);
          line-height: 1.6;
        }

        /* --- Steps --- */
        .docs-steps {
          max-width: 720px;
          margin: 0 auto;
          display: flex;
          flex-direction: column;
          gap: var(--space-8);
        }
        .docs-step {
          display: flex;
          gap: var(--space-6);
          align-items: flex-start;
        }
        .docs-step__number {
          flex-shrink: 0;
          width: 40px;
          height: 40px;
          border-radius: 50%;
          background-color: var(--color-primary);
          color: var(--color-text-inverse);
          display: flex;
          align-items: center;
          justify-content: center;
          font-family: var(--font-primary);
          font-size: var(--font-size-lg);
          font-weight: var(--font-weight-semibold);
        }
        .docs-step__title {
          font-size: var(--font-size-xl);
          margin-bottom: var(--space-2);
        }
        .docs-step__description {
          font-size: var(--font-size-base);
          color: var(--color-text-secondary);
          line-height: 1.6;
        }

        /* --- Quick Links --- */
        .docs-section-title {
          text-align: center;
          margin-bottom: var(--space-8);
        }
        .docs-quick-links {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
          gap: var(--space-6);
        }
        .docs-quick-links h3 {
          margin-bottom: var(--space-3);
        }
        .docs-quick-links p {
          margin-bottom: var(--space-4);
          font-size: var(--font-size-sm);
          color: var(--color-text-secondary);
          line-height: 1.6;
        }
        .docs-quick-links .hm-code-block {
          margin-bottom: var(--space-4);
        }
        .docs-sdk-card__actions {
          margin-top: var(--space-2);
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "Documentation | Mailngine",
  meta: [
    {
      name: "description",
      content:
        "Get started with Mailngine. Learn how to send transactional and marketing emails with our API, Go SDK, Node.js SDK, and Laravel SDK.",
    },
  ],
};
