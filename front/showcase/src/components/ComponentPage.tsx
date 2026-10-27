import { CodeBlock } from "@trieoh/ui-solid";
import type { Entry, Values } from "../catalog";
import { usage } from "./usage";
import { PropertiesPanel } from "./PropertiesPanel";
import { Installation } from "./Installation";
import { ApiReference } from "./ApiReference";
import { ComponentPagination } from "./ComponentPagination";

export function ComponentPage(props: {
  entry: Entry;
  values: Values;
  previous?: Entry;
  next?: Entry;
  onUpdate: (key: string, value: string | boolean) => void;
  onReset: () => void;
}) {
  return (
    <main class="mx-auto w-full min-w-0 max-w-300 px-5 pt-6 pb-10 min-[761px]:px-9 min-[761px]:pt-12 min-[761px]:pb-18">
      <header class="mb-9">
        <div class="mb-4.5 flex items-center gap-2.5 text-[11px] font-medium">
          <span class="rounded-md border border-border bg-card px-2.25 py-1">
            {props.entry.group}
          </span>
          <span class="text-muted-foreground">Solid</span>
        </div>
        <h1 class="mb-3.5 font-heading text-[clamp(32px,4vw,44px)] leading-[1.15] font-semibold tracking-[-1.7px] wrap-anywhere">
          {props.entry.name}
        </h1>
        <p class="max-w-155 text-sm leading-[1.8] text-muted-foreground">
          {props.entry.description}
        </p>
      </header>
      <section class="min-w-0 rounded-[14px] border border-border bg-card">
        <div class="flex items-center justify-between gap-4 border-b border-border px-5 py-4">
          <h2 class="text-[13px] font-semibold">Prévia</h2>
        </div>
        <div class="flex min-h-50 items-center justify-center rounded-b-[14px] bg-background bg-[radial-gradient(var(--border)_0.75px,transparent_0.75px)] bg-size-[18px_18px] p-6 min-[761px]:min-h-85 min-[761px]:px-7 min-[761px]:py-10 has-data-wide:p-4 *:min-w-0 *:max-w-full [&>input:not([type=checkbox]):not([type=radio])]:w-full [&>input:not([type=checkbox]):not([type=radio])]:max-w-120 [&>textarea]:w-full [&>textarea]:max-w-120 [&>div:not([data-wide])]:w-full [&>div:not([data-wide])]:max-w-120">
          {props.entry.render(props.values, (key, value) =>
            props.onUpdate(key, value),
          )}
        </div>
      </section>
      <PropertiesPanel
        entry={props.entry}
        values={props.values}
        onUpdate={props.onUpdate}
        onReset={props.onReset}
      />
      <Installation />
      <section class="mt-6 min-w-0 min-[761px]:mt-8">
        <div class="mb-4 flex items-center justify-between gap-4">
          <h2 class="text-[13px] font-semibold">Código</h2>
          <span class="rounded-[5px] border border-border px-1.75 py-1 font-mono text-[10px] text-muted-foreground">
            TSX
          </span>
        </div>
        <CodeBlock code={usage(props.entry.name, props.values)} />
      </section>
      <ApiReference entry={props.entry} />
      <ComponentPagination previous={props.previous} next={props.next} />
    </main>
  );
}
