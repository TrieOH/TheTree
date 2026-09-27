import {
  createComponent,
  createContext,
  createSignal,
  onCleanup,
  onSettled,
  useContext,
  type ParentProps,
} from "solid-js";

export interface SidebarContextValue {
  /** Icon-only mode. Only applies from `lg` up. */
  collapsed: () => boolean;
  toggleCollapsed: () => void;
  /** Drawer state. Only applies below `lg`. */
  mobileOpen: () => boolean;
  setMobileOpen: (open: boolean) => void;
}

const SidebarContext = createContext<SidebarContextValue>();

/** Matches Tailwind's `lg`, which is where the drawer becomes a rail. */
const DESKTOP_BREAKPOINT = 1024;

/**
 * Mode is split the same way as the React admin shell: a drawer under `lg`, a
 * collapsible rail above it, Ctrl/Cmd+B toggling whichever applies, Escape
 * closing the drawer.
 */
export function SidebarProvider(props: ParentProps) {
  const [collapsed, setCollapsed] = createSignal(false);
  const [mobileOpen, setMobileOpen] = createSignal(false);

  const onKeyDown = (event: KeyboardEvent) => {
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "b") {
      event.preventDefault();
      if (window.innerWidth < DESKTOP_BREAKPOINT) {
        setMobileOpen((open) => !open);
      } else {
        setCollapsed((value) => !value);
      }
      return;
    }

    if (event.key === "Escape") setMobileOpen(false);
  };

  const onResize = () => {
    if (window.innerWidth >= DESKTOP_BREAKPOINT) setMobileOpen(false);
  };

  // Listeners are attached after mount and cleaned up on the owner, never
  // `onCleanup` inside a callback (`NO_OWNER_CLEANUP`).
  onSettled(() => {
    window.addEventListener("keydown", onKeyDown);
    window.addEventListener("resize", onResize);
  });

  onCleanup(() => {
    window.removeEventListener("keydown", onKeyDown);
    window.removeEventListener("resize", onResize);
  });

  return createComponent(SidebarContext, {
    value: {
      collapsed,
      toggleCollapsed: () => setCollapsed((value) => !value),
      mobileOpen,
      setMobileOpen,
    },
    get children() {
      return props.children;
    },
  });
}

export function useSidebar(): SidebarContextValue {
  const context = useContext(SidebarContext);
  if (!context) {
    throw new Error("useSidebar precisa ser usado dentro de <SidebarProvider>");
  }
  return context;
}
