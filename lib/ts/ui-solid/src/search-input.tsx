import SearchIcon from "~icons/lucide/search";
import { Input, type InputProps } from "./input";
import { cn } from "./lib/cn";

export interface SearchInputProps {
  shape?: InputProps["shape"];
  disabled?: boolean;
  required?: boolean;
  readonly?: boolean;
  id?: string;
  name?: string;
  value: string;
  onInput: NonNullable<InputProps["onInput"]>;
  placeholder?: string;
  "aria-label": string;
  class?: string;
}

export function SearchInput(props: SearchInputProps) {
  return (
    <div class={cn("relative w-full", props.class)}>
      <SearchIcon
        class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
        aria-hidden="true"
      />
      <Input
        type="search"
        shape={props.shape}
        disabled={props.disabled}
        required={props.required}
        readonly={props.readonly}
        id={props.id}
        name={props.name}
        value={props.value}
        onInput={props.onInput}
        placeholder={props.placeholder}
        aria-label={props["aria-label"]}
        class="pl-9"
      />
    </div>
  );
}
