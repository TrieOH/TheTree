import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  deleteEmailTemplate,
  listEmailTemplates,
  putEmailTemplate,
} from "@trieoh/identityx-api";
import type {
  EffectiveEmailTemplate,
  EmailTemplateBody,
  EmailTemplateKind,
} from "@trieoh/identityx-api/schemas";

export const emailTemplatesQueryOptions = (projectId: string) =>
  queryOptions({
    queryKey: ["projects", projectId, "email-templates"],
    queryFn: () =>
      listEmailTemplates(projectId).then(orvalData<EffectiveEmailTemplate[]>),
  });

export const saveEmailTemplate = createClientOnlyFn(
  (projectId: string, kind: EmailTemplateKind, body: EmailTemplateBody) =>
    putEmailTemplate(projectId, kind, body).then(
      orvalData<EffectiveEmailTemplate>,
    ),
);

export const restoreEmailTemplate = createClientOnlyFn(
  (projectId: string, kind: EmailTemplateKind) =>
    deleteEmailTemplate(projectId, kind).then(orvalData<void>),
);
