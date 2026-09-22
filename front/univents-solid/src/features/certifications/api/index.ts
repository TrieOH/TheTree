import { orvalData } from "@trieoh/api-client";
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
