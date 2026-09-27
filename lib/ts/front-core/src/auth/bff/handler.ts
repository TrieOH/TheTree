import { createIdentityXBff } from "./core";
import { createCookieBffSession, type CookieSessionOptions } from "./cookie-session";
import { BffInputError, runBffOperation } from "./operations";
import type { IdentityXTransportLogEvent } from "./types";

export interface BffHandlerConfig {
  identityX: {
    baseURL: string;
    projectId?: string;
  };
  apiBaseURL: string;
  session: CookieSessionOptions;
  observability?: {
    log?: (event: IdentityXTransportLogEvent) => void;
    logSuccesses?: boolean;
  };
}

function json(body: unknown, status: number): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: {
      "Content-Type": "application/json",
      "Cache-Control": "no-store",
    },
  });
}

/**
 * Transport for the shared BFF: one POST endpoint, `{ op, input }`, dispatched
 * through {@link runBffOperation}. The session is created per request, so two
 * concurrent callers can never share tokens.
 */
export function createBffHandler(config: BffHandlerConfig) {
  return async function handleBffRequest(request: Request): Promise<Response> {
    let payload: unknown;
    try {
      payload = await request.json();
    } catch {
      return json({ success: false, code: 400, message: "Invalid JSON body" }, 400);
    }

    const { op, input } = (payload ?? {}) as { op?: unknown; input?: unknown };
    if (typeof op !== "string") {
      return json({ success: false, code: 400, message: "Missing operation" }, 400);
    }

    const cookieSession = await createCookieBffSession(request, config.session);
    const bff = createIdentityXBff({
      identityX: config.identityX,
      apiBaseURL: config.apiBaseURL,
      observability: config.observability,
      session: () => cookieSession.session,
      request: () => request,
    });

    try {
      const result = await runBffOperation(bff, op, input);
      return cookieSession.applyTo(json(result, 200));
    } catch (error) {
      if (error instanceof BffInputError) {
        return json({ success: false, code: 400, message: error.message }, 400);
      }

      console.error("[identityx-bff] operation failed", op, error);
      return json(
        { success: false, code: 500, message: "BFF operation failed" },
        500,
      );
    }
  };
}
