import { afterEach, beforeEach, expect, it, vi } from "vitest";
import { render } from "@solidjs/web";
import { createStore } from "solid-js";
import {
  BulkActionsExample,
  DataTableExample,
  FilterBarExample,
  InfiniteScrollExample,
  PaginatedExample,
  PaginationExample,
  ToolbarExample,
} from "../src/examples/CollectionExamples";
import { collectionUsage } from "../src/components/collection-usage";
import { transformWithOxc } from "vite";
import type { Values } from "../src/catalog";

const defaults: Values = {
  page: "1",
  pageSize: "10",
  totalItems: "240",
  mode: "numbered",
  search: "",
  sort: "asc",
  layout: "list",
  disabled: false,
  busy: false,
  automatic: false,
  simulateError: false,
  dataMode: "No servidor (simulado)",
  filterValue: "",
  gap: "3",
  minItemWidth: "12rem",
  itemLabel: "projetos",
  filterPlaceholder: "Buscar projetos…",
  selectable: true,
  rowActions: true,
  loading: false,
};
let host: HTMLDivElement;
let dispose: (() => void) | undefined;
let warnings: ReturnType<typeof vi.spyOn>;
beforeEach(() => {
  host = document.createElement("div");
  document.body.append(host);
  warnings = vi.spyOn(console, "warn");
  vi.stubGlobal(
    "ResizeObserver",
    class {
      observe() {}
      disconnect() {}
    },
  );
});
afterEach(() => {
  dispose?.();
  dispose = undefined;
  host.remove();
  expect(warnings.mock.calls.flat().join(" ")).not.toMatch(
    /STRICT_READ_UNTRACKED|REACTIVE_WRITE_IN_OWNED_SCOPE|UNTRACKED_READ_AFTER_AWAIT/,
  );
  vi.unstubAllGlobals();
});
it.each([
  PaginationExample,
  ToolbarExample,
  FilterBarExample,
  BulkActionsExample,
  InfiniteScrollExample,
  DataTableExample,
])("renders %s with reactive properties", async (Example) => {
  const [values, setValues] = createStore({ ...defaults });
  dispose = render(
    () => (
      <Example
        values={values}
        update={(key, value) =>
          setValues((draft) => {
            draft[key] = value;
          })
        }
      />
    ),
    host,
  );
  await new Promise((resolve) => setTimeout(resolve, 20));
  expect(host.children.length).toBeGreaterThan(0);
});
it("server pagination example delivers only the requested page", async () => {
  const [values, setValues] = createStore({ ...defaults });
  dispose = render(
    () => (
      <PaginatedExample
        values={values}
        update={(key, value) =>
          setValues((draft) => {
            draft[key] = value;
          })
        }
      />
    ),
    host,
  );
  await new Promise((resolve) => setTimeout(resolve, 20));
  expect(host.querySelector('[role="status"]')?.textContent).toContain(
    "Carregando",
  );
  await new Promise((resolve) => setTimeout(resolve, 450));
  expect(host.querySelector('[aria-current="page"]')?.textContent).toBe("1");
  (
    host.querySelector('[aria-label="Próxima página"]') as HTMLButtonElement
  ).click();
  await new Promise((resolve) => setTimeout(resolve, 450));
  expect(host.querySelector('[aria-current="page"]')?.textContent).toBe("2");
});
it.each([
  "Pagination",
  "CollectionToolbar",
  "FilterBar",
  "BulkActions",
  "InfiniteScroll",
  "DataTable",
  "PaginatedContainer",
])("generates syntactically valid copyable code for %s", async (name) => {
  const code = collectionUsage(name, defaults)!;
  expect(code).toContain("export default function Example");
  const result = await transformWithOxc(code, `${name}.tsx`, {
    jsx: "preserve",
  });
  expect(result.code.length).toBeGreaterThan(0);
});

it("keeps DataTable search usable while refining a selection", async () => {
  const [values, setValues] = createStore({ ...defaults });
  dispose = render(
    () => (
      <DataTableExample
        values={values}
        update={(key, value) =>
          setValues((draft) => {
            draft[key] = value;
          })
        }
      />
    ),
    host,
  );
  await new Promise((resolve) => setTimeout(resolve, 20));
  const input = host.querySelector('input[type="search"]') as HTMLInputElement;
  (
    host.querySelector('tbody input[type="checkbox"]') as HTMLInputElement
  ).click();
  await new Promise((resolve) => setTimeout(resolve, 20));
  expect(input.closest("[inert]")).toBeNull();
  expect(input.closest('[aria-hidden="true"]')).toBeNull();
  input.value = "Website";
  input.dispatchEvent(new Event("input", { bubbles: true }));
  await new Promise((resolve) => setTimeout(resolve, 20));
  expect(host.querySelector('input[type="search"]')).toBe(input);
  const cells = [...host.querySelectorAll("tbody tr td:nth-child(2)")];
  expect(cells.length).toBeGreaterThan(0);
  expect(cells.every((cell) => cell.textContent?.includes("Website"))).toBe(
    true,
  );
  expect(
    host.querySelector('[aria-label="Ações em lote"]')?.textContent,
  ).toContain("1 selecionado");
});
