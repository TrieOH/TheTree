import { createSignal, Show } from "solid-js";
import {
  Button,
  DropdownMenu,
  Popover,
  DatePicker,
  DateRangePicker,
  FileUpload,
  CommandPalette,
  Timeline,
  Kbd,
  createHotkeys,
  parseHotkey,
  Field,
  Input,
  Switch,
  type Shape,
  type DropdownMenuProps,
  type CommandItem,
} from "@trieoh/ui-solid";
import Ellipsis from "~icons/lucide/ellipsis";
import Pencil from "~icons/lucide/pencil";
import Copy from "~icons/lucide/copy";
import Trash from "~icons/lucide/trash-2";
import Sliders from "~icons/lucide/sliders-horizontal";
import Search from "~icons/lucide/search";
import Calendar from "~icons/lucide/calendar-days";
import CreditCard from "~icons/lucide/credit-card";
import FileDown from "~icons/lucide/file-down";
import Ticket from "~icons/lucide/ticket";
import type { Values } from "../catalog";
type Props = {
  values: Values;
  update: (key: string, value: string | boolean) => void;
};
// Fixtures
const localDate = (value: unknown) => {
  if (typeof value !== "string" || !value) return null;
  const [year, month, day] = value.split("-").map(Number);
  const date = new Date(year, month - 1, day);
  return Number.isFinite(date.getTime()) ? date : null;
};
const dateValue = (date: Date | null) =>
  date
    ? `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`
    : "";
