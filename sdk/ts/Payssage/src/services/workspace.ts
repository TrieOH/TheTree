import { client } from "./utils/_client";
import type { ConnectRequestI } from "./utils/types";

export const workspaceService = {
  /**
   * Connect a seller account to a workspace
   * @param name Workspace name
   * @param provider Provider name, eg. mercadopago
   * @param body - { final_redirect_url: string, provider_redirect_url: string }
   * @returns The ApiResponse null
   */
  connectProvider: (
    name: string,
    provider: string,
    body: ConnectRequestI,
  ) =>
    client.post<null>(`/workspaces/${name}/providers/${provider}/connect`, body),

  /**
   * Connect a seller account to a workspace
   * @param name Workspace name
   * @param credential_id Credential ID
   * @returns 
   */
  disconnectProvider: (
    name: string,
    credential_id: string,
  ) =>
    client.delete<null>(`/workspaces/${name}/providers/${credential_id}/disconnect`),
};
