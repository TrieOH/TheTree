// Pages that share a payload read it through here: one request per path per
// 30 s, so moving between pages doesn't refetch. A failure isn't kept.
const FRESH_MS = 30_000;
const cache = new Map<string, { at: number; data: Promise<unknown> }>();

export function cachedJson<T>(path: string): Promise<T> {
  const hit = cache.get(path);
  if (hit && Date.now() - hit.at < FRESH_MS) return hit.data as Promise<T>;
  const data = fetch(path).then((response) => {
    if (!response.ok) throw new Error(`GET ${path}: ${response.status}`);
    return response.json() as Promise<T>;
  });
  cache.set(path, { at: Date.now(), data });
  data.catch(() => {
    if (cache.get(path)?.data === data) cache.delete(path);
  });
  return data;
}
