import type { QueryClient } from "@trieoh/front-core-solid";

/**
 * How long a route transition may wait for data before rendering anyway. Long
 * enough that a healthy API answers inside it (so the page paints with content
 * and no skeleton flash), short enough that a slow one never holds navigation.
 */
export const LOADER_WAIT_MS = 250;

type Options = Parameters<QueryClient["prefetchQuery"]>[0];

const delay = (ms: number) =>
  new Promise<void>((resolve) => setTimeout(resolve, ms));

/**
 * Loader-side data fetching that never blocks navigation for long and never
 * fails the route.
 *
 *  - cached data (even stale) is used as is, and revalidated in the background;
 *  - otherwise the request is started and given {@link LOADER_WAIT_MS} to land.
 *    Whatever happens after that, the page renders and reads the query itself.
 *
 * A loader that `await`s `prefetchQuery` instead holds the previous page on
 * screen for as long as the API takes, then flashes the pending spinner.
 */
export function warmQueries(
  client: QueryClient,
  queries: readonly Options[],
  waitMs: number = LOADER_WAIT_MS,
): Promise<void> {
  const inFlight: Promise<unknown>[] = [];

  for (const options of queries) {
    // `prefetchQuery` never rejects, and resolves at once for fresh data.
    const request = client.prefetchQuery(options);
    if (client.getQueryData(options.queryKey) === undefined) inFlight.push(request);
  }

  if (inFlight.length === 0) return Promise.resolve();
  return Promise.race([Promise.all(inFlight).then(() => undefined), delay(waitMs)]);
}

/**
 * The same bound for loaders that chain requests (a list, then the details of
 * each item), where the chain can't be expressed as independent queries.
 * The work keeps running after the bound; it only stops the navigation waiting.
 */
export function boundLoader(
  work: Promise<unknown>,
  waitMs: number = LOADER_WAIT_MS,
): Promise<void> {
  return Promise.race([
    work.then(() => undefined, () => undefined),
    delay(waitMs),
  ]);
}

export const warmQuery = (
  client: QueryClient,
  options: Options,
  waitMs?: number,
) => warmQueries(client, [options], waitMs);

/**
 * Reads a query for a guard: the cached value if there is one (refreshed in
 * the background), otherwise waits up to `waitMs` for the network. Resolves to
 * `undefined` when the answer is unknown — never rejects, so a guard can decide
 * to fail open instead of taking the app down with the API.
 */
export async function peekQuery<T>(
  client: QueryClient,
  options: Options & { queryKey: readonly unknown[] },
  waitMs: number = LOADER_WAIT_MS,
): Promise<T | undefined> {
  const cached = client.getQueryData<T>(options.queryKey);
  if (cached !== undefined) {
    void client.prefetchQuery(options);
    return cached;
  }

  await Promise.race([client.prefetchQuery(options), delay(waitMs)]);
  return client.getQueryData<T>(options.queryKey);
}
