import { Callout } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import { RELATION_LABELS, relationsOf } from "../adr";
import type { DocMeta } from "../types";
import { AdrRef } from "./AdrRelations";

/** Shown above a decision that a later one replaced. */
export function SupersededNotice(props: { meta: DocMeta }) {
  return (
    <Show when={props.meta.superseded_by}>
      {(number) => (
        <Callout tone="warning" class="mb-8" title="This decision was replaced">
          Read <AdrRef number={number()} titled /> for how things work now.
        </Callout>
      )}
    </Show>
  );
}

/** The summary as a lead, then what the decision touches and how it links to others. */
export function AdrFacts(props: { meta: DocMeta }) {
  // The replacement already has the notice above; don't list it twice.
  const relations = () => relationsOf(props.meta).filter(([kind]) => kind !== "superseded_by");

  return (
    <div class="max-w-[68ch]">
      <Show when={props.meta.summary}>
        <p class="mb-8 text-[17px] leading-relaxed text-foreground">{props.meta.summary}</p>
      </Show>
      <Show when={props.meta.scope?.length || relations().length}>
        <dl class="mb-10 grid gap-x-6 gap-y-3 border-y border-border py-4 text-[14px] sm:grid-cols-[110px_minmax(0,1fr)]">
          <Show when={props.meta.scope?.length}>
            <dt class="text-muted-foreground">Applies to</dt>
            <dd>
              <ul class="flex flex-wrap gap-x-4 gap-y-1">
                <For each={props.meta.scope}>
                  {(path) => <li class="font-mono text-[12.5px] wrap-anywhere">{path}</li>}
                </For>
              </ul>
            </dd>
          </Show>
          <For each={relations()}>
            {([kind, numbers]) => (
              <>
                <dt class="text-muted-foreground">{RELATION_LABELS[kind]}</dt>
                <dd>
                  <ul class="grid gap-1">
                    <For each={numbers}>{(number) => <li><AdrRef number={number} titled /></li>}</For>
                  </ul>
                </dd>
              </>
            )}
          </For>
        </dl>
      </Show>
    </div>
  );
}
