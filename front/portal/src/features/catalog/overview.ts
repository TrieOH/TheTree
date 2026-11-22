import type { AttentionItem } from "@/shared/ui/attention";
import { needsAttention as deliveryAttention } from "@/features/delivery/attention";
import type { DeliveryPayload } from "@/features/delivery/model";
import { needsAttention as metricsAttention } from "@/features/metrics/attention";
import type { ColdStartDay, MetricsPayload } from "@/features/metrics/model";
import type { Service } from "./services";

// What a service's own page pulls out of the Metrics and Delivery payloads.

/** Which page explains an item on a service page. */
export type OverviewSource = "metrics" | "delivery";

/** The Metrics and Delivery items that are about this service, its front or its Workers. */
export function serviceAttention(
  service: Service,
  metrics: MetricsPayload | null,
  delivery: DeliveryPayload | null,
  now = new Date(),
): AttentionItem<OverviewSource>[] {
  const id = service.id;
  const workers = [service.analytics.apiWorker, service.analytics.frontWorker];
  const mine = (itemId: string) =>
    itemId.endsWith(`-${id}`) ||
    itemId.endsWith(`-${id}-ui`) ||
    itemId.startsWith(`slow-${id}-`) ||
    itemId.startsWith(`tag-${id}/`) ||
    itemId.startsWith(`tag-${id}-ui/`) ||
    workers.some((script) => itemId === `worker-${script}`);
  const items: AttentionItem<OverviewSource>[] = [
    ...(metrics?.syncedAt ? metricsAttention(metrics, now) : []).filter((i) => mine(i.id)).map((i) => ({ ...i, tab: "metrics" as const })),
    ...(delivery ? deliveryAttention(delivery, now) : []).filter((i) => mine(i.id)).map((i) => ({ ...i, tab: "delivery" as const })),
  ];
  return items.sort((a, b) => (a.tone === b.tone ? 0 : a.tone === "destructive" ? -1 : 1));
}

export interface ColdStartBar {
  /** YYYY-MM-DD */
  day: string;
  count: number;
  failed: number;
}

/** One bar per UTC day over the last `days` days, oldest first; days without starts are zero. */
export function coldStartDays(rows: readonly ColdStartDay[], service: string, days = 30, now = new Date()): ColdStartBar[] {
  return Array.from({ length: days }, (_, i) => {
    const date = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate() - (days - 1 - i)));
    const day = date.toISOString().slice(0, 10);
    const row = rows.find((r) => r.service === service && r.day === day);
    return { day, count: row?.count ?? 0, failed: row?.failed ?? 0 };
  });
}
