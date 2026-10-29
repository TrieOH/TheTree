import type { JSX } from "@solidjs/web";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { captureArtifact, expectNoDiagnostics } from "@solidjs/diagnostics";
import { createMemo, createSignal, flush } from "solid-js";
import { afterEach, describe, expect, it, vi } from "vitest";

const location = vi.hoisted(() => ({ read: () => ({ pathname: "/" }) }));
vi.mock("@tanstack/solid-router", () => ({
  useLocation: () => location.read,
  Link: (props: { to: string; class?: string; children?: JSX.Element }) => (
    <a href={props.to} class={props.class}>{props.children}</a>
  ),
}));
import { Breadcrumb } from "@/shared/ui/Breadcrumb";

afterEach(cleanup);

describe("Breadcrumb navigation identity", () => {
  it("retains unchanged ancestor links while navigating between sibling routes", async () => {
    const { artifact } = await captureArtifact(() => {
      const [pathname, setPathname] = createSignal("/admin/events/event-1/editions", { name: "breadcrumbPath" });
      location.read = createMemo(() => ({ pathname: pathname() }));
      render(() => <Breadcrumb />);
      flush();
      const adminLink = screen.getByRole("link", { name: "Admin" });
      const eventLink = screen.getByRole("link", { name: "Event-1" });
      for (const page of ["products", "editions", "programs", "editions"]) {
        setPathname(`/admin/events/event-1/${page}`);
        flush();
        expect(screen.getByRole("link", { name: "Admin" })).toBe(adminLink);
        expect(screen.getByRole("link", { name: "Event-1" })).toBe(eventLink);
      }
      expect(screen.getByText("Edições")).toBeInTheDocument();
      cleanup();
    }, { scenario: "breadcrumb-sibling-navigation" });
    expectNoDiagnostics(artifact);
  });
});
