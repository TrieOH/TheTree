import { QueryClient } from "@tanstack/react-query";
import { describe, expect, it } from "vitest";
import type { BadgeTemplate } from "../model";
import { removeBadgeTemplateCache, syncBadgeTemplateCache } from "./cache";
import { badgeKeys } from "./query-keys";

const template = (overrides: Partial<BadgeTemplate> = {}): BadgeTemplate => ({
  id: "template-1",
  edition_id: "edition-1",
  name: "Badge",
  ticket_type_id: null,
  origin: null,
  design_data: {
    canvas: { width: 321, height: 204 },
    backgroundColor: "#ffffff",
    background: null,
    elements: [],
  },
  created_at: "2026-01-01T00:00:00Z",
  updated_at: null,
  ...overrides,
});

describe("badge template cache synchronization", () => {
  it("does not create a partial template list or detail", () => {
    const queryClient = new QueryClient();

    syncBadgeTemplateCache(queryClient, template());

    expect(
      queryClient.getQueryData(badgeKeys.byEdition("edition-1")),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(badgeKeys.detail("template-1")),
    ).toBeUndefined();
  });

  it("updates loaded templates and invalidates rendered badges", () => {
    const queryClient = new QueryClient();
    const updated = template({ name: "Updated badge" });
    const printKey = badgeKeys.print("edition-1");
    const userKey = badgeKeys.user("actor-1");
    queryClient.setQueryData(badgeKeys.byEdition("edition-1"), [template()]);
    queryClient.setQueryData(badgeKeys.detail("template-1"), template());
    queryClient.setQueryData(printKey, []);
    queryClient.setQueryData(userKey, {});

    syncBadgeTemplateCache(queryClient, updated);

    expect(queryClient.getQueryData(badgeKeys.byEdition("edition-1"))).toEqual([
      updated,
    ]);
    expect(queryClient.getQueryData(badgeKeys.detail("template-1"))).toEqual(
      updated,
    );
    expect(queryClient.getQueryState(printKey)?.isInvalidated).toBe(true);
    expect(queryClient.getQueryState(userKey)?.isInvalidated).toBe(true);
  });

  it("removes the template list item and detail", () => {
    const queryClient = new QueryClient();
    queryClient.setQueryData(badgeKeys.byEdition("edition-1"), [template()]);
    queryClient.setQueryData(badgeKeys.detail("template-1"), template());

    removeBadgeTemplateCache(queryClient, "edition-1", "template-1");

    expect(queryClient.getQueryData(badgeKeys.byEdition("edition-1"))).toEqual(
      [],
    );
    expect(
      queryClient.getQueryData(badgeKeys.detail("template-1")),
    ).toBeUndefined();
  });
});
