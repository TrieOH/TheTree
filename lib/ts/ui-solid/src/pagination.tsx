import ChevronLeft from "~icons/lucide/chevron-left";
import ChevronRight from "~icons/lucide/chevron-right";
import ChevronsLeft from "~icons/lucide/chevrons-left";
import ChevronsRight from "~icons/lucide/chevrons-right";
import { For, Show, Switch, Match, createMemo } from "solid-js";
import { Button } from "./button";
import { Combobox } from "./combobox";
import { cn } from "./lib/cn";

export interface PaginationProps {
  /** One-based, controlled page. */
  page: number;
  onPageChange: (page: number) => void;
  pageSize?: number;
  /** Omit for cursor-based navigation with no known total. */
  totalItems?: number;
  hasNextPage?: boolean;
  hasPreviousPage?: boolean;
  pageSizeOptions?: readonly number[];
  onPageSizeChange?: (size: number) => void;
  /** Defaults to true when a change handler is supplied. */
  showPageSizeSelector?: boolean;
  disabled?: boolean;
  itemLabel?: string;
  class?: string;
  "aria-label"?: string;
}

/** Bounded window: first, last and at most five neighbouring pages. */
function pageNumbers(
  page: number,
  total: number,
): (number | "start-gap" | "end-gap")[] {
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);
  const start = Math.max(2, Math.min(page - 2, total - 5));
  const end = Math.min(total - 1, Math.max(page + 2, 6));
  return [
    1,
    ...(start > 2 ? ["start-gap" as const] : []),
    ...Array.from({ length: end - start + 1 }, (_, i) => start + i),
    ...(end < total - 1 ? ["end-gap" as const] : []),
    total,
  ];
}

export function Pagination(props: PaginationProps) {
  const size = () => Math.max(1, Math.floor(props.pageSize ?? 10));
  const total = () => Math.max(0, props.totalItems ?? 0);
  const pages = () => Math.max(1, Math.ceil(total() / size()));
  const current = () =>
    Math.max(
      1,
      Math.min(
        Math.floor(props.page),
        props.totalItems === undefined ? Infinity : pages(),
      ),
    );
  const previous = () => props.hasPreviousPage ?? current() > 1;
  const next = () =>
    props.hasNextPage ??
    (props.totalItems !== undefined && current() < pages());
  const entries = createMemo(() => pageNumbers(current(), pages()));
  const go = (page: number) => {
    if (!props.disabled)
      props.onPageChange(
        Math.max(
          1,
          Math.min(page, props.totalItems === undefined ? Infinity : pages()),
        ),
      );
  };
  return (
    <nav
      aria-label={props["aria-label"] ?? "Paginação"}
      class={cn(
        "@container/pagination flex min-w-0 flex-wrap items-center justify-between gap-3",
        props.class,
      )}
    >
      <div class="flex min-w-0 flex-1 basis-full flex-wrap items-center justify-between gap-x-6 gap-y-3 text-xs text-muted-foreground @[800px]/pagination:basis-0">
        <span class="min-w-0 wrap-anywhere" aria-live="polite">
          <Show
            when={props.totalItems !== undefined}
            fallback={<>Página {current()}</>}
          >
            Mostrando {total() === 0 ? 0 : (current() - 1) * size() + 1}–
            {Math.min(current() * size(), total())} de {total()}{" "}
            {props.itemLabel ?? "itens"}
          </Show>
        </span>
        <Show
          when={props.onPageSizeChange && props.showPageSizeSelector !== false}
        >
          <Combobox
            size="sm"
            shape="soft"
            triggerClass="border-input shadow-xs"
            searchable={false}
            aria-label="Itens por página"
            placeholder="Itens por página"
            class="ml-auto w-32 max-w-full"
            value={String(size())}
            options={[
              ...new Set([
                ...(props.pageSizeOptions ?? [10, 25, 50, 100]),
                size(),
              ]),
            ]
              .filter((n) => Number.isInteger(n) && n > 0)
              .sort((a, b) => a - b)
              .map((value) => ({
                value: String(value),
                label: `${value} por página`,
              }))}
            disabled={props.disabled}
            onChange={(value) => props.onPageSizeChange?.(Number(value))}
          />
        </Show>
      </div>
      <div class="ml-auto flex shrink-0 items-center justify-end gap-1">
        <Show when={props.totalItems !== undefined}>
          <Button
            size="sm"
            variant="outline"
            shape="subtle"
            class="size-7 p-0 text-muted-foreground"
            aria-label="Primeira página"
            disabled={props.disabled || current() === 1}
            onClick={() => go(1)}
          >
            <PageArrow direction="first" />
          </Button>
        </Show>
        <Button
          size="sm"
          variant="outline"
          shape="subtle"
          disabled={props.disabled || !previous()}
          class="size-7 p-0 text-muted-foreground"
          aria-label="Página anterior"
          onClick={() => go(current() - 1)}
        >
          <PageArrow direction="previous" />
        </Button>
        <Show when={props.totalItems !== undefined}>
          <div class="hidden items-center gap-1 @[480px]/pagination:flex">
            <For each={entries()}>
              {(entry) => (
                <Show
                  when={typeof entry === "number"}
                  fallback={
                    <span class="px-1 text-muted-foreground" aria-hidden="true">
                      …
                    </span>
                  }
                >
                  <Button
                    size="sm"
                    variant="outline"
                    shape="subtle"
                    class={cn(
                      "h-7 min-w-7 px-1.5 text-xs text-muted-foreground",
                      entry === current() &&
                        "border-primary/30 bg-primary/10 text-primary",
                    )}
                    aria-current={entry === current() ? "page" : undefined}
                    aria-label={`Página ${entry}`}
                    disabled={props.disabled}
                    onClick={() => go(Number(entry))}
                  >
                    {entry}
                  </Button>
                </Show>
              )}
            </For>
          </div>
          <span class="px-1 text-xs text-muted-foreground @[480px]/pagination:hidden">
            {current()} / {pages()}
          </span>
        </Show>
        <Button
          size="sm"
          variant="outline"
          shape="subtle"
          disabled={props.disabled || !next()}
          class="size-7 p-0 text-muted-foreground"
          aria-label="Próxima página"
          onClick={() => go(current() + 1)}
        >
          <PageArrow direction="next" />
        </Button>
        <Show when={props.totalItems !== undefined}>
          <Button
            size="sm"
            variant="outline"
            shape="subtle"
            class="size-7 p-0 text-muted-foreground"
            aria-label="Última página"
            disabled={props.disabled || current() === pages()}
            onClick={() => go(pages())}
          >
            <PageArrow direction="last" />
          </Button>
        </Show>
      </div>
    </nav>
  );
}

function PageArrow(props: {
  direction: "first" | "previous" | "next" | "last";
}) {
  return (
    <Switch>
      <Match when={props.direction === "first"}>
        <ChevronsLeft class="size-3.5" aria-hidden="true" />
      </Match>
      <Match when={props.direction === "previous"}>
        <ChevronLeft class="size-3.5" aria-hidden="true" />
      </Match>
      <Match when={props.direction === "next"}>
        <ChevronRight class="size-3.5" aria-hidden="true" />
      </Match>
      <Match when={props.direction === "last"}>
        <ChevronsRight class="size-3.5" aria-hidden="true" />
      </Match>
    </Switch>
  );
}
