import { afterEach, beforeEach, expect, it, vi } from "vitest";
import { render } from "@solidjs/web";
import { createSignal } from "solid-js";
import { captureArtifact, expectNoDiagnostics } from "@solidjs/diagnostics";
import {
  MetricDelta,
  DashboardStatCard,
  DashboardPanel,
  DashboardBarList,
  GenericChart,
  ChartCard,
  Sparkline,
  type ChartDatum,
} from "@trieoh/ui-solid";
import { fillContinuousSeries } from "../../../lib/ts/ui-solid/src/chart/data-continuity";
import { computeRangeWindow } from "../../../lib/ts/ui-solid/src/chart/range";
import { dashboardUsage } from "../src/components/dashboard-usage";
import { transformWithOxc } from "vite";

const adapters = vi.hoisted(() => ({
  calls: [] as {
    mount: ReturnType<typeof vi.fn>;
    update: ReturnType<typeof vi.fn>;
    destroy: ReturnType<typeof vi.fn>;
    options: Record<string, unknown>;
  }[],
}));
vi.mock(
  "../../../lib/ts/ui-solid/node_modules/@tanstack/charts/dist/adapter-renderer.js",
  () => ({
    createChartRendererAdapter: (options: Record<string, unknown>) => {
      const instance = {
        mount: vi.fn(),
        update: vi.fn(),
        destroy: vi.fn(),
        options,
      };
      adapters.calls.push(instance);
      return instance;
    },
  }),
);
let host: HTMLDivElement;
let dispose: (() => void) | undefined;
const settle = () => new Promise((resolve) => setTimeout(resolve, 20));
const rows: ChartDatum[] = [
  { date: new Date(2026, 10, 1), series: "receita", value: 100 },
  { date: new Date(2026, 10, 2), series: "receita", value: 125 },
  { date: new Date(2026, 10, 1), series: "despesas", value: 40 },
];
beforeEach(() => {
  host = document.createElement("div");
  document.body.append(host);
  adapters.calls.length = 0;
});
afterEach(() => {
  dispose?.();
  dispose = undefined;
  host.remove();
  vi.restoreAllMocks();
});

