import { component$ } from "@builder.io/qwik";
import { Card } from "~/components/ui";

const features = [
  {
    icon: "rocket_launch",
    title: "Reliable Delivery",
    description:
      "Built on Postfix MTA with automatic retry, bounce handling, and IP warmup for maximum deliverability.",
  },
  {
    icon: "dns",
    title: "Domain Management",
    description:
      "Auto-verify DNS records, generate DKIM keys, and configure SPF/DMARC with one-click Cloudflare integration.",
  },
  {
    icon: "analytics",
    title: "Real-time Analytics",
    description:
      "Track opens, clicks, bounces, and deliveries with detailed timeseries dashboards and event streams.",
  },
  {
    icon: "webhook",
    title: "Webhooks",
    description:
      "Get instant notifications for email events. HMAC-signed payloads with automatic retry and delivery tracking.",
  },
  {
    icon: "inbox",
    title: "Inbox Support",
    description:
      "Receive and manage inbound emails with threading, labels, search, and a full Gmail-like inbox UI.",
  },
  {
    icon: "code",
    title: "Multi-SDK Support",
    description:
      "Official SDKs for Go, Node.js, and Laravel. Send your first email in under 5 minutes.",
  },
];

export const Features = component$(() => {
  return (
    <div class="features-grid">
      {features.map((feature) => (
        <Card key={feature.icon}>
          <span
            class="material-symbols-rounded features-grid__icon"
            aria-hidden="true"
          >
            {feature.icon}
          </span>
          <h3 class="features-grid__title">{feature.title}</h3>
          <p class="features-grid__desc">{feature.description}</p>
        </Card>
      ))}
    </div>
  );
});
