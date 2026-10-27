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
  return (
    <section
      class="mt-6 min-w-0 rounded-lg border border-border bg-card px-4 pb-1.5"
      aria-label="Inspetor de componentes"
    >
      <div class="flex min-h-13.5 items-center justify-between gap-3">
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

      <div class="grid grid-cols-1 items-start gap-x-8 border-t border-border pt-1.5 min-[761px]:grid-cols-[repeat(auto-fit,minmax(240px,1fr))]">
        <For each={props.entry.controls}>
          {(control) => (
            <Show
              when={typeof props.values[control.key] === "boolean"}
              fallback={
                <Show
                  when={control.options}
                  fallback={
                    <KnobText
                      multiline={control.key === "code"}
                      label={control.label}
                      value={String(props.values[control.key])}
                      onChange={(value) => props.onUpdate(control.key, value)}
                    />
                  }
                >
                  <KnobSelect
                    label={control.label}
                    value={String(props.values[control.key])}
                    options={control.options!}
                    onChange={(value) => props.onUpdate(control.key, value)}
                  />
                </Show>
              }
            >
              <KnobSwitch
                label={control.label}
                checked={!!props.values[control.key]}
                onChange={(value) => props.onUpdate(control.key, value)}
              />
            </Show>
          )}
        </For>
      </div>
    </section>
  );
}
