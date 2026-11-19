import { afterEach, beforeEach, expect, it, vi } from "vitest";
import { createSignal } from "solid-js";
import { render } from "@solidjs/web";
import { captureArtifact, expectNoDiagnostics } from "@solidjs/diagnostics";
import {
  Button,
  Popover,
  DropdownMenu,
  DatePicker,
  DateRangePicker,
  FileUpload,
  CommandPalette,
  Timeline,
  type DateRange,
} from "@trieoh/ui-solid";
import { interactionUsage } from "../src/components/interaction-usage";
import { interactionEntries } from "../src/examples/interaction-entries";
import { transformWithOxc } from "vite";
let host: HTMLDivElement, dispose: (() => void) | undefined;
const settle = () => new Promise((resolve) => setTimeout(resolve, 25));
const key = (target: Element, key: string, extra: KeyboardEventInit = {}) =>
  target.dispatchEvent(
    new KeyboardEvent("keydown", { key, bubbles: true, ...extra }),
  );
const button = (text: string) =>
  [...document.querySelectorAll<HTMLButtonElement>("button")].find(
    (el) => el.textContent?.trim() === text,
  )!;
const day = (value: string) =>
  document.querySelector<HTMLButtonElement>(`[data-day="${value}"]`)!;
beforeEach(() => {
  host = document.createElement("div");
  document.body.append(host);
});
afterEach(() => {
  dispose?.();
  dispose = undefined;
  host.remove();
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
});
it("opens menu, skips disabled actions, selects by keyboard and restores focus", async () => {
  const select = vi.fn();
  const { artifact } = await captureArtifact(async () => {
    dispose = render(
      () => (
        <DropdownMenu
          ariaLabel="Ações"
          trigger={(trigger) => <Button {...trigger}>Abrir</Button>}
          items={[
            { id: "a", label: "Editar", onSelect: select },
            { id: "b", label: "Bloqueado", disabled: true, onSelect: select },
            { id: "s", separator: true },
            { id: "c", label: "Excluir", destructive: true, onSelect: select },
          ]}
        />
      ),
      host,
    );
    key(button("Abrir"), "ArrowDown");
    await settle();
    expect(document.activeElement?.textContent).toBe("Editar");
    key(document.activeElement!, "ArrowDown");
    expect(document.activeElement?.textContent).toBe("Excluir");
    (document.activeElement as HTMLButtonElement).click();
    await settle();
    expect(select).toHaveBeenCalledTimes(1);
    expect(document.querySelector('[role="menu"]')).toBeNull();
    expect(document.activeElement).toBe(button("Abrir"));
    button("Abrir").click();
    await settle();
    key(document.activeElement!, "Escape");
    await settle();
    expect(document.querySelector('[role="menu"]')).toBeNull();
  });
  expectNoDiagnostics(artifact);
});
it("supports controlled popovers, outside dismissal, disabled state and cleanup", async () => {
  const [open, setOpen] = createSignal(false),
    [disabled, setDisabled] = createSignal(false);
  dispose = render(
    () => (
      <Popover
        open={open()}
        onOpenChange={setOpen}
        disabled={disabled()}
        ariaLabel="Filtros"
        trigger={(trigger) => <Button {...trigger}>Filtros</Button>}
      >
        <input aria-label="Nome" />
      </Popover>
    ),
    host,
  );
  button("Filtros").click();
  await settle();
  expect(open()).toBe(true);
  expect(document.activeElement?.getAttribute("aria-label")).toBe("Nome");
  document.body.dispatchEvent(new Event("pointerdown", { bubbles: true }));
  await settle();
  expect(open()).toBe(false);
  setDisabled(true);
  button("Filtros").click();
  await settle();
  expect(
    document.querySelector('[aria-label="Filtros"][role="dialog"]'),
  ).toBeNull();
  dispose();
  dispose = undefined;
  expect(document.querySelector('[aria-label="Nome"]')).toBeNull();
});
it("selects dates via the keyboard, enforces bounds and clears", async () => {
  const [value, setValue] = createSignal<Date | null>(new Date(2026, 10, 19));
  dispose = render(
    () => (
      <DatePicker
        value={value()}
        min={new Date(2026, 10, 19)}
        max={new Date(2026, 10, 25)}
        onChange={setValue}
      />
    ),
    host,
  );
  button("19/11/2026").click();
  await settle();
  expect(day("2026-11-18").getAttribute("aria-disabled")).toBe("true");
  day("2026-11-18").click();
  expect(value()?.getDate()).toBe(19);
  day("2026-11-19").focus();
  key(day("2026-11-19"), "ArrowRight");
  await settle();
  expect(document.activeElement).toBe(day("2026-11-20"));
  day("2026-11-20").click();
  await settle();
  expect(value()?.getDate()).toBe(20);
  button("20/11/2026").click();
  await settle();
  button("Limpar data").click();
  await settle();
  expect(value()).toBeNull();
});
it("keeps range edits in a draft until apply, normalizes reverse selection and cancels", async () => {
  const initial = { from: new Date(2026, 10, 19), to: new Date(2026, 10, 20) };
  const [value, setValue] = createSignal<DateRange>(initial);
  const { artifact } = await captureArtifact(async () => {
    dispose = render(
      () => <DateRangePicker value={value()} onChange={setValue} />,
      host,
    );
    button("19/11/2026 – 20/11/2026").click();
    await settle();
    day("2026-11-25").click();
    await settle();
    expect(value()).toBe(initial);
    expect(button("Aplicar").disabled).toBe(true);
    day("2026-11-22").click();
    await settle();
    button("Cancelar").click();
    await settle();
    expect(value()).toBe(initial);
    button("19/11/2026 – 20/11/2026").click();
    await settle();
    day("2026-11-25").click();
    await settle();
    day("2026-11-22").click();
    await settle();
    expect(day("2026-11-22").parentElement?.dataset.range).toBe("start");
    expect(day("2026-11-23").parentElement?.dataset.range).toBe("middle");
    expect(day("2026-11-23").parentElement?.getAttribute("aria-selected")).toBe(
      "true",
    );
    expect(day("2026-11-25").parentElement?.dataset.range).toBe("end");
    button("Aplicar").click();
    await settle();
    expect(value().from?.getDate()).toBe(22);
    expect(value().to?.getDate()).toBe(25);
  });
  expectNoDiagnostics(artifact);
});
const choose = async (files: File[]) => {
  const input = host.querySelector('input[type="file"]')!;
  Object.defineProperty(input, "files", { value: files, configurable: true });
  input.dispatchEvent(new Event("change", { bubbles: true }));
  await settle();
};
it("validates file formats, size and counts without uploading rejected files", async () => {
  const upload = vi.fn(async () => {}),
    reject = vi.fn();
  dispose = render(
    () => (
      <FileUpload
        accept="image/*,.pdf"
        maxSize={10}
        maxFiles={1}
        onUpload={upload}
        onReject={reject}
      />
    ),
    host,
  );
  await choose([
    new File(["text"], "bad.txt", { type: "text/plain" }),
    new File(["01234567890"], "large.pdf", { type: "application/pdf" }),
    new File(["ok"], "ticket.pdf", { type: "application/pdf" }),
    new File(["ok"], "extra.pdf", { type: "application/pdf" }),
  ]);
  expect(upload).toHaveBeenCalledTimes(1);
  expect(reject.mock.calls[0][0]).toHaveLength(3);
  expect(host.textContent).toContain("Formato não permitido");
  expect(host.textContent).toContain("ticket.pdf");
});
it("shows real upload progress, retries failures and cancels removed uploads", async () => {
  let signal: AbortSignal | undefined,
    attempt = 0,
    progress: ((value: number) => void) | undefined;
  dispose = render(
    () => (
      <FileUpload
        onUpload={async (_file, context) => {
          signal = context.signal;
          progress = context.onProgress;
          if (++attempt === 1) throw new Error("Rede indisponível");
          await new Promise<void>(() => {});
        }}
      />
    ),
    host,
  );
  await choose([new File(["ok"], "ticket.pdf")]);
  expect(host.textContent).toContain("Rede indisponível");
  host
    .querySelector<HTMLButtonElement>(
      '[aria-label="Tentar novamente ticket.pdf"]',
    )!
    .click();
  await settle();
  progress!(45);
  await settle();
  expect(
    host.querySelector('[role="progressbar"]')?.getAttribute("aria-valuenow"),
  ).toBe("45");
  host
    .querySelector<HTMLButtonElement>(
      '[aria-label="Cancelar envio de ticket.pdf"]',
    )!
    .click();
  await settle();
  expect(signal?.aborted).toBe(true);
  expect(host.querySelector('[role="progressbar"]')).toBeNull();
});
it("revokes image previews and aborts pending uploads on unmount", async () => {
  const create = vi.fn(() => "blob:preview"),
    revoke = vi.fn();
  vi.stubGlobal(
    "URL",
    class extends URL {
      static createObjectURL = create;
      static revokeObjectURL = revoke;
    },
  );
  let signal: AbortSignal | undefined;
  dispose = render(
    () => (
      <FileUpload
        onUpload={async (_file, context) => {
          signal = context.signal;
          await new Promise<void>(() => {});
        }}
      />
    ),
    host,
  );
  await choose([new File(["png"], "cover.png", { type: "image/png" })]);
  expect(host.querySelector("img")?.getAttribute("src")).toBe("blob:preview");
  dispose();
  dispose = undefined;
  expect(signal?.aborted).toBe(true);
  expect(revoke).toHaveBeenCalledWith("blob:preview");
  vi.unstubAllGlobals();
});
it("opens command palette with Ctrl K, searches without accents and executes enabled actions", async () => {
  const select = vi.fn();
  dispose = render(
    () => (
      <CommandPalette
        items={[
          {
            id: "off",
            label: "Ação indisponível",
            disabled: true,
            onSelect: select,
          },
          {
            id: "events",
            label: "Edições",
            keywords: ["evento"],
            onSelect: select,
          },
        ]}
      />
    ),
    host,
  );
  key(document.body, "k", { ctrlKey: true });
  await settle();
  const input = document.querySelector<HTMLInputElement>('[role="combobox"]')!;
  expect(input).not.toBeNull();
  input.value = "edicoes";
  input.dispatchEvent(new Event("input", { bubbles: true }));
  await settle();
  expect(document.querySelectorAll('[role="option"]')).toHaveLength(1);
  key(input, "Enter");
  await settle();
  expect(select).toHaveBeenCalledTimes(1);
});
it("renders ordered timeline semantics and loading and empty states", async () => {
  const [loading, setLoading] = createSignal(true);
  dispose = render(
    () => (
      <Timeline
        loading={loading()}
        items={[
          { id: "first", title: "Criada", date: new Date(2026, 10, 19) },
          {
            id: "second",
            title: "Aprovada",
            date: "2026-11-20T12:00:00Z",
            tone: "success",
          },
        ]}
      />
    ),
    host,
  );
  expect(host.querySelector('[role="status"]')).not.toBeNull();
  setLoading(false);
  await settle();
  expect([...host.querySelectorAll("li")].map((el) => el.textContent)).toEqual([
    expect.stringContaining("Criada"),
    expect.stringContaining("Aprovada"),
  ]);
  expect(host.querySelectorAll("time[datetime]")).toHaveLength(2);
});
it.each(interactionEntries)(
  "provides matching copyable $name examples",
  async (entry) => {
    const code = interactionUsage(entry.name, {
      ...entry.defaults,
      shape: "soft",
    })!;
    expect(code).toContain("export default function Example");
    expect(code).not.toContain("props.");
    const result = await transformWithOxc(code, `${entry.name}.tsx`, {
      jsx: "preserve",
    });
    expect(result.code).toContain(entry.name);
  },
);
