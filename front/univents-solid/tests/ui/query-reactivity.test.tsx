import { cleanup, render, screen, waitFor } from "@solidjs/testing-library";
import { captureArtifact, expectNoDiagnostics } from "@solidjs/diagnostics";
import { QueryClient, TanStackQueryProvider, useQuery } from "@trieoh/front-core-solid";
import { createEffect, createSignal, flush } from "solid-js";
import { afterEach, describe, expect, it } from "vitest";

afterEach(cleanup);

describe("useQuery reactive options", () => {
  it("never renders a new profile key with the previous profile data", async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    for (const id of ["alice", "bob"]) client.setQueryData(["profile", id], { id });
    const frames: Array<readonly [string, string | undefined]> = [];
    const { artifact } = await captureArtifact(() => {
      const [actorId, setActorId] = createSignal("alice", { name: "profileActorId" });
      function Consumer() {
        const query = useQuery(() => ({
          queryKey: ["profile", actorId()],
          queryFn: async () => ({ id: "unused" }),
          enabled: false,
        }));
        createEffect(() => [actorId(), query().data?.id] as const, (pair) => { frames.push(pair); }, { name: "profileFrame" });
        return <p data-testid="profile">{actorId()}:{query().data?.id}</p>;
      }
      render(() => <TanStackQueryProvider client={client}><Consumer /></TanStackQueryProvider>);
      flush();
      for (const id of ["bob", "alice", "bob", "alice"]) {
        setActorId(id);
        flush();
        expect(screen.getByTestId("profile")).toHaveTextContent(`${id}:${id}`);
      }
      cleanup();
    }, { scenario: "cached-profile-key-navigation" });
    client.clear();
    expect(frames).toEqual([
      ["alice", "alice"], ["bob", "bob"], ["alice", "alice"],
      ["bob", "bob"], ["alice", "alice"],
    ]);
    expectNoDiagnostics(artifact);
  });

  it("shows loading immediately for an uncached key and delivers the fetch result", async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false, staleTime: Infinity } } });
    client.setQueryData(["profile", "alice"], { id: "alice" });
    let finish!: (data: { id: string }) => void;
    const pending = new Promise<{ id: string }>((resolve) => { finish = resolve; });
    const frames: Array<{ actor: string; data?: string; loading: boolean }> = [];
    const { artifact } = await captureArtifact(async () => {
      const [actor, setActor] = createSignal("alice", { name: "uncachedProfileActor" });
      function Consumer() {
        const query = useQuery(() => ({ queryKey: ["profile", actor()], queryFn: () => pending }));
        createEffect(() => ({ actor: actor(), data: query().data?.id, loading: query().isLoading }),
          (frame) => { frames.push(frame); }, { name: "uncachedProfileFrame" });
        return <p data-testid="async-profile">{query().isLoading ? "Carregando" : query().data?.id}</p>;
      }
      render(() => <TanStackQueryProvider client={client}><Consumer /></TanStackQueryProvider>);
      flush();
      setActor("bob");
      flush();
      expect(screen.getByTestId("async-profile")).toHaveTextContent("Carregando");
      finish({ id: "bob" });
      await waitFor(() => expect(screen.getByTestId("async-profile")).toHaveTextContent("bob"));
      cleanup();
    }, { scenario: "uncached-profile-loading" });
    client.clear();
    expect(frames.filter((frame) => frame.actor === "bob" && !frame.data).every((frame) => frame.loading)).toBe(true);
    expectNoDiagnostics(artifact);
  });

  it("receives cache updates and releases the observer on unmount", async () => {
    const client = new QueryClient();
    client.setQueryData(["profile", "alice"], { name: "Alice" });
    function Consumer() {
      const query = useQuery({ queryKey: ["profile", "alice"], queryFn: async () => ({ name: "unused" }), enabled: false });
      return <p>{query().data?.name}</p>;
    }
    const { artifact } = await captureArtifact(async () => {
      const view = render(() => <TanStackQueryProvider client={client}><Consumer /></TanStackQueryProvider>);
      client.setQueryData(["profile", "alice"], { name: "Alice atualizada" });
      await screen.findByText("Alice atualizada");
      view.unmount();
    }, { scenario: "query-cache-update-cleanup" });
    expect(client.getQueryCache().find({ queryKey: ["profile", "alice"] })?.getObserversCount()).toBe(0);
    client.clear();
    expectNoDiagnostics(artifact);
  });
});
