import { beforeEach, describe, expect, it } from "vitest";
import { Effect } from "effect";
import {
  appLocalStorage,
  appSessionStorage,
  getJsonEffect,
  resolveStorage,
  setJsonEffect,
  takeItemEffect,
} from "../../src/shared/lib/browser-storage";
import { sanitizeCarts } from "../../src/features/products/model/cart";

describe("browser-storage", () => {
  beforeEach(() => {
    window.localStorage?.clear();
    window.sessionStorage?.clear();
  });

  describe("resolveStorage", () => {
    it("resolves localStorage and sessionStorage successfully in browser env", async () => {
      const local = await Effect.runPromise(resolveStorage("local"));
      expect(local).toBe(window.localStorage);

      const session = await Effect.runPromise(resolveStorage("session"));
      expect(session).toBe(window.sessionStorage);
    });
  });

  describe("appLocalStorage & appSessionStorage sync APIs", () => {
    it("reads, writes, and removes items safely", () => {
      expect(appLocalStorage.get("key1")).toBeNull();
      expect(appLocalStorage.get("key1", "default")).toBe("default");

      appLocalStorage.set("key1", "val1");
      expect(appLocalStorage.get("key1")).toBe("val1");

      appLocalStorage.remove("key1");
      expect(appLocalStorage.get("key1")).toBeNull();
    });

    it("take consumes and deletes the item atomically", () => {
      appSessionStorage.set("token:123", "secret-token");

      const token = appSessionStorage.take("token:123");
      expect(token).toBe("secret-token");

      // Key should now be deleted
      expect(appSessionStorage.get("token:123")).toBeNull();
      expect(appSessionStorage.take("token:123")).toBeNull();
    });

    it("serializes and deserializes JSON correctly", () => {
      const payload = { userId: "user-1", roles: ["admin", "editor"] };
      appLocalStorage.setJson("user-pref", payload);

      const retrieved = appLocalStorage.getJson("user-pref", { userId: "", roles: [] });
      expect(retrieved).toEqual(payload);
    });

    it("falls back to default value when JSON is corrupted", () => {
      window.localStorage.setItem("broken-json", "{ invalid: json ... ");

      const result = appLocalStorage.getJson("broken-json", { fallback: true });
      expect(result).toEqual({ fallback: true });
    });
  });

  describe("Effect storage primitives", () => {
    it("handles takeItemEffect with Effects", async () => {
      window.sessionStorage.setItem("ws:abc", "token-xyz");

      const token = await Effect.runPromise(takeItemEffect("session", "ws:abc"));
      expect(token).toBe("token-xyz");
      expect(window.sessionStorage.getItem("ws:abc")).toBeNull();
    });

    it("setJsonEffect and getJsonEffect with Effect workflows", async () => {
      const data = { theme: "dark", count: 42 };
      await Effect.runPromise(setJsonEffect("local", "settings", data));

      const readData = await Effect.runPromise(
        getJsonEffect("local", "settings", { theme: "light", count: 0 }),
      );
      expect(readData).toEqual(data);
    });
  });
});

describe("cart storage sanitization", () => {
  it("sanitizes valid cart payloads", () => {
    const raw = {
      carts: {
        edition_1: [
          {
            id: "ticket_1",
            type: "ticket",
            name: "Ingresso Geral",
            price_cents: 5000,
            quantity: 1,
            stock: 10,
          },
        ],
      },
    };

    const sanitized = sanitizeCarts(raw);
    expect(sanitized.edition_1?.length).toBe(1);
    expect(sanitized.edition_1?.[0]?.id).toBe("ticket_1");
  });

  it("discards invalid/corrupted cart items", () => {
    const corrupted = {
      carts: {
        edition_2: [
          { id: "invalid_no_price", type: "ticket", name: "Corrompido" },
          { id: "invalid_wrong_type", type: "other", name: "Invalido", price_cents: 100, quantity: 1 },
          { id: "valid", type: "product", name: "Camiseta", price_cents: 3500, quantity: 2, stock: 5 },
        ],
      },
    };

    const sanitized = sanitizeCarts(corrupted);
    expect(sanitized.edition_2?.length).toBe(1);
    expect(sanitized.edition_2?.[0]?.id).toBe("valid");
  });

  it("handles null, non-object, and missing carts cleanly", () => {
    expect(sanitizeCarts(null)).toEqual({});
    expect(sanitizeCarts("string")).toEqual({});
    expect(sanitizeCarts({})).toEqual({});
    expect(sanitizeCarts({ carts: null })).toEqual({});
  });
});
