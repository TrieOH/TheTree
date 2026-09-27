import type { IdentityXBff } from "./core";
import type { ServerProxyRequest, SerializableValue } from "./types";

/** Invalid `{ op, input }` payload — a client mistake, answered with a 400. */
export class BffInputError extends Error { }

function asRecord(input: unknown): Record<string, unknown> {
  if (!input || typeof input !== "object" || Array.isArray(input)) {
    throw new BffInputError("Expected an object");
  }
  return input as Record<string, unknown>;
}

function asString(input: unknown, field: string): string {
  const value = asRecord(input)[field];
  if (typeof value !== "string" || value.length === 0) {
    throw new BffInputError(`Expected a non-empty string at "${field}"`);
  }
  return value;
}

function credentials(input: unknown): { email: string; password: string } {
  if (typeof asRecord(input).password !== "string") {
    throw new BffInputError(`Expected a string at "password"`);
  }
  return {
    email: asString(input, "email"),
    password: asRecord(input).password as string,
  };
}

function provider(input: unknown): "github" | "google" {
  const value = asString(input, "provider");
  if (value !== "github" && value !== "google") {
    throw new BffInputError(`Unsupported provider "${value}"`);
  }
  return value;
}

function proxyRequest(input: unknown): ServerProxyRequest {
  const record = asRecord(input);
  const request: ServerProxyRequest = { path: asString(input, "path") };

  if (record.target !== undefined) {
    if (record.target !== "api" && record.target !== "identityx") {
      throw new BffInputError(`Unsupported target "${String(record.target)}"`);
    }
    request.target = record.target;
  }

  if (record.method !== undefined) {
    if (!["GET", "POST", "PUT", "PATCH", "DELETE"].includes(String(record.method))) {
      throw new BffInputError(`Unsupported method "${String(record.method)}"`);
    }
    request.method = record.method as ServerProxyRequest["method"];
  }

  if (record.body !== undefined && typeof record.body === "function") {
    throw new BffInputError("Body must be serializable");
  }
  if (record.body !== undefined) {
    request.body = record.body as SerializableValue;
  }

  if (record.headers !== undefined) {
    if (typeof record.headers !== "object" || record.headers === null) {
      throw new BffInputError("Headers must be an object");
    }
    for (const value of Object.values(record.headers as Record<string, unknown>)) {
      if (typeof value !== "string") {
        throw new BffInputError("Header values must be strings");
      }
    }
    request.headers = record.headers as Record<string, string>;
  }

  return request;
}

/**
 * The BFF surface, defined once. Every transport (TanStack Start server
 * functions, an HTTP endpoint, a test) dispatches from here instead of
 * redeclaring the operations and their payload shapes.
 */
export const bffOperations: Record<
  string,
  (bff: IdentityXBff, input: unknown) => Promise<unknown>
> = {
  isSetupDone: (bff) => bff.isSetupDone(),
  setup: (bff, input) => {
    const { email, password } = credentials(input);
    return bff.setup(email, password);
  },
  login: (bff, input) => {
    const { email, password } = credentials(input);
    return bff.login(email, password);
  },
  register: (bff, input) => {
    const { email, password } = credentials(input);
    return bff.register(email, password);
  },
  loginWithProvider: (bff, input) => bff.loginWithProvider(provider(input)),
  completeProviderLogin: (bff, input) =>
    bff.completeProviderLogin(
      provider(input),
      asString(input, "code"),
      asString(input, "state"),
    ),
  logout: (bff) => bff.logout(),
  refresh: (bff) => bff.refresh(),
  restore: (bff) => bff.restore(),
  introspect: (bff, input) => {
    const apiKey = asRecord(input ?? {}).apiKey;
    if (apiKey !== undefined && typeof apiKey !== "string") {
      throw new BffInputError("apiKey must be a string");
    }
    return bff.introspect(apiKey);
  },
  request: (bff, input) => bff.request(proxyRequest(input)),
};

export function runBffOperation(
  bff: IdentityXBff,
  operation: string,
  input: unknown,
): Promise<unknown> {
  const run = bffOperations[operation];
  if (!run) throw new BffInputError(`Unknown operation "${operation}"`);
  return run(bff, input);
}
