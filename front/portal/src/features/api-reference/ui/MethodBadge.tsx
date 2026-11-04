import { cn } from "@trieoh/ui-solid";

const TONE: Record<string, string> = {
  get: "bg-ok/15 text-ok",
  post: "bg-accent/15 text-accent",
  put: "bg-warn/15 text-warn",
  patch: "bg-warn/15 text-warn",
  delete: "bg-destructive/20 text-destructive dark:text-destructive-foreground",
};

export function MethodBadge(props: { method: string; class?: string }) {
  return (
    <span
      class={cn(
        "inline-flex h-5 w-14 shrink-0 items-center justify-center rounded font-mono text-[10.5px] font-semibold uppercase",
        TONE[props.method] ?? "bg-muted text-muted-foreground",
        props.class,
      )}
    >
      {props.method}
    </span>
  );
}

/** `/wallets/{wallet_id}` with the path parameters picked out. */
export function PathText(props: { path: string }) {
  return (
    <span class="min-w-0 break-all font-mono text-[13px]">
      {props.path.split(/(\{[^}]+\})/).map((part) =>
        part.startsWith("{") ? <span class="text-accent">{part}</span> : part,
      )}
    </span>
  );
}
