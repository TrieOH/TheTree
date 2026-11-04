import { cachedJson } from "@/shared/lib/cached-json";
import type { MetricsPayload } from "./model";

export const loadMetrics = () => cachedJson<MetricsPayload>("/api/metrics");
