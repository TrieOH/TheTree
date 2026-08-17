import { browserTracingResource } from "@trieoh/front-core/tracing/browser";
import { describe, expect, it } from "vitest";

describe("browser tracing resource", () => {
  it("sets the configured service.name", () => {
    const resource = browserTracingResource("univents-web");

    expect(resource.attributes).toMatchObject({
      "service.name": "univents-web",
    });
  });
});
