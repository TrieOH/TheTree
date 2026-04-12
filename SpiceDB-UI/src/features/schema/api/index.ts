import { createServerFn } from "@tanstack/react-start";
import { z } from "zod";
import { spicedb } from "@trieoh/node-perm-sdk";
import { queryOptions } from "@tanstack/react-query";

const writeSchemaInput = z.object({
  schema: z.string().min(1),
});


export const readSchema = createServerFn({
  method: "GET",
}).handler(async () => {
  const response = await spicedb.schema.read()
  if (response.success) return response.data
  else return { schemaText: "" }
})

export const schemaQueryOptions = queryOptions({
  queryKey: ["schema"],
  queryFn: () => readSchema()
});

export const writeSchema = createServerFn({
  method: "POST",
})
  .inputValidator(writeSchemaInput)
  .handler(async ({ data }) => {
    const response = await spicedb.schema.write({ schema: data.schema })
    if (!response.success) {
      return {
        success: false,
        message: response.message,
      }
    }
    return response
  })
