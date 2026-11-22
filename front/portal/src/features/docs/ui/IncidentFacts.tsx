import { For, Show } from "solid-js";
import { formatMoment, formatSpan, phasesOf, serviceName } from "../incident";
import type { DocMeta } from "../types";

/** The summary as a lead, then when the incident ran and what it hit. */
export function IncidentFacts(props: { meta: DocMeta }) {
  const phases = () => phasesOf(props.meta);

  return (
    <div class="max-w-[68ch]">
      <Show when={props.meta.summary}>
        <p class="mb-8 text-[17px] leading-relaxed text-foreground">{props.meta.summary}</p>
      </Show>
      <Show when={phases().length > 0 || props.meta.services?.length}>
        <dl class="mb-10 grid gap-x-6 gap-y-4 border-y border-border py-4 text-[14px] sm:grid-cols-[110px_minmax(0,1fr)]">
          <Show when={phases().length > 0}>
            <dt class="text-muted-foreground">When</dt>
            <dd>
              <ol class="grid grid-cols-[auto_12px_minmax(0,1fr)] gap-x-3">
                <For each={phases()}>
                  {(phase, i) => (
                    <li class="contents">
                      <Show when={i() > 0}>
                        <span />
                        <span class="mx-auto w-px bg-border" aria-hidden="true" />
                        <span class="py-1.5 text-[13px] text-muted-foreground">
                          {formatSpan(phases()[i() - 1], phase)} later
                        </span>
                      </Show>
                      <span class="font-medium">{phase.label}</span>
                      <span class="relative flex items-center justify-center" aria-hidden="true">
                        <Show when={i() > 0}>
                          <span class="absolute top-0 left-1/2 h-1/2 w-px -translate-x-1/2 bg-border" />
                        </Show>
                        <Show when={i() < phases().length - 1}>
                          <span class="absolute bottom-0 left-1/2 h-1/2 w-px -translate-x-1/2 bg-border" />
                        </Show>
                        <span
                          class={`relative size-2.5 rounded-full border-2 border-foreground bg-background ${i() === phases().length - 1 ? "!bg-foreground" : ""}`}
                        />
                      </span>
                      <time datetime={phase.at} class="tabular-nums">
                        {formatMoment(phase, i() === 0)}
                      </time>
                    </li>
                  )}
                </For>
              </ol>
            </dd>
          </Show>
          <Show when={props.meta.services?.length}>
            <dt class="text-muted-foreground">Affected</dt>
            <dd>{props.meta.services!.map(serviceName).join(", ")}</dd>
          </Show>
        </dl>
      </Show>
    </div>
  );
}
