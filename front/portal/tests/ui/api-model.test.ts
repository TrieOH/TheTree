// @vitest-environment node
import { readFileSync } from "node:fs";
import { validate } from "@scalar/openapi-parser";
import { describe, expect, it } from "vitest";
import { parse } from "yaml";
import { buildModel, type ApiOperation } from "../../mdx/api-plugin";

const SERVICES = ["identityx", "payssage", "univents", "informd"] as const;
const load = (svc: string) =>
  parse(readFileSync(new URL(`../../../../api/${svc}/api-spec.gen.yml`, import.meta.url), "utf8"));
const operations = (svc: string): ApiOperation[] => buildModel(load(svc)).tags.flatMap((tag) => tag.operations);

describe("API model, built from the generated specs", () => {
  it.each(SERVICES)("%s: the spec is valid and every operation lands under a tag", async (svc) => {
    const doc = load(svc);
    expect((await validate(doc)).valid).toBe(true);

    const model = buildModel(doc);
    const count = Object.values(doc.paths as Record<string, object>).reduce(
      (n: number, item) => n + Object.keys(item).filter((k) => ["get", "post", "put", "patch", "delete"].includes(k)).length,
      0,
    );
    expect(new Set(model.tags.flatMap((tag) => tag.operations.map((op) => `${op.method} ${op.path}`))).size).toBe(count);
  });

  it.each(SERVICES)("%s: no real host or $ref survives, and the model stays small", (svc) => {
    const json = JSON.stringify(buildModel(load(svc)));
    expect(json).not.toMatch(/https?:\/\/(api\.|localhost)/);
    expect(json).not.toContain('"$ref"');
    expect(json.length).toBeLessThan(600_000);
  });

  it("keeps component schemas as names and resolves parameters", () => {
    const members = operations("payssage").find((op) => op.id === "listOrganizationMembers")!;

    expect(members).toMatchObject({ method: "get", path: "/organizations/{organization_id}/members", auth: ["bearerAuth"] });
    expect(members.parameters[0]).toMatchObject({ name: "organization_id", in: "path", required: true });
    expect(members.responses[0].schema).toMatchObject({ properties: { items: { type: "array", items: { ref: "OrganizationMember" } } } });
    expect(members.problems).toContain("not-found");
  });

  it("samples request bodies and renders descriptions as HTML", () => {
    const withBody = operations("payssage").find((op) => op.requestBody?.contentType === "application/json")!;

    expect(JSON.parse(withBody.requestBody!.example!)).toBeTypeOf("object");
    expect(withBody.description).toMatch(/^<p>/);
  });

  it("describes auth without hosts", () => {
    const model = buildModel(load("payssage"));
    expect(model.baseUrl).toBe("{API_URL}");
    expect(model.auth.map((a) => a.summary)).toContain("Bearer token");
  });
});
