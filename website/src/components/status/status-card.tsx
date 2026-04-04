import { component$ } from "@builder.io/qwik";
import { Card } from "~/components/ui";
import "./status-card.css";

export type ServiceStatus = "up" | "down" | "unknown";

export interface StatusCardProps {
  name: string;
  status: ServiceStatus;
  responseTime?: string;
}

export const StatusCard = component$<StatusCardProps>(
  ({ name, status, responseTime }) => {
    const statusLabel =
      status === "up"
        ? "Operational"
        : status === "down"
          ? "Outage"
          : "Unknown";

    return (
      <Card class="status-card">
        <div class="status-card__row">
          <div class="status-card__info">
            <span class={`status-card__dot status-card__dot--${status}`} />
            <span class="status-card__name">{name}</span>
          </div>
          <div class="status-card__meta">
            {responseTime && (
              <span class="status-card__response-time">{responseTime}</span>
            )}
            <span class={`status-card__label status-card__label--${status}`}>
              {statusLabel}
            </span>
          </div>
        </div>
      </Card>
    );
  },
);
