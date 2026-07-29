import type { CertificationTargetType } from "../model";

export const certificationKeys = {
  all: ["certifications"] as const,

  templates: () => [...certificationKeys.all, "templates"] as const,
  templateLists: () => [...certificationKeys.templates(), "list"] as const,
  templatesByEdition: (eventId: string, editionId: string) =>
    [...certificationKeys.templateLists(), eventId, editionId] as const,
  templateById: (eventId: string, editionId: string, templateId: string) =>
    [
      ...certificationKeys.templates(),
      "detail",
      eventId,
      editionId,
      templateId,
    ] as const,

  issued: () => [...certificationKeys.all, "issued"] as const,
  issuedById: (eventId: string, editionId: string, certificationId: string) =>
    [
      ...certificationKeys.issued(),
      "detail",
      eventId,
      editionId,
      certificationId,
    ] as const,
  issuedByTarget: (targetType: CertificationTargetType, targetId: string) =>
    [...certificationKeys.issued(), "target", targetType, targetId] as const,
  issuedByUser: (userId: string) =>
    [...certificationKeys.issued(), "user", userId] as const,

  verification: (hash: string) =>
    [...certificationKeys.all, "verification", hash] as const,
};
