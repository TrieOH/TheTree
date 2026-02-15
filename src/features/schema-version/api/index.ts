import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { DetailedSchemaVersion, SchemaVersion, SchemaVersionFields, VersionDraft } from "../model/types";

export const createSchemaVersionDraftFn = createClientOnlyFn((versionData: VersionDraft) => {
  const { project_id, schema_id } = versionData
  return authFetcher<SchemaVersion>(`/projects/${project_id}/schemas/${schema_id}/versions/draft`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
  });
});


export const publishSchemaVersionFn = createClientOnlyFn((fieldsData: VersionDraft) => {
  const { project_id, schema_id } = fieldsData
  return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/versions/publish`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  });
});

export const publishSchemaVersionFieldFn = createClientOnlyFn((fieldsData: SchemaVersionFields) => {
  const { project_id, schema_id, version, fields } = fieldsData
  return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({fields})
  });
});


export const getLatestSchemaVersionFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["latestSchemaVersion", string, string];
}) => {
  const [, projectId, schemaId] = queryKey;
  return await tanstackQueryFetcher<SchemaVersion>(
    `/projects/${projectId}/schemas/${schemaId}/versions/latest`
  );
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
  return await tanstackQueryFetcher<SchemaVersion>(
    `/projects/${projectId}/schemas/${schemaId}/versions/current`
  );
});

export const currentSchemaVersionQueryOptions = (project_id: string, schema_id: string) => {
  return queryOptions({
    queryKey: ['currentSchemaVersion', project_id, schema_id],
    queryFn: getCurrentSchemaVersionFn,
    enabled: !!project_id && !!schema_id
  })
}

export const getSchemaVersionByIdFn = createClientOnlyFn(async ({
  queryKey,
}: {
  queryKey: ["schemaVersionById", string, string, number];
}) => {
  const [, projectId, schemaId, version] = queryKey;
  return await tanstackQueryFetcher<DetailedSchemaVersion>(
    `/projects/${projectId}/schemas/${schemaId}/v${version}`
  );
});

export const schemaVersionByIdQueryOptions = (project_id: string, schema_id: string, version: number) => {
  return queryOptions({
    queryKey: ['schemaVersionById', project_id, schema_id, version],
    queryFn: getSchemaVersionByIdFn,
    enabled: Boolean(project_id && schema_id && version !== undefined)
  })
}