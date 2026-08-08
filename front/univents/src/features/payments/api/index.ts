import { createClientOnlyFn } from "@tanstack/react-start";
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
import type { PaymentProviderI } from "../model";

export const connectEventSellerFn = createClientOnlyFn(
  (eventId: string, provider: PaymentProviderI) =>
    connectEventPayments(eventId, { provider }).then(
      orvalData<ConnectEventPaymentsResult>,
    ),
);

export const completeEventSellerFn = createClientOnlyFn(
  (eventId: string, data: CompleteEventPaymentsRequest) =>
    completeEventPayments(eventId, data).then(orvalData<Event>),
);

export const disconnectEventSellerFn = createClientOnlyFn((eventId: string) =>
  disconnectEventPayments(eventId).then(orvalData<null>),
);
