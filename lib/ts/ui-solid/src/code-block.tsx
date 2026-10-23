import { createMemo, For } from "solid-js";
import { Button } from "./button";
import { toast } from "./toast";
import { cn } from "./lib/cn";

export interface CodeBlockProps {
  code: string;
  filename?: string;
  class?: string;
}

function tokenClass(value: string) {
  if (value.startsWith("//") || value.startsWith("/*"))
    return "text-muted-foreground italic";
  if (/^["'`]/.test(value)) return "text-code-string";
  if (
    /^(import|from|export|function|return|const|let|async|await|if|else|true|false|null|undefined)$/.test(
      value,
    )
  )
    return "text-code-keyword";
  if (/^[A-Z]/.test(value)) return "text-code-component";
  if (/^\d+$/.test(value)) return "text-code-number";
  return "";
}

// Lightweight lexical color, rendered as text nodes so code is never interpreted as HTML.
function tokens(code: string) {
  return code
    .split(
      /(\/\/[^\n]*|\/\*[\s\S]*?\*\/|"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|`(?:\\.|[^`\\])*`|\b(?:import|from|export|function|return|const|let|async|await|if|else|true|false|null|undefined)\b|\b\d+\b|\b[A-Z][A-Za-z0-9]*\b)/g,
    )
    .filter(Boolean)
    .map((value, index) => ({
      id: index,
      value,
      class: tokenClass(value),
    }));
}

export function CodeBlock(props: CodeBlockProps) {
  const highlighted = createMemo(() => tokens(props.code));
  const copy = async () => {
    try {
      await navigator.clipboard.writeText(props.code);
      toast.success("Código copiado");
    } catch {
      toast.error("Não foi possível copiar o código");
    }
  };
  return (
    <figure
      class={cn(
        "min-w-0 overflow-hidden rounded-lg border border-border bg-code-background text-foreground",
        props.class,
      )}
    >
      <figcaption class="flex min-h-11 items-center justify-between gap-3 border-b border-border px-4 text-xs text-muted-foreground">
        <span class="font-mono">{props.filename ?? "example.tsx"}</span>
        <Button
          variant="ghost"
          size="sm"
          onClick={copy}
          aria-label="Copiar código"
          class="gap-1.5 text-xs"
        >
          <svg
            aria-hidden="true"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
            class="size-3.5"
          >
            <rect x="8" y="8" width="12" height="12" rx="2" />
            <path d="M16 8V4a1 1 0 0 0-1-1H4a1 1 0 0 0-1 1v11a1 1 0 0 0 1 1h4" />
          </svg>
          <span>Copiar</span>
        </Button>
      </figcaption>
      <pre
        tabindex={0}
        aria-label="Exemplo de código"
        class="overflow-x-auto px-5 py-5 font-mono text-[13px] leading-6 outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring"
      >
        <code>
          <For each={highlighted()} keyed={(token) => token.id}>
            {(token) => <span class={token().class}>{token().value}</span>}
          </For>
        </code>
      </pre>
    </figure>
  );
}
