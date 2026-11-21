import { afterEach, beforeEach, expect, it, vi } from "vitest";
import { render } from "@solidjs/web";
import { createSignal } from "solid-js";
import {
  Badge,
  Callout,
  Card,
  CopyButton,
  Drawer,
  EmptyState,
  Meter,
  PageHeader,
  PageSection,
  Sparkline,
  StatusDot,
  Table,
  ToggleGroup,
} from "@trieoh/ui-solid";

let host: HTMLDivElement;
let dispose: (() => void) | undefined;
const settle = () => new Promise((resolve) => setTimeout(resolve, 0));
beforeEach(() => {
  host = document.createElement("div");
  document.body.append(host);
});
afterEach(() => {
  dispose?.();
  dispose = undefined;
  host.remove();
  document.body.innerHTML = "";
});
const mount = (view: () => import("@solidjs/web").JSX.Element) => {
  dispose = render(view, host);
};

it("keeps the existing Badge variants and adds status ones", () => {
  mount(() => (
    <>
      <Badge>Padrão</Badge>
      <Badge variant="success">Ok</Badge>
      <Badge variant="destructive">Falhou</Badge>
    </>
  ));
  const [plain, ok, bad] = host.querySelectorAll("span");
  expect(plain.className).toContain("bg-muted");
  expect(ok.className).toContain("text-success");
  expect(bad.className).toContain("text-destructive");
});

it("renders Card as an article unless it has a link", () => {
  mount(() => (
    <>
      <Card>Sem link</Card>
      <Card href="https://trieoh.com" external>
        Externo
      </Card>
      <Card
        href="/docs"
        renderLink={(link) => (
          <a data-router href={link.href} class={link.class}>
            {link.children}
          </a>
        )}
      >
        Interno
      </Card>
    </>
  ));
  expect(host.querySelector("article")?.textContent).toBe("Sem link");
  const [external, internal] = host.querySelectorAll("a");
  expect(external.target).toBe("_blank");
  expect(external.rel).toContain("noopener");
  expect(internal.hasAttribute("data-router")).toBe(true);
  expect(internal.getAttribute("href")).toBe("/docs");
});

it("draws Sparkline bars, with a zero as a baseline dot", () => {
  mount(() => <Sparkline type="bar" data={[0, 5, 10]} label="Por hora" />);
  const chart = host.querySelector('[role="img"]')!;
  expect(chart.getAttribute("aria-label")).toBe("Por hora");
  const bars = chart.querySelectorAll("span");
  expect(bars).toHaveLength(3);
  expect(bars[0].className).toContain("h-px");
  expect((bars[2] as HTMLElement).style.height).toBe("100%");
});

it("keeps Drawer a bottom sheet by default and slides from the side when asked", async () => {
  mount(() => (
    <>
      <Drawer open onClose={() => {}} title="Base">
        <p>a</p>
      </Drawer>
      <Drawer open side="left" onClose={() => {}} title="Lado">
        <p>b</p>
      </Drawer>
    </>
  ));
  await settle();
  const [sheet, side] = document.querySelectorAll('[role="dialog"]');
  expect(sheet.className).toContain("rounded-t-2xl");
  expect(sheet.querySelector(".rounded-full")).not.toBeNull();
  expect(side.className).toContain("left-0");
  expect(side.className).not.toContain("rounded-t-2xl");
});

it("renders the inline EmptyState as a compact dashed box", () => {
  mount(() => <EmptyState variant="inline" title="Nada ainda." />);
  const box = host.firstElementChild!;
  expect(box.className).toContain("border-dashed");
  expect(box.textContent).toBe("Nada ainda.");
});

it("copies text read at click time and confirms it", async () => {
  const writeText = vi.fn().mockResolvedValue(undefined);
  Object.assign(navigator, { clipboard: { writeText } });
  const onCopied = vi.fn();
  const [value, setValue] = createSignal("antes");
  mount(() => <CopyButton text={value} onCopied={onCopied} />);
  setValue("depois");
  await settle();
  const button = host.querySelector("button")!;
  expect(button.getAttribute("aria-label")).toBe("Copiar");
  button.click();
  await settle();
  expect(writeText).toHaveBeenCalledWith("depois");
  expect(onCopied).toHaveBeenCalledWith("depois");
  expect(button.getAttribute("aria-label")).toBe("Copiado");
});

it("picks one ToggleGroup item at a time", async () => {
  const [value, setValue] = createSignal("all");
  mount(() => (
    <ToggleGroup
      ariaLabel="Filtro"
      value={value()}
      onChange={setValue}
      items={[
        { value: "all", label: "Todos" },
        { value: "bugs", label: "Bugs" },
      ]}
    />
  ));
  const [all, bugs] = host.querySelectorAll("button");
  expect(all.getAttribute("aria-pressed")).toBe("true");
  bugs.click();
  await settle();
  expect(value()).toBe("bugs");
  expect(bugs.getAttribute("aria-pressed")).toBe("true");
  expect(all.getAttribute("aria-pressed")).toBe("false");
});

it("renders a Table, or its empty state", async () => {
  const [rows, setRows] = createSignal([
    { name: "IdentityX", p95: 900 },
    { name: "Payssage", p95: 100 },
  ]);
  mount(() => (
    <Table
      caption="Serviços"
      items={rows()}
      emptyState={<p>Vazio</p>}
      columns={[
        { id: "name", header: "Serviço", cell: (row) => row.name },
        {
          id: "p95",
          header: "p95",
          align: "right",
          cell: (row) => row.p95,
          class: (row) => row.p95 > 500 && "text-warning",
        },
      ]}
    />
  ));
  expect(host.querySelector("caption")?.textContent).toBe("Serviços");
  const cells = host.querySelectorAll("tbody td");
  expect(cells).toHaveLength(4);
  expect(cells[1].className).toContain("text-right");
  expect(cells[1].className).toContain("text-warning");
  expect(cells[3].className).not.toContain("text-warning");
  setRows([]);
  await settle();
  expect(host.querySelector("table")).toBeNull();
  expect(host.textContent).toBe("Vazio");
});

it("reports Meter value and places its markers", () => {
  mount(() => (
    <Meter
      value={30}
      max={60}
      label="Uso"
      markers={[{ value: 45, label: "Limite", tone: "destructive" }]}
    />
  ));
  const meter = host.querySelector('[role="meter"]')!;
  expect(meter.getAttribute("aria-label")).toBe("Uso");
  expect(meter.getAttribute("aria-valuenow")).toBe("30");
  expect((meter.firstElementChild as HTMLElement).style.width).toBe("50%");
  expect((meter.querySelector("[title=Limite]") as HTMLElement).style.left).toBe("75%");
});

it("labels StatusDot and renders Callout, PageHeader and PageSection", () => {
  mount(() => (
    <>
      <StatusDot tone="idle" label="Ocioso" />
      <Callout tone="warning" title="Cuidado">
        <p>Texto</p>
      </Callout>
      <PageHeader title="Entrega" description="Deploys" />
      <PageSection title="Deploys" aside="30 dias">
        <p>Corpo</p>
      </PageSection>
    </>
  ));
  expect(host.querySelector('[role="img"]')?.getAttribute("aria-label")).toBe("Ocioso");
  expect(host.querySelector("aside")?.textContent).toContain("Cuidado");
  expect(host.querySelector("h1")?.textContent).toBe("Entrega");
  expect(host.querySelector("h2")?.textContent).toBe("Deploys");
});
