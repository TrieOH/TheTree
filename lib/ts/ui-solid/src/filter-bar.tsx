import XIcon from "~icons/lucide/x";
import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { Button } from "./button";
import { cn } from "./lib/cn";

export interface ActiveFilter {
  id: string;
  label: string;
  value: string;
}
export interface FilterBarProps {
  filters: readonly ActiveFilter[];
  onRemove: (id: string) => void;
  onClear: () => void;
  /** Optional controls for adding filters. */
  children?: JSX.Element;
  disabled?: boolean;
  class?: string;
}
export function FilterBar(props: FilterBarProps) {
  return (
    <div class={cn("flex min-w-0 flex-wrap items-center gap-2", props.class)}>
      {props.children}
      <Show when={props.filters.length}>
        <div
          role="group"
          aria-label="Filtros ativos"
          class="flex min-w-0 flex-wrap items-center gap-2"
        >
          <For each={props.filters}>
            {(filter) => (
              <Button
                size="sm"
                class="h-auto min-h-8 max-w-full whitespace-normal py-1 text-left wrap-anywhere"
                variant="outline"
                disabled={props.disabled}
                aria-label={`Remover filtro ${filter.label}: ${filter.value}`}
                onClick={() => props.onRemove(filter.id)}
              >
                {filter.label}: {filter.value}
                <XIcon class="size-3 shrink-0" aria-hidden="true" />
              </Button>
            )}
          </For>
          <Button
            size="sm"
            variant="ghost"
            disabled={props.disabled}
            onClick={() => props.onClear()}
          >
            Limpar todos
          </Button>
        </div>
      </Show>
    </div>
  );
}
