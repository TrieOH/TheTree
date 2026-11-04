import type { JSX } from "@solidjs/web";
import { Tabs } from "@trieoh/ui-solid";

/**
 * `<Tabs labels={["pnpm", "just"]}><Tab>…</Tab><Tab>…</Tab></Tabs>`: the
 * labels go on Tabs because MDX hands over children already rendered.
 */
export function DocTabs(props: { labels: string[]; children?: JSX.Element }) {
  const panels = () => {
    const children = props.children;
    return (Array.isArray(children) ? children : [children]).filter(
      (child) => !(typeof child === "string" && child.trim() === ""),
    );
  };
  return (
    <div class="not-prose my-5">
      <Tabs
        ariaLabel={props.labels.join(", ")}
        items={props.labels.map((label, i) => ({ id: String(i), label, content: panels()[i] }))}
        keepMounted
        panelClass="prose-doc pt-3"
      />
    </div>
  );
}

export function Tab(props: { children?: JSX.Element }) {
  return <div>{props.children}</div>;
}
