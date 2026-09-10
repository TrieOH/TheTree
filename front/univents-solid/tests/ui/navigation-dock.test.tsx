import { render, screen } from "@solidjs/testing-library";
import { beforeEach, describe, expect, it, vi } from "vitest";

const h = vi.hoisted(() => ({
  animate: vi.fn(() => ({ stop: vi.fn() })),
  setPath: (_pathname: string) => { },
}));

vi.mock("motion", async (importOriginal) => {
  const actual = await importOriginal<typeof import("motion")>();
  return { ...actual, animate: h.animate };
});

vi.mock("@tanstack/solid-router", async () => {
  const { createSignal } = await import("solid-js");
  const [location, setLocation] = createSignal({
    pathname: "/profile/marcos",
    href: "/profile/marcos",
  });
  h.setPath = (pathname: string) => setLocation({ pathname, href: pathname });

  return {
    useLocation: () => location,
    useNavigate: () => () => undefined,
  };
});

vi.mock("@trieoh/identityx-sdk-ts-solid", () => ({
  useAuth: () => ({ isAuthenticated: () => false }),
}));

vi.mock("@/features/auths/hooks/use-session-actions", () => ({
  useSessionActions: () => ({ logoutTo: async () => { } }),
}));

const { NavigationDock } = await import("@/widgets/ui/NavigationDock");

describe("NavigationDock", () => {
  beforeEach(() => {
    h.animate.mockClear();
  });

  it("reappears when leaving an immersive route and coming back", async () => {
    h.setPath("/profile/edit");
    render(() => <NavigationDock />);

    await Promise.resolve();
    expect(screen.queryByRole("navigation")).toBeNull();

    h.setPath("/profile/marcos");
    await Promise.resolve();

    expect(screen.getAllByRole("button").length).toBeGreaterThan(0);
    expect(screen.getAllByLabelText("Home").length).toBeGreaterThan(0);
  });
});
