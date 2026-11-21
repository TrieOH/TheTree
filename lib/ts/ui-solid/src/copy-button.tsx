import CheckIcon from "~icons/lucide/check";
import CopyIcon from "~icons/lucide/copy";
import { createSignal, onCleanup } from "solid-js";
import { Button } from "./button";
import { cn } from "./lib/cn";
import { type Shape } from "./lib/shape";

export interface CopyButtonProps {
  /** Text to copy; a function is read only when clicked. */
  text: string | (() => string);
  label?: string;
  copiedLabel?: string;
  size?: "sm" | "md";
  shape?: Shape;
  onCopied?: (text: string) => void;
  onError?: (error: unknown) => void;
  class?: string;
}

const sizes = {
  sm: { button: "size-7", icon: "size-3.5" },
  md: { button: "size-8", icon: "size-4" },
};

/** Copies to the clipboard and shows a check for a moment. */
export function CopyButton(props: CopyButtonProps) {
  const [copied, setCopied] = createSignal(false, { name: "copied" });
  let timer: ReturnType<typeof setTimeout> | undefined;
  onCleanup(() => clearTimeout(timer));
  const size = () => sizes[props.size ?? "md"];
  const label = () =>
    copied() ? (props.copiedLabel ?? "Copiado") : (props.label ?? "Copiar");
  const copy = async () => {
    const text = typeof props.text === "function" ? props.text() : props.text;
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      clearTimeout(timer);
      timer = setTimeout(() => setCopied(false), 1500);
      props.onCopied?.(text);
    } catch (error) {
      props.onError?.(error);
    }
  };
  return (
    <Button
      variant="ghost"
      size="icon"
      shape={props.shape}
      onClick={copy}
      aria-label={label()}
      title={label()}
      class={cn("shrink-0 text-muted-foreground", size().button, props.class)}
    >
      {copied() ? (
        <CheckIcon class={size().icon} aria-hidden="true" />
      ) : (
        <CopyIcon class={size().icon} aria-hidden="true" />
      )}
    </Button>
  );
}
