import { cleanup, render, screen } from "@solidjs/testing-library";
import { captureArtifact, expectNoDiagnostics } from "@solidjs/diagnostics";
import { PaginatedContainer } from "@trieoh/ui-solid";
import { For, createSignal, flush } from "solid-js";
import { afterEach, describe, expect, it } from "vitest";

afterEach(cleanup);

describe("PaginatedContainer diagnostics", () => {
  it("mounts measured lists without untracked reactive reads", async () => {
    const { artifact } = await captureArtifact(() => {
      const [items, setItems] = createSignal([{ id: "a", name: "Alice" }]);
      render(() => <PaginatedContainer items={items()} sortFields={[{ key: "name", label: "Nome" }]}
        renderItems={(slice) => <For each={slice}>{(item) => <p>{item.name}</p>}</For>} />);
      flush();
      expect(screen.getByText("Alice")).toBeInTheDocument();
      setItems([{ id: "a", name: "Alice" }, { id: "b", name: "Bob" }]);
      flush();
      cleanup();
    }, { scenario: "paginated-list-mount", attribution: false });
    expectNoDiagnostics(artifact);
  });
});
