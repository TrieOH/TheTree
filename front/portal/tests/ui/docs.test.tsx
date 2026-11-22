import { render } from "@solidjs/testing-library";
import { describe, expect, it } from "vitest";
import { docsIn, findDoc, neighboursOf } from "../../src/features/docs/catalog";

describe("docs catalog", () => {
  it("collects every ADR from the repo's docs/adr, newest first", () => {
    const decisions = docsIn("decisions");

    expect(decisions.length).toBeGreaterThanOrEqual(19);
    expect(decisions[0].slug > decisions[1].slug).toBe(true);
    expect(decisions.every((doc) => doc.meta.number?.startsWith("ADR-"))).toBe(true);
  });

  it("renders an ADR through the MDX pipeline without its preamble", async () => {
    const doc = findDoc("decisions", "adr-0019");
    expect(doc?.meta).toMatchObject({
      number: "ADR-0019",
      status: "accepted",
      date: "2026-10-29",
      amends: ["ADR-0016", "ADR-0017"],
    });
    expect(doc?.meta.title).toBeTruthy();
    expect(doc?.meta.summary).toBeTruthy();

    const { default: Content } = await doc!.load();
    const { container } = render(() => <Content />);

    expect(container.querySelector("h2")).toHaveTextContent("Why");
    expect(container.querySelector("code")).not.toBeNull();
    expect(container).not.toHaveTextContent("number: ADR-0019");
    expect(container).not.toHaveTextContent("status: accepted");
  });
});

describe("doc neighbours", () => {
  it("links each doc to the ones beside it in its category", () => {
    const decisions = docsIn("decisions");
    const middle = decisions[1];

    expect(neighboursOf(middle)).toEqual({
      previous: { category: "decisions", slug: decisions[0].slug, meta: decisions[0].meta },
      next: { category: "decisions", slug: decisions[2].slug, meta: decisions[2].meta },
    });
    expect(neighboursOf(decisions[0]).previous).toBeUndefined();
    expect(neighboursOf(decisions.at(-1)!).next).toBeUndefined();
  });
});
