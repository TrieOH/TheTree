import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { SchemaVersion, VersionDraft } from "../model/types";

export const createSchemaVersionDraftFn = createClientOnlyFn((versionData: VersionDraft) => {
  const { project_id, schema_id } = versionData
  return authFetcher<SchemaVersion>(`/projects/${project_id}/schemas/${schema_id}/versions/draft`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
  });
});


export const getLatestSchemaVersionFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["latestSchemaVersion", string, string];
}) => {
  const [, projectId, schemaId] = queryKey;

  try {
    return await tanstackQueryFetcher<SchemaVersion>(
      `/projects/${projectId}/schemas/${schemaId}/versions/latest`
    );
  } catch (error) {
    console.error("Error fetching latest schema version:", error);
    throw error;
  }
});

export const latestSchemaVersionQueryOptions = (project_id: string, schema_id: string) => {
  return queryOptions({
    queryKey: ['latestSchemaVersion', project_id, schema_id],
    queryFn: getLatestSchemaVersionFn,
    enabled: !!project_id && !!schema_id
  })
}

export const getCurrentSchemaVersionFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["currentSchemaVersion", string, string];
}) => {
  const [, projectId, schemaId] = queryKey;

  try {
    return await tanstackQueryFetcher<SchemaVersion>(
      `/projects/${projectId}/schemas/${schemaId}/versions/current`
    );
  } catch (error) {
    console.error("Error fetching current schema version:", error);
    throw error;
  }
});

export const currentSchemaVersionQueryOptions = (project_id: string, schema_id: string) => {
  return queryOptions({
    queryKey: ['currentSchemaVersion', project_id, schema_id],
    queryFn: getCurrentSchemaVersionFn,
    enabled: !!project_id && !!schema_id
  })
}