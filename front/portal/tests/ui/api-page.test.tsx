import { cleanup, fireEvent, render, within } from "@solidjs/testing-library";
import { RouterProvider, createMemoryHistory, createRootRoute, createRouter } from "@tanstack/solid-router";
import { afterEach, describe, expect, it } from "vitest";
import { loadApi, typeLabel } from "../../src/features/api-reference/model";
import { OperationCard } from "../../src/features/api-reference/ui/OperationCard";

window.scrollTo = () => {};
afterEach(cleanup);

describe("typeLabel", () => {
  it("reads like TypeScript", () => {
    expect(typeLabel({ ref: "Wallet" })).toBe("Wallet");
    expect(typeLabel({ type: "array", items: { ref: "Wallet" } })).toBe("Wallet[]");
    expect(typeLabel({ type: "string", format: "uuid", nullable: true })).toBe("string (uuid) | null");
    expect(typeLabel({ enum: ["open", "closed"] })).toBe('"open" | "closed"');
    expect(typeLabel({ type: "object", additionalProperties: { type: "string" } })).toBe("Record<string, string>");
  });
});

describe("API reference", () => {
  it("loads a service's model built from its spec", async () => {
    const api = await loadApi("payssage");

    expect(api?.title).toBe("Payssage API");
    expect(api?.tags.map((tag) => tag.name)).toContain("wallets");
    expect(api?.baseUrl).toBe("{API_URL}");
  });

  it("shows an endpoint collapsed, then everything about it", async () => {
    const api = (await loadApi("payssage"))!;
    const op = api.tags.flatMap((tag) => tag.operations).find((o) => o.id === "listOrganizationMembers")!;
    const router = createRouter({
      routeTree: createRootRoute({ component: () => <OperationCard op={op} service="payssage" /> }),
      history: createMemoryHistory({ initialEntries: ["/"] }),
    });
    const { container } = render(() => <RouterProvider router={router} />);
    const view = within(container);
    const header = await view.findByRole("button", { expanded: false });

    expect(header.textContent).toContain("/organizations/{organization_id}/members");
    expect(view.queryByText("Parameters")).toBeNull();

    fireEvent.click(header);
    // Solid 2 applies the update on the next tick.
    expect(await view.findByText("Parameters")).toBeInTheDocument();
    expect(view.getByText("organization_id")).toBeInTheDocument();
    expect(view.getByText(/Auth: bearerAuth/)).toBeInTheDocument();
    // Its problems, shared ones included, through the problems list.
    expect(container.querySelector("a[href='/problems/shared/not-found']")).not.toBeNull();
    // The response's items link to their model.
    expect(container.querySelector("a[href='#model-OrganizationMember']")?.textContent).toBe("OrganizationMember[]");
  });
});
