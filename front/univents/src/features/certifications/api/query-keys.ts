import type { CertificationTargetType } from "../model";

export const certificationKeys = {
  all: ["certifications"] as const,

  templates: () => [...certificationKeys.all, "templates"] as const,
  templateLists: () => [...certificationKeys.templates(), "list"] as const,
  templatesByEdition: (editionId: string) =>
    [...certificationKeys.templateLists(), editionId] as const,
  templateById: (eventId: string, editionId: string, templateId: string) =>
    [
      ...certificationKeys.templates(),
      "detail",
      eventId,
      editionId,
      templateId,
    ] as const,

  issued: () => [...certificationKeys.all, "issued"] as const,
  issuedById: (certificationId: string) =>
    [...certificationKeys.issued(), certificationId] as const,
  issuedByTarget: (targetType: CertificationTargetType, targetId: string) =>
    [...certificationKeys.issued(), "target", targetType, targetId] as const,
  issuedByUser: (userId: string) =>
    [...certificationKeys.issued(), "user", userId] as const,
  issuedByEdition: (editionId: string) =>
    [...certificationKeys.issued(), "edition", editionId] as const,
  emissionErrorsByEdition: (editionId: string) =>
    [...certificationKeys.all, "emission-errors", editionId] as const,

  verification: (hash: string) =>
    [...certificationKeys.all, "verification", hash] as const,
};
