import { QueryClient } from "@tanstack/react-query";
import { describe, expect, it } from "vitest";
import type { CertificationTemplateI } from "../model";
import {
  removeCertificationTemplateCache,
  syncCertificationTemplateCache,
} from "./cache";
import { certificationKeys } from "./query-keys";

const template = (
  overrides: Partial<CertificationTemplateI> = {},
): CertificationTemplateI => ({
  id: "template-1",
  edition_id: "edition-1",
  kind: "edition_attendance",
  name: "Certificate",
  description: null,
  design_data: { background: null, elements: [] },
  created_at: "2026-01-01T00:00:00Z",
  ...overrides,
});

describe("certification template cache synchronization", () => {
  it("does not create a partial template list or detail", () => {
    const queryClient = new QueryClient();

    syncCertificationTemplateCache(queryClient, template());

    expect(
      queryClient.getQueryData(
        certificationKeys.templatesByEdition("edition-1"),
      ),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(certificationKeys.templateById("template-1")),
    ).toBeUndefined();
  });

  it("updates loaded template lists and details", () => {
    const queryClient = new QueryClient();
    const updated = template({ name: "Updated certificate" });
    queryClient.setQueryData(
      certificationKeys.templatesByEdition("edition-1"),
      [template()],
    );
    queryClient.setQueryData(
      certificationKeys.templateById("template-1"),
      template(),
    );

    syncCertificationTemplateCache(queryClient, updated);

    expect(
      queryClient.getQueryData(
        certificationKeys.templatesByEdition("edition-1"),
      ),
    ).toEqual([updated]);
    expect(
      queryClient.getQueryData(certificationKeys.templateById("template-1")),
    ).toEqual(updated);
  });

  it("removes the template and its links", () => {
    const queryClient = new QueryClient();
    queryClient.setQueryData(
      certificationKeys.templatesByEdition("edition-1"),
      [template()],
    );
    queryClient.setQueryData(
      certificationKeys.templateById("template-1"),
      template(),
    );
    queryClient.setQueryData(certificationKeys.templateLinks("template-1"), []);

    removeCertificationTemplateCache(queryClient, "edition-1", "template-1");

    expect(
      queryClient.getQueryData(
        certificationKeys.templatesByEdition("edition-1"),
      ),
    ).toEqual([]);
    expect(
      queryClient.getQueryData(certificationKeys.templateById("template-1")),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(certificationKeys.templateLinks("template-1")),
    ).toBeUndefined();
  });
});
