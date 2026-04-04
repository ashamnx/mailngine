import { component$, type QRL, Slot } from "@builder.io/qwik";

export type ButtonVariant = "primary" | "secondary" | "text";
export type ButtonSize = "sm" | "md" | "lg";

export interface ButtonProps {
  variant?: ButtonVariant;
  size?: ButtonSize;
  href?: string;
  type?: "button" | "submit" | "reset";
  onClick$?: QRL<() => void>;
  disabled?: boolean;
  class?: string;
}

export const Button = component$<ButtonProps>(
  ({
    variant = "primary",
    size = "md",
    href,
    type = "button",
    disabled = false,
    class: className,
  }) => {
    const baseClass = [
      "hm-button",
      `hm-button--${variant}`,
      `hm-button--${size}`,
      className,
    ]
      .filter(Boolean)
      .join(" ");

    if (href) {
      return (
        <a href={href} class={baseClass}>
          <Slot />
        </a>
      );
    }

    return (
      <button type={type} class={baseClass} disabled={disabled}>
        <Slot />
      </button>
    );
  },
);

/* --- Styles scoped via class names for this component --- */
/* Included as a <style> import or global; using CSS below. */
