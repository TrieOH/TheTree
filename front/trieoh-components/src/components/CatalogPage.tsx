import { For, Show } from "solid-js";
import { Button, Card, Link, SearchInput } from "@trieoh/ui-solid";
import IconArrowUpRight from "~icons/lucide/arrow-up-right";
import type { Entry, Values } from "../catalog";
import { FitPreview } from "./FitPreview";

export function CatalogPage(props: {
  title: string;
  entries: readonly Entry[];
  query: string;
  values: Record<string, Values>;
  onQueryChange: (query: string) => void;
  onClearFilters: () => void;
  onUpdate: (name: string, key: string, value: string | boolean) => void;
}) {
  return (
    <main class="min-w-0 px-5 pt-7.5 min-[761px]:px-6.25 min-[761px]:pt-8 lg:px-10 lg:pt-10">
      <div class="mb-7">
        <h1 class="mb-2 text-[28px] leading-tight font-semibold tracking-[-0.7px]">
          {props.title}
        </h1>
        <p class="max-w-160 text-[13px] leading-[1.6] text-muted-foreground">
          A UI do Univents, componente por componente. Teste estados e
          variantes, ajuste as propriedades e copie o código em Solid.
        </p>
      </div>
      <div class="mb-6 flex items-center justify-between gap-5">
        <SearchInput
          class="max-w-87.5"
          aria-label="Buscar componentes"
          placeholder="Buscar componente…"
          value={props.query}
          onInput={(e) => props.onQueryChange(e.currentTarget.value)}
        />
        <span
          class="text-[11px] whitespace-nowrap text-muted-foreground max-[440px]:hidden"
          aria-live="polite"
        >
          {props.entries.length} componentes
        </span>
      </div>
      <div class="grid grid-cols-1 border-t border-border min-[641px]:grid-cols-[repeat(auto-fill,minmax(240px,1fr))] min-[641px]:gap-4.25 min-[641px]:border-t-0">
        <For each={props.entries}>
          {(entry, index) => (
            <Card class="block self-start overflow-hidden rounded-none border-0 border-b border-border bg-transparent shadow-none hover:border-muted-foreground min-[641px]:grid min-[641px]:aspect-square min-[641px]:grid-rows-[auto_minmax(0,1fr)_auto] min-[641px]:rounded-xl min-[641px]:border min-[641px]:bg-card min-[641px]:shadow-xs min-[641px]:contain-[size]">
              <header class="hidden h-10.75 items-center gap-2.5 px-3.5 min-[641px]:flex">
                <span class="font-mono text-[10px] text-muted-foreground opacity-60">
                  {String(index() + 1).padStart(2, "0")}
                </span>
                <span class="text-[10px] text-muted-foreground">
                  {entry.group}
                </span>
                <Link
                  class="ml-auto inline-flex rounded p-1 text-foreground no-underline hover:bg-muted hover:no-underline"
                  aria-label={`Ver documentação de ${entry.name}`}
                  href={`#/components/${entry.name}`}
                >
                  <IconArrowUpRight class="size-3.5" />
                </Link>
              </header>
              <div class="relative hidden min-h-0 overflow-hidden bg-background min-[641px]:block">
                <FitPreview>
                  {entry.render(props.values[entry.name], (key, value) =>
                    props.onUpdate(entry.name, key, value),
                  )}
                </FitPreview>
              </div>
              <Link
                class="flex h-auto w-full flex-col gap-1.5 overflow-hidden px-1 py-3.5 text-foreground no-underline hover:no-underline min-[641px]:h-21 min-[641px]:border-t min-[641px]:border-border min-[641px]:px-4 min-[641px]:py-3"
                href={`#/components/${entry.name}`}
              >
                <strong class="text-sm font-semibold wrap-anywhere min-[641px]:text-xs">
                  {entry.name}
                </strong>
                <span class="hidden text-[11px] leading-[1.6] text-muted-foreground min-[641px]:line-clamp-2">
                  {entry.description}
                </span>
              </Link>
            </Card>
          )}
        </For>
      </div>
      <Show when={props.entries.length === 0}>
        <div class="px-5 py-17.5 text-center text-[13px]">
          <h2>Nenhum componente encontrado</h2>
          <p class="mt-2.5 mb-5 text-muted-foreground">
            Tente outro nome ou categoria.
          </p>
          <Button variant="outline" onClick={props.onClearFilters}>
            Limpar filtros
          </Button>
        </div>
      </Show>
      <footer class="mt-9 flex justify-between gap-3 border-t border-border py-5.5 text-[10px] text-muted-foreground">
        <span>TrieOH Design System</span>
        <span>Solid · Tailwind CSS</span>
      </footer>
    </main>
  );
}
