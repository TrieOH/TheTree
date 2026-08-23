import { beforeEach, describe, expect, it, vi } from "vitest";

const session = {
  data: {
    tokens: {
      access_token: "expired.access.token",
      refresh_token: "refresh.token.value",
      access_expires_at: "2000-01-01T00:00:00.000Z",
      refresh_expires_at: "2100-01-01T00:00:00.000Z",
    },
  },
  update: vi.fn(async (data: typeof session.data) => {
    session.data = data;
  }),
  clear: vi.fn(async () => {
    session.data = undefined as never;
  }),
};

vi.mock("@tanstack/react-start/server", () => ({
  getRequest: () => new Request("https://univents.test/_server"),
  useSession: () => Promise.resolve(session),
}));

import { createTanStackIdentityXBff } from "@trieoh/front-core/auth/tanstack/server";

describe("IdentityX BFF token refresh", () => {
  beforeEach(() => {
    session.data = {
      tokens: {
        access_token: "expired.access.token",
        refresh_token: "refresh.token.value",
        access_expires_at: "2000-01-01T00:00:00.000Z",
        refresh_expires_at: "2100-01-01T00:00:00.000Z",
      },
    };
    vi.clearAllMocks();
  });

  it("shares a concurrent refresh and authenticates every API request", async () => {
    let refreshCalls = 0;
    const apiAuthorizations: Array<string | null> = [];
    vi.stubGlobal(
      "fetch",
      vi.fn(async (input: string | URL | Request, init?: RequestInit) => {
        const url = String(input);
        if (url.endsWith("/auth/refresh")) {
          refreshCalls += 1;
          const call = refreshCalls;
          await new Promise((resolve) => setTimeout(resolve, 5));
          if (call > 1) {
            return Response.json(
              { code: 401, message: "refresh token already rotated" },
              { status: 401 },
            );
          }
          return Response.json({
            code: 200,
            data: {
              access_token: "fresh.access.token",
              refresh_token: "fresh.refresh.token",
              access_expires_at: "2100-01-01T00:00:00.000Z",
              refresh_expires_at: "2100-01-01T00:00:00.000Z",
            },
          });
        }

        apiAuthorizations.push(new Headers(init?.headers).get("Authorization"));
        return Response.json({ code: 200, data: [] });
      }),
    );

    const bff = createTanStackIdentityXBff({
      identityX: { baseURL: "https://identityx.test" },
      session: { password: "a".repeat(32) },
      apiBaseURL: "https://univents-api.test",
    });

    const results = await Promise.all([
      bff.request({ path: "/events/joined" }),
      bff.request({ path: "/events/joined" }),
    ]);

    expect(results.every((result) => result.success)).toBe(true);
    expect(refreshCalls).toBe(1);
    expect(apiAuthorizations).toEqual([
      "Bearer fresh.access.token",
      "Bearer fresh.access.token",
    ]);
  });
});
