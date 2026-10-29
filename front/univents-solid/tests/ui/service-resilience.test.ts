import { QueryClient } from "@trieoh/front-core-solid";
import type { ApiTransport } from "@trieoh/api-client";
import type { BffTransport } from "@trieoh/front-core/auth/bff/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { shouldRetryRead } from "@/shared/lib/query-client";
import {
  READ_DEADLINE_MS,
  SESSION_DEADLINE_MS,
  withApiResilience,
  withBffResilience,
  withDeadline,
} from "@/shared/lib/resilient-transport";
import { boundLoader, peekQuery, warmQueries } from "@/shared/lib/route-data";
import { classifyError } from "@/shared/lib/service-error";
import {
  COOL_DOWN_MS,
  OPEN_AFTER,
  CircuitOpenError,
  createServiceHealth,
  degradedServices,
} from "@/shared/lib/service-health";

beforeEach(() => {
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
});

const envelopeError = (code: number, error_id?: string) =>
  Object.assign(new Error("failed"), { envelope: { code, error_id } });

describe("classifyError", () => {
  it("reads the envelope the fetch client and QueryError both carry", () => {
    expect(classifyError(envelopeError(503))).toMatchObject({ kind: "unavailable", transient: true });
    expect(classifyError(envelopeError(500))).toMatchObject({ kind: "unavailable", transient: true });
    expect(classifyError(envelopeError(408))).toMatchObject({ kind: "timeout", transient: true });
    expect(classifyError(envelopeError(0, "CLIENT_OFFLINE"))).toMatchObject({ kind: "offline", transient: true });
    expect(classifyError(envelopeError(0, "CLIENT_REQUEST_FAILED"))).toMatchObject({ kind: "unavailable" });
  });

  it("treats definitive answers as not transient", () => {
    expect(classifyError(envelopeError(401))).toMatchObject({ kind: "unauthorized", transient: false });
    expect(classifyError(envelopeError(403))).toMatchObject({ kind: "forbidden", transient: false });
    expect(classifyError(envelopeError(404))).toMatchObject({ kind: "not-found", transient: false });
    expect(classifyError(envelopeError(422))).toMatchObject({ kind: "client", transient: false });
  });

  it("recognises the BFF's refresh failure as the auth service being down", () => {
    const error = envelopeError(503, "AUTH_SERVICE_UNAVAILABLE");
    expect(classifyError(error)).toMatchObject({ kind: "unavailable", transient: true });
  });

  it("treats a bare fetch TypeError as the server being unreachable", () => {
    expect(classifyError(new TypeError("Failed to fetch"))).toMatchObject({
      kind: "unavailable",
      transient: true,
    });
  });

  it("does not mistake a bug in our own code for an outage", () => {
    expect(classifyError(new TypeError("x is not a function"))).toMatchObject({
      kind: "unknown",
      transient: false,
    });
  });

  it("does not throw on junk", () => {
    expect(classifyError(undefined).kind).toBe("unknown");
    expect(classifyError("boom").kind).toBe("unknown");
    expect(classifyError({}).transient).toBe(false);
  });
});

