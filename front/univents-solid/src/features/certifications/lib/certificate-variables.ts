import type { CertificationTemplateElement } from "../model";
import type { CertificateVariableKey } from "../editor/constants";
import type { CertificationTemplateDraft } from "../editor/types";

export type { CertificateVariableKey };

export type CertificateVariableValues = Partial<
  Record<CertificateVariableKey, string>
>;

export const SAMPLE_CERTIFICATE_VARIABLES: CertificateVariableValues = {
  participant_name: "Nome do Participante",
  event_name: "Nome do Evento",
  edition_name: "Nome da Edição",
  activity_name: "Nome da Atividade",
  participation_type: "edição",
  location: "Local do Evento",
  workload_hours: "8",
  participation_date: "15 de março de 2026",
  certified_at: "15 de mar. de 2026, 12:00",
  cert_hash: "UNIV-2026-CERT-SAMPLE",
  verify_url: "https://univents.app/verify/UNIV-2026-CERT-SAMPLE",
};

/**
 * Formats a duration in milliseconds into a workload hours string in pt-BR.
 * Integer hours are formatted as whole numbers ("8", "20").
 * Fractional hours are formatted with 1 decimal place with comma ("2,5", "1,5").
 */
export function formatWorkloadHours(milliseconds: number): string {
  if (!Number.isFinite(milliseconds) || milliseconds <= 0) {
    return "";
  }
  const hours = milliseconds / 3_600_000;
  if (hours <= 0) return "";
  return Number.isInteger(hours)
    ? String(hours)
    : hours.toFixed(1).replace(".", ",");
}

/**
 * Returns a human-friendly label for workload hours, e.g. "1 hora", "8 horas", "2,5 horas",
 * optionally appending "(aproximado)" when calculated from general schedule rather than confirmed attendance.
 */
export function formatWorkloadLabel(
  hours: string | number | undefined | null,
  options?: { approximate?: boolean },
): string {
  if (!hours) return "";
  const str = String(hours).trim();
  if (!str) return "";
  const base = `${str} ${str === "1" ? "hora" : "horas"}`;
  return options?.approximate ? `${base} (aproximado)` : base;
}

/**
 * Formats a date string (ISO) into a localized Brazilian long date: "15 de março de 2026".
 */
export function formatParticipationDate(dateString: string | undefined | null): string {
  if (!dateString) return "";
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) return "";
  return date.toLocaleDateString("pt-BR", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  });
}

/**
 * Formats certification issuance date with time: "15 de mar. de 2026, 12:00".
 */
