import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section, Badge, CodeBlock } from "~/components/ui";

type HttpMethod = "GET" | "POST" | "PATCH" | "DELETE";

interface Endpoint {
  method: HttpMethod;
  path: string;
  description: string;
}

interface EndpointGroup {
  id: string;
  name: string;
  endpoints: Endpoint[];
}

const API_GROUPS: EndpointGroup[] = [
  {
    id: "auth",
    name: "Auth",
    endpoints: [
      {
        method: "GET",
        path: "/v1/auth/google",
        description: "Initiate Google OAuth",
      },
      {
        method: "GET",
        path: "/v1/auth/google/callback",
        description: "OAuth callback",
      },
      { method: "POST", path: "/v1/auth/logout", description: "Logout" },
      {
        method: "GET",
        path: "/v1/auth/me",
        description: "Get current user",
      },
    ],
  },
  {
    id: "emails",
    name: "Emails",
    endpoints: [
      {
        method: "POST",
        path: "/v1/emails",
        description: "Send an email",
      },
      {
        method: "GET",
        path: "/v1/emails",
        description: "List sent emails",
      },
      {
        method: "GET",
        path: "/v1/emails/{id}",
        description: "Get email by ID",
      },
    ],
  },
  {
    id: "domains",
    name: "Domains",
    endpoints: [
      {
        method: "POST",
        path: "/v1/domains",
        description: "Add a domain",
      },
      {
        method: "GET",
        path: "/v1/domains",
        description: "List domains",
      },
      {
        method: "GET",
        path: "/v1/domains/{id}",
        description: "Get domain details",
      },
      {
        method: "PATCH",
        path: "/v1/domains/{id}",
        description: "Update domain settings",
      },
      {
        method: "DELETE",
        path: "/v1/domains/{id}",
        description: "Delete domain",
      },
      {
        method: "POST",
        path: "/v1/domains/{id}/verify",
        description: "Verify domain DNS",
      },
      {
        method: "POST",
        path: "/v1/domains/{id}/auto-dns",
        description: "Auto-configure DNS via Cloudflare",
      },
    ],
  },
  {
    id: "webhooks",
    name: "Webhooks",
    endpoints: [
      {
        method: "POST",
        path: "/v1/webhooks",
        description: "Create webhook",
      },
      {
        method: "GET",
        path: "/v1/webhooks",
        description: "List webhooks",
      },
      {
        method: "GET",
        path: "/v1/webhooks/{id}",
        description: "Get webhook",
      },
      {
        method: "PATCH",
        path: "/v1/webhooks/{id}",
        description: "Update webhook",
      },
      {
        method: "DELETE",
        path: "/v1/webhooks/{id}",
        description: "Delete webhook",
      },
      {
        method: "GET",
        path: "/v1/webhooks/{id}/deliveries",
        description: "List deliveries",
      },
    ],
  },
  {
    id: "templates",
    name: "Templates",
    endpoints: [
      {
        method: "POST",
        path: "/v1/templates",
        description: "Create template",
      },
      {
        method: "GET",
        path: "/v1/templates",
        description: "List templates",
      },
      {
        method: "GET",
        path: "/v1/templates/{id}",
        description: "Get template",
      },
      {
        method: "PATCH",
        path: "/v1/templates/{id}",
        description: "Update template",
      },
      {
        method: "DELETE",
        path: "/v1/templates/{id}",
        description: "Delete template",
      },
      {
        method: "POST",
        path: "/v1/templates/{id}/preview",
        description: "Preview template",
      },
    ],
  },
  {
    id: "api-keys",
    name: "API Keys",
    endpoints: [
      {
        method: "POST",
        path: "/v1/api-keys",
        description: "Create API key",
      },
      {
        method: "GET",
        path: "/v1/api-keys",
        description: "List API keys",
      },
      {
        method: "DELETE",
        path: "/v1/api-keys/{id}",
        description: "Revoke API key",
      },
    ],
  },
  {
    id: "suppressions",
    name: "Suppressions",
    endpoints: [
      {
        method: "GET",
        path: "/v1/suppressions",
        description: "List suppressions",
      },
      {
        method: "POST",
        path: "/v1/suppressions",
        description: "Add suppression",
      },
      {
        method: "DELETE",
        path: "/v1/suppressions/{id}",
        description: "Remove suppression",
      },
    ],
  },
  {
    id: "analytics",
    name: "Analytics",
    endpoints: [
      {
        method: "GET",
        path: "/v1/analytics/overview",
        description: "Analytics overview",
      },
      {
        method: "GET",
        path: "/v1/analytics/timeseries",
        description: "Timeseries data",
      },
      {
        method: "GET",
        path: "/v1/analytics/events",
        description: "Event breakdown",
      },
    ],
  },
  {
    id: "inbox",
    name: "Inbox",
    endpoints: [
      {
        method: "GET",
        path: "/v1/inbox/threads",
        description: "List threads",
      },
      {
        method: "GET",
        path: "/v1/inbox/threads/{id}",
        description: "Get thread",
      },
      {
        method: "DELETE",
        path: "/v1/inbox/threads/{id}",
        description: "Delete thread",
      },
      {
        method: "GET",
        path: "/v1/inbox/messages/{id}",
        description: "Get message",
      },
      {
        method: "PATCH",
        path: "/v1/inbox/messages/{id}",
        description: "Update message flags",
      },
      {
        method: "DELETE",
        path: "/v1/inbox/messages/{id}",
        description: "Delete message",
      },
      {
        method: "GET",
        path: "/v1/inbox/labels",
        description: "List labels",
      },
      {
        method: "POST",
        path: "/v1/inbox/labels",
        description: "Create label",
      },
      {
        method: "DELETE",
        path: "/v1/inbox/labels/{id}",
        description: "Delete label",
      },
      {
        method: "GET",
        path: "/v1/inbox/search",
        description: "Search messages",
      },
    ],
  },
  {
    id: "billing",
    name: "Billing",
    endpoints: [
      {
        method: "GET",
        path: "/v1/billing/usage",
        description: "Current usage",
      },
      {
        method: "GET",
        path: "/v1/billing/usage/history",
        description: "Usage history",
      },
      {
        method: "GET",
        path: "/v1/billing/plan",
        description: "Plan details",
      },
    ],
  },
  {
    id: "org",
    name: "Organization & Team",
    endpoints: [
      {
        method: "GET",
        path: "/v1/org",
        description: "Get organization",
      },
      {
        method: "PATCH",
        path: "/v1/org",
        description: "Update organization",
      },
      {
        method: "GET",
        path: "/v1/org/members",
        description: "List members",
      },
      {
        method: "POST",
        path: "/v1/org/members/invite",
        description: "Invite member",
      },
      {
        method: "PATCH",
        path: "/v1/org/members/{id}",
        description: "Update role",
      },
      {
        method: "DELETE",
        path: "/v1/org/members/{id}",
        description: "Remove member",
      },
    ],
  },
  {
    id: "audit",
    name: "Audit Logs",
    endpoints: [
      {
        method: "GET",
        path: "/v1/audit-logs",
        description: "List audit logs",
      },
      {
        method: "GET",
        path: "/v1/audit-logs/{id}",
        description: "Get audit log",
      },
    ],
  },
];

