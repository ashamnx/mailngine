import { component$ } from "@builder.io/qwik";
import {
  routeAction$,
  zod$,
  z,
  Form,
  type DocumentHead,
} from "@builder.io/qwik-city";
import { Section } from "~/components/ui/section";
import { Card } from "~/components/ui/card";
import { Button } from "~/components/ui/button";

export const useContactAction = routeAction$(
  async (data) => {
    // eslint-disable-next-line no-console
    console.log("Contact form submission:", data);
    return { success: true };
  },
  zod$({
    name: z.string().min(1, "Name is required"),
    email: z.string().email("Please enter a valid email address"),
    subject: z.string().min(1, "Subject is required"),
    message: z.string().min(10, "Message must be at least 10 characters"),
  }),
);

export default component$(() => {
  const action = useContactAction();

  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="contact-hero">
          <h1>Contact Us</h1>
          <p>
            Have a question, feature request, or need help? We'd love to hear
            from you.
          </p>
        </div>
      </Section>

      <Section padding="lg">
        <div class="contact-grid">
          <div class="contact-form-wrapper">
            {action.value?.success ? (
              <Card elevated>
                <div class="contact-success">
                  <span class="material-symbols-rounded contact-success__icon">
                    check_circle
                  </span>
                  <h3>Message sent</h3>
                  <p>
                    Thank you for reaching out. We'll get back to you within 24
                    hours.
                  </p>
                </div>
              </Card>
            ) : (
              <Form action={action}>
                <div class="contact-form">
                  <div class="contact-form__field">
                    <label for="name">Name</label>
                    <input
                      type="text"
                      id="name"
                      name="name"
                      placeholder="Your name"
                      required
                    />
                    {action.value?.fieldErrors?.name && (
                      <span class="contact-form__error">
                        {action.value.fieldErrors.name}
                      </span>
                    )}
                  </div>

                  <div class="contact-form__field">
                    <label for="email">Email</label>
                    <input
                      type="email"
                      id="email"
                      name="email"
                      placeholder="you@example.com"
                      required
                    />
                    {action.value?.fieldErrors?.email && (
                      <span class="contact-form__error">
                        {action.value.fieldErrors.email}
                      </span>
                    )}
                  </div>

                  <div class="contact-form__field">
                    <label for="subject">Subject</label>
                    <input
                      type="text"
                      id="subject"
                      name="subject"
                      placeholder="What is this about?"
                      required
                    />
                    {action.value?.fieldErrors?.subject && (
                      <span class="contact-form__error">
                        {action.value.fieldErrors.subject}
                      </span>
                    )}
                  </div>

                  <div class="contact-form__field">
                    <label for="message">Message</label>
                    <textarea
                      id="message"
                      name="message"
                      rows={6}
                      placeholder="Tell us more..."
                      required
                    />
                    {action.value?.fieldErrors?.message && (
                      <span class="contact-form__error">
                        {action.value.fieldErrors.message}
                      </span>
                    )}
                  </div>

                  <Button type="submit" variant="primary" size="lg">
                    Send Message
                  </Button>
                </div>
              </Form>
            )}
          </div>

          <div class="contact-info">
            <Card elevated>
              <h3 class="contact-info__title">Get in Touch</h3>

              <div class="contact-info__item">
                <span class="material-symbols-rounded">mail</span>
                <div>
                  <p class="contact-info__label">Email</p>
                  <a href="mailto:hello@mailngine.com">hello@mailngine.com</a>
                </div>
              </div>

              <div class="contact-info__item">
                <span class="material-symbols-rounded">schedule</span>
                <div>
                  <p class="contact-info__label">Response Time</p>
                  <p>We typically respond within 24 hours</p>
                </div>
              </div>

              <div class="contact-info__divider" />

              <h4 class="contact-info__subtitle">Follow Us</h4>
              <div class="contact-info__social">
                <a
                  href="https://github.com/mailngine"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  GitHub
                </a>
                <a
                  href="https://twitter.com/mailngine"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Twitter
                </a>
                <a
                  href="https://discord.gg/mailngine"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Discord
                </a>
              </div>
            </Card>
          </div>
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .contact-hero {
          text-align: center;
          max-width: 600px;
          margin: 0 auto;
        }
        .contact-hero h1 {
          margin-bottom: var(--space-4);
        }
        .contact-hero p {
          font-size: var(--font-size-lg);
        }

        .contact-grid {
          display: grid;
          grid-template-columns: 1fr 380px;
          gap: var(--space-10);
          align-items: start;
        }

        .contact-form {
          display: flex;
          flex-direction: column;
          gap: var(--space-5);
        }
        .contact-form__field {
          display: flex;
          flex-direction: column;
          gap: var(--space-2);
        }
        .contact-form__field label {
          font-family: var(--font-primary);
          font-weight: var(--font-weight-medium);
          font-size: var(--font-size-sm);
          color: var(--color-text-primary);
        }
        .contact-form__field input,
        .contact-form__field textarea {
          padding: var(--space-3) var(--space-4);
          border: 1px solid var(--color-border);
          border-radius: var(--radius-sm);
          font-size: var(--font-size-base);
          transition: border-color var(--transition-fast);
          background-color: var(--color-bg-app);
          color: var(--color-text-primary);
        }
        .contact-form__field input:focus,
        .contact-form__field textarea:focus {
          outline: none;
          border-color: var(--color-primary);
          box-shadow: 0 0 0 3px rgba(26, 115, 232, 0.15);
        }
        .contact-form__field textarea {
          resize: vertical;
        }
        .contact-form__error {
          font-size: var(--font-size-sm);
          color: var(--color-error);
        }

        .contact-success {
          text-align: center;
          padding: var(--space-8) var(--space-4);
        }
        .contact-success__icon {
          font-size: 48px;
          color: var(--color-success);
          margin-bottom: var(--space-4);
        }
        .contact-success h3 {
          margin-bottom: var(--space-2);
        }

        .contact-info__title {
          margin-bottom: var(--space-6);
        }
        .contact-info__item {
          display: flex;
          gap: var(--space-3);
          margin-bottom: var(--space-5);
        }
        .contact-info__item .material-symbols-rounded {
          color: var(--color-primary);
          flex-shrink: 0;
          margin-top: 2px;
        }
        .contact-info__label {
          font-weight: var(--font-weight-medium);
          color: var(--color-text-primary);
          font-size: var(--font-size-sm);
          margin-bottom: var(--space-1);
        }
        .contact-info__divider {
          height: 1px;
          background-color: var(--color-border-light);
          margin: var(--space-5) 0;
        }
        .contact-info__subtitle {
          font-size: var(--font-size-base);
          margin-bottom: var(--space-3);
        }
        .contact-info__social {
          display: flex;
          gap: var(--space-4);
        }
        .contact-info__social a {
          font-size: var(--font-size-sm);
          font-weight: var(--font-weight-medium);
        }

        @media (max-width: 768px) {
          .contact-grid {
            grid-template-columns: 1fr;
          }
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "Contact Us - Mailngine",
  meta: [
    {
      name: "description",
      content:
        "Get in touch with the Mailngine team. We're here to help with questions, feature requests, and support.",
    },
    {
      property: "og:title",
      content: "Contact Us - Mailngine",
    },
    {
      property: "og:description",
      content:
        "Get in touch with the Mailngine team. We're here to help with questions, feature requests, and support.",
    },
  ],
};
