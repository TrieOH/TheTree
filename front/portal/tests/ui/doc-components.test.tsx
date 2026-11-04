import { cleanup, render, within } from "@solidjs/testing-library";
import { RouterProvider, createMemoryHistory, createRootRoute, createRouter } from "@tanstack/solid-router";
import { afterEach, describe, expect, it } from "vitest";
import { findDoc } from "../../src/features/docs/catalog";
import { docComponents } from "../../src/features/docs/components";
import { ThemeProvider } from "../../src/shared/lib/theme";

// The router restores scroll on navigation; jsdom has no layout.
window.scrollTo = () => {};
afterEach(cleanup);

async function renderDoc(category: string, slug: string) {
  const doc = findDoc(category, slug);
  expect(doc).toBeDefined();
  const { default: Content } = await doc!.load();
  const router = createRouter({
    routeTree: createRootRoute({
      component: () => (
        <ThemeProvider>
          <article class="prose-doc">
            <Content components={docComponents} />
          </article>
        </ThemeProvider>
      ),
    }),
    history: createMemoryHistory({ initialEntries: ["/"] }),
  });
  const result = render(() => <RouterProvider router={router} />);
  const view = within(result.container);
  await view.findByText("Shape of a doc");
  return { meta: doc!.meta, container: result.container, view };
}

describe("doc components, through the real MDX pipeline", () => {
  it("lifts the header and gives sections linkable ids", async () => {
    const { meta, container } = await renderDoc("onboarding", "writing-docs");

    expect(meta).toMatchObject({ title: "Writing docs", summary: expect.stringContaining("every component") });
    expect(container.querySelector("h2#shape-of-a-doc a[href='#shape-of-a-doc']")).not.toBeNull();
    expect(container.querySelector("table")).not.toBeNull();
  });

  it("highlights code at build time and adds a copy button", async () => {
    const { container } = await renderDoc("onboarding", "writing-docs");

    const pre = container.querySelector("pre.shiki");
    expect(pre?.querySelector("span[style*='--shiki-dark']")).not.toBeNull();
    expect(pre?.parentElement?.querySelector("button[aria-label='Copy code']")).not.toBeNull();
  });

  it("renders callouts, cards, steps and tabs", async () => {
    const { container, view } = await renderDoc("onboarding", "writing-docs");

    expect(container.querySelectorAll("aside")).toHaveLength(4);
    expect(view.getByText("Wakes the API")).toBeInTheDocument();
    expect(view.getByText("Decisions").closest("a")?.getAttribute("href")).toBe("/docs/decisions");
    expect(view.getByText("TheTree on GitHub").closest("a")?.getAttribute("target")).toBe("_blank");
    expect(container.querySelectorAll(".doc-steps h3")).toHaveLength(3);
    expect(view.getAllByRole("tab").map((tab) => tab.textContent)).toEqual(["API", "Front"]);
  });

  it("turns a mermaid fence into a diagram and lists a service's problems", async () => {
    const { container } = await renderDoc("onboarding", "writing-docs");

    // jsdom can't draw the diagram itself; what matters here is that the
    // fence became the Mermaid component, not highlighted code.
    expect(container.querySelector("figure [role='img'][aria-label='Diagram']")).not.toBeNull();
    const highlighted = [...container.querySelectorAll("pre.shiki")].map((pre) => pre.textContent ?? "");
    expect(highlighted.some((code) => code.includes("flowchart"))).toBe(false);
    const slugs = [...container.querySelectorAll("tbody code")].map((code) => code.textContent);
    expect(slugs.length).toBeGreaterThan(40);
    expect(slugs).toContain("not-found");
  });
});
