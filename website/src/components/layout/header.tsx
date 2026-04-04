import { component$, useSignal } from "@builder.io/qwik";
import { Button } from "../ui/button";
import "./header.css";

const NAV_LINKS = [
  { label: "Home", href: "/" },
  { label: "Pricing", href: "/pricing" },
  { label: "Docs", href: "/docs" },
  { label: "Integrations", href: "/integrations" },
  { label: "Blog", href: "/blog" },
  { label: "Status", href: "/status" },
] as const;

export const Header = component$(() => {
  const mobileMenuOpen = useSignal(false);

  return (
    <header class="hm-header">
      <div class="hm-header__inner">
        {/* Logo */}
        <a href="/" class="hm-header__logo" aria-label="Hello Mail home">
          <span class="material-symbols-rounded hm-header__logo-icon">
            mail
          </span>
          <span class="hm-header__logo-text">Hello Mail</span>
        </a>

        {/* Desktop Navigation */}
        <nav class="hm-header__nav" aria-label="Main navigation">
          {NAV_LINKS.map((link) => (
            <a key={link.href} href={link.href} class="hm-header__nav-link">
              {link.label}
            </a>
          ))}
        </nav>

        {/* Right Actions */}
        <div class="hm-header__actions">
          <a href="/sign-in" class="hm-header__sign-in">
            Sign In
          </a>
          <Button variant="primary" size="sm" href="/get-started">
            Get Started
          </Button>
        </div>

        {/* Mobile Hamburger */}
        <button
          class="hm-header__hamburger"
          aria-label={mobileMenuOpen.value ? "Close menu" : "Open menu"}
          aria-expanded={mobileMenuOpen.value}
          onClick$={() => {
            mobileMenuOpen.value = !mobileMenuOpen.value;
          }}
        >
          <span
            class={`hm-header__hamburger-icon ${mobileMenuOpen.value ? "hm-header__hamburger-icon--open" : ""}`}
          />
        </button>
      </div>

      {/* Mobile Menu */}
      {mobileMenuOpen.value && (
        <nav class="hm-header__mobile-menu" aria-label="Mobile navigation">
          {NAV_LINKS.map((link) => (
            <a
              key={link.href}
              href={link.href}
              class="hm-header__mobile-link"
            >
              {link.label}
            </a>
          ))}
          <div class="hm-header__mobile-actions">
            <a href="/sign-in" class="hm-header__mobile-link">
              Sign In
            </a>
            <Button variant="primary" size="md" href="/get-started">
              Get Started
            </Button>
          </div>
        </nav>
      )}
    </header>
  );
});
