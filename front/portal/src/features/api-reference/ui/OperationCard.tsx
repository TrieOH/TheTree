import { CodeBlock, cn } from "@trieoh/ui-solid";
import { For, Show, createEffect, createSignal } from "solid-js";
import ChevronIcon from "~icons/lucide/chevron-right";
import LockIcon from "~icons/lucide/lock";
import type { ServiceId } from "@/features/catalog/services";
import { ProblemList } from "@/features/problems/ui/ProblemList";
import { operationAnchor, type ApiOperation, type ApiParameter } from "../model";
import { MethodBadge, PathText } from "./MethodBadge";
import { SchemaView, TypeText } from "./SchemaView";

const IN_ORDER: ApiParameter["in"][] = ["path", "query", "header", "cookie"];

function Block(props: { title: string; aside?: string; children: unknown }) {
  return (
    <section class="space-y-2">
      <h4 class="flex items-baseline gap-2 text-[12px] font-medium uppercase tracking-wide text-muted-foreground">
        {props.title}
        <Show when={props.aside}>
          <span class="font-mono normal-case tracking-normal">{props.aside}</span>
        </Show>
      </h4>
      {props.children as never}
    </section>
  );
}

/** One endpoint: a one-line header that opens to everything about it. */
export function OperationCard(props: { op: ApiOperation; service: ServiceId; open?: boolean }) {
  const [open, setOpen] = createSignal(false, { name: "operation-open" });
  // A link straight to this endpoint (#op-…) opens it.
  createEffect(
    () => props.open ?? false,
    (wanted) => {
      if (wanted) setOpen(true);
    },
  );
  const params = () =>
    IN_ORDER.map((where) => ({ where, list: props.op.parameters.filter((param) => param.in === where) })).filter(
      (group) => group.list.length > 0,
    );

  return (
    <article id={operationAnchor(props.op)} class="scroll-mt-20 rounded-xl border border-border">
      <button
        type="button"
        aria-expanded={open() ? "true" : "false"}
        onClick={() => setOpen(!open())}
        class="flex w-full items-center gap-3 px-3 py-2.5 text-left hover:bg-foreground/5"
      >
        <MethodBadge method={props.op.method} />
        <PathText path={props.op.path} />
        <span class="ml-auto hidden truncate pl-3 text-[13px] text-muted-foreground md:inline">{props.op.summary}</span>
        <Show when={props.op.auth.length > 0}>
          <LockIcon aria-label="Needs authentication" class="size-3.5 shrink-0 text-muted-foreground" />
        </Show>
        <ChevronIcon class={cn("size-4 shrink-0 text-muted-foreground transition-transform", open() && "rotate-90")} />
      </button>

      <Show when={open()}>
        <div class="space-y-6 border-t border-border px-4 py-4">
          <div>
            <p class="font-medium md:hidden">{props.op.summary}</p>
            <Show when={props.op.deprecated}>
              <p class="mb-2 text-[13px] text-warn">Deprecated.</p>
            </Show>
            <div class="prose-doc text-[13.5px]" innerHTML={props.op.description} />
            <p class="mt-2 text-[12.5px] text-muted-foreground">
              {props.op.auth.length ? `Auth: ${props.op.auth.join(" or ")}` : "Public: no authentication."}
              <Show when={props.op.scope}> · scope {props.op.scope}</Show>
              {" · "}
              <code class="font-mono">{props.op.id}</code>
            </p>
          </div>

          <Show when={params().length > 0}>
            <Block title="Parameters">
              <ul class="divide-y divide-border rounded-lg border border-border">
                <For each={params().flatMap((group) => group.list)}>
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
                <Show when={body().example}>{(example) => <CodeBlock code={example()} filename="example.json" />}</Show>
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
            <Block title="Problems" aside={`${props.op.problems.length}`}>
              <ProblemList set={props.service} shared="include" only={props.op.problems} />
            </Block>
          </Show>
        </div>
      </Show>
    </article>
  );
}
