import type { RequestHandler } from "@builder.io/qwik-city";
import { posts } from "~/content/blog";

const ORIGIN = "https://mailngine.com";

interface SitemapEntry {
  loc: string;
  lastmod: string;
  changefreq?: string;
  priority?: string;
}

const staticPages: SitemapEntry[] = [
  { loc: "/", lastmod: "2026-04-04", changefreq: "weekly", priority: "1.0" },
  {
    loc: "/pricing/",
    lastmod: "2026-04-04",
    changefreq: "monthly",
    priority: "0.9",
  },
  {
    loc: "/docs/",
    lastmod: "2026-04-04",
    changefreq: "weekly",
    priority: "0.9",
  },
  {
    loc: "/docs/api/",
    lastmod: "2026-04-04",
    changefreq: "weekly",
    priority: "0.8",
  },
  {
    loc: "/integrations/",
    lastmod: "2026-04-04",
    changefreq: "monthly",
    priority: "0.7",
  },
  {
    loc: "/blog/",
    lastmod: "2026-04-03",
    changefreq: "weekly",
    priority: "0.7",
  },
  {
    loc: "/about/",
    lastmod: "2026-04-04",
    changefreq: "monthly",
    priority: "0.5",
  },
  {
    loc: "/changelog/",
    lastmod: "2026-04-04",
    changefreq: "weekly",
    priority: "0.5",
  },
  {
    loc: "/status/",
    lastmod: "2026-04-04",
    changefreq: "daily",
    priority: "0.4",
  },
  {
    loc: "/contact/",
    lastmod: "2026-04-04",
    changefreq: "monthly",
    priority: "0.4",
  },
  {
    loc: "/careers/",
    lastmod: "2026-04-04",
    changefreq: "monthly",
    priority: "0.3",
  },
  {
    loc: "/privacy/",
    lastmod: "2026-04-01",
    changefreq: "yearly",
    priority: "0.2",
  },
  {
    loc: "/terms/",
    lastmod: "2026-04-01",
    changefreq: "yearly",
    priority: "0.2",
  },
];

export const onGet: RequestHandler = async ({ send }) => {
  const blogEntries: SitemapEntry[] = posts.map((post) => ({
    loc: `/blog/${post.slug}/`,
    lastmod: post.date,
    priority: "0.6",
  }));

  const entries = [...staticPages, ...blogEntries];

  const xml = [
    '<?xml version="1.0" encoding="UTF-8"?>',
    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
    ...entries.map(
      (entry) =>
        `<url><loc>${ORIGIN}${entry.loc}</loc><lastmod>${entry.lastmod}</lastmod>${entry.changefreq ? `<changefreq>${entry.changefreq}</changefreq>` : ""}${entry.priority ? `<priority>${entry.priority}</priority>` : ""}</url>`,
    ),
    "</urlset>",
  ].join("\n");

  send(
    new Response(xml, {
      headers: {
        "Content-Type": "application/xml",
        "Cache-Control": "public, max-age=3600",
      },
    }),
  );
};
