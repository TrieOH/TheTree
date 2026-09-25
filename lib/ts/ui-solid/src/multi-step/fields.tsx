import type { JSX } from "@solidjs/web";
import { For, Match, Switch } from "solid-js";
import { cn } from "../lib/cn";
import { MultiStepCustomField } from "./fields/custom-field";
import { MultiStepTextField } from "./fields/text-field";
import { MultiStepTextareaField } from "./fields/textarea-field";
import type {
  CustomMultiStepField,
  MultiStepFieldsProps,
  TextMultiStepField,
  TextareaMultiStepField,
} from "./types";

export { MultiStepTextField } from "./fields/text-field";
export { MultiStepTextareaField } from "./fields/textarea-field";
export { MultiStepCustomField } from "./fields/custom-field";

/**
 * Modular subcomponent for rendering configured form fields in a responsive 1/2-column grid.
 * Segregates field types cleanly using Solid Switch/Match.
 */
export function MultiStepFields<T = Record<string, unknown>>(
  props: MultiStepFieldsProps<T>,
): JSX.Element {
  return (
    <div class={cn("grid grid-cols-1 sm:grid-cols-2 gap-4 w-full max-w-full min-w-0", props.class)}>
      <For each={props.fields}>
        {(field) => {
          const isHalf = () => field.layout === "half";
          const errorMsg = () => props.errors?.[field.name];
          const val = () => (props.values as Record<string, unknown> | undefined)?.[field.name];

          return (
            <div class={cn("min-w-0 w-full", isHalf() ? "sm:col-span-1" : "sm:col-span-2")}>
              <Switch>
                <Match when={field.kind === "textarea"}>
                  <MultiStepTextareaField
                    field={field as TextareaMultiStepField<T>}
                    value={val()}
                    error={errorMsg()}
                    onChange={(v) => props.onChange?.(field.name, v)}
                  />
                </Match>

                <Match when={field.kind === "custom"}>
                  <MultiStepCustomField
                    field={field as CustomMultiStepField<T>}
                    value={val()}
                    error={errorMsg()}
                    onChange={(v) => props.onChange?.(field.name, v)}
                  />
                </Match>

                <Match when={true}>
                  <MultiStepTextField
                    field={field as TextMultiStepField<T>}
                    value={val()}
                    error={errorMsg()}
                    onChange={(v) => props.onChange?.(field.name, v)}
                  />
                </Match>
              </Switch>
            </div>
          );
        }}
      </For>
    </div>
  );
}
