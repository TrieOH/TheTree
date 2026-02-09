import { createClientOnlyFn } from "@tanstack/react-start";
import type { Schema, SchemaCRUD } from "../model/types";
import { authFetcher } from "@/shared/lib/api/fetch";


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
    body: JSON.stringify({ ...dataToSend, schema_type: "core" }),
  });
});

