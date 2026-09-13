import { useQueryClient } from "@trieoh/front-core-solid";
import type { StoreStockItem } from "@trieoh/univents-api/schemas";
import { createEffect, createSignal, type Accessor } from "solid-js";
import { Data, Effect, Fiber, Schedule } from "effect";
import { cart } from "../model/cart";

export class InventoryStreamConnectionError extends Data.TaggedError(
  "InventoryStreamConnectionError",
)<{
  readonly message?: string;
  readonly cause?: unknown;
}> { }

export class InventoryStreamParseError extends Data.TaggedError(
  "InventoryStreamParseError",
)<{
  readonly rawData: string;
  readonly eventType: "snapshot" | "stock";
  readonly cause: unknown;
}> { }

export const isValidStoreStockItem = (val: unknown): val is StoreStockItem =>
  typeof val === "object" &&
  val !== null &&
  typeof (val as StoreStockItem).id === "string" &&
  ("stock" in val
    ? (val as StoreStockItem).stock === null ||
    typeof (val as StoreStockItem).stock === "number"
    : false);

export const parseSnapshotData = (
  rawData: string,
): Effect.Effect<StoreStockItem[], InventoryStreamParseError> =>
  Effect.try({
    try: () => {
      const parsed = JSON.parse(rawData);
      if (!Array.isArray(parsed)) return [];
      return parsed.filter(isValidStoreStockItem);
    },
    catch: (cause) =>
      new InventoryStreamParseError({
        rawData,
        eventType: "snapshot",
        cause,
      }),
  });

export const parseStockData = (
  rawData: string,
): Effect.Effect<StoreStockItem, InventoryStreamParseError> =>
  Effect.try({
    try: () => {
      const parsed = JSON.parse(rawData);
      if (!isValidStoreStockItem(parsed)) {
        throw new Error("Item de estoque inválido.");
      }
      return parsed;
    },
    catch: (cause) =>
      new InventoryStreamParseError({
        rawData,
        eventType: "stock",
        cause,
      }),
  });

export interface InventoryStreamSessionParams {
  readonly editionId: string;
  readonly onStockUpdate: (item: StoreStockItem) => void;
  readonly onConnectionChange?: (connected: boolean) => void;
  readonly retrySchedule?: Schedule.Schedule<unknown, unknown>;
  readonly apiBaseUrl?: string;
  readonly EventSourceClass?: typeof EventSource;
}

export const defaultInventoryStreamRetrySchedule = Schedule.exponential("500 millis").pipe(
  Schedule.jittered,
  Schedule.upTo({ times: 10 }),
);

const safeCloseEventSource = (stream: EventSource) => {
  try {
    stream.close();
  } catch {
    // Socket was already closed or disposed
  }
};

export const inventoryStreamSessionEffect = (
  params: InventoryStreamSessionParams,
): Effect.Effect<void, unknown> =>
  Effect.gen(function* () {
    const singleAttempt = Effect.scoped(
      Effect.gen(function* () {
        const baseUrl =
          params.apiBaseUrl ??
          import.meta.env.VITE_API_URL ??
          "http://localhost:8081";
        const url = new URL(`/editions/${params.editionId}/store/stream`, baseUrl);

        const ES = params.EventSourceClass ?? globalThis.EventSource;
        if (!ES) {
          return;
        }

        const stream = yield* Effect.acquireRelease(
          Effect.sync(() => new ES(url)),
          (es) =>
            Effect.sync(() => {
              safeCloseEventSource(es);
              params.onConnectionChange?.(false);
            }),
        );

        yield* Effect.callback<void, InventoryStreamConnectionError>((resume, signal) => {
          stream.onopen = () => {
            params.onConnectionChange?.(true);
          };

          const handleSnapshot = (event: MessageEvent<string>) => {
            Effect.runSync(
              parseSnapshotData(String(event.data)).pipe(
                Effect.tap((items) =>
                  Effect.sync(() => {
                    for (const item of items) {
                      params.onStockUpdate(item);
                    }
                  }),
                ),
                Effect.orElseSucceed(() => []),
              ),
            );
          };

          const handleStock = (event: MessageEvent<string>) => {
            Effect.runSync(
              parseStockData(String(event.data)).pipe(
                Effect.tap((item) =>
                  Effect.sync(() => {
                    params.onStockUpdate(item);
                  }),
                ),
                Effect.orElseSucceed(() => null),
              ),
            );
          };

          stream.addEventListener("snapshot", handleSnapshot as EventListener);
          stream.addEventListener("stock", handleStock as EventListener);

          stream.onerror = (err) => {
            params.onConnectionChange?.(false);
            resume(
              Effect.fail(
                new InventoryStreamConnectionError({
                  message: "Erro na conexão SSE do inventário.",
                  cause: err,
                }),
              ),
            );
          };

          signal.addEventListener("abort", () => {
            safeCloseEventSource(stream);
          });
        });
      }),
    );

    const schedule =
      params.retrySchedule ?? defaultInventoryStreamRetrySchedule;

    yield* Effect.retry(singleAttempt, {
      schedule,
      while: (err: unknown) =>
        typeof err === "object" &&
        err !== null &&
        "_tag" in err &&
        (err as { _tag: string })._tag === "InventoryStreamConnectionError",
    });
  });

export interface UseInventoryStreamOptions {
  readonly onConnectionChange?: (connected: boolean) => void;
  readonly retrySchedule?: Schedule.Schedule<unknown, unknown>;
  readonly apiBaseUrl?: string;
  readonly EventSourceClass?: typeof EventSource;
}

export function useInventoryStream(
  editionId: string | (() => string),
  initialStock: StoreStockItem[],
  options?: UseInventoryStreamOptions,
): Accessor<StoreStockItem[]> {
  const queryClient = useQueryClient();
  const [stockItems, setStockItems] = createSignal<StoreStockItem[]>(initialStock);

  const resolveId = typeof editionId === "function" ? editionId : () => editionId;

  createEffect(
    () => resolveId(),
    (id) => {
      if (!id) return;

      const apply = (item: StoreStockItem) => {
        cart.stock(id, item.id, item.stock);
        setStockItems((current) => {
          const exists = current.some((entry) => entry.id === item.id);
          return exists
            ? current.map((entry) => (entry.id === item.id ? item : entry))
            : [...current, item];
        });
        queryClient.setQueryData<StoreStockItem[]>(
          ["store", "stock", id],
          (current = []) => {
            const exists = current.some((entry) => entry.id === item.id);
            return exists
              ? current.map((entry) => (entry.id === item.id ? item : entry))
              : [...current, item];
          },
        );
      };

      const program = inventoryStreamSessionEffect({
        editionId: id,
        onStockUpdate: apply,
        onConnectionChange: options?.onConnectionChange,
        retrySchedule: options?.retrySchedule,
        apiBaseUrl: options?.apiBaseUrl,
        EventSourceClass: options?.EventSourceClass,
      });

      const fiber = Effect.runFork(program);

      return () => {
        Effect.runFork(Fiber.interrupt(fiber));
        options?.onConnectionChange?.(false);
      };
    },
  );

  return stockItems;
}
