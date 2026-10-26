import type { JSX } from "@solidjs/web";
import { captureArtifact, expectNoDiagnostics } from "@solidjs/diagnostics";
import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { QueryClient, TanStackQueryProvider, useQuery } from "@trieoh/front-core-solid";
import { createSignal, flush } from "solid-js";
import { afterEach, describe, expect, it, vi } from "vitest";
import { serviceHealth } from "@/shared/lib/service-health";
import {
  QueryBoundary,
  combineQueries,
  type QueryLike,
} from "@/shared/ui/QueryBoundary";

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

const failure = (code: number, error_id?: string) =>
  Object.assign(new Error("failed"), { envelope: { code, error_id } });

function setup(initial: Partial<QueryLike<string>> = {}, extra: { notFound?: JSX.Element } = {}) {
  const refetch = vi.fn();
  const [query, setQuery] = createSignal<QueryLike<string>>({
    data: undefined,
    error: null,
    isFetching: true,
    errorUpdateCount: 0,
    refetch,
    ...initial,
  });
  const update = (next: Partial<QueryLike<string>>) => {
    setQuery((current) => ({ ...current, ...next }));
    flush();
  };

  render(() => (
    <QueryBoundary
      query={query}
      what="os dados"
      skeleton={<p>skeleton</p>}
      notFound={extra.notFound}
    >
      {(data) => <p>content:{data()}</p>}
    </QueryBoundary>
  ));
  flush();
  return { update, refetch };
}

