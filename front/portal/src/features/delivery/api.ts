import { cachedJson } from "@/shared/lib/cached-json";
import type { DeliveryPayload } from "./model";

export const loadDelivery = () => cachedJson<DeliveryPayload>("/api/delivery");
