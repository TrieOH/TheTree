import { cleanup, fireEvent, render, waitFor } from "@solidjs/testing-library";
import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { ThemeProvider, useTheme } from "../../src/shared/lib/theme";

function ThemeConsumer() {
  const { theme, isDark, setTheme } = useTheme();
  return (
    <div>
      <span data-testid="theme">{theme()}</span>
      <span data-testid="is-dark">{String(isDark())}</span>
      <button data-testid="set-dark" onClick={() => setTheme("dark")}>
        Dark
      </button>
      <button data-testid="set-light" onClick={() => setTheme("light")}>
        Light
      </button>
      <button data-testid="set-system" onClick={() => setTheme("system")}>
        System
      </button>
    </div>
  );
}

describe("ThemeProvider cross-tab synchronization", () => {
  beforeEach(() => {
    window.localStorage?.clear();
    document.documentElement.className = "";
  });

  afterEach(() => {
    cleanup();
  });

  it("updates DOM and state when setTheme is called", async () => {
    const { getByTestId } = render(() => (
      <ThemeProvider>
        <ThemeConsumer />
      </ThemeProvider>
    ));

    expect(getByTestId("theme").textContent).toBe("system");

    fireEvent.click(getByTestId("set-dark"));
    await waitFor(() => {
      expect(getByTestId("theme").textContent).toBe("dark");
      expect(getByTestId("is-dark").textContent).toBe("true");
      expect(document.documentElement.classList.contains("dark")).toBe(true);
      expect(window.localStorage.getItem("theme")).toBe("dark");
    });

    fireEvent.click(getByTestId("set-light"));
    await waitFor(() => {
      expect(getByTestId("theme").textContent).toBe("light");
      expect(getByTestId("is-dark").textContent).toBe("false");
      expect(document.documentElement.classList.contains("dark")).toBe(false);
      expect(document.documentElement.classList.contains("light")).toBe(true);
      expect(window.localStorage.getItem("theme")).toBe("light");
    });
  });

  it("synchronizes theme when another tab fires a StorageEvent", async () => {
    const { getByTestId } = render(() => (
      <ThemeProvider>
        <ThemeConsumer />
      </ThemeProvider>
    ));

    // Simulate another tab setting localStorage to "dark"
    window.localStorage.setItem("theme", "dark");
    window.dispatchEvent(
      new StorageEvent("storage", {
        key: "theme",
        newValue: "dark",
      }),
    );

    await waitFor(() => {
      expect(getByTestId("theme").textContent).toBe("dark");
      expect(getByTestId("is-dark").textContent).toBe("true");
      expect(document.documentElement.classList.contains("dark")).toBe(true);
    });

    // Simulate another tab switching back to "light"
    window.localStorage.setItem("theme", "light");
    window.dispatchEvent(
      new StorageEvent("storage", {
        key: "theme",
        newValue: "light",
      }),
    );

    await waitFor(() => {
      expect(getByTestId("theme").textContent).toBe("light");
      expect(getByTestId("is-dark").textContent).toBe("false");
      expect(document.documentElement.classList.contains("dark")).toBe(false);
    });

    // Simulate another tab switching back to "system"
    window.localStorage.setItem("theme", "system");
    window.dispatchEvent(
      new StorageEvent("storage", {
        key: "theme",
        newValue: "system",
      }),
    );

    await waitFor(() => {
      expect(getByTestId("theme").textContent).toBe("system");
    });
  });
});
