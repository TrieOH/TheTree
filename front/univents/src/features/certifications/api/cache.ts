import type { QueryClient } from "@tanstack/react-query";
import type { CertificationTemplateI } from "../model";
import { certificationKeys } from "./query-keys";

function upsertTemplate(
  templates: CertificationTemplateI[] | undefined,
  template: CertificationTemplateI,
) {
  if (!templates) return templates;
  const index = templates.findIndex((item) => item.id === template.id);
  if (index === -1) return [template, ...templates];
  const next = [...templates];
  next[index] = template;
  return next;
}

export function syncCertificationTemplateCache(
  queryClient: QueryClient,
  template: CertificationTemplateI,
) {
  queryClient.setQueryData<CertificationTemplateI[]>(
    certificationKeys.templatesByEdition(template.edition_id),
    (old) => upsertTemplate(old, template),
  );
  queryClient.setQueryData<CertificationTemplateI>(
    certificationKeys.templateById(template.id),
    (old) => (old ? template : old),
  );
}

export function removeCertificationTemplateCache(
  queryClient: QueryClient,
  editionId: string,
  templateId: string,
) {
  queryClient.setQueryData<CertificationTemplateI[]>(
    certificationKeys.templatesByEdition(editionId),
    (old) => old?.filter((template) => template.id !== templateId),
  );
  queryClient.removeQueries({
    queryKey: certificationKeys.templateById(templateId),
  });
  queryClient.removeQueries({
    queryKey: certificationKeys.templateLinks(templateId),
  });
}
