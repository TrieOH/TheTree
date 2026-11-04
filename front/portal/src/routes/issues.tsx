import { createFileRoute } from "@tanstack/solid-router";
import { Tabs } from "@trieoh/ui-solid";
import { NotConnected, PageHeader } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/issues")({
  head: () => ({ meta: [{ title: "Issues · TrieOH dev" }] }),
  component: IssuesPage,
});

const TABS = [
  { id: "current", label: "Current cycle" },
  { id: "projects", label: "Projects" },
  { id: "bugs", label: "Bugs" },
  { id: "done", label: "Done" },
] as const;

function IssuesPage() {
  return (
    <>
      <PageHeader
        title="Issues"
        description="Synced from Linear as they change. Open any issue in Linear to edit it."
      />
      <Tabs
        ariaLabel="Issue views"
        items={TABS.map((tab) => ({ ...tab, content: <NotConnected source="Linear" /> }))}
        panelClass="pt-6"
      />
    </>
  );
}
