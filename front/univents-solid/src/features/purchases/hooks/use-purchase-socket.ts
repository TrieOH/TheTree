import { createEffect, createSignal } from "solid-js";
import { Data, Effect, Fiber, Schedule } from "effect";
import { appSessionStorage } from "@/shared/lib/browser-storage";
import { purchaseWsToken } from "../api";

export class PurchaseWsTokenError extends Data.TaggedError("PurchaseWsTokenError")<{
  readonly cause: unknown;
  readonly message: string;
}> {}

export class PurchaseSocketConnectionError extends Data.TaggedError("PurchaseSocketConnectionError")<{
  readonly message?: string;
}> {}

export const isTerminalPurchaseFrame = (frame: {
  type?: string;
  payload?: { status?: string };
}): boolean => {
  const terminalTypes = [
    "purchase.confirmed",
    "purchase.expired",
    "purchase.cancelled",
    "purchase.failed",
    "purchase.rejected",
  ];
  return (
    terminalTypes.includes(frame.type ?? "") ||
    (frame.type === "purchase.snapshot" && frame.payload?.status !== "pending")
  );
};

export const defaultPurchaseSocketRetrySchedule = Schedule.exponential("500 millis").pipe(
  Schedule.jittered,
  Schedule.upTo({ times: 10 }),
);

export const resolveWsTokenEffect = (
  purchaseId: string,
): Effect.Effect<string, PurchaseWsTokenError> =>
  Effect.gen(function* () {
    const storedToken = appSessionStorage.take(`purchase-ws:${purchaseId}`);
    if (storedToken) return storedToken;

    const res = yield* Effect.tryPromise({
      try: () => purchaseWsToken(purchaseId),
      catch: (cause) =>
        new PurchaseWsTokenError({
          cause,
          message: "Não foi possível obter o token do WebSocket.",
        }),
    });

    return res.token;
  });

export interface PurchaseSocketSessionParams {
  readonly purchaseId: string;
  readonly onChange: (purchaseId: string) => void;
  readonly onConnectionChange?: (connected: boolean) => void;
  readonly retrySchedule?: Schedule.Schedule<unknown, unknown>;
  readonly apiBaseUrl?: string;
  readonly WebSocketClass?: typeof WebSocket;
}

const safeCloseSocket = (socket: WebSocket) => {
  try {
    socket.close();
  } catch {
    // Socket was already closed or disposed
  }
};

export const purchaseSocketSessionEffect = (
  params: PurchaseSocketSessionParams,
): Effect.Effect<void, unknown> =>
  Effect.gen(function* () {
    const singleAttempt = Effect.scoped(
      Effect.gen(function* () {
        const token = yield* resolveWsTokenEffect(params.purchaseId);
        const baseUrl =
          params.apiBaseUrl ??
          import.meta.env.VITE_API_URL ??
          "http://localhost:8081";
        const url = new URL("ws", baseUrl);
        url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
        url.searchParams.set("token", token);

        const WS = params.WebSocketClass ?? globalThis.WebSocket;
        const socket = yield* Effect.acquireRelease(
          Effect.sync(() => new WS(url)),
          (ws) =>
            Effect.sync(() => {
              safeCloseSocket(ws);
              params.onConnectionChange?.(false);
            }),
        );

        yield* Effect.callback<void, PurchaseSocketConnectionError>((resume, signal) => {
          socket.onopen = () => {
            params.onConnectionChange?.(true);
          };

          socket.onmessage = (event) => {
            try {
              const frame = JSON.parse(String(event.data)) as {
                type?: string;
                payload?: { status?: string };
              };
              if (isTerminalPurchaseFrame(frame)) {
                resume(Effect.void);
                return;
              }
            } catch {
              // Frame was not JSON (e.g. heartbeat ping/pong).
              // Proceed to onChange notification without failing.
            } finally {
              params.onChange(params.purchaseId);
            }
          };

          socket.onerror = () => {
            // onclose will handle disconnection and schedule retry
          };

          socket.onclose = () => {
            params.onConnectionChange?.(false);
            resume(
              Effect.fail(
                new PurchaseSocketConnectionError({
                  message: "WebSocket fechado inesperadamente.",
                }),
              ),
            );
          };

          signal.addEventListener("abort", () => {
            safeCloseSocket(socket);
          });
        });
      }),
    );

    const schedule =
      params.retrySchedule ?? defaultPurchaseSocketRetrySchedule;

    yield* Effect.retry(singleAttempt, {
      schedule,
      while: (err: unknown) =>
        typeof err === "object" &&
        err !== null &&
        "_tag" in err &&
        ((err as { _tag: string })._tag === "PurchaseWsTokenError" ||
          (err as { _tag: string })._tag === "PurchaseSocketConnectionError"),
    });
  });

export function usePurchaseSocket(
  purchaseId: () => string,
  onChange: (purchaseId: string) => void,
  options?: {
    retrySchedule?: Schedule.Schedule<unknown, unknown>;
    apiBaseUrl?: string;
    WebSocketClass?: typeof WebSocket;
  },
) {
  const [connected, setConnected] = createSignal(false);

  createEffect(
    () => purchaseId(),
    (id) => {
      if (!id) return;

      const program = purchaseSocketSessionEffect({
        purchaseId: id,
        onChange,
        onConnectionChange: setConnected,
        retrySchedule: options?.retrySchedule,
        apiBaseUrl: options?.apiBaseUrl,
        WebSocketClass: options?.WebSocketClass,
      });

      const fiber = Effect.runFork(program);

      return () => {
        Effect.runFork(Fiber.interrupt(fiber));
        setConnected(false);
      };
    },
  );

  return connected;
}
