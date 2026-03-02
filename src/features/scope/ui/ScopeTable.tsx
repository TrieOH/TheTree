import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { Edit, Globe } from "lucide-react";
import { formatDate } from "../../../shared/lib/date-utils";
import ScopeDialog from "./ScopeDialog";
import { scopeActions } from "../store";
import { useQuery } from "@tanstack/react-query";
import { scopesQueryOptions } from "../api";
import { Badge } from "@/shared/ui/shadcn/badge";
import type { Scope } from "../model/types";
import TruncatedId from "@/shared/ui/TruncatedId";
import { MetadataVisualizer } from "@/shared/ui/MetadataVisualizer";

interface PropsI {
  project_id: string;
}

/**
 * Extended scope type to include flattened meta fields for table filtering.
 */
interface FlattenedScope extends Scope {
  status: string;
}

export default function ScopeTable({ project_id }: PropsI) {
  const { data = [] } = useQuery(scopesQueryOptions(project_id))

  const tableData: FlattenedScope[] = data.map((s) => {
    return {
      ...s,
      status: s.meta?.status || "active",
    };
  });

  return (
    <>
      <CustomDataTable<FlattenedScope>
        data={tableData} 
        searchPlaceholder="Search scopes by name, ID or status..."
        columns={[
          {
            key: "name",
            header: "Scope Identity",
            sortable: true,
            searchableTextExtractor: (_, row) => `${row.name} ${row.status} ${row.id}`,
            render: (value, row) => <MetadataVisualizer name={String(value)} meta={row.meta} />
          },
          {
            key: "type",
            header: "Category",
            sortable: true,
            render: (value) => {
              const type = value as Scope['type'];
              let variant: "default" | "secondary" | "destructive" | "outline" = "default";
              let displayType = type;
              switch (type) {
                case "global":
                  variant = "outline";
                  break;
                case "project_root":
                  variant = "secondary";
                  displayType = "Root";
                  break;
                case "project_scope":
                  variant = "default";
                  displayType = "Scope";
                  break;
              }
              return <Badge variant={variant}>{displayType}</Badge>;
            },
          },
          {
            key: "external_id",
            header: "External ID",
            sortable: true,
            render: (value) => (value ? (value as string) : "N/A"),
          },
          {
            key: "id",
            header: "ID",
            sortable: true,
            render: (value) => <TruncatedId id={value as string} />,
          },
          {
            key: "created_at",
            header: "Created At",
            sortable: true,
            render: (value) => formatDate(value as string),
          },
        ]}
        filters={[
          {
            key: "status",
            type: "select",
            label: "Status",
            placeholder: "Filter by status",
            options: [
              { label: "Active", value: "active" },
              { label: "Restricted", value: "restricted" },
              { label: "Beta", value: "beta" },
              { label: "Deprecated", value: "deprecated" }
            ]
          }
        ]}
        rowActions={[
          {
            label: "Update",
            onClick: scopeActions.openEdit,
            icon: Edit,
            variant: "ghost-primary",
          },
        ]}
        tableActions={[
          {
            label: "Create Scope",
            icon: Globe,
            onClick: scopeActions.openCreate,
            variant: "solid"
          }
        ]}
      />
      <ScopeDialog project_id={project_id}/>
    </>
  )
}
