import type { RouterSession } from "@trieoh/front-core-solid";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { requireConfiguredProfile } from "@/features/auths/lib/route-guard";
import { profileKeys } from "@/features/profile/api";
import { appQueryClient } from "@/shared/lib/query-client";

const ACTOR = "11111111-1111-4111-8111-111111111111";

const session = (getActorProfile: () => Promise<unknown>): RouterSession =>
  ({
    isAuthenticated: true,
    service: {
      profile: () => ({ id: ACTOR }),
      getActorProfile,
    },
  }) as unknown as RouterSession;

const args = (context: { session?: RouterSession }, pathname = "/events") => ({
  context,
  location: { href: pathname, pathname, search: {} },
});

beforeEach(() => {
  appQueryClient.clear();
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
});

describe("requireConfiguredProfile", () => {
  it("does nothing for visitors", async () => {
    await expect(requireConfiguredProfile(args({}))).resolves.toBeUndefined();
  });

  it("sends a signed-in user with no profile to setup", async () => {
    const guard = requireConfiguredProfile(
      args({ session: session(async () => ({ success: false, code: 404 })) }),
    );
    // Attach the assertion before timers run so the rejection is never unhandled.
    const assertion = expect(guard).rejects.toMatchObject({
      options: { to: "/profile/setup" },
    });
    await vi.advanceTimersByTimeAsync(2_000);
    await assertion;
  });

  it("fails open when the profile service is down, instead of taking every page with it", async () => {
    const getActorProfile = vi.fn(async () => ({ success: false, code: 503 }));
    const guard = requireConfiguredProfile(args({ session: session(getActorProfile) }));
    await vi.advanceTimersByTimeAsync(2_000);
    await expect(guard).resolves.toBeUndefined();
  });

  it("fails open when the profile service hangs, within the wait bound", async () => {
    const guard = requireConfiguredProfile(
      args({ session: session(() => new Promise(() => undefined)) }),
    );
    await vi.advanceTimersByTimeAsync(1_200);
    await expect(guard).resolves.toBeUndefined();
  });

  it("answers from the cache without waiting on the network", async () => {
    const getActorProfile = vi.fn(() => new Promise(() => undefined));
    appQueryClient.setQueryData(profileKeys.detail(ACTOR), { actor_id: ACTOR, handle: "a" });

    // No timers advanced: it must already be resolved.
    await expect(
      requireConfiguredProfile(args({ session: session(getActorProfile) })),
    ).resolves.toBeUndefined();
  });

  it("redirects away from setup once the profile is known to exist", async () => {
    appQueryClient.setQueryData(profileKeys.detail(ACTOR), { actor_id: ACTOR, handle: "a" });
    await expect(
      requireConfiguredProfile(
        args({ session: session(async () => ({})) }, "/profile/setup"),
      ),
    ).rejects.toMatchObject({ options: { to: "/profile" } });
  });
});
