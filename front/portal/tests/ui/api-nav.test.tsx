import { cleanup, fireEvent, render, waitFor, within } from "@solidjs/testing-library";
import { afterEach, describe, expect, it } from "vitest";
import { loadApi } from "../../src/features/api-reference/model";
import { ApiNav } from "../../src/features/api-reference/ui/ApiNav";

afterEach(cleanup);

describe("API route list", () => {
  it("highlights the group being read without opening it, and the route once opened", async () => {
    const api = (await loadApi("payssage"))!;
    const tag = api.tags.find((t) => t.name === "orgs")!;
    const active = tag.operations[0];
    const { container } = render(() => <ApiNav api={api} active={active.id} />);
    const view = within(container);

    const group = container.querySelector<HTMLElement>(`[data-nav-tag="orgs"]`)!;
    expect(group.getAttribute("aria-current")).toBe("location");
    // Closed: no routes listed, nothing opened on its own.
    expect(container.querySelector(`[data-nav-op="${active.id}"]`)).toBeNull();
    expect(container.querySelectorAll("[aria-current]")).toHaveLength(1);

    fireEvent.click(view.getByRole("button", { name: "Open orgs" }));
    await waitFor(() =>
      expect(container.querySelector(`[data-nav-op="${active.id}"]`)?.getAttribute("aria-current")).toBe("location"),
    );
    // Open, the route carries the location instead of the group.
    expect(group.getAttribute("aria-current")).toBeNull();
    // Each route shows its method, colored per method.
    expect(container.querySelector(`[data-nav-op="${active.id}"] span`)?.className).toMatch(/text-(ok|accent|warn|destructive)/);
  });
});
