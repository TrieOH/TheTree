import type { IssuesPayload } from "./model";

// Every page that shows issues reads the same payload; keep it for a short
// while so moving between Home, Issues and a service page is one request.
const FRESH_MS = 30_000;
let cached: { at: number; data: Promise<IssuesPayload> } | undefined;

export function loadIssues(): Promise<IssuesPayload> {
  if (cached && Date.now() - cached.at < FRESH_MS) return cached.data;
  const data = fetch("/api/issues").then((response) => {
    if (!response.ok) throw new Error(`GET /api/issues: ${response.status}`);
    return response.json() as Promise<IssuesPayload>;
  });
  cached = { at: Date.now(), data };
  // A failure shouldn't stick for the whole window.
  data.catch(() => {
    if (cached?.data === data) cached = undefined;
  });
  return data;
}
