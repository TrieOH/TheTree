import { createServerFn } from "@tanstack/react-start";
import { z } from "zod";
import { spicedb } from "@trieoh/node-perm-sdk";
import { queryOptions } from "@tanstack/react-query";


export const readSchema = createServerFn({
  method: "GET",
})
  .inputValidator((envId: string) => envId)
  .handler(async ({ data: envId }) => {
    const response = await spicedb.schema(envId).read()
    if (response.success) return response.data
    else if (response.code === 5) return { schemaText: "", readAt: { token: "" } }
    throw new Error(response.message || "Erro ao ler o schema")
  })

export const schemaQueryOptions = (envId: string) => queryOptions({
  queryKey: ["schema", envId],
  queryFn: () => readSchema({ data: envId }),
  staleTime: 0,
});

export const writeSchema = createServerFn({
  method: "POST",
})
  .inputValidator(z.object({
    envId: z.string(),
    schema: z.string().min(1),
  }))
  .handler(async ({ data }) => {
    const response = await spicedb.schema(data.envId).write({ schema: data.schema })
    if (!response.success) {
      return {
        success: false,
        message: response.message,
      }
    }
    return response
  })
