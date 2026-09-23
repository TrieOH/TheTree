import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import type { RichParagraph, RichRun } from "./types";

export interface StaticTextProps {
  paragraphs: readonly RichParagraph[];
  values?: Record<string, string>;
  showVariables?: boolean;
}

export const FALLBACK_LABELS: Record<string, string> = {
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

function replaceVariables(
  text: string,
  values?: Record<string, string>,
  showVariables?: boolean,
): string {
  if (!values && !showVariables) return text;
  return text.replace(/\{\{([^}]+)\}\}/g, (match, rawKey: string) => {
    const key = rawKey.trim();
    const val = resolveValue(key, values);
    if (val !== undefined && val !== "") {
      return val;
    }
    if (showVariables) {
      return FALLBACK_LABELS[key] ?? `{{${key}}}`;
    }
    return val !== undefined ? val : match;
  });
}

function findNearestRun(
  paragraphs: readonly RichParagraph[],
  paragraph: RichParagraph,
): RichRun | undefined {
  const index = paragraphs.indexOf(paragraph);
  const before = paragraphs
    .slice(0, index)
    .reverse()
    .flatMap((p) => [...p.runs].reverse());
  const after = paragraphs.slice(index + 1).flatMap((p) => p.runs);
  return [...before, ...after][0];
}

export function StaticText(props: StaticTextProps): JSX.Element {
  const showVariables = () => props.showVariables ?? true;

  return (
    <div class="h-full w-full select-none overflow-hidden">
      <For each={props.paragraphs}>
        {(paragraph) => {
          const fallbackRun = () =>
            paragraph.runs[0] ?? findNearestRun(props.paragraphs, paragraph);

          return (
            <p
              style={{
                "text-align": paragraph.align ?? "left",
                "line-height": paragraph.lineHeight ? `${paragraph.lineHeight}` : "1.25",
                "font-size": `${fallbackRun()?.fontSize ?? 14}px`,
                "font-family": fallbackRun()?.fontFamily ?? "sans-serif",
                "font-weight": fallbackRun()?.bold ? "bold" : "normal",
                "font-style": fallbackRun()?.italic ? "italic" : "normal",
                "text-decoration": fallbackRun()?.underline ? "underline" : "none",
                color: fallbackRun()?.color ?? "#000000",
              }}
            >
              <Show
                when={paragraph.runs.length > 0}
                fallback={<span>&nbsp;</span>}
              >
                <For each={paragraph.runs}>
                  {(run) => {
                    const text = () =>
                      replaceVariables(run.text, props.values, showVariables());

                    return (
                      <span
                        style={{
                          "font-size": `${run.fontSize}px`,
                          "font-family": run.fontFamily,
                          "font-weight": run.bold ? "bold" : "normal",
                          "font-style": run.italic ? "italic" : "normal",
                          "text-decoration": run.underline ? "underline" : "none",
                          color: run.color,
                        }}
                      >
                        {text()}
                      </span>
                    );
                  }}
                </For>
              </Show>
            </p>
          );
        }}
      </For>
    </div>
  );
}
