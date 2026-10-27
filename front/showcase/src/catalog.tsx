import { createSignal, For, Show } from "solid-js";
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
  SectionTabs,
  Combobox,
  EmptyState,
  Skeleton,
  CardSkeleton,
  toast,
  MultiStepDialog,
  MultiStepStepper,
  MultiStepFields,
  MultiStepSummaryCard,
  PaginatedContainer,
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
} from "@trieoh/ui-solid";
import IconArrowUpRight from "~icons/lucide/arrow-up-right";
import IconHome from "~icons/lucide/house";
import IconCalendar from "~icons/lucide/calendar-days";
import IconUsers from "~icons/lucide/users";
import IconPlus from "~icons/lucide/plus";
import IconMenu from "~icons/lucide/menu";
import { DeviceFrame } from "./components/DeviceFrame";

export type Values = Record<string, string | boolean>;
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
const select = (
  key: string,
  label: string,
  options: readonly string[],
): Control => ({
  key,
  label,
  options,
});

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
        text("duration", "Duração (ms)"),
        text("delay", "Intervalo (ms)"),
        text("lift", "Elevação (px)"),
        text("waveWidth", "Largura da onda (%)"),
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
      defaults: { name: "Sophia Hoffmann", initials: "", src: "", size: "md" },
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
        ]),
      ],
      render: (v) => (
        <Badge
          class={String(v.class)}
          variant={v.variant as "default" | "accent" | "secondary" | "outline"}
          shape={v.shape as Shape}
        >
          {String(v.children)}
        </Badge>
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
      defaults: { children: "Conteúdo do card" },
      controls: [text("children", "Conteúdo")],
      render: (v) => (
        <Card
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
        text("scale", "Escala"),
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
      },
      controls: [
        text("placeholder", "Placeholder"),
        text("disabled", "Desabilitado"),
        text("required", "Obrigatório"),
        text("readonly", "Somente leitura"),
        select("rows", "Linhas", ["2", "3", "4", "6"]),
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
      name: "SectionTabs",
      group: "Navegação",
      description: "Abas de seção com contagem e seleção controlada.",
      defaults: {
        active: "overview",
        ariaLabel: "Navegação de exemplo",
        count: "8",
        badge: "",
      },
      controls: [
        select("active", "Aba ativa", ["overview", "activity", "settings"]),
        text("ariaLabel", "Nome acessível"),
        select("count", "Contagem", ["0", "8", "24"]),
        text("badge", "Badge de ajustes"),
      ],
      render: (v, update) => (
        <SectionTabs
          class={String(v.class)}
          ariaLabel={String(v.ariaLabel)}
          items={[
            { id: "overview", label: "Visão geral" },
            { id: "activity", label: "Atividade", count: Number(v.count) },
            {
              id: "settings",
              label: "Ajustes",
              badge: String(v.badge) || undefined,
            },
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
          description: "Ajuste os dados antes de continuar.",
          closable: true,
          preventClose: false,
          loading: false,
          confirmLabel: "Confirmar",
          cancelLabel: "Cancelar",
          fixedHeight: true,
          nextLabel: "Continuar",
          submitLabel: "Concluir",
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
                  id={String(v.id) || undefined}
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
      },
      controls: [
        text("title", "Título"),
        text("description", "Descrição"),
        text("eyebrow", "Texto superior"),
        text("showAction", "Mostrar ação"),
        text("actionLabel", "Texto da ação"),
      ],
      render: (v) => (
        <EmptyState
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
      name: "PaginatedContainer",
      group: "Layout",
      description: "Coleção com busca, ordenação e paginação.",
      defaults: {
        filterValue: "",
        pageSize: "2",
        layout: "list",
        gap: "3",
        filterPlaceholder: "Buscar projetos…",
        itemLabel: "projetos",
        minItemWidth: "12rem",
      },
      controls: [
        select("pageSize", "Itens por página", ["1", "2", "3"]),
        select("layout", "Layout", ["list", "grid"]),
        select("gap", "Espaçamento", ["0", "2", "3", "4", "6"]),
        text("filterValue", "Busca"),
        text("filterPlaceholder", "Placeholder da busca"),
        text("itemLabel", "Nome dos itens"),
        text("minItemWidth", "Largura mínima"),
      ],
      render: (v, update) => (
        <PaginatedContainer
          class={String(v.class)}
          items={[
            { name: "Website" },
            { name: "Aplicativo" },
            { name: "Design system" },
          ]}
          pageSize={Number(v.pageSize)}
          layout={v.layout as "list" | "grid"}
          gap={v.gap as "0" | "2" | "3" | "4" | "6"}
          filterValue={String(v.filterValue)}
          onFilterChange={(value) => update("filterValue", value)}
          emptyState={
            <EmptyState
              title="Nenhum projeto encontrado"
              description="Tente outro termo de busca."
            />
          }
          filterPlaceholder={String(v.filterPlaceholder)}
          itemLabel={String(v.itemLabel)}
          minItemWidth={String(v.minItemWidth)}
          sortFields={[{ key: "name", label: "Nome" }]}
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
            src={`#/preview/scrollbar?enabled=${v.enabled}&dark=${options.isDark()}`}
            class={
              options.isDetail()
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
        text("viewport", "Largura simulada (px)"),
        text("threshold", "Largura mínima (px)"),
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
      defaults: { mode: "sidebar", fade: true, activeId: "Button" },
      controls: [
        select("mode", "Modo", ["sidebar", "topbar"]),
        text("fade", "Bordas suaves"),
      ],
      render: (v, update) => {
        const sections = ["Ações", "Formulários", "Navegação", "Overlays"].map((title) => ({
          id: title,
          title,
          href: "#/components/CatalogNav",
          items: ["Primeiro", "Segundo", "Terceiro", "Quarto", "Quinto", "Sexto"].map((label) => ({
            id: `${title}-${label}`,
            label: `${title} · ${label}`,
            href: "#/components/CatalogNav",
          })),
        }));
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
              top={{ id: "all", label: "Todos os componentes", href: "#/components/CatalogNav" }}
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
