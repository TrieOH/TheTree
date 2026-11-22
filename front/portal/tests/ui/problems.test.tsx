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

  it("reads the problems a feature declares in its own spec.yml", () => {
    expect(findProblem("identityx", "setup-complete")?.status).toBe(409);
    expect(problemsOf("identityx").length).toBeGreaterThan(37);
  });

  it("hides a description that only restates the title", () => {
    expect(extraDescription({ ...findProblem("identityx", "invalid-access-token")! })).toBeNull();
    expect(extraDescription(findProblem("payssage", "provider-unauthorized")!)).toContain("refused");
  });

  it("carries a named error's detail as the example detail", () => {
    const problem = findProblem("identityx", "platform-only")!;
    expect(problem.detail).toMatch(/^this operation manages the IdentityX platform/);
    expect(JSON.parse(exampleBody(problem)).detail).toBe(problem.detail);
  });

  it("builds an example body with the problem's own members", () => {
    const body = JSON.parse(exampleBody(findProblem("payssage", "provider-unauthorized")!));
    expect(body).toMatchObject({
      // A placeholder, never a real host someone could paste and wake.
      type: "{API_URL}/problems/provider-unauthorized",
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
    expect(container.textContent).toContain("This problem adds");
    expect(container.textContent).toContain("{API_URL}/problems/provider-unauthorized");
    expect(container.textContent).not.toContain("trieoh.com");
  });

  it("expands a row's details on click", async () => {
    const { container, view } = await renderList({ set: "payssage", only: ["provider-unauthorized"] });

    expect(container.textContent).not.toContain("This problem adds");
    fireEvent.click(view.getByRole("button", { name: "Show details" }));
    await waitFor(() => expect(container.textContent).toContain("This problem adds"));
    expect(view.getByRole("button", { name: "Copy example body" })).toBeInTheDocument();
  });

  it("shows a named error's detail, and says when the detail varies", async () => {
    const named = await renderList({ set: "identityx", only: ["platform-only"], open: true });
    expect(named.container.textContent).toContain(findProblem("identityx", "platform-only")!.detail);
    expect(named.container.textContent).not.toContain("varies");
    cleanup();

    const varies = await renderList({ set: "shared", only: ["validation"], open: true });
    expect(varies.container.textContent).toContain("varies: the text of the error that matched");
  });

  it("puts descriptions inline only when asked", async () => {
    const page = await renderList({ set: "payssage", only: ["provider-unauthorized"] });
    expect(page.container.textContent).not.toContain("refused");
    cleanup();

    const inline = await renderList({ set: "payssage", only: ["provider-unauthorized"], description: "inline" });
    expect(inline.container.textContent).toContain("refused Payssage's credentials");
  });
});
