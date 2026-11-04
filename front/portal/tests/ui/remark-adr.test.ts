import { compile } from "@mdx-js/mdx";
import { describe, expect, it } from "vitest";
import { remarkAdr } from "../../mdx/remark-adr";

async function metaOf(source: string) {
  const code = String(await compile(source, { format: "md", remarkPlugins: [remarkAdr] }));
  const match = /export const meta = (\{.*?\});/s.exec(code);
  return { meta: match ? Function(`return ${match[1]}`)() : undefined, code };
}

describe("remarkAdr", () => {
  it("lifts the ADR number, title and fields out of the body", async () => {
    const { meta, code } = await metaOf(
      [
        "# ADR-0019 — Errors are RFC 9457 problems",
        "",
        "- **Status:** accepted",
        "- **Date:** 2026-10-29",
        "- **Scope:** `api/identityx`",
        "",
        "## Context",
        "",
        "Body text.",
      ].join("\n"),
    );

    expect(meta).toEqual({
      title: "Errors are RFC 9457 problems",
      number: "ADR-0019",
      status: "accepted",
      date: "2026-10-29",
      scope: "api/identityx",
    });
    expect(code).not.toContain("ADR-0019 —");
    expect(code).toContain("Context");
  });

  it("keeps a list that isn't a field list", async () => {
    const { meta, code } = await metaOf("# Runbook\n\n- step one\n- step two\n");

    expect(meta).toEqual({ title: "Runbook" });
    expect(code).toContain("step one");
  });

  it("exports an empty meta for a doc without a header", async () => {
    const { meta } = await metaOf("Just text.\n");

    expect(meta).toEqual({});
  });
});