it("compares metrics reactively without infinite percentages at zero", async () => {
  const { artifact } = await captureArtifact(
    async () => {
      const [value, setValue] = createSignal(125),
        [previous, setPrevious] = createSignal(100);
      dispose = render(
        () => <MetricDelta value={value()} previousValue={previous()} />,
        host,
      );
      await settle();
      expect(host.textContent).toContain("+25%");
      setValue(80);
      await settle();
      expect(host.textContent).toContain("-20%");
      setPrevious(0);
      await settle();
      expect(host.textContent).toContain("Sem base de comparação");
      expect(host.textContent).not.toContain("Infinity");
      setValue(0);
      await settle();
      expect(host.textContent).toContain("0%");
    },
    { scenario: "metric-delta", attribution: false },
  );
  expectNoDiagnostics(artifact);
});
it("can mark a decrease as an improvement and handle negative baselines", async () => {
  dispose = render(
    () => (
      <>
        <MetricDelta value={80} previousValue={100} improvement="decrease" />
        <MetricDelta value={-80} previousValue={-100} />
      </>
    ),
    host,
  );
  await settle();
  expect(host.textContent).toContain("-20%");
  expect(host.textContent).toContain("+20%");
  expect(host.querySelectorAll(".text-emerald-700").length).toBe(2);
});
it("composes a metric card with a comparison and a sparkline", async () => {
  const { artifact } = await captureArtifact(
    async () => {
      dispose = render(
        () => (
          <DashboardStatCard
            label="Inscrições"
            value="125"
            hint="Confirmadas"
            comparison={{ value: 125, previousValue: 100 }}
            sparkline={{ data: [2, 4, 3, 7], label: "Evolução" }}
          />
        ),
        host,
      );
      await settle();
      expect(host.textContent).toContain("Confirmadas");
      expect(host.textContent).toContain("+25%");
      await vi.waitFor(() => expect(adapters.calls).toHaveLength(1));
    },
    { scenario: "metric-card-composition", attribution: false },
  );
  expectNoDiagnostics(artifact);
});
it("renders panel loading, error and retry without losing the title", async () => {
  const [loading, setLoading] = createSignal(true),
    [error, setError] = createSignal<string>();
  const retry = vi.fn();
  dispose = render(
    () => (
      <DashboardPanel
        title="Resultados"
        loading={loading()}
        error={error()}
        onRetry={retry}
      >
        <p>Conteúdo</p>
      </DashboardPanel>
    ),
    host,
  );
  await settle();
  expect(host.querySelector("[role=status]")).not.toBeNull();
  expect(host.textContent).toContain("Resultados");
  setLoading(false);
  setError("Falha");
  await settle();
  expect(host.querySelector("[role=alert]")?.textContent).toBe("Falha");
  host.querySelector("button")!.click();
  expect(retry).toHaveBeenCalledOnce();
  setError(undefined);
  await settle();
  expect(host.textContent).toContain("Conteúdo");
});
it("clamps ranking bars and handles zero totals and invalid values", async () => {
  dispose = render(
    () => (
      <DashboardBarList
        items={[
          { id: "a", label: "Acima", value: 200 },
          { id: "b", label: "Inválido", value: NaN },
        ]}
        maxValue={100}
        total={0}
        showPercentage
      />
    ),
    host,
  );
  await settle();
  const meters = host.querySelectorAll("[role=meter]");
  expect(meters[0].getAttribute("aria-valuenow")).toBe("100");
  expect((meters[0].firstElementChild as HTMLElement).style.width).toBe("100%");
  expect(host.textContent).not.toMatch(/NaN|Infinity/);
  expect(host.textContent).toContain("0%");
});
it("mounts an initially empty chart, updates its size and destroys it on removal", async () => {
  const { artifact } = await captureArtifact(
    async () => {
      const [data, setData] = createSignal<ChartDatum[]>([]),
        [height, setHeight] = createSignal(200);
      dispose = render(
        () => (
          <GenericChart ariaLabel="Receita" data={data()} height={height()} />
        ),
        host,
      );
      await settle();
      expect(adapters.calls).toHaveLength(0);
      setData(rows);
      await settle();
      await vi.waitFor(() => expect(adapters.calls).toHaveLength(1));
      const adapter = adapters.calls[0];
      expect(adapter.mount).toHaveBeenCalledOnce();
      setHeight(350);
      await settle();
      expect(adapter.update.mock.calls.at(-1)?.[0].height).toBe(350);
      setData([]);
      await settle();
      expect(adapter.destroy).toHaveBeenCalledOnce();
      expect(host.textContent).toContain("Nenhum dado");
      setData(rows);
      await settle();
      await vi.waitFor(() => expect(adapters.calls).toHaveLength(2));
      dispose?.();
      dispose = undefined;
      expect(adapters.calls[1].destroy).toHaveBeenCalledOnce();
    },
    { scenario: "chart-async-lifecycle", attribution: false },
  );
  expectNoDiagnostics(artifact);
});
it("supports constant and one-point sparklines without leaking numeric errors", async () => {
  dispose = render(
    () => (
      <>
        <Sparkline data={[10, 10, 10]} />
        <Sparkline data={[5, NaN]} />
      </>
    ),
    host,
  );
  await settle();
  await vi.waitFor(() => expect(adapters.calls).toHaveLength(2));
  expect(host.textContent).not.toMatch(/NaN|Infinity/);
});
it("filters chart series and restores an empty selection", async () => {
  const { artifact } = await captureArtifact(
    async () => {
      dispose = render(
        () => (
          <ChartCard
            title="Receita"
            data={rows}
            initialRange="all"
            continuity={false}
            showRangeFilter={false}
          />
        ),
        host,
      );
      await settle();
      const button = (label: string) =>
        host.querySelector(`[aria-label="${label}"]`) as HTMLButtonElement;
      button("Série receita").click();
      button("Série despesas").click();
      await settle();
      expect(host.textContent).toContain("Nenhum dado");
      button("Série receita").click();
      await settle();
      expect(host.textContent).not.toContain("Nenhum dado");
    },
    { scenario: "chart-series-selection", attribution: false },
  );
  expectNoDiagnostics(artifact);
});
it("masks charts without mounting raw monetary data", async () => {
  const [masked, setMasked] = createSignal(true);
  dispose = render(
    () => (
      <ChartCard
        title="Receita privada"
        data={rows}
        isMasked={masked()}
        onToggleMask={() => setMasked(false)}
      />
    ),
    host,
  );
  await settle();
  expect(adapters.calls).toHaveLength(0);
  expect(host.textContent).toContain("Valores ocultos");
  const reveal = [...host.querySelectorAll("button")].find((button) =>
    button.textContent?.includes("Mostrar valores"),
  )!;
  reveal.click();
  await settle();
  await vi.waitFor(() => expect(adapters.calls).toHaveLength(1));
});
it("treats custom end dates as inclusive and normalizes reversed ranges", () => {
  const start = new Date(2026, 10, 10),
    end = new Date(2026, 10, 12);
  const window = computeRangeWindow("custom", { from: end, to: start });
  expect(window.windowStart?.getDate()).toBe(10);
  expect(window.windowEnd?.getHours()).toBe(23);
  expect(window.windowEnd?.getDate()).toBe(12);
  const seven = computeRangeWindow(
    "7d",
    { from: null, to: null },
    new Date(2026, 10, 19, 12),
  );
  expect(seven.windowStart?.getDate()).toBe(13);
});
it("fills continuous daily series without mutating original data or expanding huge histories", () => {
  const input = [
    { date: new Date(2026, 10, 1), value: 10, series: "a" },
    { date: new Date(2026, 10, 3), value: 20, series: "a" },
  ];
  const filled = fillContinuousSeries({
    data: input,
    seriesKeys: ["a"],
    windowStart: new Date(2026, 10, 1),
    windowEnd: new Date(2026, 10, 3),
    mode: "continuous",
  });
  expect(filled).toHaveLength(3);
  expect(filled[1].synthetic).toBe(true);
  expect(filled[1].value).toBe(10);
  expect(input).toHaveLength(2);
  const huge = fillContinuousSeries({
    data: input,
    seriesKeys: ["a"],
    windowStart: new Date(1900, 0, 1),
    windowEnd: new Date(2026, 10, 3),
    mode: "continuous",
  });
  expect(huge).toHaveLength(2);
});
it.each([
  "DashboardStatCard",
  "DashboardPanel",
  "DashboardBarList",
  "MetricDelta",
  "Sparkline",
  "GenericChart",
  "ChartCard",
])("provides copyable %s code", async (name) => {
  const code = dashboardUsage(name, {
    value: "125",
    previousValue: "100",
    label: "Receita",
    hint: "Confirmadas",
    title: "Resultado",
    description: "Resumo",
    type: name === "Sparkline" ? "area" : "line",
    height: "200",
    improvement: "increase",
    comparison: true,
    sparkline: true,
    multiSeries: true,
    showRangeFilter: true,
    showSeriesFilter: true,
    showPointsToggle: true,
  })!;
  expect(code).toContain("export default function Example");
  const result = await transformWithOxc(code, `${name}.tsx`, {
    jsx: "preserve",
  });
  expect(result.code.length).toBeGreaterThan(0);
});

it("can hide the type selector without hiding other chart filters", async () => {
  const [show, setShow] = createSignal(false);
  dispose = render(
    () => <ChartCard title="Resultado" data={rows} showTypeSelector={show()} />,
    host,
  );
  await settle();
  expect(host.querySelector('[aria-label="Tipo do gráfico"]')).toBeNull();
  expect(
    host.querySelector('[aria-label="Período do gráfico"]'),
  ).not.toBeNull();
  setShow(true);
  await settle();
  expect(host.querySelector('[aria-label="Tipo do gráfico"]')).not.toBeNull();
});

it("copies the actual preview with its formatting, data and callbacks", () => {
  const code = dashboardUsage("ChartCard", {
    type: "area",
    height: "280",
    showTypeSelector: false,
  })!;
  expect(code).toContain("Math.sin(index) * 20");
  expect(code).toContain('subtitle="Receita e despesas ao longo do período."');
  expect(code).toContain("onToggleMask={() => update");
  expect(code).toContain(
    "showTypeSelector={values().showTypeSelector !== false}",
  );
});
