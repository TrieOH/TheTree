import { describe, expect, it, vi } from "vitest";
import type { AuthService } from "@trieoh/identityx-sdk-ts-solid";
import {
  createBffAuthAdapter,
  type BffTransport,
} from "@trieoh/front-core/auth/bff/client";

const calls: Array<{ op: string; input: unknown }> = [];

const transport: BffTransport = {
  async call<T>(op: string, input?: unknown): Promise<T> {
    calls.push({ op, input });
    return {
      success: true,
      code: 200,
      data: { handle: op, actor_id: "actor-1" },
    } as T;
  },
};

/**
 * In BFF mode the browser has no token, so any SDK method left un-proxied
 * silently turns into a direct IdentityX call (wrong origin, no session). Every
 * reachable method must therefore go through `call("request", …)`.
 */
const directCall = (name: string) =>
  vi.fn(() => {
    throw new Error(`${name} must not be called directly in BFF mode`);
  });

function createAuth() {
  calls.length = 0;
  const defaultAuth = {
    getProfileByHandle: directCall("getProfileByHandle"),
    getOAuthProviders: directCall("getOAuthProviders"),
    health: directCall("health"),
    getActorProfile: directCall("getActorProfile"),
  } as unknown as AuthService;

  const adapter = createBffAuthAdapter(transport, { projectId: "project-1" });
  const auth = adapter.createAuth({
    callbacks: {},
    defaultAuth,
    getProfile: () => null,
    setProfile: () => undefined,
    setAuthenticated: () => undefined,
  });

  return { auth, defaultAuth };
}

describe("BFF auth adapter", () => {
  it("proxies getProfileByHandle instead of calling IdentityX directly", async () => {
    const { auth, defaultAuth } = createAuth();

    const result = await auth.getProfileByHandle("Humano");

    expect(result.success).toBe(true);
    expect(defaultAuth.getProfileByHandle).not.toHaveBeenCalled();
    expect(calls).toHaveLength(1);
    expect(calls[0]).toMatchObject({
      op: "request",
      input: {
        path: "/profiles/by-handle/Humano",
        target: "identityx",
        method: "GET",
      },
    });
  });

  it("proxies the OAuth provider discovery with the project id", async () => {
    const { auth, defaultAuth } = createAuth();

    await auth.getOAuthProviders();

    expect(defaultAuth.getOAuthProviders).not.toHaveBeenCalled();
    expect(calls[0]?.input).toMatchObject({
      path: "/auth/oauth-providers?project_id=project-1",
    });
  });

  it("proxies actor profiles through the project path", async () => {
    const { auth, defaultAuth } = createAuth();

    await auth.getActorProfile("actor-1");

    expect(defaultAuth.getActorProfile).not.toHaveBeenCalled();
    expect(calls[0]?.input).toMatchObject({
      path: "/projects/project-1/actors/actor-1/profile",
    });
  });
});
