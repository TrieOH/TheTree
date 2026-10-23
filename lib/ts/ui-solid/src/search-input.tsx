import { Input, type InputProps } from "./input";
import { cn } from "./lib/cn";

export interface SearchInputProps {
  value: string;
  onInput: NonNullable<InputProps["onInput"]>;
  placeholder?: string;
  "aria-label": string;
  class?: string;
}

export function SearchInput(props: SearchInputProps) {
  return (
    <div class={cn("relative w-full", props.class)}>
      <svg
        aria-hidden="true"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
      >
        <circle cx="11" cy="11" r="8" />
        <path d="m21 21-4.3-4.3" />
      </svg>
      <Input
        type="search"
        value={props.value}
        onInput={props.onInput}
        placeholder={props.placeholder}
        aria-label={props["aria-label"]}
        class="pl-9"
      />
    </div>
  );
}