describe("service health", () => {
  const failure = classifyError(envelopeError(503));

  it("opens the circuit after consecutive transient failures, then lets a probe through", () => {
    let now = 1_000;
    const health = createServiceHealth(() => now);
    const probes: string[] = [];
    health.onProbe((service) => probes.push(service));

    for (let i = 0; i < OPEN_AFTER; i++) health.failure("api", failure);
    expect(() => health.guard("api")).toThrow(CircuitOpenError);
    // The other service is unaffected: one microservice dying stalls nothing else.
    expect(() => health.guard("identityx")).not.toThrow();

    now += COOL_DOWN_MS + 1;
    expect(() => health.guard("api")).not.toThrow();
    vi.advanceTimersByTime(COOL_DOWN_MS);
    expect(probes).toEqual(["api"]);
  });

  it("closes on a single success", () => {
    const health = createServiceHealth(() => 0);
    for (let i = 0; i < OPEN_AFTER; i++) health.failure("api", failure);
    health.success("api");
    expect(() => health.guard("api")).not.toThrow();
    expect(degradedServices(health.getSnapshot())).toEqual([]);
  });

  it("ignores failures that say nothing about health", () => {
    const health = createServiceHealth(() => 0);
    for (let i = 0; i < 10; i++) health.failure("api", classifyError(envelopeError(404)));
    expect(health.getSnapshot().api.failures).toBe(0);
  });

  it("reports degraded before the circuit opens", () => {
    const health = createServiceHealth(() => 0);
    health.failure("identityx", failure);
    expect(degradedServices(health.getSnapshot())).toEqual([]);
    health.failure("identityx", failure);
    expect(degradedServices(health.getSnapshot())).toEqual(["identityx"]);
  });

  it("halfOpenAll lets an explicit retry through and re-opens on failure", () => {
    const health = createServiceHealth(() => 0);
    for (let i = 0; i < OPEN_AFTER; i++) health.failure("api", failure);
    health.halfOpenAll();
    expect(() => health.guard("api")).not.toThrow();
    health.failure("api", failure);
    expect(() => health.guard("api")).toThrow(CircuitOpenError);
  });

  it("notifies subscribers with a new snapshot", () => {
    const health = createServiceHealth(() => 0);
    const listener = vi.fn();
    const before = health.getSnapshot();
    health.subscribe(listener);
    health.failure("api", failure);
    expect(listener).toHaveBeenCalledTimes(1);
    expect(health.getSnapshot()).not.toBe(before);
  });
});

describe("withDeadline", () => {
  it("rejects with a timeout envelope but lets the request finish", async () => {
    let finish!: (value: string) => void;
    const slow = new Promise<string>((resolve) => { finish = resolve; });
    const result = withDeadline(slow, 100);
    const caught = result.catch((error: unknown) => error);
    await vi.advanceTimersByTimeAsync(100);
    expect(classifyError(await caught)).toMatchObject({ kind: "timeout", transient: true });
    // The underlying promise is untouched, so a late response is still processed.
    finish("late");
    await expect(slow).resolves.toBe("late");
  });

  it("does not leave an unhandled rejection when the request fails after the deadline", async () => {
    const rejections: unknown[] = [];
    const onRejection = (reason: unknown) => rejections.push(reason);
    process.on("unhandledRejection", onRejection);
    let fail!: (error: Error) => void;
    const slow = new Promise<string>((_, reject) => { fail = reject; });
    const caught = withDeadline(slow, 50).catch(() => "timeout");
    await vi.advanceTimersByTimeAsync(50);
    await caught;
    fail(new Error("late failure"));
    await vi.advanceTimersByTimeAsync(0);
    process.off("unhandledRejection", onRejection);
    expect(rejections).toEqual([]);
  });
});

describe("withApiResilience", () => {
  const ok = { data: [], status: 200, headers: new Headers() };
  const transport = (impl: () => Promise<unknown>) =>
    vi.fn(impl) as unknown as ApiTransport & ReturnType<typeof vi.fn>;

  it("cuts a stalled read off at the deadline and reports the service as struggling", async () => {
    const health = createServiceHealth(() => 0);
    const stalled = transport(() => new Promise(() => undefined));
    const wrapped = withApiResilience(stalled, "api", health);

    const result = wrapped("/events", { method: "GET" }).catch((error: unknown) => error);
    await vi.advanceTimersByTimeAsync(READ_DEADLINE_MS);

    expect(classifyError(await result).kind).toBe("timeout");
    expect(health.getSnapshot().api.failures).toBe(1);
  });

  it("fails reads fast once the circuit is open, without touching the network", async () => {
    const health = createServiceHealth(() => 0);
    const failing = transport(() => Promise.reject(envelopeError(503)));
    const wrapped = withApiResilience(failing, "api", health);

    for (let i = 0; i < OPEN_AFTER; i++) {
      await expect(wrapped("/events", { method: "GET" })).rejects.toBeDefined();
    }
    failing.mockClear();

    await expect(wrapped("/events", { method: "GET" })).rejects.toBeInstanceOf(CircuitOpenError);
    expect(failing).not.toHaveBeenCalled();
  });

  it("never short-circuits or times out writes", async () => {
    const health = createServiceHealth(() => 0);
    const failing = transport(() => Promise.reject(envelopeError(503)));
    const wrapped = withApiResilience(failing, "api", health);
    for (let i = 0; i < OPEN_AFTER; i++) {
      await expect(wrapped("/events", { method: "GET" })).rejects.toBeDefined();
    }
    failing.mockClear();

    await expect(wrapped("/events", { method: "POST" })).rejects.toBeDefined();
    expect(failing).toHaveBeenCalledTimes(1);
  });

  it("counts a 404 as the service being alive", async () => {
    const health = createServiceHealth(() => 0);
    const wrapped = withApiResilience(transport(() => Promise.reject(envelopeError(404))), "api", health);
    await expect(wrapped("/events/x", { method: "GET" })).rejects.toBeDefined();
    expect(health.getSnapshot().api.failures).toBe(0);
  });

  it("passes successful responses through and closes the circuit", async () => {
    const health = createServiceHealth(() => 0);
    health.failure("api", classifyError(envelopeError(503)));
    const wrapped = withApiResilience(transport(() => Promise.resolve(ok)), "api", health);
    await expect(wrapped("/events", { method: "GET" })).resolves.toBe(ok);
    expect(health.getSnapshot().api.failures).toBe(0);
  });
});

