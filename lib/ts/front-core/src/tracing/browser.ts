import { context, propagation, SpanStatusCode, trace, type Context, type Span } from "@opentelemetry/api"
import { resourceFromAttributes } from "@opentelemetry/resources"
import { ATTR_SERVICE_NAME } from "@opentelemetry/semantic-conventions"
import {
  BatchSpanProcessor,
  type ReadableSpan,
  type SpanExporter,
  type SpanProcessor,
} from "@opentelemetry/sdk-trace-base"
import { WebTracerProvider } from "@opentelemetry/sdk-trace-web"
import { ZoneContextManager } from "@opentelemetry/context-zone"
import { FetchInstrumentation } from "@opentelemetry/instrumentation-fetch"
import { registerInstrumentations } from "@opentelemetry/instrumentation"
import { JsonTraceSerializer } from "@opentelemetry/otlp-transformer/build/src/trace/json/trace"

import { TRACES_INGEST_PATH } from "./constants"

const tracer = trace.getTracer("trieoh-front")

let provider: WebTracerProvider | undefined
let pageSpan: Span | undefined
let pagePath: string | undefined
let pageSpanEnded = false

class ComponentSpanProcessor implements SpanProcessor {
  constructor(private readonly delegate: SpanProcessor) { }

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

class KeepaliveSpanExporter implements SpanExporter {
  constructor(private readonly url: string) { }

