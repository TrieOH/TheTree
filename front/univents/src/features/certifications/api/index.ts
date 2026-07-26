import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import {
  authFetcher,
  publicQueryFetcher,
  tanstackQueryFetcher,
} from "@/shared/lib/api/fetch";
import type {
  CertificationI,
  CertificationTargetType,
  CertificationTemplateCreateI,
  CertificationTemplateI,
  CertifyUserRequestI,
  SetCertificationTemplateRequestI,
  VerifyCertificationResponseI,
} from "../model";
import { certificationKeys } from "./query-keys";

export const createCertificationTemplateFn = createClientOnlyFn(
  (
    eventId: string,
    editionId: string,
    templateData: CertificationTemplateCreateI,
  ) => {
    return authFetcher.post<CertificationTemplateI>(
      `/events/${eventId}/editions/${editionId}/certification-templates`,
      templateData,
    );
  },
);

export const verifyCertificationHashFn = createClientOnlyFn((hash: string) => {
  return publicQueryFetcher<VerifyCertificationResponseI>(`/verify/${hash}`);
});

export const getAllCertificationTemplatesFn = createClientOnlyFn(
  async (eventId: string, editionId: string) => {
    return tanstackQueryFetcher<CertificationTemplateI[]>(
      `/events/${eventId}/editions/${editionId}/certification-templates`,
    );
  },
);

export const allCertificationTemplatesQueryOptions = (
  eventId: string,
  editionId: string,
) => {
  return queryOptions({
    queryKey: certificationKeys.templatesByEdition(eventId, editionId),
    queryFn: () => getAllCertificationTemplatesFn(eventId, editionId),
  });
};

export const getCertificationTemplateFn = createClientOnlyFn(
  (eventId: string, editionId: string, templateId: string) => {
    return tanstackQueryFetcher<CertificationTemplateI>(
      `/events/${eventId}/editions/${editionId}/certification-templates/${templateId}`,
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
    queryFn: () => getCertificationTemplateFn(eventId, editionId, templateId),
  });
};

export const setEditionCertificationTemplateFn = createClientOnlyFn(
  (
    eventId: string,
    editionId: string,
    data: SetCertificationTemplateRequestI,
  ) => {
    return authFetcher.patch<null>(
      `/events/${eventId}/editions/${editionId}/certification-templates/set`,
      data,
    );
  },
);

export const setActivityCertificationTemplateFn = createClientOnlyFn(
  (
    eventId: string,
    editionId: string,
    activityId: string,
    data: SetCertificationTemplateRequestI,
  ) => {
    return authFetcher.patch<null>(
      `/events/${eventId}/editions/${editionId}/activities/${activityId}/certification-templates/set`,
      data,
    );
  },
);

export const certifyUserFn = createClientOnlyFn(
  (
    eventId: string,
    editionId: string,
    userId: string,
    data: CertifyUserRequestI,
  ) => {
    return authFetcher.post<CertificationI>(
      `/events/${eventId}/editions/${editionId}/users/${userId}/certifications`,
      data,
    );
  },
);

export const getCertificationFn = createClientOnlyFn(
  (eventId: string, editionId: string, certId: string) => {
    return tanstackQueryFetcher<CertificationI>(
      `/events/${eventId}/editions/${editionId}/certifications/${certId}`,
    );
  },
);

export const certificationQueryOptions = (
  eventId: string,
  editionId: string,
  certId: string,
) => {
  return queryOptions({
    queryKey: certificationKeys.issuedById(eventId, editionId, certId),
    queryFn: () => getCertificationFn(eventId, editionId, certId),
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
  async (userId: string) => {
    return tanstackQueryFetcher<CertificationI[]>(
      `/users/${userId}/certifications`,
    );
  },
);

export const certificationsByUserQueryOptions = (userId: string) => {
  return queryOptions({
    queryKey: certificationKeys.issuedByUser(userId),
    queryFn: () => getAllCertificationsByUserFn(userId),
  });
};
