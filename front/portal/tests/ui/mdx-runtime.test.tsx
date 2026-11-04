import { render } from "@solidjs/testing-library";
import { describe, expect, it } from "vitest";
import { Fragment, jsx, jsxs } from "../../src/shared/mdx/jsx-runtime";

describe("mdx jsx-runtime", () => {
  it("renders intrinsic tags and components the way MDX calls them", () => {
    const Note = (props: { children?: unknown }) =>
      jsx("aside", { class: "note", children: props.children as never });

    const { container } = render(() =>
      jsxs(Fragment, {
        children: [
          jsx("h2", { children: "Context" }),
          jsx(Note as never, { children: jsx("p", { children: "Inside" }) }),
        ],
      }),
    );

    expect(container.querySelector("h2")).toHaveTextContent("Context");
    expect(container.querySelector("aside.note > p")).toHaveTextContent("Inside");
  });
});
