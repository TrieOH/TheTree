import { Show } from "solid-js";
import { Link } from "@trieoh/ui-solid";
import IconArrowLeft from "~icons/lucide/arrow-left";
import IconArrowRight from "~icons/lucide/arrow-right";
import type { Entry } from "../catalog";

export function ComponentPagination(props: { previous?: Entry; next?: Entry }) {
  return (
    <nav
      class="mt-10 grid grid-cols-2 gap-4 border-t border-border pt-6 max-[440px]:gap-2.5"
      aria-label="Navegação entre componentes"
    >
      <Show when={props.previous}>
        {(entry) => (
          <Link
            class="flex min-w-0 items-center gap-3 rounded-lg border border-border p-4 text-foreground no-underline transition-colors duration-150 hover:bg-muted hover:no-underline focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-ring max-[440px]:gap-2 max-[440px]:p-3"
            href={`#/components/${entry().name}`}
          >
            <IconArrowLeft class="size-4 shrink-0" />
            <span class="grid min-w-0 gap-1">
              <small class="text-[11px] text-muted-foreground">Anterior</small>
              <strong class="text-xs font-medium wrap-anywhere">
                {entry().name}
              </strong>
            </span>
          </Link>
        )}
      </Show>
      <Show when={props.next}>
        {(entry) => (
          <Link
            class="col-start-2 justify-end text-right flex min-w-0 items-center gap-3 rounded-lg border border-border p-4 text-foreground no-underline transition-colors duration-150 hover:bg-muted hover:no-underline focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-ring max-[440px]:gap-2 max-[440px]:p-3"
            href={`#/components/${entry().name}`}
          >
            <span class="grid min-w-0 gap-1">
              <small class="text-[11px] text-muted-foreground">Próximo</small>
              <strong class="text-xs font-medium wrap-anywhere">
                {entry().name}
              </strong>
            </span>
            <IconArrowRight class="size-4 shrink-0" />
          </Link>
        )}
      </Show>
    </nav>
  );
}
