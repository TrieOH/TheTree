export const certificationKeys = {
  all: ["certifications"] as const,

  templates: () => [...certificationKeys.all, "templates"] as const,
  templateLists: () => [...certificationKeys.templates(), "list"] as const,
  templatesByEdition: (editionId: string) =>
    [...certificationKeys.templateLists(), editionId] as const,
  templateById: (templateId: string) =>
    [...certificationKeys.templates(), "detail", templateId] as const,
  templateLinks: (templateId: string) =>
    [...certificationKeys.templates(), "links", templateId] as const,

  issued: () => [...certificationKeys.all, "issued"] as const,
  issuedById: (certificationId: string) =>
    [...certificationKeys.issued(), certificationId] as const,
  issuedByUser: () => [...certificationKeys.issued(), "user"] as const,
  issuedByEdition: (editionId: string) =>
    [...certificationKeys.issued(), "edition", editionId] as const,
  emissionErrorsByEdition: (editionId: string) =>
    [...certificationKeys.all, "emission-errors", editionId] as const,

  verification: (hash: string) =>
    [...certificationKeys.all, "verification", hash] as const,
};
