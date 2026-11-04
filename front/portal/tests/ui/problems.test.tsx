import { cleanup, fireEvent, render, waitFor, within } from "@solidjs/testing-library";
import { RouterProvider, createMemoryHistory, createRootRoute, createRouter } from "@tanstack/solid-router";
import { afterEach, describe, expect, it } from "vitest";
import { exampleBody, extraDescription, findProblem, problemsOf } from "../../src/features/problems/model";
import { ProblemList, type ProblemListProps } from "../../src/features/problems/ui/ProblemList";

window.scrollTo = () => {};
afterEach(cleanup);

async function renderList(props: ProblemListProps) {
  const router = createRouter({
    routeTree: createRootRoute({ component: () => <ProblemList {...props} /> }),
    history: createMemoryHistory({ initialEntries: ["/"] }),
  });
  const result = render(() => <RouterProvider router={router} />);
  const view = within(result.container);
  await view.findAllByRole("listitem");
  return { container: result.container, view };
}

const slugsIn = (container: HTMLElement) =>
  [...container.querySelectorAll("li a[href^='/problems/']")].map((a) => a.textContent);

describe("problems model", () => {
  it("keeps shared problems apart from a service's own", () => {
    const shared = problemsOf("shared").map((p) => p.slug);
    expect(shared).toContain("not-found");
    expect(problemsOf("payssage").map((p) => p.slug)).not.toContain("not-found");
  });

  it("hides a description that only restates the title", () => {
    expect(extraDescription({ ...findProblem("identityx", "invalid-access-token")! })).toBeNull();
    expect(extraDescription(findProblem("payssage", "provider-unauthorized")!)).toContain("refused");
  });

  it("builds an example body with the problem's own members", () => {
    const body = JSON.parse(exampleBody(findProblem("payssage", "provider-unauthorized")!, "payssage"));
    expect(body).toMatchObject({
      type: "https://api.payssage.trieoh.com/problems/provider-unauthorized",
      status: 401,
      provider: "…",
      reason: "…",
    });
  });
});

describe("<ProblemList>", () => {
  it("groups by status and links the shared problems by default", async () => {
    const { container } = await renderList({ set: "payssage" });

    const headings = [...container.querySelectorAll("h4")].map((h) => h.textContent);
    expect(headings[0]).toMatch(/^\d{3}/);
    expect(new Set(headings).size).toBe(headings.length);
    expect(slugsIn(container)).not.toContain("not-found");
    expect(container.querySelector("a[href='/problems/shared']")?.textContent).toMatch(/\d+ shared problems/);
  });

  it("includes shared problems on request, marked as shared", async () => {
    const { container, view } = await renderList({ set: "payssage", shared: "include" });

    expect(slugsIn(container)).toContain("not-found");
    expect(view.getAllByText("shared").length).toBe(problemsOf("shared").length);
    expect(container.querySelector("a[href='/problems/shared'][class*='underline']")).toBeNull();
  });

  it("shows only the chosen slugs, flat, with details forced open", async () => {
    const { container, view } = await renderList({
      set: "payssage",
      only: ["provider-unauthorized", "provider-conflict"],
      group: false,
      open: true,
    });

    expect(slugsIn(container)).toEqual(["provider-unauthorized", "provider-conflict"]);
    expect(container.querySelectorAll("h4")).toHaveLength(0);
    expect(view.queryByRole("button", { name: "Show details" })).toBeNull();
    expect(container.textContent).toContain('"reason"');
  });

  it("expands a row's details on click", async () => {
    const { container, view } = await renderList({ set: "payssage", only: ["provider-unauthorized"] });

    expect(container.textContent).not.toContain('"reason"');
    fireEvent.click(view.getByRole("button", { name: "Show details" }));
    // CodeBlock highlights a tick after it mounts.
    await waitFor(() => expect(container.textContent).toContain('"reason"'));
  });

  it("puts descriptions inline only when asked", async () => {
    const page = await renderList({ set: "payssage", only: ["provider-unauthorized"] });
    expect(page.container.textContent).not.toContain("refused");
    cleanup();

    const inline = await renderList({ set: "payssage", only: ["provider-unauthorized"], description: "inline" });
    expect(inline.container.textContent).toContain("refused Payssage's credentials");
  });
});
