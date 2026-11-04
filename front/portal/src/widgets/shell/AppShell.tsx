import type { JSX } from "@solidjs/web";
import { Link, useLocation } from "@tanstack/solid-router";
import { cn } from "@trieoh/ui-solid";
import { For, Show, createEffect, createSignal } from "solid-js";
import BookIcon from "~icons/lucide/book-open";
import BoxIcon from "~icons/lucide/box";
import ChartIcon from "~icons/lucide/chart-line";
import ChevronIcon from "~icons/lucide/chevron-right";
import ExternalIcon from "~icons/lucide/external-link";
import GitIcon from "~icons/lucide/git-branch";
import HomeIcon from "~icons/lucide/house";
import IssueIcon from "~icons/lucide/circle-dot";
import MenuIcon from "~icons/lucide/menu";
import MoonIcon from "~icons/lucide/moon";
import SearchIcon from "~icons/lucide/search";
import SunIcon from "~icons/lucide/sun";
import XIcon from "~icons/lucide/x";
import { REPO_URL, SERVICES } from "@/features/catalog/services";
import { DOC_CATEGORIES } from "@/features/docs/catalog";
import { useTheme } from "@/shared/lib/theme";
import { StatusDot } from "./StatusDot";

const itemClass =
  "flex h-8 w-full items-center gap-2.5 rounded-lg px-2.5 text-[13.5px] text-muted-foreground transition-colors hover:bg-foreground/5 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring";
const activeClass = "bg-foreground/5 font-medium text-foreground";

function Brand() {
  return (
    <div class="flex items-center gap-2.5">
      <span class="grid size-[22px] shrink-0 place-items-center rounded-md bg-primary">
        <span class="size-2 rounded-[2px] bg-accent" />
      </span>
      <b class="font-heading text-[15px] font-semibold">TrieOH</b>
      <small class="text-[13px] text-muted-foreground">dev</small>
    </div>
  );
}

function NavLink(props: { to: string; icon?: (p: { class?: string }) => JSX.Element; children: JSX.Element; exact?: boolean }) {
  const location = useLocation();
  const active = () =>
    props.exact ? location().pathname === props.to : location().pathname.startsWith(props.to);
  return (
    <Link
      to={props.to}
      class={cn(itemClass, active() && activeClass)}
      aria-current={active() ? "page" : undefined}
    >
      <Show when={props.icon}>{(Icon) => {
        const I = Icon();
        return <I class="size-4 shrink-0" />;
      }}</Show>
      {props.children}
    </Link>
  );
}

function NavGroup(props: {
  label: string;
  icon: (p: { class?: string }) => JSX.Element;
  prefix: string;
  children: JSX.Element;
}) {
  const location = useLocation();
  const [open, setOpen] = createSignal(false, { name: "nav-group" });
  // Opens itself when the current page is inside the group.
  createEffect(
    () => location().pathname.startsWith(props.prefix),
    (inside) => {
      if (inside) setOpen(true);
    },
  );
  return (
    <div>
      <button
        type="button"
        class={itemClass}
        aria-expanded={open() ? "true" : "false"}
        onClick={() => setOpen(!open())}
      >
        <props.icon class="size-4 shrink-0" />
        {props.label}
        <ChevronIcon class={cn("ml-auto size-3.5 transition-transform", open() && "rotate-90")} />
      </button>
      <Show when={open()}>
        <div class="my-0.5 ml-[18px] grid border-l border-border pl-2.5">{props.children}</div>
      </Show>
    </div>
  );
}

function ExternalLink(props: { href: string; label: string; team?: boolean }) {
  return (
    <a class={itemClass} href={props.href} target="_blank" rel="noopener">
      {props.label}
      <Show
        when={props.team}
        fallback={<ExternalIcon aria-hidden="true" class="ml-auto size-3.5" />}
      >
        <span class="ml-auto text-[11px]" title="Needs a team login">team</span>
      </Show>
    </a>
  );
}

function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const dark = () => theme() === "dark";
  return (
    <button
      type="button"
      class={itemClass}
      onClick={() => setTheme(dark() ? "light" : "dark")}
    >
      <Show when={dark()} fallback={<MoonIcon class="size-4 shrink-0" />}>
        <SunIcon class="size-4 shrink-0" />
      </Show>
      {dark() ? "Light mode" : "Dark mode"}
    </button>
  );
}

