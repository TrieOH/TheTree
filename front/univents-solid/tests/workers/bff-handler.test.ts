import { describe, expect, it, vi } from "vitest";
import { createBffHandler } from "@trieoh/front-core/auth/bff/handler";

const PASSWORD = "b".repeat(40);
const BFF_PATH = "/auth/bff";

const handler = createBffHandler({
  identityX: { baseURL: "https://identityx.test", projectId: "project-1" },
  apiBaseURL: "https://api.test",
  session: { name: "univents-solid-auth", password: PASSWORD, secure: true },
});

const jwt = (payload: Record<string, unknown>) =>
  [
    Buffer.from(JSON.stringify({ alg: "HS256" })).toString("base64url"),
    Buffer.from(JSON.stringify(payload)).toString("base64url"),
    "signature",
  ].join(".");

const tokens = {
  access_token: jwt({
    subject: { id: "actor-1" },
    exp: Math.floor(Date.now() / 1000) + 3600,
  }),
  refresh_token: "refresh.token.value",
  access_expires_at: new Date(Date.now() + 3600_000).toISOString(),
  refresh_expires_at: new Date(Date.now() + 86_400_000).toISOString(),
};

/** Outgoing IdentityX calls, so tests can assert on what the BFF sends. */
const calls: Array<{ url: string; traceparent: string | null }> = [];

