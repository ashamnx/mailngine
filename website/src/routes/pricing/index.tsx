import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section } from "~/components/ui";
import { PlanCard } from "~/components/pricing/plan-card";
import { FeatureMatrix } from "~/components/pricing/feature-matrix";
import { PricingFaq } from "~/components/pricing/faq";

const plans = [
  {
    name: "Free",
    price: "$0",
    period: "mo",
    features: [
      "100 emails/day",
      "1 domain",
      "Community support",
      "Basic analytics",
    ],
    highlighted: false,
    ctaText: "Get Started",
    ctaHref: "/get-started",
  },
  {
    name: "Starter",
    price: "$20",
    period: "mo",
    features: [
      "10,000 emails/mo",
      "5 domains",
      "Email support",
      "Full analytics",
      "Webhooks",
    ],
    highlighted: false,
    ctaText: "Start Free Trial",
    ctaHref: "/get-started?plan=starter",
  },
  {
    name: "Pro",
    price: "$80",
    period: "mo",
    features: [
      "100,000 emails/mo",
      "Unlimited domains",
      "Priority support",
      "Full analytics",
      "Webhooks",
      "Templates",
      "Team management",
    ],
    highlighted: true,
    ctaText: "Start Free Trial",
    ctaHref: "/get-started?plan=pro",
  },
  {
    name: "Enterprise",
    price: "Custom",
    period: undefined,
    features: [
      "Unlimited emails",
      "Dedicated IPs",
      "SLA",
      "Custom integrations",
      "Dedicated support",
    ],
    highlighted: false,
    ctaText: "Contact Sales",
    ctaHref: "/contact",
  },
];

export default component$(() => {
  return (
    <>
      {/* Header */}
      <Section padding="lg" class="pricing-hero">
        <div class="pricing-hero__inner">
          <h1 class="pricing-hero__title">Simple, transparent pricing</h1>
          <p class="pricing-hero__subtitle">
            Start free, scale as you grow. No hidden fees, no surprises.
          </p>
        </div>
      </Section>

      {/* Plan Cards */}
      <Section padding="md">
        <div class="pricing-grid">
          {plans.map((plan) => (
            <PlanCard
              key={plan.name}
              name={plan.name}
              price={plan.price}
              period={plan.period}
              features={plan.features}
              highlighted={plan.highlighted}
              ctaText={plan.ctaText}
              ctaHref={plan.ctaHref}
            />
          ))}
        </div>
      </Section>

      {/* Feature Matrix */}
      <Section background="surface" padding="lg">
        <FeatureMatrix />
      </Section>

      {/* FAQ */}
      <Section padding="lg">
        <PricingFaq />
      </Section>

      {/* CTA */}
      <Section background="primary" padding="lg">
        <div class="cta-block">
          <h2 class="cta-block__heading">
            Ready to start sending?
          </h2>
          <p class="cta-block__desc">
            Get up and running in under 5 minutes with our free plan.
          </p>
        </div>
      </Section>
    </>
  );
});

export const head: DocumentHead = {
  title: "Pricing - Hello Mail",
  meta: [
    {
      name: "description",
      content:
        "Simple, transparent pricing for Hello Mail. Start free with 100 emails/day. Scale to millions with Pro and Enterprise plans.",
    },
    {
      property: "og:title",
      content: "Pricing - Hello Mail",
    },
    {
      property: "og:description",
      content:
        "Simple, transparent pricing for email infrastructure. Start free, scale as you grow.",
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
