import { createQueryClient } from "@trieoh/front-core-solid";
import { describe, expect, it, vi } from "vitest";
import { ownProfileQueryOptions, profileKeys } from "@/features/profile/api";

describe("profile navigation cache", () => {
  it("shares in-flight and fresh profile reads, and refetches after invalidation", async () => {
    const client = createQueryClient({ maxRetries: 0 });
    const profile = { actor_id: "alice", handle: "Alice" };
    const auth = { getActorProfile: vi.fn(async () => ({ success: true, data: profile })) };
    const options = ownProfileQueryOptions("alice", auth);
    try {
      await Promise.all([client.fetchQuery(options), client.fetchQuery(options)]);
      expect(await client.fetchQuery(options)).toEqual(profile);
      expect(auth.getActorProfile).toHaveBeenCalledTimes(1);
      await client.invalidateQueries({ queryKey: profileKeys.all });
      await client.fetchQuery(options);
      expect(auth.getActorProfile).toHaveBeenCalledTimes(2);
    } finally {
      client.clear();
    }
  });

  it("keeps different accounts in separate cache entries", async () => {
    const client = createQueryClient({ maxRetries: 0 });
    const auth = { getActorProfile: vi.fn(async (id: string) => ({ success: true, data: { actor_id: id } })) };
    try {
      expect(await client.fetchQuery(ownProfileQueryOptions("alice", auth))).toEqual({ actor_id: "alice" });
      expect(await client.fetchQuery(ownProfileQueryOptions("bob", auth))).toEqual({ actor_id: "bob" });
      expect(auth.getActorProfile).toHaveBeenCalledTimes(2);
    } finally {
      client.clear();
    }
  });

  it("distinguishes a missing profile from a backend failure", async () => {
    const missing = { getActorProfile: vi.fn(async () => ({ success: false, code: 404 })) };
    const failed = { getActorProfile: vi.fn(async () => ({ success: false, code: 503 })) };
    expect(await ownProfileQueryOptions("alice", missing).queryFn()).toBeNull();
    await expect(ownProfileQueryOptions("alice", failed).queryFn()).rejects.toThrow();
  });
});
