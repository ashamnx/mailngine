import { component$ } from "@builder.io/qwik";
import { AccordionItem } from "~/components/ui";

const faqItems = [
  {
    question: "Can I upgrade or downgrade at any time?",
    answer:
      "Yes. You can switch plans at any time from your account dashboard. When upgrading, you'll be prorated for the remainder of your billing cycle. Downgrades take effect at the start of the next billing period.",
  },
  {
    question: "What happens if I exceed my email limit?",
    answer:
      "We won't cut you off. If you go over your plan's limit, emails will continue to be delivered. You'll receive a notification and have the option to upgrade or pay for the overage at the end of the billing cycle.",
  },
  {
    question: "Do you offer annual billing?",
    answer:
      "Yes. Annual billing is available for Starter and Pro plans at a 20% discount. Contact us for custom Enterprise billing arrangements.",
  },
  {
    question: "Is there a free trial for paid plans?",
    answer:
      "All paid plans come with a 14-day free trial. No credit card required to start. You can explore all features and only start paying when you're ready.",
  },
  {
    question: "What payment methods do you accept?",
    answer:
      "We accept all major credit cards (Visa, Mastercard, American Express), as well as wire transfers for Enterprise plans. All payments are processed securely through Stripe.",
  },
];

export const PricingFaq = component$(() => {
  return (
    <div class="pricing-faq">
      <h2 class="pricing-faq__heading">Frequently asked questions</h2>
      <div class="pricing-faq__list">
        {faqItems.map((item) => (
          <AccordionItem key={item.question} title={item.question}>
            <p>{item.answer}</p>
          </AccordionItem>
        ))}
      </div>
    </div>
  );
});
