import { Show, createSignal, type Component } from "solid-js";
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

interface KnobNumberProps {
  label: string;
  value: string;
  min?: number;
  max?: number;
  step?: number;
  /** An empty field is a valid answer (the prop is simply not set). */
  optional?: boolean;
  onChange: (value: string) => void;
}

/**
 * Numeric knob. The browser's number input keeps most letters out, the keys it
 * still lets through ("e", "+", and "-" when negatives make no sense) are
 * blocked, and the range is enforced without fighting the user: a value is only
 * sent to the preview while it is inside [min, max], an out-of-range draft is
 * flagged, and leaving the field snaps it to the nearest allowed value.
 */
export const KnobNumber: Component<KnobNumberProps> = (props) => {
  // undefined mirrors the current value; a string is what the user is typing.
  const [draft, setDraft] = createSignal<string | undefined>();
  const shown = () => draft() ?? props.value;

  const parsed = (raw: string) => (raw.trim() === "" ? undefined : Number(raw));
  const inRange = (value: number) =>
    Number.isFinite(value) &&
    (props.min === undefined || value >= props.min) &&
    (props.max === undefined || value <= props.max);
  const clamp = (value: number) =>
    Math.min(props.max ?? Infinity, Math.max(props.min ?? -Infinity, value));

  const invalid = () => {
    const value = parsed(shown());
    return value === undefined ? !props.optional : !inRange(value);
  };

  const hint = () => {
    if (props.min !== undefined && props.max !== undefined)
      return `Entre ${props.min} e ${props.max}`;
    if (props.min !== undefined) return `Mínimo ${props.min}`;
    if (props.max !== undefined) return `Máximo ${props.max}`;
    return undefined;
  };

  const commit = () => {
    const raw = draft();
    setDraft(undefined);
    if (raw === undefined) return;

    const value = parsed(raw);
    if (value === undefined) {
      if (props.optional) props.onChange("");
      return;
    }
    if (Number.isFinite(value)) props.onChange(String(clamp(value)));
  };

  return (
    <div class="min-w-0 [&_label]:text-xs [&_label]:font-normal [&_label]:text-muted-foreground">
      <Field
        label={props.label}
        hint={invalid() ? undefined : hint()}
        error={invalid() ? (hint() ?? "Valor inválido") : undefined}
      >
        {(ids) => (
          <Input
            shape="soft"
            type="number"
            inputmode="decimal"
            class="h-8 bg-background text-xs"
            id={ids.id}
            aria-describedby={ids.describedBy}
            aria-invalid={invalid()}
            min={props.min}
            max={props.max}
            step={props.step}
            value={shown()}
            onKeyDown={(event) => {
              if (
                event.key === "e" ||
                event.key === "E" ||
                event.key === "+" ||
                (event.key === "-" && props.min !== undefined && props.min >= 0)
              )
                event.preventDefault();
            }}
            onInput={(event) => {
              const raw = event.currentTarget.value;
              setDraft(raw);
              const value = parsed(raw);
              if (value !== undefined && inRange(value)) props.onChange(String(value));
              else if (value === undefined && props.optional) props.onChange("");
            }}
            onBlur={commit}
          />
        )}
      </Field>
    </div>
  );
};
