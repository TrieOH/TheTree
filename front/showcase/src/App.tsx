import { ScrollbarExample } from "./showcases/ScrollbarExample";
import IconRotateCcw from "~icons/lucide/rotate-ccw";
import {
  createEffect,
  createSignal,
  createStore,
  onSettled,
  For,
  Show,
} from "solid-js";
import type { JSX } from "@solidjs/web";
import {
  Button,
  Card,
  Checkbox,
  Link,
  SearchInput,
  CodeBlock,
  Input,
  Textarea,
  Label,
  Field,
  Dialog,
  AlertModal,
  Drawer,
  Tooltip,
  Accordion,
  SectionTabs,
  Combobox,
  EmptyState,
  Skeleton,
  CardSkeleton,
  Toaster,
  toast,
  MultiStepDialog,
  MultiStepStepper,
  MultiStepFields,
  MultiStepSummaryCard,
  PaginatedContainer,
  type ButtonVariant,
  type ButtonSize,
  type DialogSize,
  type AlertModalVariant,
} from "@trieoh/ui-solid";
import { usage } from "./components/usage";
import { KnobSelect, KnobSwitch, KnobText } from "./components/KnobControl";
import IconSun from "~icons/lucide/sun";
import IconMoon from "~icons/lucide/moon";
import IconArrowUpRight from "~icons/lucide/arrow-up-right";

const groups = [
  "Todos",
  "Ações",
  "Formulários",
  "Navegação",
  "Overlays",
  "Feedback",
  "Layout",
];
type Values = Record<string, string | boolean>;
type Control = { key: string; label: string; options?: readonly string[] };
type Entry = {
  name: string;
  group: string;
  description: string;
  defaults: Values;
  controls: Control[];
  render: (
    v: Values,
    update: (key: string, value: string | boolean) => void,
  ) => JSX.Element;
};
const text = (key: string, label: string): Control => ({ key, label });
const select = (key: string, label: string, options: string[]): Control => ({
  key,
  label,
  options,
});

