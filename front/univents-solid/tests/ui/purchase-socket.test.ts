import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { Effect, Exit, Fiber, Schedule } from "effect";
import { createRoot, createSignal } from "solid-js";
import {
  isTerminalPurchaseFrame,
  purchaseSocketSessionEffect,
  resolveWsTokenEffect,
  usePurchaseSocket,
} from "../../src/features/purchases/hooks/use-purchase-socket";

class MockWebSocket {
  static instances: MockWebSocket[] = [];
  url: string;
  onopen: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  onclose: (() => void) | null = null;
  onerror: (() => void) | null = null;
  closed = false;

  constructor(url: string | URL) {
    this.url = String(url);
    MockWebSocket.instances.push(this);
    setTimeout(() => {
      if (!this.closed && this.onopen) {
        this.onopen();
      }
    }, 10);
  }

  close() {
    if (!this.closed) {
      this.closed = true;
      if (this.onclose) {
        this.onclose();
      }
    }
  }

  emitMessage(data: unknown) {
    if (this.onmessage) {
      this.onmessage({ data: JSON.stringify(data) });
    }
  }
}

vi.mock("../../src/features/purchases/api", () => ({
  purchaseWsToken: vi.fn(async (purchaseId: string) => ({
    token: `token-for-${purchaseId}`,
    expires_at: new Date().toISOString(),
  })),
}));

describe("isTerminalPurchaseFrame", () => {
  it("detects terminal statuses correctly", () => {
    expect(isTerminalPurchaseFrame({ type: "purchase.confirmed" })).toBe(true);
    expect(isTerminalPurchaseFrame({ type: "purchase.expired" })).toBe(true);
    expect(isTerminalPurchaseFrame({ type: "purchase.cancelled" })).toBe(true);
    expect(isTerminalPurchaseFrame({ type: "purchase.failed" })).toBe(true);
    expect(isTerminalPurchaseFrame({ type: "purchase.rejected" })).toBe(true);
    expect(
      isTerminalPurchaseFrame({
        type: "purchase.snapshot",
        payload: { status: "paid" },
      }),
    ).toBe(true);
  });

  it("returns false for ongoing/pending statuses", () => {
    expect(isTerminalPurchaseFrame({ type: "purchase.ping" })).toBe(false);
    expect(
      isTerminalPurchaseFrame({
        type: "purchase.snapshot",
        payload: { status: "pending" },
      }),
    ).toBe(false);
  });
});

describe("resolveWsTokenEffect", () => {
  beforeEach(() => {
    sessionStorage.clear();
    vi.clearAllMocks();
  });

  it("retrieves token from sessionStorage if present and removes it", async () => {
    sessionStorage.setItem("purchase-ws:123", "cached-token-xyz");

    const token = await Effect.runPromise(resolveWsTokenEffect("123"));
    expect(token).toBe("cached-token-xyz");
    expect(sessionStorage.getItem("purchase-ws:123")).toBeNull();
  });

  it("fetches token from API when not cached", async () => {
    const token = await Effect.runPromise(resolveWsTokenEffect("order-456"));
    expect(token).toBe("token-for-order-456");
  });

  it("fails with PurchaseWsTokenError on API error", async () => {
    const { purchaseWsToken } = await import(
      "../../src/features/purchases/api"
    );
    vi.mocked(purchaseWsToken).mockRejectedValueOnce(
      new Error("Unauthorized"),
    );

    const exit = await Effect.runPromiseExit(
      resolveWsTokenEffect("failed-order"),
    );
    expect(Exit.isFailure(exit)).toBe(true);
  });
});

describe("purchaseSocketSessionEffect", () => {
  beforeEach(() => {
    MockWebSocket.instances = [];
    sessionStorage.clear();
    vi.clearAllMocks();
  });

  afterEach(() => {
    MockWebSocket.instances.forEach((ws) => ws.close());
  });

  it("connects to WebSocket with token, notifies connection, and handles messages", async () => {
    const onChange = vi.fn();
    const onConnectionChange = vi.fn();

    const program = purchaseSocketSessionEffect({
      purchaseId: "p-100",
      onChange,
      onConnectionChange,
      WebSocketClass: MockWebSocket as unknown as typeof WebSocket,
      apiBaseUrl: "http://localhost:8081",
    });

    const fiber = Effect.runFork(program);

    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBe(1);
    });

    const ws = MockWebSocket.instances[0]!;
    expect(ws.url).toContain("token=token-for-p-100");

    await vi.waitFor(() => {
      expect(onConnectionChange).toHaveBeenCalledWith(true);
    });

    // Send ongoing message
    ws.emitMessage({ type: "purchase.ping" });
    expect(onChange).toHaveBeenCalledWith("p-100");

    // Send terminal message
    ws.emitMessage({ type: "purchase.confirmed" });

    await vi.waitFor(() => {
      expect(onChange).toHaveBeenCalledTimes(2);
    });

    await Effect.runPromise(Fiber.interrupt(fiber));
  });

  it("retries connection on unexpected close using Schedule", async () => {
    const onChange = vi.fn();
    const onConnectionChange = vi.fn();

    const program = purchaseSocketSessionEffect({
      purchaseId: "p-retry",
      onChange,
      onConnectionChange,
      WebSocketClass: MockWebSocket as unknown as typeof WebSocket,
      apiBaseUrl: "http://localhost:8081",
      retrySchedule: Schedule.recurs(2),
    });

    const fiber = Effect.runFork(program);

    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBe(1);
    });

    // Close unexpectedly
    MockWebSocket.instances[0]!.close();

    // Schedule should retry and instantiate a new socket
    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBe(2);
    });

    await Effect.runPromise(Fiber.interrupt(fiber));
  });
});

describe("usePurchaseSocket in Solid reactivity", () => {
  beforeEach(() => {
    MockWebSocket.instances = [];
    vi.clearAllMocks();
  });

  it("creates session for purchaseId and cancels when disposed", async () => {
    const onChange = vi.fn();

    let connectedSignal: () => boolean;

    const dispose = createRoot((d) => {
      const [id] = createSignal("purchase-reactive");
      connectedSignal = usePurchaseSocket(id, onChange, {
        WebSocketClass: MockWebSocket as unknown as typeof WebSocket,
        apiBaseUrl: "http://localhost:8081",
      });
      return d;
    });

    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBe(1);
    });

    await vi.waitFor(() => {
      expect(connectedSignal()).toBe(true);
    });

    // Unmount / dispose
    dispose();

    await vi.waitFor(() => {
      expect(MockWebSocket.instances[0]?.closed).toBe(true);
      expect(connectedSignal()).toBe(false);
    });
  });
});
