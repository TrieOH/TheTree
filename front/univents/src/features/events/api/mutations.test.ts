import { QueryClient } from "@tanstack/react-query";
import { describe, expect, it } from "vitest";
import type { EventI } from "../model";
import { syncEventCaches } from "./cache";
import { eventKeys } from "./query-keys";

const event = (overrides: Partial<EventI> = {}): EventI => ({
  id: "event-1",
  owner_id: "owner-1",
  full_name: "Event",
  slug: "event",
  status: "active",
  created_at: "2026-01-01T00:00:00Z",
  ...overrides,
});

describe("event cache synchronization", () => {
  it("does not create partial lists when their cache is missing", () => {
    const queryClient = new QueryClient();

    syncEventCaches(queryClient, event());

    expect(queryClient.getQueryData(eventKeys.ownLists())).toBeUndefined();
    expect(queryClient.getQueryData(eventKeys.publicLists())).toBeUndefined();
    expect(queryClient.getQueryData(eventKeys.joinedLists())).toBeUndefined();
  });

  it("updates loaded lists and the loaded public detail", () => {
    const queryClient = new QueryClient();
    const previous = event({ full_name: "Old name" });
    const updated = event({ full_name: "New name" });

    queryClient.setQueryData(eventKeys.ownLists(), [previous]);
    queryClient.setQueryData(eventKeys.publicLists(), [previous]);
    queryClient.setQueryData(
      eventKeys.detail.publicBySlug(previous.slug),
      previous,
    );

    syncEventCaches(queryClient, updated);

    expect(queryClient.getQueryData(eventKeys.ownLists())).toEqual([updated]);
    expect(queryClient.getQueryData(eventKeys.publicLists())).toEqual([
      updated,
    ]);
    expect(
      queryClient.getQueryData(eventKeys.detail.publicBySlug(updated.slug)),
    ).toEqual(updated);
  });

  it("removes a draft event from an already loaded public list", () => {
    const queryClient = new QueryClient();
    const published = event();

    queryClient.setQueryData(eventKeys.publicLists(), [published]);
    syncEventCaches(queryClient, event({ status: "draft" }));

    expect(queryClient.getQueryData(eventKeys.publicLists())).toEqual([]);
  });
});
