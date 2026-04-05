import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { Section } from "~/components/ui/section";
import { Card } from "~/components/ui/card";
import { Badge } from "~/components/ui/badge";
import { posts } from "~/content/blog";

export default component$(() => {
  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="blog-hero">
          <h1>Blog</h1>
          <p>
            Product updates, engineering deep dives, and guides from the
            Mailngine team.
          </p>
        </div>
      </Section>

      <Section padding="lg">
        <div class="blog-grid">
          {posts.map((post) => (
            <a
              key={post.slug}
              href={`/blog/${post.slug}`}
              class="blog-card-link"
            >
              <Card elevated class="blog-card">
                <div class="blog-card__content">
                  <div class="blog-card__meta">
                    <time>{post.date}</time>
                    <Badge variant="neutral">{post.readTime}</Badge>
                  </div>
                  <h2 class="blog-card__title">{post.title}</h2>
                  <p class="blog-card__excerpt">{post.excerpt}</p>
                  <span class="blog-card__read-more">
                    Read more
                    <span class="material-symbols-rounded">arrow_forward</span>
                  </span>
                </div>
              </Card>
            </a>
          ))}
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .blog-hero {
          text-align: center;
          max-width: 600px;
          margin: 0 auto;
        }
        .blog-hero h1 {
          margin-bottom: var(--space-4);
        }
        .blog-hero p {
          font-size: var(--font-size-lg);
        }

        .blog-grid {
          display: grid;
          grid-template-columns: repeat(3, 1fr);
          gap: var(--space-6);
        }

        .blog-card-link {
          text-decoration: none;
          color: inherit;
          display: block;
        }
        .blog-card-link:hover {
          color: inherit;
        }

        .blog-card__content {
          display: flex;
          flex-direction: column;
          gap: var(--space-3);
        }
        .blog-card__meta {
          display: flex;
          align-items: center;
          gap: var(--space-3);
        }
        .blog-card__meta time {
          font-size: var(--font-size-sm);
          color: var(--color-text-secondary);
        }
        .blog-card__title {
          font-size: var(--font-size-xl);
          line-height: 1.3;
        }
        .blog-card__excerpt {
          font-size: var(--font-size-sm);
          line-height: 1.6;
          flex: 1;
        }
        .blog-card__read-more {
          display: inline-flex;
          align-items: center;
          gap: var(--space-1);
          font-size: var(--font-size-sm);
          font-weight: var(--font-weight-medium);
          color: var(--color-primary);
          margin-top: var(--space-2);
        }
        .blog-card__read-more .material-symbols-rounded {
          font-size: 18px;
          transition: transform var(--transition-fast);
        }
        .blog-card-link:hover .blog-card__read-more .material-symbols-rounded {
          transform: translateX(4px);
        }

        @media (max-width: 768px) {
          .blog-grid {
            grid-template-columns: 1fr;
          }
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = {
  title: "Blog - Mailngine",
  meta: [
    {
      name: "description",
      content:
        "Product updates, engineering deep dives, and guides from the Mailngine team.",
    },
    {
      property: "og:title",
      content: "Blog - Mailngine",
    },
    {
      property: "og:description",
      content:
        "Product updates, engineering deep dives, and guides from the Mailngine team.",
    },
  ],
};
