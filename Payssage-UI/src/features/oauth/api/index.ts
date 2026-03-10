import { createClientOnlyFn } from "@tanstack/react-start";
import type { OauthCallbackResponseI, OauthSetupI, OauthSetupResponseI } from "../model";
import { authFetcher } from "#/shared/lib/api/fetch";
import { env } from "#/env";


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
  return authFetcher.post<OauthSetupResponseI>(
    `/workspaces/${workspaceName}/providers/${provider}/setup`,
    {
      ...oauthData,
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

// /**
//  * Revoke an API key for the specified workspace on the server.
//  * @param apiKeyId - The ID of the API key to revoke.
//  * @param workspaceName - The workspace name
//  * @returns A promise that resolves to the API response(void).
//  */
// export const removeMarketplaceOnWorkspaceFn = createClientOnlyFn((
//   apiKeyId: string,
//   workspaceName: string
// ) => {
//   return authFetcher.delete<void>(
//     `/workspaces/${workspaceName}/keys/${apiKeyId}`
//   );
// });