import { interactionEntries } from "./examples/interaction-entries";
import { clientOnly } from "@solidjs/web";
import {
  PaginationExample,
  ToolbarExample,
  FilterBarExample,
  BulkActionsExample,
  InfiniteScrollExample,
  PaginatedExample,
  DataTableExample,
} from "./examples/CollectionExamples";
import { catalogSections, catalogTop } from "./catalog-navigation";
import { createSignal, Show } from "solid-js";
import type { JSX } from "@solidjs/web";
import {
  cn,
  shapes,
  type Shape,
  Avatar,
  Breadcrumb,
  Sidebar,
  Badge,
  Loading,
  Button,
  Card,
  Checkbox,
  Switch,
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
  Combobox,
  EmptyState,
  Skeleton,
  CardSkeleton,
  toast,
  MultiStepDialog,
  MultiStepStepper,
  MultiStepFields,
  MultiStepSummaryCard,
  type InputProps,
  type ButtonVariant,
  type ButtonSize,
  type DialogSize,
  type AlertModalVariant,
  Spinner,
  WaveText,
  WaveSpinnerLoading,
  DesktopOnly,
  CanvasControls,
  CatalogNav,
  Reveal,
  Separator,
  Tabs,
  Callout,
  CopyButton,
  Meter,
  PageHeader,
  PageSection,
  StatusDot,
  Table,
  ToggleGroup,
  type BadgeProps,
  type CalloutTone,
  type MeterTone,
  type StatusDotTone,
  type RevealDirection,
  type RevealVariant,
  type SeparatorLabelPosition,
  type SeparatorOrientation,
  type SeparatorVariant,
  type TabsActivation,
  type TabsOrientation,
  type TabsSize,
  type TabsVariant,
  type TextareaResize,
} from "@trieoh/ui-solid";
import IconArrowUpRight from "~icons/lucide/arrow-up-right";
import IconHome from "~icons/lucide/house";
import IconCalendar from "~icons/lucide/calendar-days";
import IconUsers from "~icons/lucide/users";
import IconPlus from "~icons/lucide/plus";
import IconMenu from "~icons/lucide/menu";
import { DeviceFrame } from "./components/DeviceFrame";

const StatCardExample = clientOnly(
  () => import("./examples/DashboardExamples"),
  { lazy: true, export: "StatCardExample" },
);
const MetricDeltaExample = clientOnly(
  () => import("./examples/DashboardExamples"),
  { lazy: true, export: "MetricDeltaExample" },
);
const SparklineExample = clientOnly(
  () => import("./examples/DashboardExamples"),
  { lazy: true, export: "SparklineExample" },
);
const BarListExample = clientOnly(
  () => import("./examples/DashboardExamples"),
  { lazy: true, export: "BarListExample" },
);
const PanelExample = clientOnly(() => import("./examples/DashboardExamples"), {
  lazy: true,
  export: "PanelExample",
});
const ChartExample = clientOnly(() => import("./examples/DashboardExamples"), {
  lazy: true,
  export: "ChartExample",
});

type BadgeVariant = NonNullable<BadgeProps["variant"]>;
const tableRows = [
  { name: "IdentityX", p95: 212, requests: 18_420 },
  { name: "Payssage", p95: 940, requests: 2_318 },
  { name: "Univents", p95: 388, requests: 7_904 },
];

