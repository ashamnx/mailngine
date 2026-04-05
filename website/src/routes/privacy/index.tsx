import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section } from "~/components/ui/section";

export default component$(() => {
  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="legal-hero">
          <h1>Privacy Policy</h1>
          <p>Effective date: April 1, 2026</p>
        </div>
      </Section>

      <Section padding="lg">
        <div class="legal-content">
          <nav class="legal-toc">
            <h2>Table of Contents</h2>
            <ol>
              <li>
                <a href="#information-we-collect">Information We Collect</a>
              </li>
              <li>
                <a href="#how-we-use">How We Use Your Information</a>
              </li>
              <li>
                <a href="#data-storage">Data Storage and Security</a>
              </li>
              <li>
                <a href="#data-retention">Data Retention</a>
              </li>
              <li>
                <a href="#third-party">Third-Party Services</a>
              </li>
              <li>
                <a href="#your-rights">Your Rights</a>
              </li>
              <li>
                <a href="#children">Children's Privacy</a>
              </li>
              <li>
                <a href="#changes">Changes to This Policy</a>
              </li>
              <li>
                <a href="#contact">Contact Us</a>
              </li>
            </ol>
          </nav>

          <article class="legal-body">
            <section id="information-we-collect">
              <h2>1. Information We Collect</h2>

              <h3>Account Data</h3>
              <p>
                When you create an account, we collect your name, email
                address, and authentication credentials. If you sign in via
                Google OAuth, we receive your name, email address, and profile
                picture from Google. We do not store your Google password.
              </p>

              <h3>Email Content</h3>
              <p>
                When you send emails through our Service, we temporarily
                process and store the email content, including subject lines,
                body text, attachments, and recipient addresses. This data is
                necessary to deliver your emails and provide delivery status
                information.
              </p>

              <h3>Usage Data</h3>
              <p>
                We collect information about how you interact with the Service,
                including API requests, dashboard page views, feature usage
                patterns, and error logs. This data is used to improve the
                Service and diagnose technical issues.
              </p>

              <h3>Cookies</h3>
              <p>
                We use essential cookies to maintain your authentication
                session and remember your preferences. We do not use
                third-party advertising or tracking cookies. You can configure
                your browser to block cookies, but this may affect your ability
                to use the Service.
              </p>
            </section>

            <section id="how-we-use">
              <h2>2. How We Use Your Information</h2>
              <p>We use the information we collect to:</p>
              <ul>
                <li>Provide, maintain, and improve the Service.</li>
                <li>
                  Deliver emails on your behalf and provide delivery status
                  and analytics.
                </li>
                <li>
                  Authenticate your identity and secure your account.
                </li>
                <li>
                  Communicate with you about service updates, security alerts,
                  and support requests.
                </li>
                <li>
                  Monitor for abuse and enforce our Terms of Service and
                  Acceptable Use Policy.
                </li>
                <li>
                  Generate aggregated, anonymized usage statistics to improve
                  the Service.
                </li>
              </ul>
              <p>
                We do not sell your personal information to third parties. We
                do not use your email content for advertising, profiling, or
                any purpose other than delivering the Service.
              </p>
            </section>

            <section id="data-storage">
              <h2>3. Data Storage and Security</h2>
              <p>
                Your data is stored on servers hosted by DigitalOcean in data
                centers located in the United States. All data is encrypted at
                rest using AES-256 encryption. All data transmitted between
                your applications and our Service is encrypted in transit using
                TLS 1.2 or higher.
              </p>
              <p>
                We implement industry-standard security measures including
                network firewalls, intrusion detection systems, regular
                security audits, and access controls. API keys are hashed
                before storage and cannot be retrieved in plain text after
                creation.
              </p>
              <p>
                While we take reasonable precautions to protect your data, no
                method of electronic transmission or storage is 100% secure. We
                cannot guarantee absolute security.
              </p>
            </section>

            <section id="data-retention">
              <h2>4. Data Retention</h2>
              <p>
                We retain your account data for as long as your account is
                active. Email content is retained for up to 30 days after
                delivery to support delivery troubleshooting and analytics.
                After 30 days, email content is permanently deleted.
              </p>
              <p>
                Delivery logs and analytics data (metadata such as timestamps,
                delivery status, and open/click events) are retained for up to
                12 months. When you delete your account, we remove your data
                within 30 days, except where retention is required by law.
              </p>
            </section>

            <section id="third-party">
              <h2>5. Third-Party Services</h2>
              <p>
                We use the following third-party services to operate
                Mailngine:
              </p>
              <ul>
                <li>
                  <strong>DigitalOcean</strong> -- Cloud infrastructure hosting
                  for our servers and databases.
                </li>
                <li>
                  <strong>Google OAuth</strong> -- Authentication provider for
                  user sign-in. Google receives your authentication request but
                  does not receive your Mailngine usage data.
                </li>
                <li>
                  <strong>Cloudflare</strong> -- DNS management and CDN
                  services for domain verification and asset delivery. When you
                  enable auto-DNS, Cloudflare receives your domain
                  configuration data.
                </li>
              </ul>
              <p>
                Each third-party service operates under its own privacy policy.
                We encourage you to review their policies to understand how
                they handle your data.
              </p>
            </section>

            <section id="your-rights">
              <h2>6. Your Rights</h2>
              <p>
                Under applicable data protection laws, including the General
                Data Protection Regulation (GDPR), you have the following
                rights:
              </p>
              <ul>
                <li>
                  <strong>Right of Access</strong> -- You can request a copy of
                  the personal data we hold about you.
                </li>
                <li>
                  <strong>Right to Rectification</strong> -- You can request
                  that we correct inaccurate or incomplete personal data.
                </li>
                <li>
                  <strong>Right to Erasure</strong> -- You can request that we
                  delete your personal data, subject to any legal obligations
                  requiring retention.
                </li>
                <li>
                  <strong>Right to Data Portability</strong> -- You can request
                  your data in a structured, machine-readable format for
                  transfer to another service.
                </li>
                <li>
                  <strong>Right to Object</strong> -- You can object to the
                  processing of your personal data for certain purposes.
                </li>
                <li>
                  <strong>Right to Restrict Processing</strong> -- You can
                  request that we limit the processing of your personal data
                  under certain circumstances.
                </li>
              </ul>
              <p>
                To exercise any of these rights, please contact us at{" "}
                <a href="mailto:privacy@mailngine.com">
                  privacy@mailngine.com
                </a>
                . We will respond to your request within 30 days.
              </p>
            </section>

            <section id="children">
              <h2>7. Children's Privacy</h2>
              <p>
                The Service is not intended for children under the age of 16.
                We do not knowingly collect personal information from children
                under 16. If we become aware that we have collected personal
                data from a child under 16, we will take steps to delete that
                information promptly.
              </p>
            </section>

            <section id="changes">
              <h2>8. Changes to This Policy</h2>
              <p>
                We may update this Privacy Policy from time to time. We will
                notify you of material changes by posting the updated policy on
                our website and updating the "Effective date" above. For
                significant changes, we will also send a notification to the
                email address associated with your account.
              </p>
              <p>
                Your continued use of the Service after the effective date of a
                revised policy constitutes your acceptance of the changes.
              </p>
            </section>

            <section id="contact">
              <h2>9. Contact Us</h2>
              <p>
                If you have any questions about this Privacy Policy or how we
                handle your data, please contact us at{" "}
                <a href="mailto:privacy@mailngine.com">
                  privacy@mailngine.com
                </a>{" "}
                or visit our <a href="/contact">contact page</a>.
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
  title: "Privacy Policy - Mailngine",
  meta: [
    {
      name: "description",
      content:
        "Mailngine Privacy Policy. Learn how we collect, use, and protect your data.",
    },
    {
      property: "og:title",
      content: "Privacy Policy - Mailngine",
    },
    {
      property: "og:description",
      content:
        "Mailngine Privacy Policy. Learn how we collect, use, and protect your data.",
    },
  ],
};
