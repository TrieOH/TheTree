import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  getProjectProfileSchema,
  upsertProjectProfileSchema,
} from "@trieoh/identityx-api";
import type { ProfileSchemaI, ProfileSchemaInputI } from "../model";

export const projectProfileSchemaQueryOptions = (projectId: string) =>
  queryOptions({
    queryKey: ["projects", projectId, "profile-schema"],
    queryFn: () =>
      getProjectProfileSchema(projectId).then(orvalData<ProfileSchemaI>),
    retry: false,
  });

export const upsertProjectProfileSchemaFn = createClientOnlyFn(
  (projectId: string, data: ProfileSchemaInputI) =>
    upsertProjectProfileSchema(projectId, data).then(orvalData<ProfileSchemaI>),
);
