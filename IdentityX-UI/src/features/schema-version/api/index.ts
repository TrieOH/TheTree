import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { 
  ProjectFieldDefinitionResultI, 
  SchemaVersionDraft, 
  SchemaVersionResultI,
  SchemaFieldCreateRequestI,
  FieldDefinitionResultI,
  OptionFieldCreateRequestI,
  RuleFieldCreateRequestI
} from "../model/types";

export const createSchemaVersionDraftFn = createClientOnlyFn((versionData: SchemaVersionDraft) => {
  const { project_id, schema_id } = versionData
  return authFetcher<SchemaVersionResultI>(
    `/projects/${project_id}/schemas/${schema_id}/versions/draft`, 
    {
      method: "POST",
      headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    }
  );
});


export const publishSchemaVersionFn = createClientOnlyFn((fieldsData: SchemaVersionDraft) => {
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
  return await tanstackQueryFetcher<SchemaVersionResultI>(
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
  return await tanstackQueryFetcher<SchemaVersionResultI>(
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
  return await tanstackQueryFetcher<ProjectFieldDefinitionResultI>(
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

interface CreateSchemaVersionFieldRequestI {
  project_id: string;
  schema_id: string;
  version: number;
  fields: SchemaFieldCreateRequestI[];
}

export const createSchemaVersionFieldFn = createClientOnlyFn((data: CreateSchemaVersionFieldRequestI) => {
  const { project_id, schema_id, version, fields } = data
  return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({fields})
  });
});

interface SetSchemaVersionFieldsRequestI {
  fields: FieldDefinitionResultI[]
  project_id: string;
  schema_id: string;
  version: number;
}

export const setSchemaVersionFieldsFn = createClientOnlyFn((data: SetSchemaVersionFieldsRequestI) => {
  const { project_id, schema_id, version, fields } = data
  return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
    body: JSON.stringify({fields})
  });
});

interface DeleteSchemaVersionFieldRequestI {
  project_id: string;
  schema_id: string;
  version: number;
  field_id: string;
}

export const deleteSchemaVersionFieldFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id }: DeleteSchemaVersionFieldRequestI) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}`, {
      method: "DELETE",
    });
  }
);

interface SetSchemaFieldOptionsRequestI {
  project_id: string;
  schema_id: string;
  version: number;
  field_id: string;
  options: OptionFieldCreateRequestI[];
}

export const setSchemaFieldOptionsFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, options }: SetSchemaFieldOptionsRequestI) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/options`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({options}),
    });
  }
);

interface DeleteSchemaFieldOptionRequestI {
  project_id: string;
  schema_id: string;
  version: number;
  field_id: string;
  option_id: string;
}

export const deleteSchemaFieldOptionFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, option_id }: DeleteSchemaFieldOptionRequestI) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/options/${option_id}`, {
      method: "DELETE",
    });
  }
);

interface SetSchemaFieldRulesRequestI {
  project_id: string;
  schema_id: string;
  version: number;
  field_id: string;
  rules: RuleFieldCreateRequestI[];
}

export const setSchemaFieldVisibilityRulesFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, rules }: SetSchemaFieldRulesRequestI) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/visibility-rules`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({visibility_rules: rules}),
    });
  }
);

// export const deleteSchemaFieldVisibilityRuleFn = createClientOnlyFn(
//   ({ project_id, schema_id, version, field_id, rule_id }: { project_id: string; schema_id: string; version: number; field_id: string; rule_id: string }) => {
//     return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/visibility-rules/${rule_id}`, {
//       method: "DELETE",
//     });
//   }
// );

export const setSchemaFieldRequiredRulesFn = createClientOnlyFn(
  ({ project_id, schema_id, version, field_id, rules }: SetSchemaFieldRulesRequestI) => {
    return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/required-rules`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({required_rules: rules}),
    });
  }
);

// export const deleteSchemaFieldRequiredRuleFn = createClientOnlyFn(
//   ({ project_id, schema_id, version, field_id, rule_id }: { project_id: string; schema_id: string; version: number; field_id: string; rule_id: string }) => {
//     return authFetcher<string>(`/projects/${project_id}/schemas/${schema_id}/v${version}/fields/${field_id}/required-rules/${rule_id}`, {
//       method: "DELETE",
//     });
//   }
// );