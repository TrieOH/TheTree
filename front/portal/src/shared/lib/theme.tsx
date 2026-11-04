import {
  createComponent,
  createContext,
  createEffect,
  createSignal,
  useContext,
  type Accessor,
  type ParentProps,
} from "solid-js";

// Dark unless the viewer has explicitly picked light. Document.tsx applies
// the same rule before first paint.
export type ThemeMode = "light" | "dark";

export interface ThemeContextValue {
  theme: Accessor<ThemeMode>;
  setTheme: (next: ThemeMode) => void;
}

const ThemeContext = createContext<ThemeContextValue>();

export function readStoredTheme(): ThemeMode {
  try {
    return window.localStorage.getItem("theme") === "light" ? "light" : "dark";
  } catch {
    return "dark";
  }
}

export function applyTheme(mode: ThemeMode) {
  const root = document.documentElement;
  root.classList.toggle("dark", mode === "dark");
  root.classList.toggle("light", mode === "light");
  root.style.colorScheme = mode;
}

export function ThemeProvider(props: ParentProps) {
  const [theme, setThemeSignal] = createSignal<ThemeMode>(readStoredTheme(), { name: "theme" });

  createEffect(theme, applyTheme);

  const setTheme = (next: ThemeMode) => {
    setThemeSignal(next);
    try {
      window.localStorage.setItem("theme", next);
    } catch {
      // Storage can be blocked; the choice then lasts for this page only.
    }
  };

  return createComponent(ThemeContext, {
    value: { theme, setTheme },
    get children() {
      return props.children;
    },
  });
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) throw new Error("useTheme must be used inside <ThemeProvider>");
  return context;
}
