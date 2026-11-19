import type { Values } from "../catalog";
import exampleSource from "../examples/DashboardExamples.tsx?raw";

const examples: Record<string, string> = {
  DashboardStatCard: "StatCardExample",
  DashboardPanel: "PanelExample",
  DashboardBarList: "BarListExample",
  MetricDelta: "MetricDeltaExample",
  Sparkline: "SparklineExample",
  GenericChart: "ChartExample",
  ChartCard: "ChartExample",
};

/** Generate copyable code from the actual preview, so its content and behavior stay aligned. */
export function dashboardUsage(
  name: string,
  values: Values,
): string | undefined {
  const example = examples[name];
  if (!example) return undefined;
  const start = exampleSource.indexOf(`export function ${example}(`);
  const next = exampleSource.indexOf("\nexport function ", start + 1);
  const body = exampleSource
    .slice(start, next < 0 ? undefined : next)
    .replace(
      /export function \w+\(props: Props(?: & \{ bare\?: boolean \})?\)/,
      "export default function Example()",
    )
    .replaceAll("props.values", "values()")
    .replaceAll("props.update", "update")
    .replaceAll("props.bare", String(name === "GenericChart"))
    .replace(
      "export default function Example() {",
      `export default function Example() {\n  const [values, setValues] = createSignal(${JSON.stringify(values, null, 2)});\n  const update = (key: string, value: string | boolean) =>\n    setValues(previous => ({ ...previous, [key]: value }));`,
    );
  // Keep the same fixtures, formatting, wrappers, callbacks and component props as the preview.
  const fixtureSource = exampleSource.slice(
    exampleSource.indexOf("const history ="),
    exampleSource.indexOf("export function StatCardExample"),
  );
  const fixtures = fixtureSource
    .split(/(?=^(?:const |export function ))/m)
    .filter((source) => {
      const identifier = source.match(/^(?:const |export function )(\w+)/)?.[1];
      return identifier && new RegExp(`\\b${identifier}\\b`).test(body);
    })
    .join("");
  const used = (identifier: string) =>
    new RegExp(`\\b${identifier}\\b`).test(body + fixtures);
  const imports = exampleSource
    .slice(0, exampleSource.indexOf("import type { Values }"))
    .replace("import { createMemo }", "import { createMemo, createSignal }")
    .replace(
      /import \{([^}]+)\} from ([^;]+);/g,
      (_, names: string, module: string) => {
        const selected = names
          .split(",")
          .map((item) => item.trim())
          .filter((item) => used(item.replace(/^type /, "")));
        return selected.length
          ? `import { ${selected.join(", ")} } from ${module};`
          : "";
      },
    )
    .replace(/import (\w+) from ([^;]+);/g, (statement, identifier: string) =>
      used(identifier) ? statement : "",
    );
  return `${imports.trim()}\n\n${fixtures}\n${body}`;
}
