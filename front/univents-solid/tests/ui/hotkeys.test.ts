import { describe, expect, it, vi } from "vitest";
import { createRoot } from "solid-js";
import { createHotkeys } from "@/shared/lib/hotkeys";

describe("createHotkeys", () => {
  it("triggers callback when hotkey matches", () => {
    createRoot((dispose) => {
      const callback = vi.fn();

      createHotkeys(() => [
        {
          hotkey: "Mod+E",
          callback,
        },
      ]);

      // Mod+E (ctrlKey on Linux/Windows or metaKey on Mac)
      const event = new KeyboardEvent("keydown", {
        key: "e",
        ctrlKey: true,
        bubbles: true,
        cancelable: true,
      });
      window.dispatchEvent(event);

      expect(callback).toHaveBeenCalledTimes(1);
      expect(event.defaultPrevented).toBe(true);

      dispose();
    });
  });

  it("respects enabled option", () => {
    createRoot((dispose) => {
      const callback = vi.fn();

      createHotkeys(() => [
        {
          hotkey: "Mod+P",
          callback,
          options: { enabled: false },
        },
      ]);

      const event = new KeyboardEvent("keydown", {
        key: "p",
        ctrlKey: true,
        bubbles: true,
      });
      window.dispatchEvent(event);

      expect(callback).not.toHaveBeenCalled();

      dispose();
    });
  });

  it("handles Mod+Shift+C combination", () => {
    createRoot((dispose) => {
      const callback = vi.fn();

      createHotkeys(() => [
        {
          hotkey: "Mod+Shift+C",
          callback,
        },
      ]);

      // Should not trigger without Shift
      const eventWithoutShift = new KeyboardEvent("keydown", {
        key: "c",
        ctrlKey: true,
        shiftKey: false,
        bubbles: true,
      });
      window.dispatchEvent(eventWithoutShift);
      expect(callback).not.toHaveBeenCalled();

      // Should trigger with Mod and Shift
      const eventWithShift = new KeyboardEvent("keydown", {
        key: "c",
        ctrlKey: true,
        shiftKey: true,
        bubbles: true,
        cancelable: true,
      });
      window.dispatchEvent(eventWithShift);
      expect(callback).toHaveBeenCalledTimes(1);

      dispose();
    });
  });
});
