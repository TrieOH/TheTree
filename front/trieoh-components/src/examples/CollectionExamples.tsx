import ArchiveIcon from "~icons/lucide/archive";
import DownloadIcon from "~icons/lucide/download";
import Folder from "~icons/lucide/folder";
import Plus from "~icons/lucide/plus";
import Filter from "~icons/lucide/list-filter";
import {
  createEffect,
  createMemo,
  createSignal,
  For,
  onCleanup,
  Show,
  untrack,
} from "solid-js";
import {
  Badge,
  Button,
  BulkActions,
  CollectionToolbar,
  Combobox,
  DataTable,
  FilterBar,
  InfiniteScroll,
  Pagination,
  PaginatedContainer,
  type ActiveFilter,
  type DataTableSort,
  type GapSize,
  type SortState,
} from "@trieoh/ui-solid";
import type { Values } from "../catalog";

type ExampleProps = {
  values: Values;
  update: (key: string, value: string | boolean) => void;
};
const projects = Array.from({ length: 48 }, (_, i) => ({
  id: i + 1,
  name: `${["Website", "Aplicativo", "Design system"][i % 3]} ${String(i + 1).padStart(2, "0")}`,
  status: i % 3 === 0 ? "Arquivado" : "Ativo",
}));

export function PaginationExample(props: ExampleProps) {
  const size = () => Number(props.values.pageSize);
  const total = () => Number(props.values.totalItems);
  const page = () =>
    Math.min(
      Number(props.values.page),
      Math.max(1, Math.ceil(total() / size())),
    );
  return (
    <div data-wide class="w-full">
      <Pagination
        showPageSizeSelector={props.values.showPageSizeSelector !== false}
        page={page()}
        pageSize={size()}
        totalItems={props.values.mode === "cursor" ? undefined : total()}
        hasNextPage={
          props.values.mode === "cursor" ? page() * size() < total() : undefined
        }
        onPageChange={(page) => props.update("page", String(page))}
        onPageSizeChange={(size) => {
          props.update("pageSize", String(size));
          props.update("page", "1");
        }}
        disabled={!!props.values.disabled}
      />
    </div>
  );
}

export function ToolbarExample(props: ExampleProps) {
  const [page, setPage] = createSignal(1);
  const current = () =>
    Math.min(page(), Math.max(1, Math.ceil(filtered().length / 3)));
  const [items, setItems] = createSignal(projects.slice(0, 6));
  const filtered = createMemo(() => {
    const search = String(props.values.search).toLocaleLowerCase();
    const status = props.values.status ?? "Todos";
    const result = items().filter(
      (item) =>
        item.name.toLocaleLowerCase().includes(search) &&
        (status === "Todos" || item.status === status),
    );
    return result.sort(
      (a, b) =>
        a.name.localeCompare(b.name) * (props.values.sort === "desc" ? -1 : 1),
    );
  });
  return (
    <div data-wide class="@container/projects w-full min-w-0 space-y-4">
      <header class="flex flex-wrap items-center justify-between gap-2">
        <div>
          <h3 class="text-sm font-medium">Biblioteca de projetos</h3>
          <p class="mt-1 text-xs text-muted-foreground">
            Busque, filtre ou mude a visualização.
          </p>
        </div>
        <Button
          size="sm"
          variant="outline"
          disabled={!!props.values.disabled}
          onClick={() => {
            const id = items().length + 1;
            setItems((previous) => [
              ...previous,
              { id, name: `Novo projeto ${id}`, status: "Ativo" },
            ]);
          }}
        >
          <Plus class="size-3.5" />
          Projeto
        </Button>
      </header>
      <CollectionToolbar
        pagination={
          props.values.withPagination
            ? {
                page: current(),
                pageSize: 3,
                totalItems: filtered().length,
                onPageChange: setPage,
                itemLabel: "projetos",
              }
            : undefined
        }
        filterBar={
          props.values.status && props.values.status !== "Todos"
            ? {
                filters: [
                  {
                    id: "status",
                    label: "Status",
                    value: String(props.values.status),
                  },
                ],
                onRemove: () => props.update("status", "Todos"),
                onClear: () => props.update("status", "Todos"),
              }
            : undefined
        }
        search={String(props.values.search)}
        onSearchChange={(value) => props.update("search", value)}
        totalItems={filtered().length}
        itemLabel="projetos"
        sort={String(props.values.sort)}
        sortOptions={[
          { value: "asc", label: "Nome: A → Z" },
          { value: "desc", label: "Nome: Z → A" },
        ]}
        onSortChange={(value) => props.update("sort", value)}
        layout={props.values.layout as "list" | "grid"}
        onLayoutChange={(value) => props.update("layout", value)}
        disabled={!!props.values.disabled}
        filters={
          <Combobox
            class="w-8 shrink-0"
            shape="soft"
            size="sm"
            iconOnly
            icon={<Filter class="size-4" />}
            triggerClass="border-input shadow-xs"
            aria-label="Filtrar por status"
            placeholder="Status"
            value={String(props.values.status ?? "Todos")}
            options={["Todos", "Ativo", "Arquivado"].map((value) => ({
              value,
              label: value,
            }))}
            onChange={(value) => props.update("status", value)}
            disabled={!!props.values.disabled}
            searchable={false}
          />
        }
      />
      <div
        class={
          props.values.layout === "grid"
            ? "grid grid-cols-1 gap-2 @[400px]/projects:grid-cols-2"
            : "divide-y divide-border rounded-md border border-border"
        }
      >
        <For
          each={
            props.values.withPagination
              ? filtered().slice((current() - 1) * 3, current() * 3)
              : filtered()
          }
        >
          {(item) => (
            <div
              class={
                props.values.layout === "grid"
                  ? "flex min-w-0 flex-col items-start gap-3 rounded-md border border-border p-4"
                  : "flex min-w-0 flex-wrap items-center gap-3 p-3"
              }
            >
              <Folder class="size-4 shrink-0 text-muted-foreground" />
              <span class="min-w-0 flex-1 text-xs font-medium wrap-anywhere">
                {item.name}
              </span>
              <span class="text-[11px] text-muted-foreground">
                {item.status}
              </span>
            </div>
          )}
        </For>
      </div>
      <Show when={!filtered().length}>
        <p class="py-6 text-center text-sm text-muted-foreground">
          Nenhum projeto encontrado.
        </p>
      </Show>
    </div>
  );
}

