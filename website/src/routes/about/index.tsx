import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section } from "~/components/ui/section";
import { Card } from "~/components/ui/card";

const values = [
  {
    icon: "verified",
    title: "Reliability",
    description: "Your emails deserve 99.9% uptime",
  },
  {
    icon: "code",
    title: "Developer Experience",
    description: "APIs that feel natural in every language",
  },
  {
    icon: "visibility",
    title: "Transparency",
    description: "Open about our infrastructure and pricing",
  },
  {
    icon: "lock",
    title: "Privacy",
    description: "Your data stays yours. Always.",
  },
];

const team = [
  { name: "Alex Chen", role: "Founder & CEO" },
  { name: "Sarah Kim", role: "Head of Engineering" },
  { name: "James Patel", role: "Infrastructure Lead" },
  { name: "Maria Gonzalez", role: "Developer Experience" },
];

export default component$(() => {
  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="about-hero">
          <h1>Built for developers who care about email</h1>
          <p>
            Mailngine started with a simple frustration: email infrastructure
            should not be this hard. We are a small, focused team building the
            email platform we always wished existed -- reliable, transparent, and
            designed for developers first.
          </p>
        </div>
      </Section>

      <Section padding="lg">
        <div class="about-values">
          <h2 class="about-section-title">Our Values</h2>
          <div class="about-values__grid">
            {values.map((value) => (
              <Card key={value.title} elevated>
                <div class="about-value-card">
                  <span class="material-symbols-rounded about-value-card__icon">
                    {value.icon}
                  </span>
                  <h3>{value.title}</h3>
                  <p>{value.description}</p>
                </div>
              </Card>
            ))}
          </div>
        </div>
      </Section>

      <Section background="surface" padding="lg">
        <div class="about-team">
          <h2 class="about-section-title">Meet the Team</h2>
          <p class="about-team__intro">
            A small crew with big ambitions. We have built infrastructure at
            scale before, and now we are channeling that experience into
            Mailngine.
          </p>
          <div class="about-team__grid">
            {team.map((member) => (
              <div key={member.name} class="about-team__member">
                <div class="about-team__avatar" aria-hidden="true">
                  <span class="material-symbols-rounded">person</span>
                </div>
                <h3>{member.name}</h3>
                <p>{member.role}</p>
              </div>
            ))}
          </div>
        </div>
      </Section>

      <Section padding="lg">
        <div class="about-oss">
          <span class="material-symbols-rounded about-oss__icon">
            open_in_new
          </span>
          <h2>Built with open source</h2>
          <p>
            Mailngine is built on open source tools we trust: Postfix for SMTP
            delivery, OpenDKIM for message signing, Valkey for caching, and
            PostgreSQL for data storage. We believe in giving back to the
            communities that make our work possible, and we open source our SDKs
            and client libraries.
          </p>
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .about-hero {
          text-align: center;
          max-width: 720px;
          margin: 0 auto;
        }
        .about-hero h1 {
          font-size: var(--font-size-4xl);
          margin-bottom: var(--space-6);
        }
        .about-hero p {
          font-size: var(--font-size-lg);
          line-height: 1.7;
        }

        .about-section-title {
          text-align: center;
          margin-bottom: var(--space-10);
        }

        .about-values__grid {
          display: grid;
          grid-template-columns: repeat(4, 1fr);
          gap: var(--space-6);
        }
        .about-value-card {
          text-align: center;
          padding: var(--space-4) 0;
        }
        .about-value-card__icon {
          font-size: 40px;
          color: var(--color-primary);
          margin-bottom: var(--space-4);
        }
        .about-value-card h3 {
          font-size: var(--font-size-lg);
          margin-bottom: var(--space-2);
        }
        .about-value-card p {
          font-size: var(--font-size-sm);
        }

        .about-team {
          text-align: center;
        }
        .about-team__intro {
          max-width: 600px;
          margin: 0 auto var(--space-10);
          font-size: var(--font-size-lg);
        }
        .about-team__grid {
          display: grid;
          grid-template-columns: repeat(4, 1fr);
          gap: var(--space-8);
        }
        .about-team__member {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: var(--space-2);
        }
        .about-team__avatar {
          width: 80px;
          height: 80px;
          border-radius: 50%;
          background-color: var(--color-primary-light);
          display: flex;
          align-items: center;
          justify-content: center;
          margin-bottom: var(--space-2);
        }
        .about-team__avatar .material-symbols-rounded {
          font-size: 36px;
          color: var(--color-primary);
        }
        .about-team__member h3 {
          font-size: var(--font-size-base);
        }
        .about-team__member p {
          font-size: var(--font-size-sm);
        }

        .about-oss {
          text-align: center;
          max-width: 640px;
          margin: 0 auto;
        }
        .about-oss__icon {
          font-size: 40px;
          color: var(--color-primary);
          margin-bottom: var(--space-4);
        }
        .about-oss h2 {
          margin-bottom: var(--space-4);
        }
        .about-oss p {
          font-size: var(--font-size-lg);
          line-height: 1.7;
        }

        @media (max-width: 768px) {
          .about-hero h1 {
            font-size: var(--font-size-3xl);
          }
          .about-values__grid {
            grid-template-columns: repeat(2, 1fr);
          }
          .about-team__grid {
            grid-template-columns: repeat(2, 1fr);
          }
        }
        @media (max-width: 480px) {
          .about-values__grid,
          .about-team__grid {
            grid-template-columns: 1fr;
          }
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "About - Mailngine",
  meta: [
    {
      name: "description",
      content:
        "Learn about Mailngine, the email infrastructure platform built for developers who care about reliability, transparency, and great developer experience.",
    },
    {
      property: "og:title",
      content: "About - Mailngine",
    },
    {
      property: "og:description",
      content:
        "Learn about Mailngine, the email infrastructure platform built for developers who care about reliability, transparency, and great developer experience.",
    },
  ],
};
