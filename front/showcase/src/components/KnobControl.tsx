import { Show, type Component } from "solid-js";
import { Combobox, Field, Input, Switch, Textarea } from "@trieoh/ui-solid";

interface KnobSwitchProps {
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
}

export const KnobSwitch: Component<KnobSwitchProps> = (props) => (
  <label class="flex min-h-11 cursor-pointer items-center justify-between gap-6 py-2">
    <span class="text-xs">{props.label}</span>
    <Switch
      aria-label={props.label}
      checked={props.checked}
      disabled={props.disabled}
      onChange={props.onChange}
    />
  </label>
);

interface KnobSelectProps<T extends string> {
  label: string;
  value: T;
  options: readonly T[];
  onChange: (value: T) => void;
}

export function KnobSelect<T extends string>(props: KnobSelectProps<T>) {
  return (
    <Field
      label={props.label}
      class="min-w-0 [&_label]:text-xs [&_label]:font-normal [&_label]:text-muted-foreground"
    >
      {(ids) => (
        <Combobox
          shape="soft"
          triggerClass="border-input shadow-xs"
          id={ids.id}
          aria-describedby={ids.describedBy}
          size="sm"
          aria-label={props.label}
          placeholder={props.label}
          value={props.value}
          options={props.options.map((value) => ({
            value,
            label: value || "(vazio)",
          }))}
          onChange={(value) => props.onChange(value as T)}
        />
      )}
    </Field>
  );
}

interface KnobTextProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  multiline?: boolean;
}

export const KnobText: Component<KnobTextProps> = (props) => (
  <div
    class={
      props.multiline
        ? "min-w-0 [&_label]:text-xs [&_label]:font-normal [&_label]:text-muted-foreground col-span-full"
        : "min-w-0 [&_label]:text-xs [&_label]:font-normal [&_label]:text-muted-foreground"
    }
  >
    <Field label={props.label}>
      {(ids) => (
        <Show
          when={props.multiline}
          fallback={
            <Input
              shape="soft"
              class="h-8 bg-background text-xs"
              id={ids.id}
              aria-describedby={ids.describedBy}
              value={props.value}
              placeholder={props.placeholder}
              onInput={(event) => props.onChange(event.currentTarget.value)}
            />
          }
        >
          <Textarea
            shape="soft"
            id={ids.id}
            aria-describedby={ids.describedBy}
            value={props.value}
            rows={10}
            class="h-auto min-h-48 resize-y font-mono text-[11px]"
            onInput={(event) => props.onChange(event.currentTarget.value)}
          />
        </Show>
      )}
    </Field>
  </div>
);
