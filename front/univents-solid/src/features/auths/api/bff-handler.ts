import { BFF_PATH } from "@trieoh/front-core/auth/bff/constants";
import { createBffHandler } from "@trieoh/front-core/auth/bff/handler";
import { emitServerSpan } from "@trieoh/front-core/tracing/server";

export { BFF_PATH };

/**
 * The BFF only needs infrastructure config: IdentityX URLs, and the secret that
 * seals the session cookie. Endpoints, envelopes and refresh stay in the shared
 * core.
 */
export function handleBffRequest(
  request: Request,
  env: Env,
  ctx?: ExecutionContext,
): Promise<Response> {
  return createBffHandler({
    identityX: {
      baseURL: env.VITE_AUTH_API_URL,
      projectId: env.VITE_TRIEOH_AUTH_PROJECT_ID,
    },
    apiBaseURL: env.VITE_API_URL,
    session: {
      name: "univents-solid-auth",
      password: env.AUTH_SESSION_PASSWORD,
      secure: true,
    },
    // One span per outbound IdentityX call, linked to the caller's trace.
    // `waitUntil` keeps the isolate alive until the collector answers.
    observability: {
      log: (event) => {
        const span = emitServerSpan(
          {
            name: `bff.${event.operation}`,
            traceparent: event.traceparent,
            durationMs: event.duration_ms,
            success: event.success,
            status: event.status,
            errorId: event.error_id,
            attributes: {
              "bff.operation": event.operation,
              "bff.method": event.method,
              "bff.path": event.path,
            },
          },
          env,
        );

        if (ctx) ctx.waitUntil(span);
        else void span;
      },
    },
  })(request);
}
