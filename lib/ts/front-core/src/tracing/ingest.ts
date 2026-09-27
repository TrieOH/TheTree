export interface TracesIngestEnv {
  TRACES_ENABLED?: string | null;

  // OpenTelemetry official standard
  OTEL_EXPORTER_OTLP_ENDPOINT?: string | null;
  OTEL_EXPORTER_OTLP_HEADERS?: string | null;
  OTEL_EXPORTER_OTLP_PROTOCOL?: string | null;

  // Legacy fallback
  TRACES_OTLP_USER?: string | null;
  TRACES_OTLP_PASSWORD?: string | null;
  TRACES_OTLP_URL?: string | null;
}

const DEFAULT_TRACES_INGEST_URL =
  "https://traces.trieoh.com/insert/opentelemetry/v1/traces";

/**
 * Parses comma-separated key=value pairs into a headers record.
 * Handles both quoted and unquoted values (e.g. `Authorization=Basic dXNlcjpwYXNz,X-Scope=123`).
 */
export function parseOtlpHeaders(raw?: string | null): Record<string, string> {
  if (!raw) return {};
  const headers: Record<string, string> = {};
  for (const item of raw.split(",")) {
    const trimmed = item.trim();
    if (!trimmed) continue;
    const eqIdx = trimmed.indexOf("=");
    if (eqIdx !== -1) {
      const key = trimmed.slice(0, eqIdx).trim();
      let val = trimmed.slice(eqIdx + 1).trim();
      if (
        (val.startsWith('"') && val.endsWith('"')) ||
        (val.startsWith("'") && val.endsWith("'"))
      ) {
        val = val.slice(1, -1);
      }
      if (key) {
        headers[key] = val;
      }
    }
  }
  return headers;
}

export function resolveOtlpConfig(env: TracesIngestEnv): {
  targetUrl: string;
  headers: Record<string, string>;
  hasCredentials: boolean;
} {
  const targetUrl =
    env.OTEL_EXPORTER_OTLP_ENDPOINT?.trim() ||
    env.TRACES_OTLP_URL?.trim() ||
    DEFAULT_TRACES_INGEST_URL;

  const headers = parseOtlpHeaders(env.OTEL_EXPORTER_OTLP_HEADERS);

  const hasAuthHeader = Boolean(
    headers["Authorization"] || headers["authorization"],
  );

  const user = env.TRACES_OTLP_USER?.trim();
  const password = env.TRACES_OTLP_PASSWORD?.trim();

  if (!hasAuthHeader && user && password) {
    headers["Authorization"] = `Basic ${btoa(`${user}:${password}`)}`;
  }

  const hasCredentials =
    hasAuthHeader ||
    Boolean(user && password) ||
    Object.keys(headers).length > 0;

  return {
    targetUrl,
    headers,
    hasCredentials,
  };
}

export async function handleTracesIngest(
  request: Request,
  env: TracesIngestEnv,
): Promise<Response> {
  if (env.TRACES_ENABLED === "false") {
    return new Response(null, { status: 204 });
  }

  if (request.method !== "POST") {
    return new Response(null, { status: 405 });
  }

  const body = await request.arrayBuffer();
  if (body.byteLength === 0) {
    console.warn(
      "[tracing] empty body received for ingest; returning 204",
    );
    return new Response(null, { status: 204 });
  }

  const { targetUrl, headers, hasCredentials } = resolveOtlpConfig(env);

  if (!hasCredentials) {
    console.error(
      `[tracing] OTLP credentials not configured (set OTEL_EXPORTER_OTLP_HEADERS or TRACES_OTLP_USER/PASSWORD, TRACES_ENABLED=${env.TRACES_ENABLED})`,
    );
    return new Response(null, { status: 503 });
  }

  const contentType = request.headers.get("Content-Type") ?? "application/json";

  try {
    const upstream = await fetch(targetUrl, {
      method: "POST",
      headers: {
        ...headers,
        "Content-Type": contentType,
      },
      body,
    });
    if (!upstream.ok) {
      console.error("[tracing] upstream rejected", {
        status: upstream.status,
        contentType,
        bodyBytes: body.byteLength,
      });
    }
    return new Response(null, { status: upstream.ok ? 204 : 502 });
  } catch (error) {
    console.error("[tracing] ingest failed", {
      error: error instanceof Error ? error.message : String(error),
      bodyBytes: body.byteLength,
    });
    return new Response(null, { status: 502 });
  }
}
