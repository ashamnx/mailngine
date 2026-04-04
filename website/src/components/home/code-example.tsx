import { component$, useSignal } from "@builder.io/qwik";
import { CodeBlock } from "~/components/ui";

const tabs = [
  {
    label: "Go",
    language: "go",
    code: `client := hellomail.New("hm_live_...")
email, err := client.Emails.Send(ctx, &hellomail.SendEmailParams{
    From:    "hello@example.com",
    To:      []string{"user@example.com"},
    Subject: "Welcome!",
    HTML:    "<h1>Hello from Go</h1>",
})`,
  },
  {
    label: "Node.js",
    language: "typescript",
    code: `const client = new HelloMail('hm_live_...');
const email = await client.emails.send({
    from: 'hello@example.com',
    to: ['user@example.com'],
    subject: 'Welcome!',
    html: '<h1>Hello from Node.js</h1>',
});`,
  },
  {
    label: "Laravel",
    language: "php",
    code: `$client = new HelloMail('hm_live_...');
$email = $client->emails()->send([
    'from' => 'hello@example.com',
    'to' => ['user@example.com'],
    'subject' => 'Welcome!',
    'html' => '<h1>Hello from Laravel</h1>',
]);`,
  },
];

export const CodeExample = component$(() => {
  const activeTab = useSignal(0);

  return (
    <div class="code-example">
      <h2 class="code-example__heading">Send your first email</h2>
      <p class="code-example__subheading">
        Get started with just a few lines of code using your preferred language.
      </p>

      <div class="code-example__tabs" role="tablist">
        {tabs.map((tab, index) => (
          <button
            key={tab.label}
            role="tab"
            aria-selected={activeTab.value === index}
            class={[
              "code-example__tab",
              activeTab.value === index ? "code-example__tab--active" : "",
            ]
              .filter(Boolean)
              .join(" ")}
            onClick$={() => {
              activeTab.value = index;
            }}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div class="code-example__panel" role="tabpanel">
        <CodeBlock
          code={tabs[activeTab.value].code}
          language={tabs[activeTab.value].language}
        />
      </div>
    </div>
  );
});
