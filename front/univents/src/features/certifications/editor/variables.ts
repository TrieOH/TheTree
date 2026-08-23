import type { CertificationTemplateElement } from "../model";
import type { CertificateVariableKey } from "./constants";
import type { CertificationTemplateDraft } from "./types";

export type CertificateVariableValues = Partial<
  Record<CertificateVariableKey, string>
>;

export function replaceCertificateVariables(
  text: string,
  values: CertificateVariableValues,
): string {
  return text.replace(
    /\{\{(participant_name|event_name|edition_name|activity_name|participation_type|location|workload_hours|participation_date|certified_at|cert_hash|verify_url)\}\}/g,
    (token, key: CertificateVariableKey) => values[key] ?? token,
  );
}

function resolveElement(
  element: CertificationTemplateElement,
  values: CertificateVariableValues,
): CertificationTemplateElement {
  if (element.type === "text") {
    return {
      ...element,
      paragraphs: element.paragraphs.map((paragraph) => ({
        ...paragraph,
        runs: paragraph.runs.map((run) => ({
          ...run,
          text: replaceCertificateVariables(run.text, values),
        })),
      })),
    };
  }

  if (element.type === "hash") {
    return {
      ...element,
      hashLabel: replaceCertificateVariables(element.hashLabel, values),
      hash: replaceCertificateVariables(element.hash, values),
      linkLabel: replaceCertificateVariables(element.linkLabel, values),
      url: replaceCertificateVariables(element.url, values),
    };
  }

  return { ...element };
}

export function resolveCertificationTemplate(
  template: CertificationTemplateDraft,
  values: CertificateVariableValues,
): CertificationTemplateDraft {
  const designData = template.design_data ?? {
    canvas: undefined,
    background: null,
    elements: [],
  };
  return {
    ...template,
    design_data: {
      ...designData,
      elements: (designData.elements ?? []).map((element) =>
        resolveElement(element, values),
      ),
    },
  };
}
