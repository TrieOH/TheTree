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
  allEditionProgramLinks: () =>
    [...certificationKeys.templates(), "edition-program-links"] as const,
  editionProgramLinks: (editionId: string, templateIdsKey?: string) =>
    [
      ...certificationKeys.allEditionProgramLinks(),
      editionId,
      ...(templateIdsKey ? [templateIdsKey] : []),
    ] as const,

  issued: () => [...certificationKeys.all, "issued"] as const,
  issuedById: (certificationId: string) =>
    [...certificationKeys.issued(), certificationId] as const,
  issuedByUser: (actorId?: string) =>
    [...certificationKeys.issued(), "user", actorId ?? "mine"] as const,
  issuedByEdition: (editionId: string) =>
    [...certificationKeys.issued(), "edition", editionId] as const,
  emissionErrorsByEdition: (editionId: string) =>
    [...certificationKeys.all, "emission-errors", editionId] as const,

  verification: (hash: string) =>
    [...certificationKeys.all, "verification", hash] as const,

  mine: () => ["certifications", "mine"] as const,
};
