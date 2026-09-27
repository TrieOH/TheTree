import type { AuthTokens } from "@trieoh/identityx-sdk-ts";
import type { BffSession } from "./core";

const encoder = new TextEncoder();
const decoder = new TextDecoder();

/** Imported keys are cached per password: importing per request is pure cost. */
const signingKeys = new Map<string, Promise<CryptoKey>>();

function signingKey(password: string): Promise<CryptoKey> {
  let key = signingKeys.get(password);
  if (!key) {
    key = crypto.subtle.importKey(
      "raw",
      encoder.encode(password),
      { name: "HMAC", hash: "SHA-256" },
      false,
      ["sign", "verify"],
    );
    signingKeys.set(password, key);
  }
  return key;
}

function toBase64Url(bytes: Uint8Array): string {
  let binary = "";
  for (const byte of bytes) binary += String.fromCharCode(byte);
  return btoa(binary).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

// No explicit return type on purpose: callers pass this straight to
// `crypto.subtle.verify`, which requires `Uint8Array<ArrayBuffer>`.
function fromBase64Url(value: string) {
  const base64 = value
    .replace(/-/g, "+")
    .replace(/_/g, "/")
    .padEnd(Math.ceil(value.length / 4) * 4, "=");
  const binary = atob(base64);
  const bytes = new Uint8Array(binary.length);
  for (let index = 0; index < binary.length; index += 1) {
    bytes[index] = binary.charCodeAt(index);
  }
  return bytes;
}

async function seal(password: string, value: string): Promise<string> {
  const payload = toBase64Url(encoder.encode(value));
  const signature = await crypto.subtle.sign(
    "HMAC",
    await signingKey(password),
    encoder.encode(payload),
  );
  return `${payload}.${toBase64Url(new Uint8Array(signature))}`;
}

async function open(password: string, sealed: string): Promise<string | null> {
  const separator = sealed.lastIndexOf(".");
  if (separator <= 0) return null;

  const payload = sealed.slice(0, separator);
  const signature = sealed.slice(separator + 1);

  try {
    const valid = await crypto.subtle.verify(
      "HMAC",
      await signingKey(password),
      fromBase64Url(signature),
      encoder.encode(payload),
    );
    return valid ? decoder.decode(fromBase64Url(payload)) : null;
  } catch {
    return null;
  }
}

function readCookie(request: Request, name: string): string | null {
  const header = request.headers.get("Cookie");
  if (!header) return null;

  for (const part of header.split(";")) {
    const separator = part.indexOf("=");
    if (separator === -1) continue;
    if (part.slice(0, separator).trim() !== name) continue;
    return part.slice(separator + 1).trim() || null;
  }
  return null;
}

function serializeCookie(
  name: string,
  value: string,
  options: {
    maxAge: number;
    secure: boolean;
    sameSite: "Lax" | "Strict" | "None";
  },
): string {
  const attributes = [
    `${name}=${value}`,
    "Path=/",
    "HttpOnly",
    `SameSite=${options.sameSite}`,
    `Max-Age=${options.maxAge}`,
  ];
  if (options.secure) attributes.push("Secure");
  return attributes.join("; ");
}

export interface CookieSessionOptions {
  /** Stable random secret of at least 32 characters. */
  password: string;
  name?: string;
  maxAge?: number;
  secure?: boolean;
  sameSite?: "Lax" | "Strict" | "None";
}

export interface CookieBffSession {
  /** Handed to `createIdentityXBff`, which only sees read/write/clear. */
  session: BffSession;
  /** Attaches the `Set-Cookie` header when the session was written or cleared. */
  applyTo(response: Response): Response;
}

/**
 * Stateless `BffSession` backed by a signed cookie: the tokens travel sealed by
 * HMAC-SHA256, so the server keeps no session storage. Rotating the password
 * invalidates every session.
 */
export async function createCookieBffSession(
  request: Request,
  options: CookieSessionOptions,
): Promise<CookieBffSession> {
  if (options.password.length < 32) {
    throw new Error("IdentityX session password must contain at least 32 characters");
  }

  const name = options.name ?? "trieoh-auth";
  const maxAge = options.maxAge ?? 60 * 60 * 24 * 30;
  const secure = options.secure ?? true;
  const sameSite = options.sameSite ?? "Lax";

  const stored = readCookie(request, name);
  let data: { tokens?: AuthTokens } | null = null;

  if (stored) {
    const opened = await open(options.password, stored);
    if (opened) {
      try {
        data = JSON.parse(opened) as { tokens?: AuthTokens };
      } catch {
        data = null;
      }
    }
  }

  /** `undefined` = untouched, `null` = cleared, string = sealed value. */
  let staged: string | null | undefined;

  return {
    session: {
      read: () => data,
      write: async (tokens) => {
        data = { tokens };
        staged = await seal(options.password, JSON.stringify({ tokens }));
      },
      clear: () => {
        data = null;
        staged = null;
      },
    },
    applyTo(response) {
      if (staged === undefined) return response;

      const headers = new Headers(response.headers);
      headers.append(
        "Set-Cookie",
        serializeCookie(name, staged ?? "", {
          maxAge: staged === null ? 0 : maxAge,
          secure,
          sameSite,
        }),
      );
      return new Response(response.body, {
        status: response.status,
        statusText: response.statusText,
        headers,
      });
    },
  };
}
