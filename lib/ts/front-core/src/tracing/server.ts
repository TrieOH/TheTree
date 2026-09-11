import type { TracesIngestEnv } from "./ingest";

const DEFAULT_TRACES_INGEST_URL =
  "https://traces.trieoh.com/insert/opentelemetry/v1/traces";

const TRACEPARENT = /^([\da-f]{2})-([\da-f]{32})-([\da-f]{16})-([\da-f]{2})$/i;

export interface ServerSpanEvent {
  /** Span name, e.g. `bff.request` or `bff.refresh`. */
  name: string;
  /** Incoming W3C traceparent, so the span hangs off the caller's trace. */
  traceparent?: string | null;
  durationMs: number;
  success?: boolean;
  status?: number;
  errorId?: string;
  attributes?: Record<string, string | number | boolean>;
  serviceName?: string;
  /** Overridable for tests. */
  now?: number;
}

const randomHex = (bytes: number): string =>
  Array.from({ length: bytes }, () =>
    Math.floor(Math.random() * 256).toString(16).padStart(2, "0"),
  ).join("");

function attributeValue(value: string | number | boolean) {
  if (typeof value === "number") {
    return Number.isInteger(value)
      ? { intValue: String(value) }
      : { doubleValue: value };
  }
  if (typeof value === "boolean") return { boolValue: value };
  return { stringValue: value };
}

/**
 * Emits one span to the same OTLP collector the browser ingest uses, linking it
 * to the caller's trace when a `traceparent` is available.
 *
 * Never throws: tracing must not be able to fail a request, so a missing
 * configuration or a dead collector is a no-op.
 */
export async function emitServerSpan(
  event: ServerSpanEvent,
  env: TracesIngestEnv,
): Promise<void> {
  if (env.TRACES_ENABLED === "false") return;

  const user = env.TRACES_OTLP_USER;
  const password = env.TRACES_OTLP_PASSWORD;
  if (!user || !password) return;

  const match = event.traceparent?.match(TRACEPARENT);
  const traceId = match?.[2] ?? randomHex(16);
  const parentSpanId = match?.[3];
  const endedAt = event.now ?? Date.now();
  const startedAt = endedAt - Math.max(event.durationMs, 0);

  const attributes = Object.entries({
    ...(event.status === undefined ? {} : { "http.response.status_code": event.status }),
    ...(event.errorId ? { "bff.error_id": event.errorId } : {}),
    ...event.attributes,
  }).map(([key, value]) => ({ key, value: attributeValue(value) }));

  const payload = {
    resourceSpans: [
      {
        resource: {
          attributes: [
            {
              key: "service.name",
              value: { stringValue: event.serviceName ?? "univents-solid-bff" },
            },
          ],
        },
        scopeSpans: [
          {
            scope: { name: "@trieoh/front-core/auth/bff" },
            spans: [
              {
                traceId,
                spanId: randomHex(8),
                ...(parentSpanId ? { parentSpanId } : {}),
                name: event.name,
                kind: 2,
                startTimeUnixNano: String(startedAt * 1_000_000),
                endTimeUnixNano: String(endedAt * 1_000_000),
                attributes,
                status: { code: event.success === false ? 2 : 1 },
              },
            ],
          },
        ],
      },
    ],
  };

  try {
    await fetch(env.TRACES_OTLP_URL ?? DEFAULT_TRACES_INGEST_URL, {
      method: "POST",
      headers: {
        Authorization: `Basic ${btoa(`${user}:${password}`)}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });
  } catch (error) {
    console.error("[tracing] server span dropped", {
      name: event.name,
      error: error instanceof Error ? error.message : String(error),
    });
  }
}
