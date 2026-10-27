import { shapeClass, type Shape } from "./lib/shape";
import { createEffect, createMemo, createSignal, For, Show } from "solid-js";
import CopyIcon from "~icons/lucide/copy";
import { Button } from "./button";
import { toast } from "./toast";
import { cn } from "./lib/cn";

export interface CodeBlockProps {
  shape?: Shape;
  code: string;
  filename?: string;
  showLineNumbers?: boolean;
  /** Format JS/TS examples. Invalid or unsupported source is displayed unchanged. */
  format?: boolean;
  /** Maximum line length before wrapping in prettier formatting. Defaults to 60. */
  printWidth?: number;
  class?: string;
}

let formatter:
  | Promise<{
      prettier: typeof import("prettier/standalone");
      plugins: import("prettier").Plugin[];
    }>
  | undefined;

async function formatCode(code: string, filename: string, printWidth = 60) {
  if (!/\.[cm]?[jt]sx?$/i.test(filename)) return code;
  formatter ??= Promise.all([
    import("prettier/standalone"),
    import("prettier/plugins/typescript"),
    import("prettier/plugins/estree"),
  ])
    .then(([prettier, typescript, estree]) => ({
      prettier,
      plugins: [typescript.default, estree.default],
    }))
    .catch((error) => {
      formatter = undefined;
      throw error;
    });
  const { prettier, plugins } = await formatter;
  return prettier.format(code, {
    parser: "typescript",
    plugins,
    filepath: filename,
    printWidth,
    tabWidth: 2,
  });
}

function tokenClass(value: string) {
  if (value.startsWith("//") || value.startsWith("/*"))
    return "text-muted-foreground italic";
  if (/^["'`]/.test(value))
    return "text-[color:var(--code-string,var(--foreground))]";
  if (
    /^(import|from|export|default|function|return|const|let|var|async|await|if|else|true|false|null|undefined|type|interface|extends|new|throw|try|catch|class|as|typeof)$/.test(
      value,
    )
  )
    return "text-[color:var(--code-keyword,var(--foreground))]";
  if (/^[A-Z]/.test(value))
    return "text-[color:var(--code-component,var(--foreground))]";
  if (/^\d/.test(value))
    return "text-[color:var(--code-number,var(--foreground))]";
  return "";
}

/** Keep multiline comments and strings intact while splitting tokens into numbered lines. */
function highlightLines(code: string) {
  const lines: {
    number: number;
    tokens: { id: number; value: string; class: string }[];
  }[] = [{ number: 1, tokens: [] }];
  const tokens = code.split(
    /(\/\/[^\n]*|\/\*[\s\S]*?\*\/|"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|`(?:\\.|[^`\\])*`|\b(?:import|from|export|default|function|return|const|let|var|async|await|if|else|true|false|null|undefined|type|interface|extends|new|throw|try|catch|class|as|typeof)\b|\b\d+(?:\.\d+)?\b|\b[A-Z][A-Za-z0-9]*\b)/g,
  );
  for (const token of tokens) {
    const parts = token.split("\n");
    parts.forEach((value, index) => {
      if (index) lines.push({ number: lines.length + 1, tokens: [] });
      if (value) {
        const line = lines.at(-1)!;
        line.tokens.push({
          id: line.tokens.length,
          value,
          class: tokenClass(token),
        });
      }
    });
  }
  return lines;
}

export function CodeBlock(props: CodeBlockProps) {
  const [formatted, setFormatted] = createSignal<{
    source: string;
    filename: string;
    code: string;
  }>();
  const filename = () => props.filename || "example.tsx";
  createEffect(
    () =>
      [
        props.code,
        filename(),
        props.format !== false,
        props.printWidth ?? 60,
      ] as const,
    ([source, file, shouldFormat, width]) => {
      if (!shouldFormat) {
        setFormatted(undefined);
        return;
      }
      let cancelled = false;
      void formatCode(source, file, width)
        .then((code) => {
          if (!cancelled) setFormatted({ source, filename: file, code });
        })
        .catch(() => {
          /* Incomplete examples remain readable while editing. */
        });
      return () => {
        cancelled = true;
      };
    },
  );
  const displayed = createMemo(() => {
    const result = formatted();
    const code =
      props.format !== false &&
      result?.source === props.code &&
      result.filename === filename()
        ? result.code
        : props.code;
    return code.replace(/\r\n?/g, "\n").replace(/\n$/, "");
  });
  const lines = createMemo(() => highlightLines(displayed()));
  const copy = async () => {
    try {
      await navigator.clipboard.writeText(displayed());
      toast.success("Código copiado", { description: filename() });
    } catch {
      toast.error("Não foi possível copiar", {
        description: "Selecione o código e copie manualmente.",
      });
    }
  };
  return (
    <figure
      class={cn(
        "m-0 w-full min-w-0 overflow-hidden rounded-lg border border-border bg-[var(--code-background,var(--muted))] text-foreground",
        shapeClass(props.shape, "surface"),
        props.class,
      )}
    >
      <figcaption class="flex min-h-11 items-center justify-between gap-3 border-b border-border bg-muted/30 px-4 text-xs text-muted-foreground">
        <span class="min-w-0 truncate font-mono" title={filename()}>
          {filename()}
        </span>
        <Button
          variant="ghost"
          size="icon"
          onClick={copy}
          aria-label="Copiar código"
          title="Copiar código"
          class="size-8 shrink-0"
        >
          <CopyIcon class="size-4" />
        </Button>
      </figcaption>
      <pre
        tabindex={0}
        aria-label={`Código de ${filename()}`}
        class="m-0 overflow-x-auto py-4 font-mono text-[13px] leading-6 outline-none [tab-size:2] focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring"
      >
        <code class="grid w-max min-w-full">
          <For each={lines()} keyed={(line) => line.number}>
            {(line) => (
              <span class="flex" data-code-line={line().number}>
                <Show when={props.showLineNumbers !== false}>
                  <span
                    aria-hidden="true"
                    class="sticky left-0 shrink-0 select-none border-r border-border/50 bg-[var(--code-background,var(--muted))] px-3 text-right text-[11px] text-muted-foreground"
                    style={{
                      width: `${Math.max(2, String(lines().length).length) * 8 + 24}px`,
                    }}
                  >
                    {line().number}
                  </span>
                </Show>
                <span class="flex-1 px-4 whitespace-pre">
                  <For each={line().tokens} keyed={(token) => token.id}>
                    {(token) => (
                      <span class={token().class}>{token().value}</span>
                    )}
                  </For>
                </span>
              </span>
            )}
          </For>
        </code>
      </pre>
    </figure>
  );
}
