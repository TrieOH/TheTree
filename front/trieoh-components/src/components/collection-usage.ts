import type { Values } from "../catalog";

/** Complete, copyable collection examples with framework imports kept separate. */
export function collectionUsage(
  name: string,
  values: Values,
): string | undefined {
  const wrap = (ui: string, solid: string, body: string, extraImports = "") =>
    [
      `import { ${ui} } from "@trieoh/ui-solid";`,
      ...(solid ? [`import { ${solid} } from "solid-js";`] : []),
      ...(extraImports ? [extraImports] : []),
      "",
      `export default function Example() {`,
      body,
      "}",
    ].join(String.fromCharCode(10));
  switch (name) {
    case "Pagination": {
      const cursor = values.mode === "cursor";
      return wrap(
        "Pagination",
        "createSignal",
        `  const [page, setPage] = createSignal(${Number(values.page)});
  const [size, setSize] = createSignal(${Number(values.pageSize)});
  return <Pagination page={page()} pageSize={size()} ${cursor ? `hasNextPage={page() * size() < ${Number(values.totalItems)}}` : `totalItems={${Number(values.totalItems)}}`}
    showPageSizeSelector={${values.showPageSizeSelector !== false}}
    onPageChange={setPage}
    onPageSizeChange={size => { setSize(size); setPage(1); }}
    disabled={${!!values.disabled}} />;`,
      );
    }
    case "CollectionToolbar":
      return wrap(
        "CollectionToolbar",
        "createSignal",
        `  const [search, setSearch] = createSignal(${JSON.stringify(values.search)});
  const [sort, setSort] = createSignal(${JSON.stringify(values.sort)});
  const [layout, setLayout] = createSignal<"list" | "grid">(${JSON.stringify(values.layout)});
  const [page, setPage] = createSignal(1);
  // Apply these controlled values to your local collection or server query.
  return <CollectionToolbar ${values.withPagination ? "pagination={{ page: page(), pageSize: 3, totalItems: 6, onPageChange: setPage }}" : ""} search={search()} onSearchChange={setSearch}
    sort={sort()} onSortChange={setSort}
    sortOptions={[{ value: "asc", label: "Nome: A → Z" }, { value: "desc", label: "Nome: Z → A" }]}
    layout={layout()} onLayoutChange={setLayout} disabled={${!!values.disabled}} />;`,
      );
    case "FilterBar":
      return wrap(
        "FilterBar, type ActiveFilter",
        "createSignal",
        `  const [filters, setFilters] = createSignal<ActiveFilter[]>([
    { id: "status", label: "Status", value: "Ativo" },
    { id: "category", label: "Categoria", value: "Website" },
  ]);
  return <FilterBar filters={filters()} disabled={${!!values.disabled}}
    onRemove={id => setFilters(filters => filters.filter(filter => filter.id !== id))}
    onClear={() => setFilters([])} />;`,
      );
    case "BulkActions":
      return wrap(
        "BulkActions",
        "createSignal",
        `  const [selected, setSelected] = createSignal([1, 2, 3]);
  const [message, setMessage] = createSignal("");
  return <>
    <BulkActions iconOnly={${!!values.iconOnly}} selectedCount={selected().length} busy={${!!values.busy}}
      onClearSelection={() => setSelected([])}
      actions={[{ id: "archive", label: "Arquivar", icon: <ArchiveIcon class="size-4" />, onClick: () => {
        setMessage(selected().length + " projetos arquivados.");
        setSelected([]);
      } }]} />
    <p role="status">{message()}</p>
  </>;`,
        'import ArchiveIcon from "~icons/lucide/archive";',
      );
    case "InfiniteScroll":
      return wrap(
        "InfiniteScroll",
        "createSignal, For",
        `  const all = Array.from({ length: 48 }, (_, index) => index + 1);
  const [items, setItems] = createSignal(all.slice(0, 8));
  const [root, setRoot] = createSignal<HTMLDivElement>();
  const loadMore = async () => {
    // Replace this simulated request with your backend cursor request.
    const offset = items().length;
    await new Promise(resolve => setTimeout(resolve, 500));
    ${values.simulateError ? 'throw new Error("Falha simulada.");' : "setItems(previous => [...previous, ...all.slice(offset, offset + 8)]);"}
  };
  return <div ref={setRoot} class="h-80 overflow-auto">
    <InfiniteScroll scrollElement={root()} hasMore={items().length < all.length}
      automatic={${!!values.automatic}} onLoadMore={loadMore}>
      <For each={items()}>{item => <div class="border-b p-4">Item {item}</div>}</For>
    </InfiniteScroll>
  </div>;`,
      );
    case "DataTable":
      return wrap(
        "DataTable, type DataTableSort",
        "createSignal, createMemo",
        `  const items = Array.from({ length: 18 }, (_, index) => ({ id: index + 1, name: "Projeto " + (index + 1) }));
  const [page, setPage] = createSignal(1);
  const [size, setSize] = createSignal(5);
  const [selected, setSelected] = createSignal<(string | number)[]>([]);
  const [sort, setSort] = createSignal<DataTableSort>({ column: "name", direction: "asc" });
  const sorted = createMemo(() => [...items].sort((a, b) => a.name.localeCompare(b.name, undefined, { numeric: true }) * (sort().direction === "asc" ? 1 : -1)));
  return <DataTable caption="Projetos" getRowKey={item => item.id}
    items={sorted().slice((page() - 1) * size(), page() * size())}
    columns={[{ id: "name", header: "Projeto", cell: item => item.name, sortable: true }]}
    selectedKeys={selected()} ${values.selectable ? "onSelectionChange={setSelected}" : ""}
    sort={sort()} onSortChange={value => { setSort(value); setPage(1); }}
    loading={${!!values.loading}} ${values.simulateError ? 'error="Falha simulada."' : ""}
    pagination={{ showPageSizeSelector: ${values.showPageSizeSelector !== false}, page: page(), pageSize: size(), totalItems: items.length,
      onPageChange: setPage, onPageSizeChange: size => { setSize(size); setPage(1); } }} />;`,
      );
    case "PaginatedContainer": {
      const server = values.dataMode === "No servidor (simulado)";
      if (!server)
        return wrap(
          "PaginatedContainer",
          "createSignal, For",
          `  const items = Array.from({ length: 48 }, (_, index) => ({ name: "Projeto " + (index + 1) }));
  const [search, setSearch] = createSignal(${JSON.stringify(values.filterValue)});
  const [size, setSize] = createSignal(${Number(values.pageSize)});
  return <PaginatedContainer items={items} pageSize={size()} onPageSizeChange={setSize}
    showPageSizeSelector={${values.showPageSizeSelector !== false}} layout=${JSON.stringify(values.layout)}
    filterValue={search()} onFilterChange={setSearch} sortFields={[{ key: "name", label: "Nome" }]}
    emptyState={<p>Nenhum projeto encontrado.</p>}
    renderItems={items => <For each={items}>{item => <div class="rounded-md border p-3">{item.name}</div>}</For>} />;`,
        );
      return wrap(
        "PaginatedContainer, type SortState",
        "createEffect, createSignal, For",
        `  type Project = { name: string };
  const all = Array.from({ length: 48 }, (_, index) => ({ name: "Projeto " + (index + 1) }));
  const [page, setPage] = createSignal(1);
  const [size, setSize] = createSignal(${Number(values.pageSize)});
  const [search, setSearch] = createSignal("");
  const [sort, setSort] = createSignal<SortState<Project>>({ field: "name", direction: "asc" });
  const [result, setResult] = createSignal<{ items: Project[]; total: number }>({ items: [], total: 0 });
  const [loading, setLoading] = createSignal(false);
  createEffect(() => ({ page: page(), size: size(), search: search(), sort: sort() }), request => {
    setLoading(true);
    // Simulated server: replace with an abortable request returning one page.
    const timer = setTimeout(() => {
      const filtered = all.filter(item => item.name.toLowerCase().includes(request.search.toLowerCase()));
      filtered.sort((a, b) => a.name.localeCompare(b.name) * (request.sort.direction === "asc" ? 1 : -1));
      setResult({ items: filtered.slice((request.page - 1) * request.size, request.page * request.size), total: filtered.length });
      setLoading(false);
    }, 400);
    return () => clearTimeout(timer);
  });
  return <PaginatedContainer items={result().items} loading={loading()} layout=${JSON.stringify(values.layout)}
    showPageSizeSelector={${values.showPageSizeSelector !== false}}
    pagination={{ page: page(), pageSize: size(), totalItems: result().total,
      onPageChange: setPage, onPageSizeChange: size => { setSize(size); setPage(1); } }}
    filterValue={search()} onFilterChange={setSearch}
    sortFields={[{ key: "name", label: "Nome" }]} sort={sort()} onSortChange={setSort}
    renderItems={items => <For each={items}>{item => <div class="rounded-md border p-3">{item.name}</div>}</For>} />;`,
      );
    }
  }
}