const initialFilters: ActiveFilter[] = [
  { id: "status", label: "Status", value: "Ativo" },
  { id: "category", label: "Categoria", value: "Website" },
];
export function FilterBarExample(props: ExampleProps) {
  const [filters, setFilters] = createSignal(initialFilters);
  const items = () =>
    projects
      .slice(0, 6)
      .filter((item) =>
        filters().every((filter) =>
          filter.id === "status"
            ? item.status === filter.value
            : item.name.startsWith(filter.value),
        ),
      );
  return (
    <div class="w-full space-y-4">
      <FilterBar
        filters={filters()}
        disabled={!!props.values.disabled}
        onRemove={(id) =>
          setFilters((previous) =>
            previous.filter((filter) => filter.id !== id),
          )
        }
        onClear={() => setFilters([])}
      >
        <Button
          size="sm"
          variant="outline"
          disabled={
            !!props.values.disabled ||
            filters().length === initialFilters.length
          }
          onClick={() => setFilters(initialFilters)}
        >
          Aplicar filtros de exemplo
        </Button>
      </FilterBar>
      <For each={items()}>
        {(item) => (
          <div class="rounded-md border border-border p-3 text-xs">
            {item.name} · {item.status}
          </div>
        )}
      </For>
    </div>
  );
}

export function BulkActionsExample(props: ExampleProps) {
  const [selected, setSelected] = createSignal(3);
  const [message, setMessage] = createSignal("");
  return (
    <div class="w-full space-y-3">
      <BulkActions
        iconOnly={!!props.values.iconOnly}
        selectedCount={selected()}
        busy={!!props.values.busy}
        onClearSelection={() => setSelected(0)}
        actions={[
          {
            id: "archive",
            label: "Arquivar",
            icon: <ArchiveIcon class="size-4" />,
            onClick: () => {
              setMessage(`${selected()} projetos arquivados.`);
              setSelected(0);
            },
          },
          {
            id: "export",
            label: "Exportar",
            icon: <DownloadIcon class="size-4" />,
            onClick: () =>
              setMessage(`${selected()} projetos preparados para exportação.`),
          },
        ]}
      />
      <Show when={!selected()}>
        <Button
          size="sm"
          variant="outline"
          onClick={() => {
            setSelected(3);
            setMessage("");
          }}
        >
          Selecionar 3 projetos
        </Button>
      </Show>
      <p role="status" class="text-xs text-muted-foreground">
        {message()}
      </p>
    </div>
  );
}

