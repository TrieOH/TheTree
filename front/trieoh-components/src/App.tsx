import {
  createEffect,
  createSignal,
  createStore,
  onSettled,
  Show,
} from "solid-js";
import { CatalogNav, OverlayScrollbar, Toaster, toast } from "@trieoh/ui-solid";
import { createCatalog, type Values } from "./catalog";
import { ScrollbarExample } from "./examples/ScrollbarExample";
import { StudioHeader } from "./components/StudioHeader";
import { CatalogPage } from "./components/CatalogPage";
import { ComponentPage } from "./components/ComponentPage";

import {
  ALL_GROUPS,
  groups,
  catalogTop,
  catalogSections,
} from "./catalog-navigation";

export default function App() {
  if (window.location.hash.startsWith("#/preview/scrollbar"))
    return <ScrollbarExample />;
  const [dark, setDark] = createSignal(false);
  const [query, setQuery] = createSignal("");
  const [route, setRoute] = createSignal(window.location.hash.slice(1));
  const group = () => {
    const match = route().match(/^\/group\/(.+)$/);
    return match ? decodeURIComponent(match[1]) : ALL_GROUPS;
  };
  const selected = () => route().replace("/components/", "");
  const detail = () =>
    entries.some((entry) => route() === `/components/${entry.name}`);
  onSettled(() => {
    const sync = () => {
      setRoute(window.location.hash.slice(1));
      closeOverlays();
      window.scrollTo(0, 0);
    };
    window.addEventListener("hashchange", sync);
    return () => window.removeEventListener("hashchange", sync);
  });
  createEffect(
    () => dark(),
    (value) => {
      document.documentElement.classList.toggle("dark", value);
    },
  );

  const { entries, closeOverlays } = createCatalog({
    isDetail: detail,
    isDark: dark,
  });
  const [values, setValues] = createStore<Record<string, Values>>(
    Object.fromEntries(
      entries.map((entry) => [entry.name, { ...entry.defaults }]),
    ),
  );
  const update = (name: string, key: string, value: string | boolean) => {
    setValues((draft) => {
      draft[name][key] = value;
    });
  };
  const navTop = catalogTop;
  const navSections = catalogSections(entries);
  const active = () => entries.find((entry) => entry.name === selected())!;
  const orderedEntries = groups.flatMap((title) =>
    entries.filter((entry) => entry.group === title),
  );
  const activeIndex = () =>
    orderedEntries.findIndex((entry) => entry.name === selected());
  const previous = () => orderedEntries[activeIndex() - 1];
  const next = () => orderedEntries[activeIndex() + 1];
  const filtered = () =>
    entries.filter(
      (entry) =>
        (group() === ALL_GROUPS || entry.group === group()) &&
        `${entry.name} ${entry.description} ${entry.group}`
          .toLocaleLowerCase("pt-BR")
          .includes(query().toLocaleLowerCase("pt-BR")),
    );

  return (
    <div>
      <OverlayScrollbar />
      <Toaster position="bottom-right" />
      <StudioHeader
        navTop={navTop}
        sections={navSections}
        activeId={detail() ? selected() : group()}
        currentName={detail() ? active().name : undefined}
        dark={dark()}
        onToggleTheme={() => setDark(!dark())}
      />
      <div class="mx-auto grid max-w-360 grid-cols-1 items-start lg:grid-cols-[248px_minmax(0,1fr)]">
        <CatalogNav
          mode="sidebar"
          divider
          class="sticky top-17 hidden pt-8 pb-6 pl-10 lg:block"
          top={navTop}
          activeId={detail() ? selected() : group()}
          sections={navSections}
        />
        <Show
          when={detail()}
          fallback={
            <CatalogPage
              title={group() === ALL_GROUPS ? "Componentes" : group()}
              entries={filtered()}
              query={query()}
              values={values}
              onQueryChange={setQuery}
              onUpdate={update}
              onClearFilters={() => {
                setQuery("");
                window.location.hash = "#/";
              }}
            />
          }
        >
          <ComponentPage
            entry={active()}
            values={values[selected()]}
            previous={previous()}
            next={next()}
            onUpdate={(key, value) => update(selected(), key, value)}
            onReset={() => {
              setValues((draft) => {
                Object.assign(draft[selected()], active().defaults);
              });
              toast.success("Propriedades restauradas");
            }}
          />
        </Show>
      </div>
    </div>
  );
}
