import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import type { RichParagraph, RichRun } from "./types";

export interface StaticTextProps {
  paragraphs: readonly RichParagraph[];
  values?: Record<string, string>;
  showVariables?: boolean;
}

const FALLBACK_LABELS: Record<string, string> = {
  event_name: "Nome do evento",
  edition_name: "Nome da edição",
  ticket_name: "Nome do ingresso",
  ticket_type: "Nome do ingresso",
  ticket: "Nome do ingresso",
  participant_name: "Nome do participante",
  name: "Nome do participante",
  location: "Local da edição",
  checkin_url: "Link de check-in",
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
      return FALLBACK_LABELS[key] ?? key;
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
    .find((p) => p.runs.length > 0);
  if (before && before.runs.length > 0) {
    return before.runs[before.runs.length - 1];
  }
  const after = paragraphs
    .slice(index + 1)
    .find((p) => p.runs.length > 0);
  if (after && after.runs.length > 0) {
    return after.runs[0];
  }
  return undefined;
}

export function StaticText(props: StaticTextProps): JSX.Element {
  return (
    <div
      class="h-full w-full overflow-hidden select-none whitespace-pre-wrap wrap-break-word"
      style={{
        "line-height": 1.25,
        "overflow-wrap": "anywhere",
        "word-break": "break-word",
      }}
    >
      <For each={props.paragraphs}>
        {(paragraph) => {
          const isEmpty = () =>
            paragraph.runs.length === 0 ||
            paragraph.runs.every((r) => !r.text || r.text === "");

          const nearestRun = () =>
            paragraph.runs[0] ?? findNearestRun(props.paragraphs, paragraph);

          return (
            <div
              style={{
                "text-align": paragraph.align,
                "line-height": paragraph.lineHeight ?? 1.25,
                "font-size": `${nearestRun()?.fontSize ?? 16}px`,
                "font-family": nearestRun()?.fontFamily,
                margin: 0,
                padding: 0,
              }}
            >
              <Show
                when={!isEmpty()}
                fallback={<br />}
              >
                <For each={paragraph.runs}>
                  {(run) => (
                    <span
                      style={{
                        "font-size": `${run.fontSize}px`,
                        "font-family": run.fontFamily,
                        color: run.color,
                        "font-weight": run.bold ? 700 : 400,
                        "font-style": run.italic ? "italic" : "normal",
                        "text-decoration": run.underline
                          ? "underline"
                          : "none",
                      }}
                    >
                      {replaceVariables(
                        run.text,
                        props.values,
                        props.showVariables,
                      )}
                    </span>
                  )}
                </For>
              </Show>
            </div>
          );
        }}
      </For>
    </div>
  );
}
