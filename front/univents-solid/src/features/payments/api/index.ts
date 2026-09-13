import { orvalData } from "@trieoh/api-client";
import {
  completeEventPayments,
  connectEventPayments,
  disconnectEventPayments,
} from "@trieoh/univents-api";
import type {
  CompleteEventPaymentsRequest,
  ConnectEventPaymentsResult,
  Event,
} from "@trieoh/univents-api/schemas";

export type PaymentProviderI = "mercadopago";

export const connectEventSellerFn = (
  eventId: string,
  provider: PaymentProviderI,
) =>
  connectEventPayments(eventId, { provider }).then(
    orvalData<ConnectEventPaymentsResult>,
  );

export const completeEventSellerFn = (
  eventId: string,
  data: CompleteEventPaymentsRequest,
) => completeEventPayments(eventId, data).then(orvalData<Event>);

export const disconnectEventSellerFn = (eventId: string) =>
  disconnectEventPayments(eventId).then(orvalData<null>);
