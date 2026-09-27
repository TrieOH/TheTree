import { describe, expect, it } from "vitest";
import {
  buildCertificateVariables,
  calculateWorkloadHours,
  formatCertifiedAt,
  formatParticipationDate,
  formatWorkloadHours,
  formatWorkloadLabel,
  replaceCertificateVariables,
  resolveFirstOccurrenceDate,
  resolveWorkloadOccurrences,
  type OccurrenceTimeSpan,
} from "@/features/certifications/lib/certificate-variables";

describe("Certificate Variables Utility", () => {
  describe("formatWorkloadHours", () => {
    it("returns empty string for 0 or negative values", () => {
      expect(formatWorkloadHours(0)).toBe("");
      expect(formatWorkloadHours(-1000)).toBe("");
    });

    it("formats integer hours without decimal places", () => {
      expect(formatWorkloadHours(3_600_000)).toBe("1");
      expect(formatWorkloadHours(8 * 3_600_000)).toBe("8");
      expect(formatWorkloadHours(40 * 3_600_000)).toBe("40");
    });

    it("formats fractional hours with comma decimal separator in pt-BR", () => {
      expect(formatWorkloadHours(2.5 * 3_600_000)).toBe("2,5");
      expect(formatWorkloadHours(1.25 * 3_600_000)).toBe("1,3");
    });
  });

  describe("formatWorkloadLabel", () => {
    it("returns singular label for 1 hour", () => {
      expect(formatWorkloadLabel("1")).toBe("1 hora");
      expect(formatWorkloadLabel(1)).toBe("1 hora");
    });

    it("returns plural label for other hours", () => {
      expect(formatWorkloadLabel("8")).toBe("8 horas");
      expect(formatWorkloadLabel("2,5")).toBe("2,5 horas");
    });

    it("appends (aproximado) when approximate is true", () => {
      expect(formatWorkloadLabel("8", { approximate: true })).toBe("8 horas (aproximado)");
      expect(formatWorkloadLabel("1", { approximate: true })).toBe("1 hora (aproximado)");
    });

    it("returns empty string for empty inputs", () => {
      expect(formatWorkloadLabel("")).toBe("");
      expect(formatWorkloadLabel(null)).toBe("");
      expect(formatWorkloadLabel(undefined)).toBe("");
    });
  });

  describe("formatParticipationDate", () => {
    it("formats date in long Brazilian format", () => {
      const formatted = formatParticipationDate("2026-03-15T09:00:00Z");
      expect(formatted).toContain("15 de");
      expect(formatted).toContain("2026");
    });

    it("returns empty string for invalid or missing dates", () => {
      expect(formatParticipationDate(null)).toBe("");
      expect(formatParticipationDate("invalid-date")).toBe("");
    });
  });

  describe("formatCertifiedAt", () => {
    it("formats datetime with date and time", () => {
      const formatted = formatCertifiedAt("2026-03-15T14:30:00Z");
      expect(formatted).toContain("2026");
      expect(formatted.length).toBeGreaterThan(5);
    });

    it("returns empty string for invalid dates", () => {
      expect(formatCertifiedAt("")).toBe("");
      expect(formatCertifiedAt("not-a-date")).toBe("");
    });
  });

  describe("calculateWorkloadHours and resolveFirstOccurrenceDate", () => {
    const mockOccurrences: OccurrenceTimeSpan[] = [
      {
        id: "occ-2",
        starts_at: "2026-03-16T14:00:00Z",
        ends_at: "2026-03-16T18:00:00Z", // 4 hours
        program_id: "prog-1",
      },
      {
        id: "occ-1",
        starts_at: "2026-03-15T09:00:00Z",
        ends_at: "2026-03-15T13:00:00Z", // 4 hours
        program_id: "prog-1",
      },
      {
        id: "occ-3",
        starts_at: "2026-03-17T10:00:00Z",
        ends_at: "2026-03-17T12:00:00Z", // 2 hours
        program_id: "prog-2",
      },
    ];

    it("calculates total workload across all occurrences", () => {
      expect(calculateWorkloadHours(mockOccurrences)).toBe("10");
    });

    it("resolves the earliest occurrence date correctly", () => {
      const date = resolveFirstOccurrenceDate(mockOccurrences);
      expect(date).toContain("15 de");
      expect(date).toContain("2026");
    });

    it("resolves workload occurrences filtered by programId", () => {
      const filtered = resolveWorkloadOccurrences({
        occurrences: mockOccurrences,
        programId: "prog-1",
      });
      expect(filtered.length).toBe(2);
      expect(calculateWorkloadHours(filtered)).toBe("8");
    });

    it("prioritizes attendedOccurrences when provided", () => {
      const attended: OccurrenceTimeSpan[] = [
        {
          id: "occ-1",
          starts_at: "2026-03-15T09:00:00Z",
          ends_at: "2026-03-15T13:00:00Z", // 4 hours
          program_id: "prog-1",
        },
      ];

      const resolved = resolveWorkloadOccurrences({
        occurrences: mockOccurrences,
        programId: "prog-1",
        attendedOccurrences: attended,
      });

      expect(resolved.length).toBe(1);
      expect(calculateWorkloadHours(resolved)).toBe("4");
    });
  });

  describe("buildCertificateVariables", () => {
    it("builds complete variables dictionary with computed workload and dates", () => {
      const occurrences: OccurrenceTimeSpan[] = [
        {
          starts_at: "2026-05-10T08:00:00Z",
          ends_at: "2026-05-10T16:00:00Z", // 8 hours
          program_id: null,
        },
      ];

      const vars = buildCertificateVariables({
        participantName: "Maria Souza",
        eventName: "Semana de Tecnologia",
        editionName: "Edição 2026",
        locationName: "Auditório Central",
        occurrences,
        issuedAt: "2026-05-11T10:00:00Z",
        certHash: "hash-xyz-789",
        origin: "https://univents.app",
      });

      expect(vars.participant_name).toBe("Maria Souza");
      expect(vars.event_name).toBe("Semana de Tecnologia");
      expect(vars.edition_name).toBe("Edição 2026");
      expect(vars.location).toBe("Auditório Central");
      expect(vars.workload_hours).toBe("8");
      expect(vars.participation_date).toContain("10 de");
      expect(vars.certified_at).toContain("2026");
      expect(vars.cert_hash).toBe("hash-xyz-789");
      expect(vars.verify_url).toBe("https://univents.app/verify/hash-xyz-789");
      expect(vars.participation_type).toBe("edição");
    });

    it("falls back to edition schedule when no occurrences are present", () => {
      const vars = buildCertificateVariables({
        participantName: "Carlos Pereira",
        eventName: "Workshop de Cloud",
        editionName: "Turma 1",
        editionStartsAt: "2026-08-01T09:00:00Z",
        editionEndsAt: "2026-08-01T15:00:00Z", // 6 hours
        occurrences: [],
      });

      expect(vars.workload_hours).toBe("6");
      expect(vars.participation_date).toContain("01 de");
      expect(vars.participation_date).toContain("2026");
    });

    it("replaces variables in template text correctly", () => {
      const vars = {
        participant_name: "Lucas Costa",
        workload_hours: "12",
        participation_date: "20 de abril de 2026",
      };

      const templateText =
        "Certificamos que {{participant_name}} concluiu {{workload_hours}} horas em {{participation_date}}.";
      const resolved = replaceCertificateVariables(templateText, vars);

      expect(resolved).toBe(
        "Certificamos que Lucas Costa concluiu 12 horas em 20 de abril de 2026.",
      );
    });
  });
});