describe("withBffResilience", () => {
  const bff = (impl: (op: string, input?: unknown) => Promise<unknown>) => {
    const call = vi.fn(impl);
    return { transport: { call } as unknown as BffTransport, call };
  };

  it("turns a stalled authenticated read into a timeout result instead of hanging", async () => {
    const health = createServiceHealth(() => 0);
    const { transport } = bff(() => new Promise(() => undefined));
    const wrapped = withBffResilience(transport, health);

    const result = wrapped.call<{ success: boolean; code: number; error_id?: string }>("request", { path: "/me" });
    await vi.advanceTimersByTimeAsync(READ_DEADLINE_MS);

    await expect(result).resolves.toMatchObject({ success: false, code: 408, error_id: "CLIENT_TIMEOUT" });
    expect(health.getSnapshot().api.failures).toBe(1);
  });

  it("attributes IdentityX-targeted requests to IdentityX, not the API", async () => {
    const health = createServiceHealth(() => 0);
    const { transport } = bff(() => Promise.reject(new TypeError("Failed to fetch")));
    const wrapped = withBffResilience(transport, health);

    await wrapped.call("request", { path: "/actors/1/profile", target: "identityx" });

    expect(health.getSnapshot().identityx.failures).toBe(1);
    expect(health.getSnapshot().api.failures).toBe(0);
  });

  it("blames IdentityX when a refresh failure fails an API request", async () => {
    const health = createServiceHealth(() => 0);
    const { transport } = bff(async () => ({
      success: false,
      code: 503,
      error_id: "AUTH_SERVICE_UNAVAILABLE",
    }));
    await withBffResilience(transport, health).call("request", { path: "/events" });

    expect(health.getSnapshot().identityx.failures).toBe(1);
    expect(health.getSnapshot().api.failures).toBe(0);
  });

  it("fails reads fast while the circuit is open, but lets writes through", async () => {
    const health = createServiceHealth(() => 0);
    const { transport, call } = bff(async () => ({ success: false, code: 503 }));
    const wrapped = withBffResilience(transport, health);
    for (let i = 0; i < OPEN_AFTER; i++) await wrapped.call("request", { path: "/events" });
    call.mockClear();

    await expect(wrapped.call("request", { path: "/events" })).resolves.toMatchObject({
      success: false,
      error_id: "CIRCUIT_OPEN",
    });
    expect(call).not.toHaveBeenCalled();

    await wrapped.call("request", { path: "/events", method: "POST" });
    expect(call).toHaveBeenCalledTimes(1);
  });

  it("never short-circuits session restore, so an outage can't log people out", async () => {
    const health = createServiceHealth(() => 0);
    const { transport, call } = bff(async () => ({ isAuthenticated: true, profile: null }));
    for (let i = 0; i < OPEN_AFTER; i++) {
      health.failure("identityx", classifyError(envelopeError(503)));
    }
    const wrapped = withBffResilience(transport, health);

    await expect(wrapped.call("restore")).resolves.toMatchObject({ isAuthenticated: true });
    expect(call).toHaveBeenCalledTimes(1);
  });

  it("still rejects a stalled restore, which the SDK knows how to handle", async () => {
    const health = createServiceHealth(() => 0);
    const { transport } = bff(() => new Promise(() => undefined));
    const result = withBffResilience(transport, health).call("restore").catch((error: unknown) => error);
    await vi.advanceTimersByTimeAsync(SESSION_DEADLINE_MS);
    expect(classifyError(await result).kind).toBe("timeout");
  });

  it("does not time out writes", async () => {
    const health = createServiceHealth(() => 0);
    let finish!: (value: unknown) => void;
    const { transport } = bff(() => new Promise((resolve) => { finish = resolve; }));
    const result = withBffResilience(transport, health).call("request", { path: "/x", method: "POST" });
    await vi.advanceTimersByTimeAsync(120_000);
    finish({ success: true, code: 200 });
    await expect(result).resolves.toMatchObject({ success: true });
  });
});

