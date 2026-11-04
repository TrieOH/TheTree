import { CodeBlock, cn } from "@trieoh/ui-solid";
import type { JSX } from "@solidjs/web";
import { For, Show, createSignal } from "solid-js";
import ChevronIcon from "~icons/lucide/chevron-right";
import type { ServiceId } from "@/features/catalog/services";
import { ProblemList } from "@/features/problems/ui/ProblemList";
import { operationAnchor, type ApiAuth, type ApiOperation, type ApiParameter } from "../model";
import { MethodBadge, PathText } from "./MethodBadge";
import { SchemaView, TypeText } from "./SchemaView";
import { TryPanel } from "./TryPanel";

const IN_ORDER: ApiParameter["in"][] = ["path", "query", "header", "cookie"];

function Block(props: { title: string; aside?: string; children: JSX.Element }) {
  return (
    <section class="space-y-2">
      <h4 class="flex items-baseline gap-2 text-[12px] font-medium uppercase tracking-wide text-muted-foreground">
        {props.title}
        <Show when={props.aside}>
          <span class="font-mono normal-case tracking-normal">{props.aside}</span>
        </Show>
      </h4>
      {props.children}
    </section>
  );
}

/** One endpoint, written out in full. Below `xl` it carries its own try-it panel. */
export function OperationSection(props: { op: ApiOperation; service: ServiceId; auth: readonly ApiAuth[] }) {
  const [problemsOpen, setProblemsOpen] = createSignal(false, { name: "problems-open" });
  const [tryOpen, setTryOpen] = createSignal(false, { name: "try-open" });
  const params = () => IN_ORDER.flatMap((where) => props.op.parameters.filter((param) => param.in === where));

  return (
    <article id={operationAnchor(props.op)} data-operation={props.op.id} class="scroll-mt-6 space-y-5 border-b border-border pb-10">
      <header class="space-y-2">
        <h3 class="font-heading text-[17px] font-semibold">{props.op.summary}</h3>
        <div class="flex flex-wrap items-center gap-2">
          <MethodBadge method={props.op.method} />
          <PathText path={props.op.path} />
        </div>
        <p class="text-[12.5px] text-muted-foreground">
          {props.op.auth.length ? `Auth: ${props.op.auth.join(" or ")}` : "Public: no authentication."}
          <Show when={props.op.scope}> · scope {props.op.scope}</Show>
          <Show when={props.op.deprecated}> · <span class="text-warn">deprecated</span></Show>
        </p>
      </header>

      <div class="prose-doc text-[13.5px]" innerHTML={props.op.description} />

      <Show when={params().length > 0}>
        <Block title="Parameters">
          <ul class="divide-y divide-border rounded-lg border border-border">
            <For each={params()}>
              {(param) => (
                <li class="px-3 py-2 text-[13px]">
                  <div class="flex flex-wrap items-baseline gap-x-2">
                    <code class="font-mono text-[12.5px] font-medium">{param.name}</code>
                    <span class="text-[11px] text-muted-foreground">{param.in}</span>
                    <TypeText node={param.schema} />
                    <Show when={param.required}>
                      <span class="text-[11px] text-warn">required</span>
                    </Show>
                  </div>
                  <Show when={param.description}>
                    <div class="prose-doc mt-0.5 text-muted-foreground" innerHTML={param.description} />
                  </Show>
                </li>
              )}
            </For>
          </ul>
        </Block>
      </Show>

      <Show when={props.op.requestBody}>
        {(body) => (
          <Block title="Request body" aside={`${body().contentType}${body().required ? "" : " · optional"}`}>
            <SchemaView node={body().schema} />
          </Block>
        )}
      </Show>

      <For each={props.op.responses}>
        {(response) => (
          <Block title={`Response ${response.status}`} aside={response.contentType}>
            <Show when={response.description}>
              <p class="text-[13px] text-muted-foreground">{response.description}</p>
            </Show>
            <Show when={response.schema}>
              <SchemaView node={response.schema} />
            </Show>
            <Show when={response.example}>{(example) => <CodeBlock code={example()} filename="example.json" />}</Show>
          </Block>
        )}
      </For>

      <Show when={props.op.problems.length > 0}>
        <section>
          <button
            type="button"
            aria-expanded={problemsOpen() ? "true" : "false"}
            onClick={() => setProblemsOpen(!problemsOpen())}
            class="flex items-center gap-1.5 text-[12px] font-medium uppercase tracking-wide text-muted-foreground hover:text-foreground"
          >
            <ChevronIcon class={cn("size-3.5 transition-transform", problemsOpen() && "rotate-90")} />
            Problems <span class="font-mono normal-case">{props.op.problems.length}</span>
          </button>
          <Show when={problemsOpen()}>
            <ProblemList set={props.service} shared="include" only={props.op.problems} />
          </Show>
        </section>
      </Show>

      <section class="xl:hidden">
        <button
          type="button"
          aria-expanded={tryOpen() ? "true" : "false"}
          onClick={() => setTryOpen(!tryOpen())}
          class="flex items-center gap-1.5 text-[12px] font-medium uppercase tracking-wide text-muted-foreground hover:text-foreground"
        >
          <ChevronIcon class={cn("size-3.5 transition-transform", tryOpen() && "rotate-90")} />
          Try it
        </button>
        <Show when={tryOpen()}>
          <TryPanel op={props.op} auth={props.auth} class="mt-3 rounded-2xl border border-border bg-card p-4" />
        </Show>
      </section>
    </article>
  );
}
