import { describe, expect, it } from "vitest";
import type { ApiAuth, ApiOperation } from "../../src/features/api-reference/model";
import { buildRequest, sendCheck, toFetch, type TryValues } from "../../src/features/api-reference/try/request";
import { snippet } from "../../src/features/api-reference/try/snippets";

const bearer: ApiAuth = { name: "bearerAuth", summary: "Bearer token", description: "", header: "Authorization", prefix: "Bearer ", placeholder: "<token>" };
const apiKey: ApiAuth = { name: "apiKeyAuth", summary: "API key", description: "", header: "X-API-KEY", prefix: "", placeholder: "<api-key>" };

const op: ApiOperation = {
  id: "createSeller",
  method: "post",
  path: "/wallets/{wallet_id}/sellers",
  summary: "Create a seller",
  description: "",
  deprecated: false,
  auth: ["bearerAuth", "apiKeyAuth"],
  scope: null,
  parameters: [
    { name: "wallet_id", in: "path", required: true, description: "", schema: null },
    { name: "limit", in: "query", required: true, description: "", schema: null },
    { name: "cursor", in: "query", required: false, description: "", schema: null },
  ],
  requestBody: { contentType: "application/json", schema: null, example: '{\n  "name": "string"\n}', required: true },
  responses: [],
  problems: [],
};

const empty: TryValues = { baseUrl: "", credentials: {}, params: {}, body: null };

describe("try-it request", () => {
  it("keeps placeholders for everything not filled in", () => {
    const request = buildRequest(op, [bearer, apiKey], empty);

    expect(request.url).toBe("{API_URL}/wallets/{wallet_id}/sellers");
    expect(request.queryString).toEqual([{ name: "limit", value: "<limit>" }]);
    expect(request.headers).toContainEqual({ name: "Authorization", value: "Bearer <token>" });
    expect(request.postData?.text).toContain('"name"');
    expect(sendCheck(op, request, empty).blocker).toBe("Set the base URL to send.");
  });

  it("fills in what's typed, and uses the scheme that has a credential", () => {
    const values: TryValues = {
      baseUrl: "http://localhost:8080/",
      credentials: { apiKeyAuth: "k_123" },
      params: { wallet_id: "w 1", limit: "10" },
      body: '{"name":"Ana"}',
    };
    const request = buildRequest(op, [bearer, apiKey], values);

    expect(request.url).toBe("http://localhost:8080/wallets/w%201/sellers");
    expect(request.headers).toContainEqual({ name: "X-API-KEY", value: "k_123" });
    expect(sendCheck(op, request, values)).toEqual({ blocker: null, hint: null });
    expect(toFetch(request)).toMatchObject({
      url: "http://localhost:8080/wallets/w%201/sellers?limit=10",
      init: { method: "POST", body: '{"name":"Ana"}' },
    });
  });

  it("lets a request go without a credential, with a hint", () => {
    const values: TryValues = { ...empty, baseUrl: "http://localhost:8080", params: { wallet_id: "w", limit: "1" } };
    const request = buildRequest(op, [bearer, apiKey], values);

    expect(sendCheck(op, request, values)).toEqual({ blocker: null, hint: "No credential: the API will answer 401." });
  });

  it("blocks Send until path parameters are filled", () => {
    const values: TryValues = { ...empty, baseUrl: "http://localhost:8080", params: { limit: "1" } };
    expect(sendCheck(op, buildRequest(op, [bearer], values), values).blocker).toBe("Fill in the path parameters.");
  });

  it("turns the request into code with placeholders intact", () => {
    const curl = snippet("curl", buildRequest(op, [bearer], empty));
    expect(curl).toContain("{API_URL}/wallets/{wallet_id}/sellers");
    expect(curl).toContain("Bearer <token>");
    expect(snippet("python", buildRequest(op, [bearer], empty))).toContain("requests.post(");
  });
});