export function InfiniteScrollExample(props: ExampleProps) {
  const [count, setCount] = createSignal(8);
  const [root, setRoot] = createSignal<HTMLDivElement>();
  let disposed = false;
  onCleanup(() => {
    disposed = true;
  });
  const load = async () => {
    const simulateError = untrack(() => !!props.values.simulateError);
    await new Promise((resolve) => setTimeout(resolve, 500));
    if (disposed) return;
    if (simulateError)
      throw new Error(
        "Falha simulada. Desative o erro nas propriedades e tente novamente.",
      );
    setCount((previous) => Math.min(previous + 8, projects.length));
  };
  return (
    <div class="w-full">
      <div
        ref={setRoot}
        class="h-80 overflow-auto overscroll-contain rounded-md border border-border p-3"
      >
        <InfiniteScroll
          scrollElement={root()}
          hasMore={count() < projects.length}
          automatic={!!props.values.automatic}
          onLoadMore={load}
        >
          <div class="space-y-2">
            <For each={projects.slice(0, count())}>
              {(item) => (
                <div class="rounded-md border border-border bg-card p-3 text-xs">
                  {item.name}
                </div>
              )}
            </For>
          </div>
        </InfiniteScroll>
      </div>
    </div>
  );
}

export function PaginatedExample(props: ExampleProps) {
  const [page, setPage] = createSignal(1);
  const [sort, setSort] = createSignal<SortState<(typeof projects)[number]>>({
    field: "name",
    direction: "asc",
  });
  const [result, setResult] = createSignal({
    items: [] as typeof projects,
    total: 0,
  });
  const [loading, setLoading] = createSignal(false);
  const [retry, setRetry] = createSignal(0);
  const [error, setError] = createSignal<string>();
  const server = () => props.values.dataMode === "No servidor (simulado)";
  createEffect(
    () => ({
      server: server(),
      page: page(),
      size: Number(props.values.pageSize),
      search: String(props.values.filterValue),
      sort: sort(),
      retry: retry(),
      fail: !!props.values.simulateError,
    }),
    (request) => {
      if (!request.server) return;
      setLoading(true);
      setError(undefined);
      const timer = setTimeout(() => {
        if (request.fail) {
          setError("Não foi possível carregar os projetos.");
          setLoading(false);
          return;
        }
        const filtered = projects.filter((item) =>
          item.name
            .toLocaleLowerCase()
            .includes(request.search.toLocaleLowerCase()),
        );
        filtered.sort(
          (a, b) =>
            a.name.localeCompare(b.name) *
            (request.sort.direction === "asc" ? 1 : -1),
        );
        setResult({
          items: filtered.slice(
            (request.page - 1) * request.size,
            request.page * request.size,
          ),
          total: filtered.length,
        });
        setLoading(false);
      }, 400);
      return () => clearTimeout(timer);
    },
  );
  createEffect(
    () => [props.values.pageSize, props.values.dataMode],
    () => {
      setPage(1);
    },
  );
  return (
    <div data-wide class="w-full space-y-2">
      <p class="text-xs leading-relaxed text-muted-foreground">
        {server()
          ? "Demonstração: o servidor entrega somente os itens da página solicitada."
          : "Demonstração: os 48 projetos já estão carregados; o componente busca, ordena e divide as páginas."}
      </p>
      <PaginatedContainer
        items={server() ? result().items : projects}
        showPageSizeSelector={!!props.values.showPageSizeSelector}
        pageSizeOptions={[5, 10, 25]}
        onPageSizeChange={(size) => {
          props.update("pageSize", String(size));
          setPage(1);
        }}
        pageSize={Number(props.values.pageSize)}
        layout={props.values.layout as "list" | "grid"}
        gap={props.values.gap as GapSize}
        filterValue={String(props.values.filterValue)}
        onFilterChange={(value) => props.update("filterValue", value)}
        filterPlaceholder={String(props.values.filterPlaceholder)}
        itemLabel={String(props.values.itemLabel)}
        minItemWidth={String(props.values.minItemWidth)}
        sortFields={[{ key: "name", label: "Nome" }]}
        sort={sort()}
        onSortChange={setSort}
        pagination={
          server()
            ? {
                page: page(),
                pageSize: Number(props.values.pageSize),
                totalItems: result().total,
                onPageChange: setPage,
                onPageSizeChange: (size) => {
                  props.update("pageSize", String(size));
                  setPage(1);
                },
              }
            : undefined
        }
        loading={server() && loading()}
        error={server() ? error() : undefined}
        onRetry={() => setRetry((value) => value + 1)}
        emptyState={
          <p class="py-6 text-center text-sm text-muted-foreground">
            Nenhum projeto encontrado.
          </p>
        }
        renderItems={(items) => (
          <For each={items}>
            {(item) => (
              <div class="rounded-md border border-border p-3 text-xs">
                {item.name}
              </div>
            )}
          </For>
        )}
      />
    </div>
  );
}

