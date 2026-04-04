export interface BlogPost {
  slug: string;
  title: string;
  excerpt: string;
  date: string;
  readTime: string;
  author: string;
  content: string;
}

export const posts: BlogPost[] = [
  {
    slug: "introducing-hello-mail",
    title: "Introducing Hello Mail",
    excerpt:
      "We built Hello Mail because developers deserve email infrastructure that just works.",
    date: "2026-04-01",
    readTime: "3 min",
    author: "Hello Mail Team",
    content: `
      <p>We built Hello Mail because we were tired of fighting with email infrastructure. Every project we worked on had the same story: pick an email provider, wrestle with DNS records, debug deliverability issues, and hope for the best. We knew there had to be a better way.</p>

      <p>Hello Mail is the email platform we always wanted. It starts with a clean REST API that feels natural whether you're writing Go, Node.js, or PHP. Domain verification is guided step-by-step, and if you're on Cloudflare, we can configure your DNS records automatically. No more copy-pasting TXT records and wondering if you got the spacing right.</p>

      <p>Under the hood, Hello Mail is built on battle-tested open source components. Postfix handles SMTP delivery, OpenDKIM signs every outbound message, and a custom Go service orchestrates everything through a task queue backed by Valkey. Every email gets a delivery receipt, and our webhook system notifies your application in real time with HMAC-signed payloads you can trust.</p>

      <p>We're launching today with support for transactional email, domain management with full DKIM and SPF configuration, and SDKs for Go, Node.js, and Laravel. Analytics, templates, and team management are landing in the coming weeks. We can't wait to see what you build with Hello Mail.</p>
    `,
  },
  {
    slug: "getting-started-go-sdk",
    title: "Getting Started with the Go SDK",
    excerpt:
      "Learn how to send your first email using the Hello Mail Go SDK in under 5 minutes.",
    date: "2026-04-02",
    readTime: "5 min",
    author: "Hello Mail Team",
    content: `
      <p>The Hello Mail Go SDK is designed to feel like a native part of your Go application. There are no complex configuration objects or builder patterns to memorize. You create a client with your API key, build a message struct, and call Send. That's it.</p>

      <p>To get started, install the SDK with <code>go get github.com/hellomail/hellomail-go</code>. Then create a new client by passing your API key, which you can generate from the Hello Mail dashboard under Settings &gt; API Keys. We recommend storing the key in an environment variable rather than hardcoding it in your source code.</p>

      <p>Sending an email is straightforward. Construct a <code>hellomail.Message</code> with the sender address, recipient, subject, and either an HTML or plain text body. Pass it to <code>client.Send(ctx, message)</code> and you'll get back a response containing the message ID and delivery status. The SDK handles retries, rate limiting, and connection pooling automatically so you can focus on your application logic.</p>

      <p>For more advanced use cases, the SDK supports batch sending, template rendering with variable substitution, and attachment handling. Every method accepts a <code>context.Context</code> so you get full control over timeouts and cancellation. Check out the <a href="/docs/api">API reference</a> for the complete list of options, and join our community if you have questions or feedback.</p>
    `,
  },
  {
    slug: "domain-verification-guide",
    title: "Domain Verification Best Practices",
    excerpt:
      "A complete guide to setting up SPF, DKIM, and DMARC for your sending domain.",
    date: "2026-04-03",
    readTime: "7 min",
    author: "Hello Mail Team",
    content: `
      <p>Email authentication is the foundation of deliverability. Without proper SPF, DKIM, and DMARC records, your emails are far more likely to land in spam folders or be rejected outright. The good news is that Hello Mail handles most of the heavy lifting for you, but understanding how these protocols work will help you debug issues and make informed decisions about your email infrastructure.</p>

      <p>SPF (Sender Policy Framework) tells receiving mail servers which IP addresses are authorized to send email on behalf of your domain. When you add a domain to Hello Mail, we provide you with an SPF include record that authorizes our sending infrastructure. You add this as a TXT record on your domain, and receiving servers will check it every time they get an email claiming to be from your domain. Keep in mind that SPF has a 10-lookup limit, so if you use multiple email providers, you'll need to be mindful of how many includes you're chaining together.</p>

      <p>DKIM (DomainKeys Identified Mail) adds a cryptographic signature to every outbound email. Hello Mail generates a unique 2048-bit RSA key pair for each domain you register. The private key stays on our servers and is used to sign the DKIM header of every message. The public key is published as a DNS record on your domain so receiving servers can verify the signature. This proves that the email hasn't been tampered with in transit and that it genuinely came from an authorized sender.</p>

      <p>DMARC (Domain-based Message Authentication, Reporting, and Conformance) ties SPF and DKIM together with a policy that tells receiving servers what to do when authentication fails. We recommend starting with a policy of <code>p=none</code> so you can monitor reports without affecting delivery. Once you're confident that all legitimate email from your domain passes authentication, you can move to <code>p=quarantine</code> or <code>p=reject</code>. Hello Mail's analytics dashboard shows your DMARC alignment rate so you can track your progress over time.</p>
    `,
  },
];
