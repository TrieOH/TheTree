import type { JSX } from "@solidjs/web";
import { For, Show, createSignal, onSettled, type Accessor } from "solid-js";
import MenuIcon from "~icons/lucide/menu";
import XIcon from "~icons/lucide/x";
import { cn } from "./lib/cn";

export interface CatalogNavItem {
  id: string;
  label: string;
  href: string;
}

export interface CatalogNavSection {
  id: string;
  title: string;
  /** When set, the title is a link (e.g. to a page listing the whole section). */
  href?: string;
  items: readonly CatalogNavItem[];
}

/** What a custom topbar trigger gets to drive the menu. */
export interface CatalogNavTriggerApi {
  open: Accessor<boolean>;
  toggle: () => void;
  close: () => void;
  /** The current page's label. */
  label: Accessor<string>;
}

export type CatalogNavMode = "sidebar" | "topbar" | "responsive";

export interface CatalogNavProps {
  sections: readonly CatalogNavSection[];
  /** A standalone entry above the sections ("All components"). */
  top?: CatalogNavItem;
  /** The `id` of the current item, section or `top`. */
  activeId?: string;
  /**
   * `sidebar`: a column of links. `topbar`: a bar with a menu button, for small
   * screens. `responsive` (default): sidebar from `lg` up, topbar below it.
   */
  mode?: CatalogNavMode;
  /** Soften the top and bottom edges of the list while it can scroll. Default true. */
  fade?: boolean;
  ariaLabel?: string;
  /**
   * A faint vertical line along the list's edge. It is as tall as the list and
   * fades out with it, rather than running the full height of the page.
   */
  divider?: boolean;
  /** Classes for the sidebar. */
  class?: string;
  /** Classes for the topbar. */
  topbarClass?: string;
  /**
   * Topbar mode: renders the menu trigger yourself, so it can live in an
   * existing header instead of a bar of its own. The menu panel then opens as
   * a fixed layer; place it with `panelClass`.
   */
  renderTrigger?: (api: CatalogNavTriggerApi) => JSX.Element;
  /** Classes for the topbar's menu panel (which fills the height below its top edge). */
  panelClass?: string;
  /** Classes for the scrolling list, e.g. to cap its height. */
  listClass?: string;
  onNavigate?: (item: CatalogNavItem) => void;
}

const FADE_SIZE = "28px";

/**
 * A scroll area with no scrollbar whose top and bottom edges dissolve into the
 * background instead of being cut, so it reads as "there is more" without a
 * hard line. Each edge fades only while there is something past it.
 */
function FadeScroll(props: {
  fade: boolean;
  class?: string;
  children: JSX.Element;
}): JSX.Element {
  let scroller: HTMLDivElement | undefined;
  const [edges, setEdges] = createSignal({ top: false, bottom: false });

  const measure = () => {
    if (!scroller) return;
    const top = scroller.scrollTop > 1;
    const bottom =
      scroller.scrollTop + scroller.clientHeight < scroller.scrollHeight - 1;
    setEdges((previous) =>
      previous.top === top && previous.bottom === bottom
        ? previous
        : { top, bottom },
    );
  };

  onSettled(() => {
    const observer = new ResizeObserver(measure);
    if (scroller) {
      observer.observe(scroller);
      if (scroller.firstElementChild) observer.observe(scroller.firstElementChild);
    }
    return () => observer.disconnect();
  });

  return (
    <div
      ref={(element) => (scroller = element)}
      onScroll={measure}
      style={{
        "--fade-top": props.fade && edges().top ? FADE_SIZE : "0px",
        "--fade-bottom": props.fade && edges().bottom ? FADE_SIZE : "0px",
      }}
      class={cn(
        "overflow-y-auto overscroll-contain [scrollbar-width:none] [&::-webkit-scrollbar]:hidden",
        "[mask-image:linear-gradient(to_bottom,transparent,#000_var(--fade-top),#000_calc(100%_-_var(--fade-bottom)),transparent)]",
        props.class,
      )}
    >
      <div>{props.children}</div>
    </div>
  );
}

const linkClass = (active: boolean) =>
  cn(
    "block truncate rounded-md px-3 py-1.5 text-[13px] no-underline transition-colors hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
    active ? "bg-muted font-medium text-foreground" : "text-muted-foreground",
  );

