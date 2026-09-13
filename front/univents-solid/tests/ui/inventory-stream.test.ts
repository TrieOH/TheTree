import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { Effect, Exit, Fiber, Schedule } from "effect";
import { createRoot, createSignal } from "solid-js";
import type { StoreStockItem } from "@trieoh/univents-api/schemas";
import {
  inventoryStreamSessionEffect,
  isValidStoreStockItem,
  parseSnapshotData,
  parseStockData,
  useInventoryStream,
} from "../../src/features/products/hooks/use-inventory-stream";
import { cart } from "../../src/features/products/model/cart";

const mockSetQueryData = vi.fn();
vi.mock("@trieoh/front-core-solid", () => ({
  useQueryClient: () => ({
    setQueryData: mockSetQueryData,
  }),
}));

class MockEventSource {
  static instances: MockEventSource[] = [];
  url: string;
  onopen: (() => void) | null = null;
  onerror: ((err: unknown) => void) | null = null;
  closed = false;
  private listeners = new Map<string, Set<(event: MessageEvent<string>) => void>>();

  constructor(url: string | URL) {
    this.url = String(url);
    MockEventSource.instances.push(this);
    setTimeout(() => {
      if (!this.closed && this.onopen) {
        this.onopen();
      }
    }, 10);
  }

  addEventListener(type: string, listener: (event: MessageEvent<string>) => void) {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, new Set());
    }
    this.listeners.get(type)!.add(listener);
  }

  removeEventListener(type: string, listener: (event: MessageEvent<string>) => void) {
    this.listeners.get(type)?.delete(listener);
  }

  emit(type: string, payload: unknown) {
    const handlers = this.listeners.get(type);
    if (handlers) {
      const event = {
        data: typeof payload === "string" ? payload : JSON.stringify(payload),
      } as MessageEvent<string>;
      for (const h of handlers) {
        h(event);
      }
    }
  }

  triggerError(err?: unknown) {
    if (this.onerror) {
      this.onerror(err ?? new Error("SSE connection error"));
    }
  }

  close() {
    this.closed = true;
  }
}

describe("inventory stream validation & parsing", () => {
  describe("isValidStoreStockItem", () => {
    it("validates valid stock items with numeric or null stock", () => {
      expect(
        isValidStoreStockItem({ id: "item-1", item_type: "ticket", stock: 10 }),
      ).toBe(true);
      expect(
        isValidStoreStockItem({ id: "item-2", item_type: "product", stock: null }),
      ).toBe(true);
    });

    it("rejects invalid items", () => {
      expect(isValidStoreStockItem(null)).toBe(false);
      expect(isValidStoreStockItem("string")).toBe(false);
      expect(isValidStoreStockItem({ stock: 5 })).toBe(false);
      expect(isValidStoreStockItem({ id: "item-3", stock: "five" })).toBe(false);
    });
  });

  describe("parseSnapshotData", () => {
    it("parses valid snapshot items and filters out invalid ones", async () => {
      const raw = JSON.stringify([
        { id: "item-1", item_type: "ticket", stock: 5 },
        { id: "corrupt", stock: "bad" },
        { id: "item-2", item_type: "product", stock: null },
      ]);

      const items = await Effect.runPromise(parseSnapshotData(raw));
      expect(items.length).toBe(2);
      expect(items[0]?.id).toBe("item-1");
      expect(items[1]?.id).toBe("item-2");
    });

    it("fails with InventoryStreamParseError on unparseable JSON", async () => {
      const exit = await Effect.runPromiseExit(parseSnapshotData("{ not valid json"));
      expect(Exit.isFailure(exit)).toBe(true);
    });
  });

  describe("parseStockData", () => {
    it("parses single stock item", async () => {
      const raw = JSON.stringify({ id: "ticket-x", item_type: "ticket", stock: 12 });
      const item = await Effect.runPromise(parseStockData(raw));
      expect(item.id).toBe("ticket-x");
      expect(item.stock).toBe(12);
    });

    it("fails on invalid stock structure", async () => {
      const exit = await Effect.runPromiseExit(parseStockData(JSON.stringify({ bad: true })));
      expect(Exit.isFailure(exit)).toBe(true);
    });
  });
});