const today = () => {
  const date = new Date();
  return new Date(date.getFullYear(), date.getMonth(), date.getDate());
};
const demoUpload = async (context: {
  signal: AbortSignal;
  onProgress: (value: number) => void;
}) => {
  // Simulates transport only in this demo; production supplies its own uploader.
  for (let progress = 0; progress <= 100; progress += 20) {
    await new Promise<void>((resolve, reject) => {
      const abort = () => {
        clearTimeout(timer);
        reject(new Error("Cancelado"));
      };
      const timer = setTimeout(() => {
        context.signal.removeEventListener("abort", abort);
        resolve();
      }, 180);
      context.signal.addEventListener("abort", abort, { once: true });
      if (context.signal.aborted) abort();
    });
    context.onProgress(progress);
  }
};
const demoCommands = (select: (label: string) => void): CommandItem[] => [
  {
    id: "events",
    label: "Meus eventos",
    description: "Gerenciar eventos e edições",
    group: "Páginas",
    icon: <Calendar />,
    keywords: ["eventos", "edições"],
    onSelect: () => select("Meus eventos"),
  },
  {
    id: "payments",
    label: "Pagamentos",
    description: "Consultar compras e receitas",
    group: "Páginas",
    icon: <CreditCard />,
    onSelect: () => select("Pagamentos"),
  },
  {
    id: "new",
    label: "Criar evento",
    description: "Começar um novo evento",
    group: "Ações",
    icon: <Calendar />,
    onSelect: () => select("Criar evento"),
  },
  {
    id: "export",
    label: "Exportar relatório",
    description: "Baixar os dados do evento",
    group: "Ações",
    icon: <FileDown />,
    disabled: true,
    onSelect: () => select("Exportar relatório"),
  },
];
export function DropdownExample(props: Props) {
  const [result, setResult] = createSignal("Nenhuma ação selecionada");
  return (
    <div class="flex w-full min-w-0 flex-col items-center gap-4 p-2">
      <DropdownMenu
        ariaLabel="Ações do projeto"
        shortcutScope={
          props.values.shortcutScope as DropdownMenuProps["shortcutScope"]
        }
        align={props.values.align as "start" | "end"}
        disabled={!!props.values.disabled}
        shape={props.values.shape as Shape}
        trigger={(trigger) => (
          <Button
            {...trigger}
            variant="outline"
            shape={props.values.shape as Shape}
          >
            <Ellipsis class="size-4" />
            Ações do projeto
          </Button>
        )}
        items={[
          {
            id: "edit",
            label: "Editar projeto",
            icon: <Pencil />,
            shortcut: props.values.shortcuts
              ? { key: "e", modifiers: ["Mod"] }
              : undefined,
            onSelect: () => setResult("Editar projeto"),
          },
          {
            id: "copy",
            label: "Duplicar projeto",
            icon: <Copy />,
            disabled: !!props.values.disableDuplicate,
            onSelect: () => setResult("Duplicar projeto"),
          },
          { id: "separator", separator: true },
          {
            id: "delete",
            label: "Excluir projeto",
            icon: <Trash />,
            destructive: true,
            onSelect: () => setResult("Excluir projeto"),
          },
        ]}
      />
      <p role="status" class="text-center text-xs text-muted-foreground">
        {result()}
      </p>
    </div>
  );
}
export function PopoverExample(props: Props) {
  const [name, setName] = createSignal("Inscrições"),
    [active, setActive] = createSignal(true),
    [result, setResult] = createSignal("");
  return (
    <div class="flex w-full min-w-0 flex-col items-center gap-4 p-2">
      <Popover
        ariaLabel="Filtros avançados"
        align={props.values.align as "start" | "end"}
        disabled={!!props.values.disabled}
        shape={props.values.shape as Shape}
        trigger={(trigger) => (
          <Button {...trigger} variant="outline">
            <Sliders class="size-4" />
            Filtros avançados
          </Button>
        )}
      >
        {(close) => (
          <form
            class="space-y-4"
            onSubmit={(e) => {
              e.preventDefault();
              setResult(`${name()} · ${active() ? "Somente ativos" : "Todos"}`);
              close();
            }}
          >
            <div>
              <h3 class="text-sm font-semibold">Filtros avançados</h3>
              <p class="mt-1 text-xs text-muted-foreground">
                Refine a lista de inscrições.
              </p>
            </div>
            <Field label="Nome">
              {(ids) => (
                <Input
                  id={ids.id}
                  value={name()}
                  onInput={(e) => setName(e.currentTarget.value)}
                />
              )}
            </Field>
            <label class="flex items-center justify-between gap-3 text-sm">
              Somente ativos
              <Switch
                aria-label="Somente ativos"
                checked={active()}
                onChange={setActive}
              />
            </label>
            <Button type="submit" class="w-full">
              Aplicar filtros
            </Button>
          </form>
        )}
      </Popover>
      <p role="status" class="text-xs text-muted-foreground">
        {result() || "Abra para configurar os filtros"}
      </p>
    </div>
  );
}
export function DateExample(props: Props) {
  return (
    <div class="flex w-full min-w-0 flex-col items-start gap-2">
      <span class="text-sm font-medium">Data do evento</span>
      <DatePicker
        value={localDate(props.values.value)}
        onChange={(date) => props.update("value", dateValue(date))}
        placeholder={String(props.values.placeholder)}
        disabled={!!props.values.disabled}
        clearable={!!props.values.clearable}
        min={props.values.futureOnly ? today() : undefined}
        shape={props.values.shape as Shape}
        class="max-w-full"
      />
      <p class="text-xs text-muted-foreground">
        Selecione uma data no calendário.
      </p>
    </div>
  );
}
export function DateRangeExample(props: Props) {
  return (
    <div class="flex w-full min-w-0 flex-col items-start gap-2">
      <span class="text-sm font-medium">Período do evento</span>
      <DateRangePicker
        value={{
          from: localDate(props.values.from),
          to: localDate(props.values.to),
        }}
        onChange={(range) => {
          props.update("from", dateValue(range.from));
          props.update("to", dateValue(range.to));
        }}
        placeholder={String(props.values.placeholder)}
        disabled={!!props.values.disabled}
        clearable={!!props.values.clearable}
        min={props.values.futureOnly ? today() : undefined}
        shape={props.values.shape as Shape}
        class="max-w-full"
      />
      <p class="text-xs text-muted-foreground">
        Escolha o início e o fim, depois aplique.
      </p>
    </div>
  );
}
export function FileUploadExample(props: Props) {
  const [count, setCount] = createSignal(0);
  return (
    <div class="w-full min-w-0 space-y-3">
      <FileUpload
        label={String(props.values.label)}
        accept={String(props.values.accept)}
        maxSize={Number(props.values.maxSize) * 1024 * 1024}
        maxFiles={Number(props.values.maxFiles)}
        multiple={!!props.values.multiple}
        disabled={!!props.values.disabled}
        showPreview={!!props.values.showPreview}
        shape={props.values.shape as Shape}
        onFilesChange={(files) => setCount(files.length)}
        onUpload={
          props.values.simulateUpload
            ? async (_file, context) => {
                await demoUpload(context);
                if (props.values.simulateError)
                  throw new Error(
                    "Falha simulada. Desative o erro e tente novamente.",
                  );
              }
            : undefined
        }
      />
      <p role="status" class="text-xs text-muted-foreground">
        {count()} arquivo(s) selecionado(s) · envio simulado apenas nesta
        demonstração.
      </p>
    </div>
  );
}
export function CommandExample(props: Props) {
  const [open, setOpen] = createSignal(false),
    [result, setResult] = createSignal("Nenhuma ação selecionada");
  return (
    <div class="flex w-full min-w-0 flex-col items-center gap-4">
      <Button
        variant="outline"
        class="max-w-full gap-3"
        onClick={() => setOpen(true)}
      >
        <Search class="size-4" />
        Buscar comandos
        <Show when={props.values.hotkey}>
          <Kbd keys={parseHotkey(String(props.values.shortcut))} />
        </Show>
      </Button>
      <CommandPalette
        open={open()}
        onOpenChange={setOpen}
        hotkey={
          !!props.values.hotkey && !!parseHotkey(String(props.values.shortcut))
        }
        shortcut={parseHotkey(String(props.values.shortcut)) ?? undefined}
        resultsHeight={Number(props.values.resultsHeight)}
        placeholder={String(props.values.placeholder)}
        items={props.values.empty ? [] : demoCommands(setResult)}
      />
      <p role="status" class="text-xs text-muted-foreground">
        {result()}
      </p>
    </div>
  );
}
export function TimelineExample(props: Props) {
  return (
    <div data-wide class="w-full min-w-0">
      <Timeline
        orientation={props.values.orientation as "vertical" | "horizontal"}
        loading={!!props.values.loading}
        items={
          props.values.empty
            ? []
            : [
                {
                  id: "created",
                  title: "Inscrição criada",
                  date: new Date(2026, 10, 19, 9, 30),
                  description: props.values.descriptions
                    ? "Inscrição recebida para o encontro de tecnologia."
                    : undefined,
                  icon: <Ticket />,
                  tone: "default",
                },
                {
                  id: "payment",
                  title: "Pagamento aprovado",
                  date: new Date(2026, 10, 19, 9, 32),
                  description: props.values.descriptions
                    ? "Pagamento de R$ 120,00 confirmado."
                    : undefined,
                  icon: <CreditCard />,
                  tone: "success",
                },
                {
                  id: "ticket",
                  title: "Ingresso disponível",
                  date: new Date(2026, 10, 19, 9, 33),
                  description: props.values.descriptions
                    ? "O participante já pode acessar o ingresso."
                    : undefined,
                  icon: <Calendar />,
                  tone: "default",
                },
              ]
        }
      />
    </div>
  );
}

