import { client } from "./utils/_client";
import type {
  ConnectRequestI,
  ConnectResponseI,
  ProviderCredentialResponseI
} from "./utils/types";

export const workspaceService = {
  /**
   * Connect a seller account to a workspace
   * @param name Workspace name
   * @param provider Provider name, eg. mercadopago
   * @param body - { final_redirect_url: string, provider_redirect_url: string }
   * @returns The ApiResponse containing a ConnectResponseI
   */
  connectProvider: (
    name: string,
    provider: string,
    body: ConnectRequestI,
  ) =>
    client.post<ConnectResponseI>(`/workspaces/${name}/providers/${provider}/connect`, body),

  /**
   * Connect a seller account to a workspace
   * @param name Workspace name
   * @param credential_id Credential ID
   * @returns The ApiResponse containg a ProviderCredentialResponseI
   */
  disconnectProvider: (
    name: string,
    credential_id: string,
  ) =>
    client.delete<ProviderCredentialResponseI>(
      `/workspaces/${name}/providers/${credential_id}/disconnect`
    ),
};