describe("query retry policy", () => {
  it("retries a server failure once, but not timeouts or definitive answers", () => {
    expect(shouldRetryRead(0, envelopeError(503))).toBe(true);
    expect(shouldRetryRead(1, envelopeError(503))).toBe(false);
    expect(shouldRetryRead(0, envelopeError(408, "CLIENT_TIMEOUT"))).toBe(false);
    expect(shouldRetryRead(0, envelopeError(404))).toBe(false);
    expect(shouldRetryRead(0, envelopeError(401))).toBe(false);
    expect(shouldRetryRead(0, envelopeError(503, "CIRCUIT_OPEN"))).toBe(false);
  });
});

describe("route data helpers", () => {
  const query = (key: string, fn: () => Promise<unknown>) => ({ queryKey: [key], queryFn: fn });

  it("warmQueries stops waiting after the bound while the request keeps going", async () => {
    const client = new QueryClient();
    let finish!: (value: string) => void;
    const slow = query("slow", () => new Promise<string>((resolve) => { finish = resolve; }));

    let settled = false;
    void warmQueries(client, [slow], 250).then(() => { settled = true; });
    await vi.advanceTimersByTimeAsync(249);
    expect(settled).toBe(false);
    await vi.advanceTimersByTimeAsync(1);
    expect(settled).toBe(true);

    finish("late");
    await vi.advanceTimersByTimeAsync(0);
    expect(client.getQueryData(["slow"])).toBe("late");
  });

  it("warmQueries returns at once for cached data and refreshes it in the background", async () => {
    const client = new QueryClient({ defaultOptions: { queries: { staleTime: 0 } } });
    client.setQueryData(["events"], ["old"]);
    const fetcher = vi.fn(() => new Promise<string[]>(() => undefined));

    await warmQueries(client, [query("events", fetcher)], 250);

    expect(client.getQueryData(["events"])).toEqual(["old"]);
    expect(fetcher).toHaveBeenCalledTimes(1);
  });

  it("warmQueries resolves instead of rejecting when the API fails", async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    await expect(
      warmQueries(client, [query("broken", () => Promise.reject(envelopeError(503)))]),
    ).resolves.toBeUndefined();
  });

  it("boundLoader never rejects", async () => {
    await expect(boundLoader(Promise.reject(new Error("nope")))).resolves.toBeUndefined();
  });

  it("peekQuery serves the cache, distinguishing a known-missing value from an unknown one", async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    client.setQueryData(["profile", "missing"], null);
    const never = () => new Promise<unknown>(() => undefined);

    await expect(peekQuery(client, { queryKey: ["profile", "missing"], queryFn: never })).resolves.toBeNull();

    const unknown = peekQuery(client, { queryKey: ["profile", "unknown"], queryFn: never }, 100);
    await vi.advanceTimersByTimeAsync(100);
    await expect(unknown).resolves.toBeUndefined();
  });

  it("peekQuery resolves to unknown, not an error, when the profile service is down", async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const result = peekQuery(client, {
      queryKey: ["profile", "down"],
      queryFn: () => Promise.reject(envelopeError(503)),
    });
    await expect(result).resolves.toBeUndefined();
  });
});
