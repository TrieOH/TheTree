import { createFileRoute } from "@tanstack/solid-router";
import { Tabs } from "@trieoh/ui-solid";
import { NotConnected, PageHeader } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/delivery")({
  head: () => ({ meta: [{ title: "Delivery · TrieOH dev" }] }),
  component: DeliveryPage,
});

const TABS = [
  { id: "production", label: "Production" },
  { id: "pulls", label: "Pull requests" },
  { id: "ci", label: "CI" },
  { id: "tags", label: "Tags" },
] as const;

function DeliveryPage() {
  return (
    <>
      <PageHeader
        title="Delivery"
        description="What's running, what's about to ship and how CI is doing. From GitHub and Cloudflare."
      />
      <Tabs
        ariaLabel="Delivery views"
        items={TABS.map((tab) => ({
          ...tab,
          content: <NotConnected source="GitHub and the deploy workflows" />,
        }))}
        panelClass="pt-6"
      />
    </>
  );
}
