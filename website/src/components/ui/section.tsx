import { component$, Slot } from "@builder.io/qwik";

export type SectionBackground = "white" | "surface" | "gradient" | "primary";

export interface SectionProps {
  background?: SectionBackground;
  padding?: "sm" | "md" | "lg";
  class?: string;
  id?: string;
}

export const Section = component$<SectionProps>(
  ({ background = "white", padding = "md", id, class: className }) => {
    const sectionClass = [
      "hm-section",
      `hm-section--bg-${background}`,
      `hm-section--pad-${padding}`,
      className,
    ]
      .filter(Boolean)
      .join(" ");

    return (
      <section class={sectionClass} id={id}>
        <div class="hm-section__inner">
          <Slot />
        </div>
      </section>
    );
  },
);
