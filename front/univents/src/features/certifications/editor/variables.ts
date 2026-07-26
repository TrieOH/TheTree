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
    /\{\{(activity_name|certified_at|cert_hash|verify_url)\}\}/g,
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
  return {
    ...template,
    data: {
      ...template.data,
      elements: template.data.elements.map((element) =>
        resolveElement(element, values),
      ),
    },
  };
}
