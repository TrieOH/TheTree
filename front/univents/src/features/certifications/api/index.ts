import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type {
  CertifyUserRequestI,
  CertificationI,
  CertificationTemplateCreateI,
  CertificationTemplateI,
  SetCertificationTemplateRequestI,
} from "../model";
import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";

export const createCertificationTemplateFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  templateData: CertificationTemplateCreateI
) => {
  return authFetcher.post<CertificationTemplateI>(
    `/events/${eventId}/editions/${editionId}/certification-templates`,
    templateData
  );
});

export const getAllCertificationTemplatesFn = createClientOnlyFn(async (
  eventId: string,
  editionId: string
) => {
  try {
    return await tanstackQueryFetcher<CertificationTemplateI[]>(
      `/events/${eventId}/editions/${editionId}/certification-templates`
    );
  } catch {
    return [];
  }
});

export const allCertificationTemplatesQueryOptions = (
  eventId: string,
  editionId: string
) => {
  return queryOptions({
    queryKey: ["certifications", "templates", eventId, editionId],
    queryFn: () => getAllCertificationTemplatesFn(eventId, editionId),
  });
};

export const getCertificationTemplateFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  templateId: string
) => {
  return tanstackQueryFetcher<CertificationTemplateI>(
    `/events/${eventId}/editions/${editionId}/certification-templates/${templateId}`
  );
});

export const certificationTemplateQueryOptions = (
  eventId: string,
  editionId: string,
  templateId: string
) => {
  return queryOptions({
    queryKey: ["certifications", "template", eventId, editionId, templateId],
    queryFn: () => getCertificationTemplateFn(eventId, editionId, templateId),
  });
};

export const setEditionCertificationTemplateFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  data: SetCertificationTemplateRequestI
) => {
  return authFetcher.patch<null>(
    `/events/${eventId}/editions/${editionId}/certification-templates/set`,
    data
  );
});

export const setActivityCertificationTemplateFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  activityId: string,
  data: SetCertificationTemplateRequestI
) => {
  return authFetcher.patch<null>(
    `/events/${eventId}/editions/${editionId}/activities/${activityId}/certification-templates/set`,
    data
  );
});

export const certifyUserFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  userId: string,
  data: CertifyUserRequestI
) => {
  return authFetcher.post<CertificationI>(
    `/events/${eventId}/editions/${editionId}/users/${userId}/certifications`,
    data
  );
});

export const getCertificationFn = createClientOnlyFn((
  eventId: string,
  editionId: string,
  certId: string
) => {
  return tanstackQueryFetcher<CertificationI>(
    `/events/${eventId}/editions/${editionId}/certifications/${certId}`
  );
});

export const certificationQueryOptions = (
  eventId: string,
  editionId: string,
  certId: string
) => {
  return queryOptions({
    queryKey: ["certifications", "cert", eventId, editionId, certId],
    queryFn: () => getCertificationFn(eventId, editionId, certId),
  });
};

export const getAllCertificationsByTargetFn = createClientOnlyFn(async (
  targetType: string,
  targetId: string
) => {
  try {
    return await tanstackQueryFetcher<CertificationI[]>(
      `/certifications?target_type=${targetType}&target_id=${targetId}`
    );
  } catch {
    return [];
  }
});

export const certificationsByTargetQueryOptions = (
  targetType: string,
  targetId: string
) => {
  return queryOptions({
    queryKey: ["certifications", "target", targetType, targetId],
    queryFn: () => getAllCertificationsByTargetFn(targetType, targetId),
  });
};

export const getAllCertificationsByUserFn = createClientOnlyFn(async (userId: string) => {
  try {
    return await tanstackQueryFetcher<CertificationI[]>(
      `/users/${userId}/certifications`
    );
  } catch {
    return [];
  }
});

export const certificationsByUserQueryOptions = (userId: string) => {
  return queryOptions({
    queryKey: ["certifications", "user", userId],
    queryFn: () => getAllCertificationsByUserFn(userId),
  });
};
