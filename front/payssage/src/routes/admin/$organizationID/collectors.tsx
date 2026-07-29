import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useLayoutHeader } from "@trieoh/ui-base";
import { useMemo } from "react";
import { collectorsQueryOptions } from "#/features/collectors/api";
import { ProviderConnectSection } from "#/features/oauth/ui/provider-connect-button";
import { ProviderCredentialList } from "#/shared/ui/provider-credential-list";

export const Route = createFileRoute("/admin/$organizationID/collectors")({
  component: CollectorsPage,
});

function CollectorsPage() {
  const { organizationID } = Route.useParams();
  const { data = [] } = useQuery(collectorsQueryOptions(organizationID));
  const header = useMemo(
    () => (
      <div>
        <h1 className="text-lg font-semibold">Collectors</h1>
        <p className="text-sm text-muted-foreground">
          {data.length} provider account{data.length === 1 ? "" : "s"} shared by
          this organization.
        </p>
      </div>
    ),
    [data.length],
  );
  useLayoutHeader(header);

  return (
    <div className="space-y-6">
      <ProviderConnectSection
        flow="collector"
        organizationId={organizationID}
      />
      <ProviderCredentialList
        items={data}
        flow="collector"
        queryKey={collectorsQueryOptions(organizationID).queryKey}
      />
    </div>
  );
}
