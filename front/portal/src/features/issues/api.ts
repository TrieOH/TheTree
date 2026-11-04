import { cachedJson } from "@/shared/lib/cached-json";
import type { IssuesPayload } from "./model";

export const loadIssues = () => cachedJson<IssuesPayload>("/api/issues");
