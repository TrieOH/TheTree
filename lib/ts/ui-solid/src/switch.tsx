import { cn } from "./lib/cn";

export interface SwitchProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  "aria-label"?: string;
  "aria-labelledby"?: string;
  "aria-describedby"?: string;
  disabled?: boolean;
  required?: boolean;
  id?: string;
  name?: string;
  value?: string;
  class?: string;
}

export function Switch(props: SwitchProps) {
  return (
    <input
      type="checkbox"
      role="switch"
      id={props.id}
      name={props.name}
      value={props.value}
      checked={props.checked}
      disabled={props.disabled}
      required={props.required}
      aria-label={props["aria-label"]}
      aria-labelledby={props["aria-labelledby"]}
      aria-describedby={props["aria-describedby"]}
      onChange={(event) => props.onChange(event.currentTarget.checked)}
      class={cn(
        "relative m-0 inline-block h-[18px] w-[30px] shrink-0 appearance-none rounded-full border border-transparent bg-input align-middle transition-colors checked:bg-primary",
        "before:absolute before:left-[2px] before:top-[2px] before:size-3 before:rounded-full before:bg-background before:shadow-sm before:transition-transform before:content-[''] checked:before:translate-x-3 checked:before:bg-primary-foreground",
        "focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-ring disabled:cursor-not-allowed disabled:opacity-50 motion-reduce:transition-none motion-reduce:before:transition-none forced-colors:appearance-auto",
        props.class,
      )}
    />
  );
}