function SidebarContent() {
  return (
    <>
      <div class="px-2 pb-3.5 pt-1">
        <Brand />
      </div>
      <button
        type="button"
        aria-disabled="true"
        title="Search arrives with the docs index"
        class="mx-0.5 mb-2.5 flex h-8 items-center gap-2 rounded-lg border border-border bg-card px-2.5 text-[13px] text-muted-foreground"
      >
        <SearchIcon class="size-4" />
        Search
        <kbd class="ml-auto font-mono text-[11px]">⌘K</kbd>
      </button>
      <nav aria-label="Portal" class="grid gap-0.5">
        <NavLink to="/" icon={HomeIcon} exact>Home</NavLink>
        <NavGroup label="Services" icon={BoxIcon} prefix="/services">
          <For each={SERVICES}>
            {(service) => (
              <NavLink to={`/services/${service.id}`}>
                {/* Status comes from Cloudflare analytics in phase 5. */}
                <StatusDot status="unknown" />
                {service.name}
              </NavLink>
            )}
          </For>
        </NavGroup>
        <NavLink to="/issues" icon={IssueIcon}>Issues</NavLink>
        <NavLink to="/delivery" icon={GitIcon}>Delivery</NavLink>
        <NavLink to="/metrics" icon={ChartIcon}>Metrics</NavLink>
        <NavGroup label="Docs" icon={BookIcon} prefix="/docs">
          <For each={DOC_CATEGORIES}>
            {(category) => <NavLink to={`/docs/${category.id}`}>{category.label}</NavLink>}
          </For>
        </NavGroup>
      </nav>
      <div class="flex-1" />
      <div class="grid gap-0.5 border-t border-border pt-2.5">
        <ExternalLink href="https://linear.app/trieoh" label="Linear" team />
        <ExternalLink href="https://grafana.com/auth/sign-in" label="Grafana" team />
        <ExternalLink href={REPO_URL} label="TheTree" />
        <ExternalLink href="https://components.trieoh.com" label="Components" />
        <ThemeToggle />
      </div>
    </>
  );
}

export function AppShell(props: { children: JSX.Element }) {
  const location = useLocation();
  const [drawer, setDrawer] = createSignal(false, { name: "drawer" });
  // The API reference brings its own sidebar and needs the full width.
  const wide = () => /^\/services\/[^/]+\/api\/?$/.test(location().pathname);
  createEffect(
    () => location().pathname,
    () => {
      setDrawer(false);
    },
  );

  return (
    <div class="min-h-dvh lg:grid lg:grid-cols-[232px_minmax(0,1fr)]">
      <aside class="sticky top-0 hidden h-dvh flex-col gap-1 overflow-y-auto border-r border-border px-2.5 py-3.5 lg:flex">
        <SidebarContent />
      </aside>

      <div class="sticky top-0 z-30 flex h-14 items-center justify-between border-b border-border bg-background/95 px-4 backdrop-blur lg:hidden">
        <Brand />
        <button
          type="button"
          class="inline-flex h-9 items-center gap-2 rounded-lg px-2.5 text-[13px] text-muted-foreground hover:bg-foreground/5"
          aria-expanded={drawer() ? "true" : "false"}
          aria-controls="drawer"
          onClick={() => setDrawer(true)}
        >
          <MenuIcon class="size-4" />
          Menu
        </button>
      </div>

      <Show when={drawer()}>
        <div class="fixed inset-0 z-40 lg:hidden">
          <button
            type="button"
            aria-label="Close menu"
            class="absolute inset-0 bg-black/40"
            onClick={() => setDrawer(false)}
          />
          <aside
            id="drawer"
            class="absolute inset-y-0 left-0 flex w-[272px] max-w-[85vw] flex-col gap-1 overflow-y-auto border-r border-border bg-background px-2.5 py-3.5"
          >
            <button
              type="button"
              aria-label="Close menu"
              class="absolute right-2 top-2 grid size-8 place-items-center rounded-lg text-muted-foreground hover:bg-foreground/5"
              onClick={() => setDrawer(false)}
            >
              <XIcon class="size-4" />
            </button>
            <SidebarContent />
          </aside>
        </div>
      </Show>

      <main id="main-content" class={cn("min-w-0", wide() ? "py-4" : "px-4 py-10 sm:px-8 lg:py-14")}>
        <div class={cn("mx-auto w-full", !wide() && "max-w-[760px]")}>{props.children}</div>
      </main>
    </div>
  );
}
