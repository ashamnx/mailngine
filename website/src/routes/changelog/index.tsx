import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section } from "~/components/ui/section";
import { Badge } from "~/components/ui/badge";
import type { BadgeVariant } from "~/components/ui/badge";
import { releases } from "~/content/changelog";
import type { ChangeType } from "~/content/changelog";

const changeTypeBadge: Record<ChangeType, { variant: BadgeVariant; label: string }> = {
  added: { variant: "success", label: "Added" },
  improved: { variant: "info", label: "Improved" },
  fixed: { variant: "warning", label: "Fixed" },
};

export default component$(() => {
  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="changelog-hero">
          <h1>Changelog</h1>
          <p>
            A history of improvements, new features, and fixes shipped to Hello
            Mail.
          </p>
        </div>
      </Section>

      <Section padding="lg">
        <div class="changelog-timeline">
          {releases.map((release) => (
            <div key={release.version} class="changelog-entry">
              <div class="changelog-entry__marker" aria-hidden="true" />
              <div class="changelog-entry__content">
                <div class="changelog-entry__header">
                  <Badge variant="info">{release.version}</Badge>
                  <time class="changelog-entry__date">{release.date}</time>
                </div>
                <h2 class="changelog-entry__title">{release.title}</h2>
                <ul class="changelog-entry__changes">
                  {release.changes.map((change) => {
                    const badge = changeTypeBadge[change.type];
                    return (
                      <li key={change.text} class="changelog-change">
                        <Badge variant={badge.variant}>{badge.label}</Badge>
                        <span>{change.text}</span>
                      </li>
                    );
                  })}
                </ul>
              </div>
            </div>
          ))}
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .changelog-hero {
          text-align: center;
          max-width: 600px;
          margin: 0 auto;
        }
        .changelog-hero h1 {
          margin-bottom: var(--space-4);
        }
        .changelog-hero p {
          font-size: var(--font-size-lg);
        }

        .changelog-timeline {
          max-width: 720px;
          margin: 0 auto;
          position: relative;
          padding-left: var(--space-8);
        }
        .changelog-timeline::before {
          content: "";
          position: absolute;
          left: 7px;
          top: 0;
          bottom: 0;
          width: 2px;
          background-color: var(--color-border-light);
        }

        .changelog-entry {
          position: relative;
          padding-bottom: var(--space-12);
        }
        .changelog-entry:last-child {
          padding-bottom: 0;
        }

        .changelog-entry__marker {
          position: absolute;
          left: calc(-1 * var(--space-8) + 3px);
          top: 6px;
          width: 10px;
          height: 10px;
          border-radius: 50%;
          background-color: var(--color-primary);
          border: 2px solid var(--color-bg-app);
          z-index: 1;
        }

        .changelog-entry__content {
          padding-left: var(--space-2);
        }

        .changelog-entry__header {
          display: flex;
          align-items: center;
          gap: var(--space-3);
          margin-bottom: var(--space-3);
        }
        .changelog-entry__date {
          font-size: var(--font-size-sm);
          color: var(--color-text-secondary);
        }
        .changelog-entry__title {
          font-size: var(--font-size-xl);
          margin-bottom: var(--space-4);
        }

        .changelog-entry__changes {
          list-style: none;
          display: flex;
          flex-direction: column;
          gap: var(--space-3);
        }
        .changelog-change {
          display: flex;
          align-items: baseline;
          gap: var(--space-3);
          font-size: var(--font-size-sm);
          color: var(--color-text-secondary);
          line-height: 1.5;
        }
        .changelog-change .hm-badge {
          flex-shrink: 0;
        }

        @media (max-width: 480px) {
          .changelog-timeline {
            padding-left: var(--space-6);
          }
          .changelog-entry__marker {
            left: calc(-1 * var(--space-6) + 3px);
          }
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "Changelog - Hello Mail",
  meta: [
    {
      name: "description",
      content:
        "See what's new in Hello Mail. A history of features, improvements, and fixes.",
    },
    {
      name: "og:title",
      content: "Changelog - Hello Mail",
    },
    {
      name: "og:description",
      content:
        "See what's new in Hello Mail. A history of features, improvements, and fixes.",
    },
  ],
};
