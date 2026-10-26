import type { JSX } from "@solidjs/web";
import { useRouter } from "@tanstack/solid-router";
import { Match, Switch, createMemo, type Accessor } from "solid-js";
import { classifyError } from "@/shared/lib/service-error";
import { serviceHealth } from "@/shared/lib/service-health";
import { QueryErrorPanel } from "./QueryErrorPanel";

/** The slice of a query result the boundary reads; `useQuery`'s result satisfies it. */
export interface QueryLike<T> {
  data: T | undefined;
  error: unknown;
  isFetching: boolean;
  errorUpdateCount: number;
  refetch: () => unknown;
}

export interface QueryBoundaryProps<T> {
  query: Accessor<QueryLike<T>>;
  /** What is being loaded, as it reads after "Não foi possível carregar": "os eventos". */
  what: string;
  /** Shaped like the content, so the swap does not move the page. */
  skeleton: JSX.Element;
  /** A specific reason for an error (e.g. a 404 for this particular page). */
  errorMessage?: (error: unknown) => string | undefined;
  /** Rendered in place of the retry panel when the answer is a definitive 404. */
  notFound?: JSX.Element;
  children: (data: Accessor<T>) => JSX.Element;
  class?: string;
}

function tryGetRouter() {
  try {
    return useRouter({ warn: false });
  } catch {
    return undefined;
  }
}

/**
 * Renders one query's four states without ever blanking or remounting content:
 *
 *  - **data** (fresh, stale, or stale with a failed background refresh): the
 *    content, and it stays mounted through refetches;
 *  - **loading** with no data yet: the skeleton, faded in after a beat so a
 *    cache-fast answer never flashes it;
 *  - **failed** with no data: a retry panel scoped to this region (or `notFound`
 *    for a definitive 404). It stays up
 *    while a retry runs; a refetch resets the query to "pending", which would
 *    otherwise bounce panel → skeleton → panel on every click.
 *
 * `<Loading>` can't do this job here: `useQuery` is observer-based and never
 * suspends, so a `<Loading>` around it never shows its fallback.
 */
export function QueryBoundary<T>(props: QueryBoundaryProps<T>) {
  const router = tryGetRouter();
  const hasData = createMemo(() => props.query().data !== undefined, {
    name: "QueryBoundary.hasData",
  });
  const hasFailed = createMemo(() => props.query().errorUpdateCount > 0, {
    name: "QueryBoundary.hasFailed",
  });
  // The error is cleared while a retry runs; keep showing the last one.
  const lastError = createMemo<unknown>(
    (previous) => props.query().error ?? previous,
    { name: "QueryBoundary.lastError" },
  );

  const isNotFound = createMemo(
    () => props.notFound !== undefined && classifyError(lastError()).kind === "not-found",
    { name: "QueryBoundary.isNotFound" },
  );

  const retry = () => {
    // An explicit retry shouldn't wait out a circuit's cool-down.
    // Also notifies query-client to refetch any other failed queries on the page.
    serviceHealth.requestRetry();
    void props.query().refetch();
    if (router) void router.invalidate();
  };

  return (
    <Switch>
      <Match when={hasData()}>{props.children(() => props.query().data as T)}</Match>
      <Match when={hasFailed() && isNotFound()}>{props.notFound}</Match>
      <Match when={hasFailed()}>
        <QueryErrorPanel
          error={lastError()}
          message={props.errorMessage?.(lastError())}
          what={props.what}
          retrying={props.query().isFetching}
          onRetry={retry}
          class={props.class}
        />
      </Match>
      <Match when={true}>
        <div class="skeleton-delay">{props.skeleton}</div>
      </Match>
    </Switch>
  );
}

type DataOf<Q> = Q extends Accessor<QueryLike<infer T>> ? T : never;

/**
 * Several queries as one for a {@link QueryBoundary}: ready when all of them
 * have data, failed when any has failed, retrying re-asks all of them.
 */
export function combineQueries<const Q extends readonly Accessor<QueryLike<unknown>>[]>(
  ...queries: Q
): Accessor<QueryLike<{ -readonly [K in keyof Q]: DataOf<Q[K]> }>> {
  const list = queries as unknown as Accessor<QueryLike<unknown>>[];
  return createMemo(
    () => {
      const results = list.map((query) => query());
      const ready = results.every((result) => result.data !== undefined);
      return {
        data: ready
          ? (results.map((result) => result.data) as { -readonly [K in keyof Q]: DataOf<Q[K]> })
          : undefined,
        error: results.find((result) => result.error)?.error ?? null,
        isFetching: results.some((result) => result.isFetching),
        errorUpdateCount: results.reduce((sum, result) => sum + result.errorUpdateCount, 0),
        refetch: () => list.forEach((query) => void query().refetch()),
      };
    },
    { name: "combineQueries" },
  );
}