export default function App() {
  if (window.location.hash.startsWith("#/preview/scrollbar"))
    return <ScrollbarExample />;
  const [dark, setDark] = createSignal(false);
  const [query, setQuery] = createSignal("");
  const [group, setGroup] = createSignal("Todos");
  const [route, setRoute] = createSignal(window.location.hash.slice(1));
  const selected = () => route().replace("/components/", "");
  const detail = () =>
    entries.some((entry) => route() === `/components/${entry.name}`);
  onSettled(() => {
    const sync = () => {
      setRoute(window.location.hash.slice(1));
      setOpened("");
      window.scrollTo(0, 0);
    };
    window.addEventListener("hashchange", sync);
    return () => window.removeEventListener("hashchange", sync);
  });
  const [opened, setOpened] = createSignal("");
  const [form, setForm] = createSignal({ name: "", description: "" });
  const steps = [
    {
      id: "details",
      title: "Detalhes",
      fields: [
        {
          name: "name" as const,
          label: "Nome",
          placeholder: "Nome do projeto",
        },
        {
          name: "description" as const,
          label: "Descrição",
          placeholder: "Descrição do projeto",
        },
      ],
    },
    {
      id: "review",
      title: "Revisão",
      summary: {
        title: "Revise os dados",
        items: [
          { label: "Nome", value: () => form().name },
          { label: "Descrição", value: () => form().description },
        ],
      },
    },
  ];
  createEffect(
    () => dark(),
    (value) => {
      document.documentElement.classList.toggle("dark", value);
    },
  );

  const entries: Entry[] = [
    {
      name: "Checkbox",
      group: "Formulários",
      description: "Seleção booleana com rótulo e descrição.",
      defaults: {
        label: "Receber notificações",
        checked: true,
        disabled: false,
      },
      controls: [
        text("label", "Rótulo"),
        text("checked", "Marcado"),
        text("disabled", "Desabilitado"),
      ],
      render: (v, update) => (
        <Checkbox
          label={String(v.label)}
          checked={!!v.checked}
          disabled={!!v.disabled}
          onChange={(value) => update("checked", value)}
        />
      ),
    },
    {
      name: "SearchInput",
      group: "Formulários",
      description: "Campo de busca com ícone e nome acessível.",
      defaults: { placeholder: "Buscar projetos…", value: "" },
      controls: [text("placeholder", "Placeholder")],
      render: (v, update) => (
        <SearchInput
          aria-label="Buscar projetos"
          placeholder={String(v.placeholder)}
          value={String(v.value)}
          onInput={(e) => update("value", e.currentTarget.value)}
        />
      ),
    },
    {
      name: "Card",
      group: "Layout",
      description: "Superfície de conteúdo com borda e fundo semânticos.",
      defaults: { children: "Conteúdo do card" },
      controls: [text("children", "Conteúdo")],
      render: (v) => <Card class="w-full p-5">{String(v.children)}</Card>,
    },
    {
      name: "Link",
      group: "Navegação",
      description: "Navegação por link com foco visível.",
      defaults: { children: "Ver componentes" },
      controls: [text("children", "Texto")],
      render: (v) => <Link href="#/">{String(v.children)}</Link>,
    },
    {
      name: "CodeBlock",
      group: "Layout",
      description: "Exemplo de código com ação de copiar.",
      defaults: { code: "<Button>Salvar</Button>" },
      controls: [text("code", "Código")],
      render: (v) => <CodeBlock code={String(v.code)} />,
    },
    {
      name: "Button",
      group: "Ações",
      description: "Botão com variantes, tamanhos e estado desabilitado.",
      defaults: {
        children: "Criar projeto",
        variant: "default",
        size: "md",
        disabled: false,
      },
      controls: [
        text("children", "Texto"),
        select("variant", "Variante", [
          "default",
          "outline",
          "ghost",
          "destructive",
          "link",
        ]),
        select("size", "Tamanho", ["sm", "md", "lg", "icon"]),
        text("disabled", "Desabilitado"),
      ],
      render: (v) => (
        <Button
          variant={v.variant as ButtonVariant}
          size={v.size as ButtonSize}
          disabled={!!v.disabled}
          aria-label={String(v.children)}
        >
          {v.size === "icon" ? (
            <IconArrowUpRight class="size-4" />
          ) : (
            String(v.children)
          )}
        </Button>
      ),
    },
    {
      name: "Input",
      group: "Formulários",
      description: "Campo de texto de uma linha.",
      defaults: { placeholder: "Seu e-mail", disabled: false },
      controls: [
        text("placeholder", "Placeholder"),
        text("disabled", "Desabilitado"),
      ],
      render: (v) => (
        <Input
          placeholder={String(v.placeholder)}
          disabled={!!v.disabled}
          aria-label="Exemplo de entrada"
        />
      ),
    },
    {
      name: "Textarea",
      group: "Formulários",
      description: "Campo de texto com múltiplas linhas.",
      defaults: {
        placeholder: "Conte um pouco sobre o projeto…",
        disabled: false,
      },
      controls: [
        text("placeholder", "Placeholder"),
        text("disabled", "Desabilitado"),
      ],
      render: (v) => (
        <Textarea
          placeholder={String(v.placeholder)}
          disabled={!!v.disabled}
          rows={3}
          aria-label="Exemplo de área de texto"
        />
      ),
    },
    {
      name: "Label",
      group: "Formulários",
      description: "Rótulo associado a um campo de formulário.",
      defaults: { children: "Nome do projeto" },
      controls: [text("children", "Texto")],
      render: (v) => <Label>{String(v.children)}</Label>,
    },
    {
      name: "Field",
      group: "Formulários",
      description: "Composição de rótulo, descrição e erro de validação.",
      defaults: {
        label: "Nome completo",
        hint: "Como você gostaria de ser chamado?",
        error: "",
      },
      controls: [
        text("label", "Rótulo"),
        text("hint", "Ajuda"),
        text("error", "Erro"),
      ],
      render: (v) => (
        <Field
          label={String(v.label)}
          hint={String(v.hint)}
          error={String(v.error)}
        >
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              placeholder="Maria Silva"
            />
          )}
        </Field>
      ),
    },
    {
      name: "Combobox",
      group: "Formulários",
      description: "Seleção de uma opção com busca por texto.",
      defaults: { placeholder: "Selecione uma cidade", value: "" },
      controls: [
        text("placeholder", "Placeholder"),
        select("value", "Valor", ["", "sp", "rj", "bh"]),
      ],
      render: (v, update) => (
        <Combobox
          options={[
            { value: "sp", label: "São Paulo" },
            { value: "rj", label: "Rio de Janeiro" },
            { value: "bh", label: "Belo Horizonte" },
          ]}
          placeholder={String(v.placeholder)}
          value={String(v.value)}
          onChange={(value) => update("value", value)}
        />
      ),
    },
    {
      name: "SectionTabs",
      group: "Navegação",
      description: "Abas de seção com contagem e seleção controlada.",
      defaults: { active: "overview" },
      controls: [
        select("active", "Aba ativa", ["overview", "activity", "settings"]),
      ],
      render: (v, update) => (
        <SectionTabs
          items={[
            { id: "overview", label: "Visão geral" },
            { id: "activity", label: "Atividade", count: 8 },
            { id: "settings", label: "Ajustes" },
          ]}
          active={String(v.active)}
          onChange={(value) => update("active", value)}
        />
      ),
    },
    {
      name: "Accordion",
      group: "Navegação",
      description: "Lista de seções expansíveis.",
      defaults: { title: "Como começar?", multiple: false },
      controls: [
        text("title", "Título"),
        text("multiple", "Permitir vários abertos"),
      ],
      render: (v) => (
        <Accordion
          multiple={!!v.multiple}
          items={[
            {
              value: "start",
              title: String(v.title),
              content: "Escolha um componente e explore suas propriedades.",
            },
            {
              value: "more",
              title: "Posso personalizar?",
              content: "Sim. Use o inspetor ao lado para experimentar.",
            },
          ]}
        />
      ),
    },
    ...(["Dialog", "AlertModal", "Drawer", "MultiStepDialog"] as const).map(
      (name) => ({
        name,
        group: "Overlays",
        description: {
          Dialog: "Modal com título, conteúdo e tamanhos configuráveis.",
          AlertModal: "Modal de confirmação com variantes semânticas.",
          Drawer: "Painel sobreposto que abre pela base da tela.",
          MultiStepDialog: "Formulário dividido em etapas.",
        }[name],
        defaults: {
          title: name === "AlertModal" ? "Excluir projeto?" : "Um novo projeto",
          size: "md",
          variant: "warning",
        },
        controls: [
          text("title", "Título"),
          ...(name === "Dialog"
            ? [select("size", "Tamanho", ["sm", "md", "lg", "xl", "full"])]
            : []),
          ...(name === "AlertModal"
            ? [
                select("variant", "Variante", [
                  "default",
                  "info",
                  "warning",
                  "destructive",
                ]),
              ]
            : []),
        ],
        render: (v: Values) => (
          <>
            <Button variant="outline" onClick={() => setOpened(name)}>
              Experimentar {name} <IconArrowUpRight class="size-3.5" />
            </Button>
            <Show when={name === "Dialog"}>
              <Dialog
                open={opened() === name}
                onOpenChange={(value) => setOpened(value ? name : "")}
                title={String(v.title)}
                size={v.size as DialogSize}
              >
                <p class="text-sm text-muted-foreground">
                  Este é o conteúdo do seu dialog.
                </p>
                <Button class="mt-5" onClick={() => setOpened("")}>
                  Concluir
                </Button>
              </Dialog>
            </Show>
            <Show when={name === "AlertModal"}>
              <AlertModal
                open={opened() === name}
                onOpenChange={(value) => setOpened(value ? name : "")}
                title={String(v.title)}
                description="Esta é uma demonstração de confirmação."
                variant={v.variant as AlertModalVariant}
                onConfirm={() => {
                  setOpened("");
                }}
              />
            </Show>
            <Show when={name === "Drawer"}>
              <Drawer
                open={opened() === name}
                onClose={() => setOpened("")}
                title={String(v.title)}
              >
                <p class="text-sm text-muted-foreground">
                  Explore o conteúdo e feche quando terminar.
                </p>
              </Drawer>
            </Show>
            <Show when={name === "MultiStepDialog"}>
              <MultiStepDialog
                open={opened() === name}
                onOpenChange={(value) => setOpened(value ? name : "")}
                title={String(v.title)}
                steps={steps}
                values={form()}
                onChange={(key, value) =>
                  setForm((previous) => ({ ...previous, [key]: String(value) }))
                }
                onSubmit={() => {
                  toast.success("Projeto criado");
                  return true;
                }}
              />
            </Show>
          </>
        ),
      }),
    ),
    {
      name: "Tooltip",
      group: "Overlays",
      description: "Texto auxiliar exibido ao passar o mouse ou focar.",
      defaults: { label: "Informações adicionais sobre esta ação." },
      controls: [text("label", "Mensagem")],
      render: (v) => (
        <Tooltip label={String(v.label)}>
          <Button variant="outline">Passe o mouse ou foque aqui</Button>
        </Tooltip>
      ),
    },
    {
      name: "Toaster",
      group: "Feedback",
      description: "Notificações de sucesso, erro, aviso e informação.",
      defaults: { message: "Projeto salvo com sucesso", type: "success" },
      controls: [
        text("message", "Mensagem"),
        select("type", "Tipo", ["success", "error", "warning", "info"]),
      ],
      render: (v) => (
        <Button
          variant="outline"
          onClick={() =>
            toast[v.type as "success" | "error" | "warning" | "info"](
              String(v.message),
            )
          }
        >
          Mostrar notificação <span aria-hidden="true">↗</span>
        </Button>
      ),
    },
    {
      name: "EmptyState",
      group: "Feedback",
      description: "Mensagem para coleções sem registros.",
      defaults: {
        title: "Nenhum projeto encontrado",
        description: "Crie um projeto para começar.",
      },
      controls: [text("title", "Título"), text("description", "Descrição")],
      render: (v) => (
        <EmptyState
          title={String(v.title)}
          description={String(v.description)}
          class="w-full py-6! px-4!"
        />
      ),
    },
    {
      name: "Skeleton",
      group: "Feedback",
      description: "Placeholder de carregamento.",
      defaults: { rounded: true },
      controls: [text("rounded", "Avatar circular")],
      render: (v) => (
        <div class="flex w-full items-center gap-3">
          <Skeleton
            class={
              v.rounded ? "size-11 shrink-0 rounded-full" : "size-11 shrink-0"
            }
          />
          <div class="flex-1 space-y-2">
            <Skeleton class="h-3 w-3/4" />
            <Skeleton class="h-3 w-1/2" />
          </div>
        </div>
      ),
    },
    {
      name: "CardSkeleton",
      group: "Feedback",
      description: "Placeholder de carregamento em formato de card.",
      defaults: { compact: true },
      controls: [text("compact", "Compacto")],
      render: (v) => (
        <CardSkeleton
          class={v.compact ? "w-full [&>div:first-child]:h-12" : "w-full"}
        />
      ),
    },
    {
      name: "MultiStepStepper",
      group: "Navegação",
      description: "Indicador e navegação entre etapas.",
      defaults: { currentStep: "0" },
      controls: [select("currentStep", "Etapa atual", ["0", "1"])],
      render: (v, update) => (
        <MultiStepStepper
          steps={steps}
          currentStep={Number(v.currentStep)}
          onStepClick={(value) => update("currentStep", String(value))}
        />
      ),
    },
    {
      name: "MultiStepFields",
      group: "Formulários",
      description: "Renderização de campos a partir de configuração.",
      defaults: { label: "Nome do projeto", disabled: false },
      controls: [text("label", "Rótulo"), text("disabled", "Desabilitado")],
      render: (v) => (
        <MultiStepFields
          fields={[
            { name: "name", label: String(v.label), disabled: !!v.disabled },
          ]}
        />
      ),
    },
    {
      name: "MultiStepSummaryCard",
      group: "Layout",
      description: "Resumo de dados em pares de rótulo e valor.",
      defaults: { title: "Resumo", value: "Website" },
      controls: [text("title", "Título"), text("value", "Valor")],
      render: (v) => (
        <MultiStepSummaryCard
          summary={{
            title: String(v.title),
            items: [{ label: "Projeto", value: String(v.value) }],
          }}
        />
      ),
    },
    {
      name: "PaginatedContainer",
      group: "Layout",
      description: "Coleção com busca, ordenação e paginação.",
      defaults: { pageSize: "2" },
      controls: [select("pageSize", "Itens por página", ["1", "2", "3"])],
      render: (v) => (
        <PaginatedContainer
          items={[
            { name: "Website" },
            { name: "Aplicativo" },
            { name: "Design system" },
          ]}
          pageSize={Number(v.pageSize)}
          layout="list"
          renderItems={(items) => (
            <For each={items}>
              {(item) => <div class="sample-row">{item.name}</div>}
            </For>
          )}
        />
      ),
    },
    {
      name: "OverlayScrollbar",
      group: "Layout",
      description:
        "Barra de rolagem da página, como no Univents. Ativa em dispositivos com mouse.",
      defaults: { enabled: true },
      controls: [text("enabled", "Ativar na página")],
      render: (v) => (
        <div class="scrollbar-preview">
          <iframe
            title="Demonstração da OverlayScrollbar"
            src={`#/preview/scrollbar?enabled=${v.enabled}&dark=${dark()}`}
            class={
              detail()
                ? "scrollbar-frame"
                : "scrollbar-frame scrollbar-frame-compact"
            }
          />
          <p class="scrollbar-preview-caption">
            Role dentro da prévia para testar a barra.
          </p>
        </div>
      ),
    },
  ];
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
  const active = () => entries.find((entry) => entry.name === selected())!;
  const filtered = () =>
    entries.filter(
      (entry) =>
        (group() === "Todos" || entry.group === group()) &&
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
      <Show
        when={detail()}
        fallback={
          <div class="studio-body">
            <main class="catalog">
              <div class="intro">
                <h1>Componentes</h1>
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
              <div class="catalog-tabs">
                <SectionTabs
                  items={groups.map((id) => ({ id, label: id }))}
                  active={group()}
                  onChange={setGroup}
                  ariaLabel="Filtrar componentes"
                />
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
                        {entry.render(values[entry.name], (key, value) =>
                          update(entry.name, key, value),
                        )}
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
                      setGroup("Todos");
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
          </div>
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
  );
}