const call = (op: string, input?: unknown, cookie?: string, traceparent?: string) =>
  handler(
    new Request(`https://app.test${BFF_PATH}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(cookie ? { Cookie: cookie } : {}),
        ...(traceparent ? { traceparent } : {}),
      },
      body: JSON.stringify({ op, input }),
    }),
  );

/** IdentityX stand-in: login hands out tokens, everything else echoes. */
const stubIdentityX = () => {
  const spy = vi.fn(async (input: string | URL | Request, init?: RequestInit) => {
    const url = String(input);
    const outgoing = new Headers(init?.headers).get("traceparent");
    calls.push({ url, traceparent: outgoing });
    if (url.includes("/auth/login")) {
      return Response.json({ code: 200, data: tokens });
    }
    if (url.includes("/auth/introspect")) {
      return Response.json({ code: 200, data: { cred: { type: "token" } } });
    }
    return Response.json({ code: 404, message: `unexpected ${url}` }, { status: 404 });
  });
  calls.length = 0;
  vi.stubGlobal("fetch", spy);
  return spy;
};

describe("BFF handler", () => {
  it("rejects an unknown operation without touching IdentityX", async () => {
    const spy = stubIdentityX();

    const response = await call("nope", {});

    expect(response.status).toBe(400);
    await expect(response.json()).resolves.toMatchObject({ code: 400 });
    expect(spy).not.toHaveBeenCalled();
  });

  it("rejects a malformed input instead of guessing", async () => {
    stubIdentityX();

    const response = await call("login", { email: "someone@example.test" });

    expect(response.status).toBe(400);
    await expect(response.json()).resolves.toMatchObject({
      message: expect.stringContaining("password"),
    });
  });

  it("issues a signed session cookie on login and trusts it afterwards", async () => {
    stubIdentityX();

    const login = await call("login", {
      email: "someone@example.test",
      password: "secret",
    });
    const setCookie = login.headers.get("Set-Cookie") ?? "";

    expect(login.status).toBe(200);
    expect(setCookie).toContain("univents-solid-auth=");
    expect(setCookie).toContain("HttpOnly");
    // The tokens must not be readable from the cookie itself.
    expect(setCookie).not.toContain(tokens.access_token);

    const cookie = setCookie.split(";")[0];
    const restore = await call("restore", undefined, cookie);

    await expect(restore.json()).resolves.toMatchObject({
      isAuthenticated: true,
      profile: { id: "actor-1" },
    });
  });

  it("treats a tampered cookie as no session", async () => {
    stubIdentityX();

    const login = await call("login", {
      email: "someone@example.test",
      password: "secret",
    });
    const [name, value] = (login.headers.get("Set-Cookie") ?? "").split(";")[0].split("=");
    const forged = `${name}=${value.slice(0, -4)}AAAA`;

    const restore = await call("restore", undefined, forged);

    await expect(restore.json()).resolves.toMatchObject({ isAuthenticated: false });
  });

  it("clears the cookie on logout", async () => {
    stubIdentityX();

    const login = await call("login", {
      email: "someone@example.test",
      password: "secret",
    });
    const cookie = (login.headers.get("Set-Cookie") ?? "").split(";")[0];

    const logout = await call("logout", undefined, cookie);

    expect(logout.headers.get("Set-Cookie")).toContain("Max-Age=0");
  });

  it("carries the caller's trace to IdentityX", async () => {
    stubIdentityX();
    const traceparent = "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01";

    await call("login", { email: "someone@example.test", password: "secret" }, undefined, traceparent);

    expect(calls.map((call) => call.url).some((url) => url.includes("/auth/login"))).toBe(true);
    for (const outgoing of calls) {
      expect(outgoing.traceparent).toBe(traceparent);
    }
  });

  it("refreshes an expired session", async () => {
    const expired = {
      ...tokens,
      access_expires_at: new Date(Date.now() - 1_000).toISOString(),
    };
    const rotated = {
      access_token: jwt({
        subject: { id: "actor-1" },
        exp: Math.floor(Date.now() / 1000) + 3600,
      }),
      refresh_token: "rotated.refresh.token",
      access_expires_at: new Date(Date.now() + 3_600_000).toISOString(),
      refresh_expires_at: new Date(Date.now() + 86_400_000).toISOString(),
    };

    let refreshCalls = 0;
    vi.stubGlobal(
      "fetch",
      vi.fn(async (input: string | URL | Request) => {
        const url = String(input);
        if (url.includes("/auth/login")) return Response.json({ code: 200, data: expired });
        if (url.includes("/auth/refresh")) {
          refreshCalls += 1;
          return Response.json({ code: 200, data: rotated });
        }
        return Response.json({ code: 404, message: `unexpected ${url}` }, { status: 404 });
      }),
    );

    const login = await call("login", {
      email: "someone@example.test",
      password: "secret",
    });
    const cookie = (login.headers.get("Set-Cookie") ?? "").split(";")[0];

    const restore = await call("restore", undefined, cookie);

    expect(refreshCalls).toBe(1);
    await expect(restore.json()).resolves.toMatchObject({ isAuthenticated: true });
    // The rotated session must come back to the browser, not just sit server-side.
    expect(restore.headers.get("Set-Cookie")).toContain("univents-solid-auth=");
  });

  it("shares a single refresh between concurrent requests", async () => {
    const expired = {
      ...tokens,
      access_expires_at: new Date(Date.now() - 1_000).toISOString(),
    };
    const rotated = {
      ...tokens,
      access_token: jwt({
        subject: { id: "actor-1" },
        exp: Math.floor(Date.now() / 1000) + 3600,
      }),
      refresh_token: "rotated.refresh.token",
    };

    let refreshCalls = 0;
    vi.stubGlobal(
      "fetch",
      vi.fn(async (input: string | URL | Request) => {
        const url = String(input);
        if (url.includes("/auth/login")) return Response.json({ code: 200, data: expired });
        if (url.includes("/auth/refresh")) {
          refreshCalls += 1;
          // Widen the window so both requests are genuinely in flight.
          await new Promise((resolve) => setTimeout(resolve, 10));
          return Response.json({ code: 200, data: rotated });
        }
        return Response.json({ code: 404, message: `unexpected ${url}` }, { status: 404 });
      }),
    );

    const login = await call("login", {
      email: "someone@example.test",
      password: "secret",
    });
    const cookie = (login.headers.get("Set-Cookie") ?? "").split(";")[0];

    const [first, second] = await Promise.all([
      call("restore", undefined, cookie),
      call("restore", undefined, cookie),
    ]);

    // Single-use refresh token: a second call would log the user out.
    expect(refreshCalls).toBe(1);
    await expect(first.json()).resolves.toMatchObject({ isAuthenticated: true });
    await expect(second.json()).resolves.toMatchObject({ isAuthenticated: true });
    expect(first.headers.get("Set-Cookie")).toContain("univents-solid-auth=");
    expect(second.headers.get("Set-Cookie")).toContain("univents-solid-auth=");
  });
});