export type Values = Record<string, string | boolean>;
export type NumberRange = {
  min?: number;
  max?: number;
  step?: number;
  /** An empty field means "not set". */
  optional?: boolean;
};
type Control = {
  key: string;
  label: string;
  options?: readonly string[];
  number?: NumberRange;
};
export type Entry = {
  /** Stable canvas for interactive miniature previews. */
  previewSize?: { width: number; height: number };
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
const number = (key: string, label: string, range: NumberRange): Control => ({
  key,
  label,
  number: range,
});
const select = (
  key: string,
  label: string,
  options: readonly string[],
): Control => ({
  key,
  label,
  options,
});

/**
 * Previews Reveal. The effect only plays on entry, so the example remounts it
 * (a) from the button and (b) whenever a knob changes. The swap is instant and
 * the stage has a fixed height: unmounting first and remounting later made the
 * centred preview collapse and jump, which read as a flash instead of the
 * animation.
 */
function RevealExample(props: { v: Values }) {
  const [run, setRun] = createSignal(0);
  const stagger = () =>
    Number(props.v.stagger) > 0 ? Number(props.v.stagger) : undefined;
  const token = () =>
    [
      run(),
      props.v.variant,
      props.v.direction,
      props.v.distance,
      props.v.duration,
      props.v.delay,
      props.v.stagger,
      props.v.once,
    ].join("|");

  return (
    <div class="flex w-full flex-col items-center gap-6">
      <Button size="sm" variant="outline" onClick={() => setRun(run() + 1)}>
        Reproduzir novamente
      </Button>
      <div class="min-h-52 w-full max-w-xs">
        <Show when={token()} keyed>
          {(_token) => (
            <Reveal
              class={String(props.v.class)}
              variant={props.v.variant as RevealVariant}
              direction={props.v.direction as RevealDirection}
              distance={Number(props.v.distance)}
              duration={Number(props.v.duration)}
              delay={Number(props.v.delay)}
              stagger={stagger()}
              once={!!props.v.once}
            >
              <Card class="space-y-1 p-4">
                <h3 class="text-sm font-semibold">Conteúdo animado</h3>
                <p class="text-xs text-muted-foreground">
                  {stagger() === undefined
                    ? "O bloco inteiro entra de uma vez."
                    : "Cada filho direto entra em sequência."}
                </p>
              </Card>
              <Show when={stagger() !== undefined}>
                <Card class="mt-2 p-4 text-sm">Segundo item</Card>
                <Card class="mt-2 p-4 text-sm">Terceiro item</Card>
              </Show>
            </Reveal>
          )}
        </Show>
      </div>
    </div>
  );
}

export function createCatalog(options: {
  isDetail: () => boolean;
  isDark: () => boolean;
}) {
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
  const stepperSteps = [
    {
      id: "details",
      title: "Detalhes",
      description: "Nome e informações do projeto.",
    },
    {
      id: "access",
      title: "Acesso",
      description: "Defina quem pode participar.",
    },
    {
      id: "review",
      title: "Revisão",
      description: "Confira os dados antes de criar.",
    },
  ];
  const entries: Entry[] = [
    ...interactionEntries.map(entry => ({ ...entry, defaults: { ...entry.defaults }, controls: [...entry.controls] })),
    {
      name: "Loading",
      group: "Feedback",
      description: "Indicador de carregamento com a animação do Univents.",
      defaults: {
        label: "Carregando…",
        showLabel: true,
        size: "md",
        tone: "accent",
        layout: "inline",
      },
      controls: [
        text("label", "Texto"),
        text("showLabel", "Mostrar texto"),
        select("size", "Tamanho", ["sm", "md", "lg"]),
        select("tone", "Cor", ["default", "accent", "muted"]),
        select("layout", "Disposição", ["inline", "stacked"]),
      ],
      render: (v) => (
        <Loading
          label={String(v.label)}
          showLabel={!!v.showLabel}
          size={v.size as "sm" | "md" | "lg"}
          tone={v.tone as "default" | "accent" | "muted"}
          layout={v.layout as "inline" | "stacked"}
          class={String(v.class)}
        />
      ),
    },
    {
      name: "Spinner",
      group: "Feedback",
      description: "Indicador circular com rastro de carregamento.",
      defaults: {
        size: "3rem",
      },
      controls: [
        select("size", "Tamanho", ["1.5rem", "2rem", "3rem", "4rem", "5rem", "6rem"]),
      ],
      render: (v) => (
        <div class="flex items-center justify-center p-6">
          <Spinner
            size={String(v.size)}
            activeColor="var(--primary)"
            trackColor="var(--accent)"
          />
        </div>
      ),
    },
    {
      name: "WaveText",
      group: "Feedback",
      description: "Texto ondulado caractere por caractere via requestAnimationFrame.",
      defaults: {
        text: "Processando...",
        duration: "1400",
        delay: "500",
        lift: "12",
        waveWidth: "20",
      },
      controls: [
        text("text", "Texto"),
        number("duration", "Duração (ms)", { min: 100, max: 10000, step: 100 }),
        number("delay", "Intervalo (ms)", { min: 0, max: 5000, step: 50 }),
        number("lift", "Elevação (px)", { min: 0, max: 60, step: 1 }),
        number("waveWidth", "Largura da onda (%)", { min: 1, max: 100, step: 1 }),
      ],
      render: (v) => (
        <div class="flex items-center justify-center p-6">
          <WaveText
            text={String(v.text)}
            duration={Number(v.duration) || 1400}
            delay={Number(v.delay) || 500}
            lift={Number(v.lift) || 12}
            waveWidth={Number(v.waveWidth) || 20}
          />
        </div>
      ),
    },
    {
      name: "WaveSpinnerLoading",
      group: "Feedback",
      description: "Carregamento composto por Spinner e texto com ondulação senoidal.",
      defaults: {
        text: "Processando pagamento...",
        spinnerSize: "5rem",
      },
      controls: [
        text("text", "Texto"),
        select("spinnerSize", "Tamanho do Spinner", ["3rem", "4rem", "5rem", "6rem"]),
      ],
      render: (v) => (
        <div class="flex w-full items-center justify-center p-6">
          <WaveSpinnerLoading
            text={String(v.text)}
            spinnerSize={String(v.spinnerSize)}
          />
        </div>
      ),
    },

    {
      name: "Avatar",
      group: "Layout",
      description:
        "Foto com iniciais como alternativa quando a imagem não está disponível.",
      defaults: { name: "Fulano de Tal", initials: "", src: "", size: "md" },
      controls: [
        text("name", "Nome"),
        text("initials", "Iniciais personalizadas"),
        text("src", "URL da imagem"),
        select("size", "Tamanho", ["sm", "md", "lg"]),
      ],
      render: (v) => (
        <Avatar
          class={String(v.class)}
          name={String(v.name)}
          initials={String(v.initials) || undefined}
          src={String(v.src) || undefined}
          size={v.size as "sm" | "md" | "lg"}
          shape={v.shape as Shape}
        />
      ),
    },
    {
      name: "Badge",
      group: "Feedback",
      description: "Rótulo de cargo, categoria ou estado.",
      defaults: { children: "Owner", variant: "accent" },
      controls: [
        text("children", "Texto"),
        select("variant", "Variante", [
          "default",
          "accent",
          "secondary",
          "outline",
          "muted",
          "success",
          "warning",
          "destructive",
        ]),
      ],
      render: (v) => (
        <Badge
          class={String(v.class)}
          variant={v.variant as BadgeVariant}
          shape={v.shape as Shape}
        >
          {String(v.children)}
        </Badge>
      ),
    },

    {
      name: "StatusDot",
      group: "Feedback",
      description:
        "Luz de estado compacta; idle é vazada, para quando nada está acontecendo.",
      defaults: { tone: "success", label: "Operacional", size: "sm" },
      controls: [
        select("tone", "Tom", ["success", "warning", "destructive", "muted", "idle"]),
        text("label", "Rótulo"),
        select("size", "Tamanho", ["sm", "md"]),
      ],
      render: (v) => (
        <span class="inline-flex items-center gap-1.5 text-xs text-muted-foreground">
          <StatusDot
            class={String(v.class)}
            tone={v.tone as StatusDotTone}
            label={String(v.label)}
            size={v.size as "sm" | "md"}
          />
          {String(v.label)}
        </span>
      ),
    },
    {
      name: "Callout",
      group: "Feedback",
      description: "Nota destacada dentro de uma página: contexto, dica, aviso ou perigo.",
      defaults: {
        tone: "warning",
        title: "Acorda a API",
        children: "Esta chamada liga o container se ele estiver dormindo.",
      },
      controls: [
        select("tone", "Tom", ["info", "success", "warning", "destructive"]),
        text("title", "Título"),
        text("children", "Conteúdo"),
      ],
      render: (v) => (
        <Callout
          class={cn("w-full max-w-md", String(v.class))}
          shape={v.shape as Shape}
          tone={v.tone as CalloutTone}
          title={String(v.title) || undefined}
        >
          <p>{String(v.children)}</p>
        </Callout>
      ),
    },
    {
      name: "Meter",
      group: "Feedback",
      description:
        "Quanto de algo já foi usado, com marcas opcionais para limites ou ritmo esperado.",
      defaults: {
        value: "62",
        label: "Horas de container",
        tone: "success",
        showMarkers: true,
      },
      controls: [
        number("value", "Valor (%)", { min: 0, max: 100 }),
        text("label", "Rótulo"),
        select("tone", "Tom", ["default", "accent", "success", "warning", "destructive"]),
        text("showMarkers", "Mostrar marcas"),
      ],
      render: (v) => (
        <Meter
          class={cn("w-full max-w-sm", String(v.class))}
          value={Number(v.value)}
          max={100}
          label={String(v.label)}
          valueLabel={`${Number(v.value)}%`}
          tone={v.tone as MeterTone}
          markers={
            v.showMarkers
              ? [
                  { value: 50, label: "Ritmo esperado" },
                  { value: 90, label: "Limite", tone: "destructive" },
                ]
              : []
          }
        />
      ),
    },
    {
      name: "CopyButton",
      group: "Ações",
      description: "Copia um texto para a área de transferência e confirma com um check.",
      defaults: { text: "pnpm add @trieoh/ui-solid", size: "md" },
      controls: [text("text", "Texto"), select("size", "Tamanho", ["sm", "md"])],
      render: (v) => (
        <div class="flex items-center gap-2 rounded-md border border-border px-3 py-1 font-mono text-xs">
          <span>{String(v.text)}</span>
          <CopyButton
            class={String(v.class)}
            shape={v.shape as Shape}
            text={String(v.text)}
            size={v.size as "sm" | "md"}
          />
        </div>
      ),
    },
    {
      name: "ToggleGroup",
      group: "Formulários",
      description: "Linha de opções em que uma fica marcada, como um filtro.",
      defaults: { value: "all", size: "sm" },
      controls: [select("size", "Tamanho", ["sm", "md"])],
      render: (v, update) => (
        <ToggleGroup
          class={String(v.class)}
          shape={v.shape as Shape}
          ariaLabel="Filtrar por serviço"
          size={v.size as "sm" | "md"}
          value={String(v.value)}
          onChange={(value) => update("value", value)}
          items={[
            { value: "all", label: "Todos" },
            { value: "identityx", label: "IdentityX" },
            { value: "payssage", label: "Payssage" },
            { value: "univents", label: "Univents" },
          ]}
        />
      ),
    },
    {
      name: "PageHeader",
      group: "Layout",
      description: "Bloco de título no topo de uma página, com descrição e ações.",
      defaults: {
        title: "Entrega",
        description: "Deploys, pull requests e CI dos últimos 30 dias.",
        showActions: false,
      },
      controls: [
        text("title", "Título"),
        text("description", "Descrição"),
        text("showActions", "Mostrar ações"),
      ],
      render: (v) => (
        <PageHeader
          class={cn("w-full max-w-lg", String(v.class))}
          title={String(v.title)}
          description={String(v.description) || undefined}
          actions={
            v.showActions ? (
              <Button size="sm" variant="outline">
                Exportar
              </Button>
            ) : undefined
          }
        />
      ),
    },
    {
      name: "PageSection",
      group: "Layout",
      description: "Bloco com título dentro de uma página, com uma nota à direita.",
      defaults: { title: "Deploys", aside: "últimos 30 dias" },
      controls: [text("title", "Título"), text("aside", "Nota")],
      render: (v) => (
        <PageSection
          class={cn("mb-0 w-full max-w-lg", String(v.class))}
          title={String(v.title)}
          aside={String(v.aside) || undefined}
        >
          <p class="text-sm text-muted-foreground">Conteúdo da seção.</p>
        </PageSection>
      ),
    },
    {
      name: "Table",
      group: "Layout",
      description:
        "Tabela simples só de leitura; para ordenar, selecionar ou paginar use DataTable.",
      previewSize: { width: 520, height: 220 },
      defaults: { density: "default", empty: false },
      controls: [
        select("density", "Densidade", ["default", "compact"]),
        text("empty", "Sem dados"),
      ],
      render: (v) => (
        <Table
          class={cn("w-full max-w-lg", String(v.class))}
          caption="Serviços"
          density={v.density as "default" | "compact"}
          items={v.empty ? [] : tableRows}
          emptyState={
            <EmptyState variant="inline" title="Nenhum serviço ainda." />
          }
          columns={[
            {
              id: "name",
              header: "Serviço",
              cell: (row) => row.name,
              class: "font-medium",
            },
            {
              id: "p95",
              header: "p95",
              align: "right",
              cell: (row) => `${row.p95} ms`,
              class: (row) => (row.p95 > 800 ? "tabular-nums text-warning" : "tabular-nums"),
            },
            {
              id: "requests",
              header: "Requisições",
              align: "right",
              cell: (row) => row.requests.toLocaleString("pt-BR"),
              class: "tabular-nums",
            },
          ]}
        />
      ),
    },

    {
      name: "Switch",
      group: "Formulários",
      description: "Alterna entre ligado e desligado.",
      defaults: { checked: true, disabled: false },
      controls: [text("checked", "Ligado"), text("disabled", "Desabilitado")],
      render: (v, update) => (
        <Switch
          class={String(v.class)}
          aria-label="Receber notificações"
          checked={!!v.checked}
          disabled={!!v.disabled}
          onChange={(checked) => update("checked", checked)}
        />
      ),
    },

    {
      name: "Checkbox",
      group: "Formulários",
      description: "Seleção booleana com rótulo e descrição.",
      defaults: {
        label: "Receber notificações",
        description: "Novidades e atualizações dos seus projetos.",
        checked: true,
        position: "right",
        alignment: "center",
        size: "md",
        variant: "card",
        controlShape: "subtle",
        indeterminate: false,
        required: false,
        name: "notifications",
        value: "yes",
        disabled: false,
      },
      controls: [
        text("label", "Rótulo"),
        text("description", "Descrição"),
        text("checked", "Marcado"),
        text("disabled", "Desabilitado"),
        select("position", "Posição", ["left", "right"]),
        select("alignment", "Alinhamento", ["start", "center", "end"]),
        select("size", "Tamanho", ["sm", "md", "lg"]),
        select("variant", "Variante", ["card", "plain"]),
        select("controlShape", "Shape do controle", shapes),
        text("indeterminate", "Estado misto"),
        text("required", "Obrigatório"),
        text("name", "Nome no formulário"),
        text("value", "Valor no formulário"),
      ],
      render: (v, update) => (
        <Checkbox
          class={String(v.class)}
          shape={v.shape as Shape}
          label={String(v.label)}
          description={String(v.description)}
          position={v.position as "left" | "right"}
          alignment={v.alignment as "start" | "center" | "end"}
          size={v.size as "sm" | "md" | "lg"}
          variant={v.variant as "card" | "plain"}
          controlShape={v.controlShape as Shape}
          indeterminate={!!v.indeterminate}
          required={!!v.required}
          name={String(v.name)}
          value={String(v.value)}
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
      defaults: {
        placeholder: "Buscar projetos…",
        value: "",
        disabled: false,
        readonly: false,
        required: false,
        name: "search",
      },
      controls: [
        text("placeholder", "Placeholder"),
        text("disabled", "Desabilitado"),
        text("readonly", "Somente leitura"),
        text("required", "Obrigatório"),
        text("name", "Nome"),
      ],
      render: (v, update) => (
        <SearchInput
          class={String(v.class)}
          shape={v.shape as Shape}
          required={!!v.required}
          name={String(v.name)}
          disabled={!!v.disabled}
          readonly={!!v.readonly}
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
      defaults: { children: "Conteúdo do card", href: "", external: false },
      controls: [
        text("children", "Conteúdo"),
        text("href", "Link"),
        text("external", "Abrir em nova aba"),
      ],
      render: (v) => (
        <Card
          href={String(v.href) || undefined}
          external={!!v.external}
          shape={v.shape as Shape}
          class={cn("w-full p-5", String(v.class))}
        >
          {String(v.children)}
        </Card>
      ),
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
      defaults: {
        code: 'import { Button } from "@trieoh/ui-solid";\n\nexport function SaveButton() {\n  return (\n    <Button onClick={() => console.log("Salvo")}>\n      Salvar projeto\n    </Button>\n  );\n}',
        filename: "save-button.tsx",
        showLineNumbers: true,
        format: true,
      },
      controls: [
        text("code", "Código"),
        text("filename", "Nome do arquivo"),
        text("showLineNumbers", "Números de linha"),
        text("format", "Formatar código"),
      ],
      render: (v) => (
        <CodeBlock
          class={String(v.class)}
          shape={v.shape as Shape}
          code={String(v.code)}
          showLineNumbers={!!v.showLineNumbers}
          format={!!v.format}
          filename={String(v.filename)}
        />
      ),
    },
    {
      name: "Button",
      group: "Ações",
      description: "Botão com variantes, tamanhos e estado desabilitado.",
      defaults: {
        children: "Criar projeto",
        variant: "default",
        size: "md",
        type: "button",
        title: "",
        name: "action",
        value: "save",
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
        select("type", "Tipo HTML", ["button", "submit", "reset"]),
        text("title", "Descrição (title)"),
        text("name", "Nome"),
        text("value", "Valor"),
      ],
      render: (v) => (
        <Button
          class={String(v.class)}
          shape={v.shape as Shape}
          type={v.type as "button" | "submit" | "reset"}
          title={String(v.title)}
          name={String(v.name)}
          value={String(v.value)}
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
      name: "CanvasControls",
      group: "Ações",
      description: "Barra flutuante de zoom (+, -, %) e ajuste para editores de canvas.",
      defaults: {
        scale: "1",
      },
      controls: [
        number("scale", "Escala", { min: 0.25, max: 4, step: 0.05 }),
      ],
      render: (v, update) => {
        const scaleVal = () => Number(v.scale) || 1;
        return (
          <div class="relative h-44 w-full rounded-xl border border-border bg-muted/20 flex items-center justify-center overflow-hidden">
            <div
              style={{ transform: `scale(${scaleVal()})` }}
              class="px-4 py-2 bg-card border border-border rounded-lg shadow-sm text-xs font-medium transition-transform"
            >
              Elemento no Canvas ({Math.round(scaleVal() * 100)}%)
            </div>
            <CanvasControls
              scale={scaleVal()}
              onZoomIn={() => update("scale", String(Math.min(scaleVal() + 0.1, 2).toFixed(1)))}
              onZoomOut={() => update("scale", String(Math.max(scaleVal() - 0.1, 0.2).toFixed(1)))}
              onReset={() => update("scale", "1")}
            />
          </div>
        );
      },
    },
    {
      name: "Input",
      group: "Formulários",
      description: "Campo de texto de uma linha.",
      defaults: {
        placeholder: "Seu e-mail",
        disabled: false,
        type: "text",
        required: false,
        readonly: false,
        value: "",
        name: "example",
        min: "",
        max: "",
        step: "",
        minlength: "",
        maxlength: "",
        autocomplete: "off",
        inputmode: "text",
      },
      controls: [
        text("placeholder", "Placeholder"),
        text("disabled", "Desabilitado"),
        select("type", "Tipo", [
          "text",
          "email",
          "password",
          "search",
          "number",
          "url",
          "tel",
          "date",
          "datetime-local",
        ]),
        text("required", "Obrigatório"),
        text("readonly", "Somente leitura"),
        text("name", "Nome"),
        text("value", "Valor"),
        text("min", "Mínimo"),
        text("max", "Máximo"),
        text("step", "Passo"),
        select("minlength", "Caracteres mínimos", ["", "2", "8"]),
        select("maxlength", "Caracteres máximos", ["", "16", "32", "120"]),
        text("autocomplete", "Autocomplete"),
        select("inputmode", "Teclado virtual", [
          "text",
          "email",
          "tel",
          "numeric",
          "decimal",
          "url",
          "search",
          "none",
        ]),
      ],
      render: (v, update) => (
        <Input
          class={String(v.class)}
          shape={v.shape as Shape}
          type={v.type as InputProps["type"]}
          min={String(v.min) || undefined}
          max={String(v.max) || undefined}
          step={String(v.step) || undefined}
          minlength={v.minlength ? Number(v.minlength) : undefined}
          maxlength={v.maxlength ? Number(v.maxlength) : undefined}
          autocomplete={String(v.autocomplete)}
          inputmode={v.inputmode as InputProps["inputmode"]}
          required={!!v.required}
          readonly={!!v.readonly}
          name={String(v.name)}
          value={String(v.value)}
          onInput={(e) => update("value", e.currentTarget.value)}
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
        required: false,
        readonly: false,
        rows: "3",
        value: "",
        name: "description",
        autosize: false,
        maxRows: "6",
        showCount: false,
        maxlength: "",
        resize: "",
      },
      controls: [
        text("placeholder", "Placeholder"),
        text("disabled", "Desabilitado"),
        text("required", "Obrigatório"),
        text("readonly", "Somente leitura"),
        select("rows", "Linhas", ["2", "3", "4", "6"]),
        text("autosize", "Cresce com o conteúdo"),
        select("maxRows", "Máx. de linhas (autosize)", ["3", "6", "10"]),
        text("showCount", "Mostrar contador (automático com limite)"),
        number("maxlength", "Limite de caracteres", { min: 1, max: 2000, step: 1, optional: true }),
        select("resize", "Redimensionar", ["", "vertical", "none", "horizontal", "both"]),
        text("value", "Valor"),
        text("name", "Nome"),
      ],
      render: (v, update) => (
        <Textarea
          class={String(v.class)}
          shape={v.shape as Shape}
          placeholder={String(v.placeholder)}
          disabled={!!v.disabled}
          name={String(v.name)}
          rows={Number(v.rows)}
          required={!!v.required}
          readonly={!!v.readonly}
          autosize={!!v.autosize}
          maxRows={Number(v.maxRows)}
          showCount={v.showCount ? true : undefined}
          maxlength={Number(v.maxlength) > 0 ? Number(v.maxlength) : undefined}
          resize={(String(v.resize) || undefined) as TextareaResize | undefined}
          value={String(v.value)}
          onInput={(e) => update("value", e.currentTarget.value)}
          aria-label="Exemplo de área de texto"
        />
      ),
    },
    {
      name: "Label",
      group: "Formulários",
      description: "Rótulo associado a um campo de formulário.",
      defaults: { children: "Nome do projeto", for: "label-example-input" },
      controls: [
        text("children", "Texto"),
        text("for", "Campo associado (for)"),
      ],
      render: (v) => (
        <div class="w-full space-y-2">
          <Label class={String(v.class)} for={String(v.for)}>
            {String(v.children)}
          </Label>
          <Input id={String(v.for)} placeholder="Clique no rótulo para focar" />
        </div>
      ),
    },
    {
      name: "Field",
      group: "Formulários",
      description: "Composição de rótulo, descrição e erro de validação.",
      defaults: {
        label: "Nome completo",
        hint: "Como você gostaria de ser chamado?",
        error: "",
        required: false,
        name: "full-name",
      },
      controls: [
        text("label", "Rótulo"),
        text("hint", "Ajuda"),
        text("error", "Erro"),
        text("required", "Obrigatório"),
        text("name", "Nome do campo"),
      ],
      render: (v) => (
        <Field
          class={String(v.class)}
          name={String(v.name)}
          required={!!v.required}
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
      defaults: {
        placeholder: "Selecione uma cidade",
        value: "",
        disabled: false,
        searchPlaceholder: "Buscar…",
        emptyMessage: "Nenhum resultado encontrado",
      },
      controls: [
        text("placeholder", "Placeholder"),
        text("disabled", "Desabilitado"),
        text("searchPlaceholder", "Placeholder de busca"),
        text("emptyMessage", "Mensagem sem resultados"),
        select("value", "Valor", ["", "sp", "rj", "bh"]),
      ],
      render: (v, update) => (
        <Combobox
          class={String(v.class)}
          shape={v.shape as Shape}
          disabled={!!v.disabled}
          searchPlaceholder={String(v.searchPlaceholder)}
          emptyMessage={String(v.emptyMessage)}
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
          class={String(v.class)}
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
    {
      name: "Reveal",
      group: "Layout",
      description:
        "Entrada animada quando o conteúdo aparece na tela, com Web Animations nativo e um único observador compartilhado.",
      defaults: {
        variant: "slide",
        direction: "up",
        distance: "15",
        duration: "0.4",
        delay: "0",
        stagger: "0",
        once: true,
      },
      controls: [
        select("variant", "Variante", ["slide", "fade", "scale", "blur"]),
        select("direction", "Direção (slide)", ["up", "down", "left", "right"]),
        number("distance", "Distância (px)", { min: 0, max: 200, step: 1 }),
        number("duration", "Duração (s)", { min: 0.05, max: 5, step: 0.05 }),
        number("delay", "Atraso (s)", { min: 0, max: 5, step: 0.05 }),
        number("stagger", "Intervalo entre filhos (s, 0 = bloco único)", { min: 0, max: 1, step: 0.01 }),
        text("once", "Animar só na primeira vez"),
      ],
      render: (v) => <RevealExample v={v} />,
    },
    {
      name: "Tabs",
      group: "Navegação",
      description:
        "Abas acessíveis com painéis, navegação por setas e três variantes visuais.",
      defaults: {
        variant: "underline",
        size: "md",
        orientation: "horizontal",
        activation: "automatic",
        fullWidth: false,
        panels: true,
        keepMounted: false,
        count: "8",
      },
      controls: [
        select("variant", "Variante", ["underline", "pills", "segmented"]),
        select("size", "Tamanho", ["sm", "md"]),
        select("orientation", "Orientação", ["horizontal", "vertical"]),
        select("activation", "Ativação", ["automatic", "manual"]),
        text("fullWidth", "Ocupar toda a largura"),
        text("panels", "Painéis embutidos (desligado: só o seletor)"),
        text("keepMounted", "Manter painéis montados"),
        select("count", "Contagem", ["0", "8", "24"]),
      ],
      render: (v) => {
        const panel = (text: string) =>
          v.panels ? (
            <p class="py-2 text-sm text-muted-foreground">{text}</p>
          ) : undefined;

        return (
        <Tabs
          class={String(v.class)}
          ariaLabel="Exemplo de abas"
          variant={v.variant as TabsVariant}
          size={v.size as TabsSize}
          orientation={v.orientation as TabsOrientation}
          activation={v.activation as TabsActivation}
          fullWidth={!!v.fullWidth}
          keepMounted={!!v.keepMounted}
          defaultValue="overview"
          items={[
            {
              id: "overview",
              label: "Visão geral",
              content: panel("Resumo do projeto: use as setas para trocar de aba."),
            },
            {
              id: "activity",
              label: "Atividade",
              count: Number(v.count),
              content: panel("Últimas alterações feitas pela equipe."),
            },
            {
              id: "settings",
              label: "Ajustes",
              content: panel("Preferências e permissões."),
            },
            {
              id: "billing",
              label: "Cobrança",
              disabled: true,
              content: panel("Indisponível."),
            },
          ]}
        />
        );
      },
    },
    {
      name: "Separator",
      group: "Layout",
      description:
        "Linha divisória horizontal ou vertical, com rótulo opcional e três estilos.",
      defaults: {
        orientation: "horizontal",
        variant: "solid",
        label: "",
        labelPosition: "center",
        decorative: true,
      },
      controls: [
        select("orientation", "Orientação", ["horizontal", "vertical"]),
        select("variant", "Estilo", ["solid", "dashed", "dotted"]),
        text("label", "Rótulo"),
        select("labelPosition", "Posição do rótulo", ["start", "center", "end"]),
        text("decorative", "Decorativo (ignorado por leitores de tela)"),
      ],
      render: (v) => (
        <div
          data-wide
          class={
            v.orientation === "vertical"
              ? "flex h-40 items-stretch justify-center gap-6 text-sm"
              : "flex w-full max-w-md flex-col gap-3 text-sm"
          }
        >
          <span>Entrar</span>
          <Separator
            class={String(v.class)}
            orientation={v.orientation as SeparatorOrientation}
            variant={v.variant as SeparatorVariant}
            label={String(v.label) || undefined}
            labelPosition={v.labelPosition as SeparatorLabelPosition}
            decorative={!!v.decorative}
          />
          <span>Criar conta</span>
        </div>
      ),
    },
    ...(["Dialog", "AlertModal", "Drawer", "MultiStepDialog"] as const).map(
      (name) => ({
        name,
        group: "Overlays",
        description: {
          Dialog: "Modal com título, conteúdo e tamanhos configuráveis.",
          AlertModal: "Modal de confirmação com variantes semânticas.",
          Drawer: "Painel sobreposto que abre pela base ou pela lateral da tela.",
          MultiStepDialog: "Formulário dividido em etapas.",
        }[name],
        defaults: {
          title: name === "AlertModal" ? "Excluir projeto?" : "Um novo projeto",
          size: "md",
          variant: "warning",
          description: "Ajuste os dados antes de continuar.",
          closable: true,
          preventClose: false,
          loading: false,
          confirmLabel: "Confirmar",
          cancelLabel: "Cancelar",
          fixedHeight: true,
          nextLabel: "Continuar",
          submitLabel: "Concluir",
          side: "bottom",
        },
        controls: [
          text("title", "Título"),
          text("description", "Descrição"),
          ...(name === "Dialog"
            ? [
              text("closable", "Botão de fechar"),
              text("preventClose", "Impedir fechamento externo"),
            ]
            : []),
          ...(name === "AlertModal"
            ? [
              text("loading", "Carregando"),
              text("confirmLabel", "Texto de confirmação"),
              text("cancelLabel", "Texto de cancelar"),
            ]
            : []),
          ...(name === "MultiStepDialog"
            ? [
              text("loading", "Carregando"),
              text("fixedHeight", "Altura fixa"),
              text("nextLabel", "Texto de avançar"),
              text("submitLabel", "Texto de concluir"),
            ]
            : []),
          ...(name === "Dialog" || name === "MultiStepDialog"
            ? [select("size", "Tamanho", ["sm", "md", "lg", "xl", "full"])]
            : []),
          ...(name === "Drawer"
            ? [select("side", "Lado", ["bottom", "left", "right"])]
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
                class={String(v.class)}
                shape={v.shape as Shape}
                open={opened() === name}
                onOpenChange={(value) => setOpened(value ? name : "")}
                title={String(v.title)}
                size={v.size as DialogSize}
                description={String(v.description)}
                closable={!!v.closable}
                preventClose={!!v.preventClose}
              >
                <p class="text-sm text-muted-foreground">
                  Este é o conteúdo do seu dialog.
                </p>
                <Button
                  class="mt-5"
                  onClick={() => setOpened("")}
                >
                  Concluir
                </Button>
              </Dialog>
            </Show>
            <Show when={name === "AlertModal"}>
              <AlertModal
                class={String(v.class)}
                shape={v.shape as Shape}
                open={opened() === name}
                onOpenChange={(value) => setOpened(value ? name : "")}
                title={String(v.title)}
                description={String(v.description)}
                loading={!!v.loading}
                children={
                  v.loading ? (
                    <Button variant="outline" onClick={() => setOpened("")}>
                      Encerrar simulação
                    </Button>
                  ) : undefined
                }
                confirmLabel={String(v.confirmLabel)}
                cancelLabel={String(v.cancelLabel)}
                variant={v.variant as AlertModalVariant}
                onConfirm={() => {
                  setOpened("");
                }}
              />
            </Show>
            <Show when={name === "Drawer"}>
              <Drawer
                side={v.side as "bottom" | "left" | "right"}
                class={String(v.class)}
                shape={v.shape as Shape}
                open={opened() === name}
                onClose={() => setOpened("")}
                title={String(v.title)}
              >
                <p class="text-sm text-muted-foreground">
                  {String(v.description)}
                </p>
              </Drawer>
            </Show>
            <Show when={name === "MultiStepDialog"}>
              <MultiStepDialog
                class={String(v.class)}
                shape={v.shape as Shape}
                open={opened() === name}
                onOpenChange={(value) => setOpened(value ? name : "")}
                title={String(v.title)}
                size={v.size as DialogSize}
                description={String(v.description)}
                loading={!!v.loading}
                fixedHeight={!!v.fixedHeight}
                nextLabel={String(v.nextLabel)}
                submitLabel={String(v.submitLabel)}
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
      defaults: {
        label: "Informações adicionais sobre esta ação.",
        placement: "top",
        align: "center",
      },
      controls: [
        text("label", "Mensagem"),
        select("placement", "Posição", ["top", "bottom"]),
        select("align", "Alinhamento", ["start", "center", "end"]),
      ],
      render: (v) => (
        <Tooltip
          shape={v.shape as Shape}
          class={String(v.class)}
          placement={v.placement as "top" | "bottom"}
          align={v.align as "start" | "center" | "end"}
          label={String(v.label)}
        >
          <Button variant="outline">Passe o mouse ou foque aqui</Button>
        </Tooltip>
      ),
    },
    {
      name: "Toaster",
      group: "Feedback",
      description: "Notificações de sucesso, erro, aviso e informação.",
      defaults: {
        message: "Projeto salvo com sucesso",
        type: "success",
        description: "",
        duration: "4000",
        position: "bottom-right",
        dismissible: true,
      },
      controls: [
        text("message", "Mensagem"),
        text("description", "Descrição"),
        select("position", "Posição", [
          "top-left",
          "top-center",
          "top-right",
          "bottom-left",
          "bottom-center",
          "bottom-right",
        ]),
        select("duration", "Duração (ms)", ["2000", "4000", "8000"]),
        text("dismissible", "Permitir fechar"),
        select("type", "Tipo", ["success", "error", "warning", "info"]),
      ],
      render: (v) => (
        <Button
          variant="outline"
          onClick={() =>
            toast[v.type as "success" | "error" | "warning" | "info"](
              String(v.message),
              {
                description: String(v.description),
                duration: Number(v.duration),
                position: v.position as
                  | "top-left"
                  | "top-center"
                  | "top-right"
                  | "bottom-left"
                  | "bottom-center"
                  | "bottom-right",
                dismissible: !!v.dismissible,
              },
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
        eyebrow: "Projetos",
        actionLabel: "Criar projeto",
        showAction: false,
        variant: "default",
      },
      controls: [
        select("variant", "Variante", ["default", "inline"]),
        text("title", "Título"),
        text("description", "Descrição"),
        text("eyebrow", "Texto superior"),
        text("showAction", "Mostrar ação"),
        text("actionLabel", "Texto da ação"),
      ],
      render: (v) => (
        <EmptyState
          variant={v.variant as "default" | "inline"}
          shape={v.shape as Shape}
          eyebrow={String(v.eyebrow)}
          action={
            v.showAction ? (
              <Button onClick={() => toast.info("Ação de exemplo")}>
                {String(v.actionLabel)}
              </Button>
            ) : undefined
          }
          title={String(v.title)}
          description={String(v.description)}
          class={cn("w-full py-6! px-4!", String(v.class))}
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
            shape={v.shape as Shape}
            class={cn(
              v.rounded ? "size-11 shrink-0 rounded-full" : "size-11 shrink-0",
              String(v.class),
            )}
          />
          <div class="flex-1 space-y-2">
            <Skeleton shape={v.shape as Shape} class="h-3 w-3/4" />
            <Skeleton shape={v.shape as Shape} class="h-3 w-1/2" />
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
          shape={v.shape as Shape}
          class={cn(
            v.compact ? "w-full [&>div:first-child]:h-12" : "w-full",
            String(v.class),
          )}
        />
      ),
    },
    {
      name: "Sidebar",
      group: "Navegação",
      description: "Painel lateral compacto, com seleção discreta, seções e modo recolhido.",
      defaults: {
        title: "Workspace",
        subtitle: "Conta pessoal",
        collapsed: false,
        activeId: "overview",
        mobile: false,
      },
      controls: [
        text("mobile", "Forçar mobile"),
        text("title", "Título"),
        text("subtitle", "Subtítulo"),
        text("collapsed", "Recolhido"),
        select("activeId", "Página ativa", ["overview", "events", "members"]),
      ],
      render: (v, update) => {
        const [menuOpen, setMenuOpen] = createSignal(false);
        const mobile = () => !!v.mobile;
        const titles: Record<string, string> = {
          overview: "Visão geral",
          events: "Eventos",
          members: "Membros",
        };
        const panel = () => (
          <Sidebar
            class={cn("bg-card", String(v.class))}
            brand={<IconHome class="size-5" />}
            title={String(v.title)}
            subtitle={String(v.subtitle)}
            collapsed={mobile() ? false : !!v.collapsed}
            onCollapsedChange={(value) =>
              mobile() ? setMenuOpen(false) : update("collapsed", value)
            }
            activeId={String(v.activeId)}
            onSelect={(item) => {
              update("activeId", item.id);
              setMenuOpen(false);
            }}
            sections={[{
              title: "Workspace",
              items: [
                { id: "overview", label: "Visão geral", icon: IconHome },
                { id: "events", label: "Eventos", icon: IconCalendar, badge: 8 },
                { id: "members", label: "Membros", icon: IconUsers, badge: 24 },
              ],
            }]}
            action={
              <Button
                variant="outline"
                class={cn("rounded-full shadow-none", !mobile() && v.collapsed ? "size-11 p-0" : "w-full justify-start gap-3")}
                aria-label="Novo projeto"
                onClick={() => toast.info("Criar novo projeto")}
              >
                <IconPlus class="size-5 shrink-0" />
                <Show when={mobile() || !v.collapsed}>Novo projeto</Show>
              </Button>
            }
            footer={<p class={cn("text-xs text-muted-foreground", !mobile() && v.collapsed ? "text-center" : "px-3")}>{!mobile() && v.collapsed ? "P" : "Plano pessoal"}</p>}
          />
        );
        const content = () => (
          <div class="min-w-0 flex-1 p-6">
            <p class="text-xs text-muted-foreground">Seu espaço de trabalho</p>
            <h3 class="mt-2 text-lg font-semibold tracking-tight">{titles[String(v.activeId)]}</h3>
            <div class="mt-6 space-y-3" aria-hidden="true">
              <div class="h-24 rounded-2xl border border-border/60 bg-muted/30" />
              <div class="h-3 w-3/4 rounded-full bg-muted" />
              <div class="h-3 w-1/2 rounded-full bg-muted" />
            </div>
          </div>
        );
        return (
          <DeviceFrame
            width={mobile() ? 390 : 1024}
            height={mobile() ? 844 : 640}
            fixedWidth={options.isDetail() ? undefined : 560}
          >
            <Show
              when={mobile()}
              fallback={
                <div class="flex h-full w-full">
                  {panel()}
                  {content()}
                </div>
              }
            >
              <div class="relative flex h-full w-full flex-col overflow-hidden">
                <header class="flex h-16 shrink-0 items-center gap-2 border-b border-border/60 bg-card/95 px-3 shadow-sm shadow-black/5">
                  <button
                    type="button"
                    aria-label="Abrir menu"
                    onClick={() => setMenuOpen(true)}
                    class="inline-flex size-10 cursor-pointer items-center justify-center rounded-xl text-foreground transition-colors hover:bg-muted/60"
                  >
                    <IconMenu class="size-5" />
                  </button>
                  <span class="truncate text-sm font-semibold">{String(v.title)}</span>
                </header>
                {content()}
                <button
                  type="button"
                  aria-label="Fechar menu"
                  tabindex={menuOpen() ? 0 : -1}
                  onClick={() => setMenuOpen(false)}
                  class={cn(
                    "absolute inset-0 z-10 bg-black/25 transition-opacity duration-200",
                    menuOpen() ? "opacity-100" : "pointer-events-none invisible opacity-0",
                  )}
                />
                <div
                  inert={!menuOpen()}
                  class={cn(
                    "absolute inset-y-0 left-0 z-20 h-full transition-transform duration-200 motion-reduce:transition-none",
                    menuOpen() ? "translate-x-0" : "-translate-x-full",
                  )}
                >
                  {panel()}
                </div>
              </div>
            </Show>
          </DeviceFrame>
        );
      },
    },
    {
      name: "Breadcrumb",
      group: "Navegação",
      description: "Trilha de navegação com links e indicação da página atual.",
      defaults: { currentLabel: "Edições", ariaLabel: "Trilha de navegação" },
      controls: [text("currentLabel", "Página atual"), text("ariaLabel", "Rótulo acessível")],
      render: (v) => (
        <Breadcrumb
          class={cn("w-full", String(v.class))}
          ariaLabel={String(v.ariaLabel)}
          items={[
            { path: "#admin", label: "Admin" },
            { path: "#events", label: "Eventos" },
            { path: "#editions", label: String(v.currentLabel) },
          ]}
        />
      ),
    },
    {
      name: "MultiStepStepper",
      group: "Navegação",
      description: "Indicador e navegação entre etapas.",
      defaults: {
        currentStep: "1",
        orientation: "auto",
        showDescriptions: false,
        interactive: true,
        ariaLabel: "Etapas do projeto",
      },
      controls: [
        select("currentStep", "Etapa atual", ["0", "1", "2"]),
        select("orientation", "Orientação", ["auto", "horizontal", "vertical"]),
        text("showDescriptions", "Mostrar descrições"),
        text("interactive", "Voltar às etapas concluídas"),
        text("ariaLabel", "Rótulo acessível"),
      ],
      render: (v, update) => (
        <MultiStepStepper
          class={String(v.class)}
          steps={stepperSteps}
          orientation={v.orientation as "auto" | "horizontal" | "vertical"}
          showDescriptions={!!v.showDescriptions}
          shape={v.shape as Shape}
          ariaLabel={String(v.ariaLabel)}
          currentStep={Number(v.currentStep)}
          onStepClick={
            v.interactive
              ? (value) => update("currentStep", String(value))
              : undefined
          }
        />
      ),
    },
    {
      name: "MultiStepFields",
      group: "Formulários",
      description: "Renderização de campos a partir de configuração.",
      defaults: {
        label: "Nome do projeto",
        disabled: false,
        required: false,
        kind: "text",
        placeholder: "Digite aqui…",
        hint: "",
        layout: "full",
      },
      controls: [
        text("label", "Rótulo"),
        text("disabled", "Desabilitado"),
        text("required", "Obrigatório"),
        select("kind", "Tipo do campo", [
          "text",
          "email",
          "number",
          "textarea",
        ]),
        text("placeholder", "Placeholder"),
        text("hint", "Ajuda"),
        select("layout", "Layout", ["full", "half"]),
      ],
      render: (v) => (
        <MultiStepFields
          class={String(v.class)}
          fields={[
            {
              name: "name",
              label: String(v.label),
              disabled: !!v.disabled,
              required: !!v.required,
              kind: v.kind as "text" | "email" | "number" | "textarea",
              placeholder: String(v.placeholder),
              hint: String(v.hint),
              layout: v.layout as "full" | "half",
            },
          ]}
        />
      ),
    },
    {
      name: "MultiStepSummaryCard",
      group: "Layout",
      description: "Resumo de dados em pares de rótulo e valor.",
      defaults: {
        title: "Resumo",
        value: "Website",
        badge: "",
        mono: false,
        fullWidth: false,
      },
      controls: [
        text("title", "Título"),
        text("value", "Valor"),
        text("badge", "Badge"),
        text("mono", "Monoespaçado"),
        text("fullWidth", "Largura total"),
      ],
      render: (v) => (
        <MultiStepSummaryCard
          shape={v.shape as Shape}
          class={String(v.class)}
          summary={{
            title: String(v.title),
            badge: String(v.badge),
            items: [
              {
                label: "Projeto",
                value: String(v.value),
                mono: !!v.mono,
                fullWidth: !!v.fullWidth,
              },
            ],
          }}
        />
      ),
    },
    {
      name: "VirtualizedContainer",
      group: "Layout",
      description: "Lista ou grid com rolagem contínua; monta apenas os itens próximos à tela.",
      defaults: { mode: "Virtualizado", count: "10.000", columns: "1", overscan: "2" },
      controls: [
        select("mode", "Modo de renderização", ["Virtualizado", "Normal"]),
        select("count", "Quantidade de itens", ["1.000", "10.000", "50.000"]),
        select("columns", "Colunas", ["1", "2", "3"]),
        select("overscan", "Linhas antecipadas", ["1", "2", "4"]),
      ],
      render: (v) => (
        <iframe
          title="Comparação de lista normal e virtualizada"
          src={`#/preview/virtualized?mode=${v.mode}&count=${Number(String(v.count).replaceAll(".", ""))}&columns=${v.columns}&overscan=${v.overscan}&dark=${options.isDark()}`}
          class={options.isDetail() ? "h-120 w-full rounded-md border border-border" : "h-40 w-80 border border-border"}
        />
      ),
    },
    {
      name: "DashboardStatCard",
      group: "Dashboard",
      description:
        "Card de métrica com ícone, comparação com o período anterior e histórico compacto.",
      previewSize: { width: 340, height: 270 },
      defaults: {
        label: "Receita aprovada",
        value: "24850",
        previousValue: "19800",
        hint: "Receita aprovada no período",
        comparison: true,
        sparkline: true,
        loading: false,
      },
      controls: [
        text("label", "Título"),
        number("value", "Valor", { min: 0, max: 1000000 }),
        number("previousValue", "Valor anterior", { min: 0, max: 1000000 }),
        text("hint", "Descrição"),
        text("comparison", "Mostrar comparação"),
        text("sparkline", "Mostrar histórico"),
        text("loading", "Carregando"),
      ],
      render: (values, update) => (
        <StatCardExample
          fallback={<Loading size="sm" tone="muted" />}
          values={values}
          update={update}
        />
      ),
    },
    {
      name: "MetricDelta",
      group: "Dashboard",
      description:
        "Variação percentual com tendência e indicação de melhora, inclusive para métricas em que diminuir é melhor.",
      previewSize: { width: 400, height: 100 },
      defaults: {
        value: "125",
        previousValue: "100",
        improvement: "increase",
        label: "vs. período anterior",
      },
      controls: [
        number("value", "Valor atual", { min: -100000, max: 100000 }),
        number("previousValue", "Valor anterior", { min: -100000, max: 100000 }),
        select("improvement", "Melhora quando", ["increase", "decrease"]),
        text("label", "Comparação"),
      ],
      render: (values, update) => (
        <MetricDeltaExample
          fallback={<Loading size="sm" tone="muted" />}
          values={values}
          update={update}
        />
      ),
    },
    {
      name: "Sparkline",
      group: "Dashboard",
      description:
        "Histórico compacto em linha, área ou barras para complementar métricas, sem eixos nem controles.",
      previewSize: { width: 420, height: 120 },
      defaults: { type: "area", height: "80", empty: false, flat: false, showTooltip: true, fill: false },
      controls: [
        select("type", "Tipo", ["line", "area", "bar"]),
        number("height", "Altura", { min: 32, max: 180 }),
        text("empty", "Sem dados"),
        text("flat", "Valores constantes"),
        text("showTooltip", "Tooltip de valores"),
        text("fill", "Barras ocupam a largura"),
      ],
      render: (values, update) => (
        <SparklineExample
          fallback={<Loading size="sm" tone="muted" />}
          values={values}
          update={update}
        />
      ),
    },
    {
      name: "DashboardPanel",
      group: "Dashboard",
      description:
        "Seção responsiva de dashboard com título, ações, conteúdo e estados de carregamento ou erro.",
      previewSize: { width: 620, height: 480 },
      defaults: {
        title: "Visão geral",
        description: "Acompanhe os resultados do evento.",
        loading: false,
        error: false,
      },
      controls: [
        text("title", "Título"),
        text("description", "Descrição"),
        text("loading", "Carregando"),
        text("error", "Simular erro"),
      ],
      render: (values, update) => (
        <PanelExample
          fallback={<Loading size="sm" tone="muted" />}
          values={values}
          update={update}
        />
      ),
    },
    {
      name: "DashboardBarList",
      group: "Dashboard",
      description:
        "Ranking com barras proporcionais, valores e percentuais. Útil para canais de venda, status e produtos.",
      previewSize: { width: 420, height: 250 },
      defaults: { showPercentage: true, empty: false, loading: false },
      controls: [
        text("showPercentage", "Mostrar percentuais"),
        text("empty", "Sem dados"),
        text("loading", "Carregando"),
      ],
      render: (values, update) => (
        <BarListExample
          fallback={<Loading size="sm" tone="muted" />}
          values={values}
          update={update}
        />
      ),
    },
    {
      name: "GenericChart",
      group: "Dashboard",
      description:
        "Gráfico responsivo de linha, área, barras ou dispersão com séries, eixos, tooltip e formatação.",
      previewSize: { width: 560, height: 320 },
      defaults: {
        type: "line",
        height: "280",
        curveStyle: "smooth",
        multiSeries: true,
        showPoints: false,
        empty: false,
      },
      controls: [
        select("type", "Tipo", ["line", "area", "bar", "scatter"]),
        number("height", "Altura", { min: 120, max: 500 }),
        select("curveStyle", "Curva", ["smooth", "straight"]),
        text("multiSeries", "Múltiplas séries"),
        text("showPoints", "Mostrar pontos"),
        text("empty", "Sem dados"),
      ],
      render: (values, update) => (
        <ChartExample
          fallback={<Loading size="sm" tone="muted" />}
          bare
          values={values}
          update={update}
        />
      ),
    },
    {
      name: "ChartCard",
      group: "Dashboard",
      description:
        "Card de gráfico completo com filtros de período e séries, troca de tipo e estados de carregamento, erro ou valores ocultos.",
      previewSize: { width: 640, height: 500 },
      defaults: {
        title: "Resultado por período",
        type: "area",
        height: "280",
        multiSeries: true,
        showTypeSelector: true,
        showRangeFilter: true,
        showSeriesFilter: true,
        showSearchFilter: false,
        showPointsToggle: true,
        loading: false,
        error: false,
        empty: false,
        masked: false,
      },
      controls: [
        text("title", "Título"),
        select("type", "Tipo", ["line", "area", "bar", "scatter"]),
        number("height", "Altura", { min: 120, max: 500 }),
        text("multiSeries", "Múltiplas séries"),
        text("showTypeSelector", "Seletor de tipo"),
        text("showRangeFilter", "Filtro de período"),
        text("showSeriesFilter", "Filtro de séries"),
        text("showSearchFilter", "Buscar série"),
        text("showPointsToggle", "Alternar pontos"),
        text("loading", "Carregando"),
        text("error", "Simular erro"),
        text("empty", "Sem dados"),
        text("masked", "Ocultar valores"),
      ],
      render: (values, update) => (
        <ChartExample
          fallback={<Loading size="sm" tone="muted" />}
          values={values}
          update={update}
        />
      ),
    },
    {
      name: "Pagination",
      previewSize: { width: 480, height: 120 },
      group: "Navegação",
      description:
        "Paginação independente com números, reticências, tamanho de página ou navegação por cursor.",
      defaults: {
        page: "1",
        pageSize: "10",
        totalItems: "240",
        mode: "numbered",
        showPageSizeSelector: true,
        disabled: false,
      },
      controls: [
        number("page", "Página", { min: 1, max: 24 }),
        select("pageSize", "Itens por página", ["10", "25", "50", "100"]),
        number("totalItems", "Total de itens", { min: 0, max: 10000 }),
        select("mode", "Navegação", ["numbered", "cursor"]),
        text("showPageSizeSelector", "Seletor de itens por página"),
        text("disabled", "Desabilitado"),
      ],
      render: (values, update) => (
        <PaginationExample values={values} update={update} />
      ),
    },
    {
      name: "CollectionToolbar",
      previewSize: { width: 560, height: 520 },
      group: "Layout",
      description:
        "Busca, ordenação, contador e alternância de visualização para coleções.",
      defaults: {
        search: "",
        sort: "asc",
        layout: "list",
        status: "Todos",
        disabled: false,
        withPagination: false,
      },
      controls: [
        text("search", "Busca"),
        select("status", "Status", ["Todos", "Ativo", "Arquivado"]),
        select("sort", "Ordenação", ["asc", "desc"]),
        select("layout", "Visualização", ["list", "grid"]),
        text("disabled", "Desabilitado"),
        text("withPagination", "Exibir paginação"),
      ],
      render: (values, update) => (
        <ToolbarExample values={values} update={update} />
      ),
    },
    {
      name: "FilterBar",
      previewSize: { width: 480, height: 320 },
      group: "Layout",
      description:
        "Filtros ativos removíveis, com limpeza individual ou de todos os filtros.",
      defaults: { disabled: false },
      controls: [text("disabled", "Desabilitado")],
      render: (values, update) => (
        <FilterBarExample values={values} update={update} />
      ),
    },
    {
      name: "BulkActions",
      previewSize: { width: 480, height: 160 },
      group: "Ações",
      description:
        "Ações em lote sobre a seleção, com contador e limpeza de seleção.",
      defaults: { busy: false, iconOnly: false },
      controls: [
        text("busy", "Operação em andamento"),
        text("iconOnly", "Somente ícones"),
      ],
      render: (values, update) => (
        <BulkActionsExample values={values} update={update} />
      ),
    },
    {
      name: "InfiniteScroll",
      previewSize: { width: 480, height: 360 },
      group: "Layout",
      description:
        "Carrega mais itens ao aproximar-se do fim, com nova tentativa, carregamento manual e fim da lista.",
      defaults: { automatic: true, simulateError: false },
      controls: [
        text("automatic", "Carregamento automático"),
        text("simulateError", "Simular erro"),
      ],
      render: (values, update) => (
        <InfiniteScrollExample values={values} update={update} />
      ),
    },
    {
      name: "DataTable",
      previewSize: { width: 560, height: 540 },
      group: "Layout",
      description:
        "Tabela com seleção entre páginas, ordenação, ações por registro e paginação controlada.",
      defaults: {
        selectable: true,
        rowActions: true,
        showPageSizeSelector: true,
        loading: false,
        simulateError: false,
      },
      controls: [
        text("selectable", "Permitir seleção"),
        text("rowActions", "Ações por registro"),
        text("showPageSizeSelector", "Seletor de itens por página"),
        text("loading", "Carregando"),
        text("simulateError", "Simular erro"),
      ],
      render: (values, update) => (
        <DataTableExample values={values} update={update} />
      ),
    },
    {
      name: "PaginatedContainer",
      previewSize: { width: 560, height: 440 },
      group: "Layout",
      description: "Coleção com busca, ordenação e paginação.",
      defaults: {
        dataMode: "No componente",
        showPageSizeSelector: true,
        simulateError: false,
        filterValue: "",
        pageSize: "5",
        layout: "list",
        gap: "3",
        filterPlaceholder: "Buscar projetos…",
        itemLabel: "projetos",
        minItemWidth: "12rem",
      },
      controls: [
        select("dataMode", "Onde paginar (demonstração)", [
          "No componente",
          "No servidor (simulado)",
        ]),
        text("showPageSizeSelector", "Seletor de itens por página"),
        text("simulateError", "Simular erro do servidor"),
        select("pageSize", "Itens por página", ["5", "10", "25"]),
        select("layout", "Layout", ["list", "grid"]),
        select("gap", "Espaçamento", ["0", "2", "3", "4", "6"]),
        text("filterValue", "Busca"),
        text("filterPlaceholder", "Placeholder da busca"),
        text("itemLabel", "Nome dos itens"),
        text("minItemWidth", "Largura mínima"),
      ],
      render: (values, update) => (
        <PaginatedExample values={values} update={update} />
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
        <div class="w-full">
          <iframe
            title="Demonstração da OverlayScrollbar"
            src={`#/preview/scrollbar?enabled=${v.enabled}&dark=${options.isDark()}`}
            class={
              options.isDetail()
                ? "block h-80 w-full rounded-md border border-border bg-background"
                : "block h-37.5 w-full rounded-md border border-border bg-background"
            }
          />
          <p class="mt-3 text-[11px] text-muted-foreground">
            Role dentro da prévia para testar a barra.
          </p>
        </div>
      ),
    },
    {
      name: "DesktopOnly",
      group: "Layout",
      description: "Guarda responsiva que orienta o usuário a abrir no computador em resoluções menores.",
      defaults: {
        threshold: "1200",
        viewport: "1440",
        title: "Disponível apenas no computador",
        description: "Esta ferramenta foi projetada para telas a partir de 1200px de largura.",
        actionLabel: "Voltar à página anterior",
      },
      controls: [
        number("viewport", "Largura simulada (px)", { min: 320, max: 2560, step: 10 }),
        number("threshold", "Largura mínima (px)", { min: 320, max: 2560, step: 10 }),
        text("title", "Título"),
        text("description", "Descrição"),
        text("actionLabel", "Texto da ação"),
      ],
      render: (v, update) => {
        const viewport = () => Number(v.viewport) || 1440;
        const threshold = () => Number(v.threshold) || 1200;
        const resize = (width: number) =>
          update("viewport", String(Math.min(1920, Math.max(320, width))));
        return (
          <div data-wide class="relative flex h-112 w-full min-w-105 items-center justify-center overflow-hidden rounded-xl border border-border bg-muted/20 px-4 pt-8 pb-16">
            <div
              class="h-full min-w-60 overflow-hidden rounded-lg border border-border bg-background shadow-sm"
              style={{ width: `${Math.max((viewport() / 1920) * 100, 12)}%` }}
            >
              <DesktopOnly
                width={viewport()}
                threshold={threshold()}
                title={String(v.title)}
                description={String(v.description)}
                actionLabel={String(v.actionLabel)}
                class="min-h-0 h-full px-4 py-3"
                onAction={() => toast.info("Ação de voltar acionada!")}
              >
                <div class="flex h-full items-center justify-center p-6 text-center">
                  <span class="text-sm font-medium text-foreground">
                    Tela de {viewport()}px: atende ao mínimo de {threshold()}px. Diminua a largura para ver o aviso.
                  </span>
                </div>
              </DesktopOnly>
            </div>
            <span class="pointer-events-none absolute left-3 top-2 rounded bg-card/90 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
              {viewport()}px
            </span>
            <CanvasControls
              scale={viewport() / threshold()}
              onZoomIn={() => resize(viewport() + 120)}
              onZoomOut={() => resize(viewport() - 120)}
              onReset={() => resize(1440)}
            />
          </div>
        );
      },
    },
    {
      name: "CatalogNav",
      group: "Navegação",
      description: "Navegação lateral de catálogo ou documentação, com topbar para telas pequenas e bordas que se dissolvem.",
      defaults: { mode: "sidebar", fade: true, activeId: "CatalogNav" },
      controls: [
        select("mode", "Modo", ["sidebar", "topbar"]),
        text("fade", "Bordas suaves"),
      ],
      render: (v, update) => {
        const sections = catalogSections(entries);
        return (
          <div
            class={cn(
              "relative h-80 overflow-hidden rounded-xl border border-border bg-background",
              v.mode === "topbar" ? "w-80" : "w-64 p-2",
            )}
          >
            <CatalogNav
              mode={v.mode === "topbar" ? "topbar" : "sidebar"}
              fade={!!v.fade}
              sections={sections}
              top={catalogTop}
              activeId={String(v.activeId)}
              listClass="max-h-72"
              onNavigate={(item) => update("activeId", item.id)}
            />
          </div>
        );
      },
    },
  ];
  const shapeDefaults: Record<string, Shape> = {
    DropdownMenu: "round",
    Popover: "round",
    DatePicker: "soft",
    DateRangePicker: "soft",
    FileUpload: "round",
    DashboardStatCard: "round",
    ChartCard: "round",
    Checkbox: "soft",
    SearchInput: "soft",
    Button: "soft",
    Card: "round",
    Input: "soft",
    Textarea: "soft",
    Combobox: "soft",
    CodeBlock: "soft",
    Tooltip: "soft",
    MultiStepSummaryCard: "round",
    MultiStepStepper: "pill",
    EmptyState: "round",
    Skeleton: "soft",
    CardSkeleton: "soft",
    Dialog: "round",
    AlertModal: "round",
    Drawer: "round",
    MultiStepDialog: "round",
    Avatar: "soft",
    Badge: "subtle",
    Callout: "round",
    CopyButton: "soft",
    ToggleGroup: "pill",
  };
  for (const entry of entries) {
    const shape = shapeDefaults[entry.name];
    if (shape) {
      entry.defaults.shape = shape;
      entry.controls.push(select("shape", "Shape", shapes));
    }
  }
  for (const entry of entries) {
    entry.defaults.class = "";
  }
  return {
    entries,
    closeOverlays: () => {
      setOpened("");
    },
  };
}
