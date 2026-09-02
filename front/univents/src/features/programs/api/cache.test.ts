import { QueryClient } from "@tanstack/react-query";
import { describe, expect, it } from "vitest";
import type { OccurrenceI, ProgramI } from "../model";
import {
  removeOccurrenceCaches,
  removeProgramCaches,
  syncOccurrenceCache,
  syncProgramCache,
} from "./cache";
import { programKeys } from "./query-keys";

const program = (overrides: Partial<ProgramI> = {}): ProgramI => ({
  id: "program-1",
  edition_id: "edition-1",
  kind: "activity",
  name: "Program",
  description: undefined,
  min_access_level: 0,
  staff_only: false,
  price: 0,
  banner_url: null,
  created_at: "2026-01-01T00:00:00Z",
  updated_at: null,
  deleted_at: null,
  ...overrides,
});

const occurrence = (overrides: Partial<OccurrenceI> = {}): OccurrenceI => ({
  id: "occurrence-1",
  program_id: "program-1",
  edition_id: "edition-1",
  starts_at: "2026-01-01T10:00:00Z",
  ends_at: "2026-01-01T11:00:00Z",
  max_capacity: 10,
  created_at: "2026-01-01T00:00:00Z",
  updated_at: null,
  deleted_at: null,
  ...overrides,
});

describe("program cache synchronization", () => {
  it("does not create partial program or occurrence lists", () => {
    const queryClient = new QueryClient();

    syncProgramCache(queryClient, program());
    syncOccurrenceCache(queryClient, occurrence());

    expect(
      queryClient.getQueryData(programKeys.byEdition("edition-1")),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(programKeys.occurrences("edition-1")),
    ).toBeUndefined();
  });

  it("updates loaded program and occurrence lists", () => {
    const queryClient = new QueryClient();
    const updatedProgram = program({ name: "Updated program" });
    const updatedOccurrence = occurrence({ max_capacity: 20 });

    queryClient.setQueryData(programKeys.byEdition("edition-1"), [program()]);
    queryClient.setQueryData(programKeys.occurrences("edition-1"), [
      occurrence(),
    ]);

    syncProgramCache(queryClient, updatedProgram);
    syncOccurrenceCache(queryClient, updatedOccurrence);

    expect(
      queryClient.getQueryData(programKeys.byEdition("edition-1")),
    ).toEqual([updatedProgram]);
    expect(
      queryClient.getQueryData(programKeys.occurrences("edition-1")),
    ).toEqual([updatedOccurrence]);
  });

  it("removes an occurrence and its participants", () => {
    const queryClient = new QueryClient();
    queryClient.setQueryData(programKeys.occurrences("edition-1"), [
      occurrence(),
    ]);
    queryClient.setQueryData(programKeys.participants("occurrence-1"), []);

    removeOccurrenceCaches(queryClient, occurrence());
    expect(
      queryClient.getQueryData(programKeys.participants("occurrence-1")),
    ).toBeUndefined();
  });

  it("removes a program, its occurrences and their participants", () => {
    const queryClient = new QueryClient();
    queryClient.setQueryData(programKeys.byEdition("edition-1"), [program()]);

    queryClient.setQueryData(programKeys.occurrences("edition-1"), [
      occurrence(),
    ]);
    queryClient.setQueryData(programKeys.participants("occurrence-1"), []);
    removeProgramCaches(queryClient, program());

    expect(
      queryClient.getQueryData(programKeys.byEdition("edition-1")),
    ).toEqual([]);
    expect(
      queryClient.getQueryData(programKeys.occurrences("edition-1")),
    ).toEqual([]);
    expect(
      queryClient.getQueryData(programKeys.participants("occurrence-1")),
    ).toBeUndefined();
  });
});
