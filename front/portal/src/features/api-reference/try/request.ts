import type { ApiAuth, ApiOperation } from "../model";

// Turns what's been typed into the try-it fields into one request, used both
// for the snippets and for Send. Anything left empty stays a placeholder, so
// a snippet is right even before a single field is filled.

export interface TryValues {
  /** Where requests go; empty keeps `{API_URL}` and disables Send. */
  baseUrl: string;
  /** Credentials by security scheme name. */
  credentials: Record<string, string>;
  /** Path, query and header parameters by name. */
  params: Record<string, string>;
  /** The request body as typed; null means the sampled example. */
  body: string | null;
}

export interface HarRequest {
  method: string;
  url: string;
  headers: { name: string; value: string }[];
  queryString: { name: string; value: string }[];
  postData?: { mimeType: string; text: string };
}

export const PLACEHOLDER_BASE = "{API_URL}";

/** The scheme this request authenticates with: the first one with a value, else the first allowed. */
export function authFor(op: ApiOperation, schemes: readonly ApiAuth[], values: TryValues): ApiAuth | null {
  const allowed = schemes.filter((scheme) => op.auth.includes(scheme.name));
  return allowed.find((scheme) => values.credentials[scheme.name]) ?? allowed[0] ?? null;
}

export function buildRequest(op: ApiOperation, schemes: readonly ApiAuth[], values: TryValues): HarRequest {
  const base = values.baseUrl.trim().replace(/\/+$/, "") || PLACEHOLDER_BASE;
  const path = op.path.replace(/\{([^}]+)\}/g, (whole, name: string) => {
    const value = values.params[name]?.trim();
    return value ? encodeURIComponent(value) : whole;
  });

  const queryString = op.parameters
    .filter((param) => param.in === "query")
    .flatMap((param) => {
      const value = values.params[param.name]?.trim();
      if (value) return [{ name: param.name, value }];
      return param.required ? [{ name: param.name, value: `<${param.name}>` }] : [];
    });

  const headers: HarRequest["headers"] = [];
  const auth = authFor(op, schemes, values);
  if (auth) {
    const value = values.credentials[auth.name]?.trim();
    headers.push({ name: auth.header, value: `${auth.prefix}${value || auth.placeholder}` });
  }
  for (const param of op.parameters.filter((param) => param.in === "header")) {
    const value = values.params[param.name]?.trim();
    if (value || param.required) headers.push({ name: param.name, value: value || `<${param.name}>` });
  }

  const body = op.requestBody;
  const text = values.body ?? body?.example ?? null;
  if (body && text !== null) headers.push({ name: "Content-Type", value: body.contentType });

  return {
    method: op.method.toUpperCase(),
    url: `${base}${path}`,
    headers,
    queryString,
    ...(body && text !== null ? { postData: { mimeType: body.contentType, text } } : {}),
  };
}

/** Why Send can't go yet (blocker), and what will likely go wrong if it does (hint). */
export function sendCheck(op: ApiOperation, request: HarRequest, values: TryValues): { blocker: string | null; hint: string | null } {
  const base = values.baseUrl.trim().replace(/\/+$/, "");
  let blocker: string | null = null;
  if (!base) blocker = "Set the base URL to send.";
  else if (/\{[^}]+\}/.test(request.url.slice(base.length))) blocker = "Fill in the path parameters.";
  else if (request.queryString.some((query) => /^<.+>$/.test(query.value))) blocker = "Fill in the required query parameters.";
  const hint =
    op.auth.length && request.headers.some((header) => /<(token|api-key)>$/.test(header.value))
      ? "No credential: the API will answer 401."
      : null;
  return { blocker, hint };
}

export function toFetch(request: HarRequest): { url: string; init: RequestInit } {
  const query = new URLSearchParams(request.queryString.map((q) => [q.name, q.value])).toString();
  return {
    url: query ? `${request.url}?${query}` : request.url,
    init: {
      method: request.method,
      headers: Object.fromEntries(request.headers.map((header) => [header.name, header.value])),
      body: request.postData?.text,
    },
  };
}