export function DataTableExample(props: ExampleProps) {
  const [items, setItems] = createSignal(projects.slice(0, 18));
  const [selected, setSelected] = createSignal<(string | number)[]>([]);
  const [page, setPage] = createSignal(1);
  const [size, setSize] = createSignal(5);
  const [search, setSearch] = createSignal("");
  const [sort, setSort] = createSignal<DataTableSort>({
    column: "name",
    direction: "asc",
  });
  const filtered = createMemo(() =>
    items()
      .filter((item) =>
        item.name.toLocaleLowerCase().includes(search().toLocaleLowerCase()),
      )
      .sort(
        (a, b) =>
          String(a[sort().column as "name" | "status"]).localeCompare(
            String(b[sort().column as "name" | "status"]),
          ) * (sort().direction === "asc" ? 1 : -1),
      ),
  );
  const current = () =>
    Math.min(page(), Math.max(1, Math.ceil(filtered().length / size())));
  return (
    <div data-wide class="w-full space-y-4">
      <DataTable
        toolbar={
          <CollectionToolbar
            search={search()}
            onSearchChange={(value) => {
              setSearch(value);
              setPage(1);
            }}
          />
        }
        bulkActions={[
          {
            id: "archive",
            label: "Arquivar",
            icon: <ArchiveIcon class="size-4" />,
            onClick: () => {
              const keys = selected();
              setItems((previous) =>
                previous.map((item) =>
                  keys.includes(item.id)
                    ? { ...item, status: "Arquivado" }
                    : item,
                ),
              );
              setSelected([]);
            },
          },
        ]}
        caption="Projetos"
        items={filtered().slice((current() - 1) * size(), current() * size())}
        getRowKey={(item) => item.id}
        columns={[
          {
            id: "name",
            header: "Projeto",
            cell: (item) => item.name,
            sortable: true,
          },
          {
            id: "status",
            header: "Status",
            cell: (item) => <Badge>{item.status}</Badge>,
            sortable: true,
          },
        ]}
        selectedKeys={selected()}
        onSelectionChange={props.values.selectable ? setSelected : undefined}
        sort={sort()}
        onSortChange={(value) => {
          setSort(value);
          setPage(1);
        }}
        loading={!!props.values.loading}
        error={
          props.values.simulateError
            ? "Não foi possível carregar os projetos."
            : undefined
        }
        onRetry={() => props.update("simulateError", false)}
        rowActions={
          props.values.rowActions
            ? (item) => (
                <Button
                  size="sm"
                  variant="ghost"
                  disabled={!!props.values.loading}
                  onClick={() =>
                    setItems((previous) =>
                      previous.map((project) =>
                        project.id === item.id
                          ? {
                              ...project,
                              status:
                                project.status === "Ativo"
                                  ? "Arquivado"
                                  : "Ativo",
                            }
                          : project,
                      ),
                    )
                  }
                >
                  {item.status === "Ativo" ? "Arquivar" : "Restaurar"}
                </Button>
              )
            : undefined
        }
        pagination={{
          showPageSizeSelector: props.values.showPageSizeSelector !== false,
          itemLabel: "projetos",
          page: current(),
          pageSize: size(),
          totalItems: filtered().length,
          onPageChange: setPage,
          onPageSizeChange: (value) => {
            setSize(value);
            setPage(1);
          },
          pageSizeOptions: [5, 10, 25],
        }}
      />
    </div>
  );
}
