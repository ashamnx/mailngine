import { component$ } from "@builder.io/qwik";

export interface CodeBlockProps {
  code: string;
  language?: string;
}

/**
 * Simple code block with CSS-based syntax highlighting.
 * Avoids heavy JS libraries to keep the bundle lean and resumable.
 */
export const CodeBlock = component$<CodeBlockProps>(
  ({ code, language = "plaintext" }) => {
    const highlighted = applySyntaxHighlighting(code);

    return (
      <div class="hm-code-block">
        {language !== "plaintext" && (
          <div class="hm-code-block__lang">{language}</div>
        )}
        <pre class="hm-code-block__pre">
          <code dangerouslySetInnerHTML={highlighted} />
        </pre>
      </div>
    );
  },
);

/**
 * Lightweight syntax coloring. Runs at SSR time (no client JS cost).
 * Keywords = blue, strings = green, comments = gray.
 */
function applySyntaxHighlighting(code: string): string {
  let escaped = code
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");

  // Single-line comments: // ...
  escaped = escaped.replace(
    /(\/\/.*)/g,
    '<span class="hm-code--comment">$1</span>',
  );

  // Multi-line comments: /* ... */
  escaped = escaped.replace(
    /(\/\*[\s\S]*?\*\/)/g,
    '<span class="hm-code--comment">$1</span>',
  );

  // Strings: double-quoted and single-quoted
  escaped = escaped.replace(
    /("(?:[^"\\]|\\.)*")/g,
    '<span class="hm-code--string">$1</span>',
  );
  escaped = escaped.replace(
    /('(?:[^'\\]|\\.)*')/g,
    '<span class="hm-code--string">$1</span>',
  );

  // Template literals
  escaped = escaped.replace(
    /(`(?:[^`\\]|\\.)*`)/g,
    '<span class="hm-code--string">$1</span>',
  );

  // Keywords
  const keywords =
    /\b(const|let|var|function|return|import|export|from|if|else|for|while|class|new|this|async|await|try|catch|throw|default|type|interface)\b/g;
  escaped = escaped.replace(
    keywords,
    '<span class="hm-code--keyword">$1</span>',
  );

  return escaped;
}
