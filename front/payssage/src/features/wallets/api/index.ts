import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import type {
  WalletBindCollectorI,
  WalletCreateI,
  WalletI,
  WalletSetFeeBpsI,
  WalletSetSandboxI,
} from "../model";

/**
 * Creates a new wallet on the server.
 * @param walletData - The data for the new wallet (orgId is optional).
 * @returns A promise that resolves to the API response containing the newly created wallet.
 */
export const createWalletFn = createClientOnlyFn(
  (walletData: WalletCreateI) => {
    const { organization_id, ...rest } = walletData;
    const payload = {
      ...rest,
      ...(organization_id !== undefined && { organization_id }),
    };
    return authFetcher.post<WalletI>("/wallets", payload);
  },
);

/**
 * Sets the fee in basis points (BPS) for a specific wallet on the server.
 * @param walletId - The ID of the wallet to update.
 * @param data - The data containing the new fee in BPS and optional organization ID.
 * @returns A promise that resolves to the API response containing the updated wallet.
 */
export const setWalletFeeBPSFn = createClientOnlyFn(
  (walletId: string, data: WalletSetFeeBpsI) => {
    const { organization_id, fee_bps } = data;
    const payload = {
      fee_bps,
      ...(organization_id !== undefined && { organization_id }),
    };
    return authFetcher.patch(`/wallets/${walletId}/fee`, payload);
  },
);

/**
 * Sets the sandbox status for a specific wallet on the server.
 * @param walletId - The ID of the wallet to update.
 * @param data - The data containing the new sandbox status and optional organization ID.
 * @returns A promise that resolves to the API response.
 */
export const setWalletSandboxFn = createClientOnlyFn(
  (walletId: string, data: WalletSetSandboxI) => {
    const { organization_id, sandbox } = data;
    const payload = {
      sandbox,
      ...(organization_id !== undefined && { organization_id }),
    };
    return authFetcher.patch(`/wallets/${walletId}/sandbox`, payload);
  },
);

/**
 * Fetches a wallet by its ID from the server.
 * @param walletId - The ID of the wallet to fetch.
 * @returns A promise that resolves to the API response containing the wallet.
 */
export const getWalletByIdFn = createClientOnlyFn((walletId: string) => {
  return tanstackQueryFetcher<WalletI>(`/wallets/${walletId}`);
});

export const walletByIdQueryOptions = (walletId: string) =>
  queryOptions({
    queryKey: ["wallets", walletId],
    queryFn: () => getWalletByIdFn(walletId),
  });

export const bindCollectorToWalletFn = createClientOnlyFn(
  (walletId: string, payload: WalletBindCollectorI) =>
    authFetcher.post<void>(`/wallets/${walletId}/collector`, payload),
);

export const unbindCollectorFromWalletFn = createClientOnlyFn(
  (walletId: string) =>
    authFetcher.delete<void>(`/wallets/${walletId}/collector`),
);

/**
 * Fetches all wallets from the server.
 * @param orgId - The organization ID to filter wallets by (optional).
 * @returns A promise that resolves to an array of WalletI objects.
 */
export const getWalletsFn = createClientOnlyFn(async (orgId?: string) => {
  const wallets = await tanstackQueryFetcher<WalletI[]>(
    orgId ? `/organizations/${orgId}/wallets` : "/wallets",
  );
  return Array.isArray(wallets) ? wallets : [];
});

/**
 * Query options for fetching wallets, compatible with React Query's useQuery hook.
 * @param orgId - The organization ID to filter wallets by (optional).
 * @returns An object containing the query key and query function for fetching wallets.
 */
export const allWalletsQueryOptions = (orgId?: string) => {
  return queryOptions({
    queryKey: ["organizations", orgId, "wallets"],
    queryFn: () => getWalletsFn(orgId),
  });
};
