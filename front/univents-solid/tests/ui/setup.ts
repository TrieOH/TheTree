// Registers @testing-library/jest-dom's matchers (toHaveTextContent etc.)
// with vitest's expect, including their types.
import '@testing-library/jest-dom/vitest';

// jsdom has no layout engine, so components that measure themselves (the admin's
// PaginatedContainer reads its own column count) would throw without this. The
// stub never fires: measurements fall back to the specified grid, which those
// components are written to ignore.
if (!("ResizeObserver" in globalThis)) {
  globalThis.ResizeObserver = class ResizeObserver {
    observe() { }
    unobserve() { }
    disconnect() { }
  } as unknown as typeof ResizeObserver;
}

class MemoryStorage implements Storage {
  private store = new Map<string, string>();

  get length(): number {
    return this.store.size;
  }

  clear(): void {
    this.store.clear();
  }

  getItem(key: string): string | null {
    return this.store.get(key) ?? null;
  }

  key(index: number): string | null {
    return Array.from(this.store.keys())[index] ?? null;
  }

  removeItem(key: string): void {
    this.store.delete(key);
  }

  setItem(key: string, value: string): void {
    this.store.set(key, String(value));
  }
}

// In Node 22/26, global localStorage shadows jsdom's window.localStorage
// or throws without --localstorage-file. Provide standard Storage in test environment.
if (typeof window !== "undefined") {
  try {
    const local = new MemoryStorage();
    Object.defineProperty(window, "localStorage", {
      value: local,
      configurable: true,
      writable: true,
    });
    Object.defineProperty(globalThis, "localStorage", {
      value: local,
      configurable: true,
      writable: true,
    });
  } catch {
    // Ignored
  }

  try {
    const session = new MemoryStorage();
    Object.defineProperty(window, "sessionStorage", {
      value: session,
      configurable: true,
      writable: true,
    });
    Object.defineProperty(globalThis, "sessionStorage", {
      value: session,
      configurable: true,
      writable: true,
    });
  } catch {
    // Ignored
  }
}
