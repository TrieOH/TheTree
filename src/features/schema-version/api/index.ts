import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { DetailedSchemaVersion, PartialVersionField, SchemaFieldOption, SchemaFieldRequiredRule, SchemaFieldVisibilityRule, SchemaVersion, SchemaVersionFields, VersionDraft } from "../model/types";

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

// Fields

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



export const createSchemaVersionFieldFn = createClientOnlyFn((fieldsData: SchemaVersionFields) => {
  const { project_id, schema_id, version, fields } = fieldsData
  return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({fields})
  });
});

export const deleteSchemaVersionFieldFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id }: { project_id: string; schema_id: string; version: number; field_id: string }) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}`, {
      method: "DELETE",
    });
  }
);

export const updateSchemaVersionFieldFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, fieldData }: { project_id: string; schema_id: string; version: number; field_id: string; fieldData: PartialVersionField }) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(fieldData),
    });
  }
);

export const setSchemaFieldOptionsFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, options }: { project_id: string; schema_id: string; version: number; field_id: string; options: SchemaFieldOption[] }) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/options`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(options),
    });
  }
);

export const deleteSchemaFieldOptionFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, option_id }: { project_id: string; schema_id: string; version: number; field_id: string; option_id: string }) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/options/${option_id}`, {
      method: "DELETE",
    });
  }
);

export const setSchemaFieldVisibilityRulesFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, rules }: { project_id: string; schema_id: string; version: number; field_id: string; rules: SchemaFieldVisibilityRule[] }) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/visibility-rules`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(rules),
    });
  }
);

export const deleteSchemaFieldVisibilityRuleFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, rule_id }: { project_id: string; schema_id: string; version: number; field_id: string; rule_id: string }) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/visibility-rules/${rule_id}`, {
      method: "DELETE",
    });
  }
);

export const setSchemaFieldRequiredRulesFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, rules }: { project_id: string; schema_id: string; version: number; field_id: string; rules: SchemaFieldRequiredRule[] }) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/required-rules`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(rules),
    });
  }
);

export const deleteSchemaFieldRequiredRuleFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, rule_id }: { project_id: string; schema_id: string; version: number; field_id: string; rule_id: string }) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/required-rules/${rule_id}`, {
      method: "DELETE",
    });
  }
);