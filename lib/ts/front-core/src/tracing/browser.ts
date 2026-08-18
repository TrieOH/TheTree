import { context, propagation, trace, type Context, type Span } from "@opentelemetry/api"
import { resourceFromAttributes } from "@opentelemetry/resources"
import { ATTR_SERVICE_NAME } from "@opentelemetry/semantic-conventions"
import {
  BatchSpanProcessor,
  type SpanProcessor,
} from "@opentelemetry/sdk-trace-base"
import { WebTracerProvider } from "@opentelemetry/sdk-trace-web"
import { ZoneContextManager } from "@opentelemetry/context-zone"
import { FetchInstrumentation } from "@opentelemetry/instrumentation-fetch"
import { registerInstrumentations } from "@opentelemetry/instrumentation"
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http"

import { TRACES_INGEST_PATH } from "./constants"

const tracer = trace.getTracer("trieoh-front")

// One root span per page session; its context is reused for every request.
let sessionSpan: Span | undefined
let provider: WebTracerProvider | undefined

class ComponentSpanProcessor implements SpanProcessor {
  constructor(private readonly delegate: SpanProcessor) {}

  onStart(
    span: Parameters<SpanProcessor["onStart"]>[0],
    parentContext: Context,
  ): void {
    span.setAttribute("component", "web")
    this.delegate.onStart(span, parentContext)
  }

  onEnd(span: Parameters<SpanProcessor["onEnd"]>[0]): void {
    this.delegate.onEnd(span)
  }

  forceFlush(): Promise<void> {
    return this.delegate.forceFlush()
  }

  shutdown(): Promise<void> {
    return this.delegate.shutdown()
  }
}

export function browserTracingResource(serviceName: string) {
  return resourceFromAttributes({ [ATTR_SERVICE_NAME]: serviceName })
}

/** Falls back to the session span as the active span whenever nothing more
 *  specific is active — guarantees every fetch, no matter where it's
 *  triggered from, is parented to the same page-session trace. */
class SessionContextManager extends ZoneContextManager {
  active(): Context {
    const current = super.active()
    if (trace.getSpan(current) || !sessionSpan) return current
    return trace.setSpan(current, sessionSpan)
  }
}

export function initSessionTrace(): void {
  if (typeof window === "undefined" || sessionSpan) return
  sessionSpan = tracer.startSpan("page-session")
}

/** Returns the session's traceparent header value. Used by the BFF client
 *  (auth/tanstack/client.ts) to attach traceparent onto the proxy request. */
export function getTraceparent(): string {
  if (!sessionSpan) initSessionTrace()
  if (!sessionSpan) return ""

  const headers: Record<string, string> = {}
  const ctx = trace.setSpan(context.active(), sessionSpan)
  propagation.inject(ctx, headers)
  return headers["traceparent"] ?? ""
}

/** Phase 2: real browser spans (fetch, document load) exported through the
 *  BFF's ingestTracesServerFn — the only escape hatch to Victoria Traces. */
export function initBrowserTracing(serviceName = "univents-web"): void {
  if (typeof window === "undefined" || provider) return

  const exporter = new OTLPTraceExporter({
    url: `${window.location.origin}${TRACES_INGEST_PATH}`,
  })

  const processor = new BatchSpanProcessor(exporter, {
    maxExportBatchSize: 30,
    scheduledDelayMillis: 5000,
  })

  provider = new WebTracerProvider({
    resource: browserTracingResource(serviceName),
    spanProcessors: [new ComponentSpanProcessor(processor)],
  })

  provider.register({ contextManager: new SessionContextManager() })

  registerInstrumentations({
    instrumentations: [
      new FetchInstrumentation({
        ignoreUrls: [new RegExp(escapeRegExp(TRACES_INGEST_PATH))],
        propagateTraceHeaderCorsUrls: [],
      }),
    ],
  })

  initSessionTrace()
  flushOnPageHide(processor)
}

function flushOnPageHide(processor: BatchSpanProcessor): void {
  const flush = () => void processor.forceFlush()
  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "hidden") flush()
  })
  window.addEventListener("pagehide", flush)
}

function escapeRegExp(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
}
