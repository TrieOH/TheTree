import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { render } from "@solidjs/web";
import { createSignal } from "solid-js";
import {
  Pagination,
  PaginatedContainer,
  DataTable,
  InfiniteScroll,
  FilterBar,
  BulkActions,
  CollectionToolbar,
} from "@trieoh/ui-solid";

let host: HTMLDivElement;
let dispose: (() => void) | undefined;
const settle = () => new Promise((resolve) => setTimeout(resolve, 10));
const button = (name: string) =>
  [...host.querySelectorAll("button")].find(
    (element) =>
      element.getAttribute("aria-label") === name ||
      element.textContent === name,
  )!;
const rows = Array.from({ length: 25 }, (_, index) => ({
  id: index + 1,
  name: `Item ${index + 1}`,
}));
let diagnostics: ReturnType<typeof vi.spyOn>;
beforeEach(() => {
  host = document.createElement("div");
  document.body.append(host);
  diagnostics = vi.spyOn(console, "warn");
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
  expect(diagnostics.mock.calls.flat().join(" ")).not.toMatch(
    /STRICT_READ_UNTRACKED|REACTIVE_WRITE_IN_OWNED_SCOPE|NO_OWNER_EFFECT/,
  );
  vi.unstubAllGlobals();
});

describe("Pagination", () => {
  it("clamps pages, renders gaps and supports controlled navigation", async () => {
    const [page, setPage] = createSignal(50);
    dispose = render(
      () => (
        <Pagination
          page={page()}
          totalItems={1000}
          pageSize={10}
          onPageChange={setPage}
        />
      ),
      host,
    );
    await settle();
    expect(host.querySelector('[aria-current="page"]')?.textContent).toBe("50");
    expect(host.querySelectorAll('span[aria-hidden="true"]').length).toBe(2);
    button("Próxima página").click();
    await settle();
    expect(host.querySelector('[aria-current="page"]')?.textContent).toBe("51");
    setPage(999);
    await settle();
    expect((button("Próxima página") as HTMLButtonElement).disabled).toBe(true);
  });
  it("uses previous/next without inventing a total for cursor data", async () => {
    const change = vi.fn();
    dispose = render(
      () => <Pagination page={2} hasNextPage={false} onPageChange={change} />,
      host,
    );
    await settle();
    expect(host.querySelector('[aria-current="page"]')).toBeNull();
    expect(button("Próxima página").disabled).toBe(true);
    button("Página anterior").click();
    expect(change).toHaveBeenCalledWith(1);
  });
  it("handles zero results without an invalid range", async () => {
    dispose = render(
      () => <Pagination page={1} totalItems={0} onPageChange={() => {}} />,
      host,
    );
    await settle();
    expect(host.textContent).toContain("0–0 de 0");
    expect(button("Página anterior").disabled).toBe(true);
    expect(button("Próxima página").disabled).toBe(true);
  });
});

describe("PaginatedContainer", () => {
  it("slices local items and resets to the first page when searching", async () => {
    const [search, setSearch] = createSignal("");
    dispose = render(
      () => (
        <PaginatedContainer
          items={rows}
          pageSize={5}
          layout="list"
          filterValue={search()}
          onFilterChange={setSearch}
          renderItems={(items) => (
            <>
              {items.map((item) => (
                <p data-row>{item.name}</p>
              ))}
            </>
          )}
        />
      ),
      host,
    );
    await settle();
    expect(host.querySelectorAll("[data-row]").length).toBe(5);
    button("Próxima página").click();
    await settle();
    expect(host.querySelector("[data-row]")?.textContent).toBe("Item 6");
    const input = host.querySelector("input")!;
    input.value = "Item 25";
    input.dispatchEvent(new Event("input", { bubbles: true }));
    await settle();
    expect(host.querySelectorAll("[data-row]").length).toBe(1);
    expect(host.textContent).toContain("Item 25");
  });
  it("does not re-slice, filter or sort a page received from the server", async () => {
    const change = vi.fn();
    dispose = render(
      () => (
        <PaginatedContainer
          items={[
            { id: 25, name: "Z" },
            { id: 24, name: "A" },
          ]}
          pageSize={2}
          filterValue="absent"
          sortFields={[{ key: "name", label: "Nome" }]}
          pagination={{
            page: 3,
            pageSize: 2,
            totalItems: 25,
            onPageChange: change,
          }}
          renderItems={(items) => (
            <>
              {items.map((item) => (
                <p data-row>{item.name}</p>
              ))}
            </>
          )}
        />
      ),
      host,
    );
    await settle();
    expect(
      [...host.querySelectorAll("[data-row]")].map((el) => el.textContent),
    ).toEqual(["Z", "A"]);
    button("Próxima página").click();
    expect(change).toHaveBeenCalledWith(4);
  });
});