describe("QueryBoundary", () => {
  it("shows the skeleton while there is nothing yet, never a blank or an empty state", () => {
    setup();
    expect(screen.getByText("skeleton")).toBeInTheDocument();
  });

  it("swaps the skeleton for the content when the data lands", () => {
    const { update } = setup();
    update({ data: "ready", isFetching: false });
    expect(screen.getByText("content:ready")).toBeInTheDocument();
    expect(screen.queryByText("skeleton")).not.toBeInTheDocument();
  });

  it("explains a failure inside its own region and offers a retry", () => {
    const { update, refetch } = setup();
    update({ error: failure(503), errorUpdateCount: 1, isFetching: false });

    expect(screen.getByText("Não foi possível carregar os dados")).toBeInTheDocument();
    expect(screen.getByText("O serviço está temporariamente indisponível.")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Tentar novamente" }));
    expect(refetch).toHaveBeenCalledTimes(1);
    expect(invalidate).toHaveBeenCalledTimes(1);
  });

  it("keeps the error panel up while a retry runs instead of flashing the skeleton", () => {
    const { update } = setup();
    update({ error: failure(503), errorUpdateCount: 1, isFetching: false });
    // A refetch clears the error and goes back to "pending" while it runs.
    update({ error: null, isFetching: true });

    expect(screen.queryByText("skeleton")).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Tentando novamente…" })).toBeDisabled();
    // ...and it doesn't fall back to an "unexpected error" message meanwhile.
    expect(screen.queryByText("Ocorreu um erro inesperado.")).not.toBeInTheDocument();
    expect(screen.getByText("Não foi possível carregar os dados")).toBeInTheDocument();
  });

  it("keeps showing the content when a background refresh fails", () => {
    const { update } = setup({ data: "ready", isFetching: false });
    update({ error: failure(503), errorUpdateCount: 1 });
    expect(screen.getByText("content:ready")).toBeInTheDocument();
    expect(screen.queryByText("Não foi possível carregar os dados")).not.toBeInTheDocument();
  });

  it("lets an explicit retry through an open circuit", () => {
    const failed = { code: 503, error_id: "CLIENT_REQUEST_FAILED", message: "", success: false as const };
    for (let i = 0; i < 3; i++) {
      serviceHealth.failure("api", {
        kind: "unavailable",
        transient: true,
        message: failed.message,
      });
    }
    expect(serviceHealth.isOpen("api")).toBe(true);

    const { update } = setup();
    update({ error: failure(503), errorUpdateCount: 1, isFetching: false });
    fireEvent.click(screen.getByRole("button", { name: "Tentar novamente" }));

    expect(serviceHealth.isOpen("api")).toBe(false);
  });

  it("renders the not-found content for a 404, with no pointless retry", () => {
    const { update } = setup({}, { notFound: <p>não existe</p> });
    update({ error: failure(404), errorUpdateCount: 1, isFetching: false });

    expect(screen.getByText("não existe")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Tentar novamente" })).not.toBeInTheDocument();
  });

  it("offers sign-in instead of retry when the session is gone", () => {
    const { update } = setup();
    update({ error: failure(401), errorUpdateCount: 1, isFetching: false });

    expect(screen.getByRole("link", { name: "Entrar" })).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Tentar novamente" })).not.toBeInTheDocument();
  });

  it("does not raise reactivity diagnostics across its states", async () => {
    const { artifact } = await captureArtifact(() => {
      const { update } = setup();
      update({ data: "ready", isFetching: false });
      update({ error: failure(503), errorUpdateCount: 1 });
      cleanup();
    }, { scenario: "query-boundary-states" });
    expectNoDiagnostics(artifact);
  });
});

describe("serviceHealth.onRecover", () => {
  it("fires onRecover when a service that had failures succeeds", () => {
    const onRecover = vi.fn();
    serviceHealth.onRecover(onRecover);

    // Initial success with 0 failures does not fire recovery
    serviceHealth.success("api");
    expect(onRecover).not.toHaveBeenCalled();

    // Now record a failure
    serviceHealth.failure("api", {
      kind: "unavailable",
      transient: true,
      message: "outage",
    });

    // When the service answers again, onRecover fires
    serviceHealth.success("api");
    expect(onRecover).toHaveBeenCalledWith("api");
    expect(onRecover).toHaveBeenCalledTimes(1);
  });
});

describe("combineQueries", () => {
  const client = () =>
    new QueryClient({ defaultOptions: { queries: { retry: false, staleTime: Infinity } } });

  it("has data only once every query has", async () => {
    let finishSecond!: (value: string) => void;

    function Probe() {
      const combined = combineQueries(
        useQuery({ queryKey: ["first"], queryFn: async () => "a" }),
        useQuery({
          queryKey: ["second"],
          queryFn: () => new Promise<string>((resolve) => { finishSecond = resolve; }),
        }),
      );
      return <p data-testid="out">{JSON.stringify(combined().data ?? null)}</p>;
    }

    render(() => (
      <TanStackQueryProvider client={client()}>
        <Probe />
      </TanStackQueryProvider>
    ));

    // The first has landed, the second hasn't: still no data, so no half-built page.
    await waitFor(() => expect(finishSecond).toBeDefined());
    await Promise.resolve();
    expect(screen.getByTestId("out")).toHaveTextContent("null");

    finishSecond("b");
    await waitFor(() => expect(screen.getByTestId("out")).toHaveTextContent('["a","b"]'));
  });

  it("surfaces the first failure and re-asks all of them on retry", async () => {
    const first = vi.fn(async () => "a");
    const second = vi.fn(async () => { throw failure(503); });
    let retry!: () => unknown;

    function Probe() {
      const combined = combineQueries(
        useQuery({ queryKey: ["first"], queryFn: first }),
        useQuery({ queryKey: ["second"], queryFn: second }),
      );
      retry = () => combined().refetch();
      return <p data-testid="out">{combined().errorUpdateCount}</p>;
    }

    render(() => (
      <TanStackQueryProvider client={client()}>
        <Probe />
      </TanStackQueryProvider>
    ));
    await waitFor(() => expect(screen.getByTestId("out")).toHaveTextContent("1"));

    retry();
    await waitFor(() => expect(first).toHaveBeenCalledTimes(2));
    expect(second).toHaveBeenCalledTimes(2);
  });
});
