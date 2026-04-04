import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section } from "~/components/ui/section";
import { Card } from "~/components/ui/card";
import { Button } from "~/components/ui/button";

const benefits = [
  {
    icon: "public",
    title: "Remote-first",
    description:
      "Work from anywhere. Our team is distributed across time zones and we've built our processes around asynchronous collaboration.",
  },
  {
    icon: "schedule",
    title: "Flexible Hours",
    description:
      "We care about outcomes, not hours logged. Structure your day in a way that works for your life.",
  },
  {
    icon: "code",
    title: "Open Source Time",
    description:
      "Spend 20% of your time contributing to open source projects. Give back to the communities that power our stack.",
  },
  {
    icon: "school",
    title: "Learning Budget",
    description:
      "Annual budget for conferences, courses, books, and anything else that helps you grow as a professional.",
  },
  {
    icon: "health_and_safety",
    title: "Health Benefits",
    description:
      "Comprehensive health, dental, and vision coverage for you and your dependents.",
  },
  {
    icon: "trending_up",
    title: "Equity",
    description:
      "Meaningful equity in the company so you share in the success you help build.",
  },
];

export default component$(() => {
  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="careers-hero">
          <h1>Join Us</h1>
          <p>Help us build the future of email infrastructure</p>
        </div>
      </Section>

      <Section padding="lg">
        <div class="careers-benefits">
          <h2 class="careers-section-title">Why Hello Mail</h2>
          <div class="careers-benefits__grid">
            {benefits.map((benefit) => (
              <Card key={benefit.title} elevated>
                <div class="careers-benefit-card">
                  <span class="material-symbols-rounded careers-benefit-card__icon">
                    {benefit.icon}
                  </span>
                  <h3>{benefit.title}</h3>
                  <p>{benefit.description}</p>
                </div>
              </Card>
            ))}
          </div>
        </div>
      </Section>

      <Section background="surface" padding="lg">
        <div class="careers-positions">
          <h2 class="careers-section-title">Open Positions</h2>
          <div class="careers-positions__empty">
            <span class="material-symbols-rounded careers-positions__icon">
              work_off
            </span>
            <p>
              No open positions right now. Check back soon or follow us for
              updates.
            </p>
          </div>
        </div>
      </Section>

      <Section background="primary" padding="lg">
        <div class="careers-cta">
          <h2>Interested?</h2>
          <p>
            Even if we don't have an open role that matches your skills, we'd
            love to hear from you. Send us an email and tell us what you're
            passionate about.
          </p>
          <Button variant="secondary" size="lg" href="mailto:careers@hellomail.dev">
            careers@hellomail.dev
          </Button>
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .careers-hero {
          text-align: center;
          max-width: 600px;
          margin: 0 auto;
        }
        .careers-hero h1 {
          font-size: var(--font-size-hero);
          margin-bottom: var(--space-4);
        }
        .careers-hero p {
          font-size: var(--font-size-lg);
        }

        .careers-section-title {
          text-align: center;
          margin-bottom: var(--space-10);
        }

        .careers-benefits__grid {
          display: grid;
          grid-template-columns: repeat(3, 1fr);
          gap: var(--space-6);
        }
        .careers-benefit-card {
          padding: var(--space-2) 0;
        }
        .careers-benefit-card__icon {
          font-size: 32px;
          color: var(--color-primary);
          margin-bottom: var(--space-3);
        }
        .careers-benefit-card h3 {
          font-size: var(--font-size-lg);
          margin-bottom: var(--space-2);
        }
        .careers-benefit-card p {
          font-size: var(--font-size-sm);
          line-height: 1.6;
        }

        .careers-positions {
          text-align: center;
        }
        .careers-positions__empty {
          max-width: 480px;
          margin: 0 auto;
          padding: var(--space-8) 0;
        }
        .careers-positions__icon {
          font-size: 48px;
          color: var(--color-text-disabled);
          margin-bottom: var(--space-4);
        }
        .careers-positions__empty p {
          font-size: var(--font-size-lg);
        }

        .careers-cta {
          text-align: center;
          max-width: 560px;
          margin: 0 auto;
        }
        .careers-cta h2 {
          color: var(--color-text-inverse);
          margin-bottom: var(--space-4);
        }
        .careers-cta p {
          margin-bottom: var(--space-8);
        }
        .careers-cta .hm-button--secondary {
          color: var(--color-text-inverse);
          border-color: var(--color-text-inverse);
        }
        .careers-cta .hm-button--secondary:hover {
          background-color: rgba(255, 255, 255, 0.15);
        }

        @media (max-width: 768px) {
          .careers-hero h1 {
            font-size: var(--font-size-3xl);
          }
          .careers-benefits__grid {
            grid-template-columns: repeat(2, 1fr);
          }
        }
        @media (max-width: 480px) {
          .careers-benefits__grid {
            grid-template-columns: 1fr;
          }
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "Careers - Hello Mail",
  meta: [
    {
      name: "description",
      content:
        "Join the Hello Mail team. We're building the future of email infrastructure with a remote-first, developer-focused culture.",
    },
    {
      name: "og:title",
      content: "Careers - Hello Mail",
    },
    {
      name: "og:description",
      content:
        "Join the Hello Mail team. We're building the future of email infrastructure with a remote-first, developer-focused culture.",
    },
  ],
};
