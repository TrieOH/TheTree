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
