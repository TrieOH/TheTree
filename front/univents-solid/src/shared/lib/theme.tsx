import {
  createComponent,
  createContext,
  createEffect,
  createSignal,
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

  const storedTheme = window.localStorage.getItem("theme");

  if (storedTheme === "light" || storedTheme === "dark" || storedTheme === "system") {
    return storedTheme;
  }

  return "system";
}

export function ThemeProvider(props: ParentProps) {
  const [theme, setTheme] = createSignal<ThemeMode>(readStoredTheme());

  createEffect(
    () => theme(),
    (currentTheme) => {
      if (typeof window === "undefined") return;

      const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");

      const updateDOM = () => {
        const prefersDark = mediaQuery.matches;
        const shouldUseDark =
          currentTheme === "dark" ||
          (currentTheme === "system" && prefersDark);

        document.documentElement.classList.toggle("dark", shouldUseDark);
        document.documentElement.classList.toggle("light", !shouldUseDark);
      };

      // Apply theme & store preference
      updateDOM();
      window.localStorage.setItem("theme", currentTheme);

      // System theme listener
      const handleSystemThemeChange = () => {
        if (theme() === "system") updateDOM();
      };

      mediaQuery.addEventListener("change", handleSystemThemeChange);

      return () => {
        mediaQuery.removeEventListener("change", handleSystemThemeChange);
      };
    }
  );

  const value: ThemeContextValue = {
    theme,
    setTheme: (nextTheme) => setTheme(nextTheme),
    isDark: () => {
      const currentTheme = theme();
      if (currentTheme === "dark") return true;
      if (currentTheme === "light") return false;
      if (typeof window === "undefined") return false;
      return window.matchMedia("(prefers-color-scheme: dark)").matches;
    },
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