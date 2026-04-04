import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section, Button, CodeBlock } from "~/components/ui";
import { Features } from "~/components/home/features";
import { CodeExample } from "~/components/home/code-example";
import { Stats } from "~/components/home/stats";

export default component$(() => {
  return (
    <>
      {/* Hero */}
      <Section class="hero-section" padding="lg">
        <div class="hero">
          <h1 class="hero__title">Email infrastructure for developers</h1>
          <p class="hero__subtitle">
            Reliable email delivery with powerful APIs. Send, receive, and
            manage email at any scale with Go, Node.js, and Laravel SDKs.
          </p>
          <div class="hero__actions">
            <Button variant="primary" size="lg" href="/pricing">
              Get Started Free
            </Button>
            <Button variant="secondary" size="lg" href="/docs">
              Read the Docs
            </Button>
          </div>
          <div class="hero__code-snippet">
            <CodeBlock
              language="typescript"
              code={`const email = await client.emails.send({
  from: 'hello@yourdomain.com',
  to: ['user@example.com'],
  subject: 'Welcome!',
  html: '<h1>Hello from Node.js</h1>',
});`}
            />
          </div>
        </div>
      </Section>

      {/* Features */}
      <Section background="surface" padding="lg" id="features">
        <h2 class="section-heading">Everything you need to send email</h2>
        <p class="section-subheading">
          From transactional emails to inbound processing, Hello Mail handles
          the full email lifecycle.
        </p>
        <Features />
      </Section>

      {/* Code Example */}
      <Section padding="lg" id="sdks">
        <CodeExample />
      </Section>

      {/* Stats */}
      <Section background="surface" padding="md" id="stats">
        <Stats />
      </Section>

      {/* CTA */}
      <Section background="primary" padding="lg">
        <div class="cta-block">
          <h2 class="cta-block__heading">Start sending emails in minutes</h2>
          <Button
            variant="secondary"
            size="lg"
            href="/pricing"
            class="cta-block__button"
          >
            Get Started Free
          </Button>
        </div>
      </Section>
    </>
  );
});

export const head: DocumentHead = {
  title: "Hello Mail - Email Infrastructure for Developers",
  meta: [
    {
      name: "description",
      content:
        "Reliable email API with Go, Node.js, and Laravel SDKs. Domain management, analytics, webhooks, and inbox support.",
    },
    {
      property: "og:title",
      content: "Hello Mail - Email Infrastructure for Developers",
    },
    {
      property: "og:description",
      content:
        "Reliable email API with Go, Node.js, and Laravel SDKs.",
    },
    {
      property: "og:type",
      content: "website",
    },
    {
      name: "twitter:card",
      content: "summary_large_image",
    },
  ],
};
