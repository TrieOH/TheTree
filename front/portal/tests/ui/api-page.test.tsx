import { cleanup, fireEvent, render, waitFor, within } from "@solidjs/testing-library";
import { RouterProvider, createMemoryHistory, createRootRoute, createRouter } from "@tanstack/solid-router";
import { afterEach, describe, expect, it } from "vitest";
import { loadApi, typeLabel } from "../../src/features/api-reference/model";
import { OperationSection } from "../../src/features/api-reference/ui/OperationSection";

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

  it("writes an endpoint out in full, problems folded", async () => {
    const api = (await loadApi("payssage"))!;
    const op = api.tags.flatMap((tag) => tag.operations).find((o) => o.id === "listOrganizationMembers")!;
    const router = createRouter({
      routeTree: createRootRoute({
        component: () => <OperationSection op={op} service="payssage" auth={api.auth} />,
      }),
      history: createMemoryHistory({ initialEntries: ["/"] }),
    });
    const { container } = render(() => <RouterProvider router={router} />);
    const view = within(container);

    expect(await view.findByText("List organization members")).toBeInTheDocument();
    expect(view.getByText("organization_id")).toBeInTheDocument();
    expect(view.getByText(/Auth: bearerAuth/)).toBeInTheDocument();
    expect(container.querySelector("a[href='#model-OrganizationMember']")?.textContent).toBe("OrganizationMember[]");

    expect(container.querySelector("a[href='/problems/shared/not-found']")).toBeNull();
    fireEvent.click(view.getByRole("button", { name: /Problems/ }));
    await waitFor(() => expect(container.querySelector("a[href='/problems/shared/not-found']")).not.toBeNull());
  });

  it("opens a try-it overlay whose snippet follows the fields", async () => {
    const api = (await loadApi("payssage"))!;
    const op = api.tags.flatMap((tag) => tag.operations).find((o) => o.id === "listOrganizationMembers")!;
    const router = createRouter({
      routeTree: createRootRoute({
        component: () => <OperationSection op={op} service="payssage" auth={api.auth} />,
      }),
      history: createMemoryHistory({ initialEntries: ["/"] }),
    });
    const { container } = render(() => <RouterProvider router={router} />);
    const view = within(document.body);

    fireEvent.click(await within(container).findByRole("button", { name: "Try it" }));
    const dialog = await view.findByRole("dialog");
    const panel = within(dialog);

    expect(panel.getByRole("button", { name: /Send/ })).toBeDisabled();
    expect(panel.getByText("Set the base URL to send.")).toBeInTheDocument();
    await waitFor(() => expect(dialog.textContent).toContain("{API_URL}/organizations/{organization_id}/members"));

    fireEvent.input(panel.getByLabelText(/organization_id/), { target: { value: "org_9" } });
    fireEvent.input(panel.getByLabelText("Base URL"), { target: { value: "http://localhost:8080" } });
    await waitFor(() => expect(dialog.textContent).toContain("http://localhost:8080/organizations/org_9/members"));
    expect(panel.getByRole("button", { name: /Send/ })).toBeEnabled();

    fireEvent.keyDown(window, { key: "Escape" });
    await waitFor(() => expect(view.queryByRole("dialog")).toBeNull());
  });
});