export function formatCertifiedAt(dateString: string | undefined | null): string {
  if (!dateString) return "";
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) return "";
  return date.toLocaleString("pt-BR", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export interface OccurrenceTimeSpan {
  id?: string;
  starts_at: string;
  ends_at: string;
  program_id?: string | null;
}

export interface ResolveWorkloadOccurrencesParams {
  occurrences: OccurrenceTimeSpan[];
  programId?: string | null;
  attendedOccurrences?: OccurrenceTimeSpan[];
}

/**
 * Filters and resolves occurrences to calculate workload and participation dates.
 * If user-specific attended occurrences are available and match the scope, they take precedence.
 * Otherwise, falls back to public edition/program occurrences.
 */
export function resolveWorkloadOccurrences(
  params: ResolveWorkloadOccurrencesParams,
): OccurrenceTimeSpan[] {
  const { occurrences, programId, attendedOccurrences } = params;

  if (attendedOccurrences && attendedOccurrences.length > 0) {
    const attended = programId
      ? attendedOccurrences.filter(
          (occ) => !occ.program_id || occ.program_id === programId,
        )
      : attendedOccurrences;
    if (attended.length > 0) {
      return attended;
    }
  }

  if (programId) {
    const programOccs = occurrences.filter((occ) => occ.program_id === programId);
    if (programOccs.length > 0) {
      return programOccs;
    }
  }

  return occurrences;
}

/**
 * Sums the total duration in hours for the given list of occurrences.
 */
export function calculateWorkloadHours(occurrences: OccurrenceTimeSpan[]): string {
  const milliseconds = occurrences.reduce((total, occ) => {
    const start = new Date(occ.starts_at).getTime();
    const end = new Date(occ.ends_at).getTime();
    if (Number.isNaN(start) || Number.isNaN(end) || end <= start) {
      return total;
    }
    return total + (end - start);
  }, 0);

  return formatWorkloadHours(milliseconds);
}

/**
 * Finds the earliest start date among occurrences and formats it.
 */
export function resolveFirstOccurrenceDate(occurrences: OccurrenceTimeSpan[]): string {
  if (!occurrences || occurrences.length === 0) return "";

  const valid = occurrences.filter((occ) => !Number.isNaN(new Date(occ.starts_at).getTime()));
  if (valid.length === 0) return "";

  const sorted = [...valid].sort(
    (a, b) => new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime(),
  );

  return formatParticipationDate(sorted[0].starts_at);
}

export interface BuildCertificateVariablesParams {
  participantName?: string | null;
  eventName?: string | null;
  editionName?: string | null;
  editionStartsAt?: string | null;
  editionEndsAt?: string | null;
  programName?: string | null;
  programId?: string | null;
  locationName?: string | null;
  occurrences?: OccurrenceTimeSpan[];
  attendedOccurrences?: OccurrenceTimeSpan[];
  issuedAt?: string | null;
  certHash?: string | null;
  origin?: string | null;
}

/**
 * Builds the complete CertificateVariableValues dictionary with all resolved variables.
 * All derived variables (workload_hours, participation_date, verify_url, participation_type, activity_name)
 * are calculated directly from domain models to ensure consistency.
 */
export function buildCertificateVariables(
  params: BuildCertificateVariablesParams,
): CertificateVariableValues {
  const isProgram = Boolean(params.programId);

  const workloadOccurrences = resolveWorkloadOccurrences({
    occurrences: params.occurrences ?? [],
    programId: params.programId,
    attendedOccurrences: params.attendedOccurrences,
  });

  let workloadHours = calculateWorkloadHours(workloadOccurrences);
  let participationDate = resolveFirstOccurrenceDate(workloadOccurrences);

  // Fallback to edition schedule if no specific occurrences were found and this is an edition cert
  if (!workloadHours && !isProgram && params.editionStartsAt && params.editionEndsAt) {
    const startMs = new Date(params.editionStartsAt).getTime();
    const endMs = new Date(params.editionEndsAt).getTime();
    if (!Number.isNaN(startMs) && !Number.isNaN(endMs) && endMs > startMs) {
      workloadHours = formatWorkloadHours(endMs - startMs);
    }
  }

  if (!participationDate && !isProgram && params.editionStartsAt) {
    participationDate = formatParticipationDate(params.editionStartsAt);
  }

  const activityName = isProgram
    ? params.programName || params.editionName || ""
    : params.editionName || "";

  const hostOrigin =
    params.origin ||
    (typeof window !== "undefined" && window.location?.origin
      ? window.location.origin
      : "");

  const verifyUrl =
    params.certHash && hostOrigin ? `${hostOrigin}/verify/${params.certHash}` : "";

  return {
    participant_name: params.participantName || "Participante",
    event_name: params.eventName || "Evento",
    edition_name: params.editionName || "",
    activity_name: activityName,
    participation_type: isProgram ? "atividade" : "edição",
    location: params.locationName || "",
    workload_hours: workloadHours,
    participation_date: participationDate,
    certified_at: params.issuedAt ? formatCertifiedAt(params.issuedAt) : "",
    cert_hash: params.certHash || "",
    verify_url: verifyUrl,
  };
}

/**
 * Replaces mustache-like variable tokens with actual certificate values.
 */
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
