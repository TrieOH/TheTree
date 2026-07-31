import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { collectorsQueryOptions } from "#/features/collectors/api";
import { ProviderConnectSection } from "#/features/oauth/ui/provider-connect-button";
import { ProviderCredentialList } from "#/shared/ui/provider-credential-list";

export const Route = createFileRoute("/admin/collectors")({
  component: CollectorsPage,
});

function CollectorsPage() {
  const { data = [] } = useQuery(collectorsQueryOptions());
  const activeCount = data.filter((collector) => !collector.revoked_at).length;
  return (
    <div>
      <div className="border-b border-border/40 bg-background px-6 py-4">
        <div>
          <h1 className="text-lg font-semibold tracking-tight">
            My collectors
          </h1>
          <p className="text-sm text-muted-foreground">
            {data.length} collector{data.length === 1 ? "" : "s"} connected ·{" "}
            {activeCount} active
          </p>
        </div>
      </div>
      <div className="space-y-6 p-6">
        <ProviderConnectSection flow="collector" />
        <ProviderCredentialList
          items={data}
          flow="collector"
          queryKey={collectorsQueryOptions().queryKey}
        />
      </div>
    </div>
  );
}
