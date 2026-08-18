export interface TracesIngestEnv {
  TRACES_ENABLED?: string | null;
  TRACES_OTLP_USER?: string | null;
  TRACES_OTLP_PASSWORD?: string | null;
  TRACES_OTLP_URL?: string | null;
}

const DEFAULT_TRACES_INGEST_URL = "https://traces.trieoh.com/insert/opentelemetry/v1/traces";

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

  const user = env.TRACES_OTLP_USER;
  const password = env.TRACES_OTLP_PASSWORD;
  if (!user || !password) {
    console.error(
      `[tracing] TRACES_OTLP_USER/PASSWORD not configured (TRACES_ENABLED=${env.TRACES_ENABLED})`,
    );
    return new Response(null, { status: 503 });
  }

  const auth = btoa(`${user}:${password}`);
  const targetUrl = env.TRACES_OTLP_URL ?? DEFAULT_TRACES_INGEST_URL;
  const contentType = request.headers.get("Content-Type") ?? "application/json";

  try {
    const upstream = await fetch(targetUrl, {
      method: "POST",
      headers: {
        Authorization: `Basic ${auth}`,
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
