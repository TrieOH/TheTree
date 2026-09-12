import { render } from "@solidjs/testing-library";
import { beforeEach, describe, expect, it, vi } from "vitest";

type Subject = { id: string; email: string | null };

const h = vi.hoisted(() => ({
  identify: vi.fn(),
  reset: vi.fn(),
  init: vi.fn(),
  loaded: { value: false },
  setSubject: (_subject: { id: string; email: string | null } | null) => { },
}));

vi.mock("posthog-js", () => ({
  default: {
    get __loaded() {
      return h.loaded.value;
    },
    init: h.init,
    identify: h.identify,
    reset: h.reset,
    get_property: () => "actor-1",
  },
}));

vi.mock("@trieoh/identityx-sdk-ts-solid", async () => {
  const { createSignal } = await import("solid-js");
  const [subject, setSubject] = createSignal<Subject | null>({
    id: "actor-1",
    email: "maria@example.com",
  });
  h.setSubject = setSubject;

  return {
    useAuth: () => ({
      auth: { profile: subject },
      isAuthenticated: () => subject() !== null,
    }),
  };
});

const { AuthenticatedPostHogProvider } = await import("@trieoh/front-core-solid");

const provider = (key: string, capturePageview?: boolean | "history_change") =>
  render(() => (
    <AuthenticatedPostHogProvider config={{ key, capturePageview }}>
      <div>child</div>
    </AuthenticatedPostHogProvider>
  ));

describe("AuthenticatedPostHogProvider", () => {
  beforeEach(() => {
    h.identify.mockClear();
    h.reset.mockClear();
    h.init.mockClear();
    h.loaded.value = false;
    h.setSubject({ id: "actor-1", email: "maria@example.com" });
  });

  it("initializes the client and renders its children", () => {
    provider("phc_test");

    expect(h.init).toHaveBeenCalledWith(
      "phc_test",
      expect.objectContaining({
        api_host: "https://us.i.posthog.com",
        // Pageviews must stay on: Web/Product Analytics read them.
        capture_pageview: "history_change",
      }),
    );
  });

  it("lets the app turn pageviews off explicitly", () => {
    provider("phc_test", false);

    expect(h.init).toHaveBeenCalledWith(
      "phc_test",
      expect.objectContaining({ capture_pageview: false }),
    );
  });

  it("keeps the placeholder key from initializing", () => {
    provider("phc_xxx");

    expect(h.init).not.toHaveBeenCalled();
  });

  it("identifies the session subject with its email", async () => {
    h.loaded.value = true;
    provider("phc_test");

    await Promise.resolve();
    expect(h.identify).toHaveBeenCalledWith("actor-1", { email: "maria@example.com" });
  });

  it("resets the person once the session is gone", async () => {
    h.loaded.value = true;
    provider("phc_test");
    await Promise.resolve();

    h.setSubject(null);
    await Promise.resolve();

    expect(h.reset).toHaveBeenCalled();
  });
});
