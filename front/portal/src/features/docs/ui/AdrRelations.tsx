import { Link } from "@tanstack/solid-router";
import { For, Show } from "solid-js";
import { RELATION_LABELS, relationsOf } from "../adr";
import { decisionByNumber } from "../catalog";
import type { DocMeta } from "../types";

/** "ADR-0011", linked when the ADR exists. With `titled`, its title follows. */
export function AdrRef(props: { number: string; titled?: boolean }) {
  const doc = () => decisionByNumber(props.number);
  return (
    <Show when={doc()} fallback={<span class="font-mono text-[12.5px]">{props.number}</span>}>
      {(found) => (
        <Link
          to="/docs/$category/$slug"
          params={{ category: "decisions", slug: found().slug }}
          title={props.titled ? undefined : found().meta.title}
          class="rounded-sm underline decoration-border underline-offset-[3px] hover:decoration-current focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
        >
          <span class="font-mono text-[12.5px]">{props.number}</span>
          <Show when={props.titled && found().meta.title}>
            {(title) => <span class="ml-2">{title()}</span>}
          </Show>
        </Link>
      )}
    </Show>
  );
}

/** One line per relation: "Amends ADR-0016, ADR-0017". */
export function AdrRelationLines(props: { meta: DocMeta; class?: string }) {
  return (
    <For each={relationsOf(props.meta)}>
      {([kind, numbers]) => (
        <span class={props.class}>
          {RELATION_LABELS[kind]}{" "}
          <For each={numbers}>
            {(number, i) => (
              <>
                <Show when={i() > 0}>{", "}</Show>
                <AdrRef number={number} />
              </>
            )}
          </For>
        </span>
      )}
    </For>
  );
}
