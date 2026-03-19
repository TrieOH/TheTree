import { createServerFn } from "@tanstack/react-start";
import { workspaceService } from "@soramux/node-payments-sdk"
import { paymentConnectSchema, paymentDisconnectSchema } from "../model";

/**
 * Connect seller to Workspace on the server.
 * @param payData - The data for the new payment connection.
 * @returns A promise that resolves to the API response containing the newly created connection.
 */
export const connectEditionSellerToWorkspaceFn = createServerFn({ method: 'POST' })
  .inputValidator(paymentConnectSchema)
  .handler(async ({ data }) => {
    const { provider, workspace_name, ...payData } = data
    return workspaceService.connectProvider(
      workspace_name,
      provider,
      payData
    );
  })

/**
 * Disconnect seller from Workspace on the server.
 * @param payData - The data for perform the payment disconnect
 * @returns A promise that resolves to the API response containing the removed connection.
 */
export const disconnectEditionSellerToWorkspaceFn = createServerFn({ method: 'POST' })
  .inputValidator(paymentDisconnectSchema)
  .handler(async ({ data }) => {
    const { workspace_name, credential_id } = data
    return workspaceService.disconnectProvider(
      workspace_name,
      credential_id,
    );
  })