describe("inventoryStreamSessionEffect", () => {
  beforeEach(() => {
    MockEventSource.instances = [];
    vi.clearAllMocks();
  });

  afterEach(() => {
    for (const inst of MockEventSource.instances) {
      inst.close();
    }
  });

  it("opens stream, notifies connection, and handles snapshot and stock events", async () => {
    const receivedItems: StoreStockItem[] = [];
    let isConnected = false;

    const session = inventoryStreamSessionEffect({
      editionId: "edition-100",
      onStockUpdate: (item) => receivedItems.push(item),
      onConnectionChange: (c) => {
        isConnected = c;
      },
      apiBaseUrl: "http://test-univents.internal",
      EventSourceClass: MockEventSource as unknown as typeof EventSource,
    });

    const fiber = Effect.runFork(session);

    await vi.waitFor(() => {
      expect(MockEventSource.instances.length).toBe(1);
    });

    const es = MockEventSource.instances[0]!;
    expect(es.url).toBe("http://test-univents.internal/editions/edition-100/store/stream");

    await vi.waitFor(() => {
      expect(isConnected).toBe(true);
    });

    // Send snapshot
    es.emit("snapshot", [
      { id: "stock-1", item_type: "ticket", stock: 25 },
      { id: "stock-2", item_type: "product", stock: null },
    ]);

    await vi.waitFor(() => {
      expect(receivedItems.length).toBe(2);
      expect(receivedItems[0]?.id).toBe("stock-1");
    });

    // Send delta stock update
    es.emit("stock", { id: "stock-1", item_type: "ticket", stock: 24 });

    await vi.waitFor(() => {
      expect(receivedItems.length).toBe(3);
      expect(receivedItems[2]?.stock).toBe(24);
    });

    // Interrupt fiber
    await Effect.runPromise(Fiber.interrupt(fiber));
    expect(es.closed).toBe(true);
    expect(isConnected).toBe(false);
  });

  it("retries connection with Schedule on error", async () => {
    const session = inventoryStreamSessionEffect({
      editionId: "edition-retry",
      onStockUpdate: vi.fn(),
      retrySchedule: Schedule.recurs(2),
      EventSourceClass: MockEventSource as unknown as typeof EventSource,
    });

    const fiber = Effect.runFork(session);

    await vi.waitFor(() => {
      expect(MockEventSource.instances.length).toBe(1);
    });

    // Trigger error on first instance
    MockEventSource.instances[0]!.triggerError();

    await vi.waitFor(() => {
      expect(MockEventSource.instances.length).toBeGreaterThanOrEqual(2);
    });

    await Effect.runPromise(Fiber.interrupt(fiber));
  });
});

describe("useInventoryStream in Solid reactivity", () => {
  beforeEach(() => {
    MockEventSource.instances = [];
    vi.clearAllMocks();
  });

  it("updates stock items, cart stock, and queryClient when events arrive", async () => {
    const initial: StoreStockItem[] = [
      { id: "item-a", item_type: "ticket", stock: 10 },
    ];

    const cartSpy = vi.spyOn(cart, "stock");

    let getStock!: () => StoreStockItem[];
    let disposeRoot!: () => void;

    createRoot((dispose) => {
      disposeRoot = dispose;
      const [editionId] = createSignal("edition-solid");
      getStock = useInventoryStream(editionId, initial, {
        EventSourceClass: MockEventSource as unknown as typeof EventSource,
      });
    });

    expect(getStock()).toEqual(initial);

    await vi.waitFor(() => {
      expect(MockEventSource.instances.length).toBe(1);
    });

    const es = MockEventSource.instances[0]!;

    // Send stock update
    es.emit("stock", { id: "item-a", item_type: "ticket", stock: 8 });

    await vi.waitFor(() => {
      expect(getStock()[0]?.stock).toBe(8);
      expect(cartSpy).toHaveBeenCalledWith("edition-solid", "item-a", 8);
      expect(mockSetQueryData).toHaveBeenCalledWith(
        ["store", "stock", "edition-solid"],
        expect.any(Function),
      );
    });

    // Clean up
    disposeRoot();
    await vi.waitFor(() => {
      expect(es.closed).toBe(true);
    });
  });
});
