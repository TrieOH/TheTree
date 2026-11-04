import { cn } from "@trieoh/ui-solid";

/** One color per method, used by badges and by the route list. */
export const METHOD_TEXT: Record<string, string> = {
  get: "text-ok",
  post: "text-accent",
  put: "text-warn",
  patch: "text-warn",
  delete: "text-destructive dark:text-destructive-foreground",
};

const METHOD_BG: Record<string, string> = {
  get: "bg-ok/15",
  post: "bg-accent/15",
  put: "bg-warn/15",
  patch: "bg-warn/15",
  delete: "bg-destructive/20",
};

export function MethodBadge(props: { method: string; class?: string }) {
  return (
    <span
      class={cn(
        "inline-flex h-6 shrink-0 items-center justify-center rounded-md px-2 font-mono text-[11px] font-semibold uppercase tracking-wide",
        METHOD_TEXT[props.method] ?? "text-muted-foreground",
        METHOD_BG[props.method] ?? "bg-muted",
        props.class,
      )}
    >
      {props.method}
    </span>
  );
}

/** `/wallets/{wallet_id}` on one line, path parameters picked out; scrolls sideways if long. */
export function PathText(props: { path: string; class?: string }) {
  return (
    <span class={cn("min-w-0 overflow-x-auto whitespace-nowrap font-mono text-[13px] [scrollbar-width:none]", props.class)}>
      {props.path.split(/(\{[^}]+\})/).map((part) =>
        part.startsWith("{") ? <span class="text-accent">{part}</span> : part,
      )}
    </span>
  );
}
