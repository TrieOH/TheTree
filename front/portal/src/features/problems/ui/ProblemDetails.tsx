import { CodeBlock } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import { exampleBody, type ProblemDoc, type ProblemSet } from "../model";

/** The problem's own members and what its body looks like on the wire. */
export function ProblemDetails(props: { problem: ProblemDoc; set: ProblemSet }) {
  return (
    <div class="space-y-3">
      <Show when={props.problem.members.length > 0}>
        <dl class="grid gap-x-4 gap-y-1 text-[13px] sm:grid-cols-[max-content_1fr]">
          <For each={props.problem.members}>
            {(member) => (
              <>
                <dt>
                  <code class="font-mono text-[12px]">{member.name}</code>
                  <span class="ml-1.5 text-[11px] text-muted-foreground">{member.type}</span>
                </dt>
                <dd class="text-muted-foreground">{member.description}</dd>
              </>
            )}
          </For>
        </dl>
      </Show>
      <CodeBlock code={exampleBody(props.problem, props.set)} filename="application/problem+json" />
    </div>
  );
}
