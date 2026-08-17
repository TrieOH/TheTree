import { handleTracesIngest } from "@trieoh/front-core/tracing/ingest";
import { describe, expect, it, vi } from "vitest";

const env = {
  TRACES_OTLP_USER: "trace-user",
  TRACES_OTLP_PASSWORD: "trace-password",
  TRACES_OTLP_URL: "https://traces.example.test/v1/traces",
};

describe("handleTracesIngest", () => {
  it("forwards the OTLP body and authentication to the upstream", async () => {
    const upstreamFetch = vi
      .fn<typeof fetch>()
      .mockResolvedValue(new Response(null, { status: 204 }));
    vi.stubGlobal("fetch", upstreamFetch);

    const body = '{"resourceSpans":[]}';
    const response = await handleTracesIngest(
      new Request("https://app.test/api/traces/ingest", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body,
      }),
      env,
    );

    expect(response.status).toBe(204);
    expect(upstreamFetch).toHaveBeenCalledOnce();

    const [url, init] = upstreamFetch.mock.calls[0];
    expect(url).toBe(env.TRACES_OTLP_URL);
    expect(init?.method).toBe("POST");
    expect(init?.headers).toEqual({
      Authorization: `Basic ${btoa("trace-user:trace-password")}`,
      "Content-Type": "application/json",
    });
    expect(await new Response(init?.body).text()).toBe(body);
  });

  it("does not call the upstream when credentials are missing", async () => {
    const upstreamFetch = vi.fn<typeof fetch>();
    vi.stubGlobal("fetch", upstreamFetch);

    const response = await handleTracesIngest(
      new Request("https://app.test/api/traces/ingest", {
        method: "POST",
        body: "payload",
      }),
      {},
    );

    expect(response.status).toBe(204);
    expect(upstreamFetch).not.toHaveBeenCalled();
  });

  it("returns 502 when the upstream rejects the batch", async () => {
    vi.stubGlobal(
      "fetch",
      vi
        .fn<typeof fetch>()
        .mockResolvedValue(new Response("unauthorized", { status: 401 })),
    );

    const response = await handleTracesIngest(
      new Request("https://app.test/api/traces/ingest", {
        method: "POST",
        body: "payload",
      }),
      env,
    );

    expect(response.status).toBe(502);
  });
});
