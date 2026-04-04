import { component$ } from "@builder.io/qwik";
import { routeLoader$, type DocumentHead } from "@builder.io/qwik-city";
import { Section } from "~/components/ui/section";
import { Button } from "~/components/ui/button";
import { Badge } from "~/components/ui/badge";
import { posts } from "~/content/blog";

export const useBlogPost = routeLoader$(({ params, status }) => {
  const post = posts.find((p) => p.slug === params.slug);

  if (!post) {
    status(404);
    return null;
  }

  return post;
});

export default component$(() => {
  const postSignal = useBlogPost();
  const post = postSignal.value;

  if (!post) {
    return (
      <Section padding="lg">
        <div class="blog-post-not-found">
          <h1>Post not found</h1>
          <p>The blog post you're looking for doesn't exist.</p>
          <Button variant="primary" href="/blog">
            Back to Blog
          </Button>
        </div>
      </Section>
    );
  }

  return (
    <>
      <Section background="gradient" padding="lg">
        <div class="blog-post-hero">
          <a href="/blog" class="blog-post-hero__back">
            <span class="material-symbols-rounded">arrow_back</span>
            Back to Blog
          </a>
          <h1>{post.title}</h1>
          <div class="blog-post-hero__meta">
            <time>{post.date}</time>
            <span class="blog-post-hero__separator" aria-hidden="true" />
            <span>{post.author}</span>
            <span class="blog-post-hero__separator" aria-hidden="true" />
            <Badge variant="neutral">{post.readTime}</Badge>
          </div>
        </div>
      </Section>

      <Section padding="lg">
        <article
          class="blog-post-body"
          dangerouslySetInnerHTML={post.content}
        />
      </Section>

      <Section background="surface" padding="md">
        <div class="blog-post-footer">
          <Button variant="secondary" href="/blog">
            <span class="material-symbols-rounded">arrow_back</span>
            All Posts
          </Button>
        </div>
      </Section>

      <style
        dangerouslySetInnerHTML={`
        .blog-post-not-found {
          text-align: center;
          padding: var(--space-16) 0;
        }
        .blog-post-not-found h1 {
          margin-bottom: var(--space-4);
        }
        .blog-post-not-found p {
          margin-bottom: var(--space-8);
          font-size: var(--font-size-lg);
        }

        .blog-post-hero {
          max-width: 720px;
          margin: 0 auto;
        }
        .blog-post-hero__back {
          display: inline-flex;
          align-items: center;
          gap: var(--space-2);
          font-size: var(--font-size-sm);
          font-weight: var(--font-weight-medium);
          margin-bottom: var(--space-6);
        }
        .blog-post-hero__back .material-symbols-rounded {
          font-size: 18px;
        }
        .blog-post-hero h1 {
          font-size: var(--font-size-4xl);
          margin-bottom: var(--space-6);
        }
        .blog-post-hero__meta {
          display: flex;
          align-items: center;
          gap: var(--space-3);
          color: var(--color-text-secondary);
          font-size: var(--font-size-sm);
        }
        .blog-post-hero__separator {
          width: 4px;
          height: 4px;
          border-radius: 50%;
          background-color: var(--color-text-disabled);
        }

        .blog-post-body {
          max-width: 720px;
          margin: 0 auto;
        }
        .blog-post-body p {
          margin-bottom: var(--space-6);
          font-size: var(--font-size-lg);
          line-height: 1.8;
        }
        .blog-post-body code {
          font-family: var(--font-mono);
          font-size: var(--font-size-sm);
          background-color: var(--color-bg-surface-alt);
          padding: 2px var(--space-2);
          border-radius: var(--radius-sm);
        }
        .blog-post-body a {
          text-decoration: underline;
          text-underline-offset: 2px;
        }
        .blog-post-body h2 {
          margin-top: var(--space-10);
          margin-bottom: var(--space-4);
        }
        .blog-post-body h3 {
          margin-top: var(--space-8);
          margin-bottom: var(--space-3);
        }
        .blog-post-body ul,
        .blog-post-body ol {
          padding-left: var(--space-6);
          margin-bottom: var(--space-6);
        }
        .blog-post-body li {
          margin-bottom: var(--space-2);
          line-height: 1.7;
          color: var(--color-text-secondary);
        }

        .blog-post-footer {
          display: flex;
          justify-content: center;
        }

        @media (max-width: 768px) {
          .blog-post-hero h1 {
            font-size: var(--font-size-3xl);
          }
          .blog-post-body p {
            font-size: var(--font-size-base);
          }
        }
      `}
      />
    </>
  );
});

export const head: DocumentHead = ({ resolveValue }) => {
  const post = resolveValue(useBlogPost);

  if (!post) {
    return {
      title: "Post Not Found - Hello Mail",
    };
  }

  return {
    title: `${post.title} - Hello Mail Blog`,
    meta: [
      {
        name: "description",
        content: post.excerpt,
      },
      {
        name: "og:title",
        content: `${post.title} - Hello Mail Blog`,
      },
      {
        name: "og:description",
        content: post.excerpt,
      },
    ],
  };
};
