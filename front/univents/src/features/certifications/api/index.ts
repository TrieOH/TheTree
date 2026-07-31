import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";

import {
  authFetcher,
  publicQueryFetcher,
  tanstackQueryFetcher,
} from "@/shared/lib/api/fetch";
import type {
  CertificationEmissionErrorI,
  CertificationI,
  CertificationTargetType,
  CertificationTemplateCreateI,
  CertificationTemplateI,
  CertificationTemplateProgramI,
  VerifyCertificationResponseI,
} from "../model";
import { certificationKeys } from "./query-keys";

export const createCertificationTemplateFn = createClientOnlyFn(
  async (editionId: string, templateData: CertificationTemplateCreateI) => {
    return authFetcher.post<CertificationTemplateI>(
      `/editions/${editionId}/certifications/templates`,
      templateData,
    );
  },
);

export const updateCertificationTemplateFn = createClientOnlyFn(
  (templateId: string, data: CertificationTemplateCreateI) =>
    authFetcher.put<CertificationTemplateI>(
      `/certifications/templates/${templateId}`,
      {
        ...data,
        design_data: data.design_data,
      },
    ),
);

export const deleteCertificationTemplateFn = createClientOnlyFn(
  (templateId: string) =>
    authFetcher.delete<null>(`/certifications/templates/${templateId}`),
);

export const verifyCertificationHashFn = createClientOnlyFn((hash: string) => {
  return publicQueryFetcher<VerifyCertificationResponseI>(`/verify/${hash}`);
});

export const getAllCertificationTemplatesFn = createClientOnlyFn(
  async (editionId: string) => {
    return tanstackQueryFetcher<CertificationTemplateI[]>(
      `/editions/${editionId}/certifications/templates`,
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
    return tanstackQueryFetcher<CertificationTemplateI>(
      `/certifications/templates/${templateId}`,
    );
  },
);

export const certificationTemplateQueryOptions = (
  eventId: string,
  editionId: string,
  templateId: string,
) => {
  return queryOptions({
    queryKey: certificationKeys.templateById(eventId, editionId, templateId),
    queryFn: () => getCertificationTemplateFn(templateId),
  });
};

export const getCertificationTemplateLinksFn = createClientOnlyFn(
  (templateId: string) =>
    tanstackQueryFetcher<CertificationTemplateProgramI[]>(
      `/certifications/templates/${templateId}/links`,
    ),
);

export const certificationTemplateLinksQueryOptions = (templateId: string) =>
  queryOptions({
    queryKey: [...certificationKeys.templates(), "links", templateId],
    queryFn: () => getCertificationTemplateLinksFn(templateId),
  });

export const linkCertificationTemplateFn = createClientOnlyFn(
  (templateId: string, programId: string) =>
    authFetcher.post<null>(`/certifications/templates/${templateId}/link`, {
      program_id: programId,
    }),
);

export const unlinkCertificationTemplateFn = createClientOnlyFn(
  (templateId: string) =>
    authFetcher.delete<null>(`/certifications/templates/${templateId}/link`),
);

export const getCertificationFn = createClientOnlyFn((certId: string) => {
  return tanstackQueryFetcher<CertificationI>(`/certifications/${certId}`);
});

export const certificationQueryOptions = (certId: string) => {
  return queryOptions({
    queryKey: certificationKeys.issuedById(certId),
    queryFn: () => getCertificationFn(certId),
  });
};

export const getAllCertificationsByTargetFn = createClientOnlyFn(
  async (targetType: CertificationTargetType, targetId: string) => {
    return tanstackQueryFetcher<CertificationI[]>(
      `/certifications?target_type=${targetType}&target_id=${targetId}`,
    );
  },
);

export const certificationsByTargetQueryOptions = (
  targetType: CertificationTargetType,
  targetId: string,
) => {
  return queryOptions({
    queryKey: certificationKeys.issuedByTarget(targetType, targetId),
    queryFn: () => getAllCertificationsByTargetFn(targetType, targetId),
  });
};

export const getAllCertificationsByUserFn = createClientOnlyFn(
  async (_userId: string) => {
    return tanstackQueryFetcher<CertificationI[]>(`/certifications`);
  },
);

export const certificationsByUserQueryOptions = (userId: string) => {
  return queryOptions({
    queryKey: certificationKeys.issuedByUser(userId),
    queryFn: () => getAllCertificationsByUserFn(userId),
  });
};

export const getCertificationsByEditionFn = createClientOnlyFn(
  (editionId: string) =>
    tanstackQueryFetcher<CertificationI[]>(
      `/editions/${editionId}/certifications`,
    ),
);

export const invalidateCertificationFn = createClientOnlyFn(
  ({ certificationId, reason }: { certificationId: string; reason?: string }) =>
    authFetcher.post<null>(
      `/certifications/${certificationId}/invalidate`,
      reason ? { reason } : {},
    ),
);

export const certificationsByEditionQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: certificationKeys.issuedByEdition(editionId),
    queryFn: () => getCertificationsByEditionFn(editionId),
  });

export const getEmissionErrorsByEditionFn = createClientOnlyFn(
  (editionId: string) =>
    tanstackQueryFetcher<CertificationEmissionErrorI[]>(
      `/editions/${editionId}/certifications/emission-errors`,
    ),
);

export const emissionErrorsByEditionQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: certificationKeys.emissionErrorsByEdition(editionId),
    queryFn: () => getEmissionErrorsByEditionFn(editionId),
  });
