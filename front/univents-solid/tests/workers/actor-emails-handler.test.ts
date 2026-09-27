import { beforeEach, describe, expect, it, vi } from "vitest";
import { clearActorCache, handleActorEmailsRequest } from "@/features/events/api/actor-emails-handler";

const mockEnv = {
  VITE_AUTH_API_URL: "https://identityx.test",
  VITE_TRIEOH_AUTH_PROJECT_ID: "proj-123",
  IDENTITYX_ACCESS_API_KEY: "secret-key-123",
} as unknown as Env;

describe("handleActorEmailsRequest", () => {
  beforeEach(() => {
    clearActorCache();
  });

  it("rejects non-POST methods", async () => {
    const req = new Request("http://localhost/api/actors/emails", { method: "GET" });
    const res = await handleActorEmailsRequest(req, mockEnv);
    expect(res.status).toBe(405);
  });

  it("returns an empty map for empty actorIds array", async () => {
    const req = new Request("http://localhost/api/actors/emails", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ actorIds: [] }),
    });
    const res = await handleActorEmailsRequest(req, mockEnv);
    expect(res.status).toBe(200);
    const data = await res.json();
    expect(data).toEqual({});
  });

  it("fetches actor emails and profile picture from IdentityX and caches results", async () => {
    const globalFetch = vi.fn().mockImplementation((url: string) => {
      if (url.includes("actor-1/profile")) {
        return Promise.resolve(
          new Response(
            JSON.stringify({
              code: 200,
              data: { actor_id: "actor-1", pfp_url: "https://example.com/alice.jpg" },
            }),
            { status: 200 },
          ),
        );
      }
      if (url.includes("actor-1")) {
        return Promise.resolve(
          new Response(
            JSON.stringify({
              code: 200,
              data: { id: "actor-1", email: "alice@example.com" },
            }),
            { status: 200 },
          ),
        );
      }
      return Promise.resolve(new Response(null, { status: 404 }));
    });
    vi.stubGlobal("fetch", globalFetch);

    const req = new Request("http://localhost/api/actors/emails", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ actorIds: ["actor-1", "actor-2"] }),
    });

    const res = await handleActorEmailsRequest(req, mockEnv);
    expect(res.status).toBe(200);
    const data = await res.json();
    expect(data).toEqual({
      "actor-1": {
        email: "alice@example.com",
        pfp_url: "https://example.com/alice.jpg",
      },
    });

    // Second call should hit the in-memory cache and not call globalFetch again for actor-1
    const fetchCountBefore = globalFetch.mock.calls.length;
    const req2 = new Request("http://localhost/api/actors/emails", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ actorIds: ["actor-1"] }),
    });
    const res2 = await handleActorEmailsRequest(req2, mockEnv);
    expect(res2.status).toBe(200);
    const data2 = await res2.json();
    expect(data2).toEqual({
      "actor-1": {
        email: "alice@example.com",
        pfp_url: "https://example.com/alice.jpg",
      },
    });
    expect(globalFetch.mock.calls.length).toBe(fetchCountBefore);

    vi.unstubAllGlobals();
  });
});
