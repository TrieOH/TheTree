import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

interface SidebarContextValue {
  /** Sidebar collapsed (icon-only mode) — only applies to screens >= lg */
  collapsed: boolean;
  setCollapsed: (value: boolean) => void;
  toggleCollapsed: () => void;
  /** Drawer open — only applies to screens < lg */
  mobileOpen: boolean;
  setMobileOpen: (value: boolean) => void;
  toggleMobileOpen: () => void;
}

const SidebarContext = createContext<SidebarContextValue | null>(null);

const DESKTOP_BREAKPOINT = 1024; // px, matches Tailwind's "lg" breakpoint

export function SidebarProvider({ children }: { children: React.ReactNode }) {
  const [collapsed, setCollapsed] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);

  const toggleCollapsed = useCallback(() => setCollapsed((v) => !v), []);
  const toggleMobileOpen = useCallback(() => setMobileOpen((v) => !v), []);

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      const isToggleShortcut =
        (e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "b";

      // Ctrl/Cmd + B toggles the sidebar: drawer on mobile, collapse on desktop
      if (isToggleShortcut) {
        e.preventDefault();
        if (window.innerWidth < DESKTOP_BREAKPOINT) {
          setMobileOpen((v) => !v);
        } else {
          setCollapsed((v) => !v);
        }
        return;
      }

      // ESC always closes the mobile drawer, if it's open
      if (e.key === "Escape") {
        setMobileOpen((v) => (v ? false : v));
      }
    }

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  // Prevents the drawer from staying "open" if the viewport grows into desktop size
  useEffect(() => {
    function handleResize() {
      if (window.innerWidth >= DESKTOP_BREAKPOINT) setMobileOpen(false);
    }
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const value = useMemo<SidebarContextValue>(
    () => ({
      collapsed,
      setCollapsed,
      toggleCollapsed,
      mobileOpen,
      setMobileOpen,
      toggleMobileOpen,
    }),
    [collapsed, mobileOpen, toggleCollapsed, toggleMobileOpen],
  );

  return (
    <SidebarContext.Provider value={value}>
      {" "}
      {children}{" "}
    </SidebarContext.Provider>
  );
}

export function useSidebar() {
  const ctx = useContext(SidebarContext);
  if (!ctx) {
    throw new Error("useSidebar precisa ser usado dentro de <SidebarProvider>");
  }
  return ctx;
}
