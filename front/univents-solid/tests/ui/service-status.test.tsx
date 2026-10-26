import type { JSX } from "@solidjs/web";
import { captureArtifact, expectNoDiagnostics } from "@solidjs/diagnostics";
import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { flush } from "solid-js";
import { afterEach, describe, expect, it, vi } from "vitest";
import { classifyError } from "@/shared/lib/service-error";
import { serviceHealth } from "@/shared/lib/service-health";
import { NotFoundPage, RouteError } from "@/widgets/ui/RouteStatus";
import { ServiceStatusBanner } from "@/widgets/ui/ServiceStatusBanner";

const invalidate = vi.fn();

vi.mock("@tanstack/solid-router", () => ({
  Link: (props: { to: string; children?: JSX.Element }) => (
    <a href={props.to}>{props.children}</a>
  ),
  useRouter: () => ({ invalidate }),
}));

afterEach(() => {
  cleanup();
  serviceHealth.reset();
  invalidate.mockClear();
});

const outage = classifyError(
  Object.assign(new Error("x"), { envelope: { code: 503 } }),
);

describe("ServiceStatusBanner", () => {
  it("stays out of the way while services are healthy", () => {
    render(() => <ServiceStatusBanner />);
    expect(screen.queryByRole("status")).not.toBeInTheDocument();
  });

  it("appears for the service that is struggling, with an accurate title", () => {
    render(() => <ServiceStatusBanner />);
    serviceHealth.failure("identityx", outage);
    serviceHealth.failure("identityx", outage);
    flush();

    const indicator = screen.getByRole("status");
    expect(indicator).toBeInTheDocument();
    expect(indicator).toHaveAttribute("title", expect.stringMatching(/sessão desconectada/i));
  });

  it("goes away on its own once the service answers again", () => {
    render(() => <ServiceStatusBanner />);
    serviceHealth.failure("api", outage);
    serviceHealth.failure("api", outage);
    flush();
    expect(screen.getByRole("status")).toBeInTheDocument();

    serviceHealth.success("api");
    flush();
    expect(screen.queryByRole("status")).not.toBeInTheDocument();
  });

  it("clicking the indicator lets requests through an open circuit", () => {
    render(() => <ServiceStatusBanner />);
    for (let i = 0; i < 3; i++) serviceHealth.failure("api", outage);
    flush();
    expect(serviceHealth.isOpen("api")).toBe(true);

    fireEvent.click(screen.getByRole("status"));
    expect(serviceHealth.isOpen("api")).toBe(false);
  });

  it("does not raise reactivity diagnostics as the health changes", async () => {
    const { artifact } = await captureArtifact(() => {
      render(() => <ServiceStatusBanner />);
      serviceHealth.failure("api", outage);
      serviceHealth.failure("api", outage);
      flush();
      serviceHealth.success("api");
      flush();
      cleanup();
    }, { scenario: "service-status-banner" });
    expectNoDiagnostics(artifact);
  });
});

describe("RouteError", () => {
  const reset = vi.fn();

  it("reads a backend failure as 'try again', without a raw message", () => {
    render(() => (
      <RouteError
        error={Object.assign(new Error("fetch exploded"), { envelope: { code: 503 } })}
        reset={reset}
        info={undefined}
      />
    ));
    expect(screen.getByRole("heading", { name: "Não foi possível abrir esta página" })).toBeInTheDocument();
    expect(screen.queryByText("fetch exploded")).not.toBeInTheDocument();
  });

  it("calls a code bug a bug", () => {
    render(() => (
      <RouteError error={new TypeError("x is not a function")} reset={reset} info={undefined} />
    ));
    expect(screen.getByRole("heading", { name: "Algo deu errado" })).toBeInTheDocument();
  });

  it("retry resets the boundary and reloads the route's data", () => {
    render(() => (
      <RouteError error={new Error("boom")} reset={reset} info={undefined} />
    ));
    fireEvent.click(screen.getByRole("button", { name: "Tentar novamente" }));
    expect(reset).toHaveBeenCalledTimes(1);
    expect(invalidate).toHaveBeenCalledTimes(1);
  });
});

describe("NotFoundPage", () => {
  it("is a real page with a way home, not a template placeholder", () => {
    render(() => <NotFoundPage />);
    expect(screen.getByRole("heading", { name: "Página não encontrada" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Ir para o início" })).toBeInTheDocument();
  });
});