it("DataTable select-all changes only the current page and preserves other pages", async () => {
  const [selected, setSelected] = createSignal<(string | number)[]>([99]);
  dispose = render(
    () => (
      <DataTable
        caption="Itens"
        items={rows.slice(0, 2)}
        columns={[
          {
            id: "name",
            header: "Nome",
            cell: (item) => item.name,
            sortable: true,
          },
        ]}
        getRowKey={(item) => item.id}
        selectedKeys={selected()}
        onSelectionChange={setSelected}
      />
    ),
    host,
  );
  await settle();
  const checkbox = host.querySelector(
    'input[type="checkbox"]',
  ) as HTMLInputElement;
  checkbox.click();
  await settle();
  expect(selected()).toEqual([99, 1, 2]);
  checkbox.click();
  await settle();
  expect(selected()).toEqual([99]);
});

it("InfiniteScroll prevents duplicate requests, recovers from rejection and reaches the end", async () => {
  let finish: () => void = () => {};
  const request = vi
    .fn()
    .mockRejectedValueOnce(new Error("Falha"))
    .mockImplementationOnce(
      () =>
        new Promise<void>((resolve) => {
          finish = resolve;
        }),
    );
  const [more, setMore] = createSignal(true);
  dispose = render(
    () => (
      <InfiniteScroll automatic={false} hasMore={more()} onLoadMore={request}>
        <p>Conteúdo</p>
      </InfiniteScroll>
    ),
    host,
  );
  await settle();
  button("Carregar mais").click();
  await settle();
  expect(host.querySelector('[role="alert"]')?.textContent).toBe("Falha");
  const retryButton = button("Tentar novamente");
  retryButton.click();
  retryButton.click();
  await settle();
  expect(host.querySelector('[role="status"]')?.textContent).toBe(
    "Carregando…",
  );
  expect(request).toHaveBeenCalledTimes(2);
  finish();
  setMore(false);
  await settle();
  expect(host.textContent).toContain("Você chegou ao fim");
  expect(host.querySelector("button")).toBeNull();
});

it("InfiniteScroll observes the chosen scroll root and disconnects on unmount", async () => {
  let intersect: IntersectionObserverCallback = () => {};
  const disconnect = vi.fn();
  const load = vi.fn();
  const observer = vi.fn().mockImplementation(function (
    callback: IntersectionObserverCallback,
  ) {
    intersect = callback;
    return { observe() {}, disconnect };
  });
  vi.stubGlobal("IntersectionObserver", observer);
  dispose = render(
    () => (
      <InfiniteScroll scrollElement={host} hasMore={true} onLoadMore={load}>
        <p>Conteúdo</p>
      </InfiniteScroll>
    ),
    host,
  );
  await settle();
  expect(observer.mock.calls[0][1]).toEqual({
    root: host,
    rootMargin: "300px",
  });
  intersect(
    [{ isIntersecting: true }] as IntersectionObserverEntry[],
    {} as IntersectionObserver,
  );
  await settle();
  expect(load).toHaveBeenCalledTimes(1);
  dispose();
  dispose = undefined;
  expect(disconnect).toHaveBeenCalled();
});

it("FilterBar removes individual filters and clears the collection", async () => {
  const remove = vi.fn(),
    clear = vi.fn();
  dispose = render(
    () => (
      <FilterBar
        filters={[{ id: "status", label: "Status", value: "Ativo" }]}
        onRemove={remove}
        onClear={clear}
      />
    ),
    host,
  );
  await settle();
  button("Remover filtro Status: Ativo").click();
  expect(remove).toHaveBeenCalledWith("status");
  button("Limpar todos").click();
  expect(clear).toHaveBeenCalled();
});

it("BulkActions hides an empty selection and blocks actions while busy", async () => {
  const action = vi.fn();
  const [count, setCount] = createSignal(0);
  dispose = render(
    () => (
      <BulkActions
        selectedCount={count()}
        busy={true}
        onClearSelection={() => {}}
        actions={[{ id: "archive", label: "Arquivar", onClick: action }]}
      />
    ),
    host,
  );
  await settle();
  expect(host.querySelector('[role="region"]')).toBeNull();
  setCount(2);
  await settle();
  button("Arquivar").click();
  expect(action).not.toHaveBeenCalled();
});

it("measures a grid that appears after asynchronous data arrives", async () => {
  const [items, setItems] = createSignal<typeof rows>([]);
  const style = vi.spyOn(window, "getComputedStyle").mockReturnValue({
    gridTemplateColumns: "200px 200px 200px",
  } as CSSStyleDeclaration);
  dispose = render(
    () => (
      <PaginatedContainer
        items={items()}
        layout="grid"
        maxRows={2}
        renderItems={(items) => (
          <>
            {items.map((item) => (
              <p data-row>{item.name}</p>
            ))}
          </>
        )}
      />
    ),
    host,
  );
  await settle();
  setItems(rows);
  await settle();
  expect(host.querySelectorAll("[data-row]").length).toBe(6);
  style.mockRestore();
});

