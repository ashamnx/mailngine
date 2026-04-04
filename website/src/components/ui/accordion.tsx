import { component$, Slot } from "@builder.io/qwik";

export interface AccordionItemProps {
  title: string;
  open?: boolean;
}

/**
 * Uses native <details>/<summary> for zero-JS accordion behavior.
 * Fully accessible and works with SSR without client-side hydration cost.
 */
export const AccordionItem = component$<AccordionItemProps>(
  ({ title, open = false }) => {
    return (
      <details class="hm-accordion" open={open}>
        <summary class="hm-accordion__summary">
          <span class="hm-accordion__title">{title}</span>
          <span class="hm-accordion__icon" aria-hidden="true" />
        </summary>
        <div class="hm-accordion__content">
          <Slot />
        </div>
      </details>
    );
  },
);
