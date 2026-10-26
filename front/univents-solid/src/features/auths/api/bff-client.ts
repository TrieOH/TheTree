import type { BffTransport } from "@trieoh/front-core/auth/bff/client";
import { BFF_PATH } from "@trieoh/front-core/auth/bff/constants";
import { getTraceparent } from "@trieoh/front-core/tracing/browser";
import { withBffResilience } from "@/shared/lib/resilient-transport";

const rawBffTransport: BffTransport = {
  async call<T>(op: string, input?: unknown): Promise<T> {
    const traceparent = getTraceparent();
    const response = await fetch(BFF_PATH, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(traceparent ? { traceparent } : {}),
      },
      credentials: "same-origin",
      body: JSON.stringify({ op, input }),
    });

    // Failures still carry the envelope (400 malformed input, 500 unexpected).
    return (await response.json()) as T;
  },
};

/**
 * Browser transport for the BFF. The session lives in an HttpOnly cookie, so
 * the client only ever sends the operation name and its input: no tokens, no
 * IdentityX URL. Reads get a deadline and a circuit breaker so a stalled
 * backend can't hold a page hostage.
 */
export const bffTransport = withBffResilience(rawBffTransport);
