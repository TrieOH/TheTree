import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { CopyButton } from "@/shared/ui/CopyButton";
import { exampleBody, problemType, statusText, type ProblemDoc } from "../model";

/**
 * A problem as it arrives: the standard RFC 9457 members with the values this
 * problem gives them, then the members it adds. The JSON is a copy away.
 */
export function ProblemDetails(props: { problem: ProblemDoc }) {
  return (
    <div class="not-prose overflow-hidden rounded-xl border border-border bg-card text-[13px]">
      <header class="flex items-center gap-2 border-b border-border px-3 py-2 text-[12px] text-muted-foreground">
        <span class="font-mono">
          {props.problem.status} {statusText(props.problem.status)}
        </span>
        <span>·</span>
        <span class="font-mono">application/problem+json</span>
        <span class="ml-auto">Copy example</span>
        <CopyButton text={() => exampleBody(props.problem)} label="Copy example body" class="-my-1 -mr-1" />
      </header>
      <dl class="grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1.5 px-3 py-2.5">
        <Member name="type">
          <span class="break-all font-mono text-[12.5px]">{problemType(props.problem)}</span>
        </Member>
        <Member name="status">
          <span class="font-mono text-[12.5px] tabular-nums">{props.problem.status}</span>
        </Member>
        <Member name="title">{props.problem.title}</Member>
        <Member name="detail">
          <span class="text-muted-foreground">varies: the message of the error that matched</span>
        </Member>
        <Member name="instance">
          <span class="text-muted-foreground">the request's path</span>
        </Member>
      </dl>
      <Show when={props.problem.members.length > 0}>
        <div class="border-t border-dashed border-border px-3 pb-2.5 pt-2">
          <p class="mb-1.5 text-[11px] font-medium uppercase tracking-wide text-muted-foreground">
            This problem adds
          </p>
          <dl class="grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1.5">
            <For each={props.problem.members}>
              {(member) => (
                <Member name={member.name} type={member.type}>
                  <span class="text-muted-foreground">{member.description}</span>
                </Member>
              )}
            </For>
          </dl>
        </div>
      </Show>
    </div>
  );
}

function Member(props: { name: string; type?: string; children: JSX.Element }) {
  return (
    <>
      <dt class="flex items-baseline gap-1.5 font-mono text-[12.5px]">
        {props.name}
        <Show when={props.type}>
          <span class="rounded bg-muted px-1 text-[10.5px] text-muted-foreground">{props.type}</span>
        </Show>
      </dt>
      <dd class="min-w-0">{props.children}</dd>
    </>
  );
}
