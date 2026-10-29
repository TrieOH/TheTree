import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { itemsOf } from "@trieoh/api-client";
import { authFetcher } from "#/shared/lib/api/fetch";
import type {
  CreateIntentRequest,
  CreateTestModeIntentRequest,
  Intent,
} from "../model";

export const listAllIntentsQueryOptions = () =>
  queryOptions({
    queryKey: ["intents", "personal"],
    queryFn: async () => {
      const response = await authFetcher.get<unknown>("/intents");
      if (!response.success) throw response;
      return itemsOf<Intent>(response.data);
    },
  });

export const listAllByWalletIntentsQueryOptions = (walletId: string) =>
  queryOptions({
    queryKey: ["intents", walletId],
    queryFn: async () => {
      const response = await authFetcher.get<unknown>(
        `/wallets/${walletId}/intents`,
      );
      if (!response.success) throw response;
      return itemsOf<Intent>(response.data);
    },
  });

export const listAllByOrgIntentsQueryOptions = (organizationId: string) =>
  queryOptions({
    queryKey: ["intents", organizationId],
    queryFn: async () => {
      const response = await authFetcher.get<unknown>(
        `/organizations/${organizationId}/intents`,
      );
      if (!response.success) throw response;
      return itemsOf<Intent>(response.data);
    },
  });

export const createTestModeWalletIntentFn = createClientOnlyFn(
  (payload: CreateTestModeIntentRequest) =>
    authFetcher.post<Intent>("/testmode/intents/create", payload),
);

export const createWalletIntentFn = createClientOnlyFn(
  (walletId: string, payload: CreateIntentRequest) =>
    authFetcher.post<Intent>(`/wallets/${walletId}/intents`, payload),
);

export const cancelIntentFn = createClientOnlyFn((intentId: string) =>
  authFetcher.post<Intent>(`/intents/${intentId}/cancel`),
);
