import { cn } from "@trieoh/ui-solid";
import { createSignal } from "solid-js";
import CheckIcon from "~icons/lucide/check";
import CopyIcon from "~icons/lucide/copy";

/** Copies `text()` (read when clicked) and shows a check for a moment. */
export function CopyButton(props: { text: () => string; label?: string; class?: string }) {
  const [copied, setCopied] = createSignal(false, { name: "copied" });
  const copy = async () => {
    await navigator.clipboard.writeText(props.text());
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };
  return (
    <button
      type="button"
      onClick={copy}
      aria-label={copied() ? "Copied" : (props.label ?? "Copy")}
      class={cn(
        "grid size-7 shrink-0 place-items-center rounded-md text-muted-foreground transition-colors hover:bg-foreground/5 hover:text-foreground",
        props.class,
      )}
    >
      {copied() ? <CheckIcon class="size-3.5" /> : <CopyIcon class="size-3.5" />}
    </button>
  );
}
