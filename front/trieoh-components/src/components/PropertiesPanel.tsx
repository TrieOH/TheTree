import { For, Show } from "solid-js";
import { Button, Tooltip } from "@trieoh/ui-solid";
import IconRotateCcw from "~icons/lucide/rotate-ccw";
import type { Entry, Values } from "../catalog";
import { KnobSelect, KnobSwitch, KnobText } from "./KnobControl";

export function PropertiesPanel(props: {
  entry: Entry;
  values: Values;
  onUpdate: (key: string, value: string | boolean) => void;
  onReset: () => void;
}) {
  const fields = () =>
    props.entry.controls.filter(
      (control) => typeof props.entry.defaults[control.key] !== "boolean",
    );
  const switches = () =>
    props.entry.controls.filter(
      (control) => typeof props.entry.defaults[control.key] === "boolean",
    );

  return (
    <section
      class="mt-6 min-w-0 rounded-lg border border-border bg-card"
      aria-label="Inspetor de componentes"
    >
      <div class="flex min-h-13.5 items-center justify-between gap-3 border-b border-border px-5">
        <span class="text-[13px] font-semibold">Propriedades</span>
        <Tooltip align="end" label="Restaurar valores padrão">
          <Button
            variant="ghost"
            size="icon"
            class="size-7 text-muted-foreground"
            aria-label="Restaurar valores padrão"
            disabled={Object.keys(props.entry.defaults).every(
              (key) => props.values[key] === props.entry.defaults[key],
            )}
            onClick={props.onReset}
          >
            <IconRotateCcw class="size-3.5" />
          </Button>
        </Tooltip>
      </div>

      <Show when={fields().length > 0}>
        <div class="grid grid-cols-1 items-start gap-x-6 gap-y-5 p-5 min-[761px]:grid-cols-2">
          <For each={fields()} keyed={(control) => control.key}>
            {(control) => (
              <Show
                when={control().options}
                fallback={
                  <KnobText
                    multiline={control().key === "code"}
                    label={control().label}
                    value={String(props.values[control().key])}
                    onChange={(value) => props.onUpdate(control().key, value)}
                  />
                }
              >
                <KnobSelect
                  label={control().label}
                  value={String(props.values[control().key])}
                  options={control().options!}
                  onChange={(value) => props.onUpdate(control().key, value)}
                />
              </Show>
            )}
          </For>
        </div>
      </Show>
      <Show when={switches().length > 0}>
        <fieldset
          class={
            fields().length > 0
              ? "border-t border-border px-5 py-2"
              : "px-5 py-2"
          }
        >
          <legend class="sr-only">Estados</legend>
          <div class="divide-y divide-border/60">
            <For each={switches()} keyed={(control) => control.key}>
              {(control) => (
                <KnobSwitch
                  label={control().label}
                  checked={!!props.values[control().key]}
                  onChange={(value) => props.onUpdate(control().key, value)}
                />
              )}
            </For>
          </div>
        </fieldset>
      </Show>
    </section>
  );
}
