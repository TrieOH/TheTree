import { describe, expect, it, vi } from "vitest";
import { emitServerSpan } from "@trieoh/front-core/tracing/server";

const env = {
  TRACES_ENABLED: "true",
  TRACES_OTLP_URL: "https://collector.test/v1/traces",
  TRACES_OTLP_USER: "user",
  TRACES_OTLP_PASSWORD: "password",
};

const TRACEPARENT = "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01";

const captureSpan = () => {
  const spy = vi.fn(async (_url: string, _init?: RequestInit) =>
    Response.json({}, { status: 200 }),
  );
  vi.stubGlobal("fetch", spy);
  return spy;
};

const spanFrom = (spy: ReturnType<typeof captureSpan>) => {
  const init = spy.mock.calls[0]?.[1] as RequestInit;
  const payload = JSON.parse(String(init.body)) as {
    resourceSpans: Array<{
      resource: { attributes: Array<{ key: string; value: { stringValue: string } }> };
      scopeSpans: Array<{ spans: Array<Record<string, unknown>> }>;
    }>;
  };
  return {
    init,
    resource: payload.resourceSpans[0]!.resource,
    span: payload.resourceSpans[0]!.scopeSpans[0]!.spans[0]!,
  };
};

describe("server span emission", () => {
  it("links the span to the caller's trace", async () => {
    const spy = captureSpan();

    await emitServerSpan(
      {
        name: "bff.refresh",
        traceparent: TRACEPARENT,
        durationMs: 42,
        success: true,
        status: 200,
        attributes: { "bff.operation": "refresh" },
        now: 1_700_000_000_000,
      },
      env,
    );

    const { init, resource, span } = spanFrom(spy);

    expect(init.headers).toMatchObject({
      Authorization: `Basic ${btoa("user:password")}`,
      "Content-Type": "application/json",
    });
    expect(resource.attributes[0]).toEqual({
      key: "service.name",
      value: { stringValue: "univents-solid-bff" },
    });
    expect(span).toMatchObject({
      traceId: "0123456789abcdef0123456789abcdef",
      parentSpanId: "0123456789abcdef",
      name: "bff.refresh",
      kind: 2,
      // 42ms before `now`, in nanoseconds.
      startTimeUnixNano: String((1_700_000_000_000 - 42) * 1_000_000),
      endTimeUnixNano: String(1_700_000_000_000 * 1_000_000),
      status: { code: 1 },
    });
    expect(span.spanId).toMatch(/^[\da-f]{16}$/);
    expect(span.attributes).toEqual(
      expect.arrayContaining([
        { key: "http.response.status_code", value: { intValue: "200" } },
        { key: "bff.operation", value: { stringValue: "refresh" } },
      ]),
    );
  });

  it("starts a new trace when the caller sends no traceparent", async () => {
    const spy = captureSpan();

    await emitServerSpan({ name: "bff.login", durationMs: 5, success: false }, env);

    const { span } = spanFrom(spy);

    expect(span.traceId).toMatch(/^[\da-f]{32}$/);
    expect(span.parentSpanId).toBeUndefined();
    expect(span.status).toEqual({ code: 2 });
  });

  it("never throws when tracing is off or the collector is unreachable", async () => {
    const spy = captureSpan();

    await emitServerSpan({ name: "bff.login", durationMs: 1 }, { ...env, TRACES_ENABLED: "false" });
    await emitServerSpan({ name: "bff.login", durationMs: 1 }, { TRACES_ENABLED: "true" });

    expect(spy).not.toHaveBeenCalled();

    vi.stubGlobal("fetch", vi.fn(async (_url: string, _init?: RequestInit) => {
      throw new Error("collector down");
    }));
    await expect(
      emitServerSpan({ name: "bff.login", durationMs: 1 }, env),
    ).resolves.toBeUndefined();
  });
});
