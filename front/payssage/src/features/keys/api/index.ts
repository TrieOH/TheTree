import { createServerFn } from "@tanstack/react-start";
import { identityXAccessClient as accessClient } from "#/shared/lib/api/identityx-access.server";
import type { ApiKeyCreateI } from "../model";

export const listCapabilitiesServerFn = createServerFn({ method: "GET" })
  .validator((data: { projectId: string }) => data)
  .handler(async ({ data }) => {
    const response = await accessClient.capabilities.list(data.projectId);
    if (!response.success) return [];
    return response.data;
  });

export const createApiKeyServerFn = createServerFn({ method: "POST" })
  .validator(
    (data: { projectId: string; subjectId?: string; payload: ApiKeyCreateI }) =>
      data,
  )
  .handler(async ({ data }) => {
    const response = await accessClient.apiKeys.create(data.projectId, {
      ...data.payload,
      ...(data.subjectId ? { subject_id: data.subjectId } : {}),
      env: "production",
    });

    if (!response.success) {
      throw response;
    }

    return response.data;
  });
