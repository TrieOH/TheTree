import { render } from "@solidjs/testing-library";
import { describe, expect, it } from "vitest";
import { docsIn, findDoc } from "../../src/features/docs/catalog";

describe("docs catalog", () => {
  it("collects every ADR from the repo's docs/adr, newest first", () => {
    const decisions = docsIn("decisions");

    expect(decisions.length).toBeGreaterThanOrEqual(19);
    expect(decisions[0].slug > decisions[1].slug).toBe(true);
    expect(decisions.every((doc) => doc.meta.number?.startsWith("ADR-"))).toBe(true);
  });

  it("renders an ADR through the MDX pipeline without its header", async () => {
    const doc = findDoc("decisions", "adr-0019");
    expect(doc?.meta).toMatchObject({ number: "ADR-0019", status: "accepted", date: "2026-10-29" });

    const { default: Content } = await doc!.load();
    const { container } = render(() => <Content />);

    expect(container.querySelector("h2")).toHaveTextContent("Context");
    expect(container.querySelector("code")).not.toBeNull();
    expect(container).not.toHaveTextContent("ADR-0019 —");
    expect(container).not.toHaveTextContent("Status:");
  });
});
