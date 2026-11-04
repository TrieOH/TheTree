import { cn } from "@trieoh/ui-solid";
import { For, Show, createSignal } from "solid-js";
import ChevronIcon from "~icons/lucide/chevron-right";
import { refOf, schemaAnchor, typeLabel, type SchemaNode } from "../model";

/** A type, linking to its model when it is (or holds) a named one. */
export function TypeText(props: { node: SchemaNode | null | undefined }) {
  const ref = () => refOf(props.node);
  return (
    <Show
      when={ref()}
      fallback={<span class="font-mono text-[12px] text-muted-foreground">{typeLabel(props.node)}</span>}
    >
      {(name) => (
        <a href={`#${schemaAnchor(name())}`} class="font-mono text-[12px] text-accent hover:underline">
          {typeLabel(props.node)}
        </a>
      )}
    </Show>
  );
}

/** The properties of an object, a level at a time; named types link instead of nesting. */
export function SchemaView(props: { node: SchemaNode | null | undefined; class?: string }) {
  const node = () => props.node;
  const properties = () => {
    const current = node();
    if (!current || "ref" in current) return [];
    const object = current.type === "array" && current.items && !("ref" in current.items) ? current.items : current;
    return Object.entries(object.properties ?? {}).map(([name, child]) => ({
      name,
      child,
      required: (object.required ?? []).includes(name),
    }));
  };

  return (
    <Show when={properties().length > 0} fallback={<TypeText node={props.node} />}>
      <ul class={cn("divide-y divide-border rounded-lg border border-border", props.class)}>
        <For each={properties()}>{(property) => <PropertyRow {...property} />}</For>
      </ul>
    </Show>
  );
}

function PropertyRow(props: { name: string; child: SchemaNode; required: boolean }) {
  const [open, setOpen] = createSignal(false, { name: "schema-open" });
  const nested = () => {
    const child = props.child;
    if ("ref" in child) return false;
    const object = child.type === "array" && child.items && !("ref" in child.items) ? child.items : child;
    return Object.keys(object.properties ?? {}).length > 0;
  };
  const description = () => ("ref" in props.child ? undefined : props.child.description);

  return (
    <li class="px-3 py-2 text-[13px]">
      <div class="flex flex-wrap items-baseline gap-x-2 gap-y-0.5">
        <Show when={nested()} fallback={<span class="inline-block w-4" />}>
          <button
            type="button"
            aria-expanded={open() ? "true" : "false"}
            aria-label={open() ? `Hide ${props.name}` : `Show ${props.name}`}
            onClick={() => setOpen(!open())}
            class="grid size-4 place-items-center self-center text-muted-foreground hover:text-foreground"
          >
            <ChevronIcon class={cn("size-3.5 transition-transform", open() && "rotate-90")} />
          </button>
        </Show>
        <code class="font-mono text-[12.5px] font-medium">{props.name}</code>
        <TypeText node={props.child} />
        <Show when={props.required}>
          <span class="text-[11px] text-warn">required</span>
        </Show>
      </div>
      <Show when={description()}>
        <p class="mt-0.5 pl-6 text-muted-foreground">{description()}</p>
      </Show>
      <Show when={nested() && open()}>
        <div class="mt-2 pl-6">
          <SchemaView node={props.child} />
        </div>
      </Show>
    </li>
  );
}
