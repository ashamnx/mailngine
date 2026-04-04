import { component$, Slot } from "@builder.io/qwik";

export interface CardProps {
  elevated?: boolean;
  class?: string;
}

export const Card = component$<CardProps>(
  ({ elevated = false, class: className }) => {
    const cardClass = [
      "hm-card",
      elevated ? "hm-card--elevated" : "",
      className,
    ]
      .filter(Boolean)
      .join(" ");

    return (
      <div class={cardClass}>
        <Slot />
      </div>
    );
  },
);