  export(
    spans: ReadableSpan[],
    resultCallback: Parameters<SpanExporter["export"]>[1],
  ): void {
    let payload: Uint8Array | undefined
    try {
      payload = JsonTraceSerializer.serializeRequest(spans)
    } catch (error) {
      const normalized = error instanceof Error ? error : new Error(String(error))
      console.error("[tracing] failed to serialize spans", normalized)
      resultCallback({ code: 1, error: normalized })
      return
    }

    void fetch(this.url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: payload ? new TextDecoder().decode(payload) : undefined,
      keepalive: true,
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Trace ingest failed: ${response.status}`)
        }
        resultCallback({ code: 0 })
      })
      .catch((error: unknown) => {
        console.error("[tracing] failed to export spans", error)
        resultCallback({
          code: 1,
          error: error instanceof Error ? error : new Error(String(error)),
        })
      })
  }

  shutdown(): Promise<void> {
    return Promise.resolve()
  }
}

export function browserTracingResource(serviceName: string) {
  return resourceFromAttributes({ [ATTR_SERVICE_NAME]: serviceName })
}

export function initSessionTrace(): void {
  // Kept for backwards compatibility; page-session spans are disabled.
}

function startPageTrace(pathname = window.location.pathname): void {
  if (pageSpan) {
    pageSpan.addEvent("page.unloaded")
    if (!pageSpanEnded) pageSpan.end()
  }

  pagePath = pathname
  pageSpan = tracer.startSpan(`page: ${pathname}`)
  pageSpanEnded = false
  pageSpan.addEvent("page.loaded")
  setTimeout(() => {
    if (!pageSpan || pageSpanEnded || pagePath !== window.location.pathname) return
    pageSpan.addEvent("page.ready")
    pageSpan.end()
    pageSpanEnded = true
  }, 0)
}

class PageContextManager extends ZoneContextManager {
  active(): Context {
    const current = super.active()
    if (trace.getSpan(current) || !pageSpan) return current
    return trace.setSpan(current, pageSpan)
  }
}

/** Propagates only the currently active operation; no page-session span is used. */
export function getTraceparent(): string {
  if (!trace.getSpan(context.active())) return ""

  const headers: Record<string, string> = {}
  propagation.inject(context.active(), headers)
  return headers.traceparent ?? ""
}

export function setActiveSpanAttributes(
  attributes: Record<string, string | number | boolean>,
): void {
  trace.getSpan(context.active())?.setAttributes(attributes)
}

export function addActiveSpanEvent(
  name: string,
  attributes?: Record<string, string | number | boolean>,
): void {
  trace.getSpan(context.active())?.addEvent(name, attributes)
}

export async function withSpan<T>(
  name: string,
  fn: () => Promise<T> | T,
): Promise<T> {
  const span = tracer.startSpan(name)
  const spanContext = trace.setSpan(context.active(), span)
  const startedAt = performance.now()
  span.addEvent("operation.started")
  try {
    const result = await context.with(spanContext, fn)
    span.setAttribute("operation.duration_ms", performance.now() - startedAt)
    span.setAttribute("operation.outcome", "success")
    span.addEvent("operation.completed")
    return result
  } catch (error) {
    const normalized = error instanceof Error ? error : new Error(String(error))
    span.recordException(normalized)
    span.setAttribute("operation.duration_ms", performance.now() - startedAt)
    span.setAttribute("operation.outcome", "failure")
    span.addEvent("operation.failed", { message: normalized.message })
    span.setStatus({ code: SpanStatusCode.ERROR, message: normalized.message })
    throw error
  } finally {
    span.end()
  }
}

export function recordCompletedSpan(
  name: string,
  startedAt: number,
  attributes: Record<string, string | number | boolean>,
  error?: unknown,
): void {
  const endedAt = Date.now()
  const span = tracer.startSpan(name, { startTime: startedAt, attributes })
  span.setAttribute("operation.duration_ms", endedAt - startedAt)
  if (error === undefined) {
    span.setAttribute("operation.outcome", "success")
    span.addEvent("operation.completed")
  } else {
    const normalized = error instanceof Error ? error : new Error(String(error))
    span.recordException(normalized)
    span.setAttribute("operation.outcome", "failure")
    span.addEvent("operation.failed", { message: normalized.message })
    span.setStatus({ code: SpanStatusCode.ERROR, message: normalized.message })
  }
  span.end(endedAt)
}

/** Phase 2: real browser spans (fetch, document load) exported through the
 *  BFF's ingestTracesServerFn — the only escape hatch to Victoria Traces. */
export function initBrowserTracing(
  serviceName = "univents",
  enabled = true,
  ignoredUrls: string[] = [],
): void {
  if (typeof window === "undefined" || !enabled || provider) return

  const exporter = new KeepaliveSpanExporter(
    `${window.location.origin}${TRACES_INGEST_PATH}`,
  )

  const processor = new BatchSpanProcessor(exporter, {
    maxExportBatchSize: 30,
    scheduledDelayMillis: 5000,
  })

  provider = new WebTracerProvider({
    resource: browserTracingResource(serviceName),
    spanProcessors: [new ComponentSpanProcessor(processor)],
  })

  provider.register({ contextManager: new PageContextManager() })

  initSessionTrace()
  startPageTrace()
  registerInstrumentations({
    instrumentations: [
      new FetchInstrumentation({
        ignoreUrls: [
          new RegExp(escapeRegExp(TRACES_INGEST_PATH)),
          /\/__tsd\//,
          ...ignoredUrls.map((url) => new RegExp(escapeRegExp(url))),
        ],
        propagateTraceHeaderCorsUrls: [],
        applyCustomAttributesOnSpan: (span, request, result) => {
          const method = request.method ?? "GET"
          const rawUrl = result instanceof Response
            ? result.url
            : request instanceof Request
              ? request.url
              : undefined
          const path = rawUrl
            ? new URL(rawUrl, window.location.origin).pathname
            : undefined
          const serverFunction = path ? getServerFunctionName(path) : undefined
          span.updateName(
            serverFunction
              ? `BFF ${serverFunction}`
              : `HTTP ${method}${path ? ` ${path}` : ""}`,
          )
        },
      }),
    ],
  })
  flushOnPageHide(processor)
}

function flushOnPageHide(processor: BatchSpanProcessor): void {
  const navigate = () => {
    if (window.location.pathname !== pagePath) startPageTrace()
  }

  window.addEventListener("popstate", navigate)
  window.addEventListener("hashchange", navigate)
  const pushState = history.pushState
  const replaceState = history.replaceState
  history.pushState = function (...args) {
    pushState.apply(history, args)
    navigate()
  }
  history.replaceState = function (...args) {
    replaceState.apply(history, args)
    navigate()
  }

  const flush = () => {
    if (pageSpan) {
      pageSpan.addEvent("page.hidden")
      if (!pageSpanEnded) pageSpan.end()
      pageSpanEnded = true
      pageSpan = undefined
    }
    void processor.forceFlush()
  }
  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "hidden") void processor.forceFlush()
  })
  window.addEventListener("pagehide", flush)
}

function escapeRegExp(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
}

function getServerFunctionName(path: string): string | undefined {
  const encoded = path.match(/^\/_serverFn\/([^/]+)$/)?.[1]
  if (!encoded) return undefined

  try {
    const base64 = encoded.replace(/-/g, "+").replace(/_/g, "/")
    const decoded = atob(base64.padEnd(Math.ceil(base64.length / 4) * 4, "="))
    const payload = JSON.parse(decoded) as { export?: unknown }
    return typeof payload.export === "string"
      ? payload.export.replace(/_createServerFn_handler$/, "")
      : "server-function"
  } catch {
    return "server-function"
  }
}
