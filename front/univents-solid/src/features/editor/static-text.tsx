import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import type { RichParagraph, RichRun } from "./types";

export interface StaticTextProps {
  paragraphs: readonly RichParagraph[];
  values?: Record<string, string>;
  showVariables?: boolean;
}

const FALLBACK_LABELS: Record<string, string> = {
  event_name: "Nome do Evento",
  edition_name: "Nome da Edição",
  ticket_name: "Nome do Ingresso",
  ticket_type: "Nome do Ingresso",
  ticket: "Nome do Ingresso",
  participant_name: "Nome do Participante",
  name: "Nome do Participante",
  activity_name: "Nome da Atividade",
  program_name: "Nome da Programação",
  participation_type: "Participação",
  location: "Local do Evento",
  location_name: "Local do Evento",
  workload_hours: "10 horas",
  participation_date: "19/09/2026",
  certified_at: "19/09/2026",
  issue_date: "19/09/2026",
  cert_hash: "UNIV-2026-CERT-SAMPLE",
  verify_url: "https://univents.app/verify/sample",
  checkin_url: "https://univents.app/check-in",
};

function resolveValue(key: string, values?: Record<string, string>): string | undefined {
  if (!values) return undefined;
  if (values[key] !== undefined && values[key] !== "") return values[key];
  if (key === "ticket_name" || key === "ticket_type" || key === "ticket") {
    return values.ticket_name ?? values.ticket_type ?? values.ticket;
  }
  if (key === "participant_name" || key === "name") {
    return values.participant_name ?? values.name;
  }
  if (key === "activity_name" || key === "program_name") {
    return values.activity_name ?? values.program_name;
  }
  if (key === "location" || key === "location_name") {
    return values.location ?? values.location_name;
  }
  return undefined;
}

function renderRun(
  run: RichRun,
  getValues: () => Record<string, string> | undefined,
  getShowVariables: () => boolean,
): JSX.Element {
  const isVariable = run.text.startsWith("{{") && run.text.endsWith("}}");
  const text = () => {
    if (!isVariable) return run.text;
    const key = run.text.slice(2, -2).trim();
    const resolved = resolveValue(key, getValues());
    return resolved ?? (getShowVariables() ? (FALLBACK_LABELS[key] ?? run.text) : "");
  };

  const style: JSX.CSSProperties = {
    color: run.color,
    "font-size": `${run.fontSize}px`,
    "font-family": run.fontFamily,
    "font-weight": run.bold ? 700 : 400,
    "font-style": run.italic ? "italic" : "normal",
    "text-decoration": run.underline ? "underline" : "none",
  };

  return (
    <span
      class={isVariable ? "rounded bg-primary/10 px-0.5 text-primary" : undefined}
      style={style}
    >
      {text() || "\u200b"}
    </span>
  );
}

export function StaticText(props: StaticTextProps): JSX.Element {
  const values = () => props.values;
  const showVariables = () => props.showVariables ?? true;

  return (
    <div class="pointer-events-none w-full select-none">
      <For
        each={props.paragraphs.length > 0 ? props.paragraphs : [{ align: "left" as const, lineHeight: 1.25, runs: [] }]}
      >
        {(paragraph) => (
          <p
            style={{
              "text-align": paragraph.align,
              "line-height": paragraph.lineHeight,
              margin: 0,
              "min-height": "1.25em",
            }}
          >
            <Show
              when={paragraph.runs.length > 0}
              fallback={<span style={{ "font-size": "16px" }}>{"\u200b"}</span>}
            >
              <For each={paragraph.runs}>
                {(run) => renderRun(run, values, showVariables)}
              </For>
            </Show>
          </p>
        )}
      </For>
    </div>
  );
}
