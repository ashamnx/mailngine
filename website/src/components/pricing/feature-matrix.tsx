import { component$ } from "@builder.io/qwik";

interface FeatureRow {
  name: string;
  free: string;
  starter: string;
  pro: string;
  enterprise: string;
}

const featureRows: FeatureRow[] = [
  {
    name: "Email sending",
    free: "100/day",
    starter: "10,000/mo",
    pro: "100,000/mo",
    enterprise: "Unlimited",
  },
  {
    name: "Inbound email",
    free: "\u2014",
    starter: "\u2714",
    pro: "\u2714",
    enterprise: "\u2714",
  },
  {
    name: "Domains",
    free: "1",
    starter: "5",
    pro: "Unlimited",
    enterprise: "Unlimited",
  },
  {
    name: "Analytics",
    free: "Basic",
    starter: "Full",
    pro: "Full",
    enterprise: "Full",
  },
  {
    name: "Webhooks",
    free: "\u2014",
    starter: "\u2714",
    pro: "\u2714",
    enterprise: "\u2714",
  },
  {
    name: "Templates",
    free: "\u2014",
    starter: "\u2014",
    pro: "\u2714",
    enterprise: "\u2714",
  },
  {
    name: "Team members",
    free: "1",
    starter: "3",
    pro: "10",
    enterprise: "Unlimited",
  },
  {
    name: "API keys",
    free: "1",
    starter: "5",
    pro: "20",
    enterprise: "Unlimited",
  },
  {
    name: "Suppression lists",
    free: "\u2014",
    starter: "\u2714",
    pro: "\u2714",
    enterprise: "\u2714",
  },
  {
    name: "Audit logs",
    free: "\u2014",
    starter: "\u2014",
    pro: "\u2714",
    enterprise: "\u2714",
  },
  {
    name: "DKIM signing",
    free: "\u2714",
    starter: "\u2714",
    pro: "\u2714",
    enterprise: "\u2714",
  },
  {
    name: "Custom DNS",
    free: "\u2014",
    starter: "\u2714",
    pro: "\u2714",
    enterprise: "\u2714",
  },
  {
    name: "Dedicated IPs",
    free: "\u2014",
    starter: "\u2014",
    pro: "\u2014",
    enterprise: "\u2714",
  },
  {
    name: "SLA",
    free: "\u2014",
    starter: "\u2014",
    pro: "\u2014",
    enterprise: "\u2714",
  },
];

const plans = ["Free", "Starter", "Pro", "Enterprise"] as const;
const planKeys = ["free", "starter", "pro", "enterprise"] as const;

export const FeatureMatrix = component$(() => {
  return (
    <div class="feature-matrix-wrapper">
      <h2 class="feature-matrix__heading">Compare plans</h2>
      <div class="feature-matrix__scroll">
        <table class="feature-matrix">
          <thead>
            <tr>
              <th class="feature-matrix__feature-col">Feature</th>
              {plans.map((plan) => (
                <th key={plan} class="feature-matrix__plan-col">
                  {plan}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {featureRows.map((row) => (
              <tr key={row.name}>
                <td class="feature-matrix__feature-name">{row.name}</td>
                {planKeys.map((key) => {
                  const value = row[key];
                  const isCheck = value === "\u2714";
                  const isDash = value === "\u2014";
                  return (
                    <td
                      key={key}
                      class={[
                        "feature-matrix__cell",
                        isCheck ? "feature-matrix__cell--check" : "",
                        isDash ? "feature-matrix__cell--dash" : "",
                      ]
                        .filter(Boolean)
                        .join(" ")}
                    >
                      {value}
                    </td>
                  );
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
});