it("DataTable requests sorting and reports server errors with retry", async () => {
  const sort = vi.fn(),
    retry = vi.fn();
  const [error, setError] = createSignal<string>();
  dispose = render(
    () => (
      <DataTable
        caption="Itens"
        items={rows.slice(0, 2)}
        columns={[
          {
            id: "name",
            header: "Nome",
            cell: (item) => item.name,
            sortable: true,
          },
        ]}
        getRowKey={(item) => item.id}
        sort={{ column: "name", direction: "asc" }}
        onSortChange={sort}
        error={error()}
        onRetry={retry}
      />
    ),
    host,
  );
  await settle();
  [...host.querySelectorAll("button")]
    .find((element) => element.textContent?.startsWith("Nome"))!
    .click();
  expect(sort).toHaveBeenCalledWith({ column: "name", direction: "desc" });
  setError("Falha do servidor");
  await settle();
  expect(host.querySelector('[role="alert"]')?.textContent).toBe(
    "Falha do servidor",
  );
  button("Tentar novamente").click();
  expect(retry).toHaveBeenCalled();
});

it("preserves the legacy local footer and sort trigger without adding a size selector", async () => {
  dispose = render(
    () => (
      <PaginatedContainer
        items={rows}
        pageSize={5}
        itemLabel="projetos"
        sortFields={[{ key: "name", label: "Nome" }]}
        renderItems={(items) => (
          <>
            {items.map((item) => (
              <p>{item.name}</p>
            ))}
          </>
        )}
      />
    ),
    host,
  );
  await settle();
  expect(host.textContent).toContain("Mostrando 1–5 de 25 projetos");
  expect(host.querySelector('[aria-label="Ordenar"]')).not.toBeNull();
  expect(host.querySelector('[aria-label="Itens por página"]')).toBeNull();
});

it("can hide the optional page-size selector without losing page navigation", async () => {
  const [show, setShow] = createSignal(true);
  dispose = render(
    () => (
      <PaginatedContainer
        items={rows}
        pageSize={5}
        onPageSizeChange={() => {}}
        showPageSizeSelector={show()}
        renderItems={(items) => (
          <>
            {items.map((item) => (
              <p>{item.name}</p>
            ))}
          </>
        )}
      />
    ),
    host,
  );
  await settle();
  expect(host.querySelector('[aria-label="Itens por página"]')).not.toBeNull();
  setShow(false);
  await settle();
  expect(host.querySelector('[aria-label="Itens por página"]')).toBeNull();
  button("Próxima página").click();
  await settle();
  expect(host.textContent).toContain("Mostrando 6–10 de 25 itens");
});

it("CollectionToolbar composes optional pagination and removable filters", async () => {
  const page = vi.fn(),
    remove = vi.fn(),
    clear = vi.fn();
  const [show, setShow] = createSignal(true);
  dispose = render(
    () => (
      <CollectionToolbar
        filterBar={{
          filters: [{ id: "status", label: "Status", value: "Ativo" }],
          onRemove: remove,
          onClear: clear,
        }}
        pagination={
          show()
            ? { page: 1, pageSize: 5, totalItems: 10, onPageChange: page }
            : undefined
        }
      />
    ),
    host,
  );
  await settle();
  button("Próxima página").click();
  expect(page).toHaveBeenCalledWith(2);
  button("Remover filtro Status: Ativo").click();
  expect(remove).toHaveBeenCalledWith("status");
  setShow(false);
  await settle();
  expect(host.querySelector('[aria-label="Paginação"]')).toBeNull();
  button("Limpar todos").click();
  expect(clear).toHaveBeenCalled();
});

it("DataTable owns selection actions without replacing the mounted table", async () => {
  const [keys, setKeys] = createSignal<(string | number)[]>([]);
  const archive = vi.fn();
  dispose = render(
    () => (
      <DataTable
        caption="Itens"
        items={rows.slice(0, 2)}
        getRowKey={(item) => item.id}
        columns={[{ id: "name", header: "Nome", cell: (item) => item.name }]}
        selectedKeys={keys()}
        onSelectionChange={setKeys}
        toolbar={<p>Controles da coleção</p>}
        bulkActions={[
          {
            id: "archive",
            label: "Arquivar",
            icon: <span aria-hidden="true">A</span>,
            onClick: archive,
          },
        ]}
        bulkActionsIconOnly
      />
    ),
    host,
  );
  await settle();
  const table = host.querySelector("table");
  expect(
    host
      .querySelector('[aria-label="Ações em lote"]')
      ?.getAttribute("aria-hidden"),
  ).toBe("true");
  (
    host.querySelector("tbody input[type=checkbox]") as HTMLInputElement
  ).click();
  await settle();
  expect(host.querySelector("table")).toBe(table);
  expect(keys()).toEqual([1]);
  expect(
    host
      .querySelector('[aria-label="Ações em lote"]')
      ?.getAttribute("aria-hidden"),
  ).toBeNull();
  button("Arquivar").click();
  expect(archive).toHaveBeenCalledOnce();
  button("Limpar seleção").click();
  await settle();
  expect(keys()).toEqual([]);
});
