import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section } from "~/components/ui/section";

export default component$(() => {
  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="legal-hero">
          <h1>Terms of Service</h1>
          <p>Effective date: April 1, 2026</p>
        </div>
      </Section>

      <Section padding="lg">
        <div class="legal-content">
          <nav class="legal-toc">
            <h2>Table of Contents</h2>
            <ol>
              <li>
                <a href="#acceptance">Acceptance of Terms</a>
              </li>
              <li>
                <a href="#description">Description of Service</a>
              </li>
              <li>
                <a href="#registration">Account Registration</a>
              </li>
              <li>
                <a href="#acceptable-use">Acceptable Use Policy</a>
              </li>
              <li>
                <a href="#api-usage">API Usage and Rate Limits</a>
              </li>
              <li>
                <a href="#payment">Payment Terms</a>
              </li>
              <li>
                <a href="#data-ownership">Data Ownership</a>
              </li>
              <li>
                <a href="#liability">Limitation of Liability</a>
              </li>
              <li>
                <a href="#termination">Termination</a>
              </li>
              <li>
                <a href="#modifications">Modifications to Terms</a>
              </li>
              <li>
                <a href="#governing-law">Governing Law</a>
              </li>
              <li>
                <a href="#contact">Contact Information</a>
              </li>
            </ol>
          </nav>

          <article class="legal-body">
            <section id="acceptance">
              <h2>1. Acceptance of Terms</h2>
              <p>
                By accessing or using the Mailngine service ("Service"), you
                agree to be bound by these Terms of Service ("Terms"). If you do
                not agree to these Terms, you may not use the Service. These
                Terms apply to all users, including individuals and organizations
                that create accounts or access the Service through our API.
              </p>
            </section>

            <section id="description">
              <h2>2. Description of Service</h2>
              <p>
                Mailngine provides email infrastructure services including
                transactional email delivery, domain management, email
                analytics, webhook notifications, and related developer tools.
                The Service is offered via a web dashboard, REST API, SMTP
                relay, and client SDKs for various programming languages.
              </p>
              <p>
                We reserve the right to modify, suspend, or discontinue any
                part of the Service at any time with reasonable notice. We will
                make commercially reasonable efforts to notify you of material
                changes via email or through the dashboard.
              </p>
            </section>

            <section id="registration">
              <h2>3. Account Registration</h2>
              <p>
                To use the Service, you must create an account by providing
                accurate and complete information. You are responsible for
                maintaining the security of your account credentials, including
                API keys. You must notify us immediately of any unauthorized
                access to your account.
              </p>
              <p>
                You may not create multiple accounts to circumvent rate limits,
                usage restrictions, or billing terms. Accounts created by
                automated methods are not permitted.
              </p>
            </section>

            <section id="acceptable-use">
              <h2>4. Acceptable Use Policy</h2>
              <p>You agree not to use the Service to:</p>
              <ul>
                <li>
                  Send unsolicited bulk email (spam) or messages to recipients
                  who have not opted in to receive communications from you.
                </li>
                <li>
                  Send phishing emails or messages designed to deceive
                  recipients into revealing personal information or credentials.
                </li>
                <li>
                  Distribute malware, viruses, or any other harmful software
                  through email content or attachments.
                </li>
                <li>
                  Violate any applicable laws or regulations, including but not
                  limited to the CAN-SPAM Act, GDPR, CASL, and other anti-spam
                  and data protection legislation.
                </li>
                <li>
                  Impersonate any person or entity, or misrepresent your
                  affiliation with any person or entity.
                </li>
                <li>
                  Harvest or collect email addresses or other personal
                  information from third parties without their consent.
                </li>
              </ul>
              <p>
                Violations of this Acceptable Use Policy may result in
                immediate suspension or termination of your account without
                notice. We reserve the right to investigate suspected
                violations and cooperate with law enforcement where required.
              </p>
            </section>

            <section id="api-usage">
              <h2>5. API Usage and Rate Limits</h2>
              <p>
                The Service imposes rate limits on API requests to ensure fair
                usage and maintain service quality for all users. Current rate
                limits are published in our API documentation and may be
                adjusted at any time.
              </p>
              <p>
                You agree not to circumvent or attempt to circumvent rate
                limits through any means, including distributing requests
                across multiple API keys or accounts. Excessive API usage that
                degrades the Service for other users may result in temporary
                throttling or suspension.
              </p>
            </section>

            <section id="payment">
              <h2>6. Payment Terms</h2>
              <p>
                Certain features of the Service require a paid subscription.
                Fees are billed in advance on a monthly or annual basis as
                selected during plan enrollment. All fees are non-refundable
                except as required by law or as explicitly stated in these
                Terms.
              </p>
              <p>
                We may change pricing at any time with at least 30 days' prior
                notice. Price changes will take effect at the start of your
                next billing cycle following the notice period. If you do not
                agree with a price change, you may cancel your subscription
                before the new price takes effect.
              </p>
            </section>

            <section id="data-ownership">
              <h2>7. Data Ownership</h2>
              <p>
                You retain all ownership rights to the data you submit to the
                Service, including email content, recipient lists, templates,
                and domain configurations ("Your Data"). We do not claim any
                ownership over Your Data.
              </p>
              <p>
                You grant Mailngine a limited, non-exclusive license to
                process Your Data solely for the purpose of providing the
                Service. We will not access, use, or share Your Data for any
                other purpose, including advertising or data analytics, without
                your explicit consent.
              </p>
            </section>

            <section id="liability">
              <h2>8. Limitation of Liability</h2>
              <p>
                To the fullest extent permitted by law, Mailngine and its
                officers, directors, employees, and agents shall not be liable
                for any indirect, incidental, special, consequential, or
                punitive damages, including loss of profits, data, or goodwill,
                arising from your use of the Service.
              </p>
              <p>
                Our total cumulative liability to you for all claims arising
                from or related to the Service shall not exceed the amount you
                paid to Mailngine in the twelve (12) months preceding the
                claim.
              </p>
            </section>

            <section id="termination">
              <h2>9. Termination</h2>
              <p>
                You may terminate your account at any time through the
                dashboard settings. Upon termination, we will delete your
                account data within 30 days, except where retention is required
                by law or for legitimate business purposes such as fraud
                prevention.
              </p>
              <p>
                We may suspend or terminate your account if you violate these
                Terms, fail to pay applicable fees, or if we reasonably believe
                your use of the Service poses a risk to our infrastructure or
                other users.
              </p>
            </section>

            <section id="modifications">
              <h2>10. Modifications to Terms</h2>
              <p>
                We may update these Terms from time to time. We will notify you
                of material changes by posting the revised Terms on our website
                and updating the "Effective date" above. Your continued use of
                the Service after the effective date constitutes your acceptance
                of the revised Terms.
              </p>
            </section>

            <section id="governing-law">
              <h2>11. Governing Law</h2>
              <p>
                These Terms shall be governed by and construed in accordance
                with the laws of the State of Delaware, United States, without
                regard to its conflict of law provisions. Any disputes arising
                from these Terms or the Service shall be resolved in the state
                or federal courts located in Delaware.
              </p>
            </section>

            <section id="contact">
              <h2>12. Contact Information</h2>
              <p>
                If you have any questions about these Terms, please contact us
                at{" "}
                <a href="mailto:legal@mailngine.com">legal@mailngine.com</a> or
                visit our <a href="/contact">contact page</a>.
              </p>
            </section>
          </article>
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .legal-hero {
          text-align: center;
          max-width: 600px;
          margin: 0 auto;
        }
        .legal-hero h1 {
          margin-bottom: var(--space-3);
        }
        .legal-hero p {
          font-size: var(--font-size-lg);
          color: var(--color-text-secondary);
        }

        .legal-content {
          max-width: 800px;
          margin: 0 auto;
        }

        .legal-toc {
          background-color: var(--color-bg-surface);
          border-radius: var(--radius-md);
          padding: var(--space-6);
          margin-bottom: var(--space-10);
        }
        .legal-toc h2 {
          font-size: var(--font-size-lg);
          margin-bottom: var(--space-4);
        }
        .legal-toc ol {
          padding-left: var(--space-6);
          display: flex;
          flex-direction: column;
          gap: var(--space-2);
        }
        .legal-toc li {
          color: var(--color-text-secondary);
        }
        .legal-toc a {
          font-size: var(--font-size-sm);
        }

        .legal-body section {
          margin-bottom: var(--space-10);
        }
        .legal-body h2 {
          font-size: var(--font-size-xl);
          margin-bottom: var(--space-4);
          padding-top: var(--space-4);
        }
        .legal-body h3 {
          font-size: var(--font-size-lg);
          margin-bottom: var(--space-3);
          margin-top: var(--space-6);
        }
        .legal-body p {
          margin-bottom: var(--space-4);
          line-height: 1.7;
        }
        .legal-body ul {
          padding-left: var(--space-6);
          margin-bottom: var(--space-4);
          display: flex;
          flex-direction: column;
          gap: var(--space-2);
        }
        .legal-body li {
          color: var(--color-text-secondary);
          line-height: 1.6;
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "Terms of Service - Mailngine",
  meta: [
    {
      name: "description",
      content:
        "Mailngine Terms of Service. Read our terms and conditions for using the Mailngine email delivery platform.",
    },
    {
      property: "og:title",
      content: "Terms of Service - Mailngine",
    },
    {
      property: "og:description",
      content:
        "Mailngine Terms of Service. Read our terms and conditions for using the Mailngine email delivery platform.",
    },
  ],
};
