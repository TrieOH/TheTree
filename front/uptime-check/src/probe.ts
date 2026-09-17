import { CHECK } from "./config";

/** Raw outcome of one attempt — the check engine turns these into Points. */
export interface ProbeResult {
  ok: boolean;
  /** Epoch ms of the attempt start. */
  t: number;
  /** Round-trip ms. */
  ms: number;
  /** HTTP status or error message, present on failure. */
  err?: string;
}

/**
 * One probe of a health endpoint. Never throws — a failed probe is a
 * result (HTTP status or error message), not an exception.
 */
export async function probe(url: string): Promise<ProbeResult> {
  const t = Date.now();
  try {
    const res = await fetch(url, {
      signal: AbortSignal.timeout(CHECK.timeoutMs),
      redirect: "follow",
    });
    if (res.ok) {
      return { ok: true, t, ms: Date.now() - t };
    }
    // Anything non-2xx/3xx fails — incl. a 502 from Caddy when the
    // service behind it is dead.
    return { ok: false, t, ms: Date.now() - t, err: `HTTP ${res.status}` };
  } catch (e) {
    return { ok: false, t, ms: Date.now() - t, err: message(e) };
  }
}

function message(e: unknown): string {
  if (e instanceof Error) {
    // AbortSignal.timeout produces a TimeoutError with a noisy DOM message.
    return e.name === "TimeoutError"
      ? `timeout after ${CHECK.timeoutMs}ms`
      : e.message;
  }
  return String(e);
}