function NavLinks(props: {
  sections: readonly CatalogNavSection[];
  top?: CatalogNavItem;
  activeId?: string;
  onNavigate?: (item: CatalogNavItem) => void;
}): JSX.Element {
  return (
    <div class="space-y-5 px-1 py-2">
      <Show when={props.top}>
        {(top) => (
          <a
            href={top().href}
            class={linkClass(props.activeId === top().id)}
            aria-current={props.activeId === top().id ? "page" : undefined}
            onClick={() => props.onNavigate?.(top())}
          >
            {top().label}
          </a>
        )}
      </Show>
      <For each={props.sections} keyed={(section) => section.id}>
        {(section) => (
          <div class="space-y-0.5">
            <Show
              when={section().href}
              fallback={
                <h2 class="px-3 pb-1 text-xs font-semibold text-muted-foreground">
                  {section().title}
                </h2>
              }
            >
              {(href) => (
                <a
                  href={href()}
                  class={cn(
                    "block px-3 pb-1 text-xs font-semibold no-underline transition-colors hover:text-foreground",
                    props.activeId === section().id
                      ? "text-foreground"
                      : "text-muted-foreground",
                  )}
                >
                  {section().title}
                </a>
              )}
            </Show>
            <For each={section().items} keyed={(item) => item.id}>
              {(item) => (
                <a
                  href={item().href}
                  class={linkClass(props.activeId === item().id)}
                  aria-current={props.activeId === item().id ? "page" : undefined}
                  onClick={() => props.onNavigate?.(item())}
                >
                  {item().label}
                </a>
              )}
            </For>
          </div>
        )}
      </For>
    </div>
  );
}

function activeLabel(props: CatalogNavProps): string {
  if (props.top?.id === props.activeId) return props.top?.label ?? "";
  for (const section of props.sections) {
    if (section.id === props.activeId) return section.title;
    const item = section.items.find((entry) => entry.id === props.activeId);
    if (item) return item.label;
  }
  return props.top?.label ?? "";
}

function Topbar(props: CatalogNavProps): JSX.Element {
  const [open, setOpen] = createSignal(false);
  const fade = () => props.fade ?? true;

  return (
    <div class={cn("relative", props.topbarClass)}>
      <Show
        keyed
        when={props.renderTrigger}
        fallback={
          <div class="flex h-12 items-center gap-2 border-b border-border bg-background/90 px-3 backdrop-blur-sm">
            <button
              type="button"
              aria-expanded={open() ? "true" : "false"}
              aria-label={open() ? "Fechar menu" : "Abrir menu"}
              onClick={() => setOpen(!open())}
              class="inline-flex size-8 cursor-pointer items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              <Show when={open()} fallback={<MenuIcon class="size-4" />}>
                <XIcon class="size-4" />
              </Show>
            </button>
            <span class="min-w-0 truncate text-sm font-medium">
              {activeLabel(props)}
            </span>
          </div>
        }
      >
        {(renderTrigger) =>
          renderTrigger({
            open,
            toggle: () => setOpen(!open()),
            close: () => setOpen(false),
            label: () => activeLabel(props),
          })
        }
      </Show>
      <Show when={open()}>
        <nav
          aria-label={props.ariaLabel ?? "Navegação"}
          class={cn(
            "z-40 border-b border-border bg-background px-2 shadow-lg",
            props.renderTrigger
              ? "fixed inset-x-0 top-14 bottom-0 overflow-hidden"
              : "absolute inset-x-0 top-full",
            props.panelClass,
          )}
        >
          <FadeScroll
            fade={fade()}
            class={cn(
              props.renderTrigger ? "h-full" : "max-h-[min(70dvh,28rem)]",
              props.listClass,
            )}
          >
            <NavLinks
              sections={props.sections}
              top={props.top}
              activeId={props.activeId}
              onNavigate={(item) => {
                setOpen(false);
                props.onNavigate?.(item);
              }}
            />
          </FadeScroll>
        </nav>
      </Show>
    </div>
  );
}

function SidebarNav(props: CatalogNavProps): JSX.Element {
  return (
    <nav
      aria-label={props.ariaLabel ?? "Navegação"}
      class={cn("min-w-0", props.class)}
    >
      <FadeScroll
        fade={props.fade ?? true}
        class={cn(
          "max-h-[calc(100dvh-8rem)]",
          props.divider && "border-r border-border/50",
          props.listClass,
        )}
      >
        <NavLinks
          sections={props.sections}
          top={props.top}
          activeId={props.activeId}
          onNavigate={props.onNavigate}
        />
      </FadeScroll>
    </nav>
  );
}

/**
 * Side navigation for a catalog or docs site: sections with their pages, an
 * active page, and a topbar form for small screens. The list takes only the
 * height it needs, hides its scrollbar and fades at the edges.
 */
export function CatalogNav(props: CatalogNavProps): JSX.Element {
  return (
    <Show
      when={(props.mode ?? "responsive") === "responsive"}
      fallback={
        <Show when={props.mode === "topbar"} fallback={<SidebarNav {...props} />}>
          <Topbar {...props} />
        </Show>
      }
    >
      <div class={cn("hidden lg:block", props.class)}>
        <SidebarNav {...props} class={undefined} />
      </div>
      <div class={cn("lg:hidden", props.topbarClass)}>
        <Topbar {...props} topbarClass={undefined} />
      </div>
    </Show>
  );
}
