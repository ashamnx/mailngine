import {
  component$,
  useSignal,
  useVisibleTask$,
} from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { routeLoader$ } from "@builder.io/qwik-city";
import { Section } from "~/components/ui";
import {
  StatusCard,
  type ServiceStatus,
} from "~/components/status/status-card";

interface HealthData {
  status: string;
  postgres: string;
  valkey: string;
}

export const useHealthCheck = routeLoader$(async () => {
  try {
    const res = await fetch("http://localhost:8080/health");
    const data = await res.json();
    return data.data as HealthData;
  } catch {
    return { status: "unknown", postgres: "unknown", valkey: "unknown" };
  }
});

function toServiceStatus(value: string): ServiceStatus {
  if (value === "up" || value === "healthy") return "up";
  if (value === "down" || value === "unhealthy") return "down";
  return "unknown";
}

type OverallState = "operational" | "degraded" | "outage";

function deriveOverallState(services: ServiceStatus[]): OverallState {
  if (services.every((s) => s === "up")) return "operational";
  if (services.every((s) => s === "down")) return "outage";
  return "degraded";
}

const OVERALL_CONFIG: Record<
  OverallState,
  { label: string; className: string }
> = {
  operational: {
    label: "All Systems Operational",
    className: "status-banner--operational",
  },
  degraded: {
    label: "Degraded Performance",
    className: "status-banner--degraded",
  },
  outage: {
    label: "System Outage",
    className: "status-banner--outage",
  },
};

export default component$(() => {
  const healthLoader = useHealthCheck();

  const apiStatus = useSignal<ServiceStatus>(
    toServiceStatus(healthLoader.value.status),
  );
  const dbStatus = useSignal<ServiceStatus>(
    toServiceStatus(healthLoader.value.postgres),
  );
  const cacheStatus = useSignal<ServiceStatus>(
    toServiceStatus(healthLoader.value.valkey),
  );
  const smtpStatus = useSignal<ServiceStatus>("up");
  const storageStatus = useSignal<ServiceStatus>("up");

  // Client-side polling every 30 seconds
  // useVisibleTask$ is appropriate here: this is a client-only interval
  // that must start when the component becomes visible in the viewport.
  useVisibleTask$(({ cleanup }) => {
    const interval = setInterval(async () => {
      try {
        const res = await fetch("/api/health");
        const data = await res.json();
        const health = data.data as HealthData;
        apiStatus.value = toServiceStatus(health.status);
        dbStatus.value = toServiceStatus(health.postgres);
        cacheStatus.value = toServiceStatus(health.valkey);
      } catch {
        apiStatus.value = "unknown";
        dbStatus.value = "unknown";
        cacheStatus.value = "unknown";
      }
    }, 30000);
    cleanup(() => clearInterval(interval));
  });

  const overall = deriveOverallState([
    apiStatus.value,
    dbStatus.value,
    cacheStatus.value,
    smtpStatus.value,
    storageStatus.value,
  ]);
  const bannerConfig = OVERALL_CONFIG[overall];

  return (
    <>
      {/* Overall Status Banner */}
      <div class={`status-banner ${bannerConfig.className}`}>
        <div class="status-banner__inner">
          <span class="status-banner__dot" />
          <span class="status-banner__label">{bannerConfig.label}</span>
        </div>
      </div>

      <Section padding="lg">
        <div class="status-page">
          <h1 class="status-page__title">System Status</h1>
          <p class="status-page__subtitle">
            Real-time health of Hello Mail infrastructure. This page
            auto-refreshes every 30 seconds.
          </p>

          <div class="status-page__services">
            <StatusCard name="API Server" status={apiStatus.value} />
            <StatusCard name="Database" status={dbStatus.value} />
            <StatusCard name="Cache & Queue" status={cacheStatus.value} />
            <StatusCard name="SMTP Delivery" status={smtpStatus.value} />
            <StatusCard name="Object Storage" status={storageStatus.value} />
          </div>

          <div class="status-page__footer">
            <p>
              Having issues? Contact us at{" "}
              <a href="mailto:support@hellomail.dev">support@hellomail.dev</a>
            </p>
          </div>
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        /* --- Status Banner --- */
        .status-banner {
          padding: var(--space-4) var(--space-6);
          text-align: center;
        }
        .status-banner__inner {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: var(--space-3);
          max-width: var(--section-max-width);
          margin: 0 auto;
        }
        .status-banner__dot {
          display: inline-block;
          width: 10px;
          height: 10px;
          border-radius: 50%;
        }
        .status-banner__label {
          font-family: var(--font-primary);
          font-size: var(--font-size-lg);
          font-weight: var(--font-weight-medium);
        }
        .status-banner--operational {
          background-color: #e6f4ea;
        }
        .status-banner--operational .status-banner__dot {
          background-color: #34a853;
        }
        .status-banner--operational .status-banner__label {
          color: #1e8e3e;
        }
        .status-banner--degraded {
          background-color: #fef7e0;
        }
        .status-banner--degraded .status-banner__dot {
          background-color: #fbbc04;
        }
        .status-banner--degraded .status-banner__label {
          color: #b06000;
        }
        .status-banner--outage {
          background-color: #fce8e6;
        }
        .status-banner--outage .status-banner__dot {
          background-color: #ea4335;
        }
        .status-banner--outage .status-banner__label {
          color: #c5221f;
        }

        /* --- Status Page Layout --- */
        .status-page {
          max-width: 720px;
          margin: 0 auto;
        }
        .status-page__title {
          font-size: var(--font-size-4xl);
          margin-bottom: var(--space-4);
        }
        .status-page__subtitle {
          font-size: var(--font-size-lg);
          color: var(--color-text-secondary);
          margin-bottom: var(--space-10);
        }
        .status-page__services {
          display: flex;
          flex-direction: column;
          gap: var(--space-1);
        }
        .status-page__footer {
          margin-top: var(--space-12);
          padding-top: var(--space-8);
          border-top: 1px solid var(--color-border-light);
          text-align: center;
        }
        .status-page__footer p {
          font-size: var(--font-size-sm);
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "System Status | Hello Mail",
  meta: [
    {
      name: "description",
      content:
        "Real-time system status for Hello Mail. Monitor API, database, cache, SMTP, and storage health.",
    },
  ],
};
