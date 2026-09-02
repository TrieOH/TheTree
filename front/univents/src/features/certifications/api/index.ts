import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  createCertificationTemplate,
  deleteCertificationTemplate,
  emitProgramCertifications,
  getCertification,
  getCertificationTemplate,
  invalidateCertification,
  linkCertificationTemplate,
  listCertificationEmissionErrors,
  listCertificationTemplateLinks,
  listCertificationTemplates,
  listEditionCertifications,
  listMyCertifications,
  unlinkCertificationTemplate,
  updateCertificationTemplate,
  verifyCertification,
} from "@trieoh/univents-api";
import type {
  CertificationEmissionErrorI,
  CertificationI,
  CertificationTemplateCreateI,
  CertificationTemplateI,
  CertificationTemplateProgramI,
  VerifyCertificationResponseI,
} from "../model";
import { certificationKeys } from "./query-keys";

export const createCertificationTemplateFn = createClientOnlyFn(
  async (editionId: string, templateData: CertificationTemplateCreateI) => {
    return withSpan("action:certification-template-create", () =>
      createCertificationTemplate(editionId, templateData).then(
        orvalData<CertificationTemplateI>,
      ),
    );
  },
);

export const updateCertificationTemplateFn = createClientOnlyFn(
  (templateId: string, data: CertificationTemplateCreateI) =>
    withSpan("action:certification-template-update", () =>
      updateCertificationTemplate(templateId, {
        ...data,
        design_data: data.design_data,
      }).then(orvalData<CertificationTemplateI>),
    ),
);

export const deleteCertificationTemplateFn = createClientOnlyFn(
  (templateId: string) =>
    withSpan("action:certification-template-delete", () =>
      deleteCertificationTemplate(templateId).then(orvalData<null>),
    ),
);

export const verifyCertificationHashFn = createClientOnlyFn((hash: string) => {
  return verifyCertification(hash, { public: true }).then(
    orvalData<VerifyCertificationResponseI>,
  );
});

export const getAllCertificationTemplatesFn = createClientOnlyFn(
  async (editionId: string) => {
    return listCertificationTemplates(editionId, { public: true }).then(
      orvalData<CertificationTemplateI[]>,
    );
  },
);

export const allCertificationTemplatesQueryOptions = (editionId: string) => {
  return queryOptions({
    queryKey: certificationKeys.templatesByEdition(editionId),
    queryFn: () => getAllCertificationTemplatesFn(editionId),
  });
};

export const getCertificationTemplateFn = createClientOnlyFn(
  async (templateId: string) => {
    return getCertificationTemplate(templateId, { public: true }).then(
      orvalData<CertificationTemplateI>,
    );
  },
);

export const certificationTemplateQueryOptions = (templateId: string) => {
  return queryOptions({
    queryKey: certificationKeys.templateById(templateId),
    queryFn: () => getCertificationTemplateFn(templateId),
  });
};

export const getCertificationTemplateLinksFn = createClientOnlyFn(
  (templateId: string) =>
    listCertificationTemplateLinks(templateId, { public: true }).then(
      orvalData<CertificationTemplateProgramI[]>,
    ),
);

export const certificationTemplateLinksQueryOptions = (templateId: string) =>
  queryOptions({
    queryKey: certificationKeys.templateLinks(templateId),
    queryFn: () => getCertificationTemplateLinksFn(templateId),
  });

export const linkCertificationTemplateFn = createClientOnlyFn(
  (templateId: string, programId: string) =>
    withSpan("action:certification-template-link", () =>
      linkCertificationTemplate(templateId, { program_id: programId }).then(
        orvalData<null>,
      ),
    ),
);

export const unlinkCertificationTemplateFn = createClientOnlyFn(
  (templateId: string, programId: string) =>
    withSpan("action:certification-template-unlink", () =>
      unlinkCertificationTemplate(templateId, { program_id: programId }).then(
        orvalData<null>,
      ),
    ),
);

export const emitProgramCertificationsFn = createClientOnlyFn(
  (programId: string) =>
    withSpan("action:program-certifications-emit", () =>
      emitProgramCertifications(programId).then(orvalData),
    ),
);

export const getCertificationFn = createClientOnlyFn((certId: string) => {
  return getCertification(certId).then(orvalData<CertificationI>);
});

export const certificationQueryOptions = (certId: string) => {
  return queryOptions({
    queryKey: certificationKeys.issuedById(certId),
    queryFn: () => getCertificationFn(certId),
  });
};

export const getAllCertificationsByUserFn = createClientOnlyFn(async () => {
  return listMyCertifications().then(orvalData<CertificationI[]>);
});

export const certificationsByUserQueryOptions = () => {
  return queryOptions({
    queryKey: certificationKeys.issuedByUser(),
    queryFn: () => getAllCertificationsByUserFn(),
  });
};

export const getCertificationsByEditionFn = createClientOnlyFn(
  (editionId: string) =>
    listEditionCertifications(editionId).then(orvalData<CertificationI[]>),
);

export const invalidateCertificationFn = createClientOnlyFn(
  ({ certificationId, reason }: { certificationId: string; reason?: string }) =>
    withSpan("action:certification-invalidate", () =>
      invalidateCertification(certificationId, { reason: reason ?? "" }).then(
        orvalData<null>,
      ),
    ),
);

export const certificationsByEditionQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: certificationKeys.issuedByEdition(editionId),
    queryFn: () => getCertificationsByEditionFn(editionId),
  });

export const getEmissionErrorsByEditionFn = createClientOnlyFn(
  (editionId: string) =>
    listCertificationEmissionErrors(editionId).then(
      orvalData<CertificationEmissionErrorI[]>,
    ),
);

export const emissionErrorsByEditionQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: certificationKeys.emissionErrorsByEdition(editionId),
    queryFn: () => getEmissionErrorsByEditionFn(editionId),
  });
