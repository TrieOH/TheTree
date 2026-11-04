import { Dynamic, type JSX } from "@solidjs/web";
import { omit } from "solid-js";
import LinkIcon from "~icons/lucide/link";
import { CopyButton } from "@/shared/ui/CopyButton";

// Overrides for the plain elements MDX emits.

type Props = { id?: string; children?: JSX.Element } & Record<string, unknown>;

/** h2–h4 get a hover anchor to their id (mdx/rehype-heading-ids.ts). */
function heading(Tag: "h2" | "h3" | "h4") {
  return (props: Props) => (
    <Dynamic component={Tag} id={props.id} class="group scroll-mt-20">
      {props.children}
      {props.id ? (
        <a
          href={`#${props.id}`}
          aria-label="Link to this section"
          class="ml-2 inline-flex align-middle text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100 focus:opacity-100"
        >
          <LinkIcon class="size-3.5" />
        </a>
      ) : null}
    </Dynamic>
  );
}

/** Highlighted code (shiki, at build time) with a copy button. */
function Pre(props: Props) {
  let pre: HTMLPreElement | undefined;
  const rest = omit(props, "children");
  return (
    <div class="group relative">
      <pre ref={(el) => (pre = el)} {...(rest as JSX.HTMLAttributes<HTMLPreElement>)}>
        {props.children}
      </pre>
      <CopyButton
        text={() => pre?.textContent ?? ""}
        label="Copy code"
        class="absolute right-2 top-2 border border-border bg-background/80 opacity-0 focus:opacity-100 group-hover:opacity-100"
      />
    </div>
  );
}

export const elements = { h2: heading("h2"), h3: heading("h3"), h4: heading("h4"), pre: Pre };