export function HotkeysExample(props: Props) {
  const [count, setCount] = createSignal(0);
  // Editable catalog settings are an external string boundary; application code
  // can use a typed object directly, e.g. { key: "l", modifiers: ["Mod", "Shift"] }.
  const shortcut = () => parseHotkey(String(props.values.shortcut));
  createHotkeys(
    () => {
      const hotkey = shortcut();
      return hotkey
        ? [
            {
              hotkey,
              callback: () => setCount((value) => value + 1),
              options: { enabled: !!props.values.enabled },
            },
          ]
        : [];
    },
    {
      get ignoreInputs() {
        return !!props.values.ignoreInputs;
      },
    },
  );
  return (
    <div class="w-full min-w-0 space-y-4 rounded-xl border border-border p-4">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <p class="text-sm font-medium">Atalho de demonstração</p>
        <Kbd keys={shortcut()} />
      </div>
      <p class="text-sm text-muted-foreground">
        Pressione o atalho para incrementar o contador.
      </p>
      <Field label="Teste enquanto digita">
        {(ids) => (
          <Input id={ids.id} placeholder="O atalho ignora campos por padrão" />
        )}
      </Field>
      <div class="flex items-center justify-between gap-3">
        <p role="status" class="text-sm tabular-nums">
          Executado {count()} vez(es)
        </p>
        <Button size="sm" variant="outline" onClick={() => setCount(0)}>
          Zerar
        </Button>
      </div>
    </div>
  );
}
export function KbdExample(props: Props) {
  const [open, setOpen] = createSignal(false);
  const [result, setResult] = createSignal("");
  const shortcut = () => parseHotkey(String(props.values.keys));
  return (
    <div class="flex w-full min-w-0 flex-col items-center gap-3">
      <Button variant="outline" onClick={() => setOpen(true)}>
        <Search class="size-4" />
        Abrir busca <Kbd keys={shortcut()} />
      </Button>
      <p class="text-center text-xs text-muted-foreground">
        Pressione o atalho ou use o botão.
      </p>
      <p role="status" class="text-xs text-muted-foreground">
        {result()}
      </p>
      <CommandPalette
        open={open()}
        onOpenChange={setOpen}
        shortcut={shortcut() ?? undefined}
        hotkey={!!shortcut()}
        items={demoCommands(setResult)}
      />
    </div>
  );
}
