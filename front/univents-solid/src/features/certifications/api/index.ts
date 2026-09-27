import { orvalData } from "@trieoh/api-client";
import type { QueryClient } from "@trieoh/front-core-solid";
import { Effect } from "effect";
import {
  createCertificationTemplate,
  deleteCertificationTemplate,
  emitProgramCertifications,
  getCertification,
  getCertificationTemplate,
  invalidateCertification,
  listCertificationEmissionErrors,
  listCertificationTemplates,
  listCertificationTemplateLinks,
  listEditionCertifications,
  listMyCertifications,
  updateCertificationTemplate,
  verifyCertification,
} from "@trieoh/univents-api";
import type {
  CertificationEmissionErrorI,
  CertificationI,
  CertificationTemplateI,
  CertificationTemplateProgramI,
  VerifyCertificationResponseI,
} from "../model";
import { certificationKeys } from "./query-keys";

export const myCertificationsQueryOptions = () => ({
  queryKey: certificationKeys.mine(),
  queryFn: () => listMyCertifications().then(orvalData<CertificationI[]>),
});

export const allCertificationTemplatesQueryOptions = (editionId: string) => ({
  queryKey: certificationKeys.templatesByEdition(editionId),
  queryFn: () =>
    listCertificationTemplates(editionId, { public: true }).then(
      orvalData<CertificationTemplateI[]>,
    ),
});

export const certificationTemplateQueryOptions = (templateId: string) => ({
  queryKey: certificationKeys.templateById(templateId),
  queryFn: () =>
    getCertificationTemplate(templateId, { public: true }).then(
      orvalData<CertificationTemplateI>,
    ),
});

export const certificationsByEditionQueryOptions = (editionId: string) => ({
  queryKey: certificationKeys.issuedByEdition(editionId),
  queryFn: () =>
    listEditionCertifications(editionId).then(orvalData<CertificationI[]>),
});

export const emissionErrorsByEditionQueryOptions = (editionId: string) => ({
  queryKey: certificationKeys.emissionErrorsByEdition(editionId),
  queryFn: () =>
    listCertificationEmissionErrors(editionId).then(
      orvalData<CertificationEmissionErrorI[]>,
    ),
});

export const certificationTemplateLinksQueryOptions = (templateId: string) => ({
  queryKey: certificationKeys.templateLinks(templateId),
  queryFn: () =>
    listCertificationTemplateLinks(templateId, { public: true }).then(
      orvalData<CertificationTemplateProgramI[]>,
    ),
});

export const fetchEditionProgramTemplateLinksEffect = (
  templates: CertificationTemplateI[],
  queryClient?: QueryClient,
) =>
  Effect.gen(function* () {
    if (templates.length === 0) return new Map<string, string>();

    const results = yield* Effect.forEach(
      templates,
      (template) =>
        Effect.tryPromise({
          try: () => {
            if (queryClient) {
              return queryClient.fetchQuery(
                certificationTemplateLinksQueryOptions(template.id),
              );
            }
            return listCertificationTemplateLinks(template.id, {
              public: true,
            }).then(orvalData<CertificationTemplateProgramI[]>);
          },
          catch: () => [] as CertificationTemplateProgramI[],
        }).pipe(
          Effect.map((links) => ({
            templateId: template.id,
            links: links ?? [],
          })),
        ),
      { concurrency: "unbounded" },
    );

    const map = new Map<string, string>();
    for (const { templateId, links } of results) {
      for (const link of links) {
        map.set(link.program_id, templateId);
      }
    }
    return map;
  });

export const editionProgramTemplateLinksQueryOptions = (
  editionId: string,
  templates: CertificationTemplateI[],
  queryClient?: QueryClient,
) => {
  const templateIdsKey = templates
    .map((t) => t.id)
    .sort()
    .join(",");

  return {
    queryKey: certificationKeys.editionProgramLinks(editionId, templateIdsKey),
    queryFn: () =>
      Effect.runPromise(
        fetchEditionProgramTemplateLinksEffect(templates, queryClient),
      ),
    enabled: Boolean(editionId && templates.length > 0),
  };
};

export const certificationVerificationQueryOptions = (hash: string) => ({
  queryKey: certificationKeys.verification(hash),
  queryFn: () =>
    verifyCertification(hash, { public: true }).then(
      orvalData<VerifyCertificationResponseI>,
    ),
});

export * from "./query-keys";

export {
  createCertificationTemplate,
  deleteCertificationTemplate,
  getCertification,
  getCertificationTemplate,
  invalidateCertification,
  listCertificationEmissionErrors,
  listCertificationTemplates,
  listCertificationTemplateLinks,
  listEditionCertifications,
  listMyCertifications,
  updateCertificationTemplate,
  verifyCertification,
  emitProgramCertifications,
};
