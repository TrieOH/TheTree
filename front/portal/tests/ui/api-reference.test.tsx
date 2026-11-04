import { cleanup, render, waitFor } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ThemeProvider } from "../../src/shared/lib/theme";

const created: { config: Record<string, unknown>; updates: Record<string, unknown>[]; destroyed: boolean }[] = [];

// Scalar itself (Vue, several MB) is stubbed: what matters here is what the
// portal hands it.
vi.mock("../../src/features/api-reference/scalar", () => ({
  createApiReference: (_host: Element, config: Record<string, unknown>) => {
    const record = { config, updates: [] as Record<string, unknown>[], destroyed: false };
    created.push(record);
    return {
      updateConfiguration: (next: Record<string, unknown>) => record.updates.push(next),
      destroy: () => {
        record.destroyed = true;
      },
    };
  },
}));

const { ApiReference } = await import("../../src/features/api-reference/ApiReference");
const { loadSpec } = await import("../../src/features/api-reference/specs");

afterEach(() => {
  cleanup();
  created.length = 0;
});

describe("API reference", () => {
  it("reads each service's generated spec from the repo", async () => {
    const spec = await loadSpec("payssage");
    expect(spec).toMatch(/^openapi: 3/m);
    expect(spec).toContain("Payssage API");
  });

  it("gives Scalar the spec and switches off everything that calls out", async () => {
    render(() => (
      <ThemeProvider>
        <ApiReference service="payssage" />
      </ThemeProvider>
    ));
    await waitFor(() => expect(created).toHaveLength(1));

    const { config } = created[0];
    expect(config.content).toContain("Payssage API");
    expect(config).toMatchObject({
      hideClientButton: true,
      hideTestRequestButton: true,
      telemetry: false,
      withDefaultFonts: false,
      showDeveloperTools: "never",
      agent: { disabled: true },
      mcp: { disabled: true },
      forceDarkModeState: "dark",
    });
    expect(config).not.toHaveProperty("url");
    expect(config).not.toHaveProperty("proxyUrl");
  });

  it("destroys Scalar when the page goes away", async () => {
    const view = render(() => (
      <ThemeProvider>
        <ApiReference service="identityx" />
      </ThemeProvider>
    ));
    await waitFor(() => expect(created).toHaveLength(1));

    view.unmount();
    expect(created[0].destroyed).toBe(true);
  });
});
