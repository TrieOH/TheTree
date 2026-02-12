import { createClientOnlyFn } from "@tanstack/react-start";
import type { Schema, SchemaCRUD } from "../model/types";
import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { queryOptions } from "@tanstack/react-query";


/**
 * Creates a new schema draft on the server.
 * @param schemaData - The data for the new schema.
 * @returns A promise that resolves to the API response containing the newly created schema.
 */
export const createSchemaFn = createClientOnlyFn((schemaData: Omit<SchemaCRUD, "id">) => {
  const { project_id, ...dataToSend } = schemaData;

  return authFetcher<Schema>(`/projects/${project_id}/schemas`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({ ...dataToSend, schema_type: "context" }),
  });
});


/**
 * Fetches all schemas from the server.
 * @returns A promise that resolves to an array of Schema objects.
 */
export const getSchemasFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["schemas", string];
}) => {
  const [, projectId] = queryKey;

  try {
    return await tanstackQueryFetcher<Schema[]>(
      `/projects/${projectId}/schemas`
    );
  } catch {
    return [] as Schema[];
  }
});

export const schemasQueryOptions = (project_id: string) => {
  return queryOptions({
    queryKey: ['schemas', project_id],
    queryFn: getSchemasFn,
    enabled: !!project_id
  })
}