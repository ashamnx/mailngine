import { component$ } from "@builder.io/qwik";
import { Card, Badge, Button } from "~/components/ui";

export interface PlanCardProps {
  name: string;
  price: string;
  period?: string;
  features: string[];
  highlighted?: boolean;
  ctaText: string;
  ctaHref: string;
}

export const PlanCard = component$<PlanCardProps>(
  ({
    name,
    price,
    period,
    features,
    highlighted = false,
    ctaText,
    ctaHref,
  }) => {
    return (
      <Card
        class={
          highlighted ? "plan-card plan-card--highlighted" : "plan-card"
        }
      >
        <div class="plan-card__header">
          <div class="plan-card__name-row">
            <h3 class="plan-card__name">{name}</h3>
            {highlighted && <Badge variant="info">Recommended</Badge>}
          </div>
          <div class="plan-card__price">
            <span class="plan-card__amount">{price}</span>
            {period && <span class="plan-card__period">/{period}</span>}
          </div>
        </div>

        <ul class="plan-card__features">
          {features.map((feature) => (
            <li key={feature} class="plan-card__feature">
              <span
                class="material-symbols-rounded plan-card__check"
                aria-hidden="true"
              >
                check
              </span>
              {feature}
            </li>
          ))}
        </ul>

        <div class="plan-card__cta">
          <Button
            variant={highlighted ? "primary" : "secondary"}
            href={ctaHref}
            size="md"
            class="plan-card__button"
          >
            {ctaText}
          </Button>
        </div>
      </Card>
    );
  },
);
