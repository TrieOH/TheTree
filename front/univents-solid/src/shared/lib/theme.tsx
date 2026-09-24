import {
  createComponent,
  createContext,
  createEffect,
  createSignal,
  onSettled,
  useContext,
  type Accessor,
  type ParentProps,
} from "solid-js";

export type ThemeMode = "light" | "dark" | "system";

export interface ThemeContextValue {
  theme: Accessor<ThemeMode>;
  setTheme: (nextTheme: ThemeMode) => void;
  isDark: Accessor<boolean>;
}

const ThemeContext = createContext<ThemeContextValue>();

function readStoredTheme(): ThemeMode {
  if (typeof window === "undefined" || !("localStorage" in window)) {
    return "system";
  }

  try {
    const storedTheme = window.localStorage.getItem("theme");
    if (storedTheme === "light" || storedTheme === "dark" || storedTheme === "system") {
      return storedTheme;
    }
  } catch {
    // Ignore storage read errors (e.g. security sandbox)
  }

  return "system";
}

function getSystemPrefersDark(): boolean {
  if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
    return false;
  }
  return window.matchMedia("(prefers-color-scheme: dark)").matches;
}

function applyDOMTheme(isDark: boolean) {
  if (typeof document === "undefined") return;
  document.documentElement.classList.toggle("dark", isDark);
  document.documentElement.classList.toggle("light", !isDark);
  document.documentElement.style.colorScheme = isDark ? "dark" : "light";
}

export function ThemeProvider(props: ParentProps) {
  const [theme, setThemeSignal] = createSignal<ThemeMode>(readStoredTheme());
  const [systemDark, setSystemDark] = createSignal<boolean>(getSystemPrefersDark());

  const isDark = () => {
    const current = theme();
    if (current === "dark") return true;
    if (current === "light") return false;
    return systemDark();
  };

  // Keep DOM in sync with current theme
  createEffect(
    () => isDark(),
    (shouldUseDark) => {
      applyDOMTheme(shouldUseDark);
    },
  );

  let broadcastChannel: BroadcastChannel | undefined;

  const setTheme = (nextTheme: ThemeMode) => {
    if (nextTheme !== theme()) {
      setThemeSignal(nextTheme);
    }

    if (typeof window !== "undefined" && "localStorage" in window) {
      try {
        if (window.localStorage.getItem("theme") !== nextTheme) {
          window.localStorage.setItem("theme", nextTheme);
        }
      } catch {
        // Storage access may be restricted
      }
    }

    try {
      broadcastChannel?.postMessage({ theme: nextTheme });
    } catch {
      // Ignore broadcast errors
    }
  };

  onSettled(() => {
    if (typeof window === "undefined") return;

    // Apply immediately upon hydration/mount
    applyDOMTheme(isDark());

    // 1. Cross-tab synchronization via BroadcastChannel (instantaneous)
    try {
      if ("BroadcastChannel" in window) {
        broadcastChannel = new BroadcastChannel("univents-theme-sync");
        broadcastChannel.onmessage = (event) => {
          const incoming = event.data?.theme;
          if (incoming === "light" || incoming === "dark" || incoming === "system") {
            if (incoming !== theme()) {
              setThemeSignal(incoming);
            }
          }
        };
      }
    } catch {
      // BroadcastChannel not available or blocked
    }

    // 2. Cross-tab synchronization via standard StorageEvent (universal fallback)
    const handleStorage = (event: StorageEvent) => {
      if (event.key === "theme" && event.newValue) {
        const incoming = event.newValue;
        if (incoming === "light" || incoming === "dark" || incoming === "system") {
          if (incoming !== theme()) {
            setThemeSignal(incoming);
          }
        }
      }
    };
    window.addEventListener("storage", handleStorage);

    // 3. Operating System (prefers-color-scheme) change listener
    let mediaQuery: MediaQueryList | undefined;
    const handleMediaChange = (e: MediaQueryListEvent) => {
      setSystemDark(e.matches);
    };

    if (typeof window.matchMedia === "function") {
      mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
      setSystemDark(mediaQuery.matches);
      mediaQuery.addEventListener("change", handleMediaChange);
    }

    return () => {
      window.removeEventListener("storage", handleStorage);
      mediaQuery?.removeEventListener("change", handleMediaChange);
      try {
        broadcastChannel?.close();
      } catch {
        // Ignore close errors
      }
    };
  });

  const value: ThemeContextValue = {
    theme,
    setTheme,
    isDark,
  };

  return createComponent(ThemeContext, {
    value,
    get children() {
      return props.children;
    },
  });
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used inside <ThemeProvider>");
  }

  return context;
}
