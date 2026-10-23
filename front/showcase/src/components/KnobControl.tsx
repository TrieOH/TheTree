import { Show, type Component } from "solid-js";
import { Checkbox, Combobox, Field, Input, Textarea } from "@trieoh/ui-solid";

interface KnobSwitchProps {
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  description?: string;
}

export const KnobSwitch: Component<KnobSwitchProps> = (props) => (
  <Checkbox
    label={props.label}
    description={props.description}
    checked={props.checked}
    onChange={props.onChange}
    class="py-3 border-b border-border/50 last:border-b-0"
  />
);

interface KnobSelectProps<T extends string> {
  label: string;
  value: T;
  options: readonly T[];
  onChange: (value: T) => void;
}

export function KnobSelect<T extends string>(props: KnobSelectProps<T>) {
  return (
    <div class="py-3 border-b border-border/50 last:border-b-0">
      <Field label={props.label}>
        {(ids) => (
          <Combobox
            id={ids.id}
            aria-describedby={ids.describedBy}
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
    </div>
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
  <div class="py-3 border-b border-border/50 last:border-b-0">
    <Field label={props.label}>
      {(ids) => (
        <Show
          when={props.multiline}
          fallback={
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              value={props.value}
              placeholder={props.placeholder}
              onInput={(event) => props.onChange(event.currentTarget.value)}
            />
          }
        >
          <Textarea
            id={ids.id}
            aria-describedby={ids.describedBy}
            value={props.value}
            rows={10}
            class="h-auto min-h-48 resize-y font-mono text-xs"
            onInput={(event) => props.onChange(event.currentTarget.value)}
          />
        </Show>
      )}
    </Field>
  </div>
);
