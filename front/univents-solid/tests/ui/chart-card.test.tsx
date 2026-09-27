import { render } from "@solidjs/testing-library";
import { describe, expect, it } from "vitest";
import { ChartCard } from "../../src/widgets/chart/chart-card";

describe("ChartCard component", () => {
  it("renders chart card with title and subtitle", () => {
    const data = [
      { date: new Date("2026-08-01"), value: 100, series: "Lucro", purchases: 2 },
      { date: new Date("2026-08-02"), value: 250, series: "Lucro", purchases: 5 },
    ];

    const { container } = render(() => (
      <ChartCard
        title="Lucro"
        subtitle="Crescimento acumulado no período selecionado."
        data={data}
        allowedTypes={["line"]}
        initialRange="30d"
        continuity
        seriesLabels={{ Lucro: "Lucro acumulado" }}
        tooltipDetails={(datum) => [
          {
            label: "Compras",
            value: String(datum.purchases ?? 0),
          },
        ]}
      />
    ));

    expect(container.textContent).toContain("Lucro");
    expect(container.textContent).toContain("Crescimento acumulado no período selecionado.");
    expect(container.textContent).toContain("Últimos 30 dias");
  });
});