const METHOD_VARIANT: Record<HttpMethod, "success" | "info" | "warning" | "error"> = {
  GET: "success",
  POST: "info",
  PATCH: "warning",
  DELETE: "error",
};

const SEND_EMAIL_REQUEST = `curl -X POST https://api.mailngine.com/v1/emails \\
  -H "Authorization: Bearer mn_live_xxxxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "from": "hello@yourdomain.com",
    "to": ["user@example.com"],
    "subject": "Welcome to Mailngine",
    "html": "<h1>Welcome!</h1><p>Thanks for signing up.</p>"
  }'`;

const SEND_EMAIL_RESPONSE = `{
  "data": {
    "id": "em_abc123",
    "from": "hello@yourdomain.com",
    "to": ["user@example.com"],
    "subject": "Welcome to Mailngine",
    "status": "queued",
    "created_at": "2026-04-02T10:30:00Z"
  }
}`;

export default component$(() => {
  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="api-hero">
          <h1 class="api-hero__title">API Reference</h1>
          <p class="api-hero__subtitle">
            Complete REST API documentation for Mailngine. All endpoints
            require authentication via API key unless otherwise noted.
          </p>
          <p class="api-hero__base-url">
            Base URL: <code>https://api.mailngine.com</code>
          </p>
        </div>
      </Section>

      <Section padding="md">
        <div class="api-layout">
          {/* Sidebar / Table of Contents */}
          <aside class="api-sidebar">
            <nav class="api-toc" aria-label="API sections">
              <h4 class="api-toc__title">Resources</h4>
              <ul class="api-toc__list">
                {API_GROUPS.map((group) => (
                  <li key={group.id}>
                    <a href={`#${group.id}`} class="api-toc__link">
                      {group.name}
                    </a>
                  </li>
                ))}
              </ul>
            </nav>
          </aside>

          {/* Main Content */}
          <div class="api-main">
            {/* Send Email Example */}
            <div class="api-example" id="send-email-example">
              <h2>Quick Example: Send an Email</h2>
              <div class="api-example__blocks">
                <div>
                  <h4>Request</h4>
                  <CodeBlock code={SEND_EMAIL_REQUEST} language="bash" />
                </div>
                <div>
                  <h4>Response</h4>
                  <CodeBlock code={SEND_EMAIL_RESPONSE} language="json" />
                </div>
              </div>
            </div>

            {/* Endpoint Groups */}
            {API_GROUPS.map((group) => (
              <div key={group.id} class="api-group" id={group.id}>
                <h2 class="api-group__title">{group.name}</h2>
                <div class="api-group__endpoints">
                  {group.endpoints.map((ep) => (
                    <div
                      key={`${ep.method}-${ep.path}`}
                      class="api-endpoint"
                    >
                      <Badge variant={METHOD_VARIANT[ep.method]} class="api-endpoint__method">
                        {ep.method}
                      </Badge>
                      <code class="api-endpoint__path">{ep.path}</code>
                      <span class="api-endpoint__desc">{ep.description}</span>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .api-hero {
          text-align: center;
          max-width: 640px;
          margin: 0 auto;
        }
        .api-hero__title {
          font-size: var(--font-size-4xl);
          margin-bottom: var(--space-4);
        }
        .api-hero__subtitle {
          font-size: var(--font-size-lg);
          color: var(--color-text-secondary);
          line-height: 1.6;
          margin-bottom: var(--space-4);
        }
        .api-hero__base-url {
          font-size: var(--font-size-sm);
          color: var(--color-text-secondary);
        }
        .api-hero__base-url code {
          font-family: var(--font-mono);
          background-color: var(--color-bg-surface-alt);
          padding: var(--space-1) var(--space-2);
          border-radius: var(--radius-sm);
          color: var(--color-text-primary);
        }

        /* --- Two-Column Layout --- */
        .api-layout {
          display: grid;
          grid-template-columns: 220px 1fr;
          gap: var(--space-10);
          align-items: start;
        }
        @media (max-width: 900px) {
          .api-layout {
            grid-template-columns: 1fr;
          }
          .api-sidebar {
            display: none;
          }
        }

        /* --- Sidebar / TOC --- */
        .api-sidebar {
          position: sticky;
          top: calc(var(--nav-height) + var(--space-6));
        }
        .api-toc__title {
          font-size: var(--font-size-xs);
          text-transform: uppercase;
          letter-spacing: 0.08em;
          color: var(--color-text-disabled);
          margin-bottom: var(--space-3);
        }
        .api-toc__list {
          list-style: none;
          display: flex;
          flex-direction: column;
          gap: var(--space-1);
        }
        .api-toc__link {
          display: block;
          padding: var(--space-2) var(--space-3);
          font-size: var(--font-size-sm);
          color: var(--color-text-secondary);
          border-radius: var(--radius-sm);
          transition: background-color var(--transition-fast), color var(--transition-fast);
        }
        .api-toc__link:hover {
          background-color: var(--color-bg-surface);
          color: var(--color-primary);
        }

        /* --- Example Section --- */
        .api-example {
          margin-bottom: var(--space-12);
          padding-bottom: var(--space-10);
          border-bottom: 1px solid var(--color-border-light);
        }
        .api-example h2 {
          margin-bottom: var(--space-6);
        }
        .api-example h4 {
          margin-bottom: var(--space-3);
          font-size: var(--font-size-sm);
          color: var(--color-text-secondary);
          text-transform: uppercase;
          letter-spacing: 0.05em;
        }
        .api-example__blocks {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: var(--space-6);
        }
        @media (max-width: 768px) {
          .api-example__blocks {
            grid-template-columns: 1fr;
          }
        }

        /* --- Endpoint Groups --- */
        .api-group {
          margin-bottom: var(--space-10);
          scroll-margin-top: calc(var(--nav-height) + var(--space-6));
        }
        .api-group__title {
          font-size: var(--font-size-2xl);
          margin-bottom: var(--space-4);
          padding-bottom: var(--space-3);
          border-bottom: 1px solid var(--color-border-light);
        }
        .api-group__endpoints {
          display: flex;
          flex-direction: column;
          gap: var(--space-2);
        }

        /* --- Single Endpoint Row --- */
        .api-endpoint {
          display: flex;
          align-items: center;
          gap: var(--space-3);
          padding: var(--space-3) var(--space-4);
          border-radius: var(--radius-sm);
          transition: background-color var(--transition-fast);
        }
        .api-endpoint:hover {
          background-color: var(--color-bg-surface);
        }
        .api-endpoint__method {
          min-width: 64px;
          text-align: center;
          font-family: var(--font-mono);
          font-size: var(--font-size-xs);
          font-weight: var(--font-weight-semibold);
        }
        .api-endpoint__path {
          font-family: var(--font-mono);
          font-size: var(--font-size-sm);
          color: var(--color-text-primary);
          white-space: nowrap;
        }
        .api-endpoint__desc {
          font-size: var(--font-size-sm);
          color: var(--color-text-secondary);
          margin-left: auto;
          text-align: right;
        }
        @media (max-width: 600px) {
          .api-endpoint {
            flex-wrap: wrap;
          }
          .api-endpoint__desc {
            width: 100%;
            text-align: left;
            margin-left: calc(64px + var(--space-3));
            margin-top: var(--space-1);
          }
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "API Reference | Mailngine",
  meta: [
    {
      name: "description",
      content:
        "Complete REST API reference for Mailngine. Endpoints for emails, domains, webhooks, templates, API keys, suppressions, analytics, inbox, billing, and more.",
    },
  ],
};
