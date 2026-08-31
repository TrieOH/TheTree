import { QueryClient } from "@tanstack/react-query";
import type { ActorProfile, ProfileData } from "@trieoh/identityx-sdk-ts";
import { describe, expect, it } from "vitest";
import { invalidateProfileCaches, syncActorProfileCache } from "./cache";
import { profileKeys } from "./query-keys";

const actorProfile = (actorId: string): ActorProfile => ({
  actor_id: actorId,
  handle: "old-handle",
  pfp_url: "old.jpg",
  profile: { legalName: "Old name" },
  schema_version: 1,
  outdated: false,
  updated_at: "2026-01-01T00:00:00Z",
});

describe("profile cache synchronization", () => {
  it("updates every loaded key for the edited actor", () => {
    const queryClient = new QueryClient();
    const byId = profileKeys.detail("actor-1");
    const byHandle = profileKeys.detail("old-handle");
    const other = profileKeys.detail("actor-2");
    queryClient.setQueryData(byId, actorProfile("actor-1"));
    queryClient.setQueryData(byHandle, actorProfile("actor-1"));
    queryClient.setQueryData(other, actorProfile("actor-2"));
    const profile: ProfileData = {
      legalName: "New name",
      pfpUrl: "new.jpg",
    };

    syncActorProfileCache(queryClient, "actor-1", profile, "new-handle");

    expect(queryClient.getQueryData<ActorProfile>(byId)).toMatchObject({
      handle: "new-handle",
      pfp_url: "new.jpg",
      profile: { legalName: "New name" },
    });
    expect(queryClient.getQueryData(byHandle)).toEqual(
      queryClient.getQueryData(byId),
    );
    expect(queryClient.getQueryData(other)).toEqual(actorProfile("actor-2"));
  });

  it("does not invent a profile when its detail was not loaded", () => {
    const queryClient = new QueryClient();

    syncActorProfileCache(queryClient, "actor-1", {}, "new-handle");

    expect(
      queryClient.getQueryData(profileKeys.detail("actor-1")),
    ).toBeUndefined();
  });

  it("preserves the loaded handle when it was omitted", () => {
    const queryClient = new QueryClient();
    const detail = profileKeys.detail("actor-1");
    queryClient.setQueryData(detail, actorProfile("actor-1"));

    syncActorProfileCache(queryClient, "actor-1", { legalName: "New name" });

    expect(queryClient.getQueryData<ActorProfile>(detail)?.handle).toBe(
      "old-handle",
    );
  });

  it("invalidates derived names after an edit", () => {
    const queryClient = new QueryClient();
    const certificateName = profileKeys.certificateName("actor-1");
    const displayNames = profileKeys.displayNames(["actor-1"]);
    queryClient.setQueryData(certificateName, "Old name");
    queryClient.setQueryData(displayNames, { "actor-1": "Old name" });

    syncActorProfileCache(queryClient, "actor-1", { legalName: "New name" });

    expect(queryClient.getQueryState(certificateName)?.isInvalidated).toBe(
      true,
    );
    expect(queryClient.getQueryState(displayNames)?.isInvalidated).toBe(true);
  });

  it("invalidates all profile reads after initial setup", () => {
    const queryClient = new QueryClient();
    const detail = profileKeys.detail("actor-1");
    queryClient.setQueryData(detail, actorProfile("actor-1"));

    void invalidateProfileCaches(queryClient);

    expect(queryClient.getQueryState(detail)?.isInvalidated).toBe(true);
  });
});
