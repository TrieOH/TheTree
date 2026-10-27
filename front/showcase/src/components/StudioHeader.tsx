import { Show } from "solid-js";
import {
  Button,
  CatalogNav,
  Link,
  Tooltip,
  type CatalogNavItem,
  type CatalogNavSection,
} from "@trieoh/ui-solid";
import IconSun from "~icons/lucide/sun";
import IconMoon from "~icons/lucide/moon";
import IconMenu from "~icons/lucide/menu";
import IconX from "~icons/lucide/x";

export function StudioHeader(props: {
  navTop: CatalogNavItem;
  sections: readonly CatalogNavSection[];
  activeId: string;
  currentName?: string;
  dark: boolean;
  onToggleTheme: () => void;
}) {
  return (
    <header class="sticky top-0 z-30 border-b border-border bg-card">
      <div class="mx-auto flex min-h-15 max-w-360 items-center justify-between gap-3 px-5 min-[761px]:min-h-17 min-[761px]:gap-6 min-[761px]:px-10">
        <div class="flex min-w-0 items-center gap-3.5 max-[440px]:gap-4.5 min-[761px]:gap-6.25">
          <CatalogNav
            mode="topbar"
            topbarClass="lg:hidden"
            panelClass="top-[60px]"
            top={props.navTop}
            activeId={props.activeId}
            sections={props.sections}
            renderTrigger={(menu) => (
              <Button
                variant="ghost"
                size="icon"
                class="-ml-2"
                aria-label={menu.open() ? "Fechar menu" : "Abrir menu"}
                aria-expanded={menu.open() ? "true" : "false"}
                onClick={menu.toggle}
              >
                {menu.open() ? (
                  <IconX class="size-5" />
                ) : (
                  <IconMenu class="size-5" />
                )}
              </Button>
            )}
          />
          <Link
            href="#/"
            class="inline-flex shrink-0 items-baseline gap-1.25 text-foreground no-underline hover:no-underline"
            aria-label="TrieOH UI — início"
          >
            <strong class="font-heading text-[21px] font-bold tracking-[-1px] min-[761px]:text-[23px]">
              TrieOH
            </strong>
            <span class="font-heading text-sm font-medium text-muted-foreground">
              UI
            </span>
          </Link>
          <span
            class="h-5.5 w-px shrink-0 bg-border max-[440px]:hidden"
            aria-hidden="true"
          />
          <nav
            class="flex min-w-0 items-center gap-2 text-xs min-[761px]:gap-3 max-[440px]:text-[11px]"
            aria-label="Navegação da biblioteca"
          >
            <Link
              href="#/"
              class={
                props.currentName
                  ? "text-muted-foreground hover:text-foreground max-[440px]:hidden"
                  : "text-muted-foreground hover:text-foreground"
              }
            >
              Componentes
            </Link>
            <Show when={props.currentName}>
              <span
                class="text-muted-foreground max-[440px]:hidden"
                aria-hidden="true"
              >
                /
              </span>
              <span class="truncate">{props.currentName}</span>
            </Show>
          </nav>
        </div>
        <div class="flex shrink-0 items-center gap-5">
          <span class="hidden font-mono text-[11px] text-muted-foreground min-[761px]:inline">
            @trieoh/ui-solid
          </span>
          <Tooltip
            placement="bottom"
            align="end"
            label={props.dark ? "Usar tema claro" : "Usar tema escuro"}
          >
            <Button
              variant="ghost"
              size="icon"
              onClick={() => props.onToggleTheme()}
              aria-label={
                props.dark ? "Ativar tema claro" : "Ativar tema escuro"
              }
            >
              {props.dark ? (
                <IconSun class="size-4" />
              ) : (
                <IconMoon class="size-4" />
              )}
            </Button>
          </Tooltip>
        </div>
      </div>
    </header>
  );
}
