import {
  createEffect,
  createSignal,
  createStore,
  onSettled,
  For,
  Show,
} from "solid-js";
import {
  Button,
  Card,
  CatalogNav,
  CodeBlock,
  Link,
  SearchInput,
  Tooltip,
  Toaster,
  toast,
} from "@trieoh/ui-solid";
import IconRotateCcw from "~icons/lucide/rotate-ccw";
import IconSun from "~icons/lucide/sun";
import IconMenu from "~icons/lucide/menu";
import IconX from "~icons/lucide/x";
import IconMoon from "~icons/lucide/moon";
import IconArrowUpRight from "~icons/lucide/arrow-up-right";
import { createCatalog, type Values } from "./catalog";
import { ScrollbarExample } from "./showcases/ScrollbarExample";
import { usage } from "./components/usage";
import { FitPreview } from "./components/FitPreview";
import { KnobSelect, KnobSwitch, KnobText } from "./components/KnobControl";

const ALL_GROUPS = "Todos";
const groupHref = (group: string) =>
  group === ALL_GROUPS ? "#/" : `#/group/${encodeURIComponent(group)}`;
const groups = [
  "Ações",
  "Formulários",
  "Navegação",
  "Overlays",
  "Feedback",
  "Layout",
];
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
  const navTop = {
    id: ALL_GROUPS,
    label: "Todos os componentes",
    href: groupHref(ALL_GROUPS),
  };
  const navSections = () =>
    groups.map((title) => ({
      id: title,
      title,
      href: groupHref(title),
      items: entries
        .filter((entry) => entry.group === title)
        .map((entry) => ({
          id: entry.name,
          label: entry.name,
          href: `#/components/${entry.name}`,
        })),
    }));
  const active = () => entries.find((entry) => entry.name === selected())!;
  const filtered = () =>
    entries.filter(
      (entry) =>
        (group() === ALL_GROUPS || entry.group === group()) &&
        `${entry.name} ${entry.description} ${entry.group}`
          .toLocaleLowerCase("pt-BR")
          .includes(query().toLocaleLowerCase("pt-BR")),
    );

  return (
    <div class="studio">
      <Toaster position="bottom-right" />
      <header class="studio-header">
        <div class="header-inner">
          <div class="header-navigation">
            <CatalogNav
              mode="topbar"
              topbarClass="lg:hidden"
              panelClass="top-[60px] max-h-[calc(100dvh-60px)]"
              top={navTop}
              activeId={detail() ? selected() : group()}
              sections={navSections()}
              renderTrigger={(menu) => (
                <Button
                  variant="ghost"
                  size="icon"
                  class="-ml-2"
                  aria-label={menu.open() ? "Fechar menu" : "Abrir menu"}
                  aria-expanded={menu.open() ? "true" : "false"}
                  onClick={menu.toggle}
                >
                  {menu.open() ? (
                    <IconX class="size-5" />
                  ) : (
                    <IconMenu class="size-5" />
                  )}
                </Button>
              )}
            />
            <Link href="#/" class="brand" aria-label="TrieOH UI — início">
              <strong>TrieOH</strong>
              <span class="brand-ui">UI</span>
            </Link>
            <span class="header-divider" aria-hidden="true" />
            <nav class="header-breadcrumb" aria-label="Navegação da biblioteca">
              <Link href="#/" class="header-catalog-link">
                Componentes
              </Link>
              <Show when={detail()}>
                <span aria-hidden="true">/</span>
                <span class="header-current">{active().name}</span>
              </Show>
            </nav>
          </div>
          <div class="header-actions">
            <span class="package-tag">@trieoh/ui-solid</span>
            <Tooltip
              placement="bottom"
              align="end"
              label={dark() ? "Usar tema claro" : "Usar tema escuro"}
            >
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setDark(!dark())}
                aria-label={dark() ? "Ativar tema claro" : "Ativar tema escuro"}
              >
                {dark() ? (
                  <IconSun class="size-4" />
                ) : (
                  <IconMoon class="size-4" />
                )}
              </Button>
            </Tooltip>
          </div>
        </div>
      </header>
      <div class="studio-body">
        <CatalogNav
          mode="sidebar"
          class="sticky top-[68px] hidden px-4 pt-8 pb-6 pl-10 lg:block"
          top={navTop}
          activeId={detail() ? selected() : group()}
          sections={navSections()}
        />
        <Show
          when={detail()}
          fallback={
            <main class="catalog">
              <div class="intro">
                <h1>{group() === ALL_GROUPS ? "Componentes" : group()}</h1>
                <p>
                  A UI do Univents, componente por componente. Teste estados e
                  variantes, ajuste as propriedades e copie o código em Solid.
                </p>
              </div>
              <div class="catalog-toolbar">
                <SearchInput
                  class="catalog-search"
                  aria-label="Buscar componentes"
                  placeholder="Buscar componente…"
                  value={query()}
                  onInput={(e) => setQuery(e.currentTarget.value)}
                />
                <span class="result-count" aria-live="polite">
                  {filtered().length} componentes
                </span>
              </div>
              <div class="component-grid">
                <For each={filtered()}>
                  {(entry, index) => (
                    <Card class="component-card">
                      <header>
                        <span class="component-number">
                          {String(index() + 1).padStart(2, "0")}
                        </span>
                        <span class="component-group">{entry.group}</span>
                        <Link
                          class="inspect-button"
                          aria-label={`Ver documentação de ${entry.name}`}
                          href={`#/components/${entry.name}`}
                        >
                          <IconArrowUpRight />
                        </Link>
                      </header>
                      <div class="component-preview">
                        <FitPreview>
                          {entry.render(values[entry.name], (key, value) =>
                            update(entry.name, key, value),
                          )}
                        </FitPreview>
                      </div>
                      <Link
                        class="component-caption"
                        href={`#/components/${entry.name}`}
                      >
                        <strong>{entry.name}</strong>
                        <span>{entry.description}</span>
                      </Link>
                    </Card>
                  )}
                </For>
              </div>
              <Show when={filtered().length === 0}>
                <div class="no-results">
                  <h2>Nenhum componente encontrado</h2>
                  <p>Tente outro nome ou categoria.</p>
                  <Button
                    variant="outline"
                    onClick={() => {
                      setQuery("");
                      window.location.hash = "#/";
                    }}
                  >
                    Limpar filtros
                  </Button>
                </div>
              </Show>
              <footer class="catalog-footer">
                <span>TrieOH Design System</span>
                <span>Solid · Tailwind CSS</span>
              </footer>
            </main>
          }
        >
          <main class="component-page">
            <Link class="back-link" href="#/">
              ← Todos os componentes
            </Link>
            <header class="component-page-header">
              <span>{active().group}</span>
              <h1>{active().name}</h1>
              <p>{active().description}</p>
            </header>
            <div class="playground-layout">
              <div class="component-documentation">
                <section class="example-section">
                  <h2>Exemplo</h2>
                  <div class="detail-preview">
                    {active().render(values[selected()], (key, value) =>
                      update(selected(), key, value),
                    )}
                  </div>
                </section>
                <section class="usage-section">
                  <h2>Como usar</h2>
                  <p>
                    Importe do pacote compartilhado. O exemplo acompanha os
                    valores definidos em propriedades.
                  </p>
                  <CodeBlock code={usage(active().name, values[selected()])} />
                </section>
              </div>
              <aside class="inspector" aria-label="Inspetor de componentes">
                <div class="inspector-heading">
                  <span class="eyebrow">Propriedades</span>
                  <Tooltip align="end" label="Restaurar valores padrão">
                    <Button
                      variant="ghost"
                      size="icon"
                      class="size-7 text-muted-foreground"
                      aria-label="Restaurar valores padrão"
                      disabled={Object.keys(active().defaults).every(
                        (key) =>
                          values[selected()][key] === active().defaults[key],
                      )}
                      onClick={() => {
                        setValues((draft) => {
                          Object.assign(draft[selected()], active().defaults);
                        });
                        toast.success("Propriedades restauradas");
                      }}
                    >
                      <IconRotateCcw class="size-3.5" />
                    </Button>
                  </Tooltip>
                </div>
                <div class="inspector-section">
                  <For each={active().controls}>
                    {(control) => (
                      <Show
                        when={
                          typeof values[selected()][control.key] === "boolean"
                        }
                        fallback={
                          <Show
                            when={control.options}
                            fallback={
                              <KnobText
                                multiline={control.key === "code"}
                                label={control.label}
                                value={String(values[selected()][control.key])}
                                onChange={(value) =>
                                  update(selected(), control.key, value)
                                }
                              />
                            }
                          >
                            <KnobSelect
                              label={control.label}
                              value={String(values[selected()][control.key])}
                              options={control.options!}
                              onChange={(value) =>
                                update(selected(), control.key, value)
                              }
                            />
                          </Show>
                        }
                      >
                        <KnobSwitch
                          label={control.label}
                          checked={!!values[selected()][control.key]}
                          onChange={(value) =>
                            update(selected(), control.key, value)
                          }
                        />
                      </Show>
                    )}
                  </For>
                </div>
              </aside>
            </div>
          </main>
        </Show>
      </div>
    </div>
  );
}
