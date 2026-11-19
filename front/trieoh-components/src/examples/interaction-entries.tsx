import type { Entry } from "../catalog";
import {
  DropdownExample,
  PopoverExample,
  DateExample,
  DateRangeExample,
  FileUploadExample,
  CommandExample,
  TimelineExample,
  HotkeysExample,
  KbdExample,
} from "./InteractionExamples";
const field = (key: string, label: string) => ({ key, label });
export const interactionEntries: Entry[] = [
  {
    name: "Hotkeys",
    group: "Ações",
    description:
      "Atalhos reutilizáveis com Mod, Ctrl, Command e Shift, escopo opcional, proteção de campos e limpeza automática. Registros concorrentes executam uma única ação.",
    previewSize: { width: 420, height: 300 },
    defaults: { shortcut: "Mod+Shift+L", enabled: true, ignoreInputs: true },
    controls: [
      field("shortcut", "Atalho"),
      field("enabled", "Habilitado"),
      field("ignoreInputs", "Ignorar campos de texto"),
    ],
    render: (values, update) => (
      <HotkeysExample values={values} update={update} />
    ),
  },
  {
    name: "Kbd",
    group: "Feedback",
    description:
      "Indicação visual com teclas tipadas. Neste exemplo, a busca registra o mesmo atalho e pode ser aberta por teclado ou pelo botão.",
    previewSize: { width: 360, height: 180 },
    defaults: { keys: "Mod+K" },
    controls: [field("keys", "Teclas")],
    render: (values, update) => <KbdExample values={values} update={update} />,
  },
  {
    name: "DropdownMenu",
    group: "Ações",
    description:
      "Menu de ações com ícones, atalhos funcionais, separadores e opções destrutivas. Navegação por teclado e busca por digitação.",
    previewSize: { width: 360, height: 180 },
    defaults: {
      align: "end",
      disabled: false,
      shortcuts: true,
      shortcutScope: "menu",
      disableDuplicate: false,
    },
    controls: [
      { key: "align", label: "Alinhamento", options: ["start", "end"] },
      field("shortcuts", "Exibir atalhos"),
      {
        key: "shortcutScope",
        label: "Escopo dos atalhos",
        options: ["menu", "global", "none"],
      },
      field("disableDuplicate", "Desabilitar duplicar"),
      field("disabled", "Desabilitado"),
    ],
    render: (values, update) => (
      <DropdownExample values={values} update={update} />
    ),
  },
  {
    name: "Popover",
    group: "Overlays",
    description:
      "Superfície ancorada para filtros e pequenos formulários, com posicionamento adaptável, fechamento externo e retorno de foco.",
    previewSize: { width: 360, height: 180 },
    defaults: { align: "end", disabled: false },
    controls: [
      { key: "align", label: "Alinhamento", options: ["start", "end"] },
      field("disabled", "Desabilitado"),
    ],
    render: (values, update) => (
      <PopoverExample values={values} update={update} />
    ),
  },
  {
    name: "DatePicker",
    group: "Formulários",
    description:
      "Seleção de data em calendário com limites, limpeza e navegação por setas, Home, End e Page Up/Down.",
    previewSize: { width: 340, height: 170 },
    defaults: {
      value: "",
      placeholder: "Selecionar data",
      clearable: true,
      disabled: false,
      futureOnly: false,
    },
    controls: [
      field("value", "Data (AAAA-MM-DD)"),
      field("placeholder", "Placeholder"),
      field("clearable", "Permitir limpar"),
      field("futureOnly", "Somente datas futuras"),
      field("disabled", "Desabilitado"),
    ],
    render: (values, update) => <DateExample values={values} update={update} />,
  },
  {
    name: "DateRangePicker",
    group: "Formulários",
    description:
      "Seleção de período com destaque do intervalo, datas limites e aplicação ou cancelamento sem alterar o valor anterior.",
    previewSize: { width: 380, height: 170 },
    defaults: {
      from: "",
      to: "",
      placeholder: "Selecionar período",
      clearable: true,
      disabled: false,
      futureOnly: false,
    },
    controls: [
      field("from", "Início (AAAA-MM-DD)"),
      field("to", "Fim (AAAA-MM-DD)"),
      field("placeholder", "Placeholder"),
      field("clearable", "Permitir limpar"),
      field("futureOnly", "Somente datas futuras"),
      field("disabled", "Desabilitado"),
    ],
    render: (values, update) => (
      <DateRangeExample values={values} update={update} />
    ),
  },
  {
    name: "FileUpload",
    group: "Formulários",
    description:
      "Seleção e arraste de arquivos, validação, prévia de imagens, progresso real, cancelamento e nova tentativa. O envio é fornecido pela aplicação.",
    previewSize: { width: 440, height: 380 },
    defaults: {
      label: "Envie os arquivos do evento",
      accept: "image/*,.pdf",
      maxSize: "5",
      maxFiles: "5",
      multiple: true,
      disabled: false,
      showPreview: true,
      simulateUpload: true,
      simulateError: false,
    },
    controls: [
      field("label", "Título"),
      field("accept", "Formatos aceitos"),
      {
        key: "maxSize",
        label: "Limite por arquivo (MB)",
        number: { min: 1, max: 100 },
      },
      {
        key: "maxFiles",
        label: "Máximo de arquivos",
        number: { min: 1, max: 20 },
      },
      field("multiple", "Múltiplos arquivos"),
      field("showPreview", "Prévia de imagem"),
      field("simulateUpload", "Simular envio"),
      field("simulateError", "Simular falha no envio"),
      field("disabled", "Desabilitado"),
    ],
    render: (values, update) => (
      <FileUploadExample values={values} update={update} />
    ),
  },
  {
    name: "CommandPalette",
    group: "Navegação",
    description:
      "Busca de páginas e ações com Ctrl/⌘+K, grupos, palavras-chave e navegação por teclado. A aplicação define as ações.",
    previewSize: { width: 400, height: 180 },
    defaults: {
      placeholder: "Buscar páginas e ações…",
      hotkey: true,
      shortcut: "Mod+K",
      resultsHeight: "320",
      empty: false,
    },
    controls: [
      field("placeholder", "Placeholder"),
      field("hotkey", "Habilitar atalho"),
      field("shortcut", "Atalho"),
      {
        key: "resultsHeight",
        label: "Altura dos resultados",
        number: { min: 128, max: 600 },
      },
      field("empty", "Sem comandos"),
    ],
    render: (values, update) => (
      <CommandExample values={values} update={update} />
    ),
  },
  {
    name: "Timeline",
    group: "Layout",
    description:
      "Histórico de atividades com ícones, datas, descrições, estados semânticos e ações opcionais, mantendo a ordem dos itens fornecidos.",
    previewSize: { width: 720, height: 300 },
    defaults: {
      orientation: "vertical",
      descriptions: true,
      loading: false,
      empty: false,
    },
    controls: [
      {
        key: "orientation",
        label: "Orientação",
        options: ["vertical", "horizontal"],
      },
      field("descriptions", "Descrições"),
      field("loading", "Carregando"),
      field("empty", "Sem atividades"),
    ],
    render: (values, update) => (
      <TimelineExample values={values} update={update} />
    ),
  },
];
