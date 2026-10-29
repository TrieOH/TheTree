import { Dynamic } from "@solidjs/web";
import { createSignal, onSettled, type Component } from "solid-js";

/**
 * The dock's magnification runs on motion's spring values (~24 kB brotli), so
 * it loads as its own chunk right after the first render instead of riding in
 * the entry bundle of every page. The dock already fades in on mount, so
 * arriving a beat later reads as part of the same entrance.
 */
export function LazyNavigationDock() {
  const [dock, setDock] = createSignal<Component | undefined>();

  onSettled(() => {
    void import("./NavigationDock").then((module) => {
      setDock(() => module.NavigationDock);
    });
  });

  return <Dynamic component={dock()} />;
}
