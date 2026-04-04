import { component$, Slot } from "@builder.io/qwik";

export type BadgeVariant =
  | "success"
  | "warning"
  | "error"
  | "info"
  | "neutral";

export interface BadgeProps {
  variant?: BadgeVariant;
  class?: string;
}

export const Badge = component$<BadgeProps>(
  ({ variant = "neutral", class: className }) => {
    const badgeClass = ["hm-badge", `hm-badge--${variant}`, className]
      .filter(Boolean)
      .join(" ");

    return (
      <span class={badgeClass}>
        <Slot />
      </span>
    );
  },
);
