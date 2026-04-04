export type ChangeType = "added" | "improved" | "fixed";

export interface Change {
  type: ChangeType;
  text: string;
}

export interface Release {
  version: string;
  date: string;
  title: string;
  changes: Change[];
}

export const releases: Release[] = [
  {
    version: "1.3.0",
    date: "2026-04-03",
    title: "Templates & Team Management",
    changes: [
      {
        type: "added",
        text: "Email templates with variable rendering and preview",
      },
      { type: "added", text: "Team management with role-based access" },
      { type: "added", text: "Billing dashboard with plan comparison" },
      {
        type: "improved",
        text: "Settings page with notification preferences",
      },
    ],
  },
  {
    version: "1.2.0",
    date: "2026-04-02",
    title: "Webhooks & Analytics",
    changes: [
      {
        type: "added",
        text: "Webhook delivery with HMAC-SHA256 signing",
      },
      {
        type: "added",
        text: "Analytics dashboard with KPI cards and timeseries",
      },
      { type: "added", text: "Open and click tracking endpoints" },
      { type: "added", text: "Audit logging for all mutations" },
    ],
  },
  {
    version: "1.1.0",
    date: "2026-04-01",
    title: "Inbox & Suppression",
    changes: [
      { type: "added", text: "Gmail-like inbox with 3-panel layout" },
      { type: "added", text: "Email threading with JWZ algorithm" },
      {
        type: "added",
        text: "Suppression list with Valkey cache layer",
      },
      { type: "added", text: "Bounce and FBL processing" },
    ],
  },
  {
    version: "1.0.0",
    date: "2026-03-30",
    title: "Initial Release",
    changes: [
      { type: "added", text: "Email sending via REST API and SMTP" },
      {
        type: "added",
        text: "Domain management with DKIM key generation",
      },
      {
        type: "added",
        text: "DNS verification and Cloudflare auto-DNS",
      },
      { type: "added", text: "Google OAuth authentication" },
      { type: "added", text: "Go, Node.js, and Laravel SDKs" },
    ],
  },
];
