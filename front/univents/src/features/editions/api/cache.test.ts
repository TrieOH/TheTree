import { QueryClient } from "@tanstack/react-query";
import { describe, expect, it } from "vitest";
import type { EditionI } from "../model";
import { syncEditionCaches } from "./cache";
import { editionKeys } from "./query-keys";

const edition = (overrides: Partial<EditionI> = {}): EditionI => ({
  id: "edition-1",
  event_id: "event-1",
  name: "Edition",
  slug: "edition",
  tagline: null,
  description: null,
  is_draft: false,
  registration_opens_at: null,
  starts_at: "2026-01-01T00:00:00Z",
  ends_at: "2026-01-02T00:00:00Z",
  location_name: null,
  location_description: null,
  logo_url: null,
  banner_url: null,
  contact_email: null,
  created_by: "actor-1",
  created_at: "2025-01-01T00:00:00Z",
  updated_at: null,
  deleted_at: null,
  status: "past",
  ...overrides,
});

describe("edition cache synchronization", () => {
  it("does not create partial lists when their cache is missing", () => {
    const queryClient = new QueryClient();

    syncEditionCaches(queryClient, edition());

    expect(
      queryClient.getQueryData(editionKeys.adminListByEvent("event-1")),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(editionKeys.publicListByEvent("event-1")),
    ).toBeUndefined();
  });

  it("updates loaded admin and public lists", () => {
    const queryClient = new QueryClient();
    const previous = edition({ name: "Old name" });
    const updated = edition({ name: "New name" });

    queryClient.setQueryData(editionKeys.adminListByEvent("event-1"), [
      previous,
    ]);
    queryClient.setQueryData(editionKeys.publicListByEvent("event-1"), [
      previous,
    ]);

    syncEditionCaches(queryClient, updated);

    expect(
      queryClient.getQueryData(editionKeys.adminListByEvent("event-1")),
    ).toEqual([updated]);
    expect(
      queryClient.getQueryData(editionKeys.publicListByEvent("event-1")),
    ).toEqual([updated]);
  });

  it("removes drafts from public lists without creating missing caches", () => {
    const queryClient = new QueryClient();
    const published = edition();

    queryClient.setQueryData(editionKeys.publicListByEvent("event-1"), [
      published,
    ]);
    syncEditionCaches(
      queryClient,
      edition({ is_draft: true, status: "draft" }),
    );

    expect(
      queryClient.getQueryData(editionKeys.publicListByEvent("event-1")),
    ).toEqual([]);
  });

  it("invalidates derived public queries for published editions", () => {
    const queryClient = new QueryClient();
    queryClient.setQueryData(editionKeys.activeByEvent("event-1"), edition());
    queryClient.setQueryData(editionKeys.pastByEvent("event-1"), []);
    queryClient.setQueryData(editionKeys.upcomingByEvent("event-1"), []);

    syncEditionCaches(queryClient, edition({ name: "Updated" }));

    expect(
      queryClient.getQueryState(editionKeys.activeByEvent("event-1"))
        ?.isInvalidated,
    ).toBe(true);
    expect(
      queryClient.getQueryState(editionKeys.pastByEvent("event-1"))
        ?.isInvalidated,
    ).toBe(true);
    expect(
      queryClient.getQueryState(editionKeys.upcomingByEvent("event-1"))
        ?.isInvalidated,
    ).toBe(true);
  });
});
