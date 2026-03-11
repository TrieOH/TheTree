import { createClientOnlyFn } from "@tanstack/react-start";
import type { OauthCallbackResponseI, OauthSetupI, OauthSetupResponseI, OauthWorkspaceMarketplaceConfigI } from "../model";
import { authFetcher, tanstackQueryFetcher } from "#/shared/lib/api/fetch";
import { env } from "#/env";
import { percentageToBps } from "#/shared/lib/utils";
import { queryOptions } from "@tanstack/react-query";


/**
 * Set up OAuth for the specified workspace on the server.
 * @param oauthData - The OAuth setup data.
 * @param workspaceName - The workspace name
 * @returns A promise that resolves to the API response containing the OAuth setup information.
 */
export const setupOauthOnWorkspaceFn = createClientOnlyFn((
  oauthData: OauthSetupI,
  workspaceName: string,
  provider: string
) => {
  const { fee_percent, ...rest } = oauthData;
  return authFetcher.post<OauthSetupResponseI>(
    `/workspaces/${workspaceName}/providers/${provider}/setup`,
    {
      ...rest,
      fee_bps: percentageToBps(fee_percent),
      is_marketplace: true,
      final_redirect_url: env.VITE_MERCADO_PAGO_CALLBACK_URL
    }
  );
});

/**
 * Fetches the OAuth callback response for the specified provider.
 * @param code - The authorization code.
 * @param state - The state parameter.
 * @param provider - The OAuth provider.
 * @returns A promise that resolves to the OAuth callback response.
 */
export const getProviderCallbackFn = createClientOnlyFn((
  code: string,
  state: string,
  provider: string
) => {
  return authFetcher.get<OauthCallbackResponseI>(
    `/oauth/${provider}/callback?code=${code}&state=${state}&redirect_uri=${env.VITE_MERCADO_PAGO_CALLBACK_URL}`
  );
});

/**
 * Fetches all marketplace configs for the specified workspace from the server.
 * @param workspaceName - The workspace name
 * @returns A promise that resolves to an array of marketplace config objects.
 */
export const getAllWorkspaceMarketplaceConfigsFn = createClientOnlyFn(async (workspaceName: string) => {
  try {
    return await tanstackQueryFetcher<OauthWorkspaceMarketplaceConfigI[]>(
      `/workspaces/${workspaceName}/marketplaces`
    );
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all marketplace configs for a specific workspace, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all marketplace configs.
 */
export const allWorkspaceMarketplaceConfigsQueryOptions = (workspaceName: string) => {
  return queryOptions({
    queryKey: ['workspaces', workspaceName, "marketplace-configs"],
    queryFn: () => getAllWorkspaceMarketplaceConfigsFn(workspaceName),
  })
}

/**
 * Update OAuth for the specified workspace on the server.
 * @param oauthData - The OAuth setup data.
 * @param workspaceName - The workspace name
 * @returns A promise that resolves to the API response containing the OAuth marketplace config.
 */
export const updateOauthOnWorkspaceFn = createClientOnlyFn((
  oauthData: OauthSetupI,
  workspaceName: string,
  credential_id: string
) => {
  const { fee_percent } = oauthData;
  return authFetcher.put<OauthWorkspaceMarketplaceConfigI>(
    `/workspaces/${workspaceName}/marketplaces`,
    {
      fee_bps: percentageToBps(fee_percent),
      credential_id,
    }
  );
});


/**
 * Revoke an OAuth config for the specified workspace on the server.
 * @param credential_id - The ID of the credential to revoke.
 * @param workspaceName - The workspace name
 * @returns A promise that resolves to the API response(void).
 */
export const removeMarketplaceConfigFromWorkspaceFn = createClientOnlyFn((
  credential_id: string,
  workspaceName: string
) => {
  return authFetcher.delete<void>(
    `/workspaces/${workspaceName}/marketplaces/${credential_id}`
  );
});