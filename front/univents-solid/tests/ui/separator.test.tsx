import { cleanup, render, screen } from "@solidjs/testing-library";
import { afterEach, describe, expect, it } from "vitest";

import { Separator } from "@trieoh/ui-solid";

afterEach(cleanup);

describe("Separator", () => {
  it("is a decorative horizontal hairline by default", () => {
    const { container } = render(() => <Separator />);
    const line = container.firstElementChild!;

    expect(line).toHaveAttribute("role", "none");
    expect(line).not.toHaveAttribute("aria-orientation");
    expect(line.className).toContain("border-t");
    expect(line.className).toContain("w-full");
    expect(line.className).toContain("border-solid");
  });

  it("is announced as a separator when it is not decorative", () => {
    render(() => <Separator decorative={false} orientation="vertical" />);
    const line = screen.getByRole("separator");

    expect(line).toHaveAttribute("aria-orientation", "vertical");
    expect(line.className).toContain("border-l");
    expect(line.className).toContain("self-stretch");
  });

  it.each([
    ["dashed", "border-dashed"],
    ["dotted", "border-dotted"],
  ] as const)("supports the %s style", (variant, expected) => {
    const { container } = render(() => <Separator variant={variant} />);

    expect(container.firstElementChild!.className).toContain(expected);
  });

  it("draws a line on both sides of a centred label", () => {
    const { container } = render(() => <Separator label="ou" />);
    const root = container.firstElementChild!;

    expect(screen.getByText("ou")).toBeInTheDocument();
    expect(root.querySelectorAll("span[aria-hidden='true']")).toHaveLength(2);
  });

  it("drops the line before a start label and after an end label", () => {
    const start = render(() => <Separator label="início" labelPosition="start" />);
    expect(
      start.container.firstElementChild!.querySelectorAll("span[aria-hidden='true']"),
    ).toHaveLength(1);
    // The only line comes after the label.
    expect(start.container.firstElementChild!.lastElementChild).toHaveAttribute(
      "aria-hidden",
      "true",
    );
    cleanup();

    const end = render(() => <Separator label="fim" labelPosition="end" />);
    expect(end.container.firstElementChild!.firstElementChild).toHaveAttribute(
      "aria-hidden",
      "true",
    );
    expect(
      end.container.firstElementChild!.querySelectorAll("span[aria-hidden='true']"),
    ).toHaveLength(1);
  });

  it("stacks a vertical label between two lines", () => {
    const { container } = render(() => (
      <Separator orientation="vertical" label="ou" />
    ));
    const root = container.firstElementChild!;

    expect(root.className).toContain("flex-col");
    expect(root.querySelector("span[aria-hidden='true']")!.className).toContain(
      "border-l",
    );
  });

  it("merges a custom class", () => {
    const { container } = render(() => <Separator class="my-4" />);

    expect(container.firstElementChild!.className).toContain("my-4");
  });
});
