import { createServerFn } from "@tanstack/react-start";
import { paymentConnectSchema } from "../model";
import type { PaymentConnectResponseI } from "../model";
// import { paymentsAuthFetcher } from "@/shared/lib/api/fetch";
import type { ApiResponse } from "@soramux/node-auth-sdk";
import { env } from "@/env";


/**
 * Connect seller to Workspace on the server.
 * @param payData - The data for the new payment connection.
 * @returns A promise that resolves to the API response containing the newly created connection.
 */
export const connectEditionSellerToWorkspaceFn = createServerFn({ method: 'POST' })
  .inputValidator(paymentConnectSchema)
  .handler(async ({ data }) => {
    const { provider, workspace_name, ...payData } = data
    const url = `${env.VITE_PAY_API_URL}/workspaces/${workspace_name}/providers/${provider}/connect`
    console.log(url)
    const res = await fetch(
      url,
      {
        method: "POST",
        body: JSON.stringify(payData),
        headers: {
          "X-API-Key": env.PAY_SECRET_KEY,
          "Content-Type": "application/json"
        }
      }
    )

    const text = await res.text()
    console.log("Resposta: ", text)
    if (!res.ok) {
      console.error(`[connectFn] HTTP ${res.status}:`, text)
      throw new Error(`Erro na API: ${res.status} - ${text.slice(0, 200)}`)
    }

    try {
      const result = JSON.parse(text) as ApiResponse<PaymentConnectResponseI>
      result.success = result.code === 200
      return result
    } catch (e) {
      console.error('[connectFn] Resposta não é JSON:', text)
      throw new Error(`Resposta inválida da API: ${text.slice(0, 200)}`)
    }
  })
// export const connectEditionSellerToWorkspaceFn = createServerFn({ method: 'POST' })
//   .inputValidator(paymentConnectSchema)
//   .handler(async ({ data }) => {
//     const { provider, workspace_name, ...payData } = data
//     const res = await fetch(`${env.VITE_PAY_API_URL}/workspaces/${workspace_name}/${provider}/connect`, {
//       method: "POST",
//       body: JSON.stringify(payData),
//       headers: {
//         "X-API-Key": env.PAY_SECRET_KEY,
//         "Content-Type": "application/json"
//       }
//     })
//     const resJson = await res.json() as ApiResponse<PaymentConnectResponseI>
//     return resJson
//   })
