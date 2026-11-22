import { Link } from "@tanstack/solid-router";
import { Show } from "solid-js";
import ArrowLeftIcon from "~icons/lucide/arrow-left";
import ArrowRightIcon from "~icons/lucide/arrow-right";
import type { DocLink } from "./catalog";

// Decisions are listed newest first, so "next" would mean older; say so.
const WORDS = {
  decisions: { previous: "Newer", next: "Older" },
  default: { previous: "Previous", next: "Next" },
};

export function DocPager(props: { previous?: DocLink; next?: DocLink }) {
  const category = () => (props.previous ?? props.next)?.category;
  const words = () => (category() === "decisions" ? WORDS.decisions : WORDS.default);

  return (
    <Show when={props.previous || props.next}>
      <nav aria-label="More docs" class="mt-14 grid max-w-[68ch] grid-cols-2 gap-3 border-t border-border pt-6">
        <Show when={props.previous}>
          {(doc) => <PagerLink doc={doc()} word={words().previous} side="previous" />}
        </Show>
        <Show when={props.next}>
          {(doc) => <PagerLink doc={doc()} word={words().next} side="next" />}
        </Show>
      </nav>
    </Show>
  );
}

function PagerLink(props: { doc: DocLink; word: string; side: "previous" | "next" }) {
  const next = () => props.side === "next";

  return (
    <Link
      to="/docs/$category/$slug"
      params={{ category: props.doc.category, slug: props.doc.slug }}
      class={`group -mx-3 flex min-w-0 items-start gap-3 rounded-lg px-3 py-3 hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring ${next() ? "col-start-2 flex-row-reverse text-right" : ""}`}
    >
      <span
        class={`mt-[3px] shrink-0 text-muted-foreground transition-transform duration-150 group-hover:text-primary motion-reduce:transition-none ${next() ? "group-hover:translate-x-0.5" : "group-hover:-translate-x-0.5"}`}
      >
        {next() ? <ArrowRightIcon class="size-4" /> : <ArrowLeftIcon class="size-4" />}
      </span>
      <span class="grid min-w-0 gap-0.5">
        <span class="text-[13px] text-muted-foreground">
          {props.word}
          <Show when={props.doc.meta.number}>
            <span class="ml-2 font-mono text-[12px]">{props.doc.meta.number}</span>
          </Show>
        </span>
        <span class="font-medium leading-snug group-hover:text-primary wrap-anywhere">
          {props.doc.meta.title ?? props.doc.slug}
        </span>
      </span>
    </Link>
  );
}
