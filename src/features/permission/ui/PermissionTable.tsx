import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { permissionsQueryOptions } from "../api";
import { useQuery } from "@tanstack/react-query";
import { formatDate } from "@/shared/lib/date-utils";
import { Edit, Shield, Trash2 } from "lucide-react";
import { permissionActions } from "../store";
import PermissionDialog from "./PermissionDialog";
import { Badge } from "@/shared/ui/shadcn/badge";
import TruncatedId from "@/shared/ui/TruncatedId";
import { MetadataVisualizer } from "@/shared/ui/MetadataVisualizer";
import type { Permission } from "../model/types";

interface PropsI {
  project_id: string;
}

/**
 * Extended permission type to include flattened meta fields for table filtering.
 */
interface FlattenedPermission extends Permission {
  status: string;
}

export default function PermissionTable({ project_id }: PropsI) {
  const { data = [] } = useQuery(permissionsQueryOptions(project_id))

  const tableData: FlattenedPermission[] = data.map((p) => {
    return {
      ...p,
      status: p.meta?.status || "active",
    };
  });

  return (
    <>
      <CustomDataTable<FlattenedPermission>
        data={tableData}
        searchPlaceholder="Search permissions by object, action or status..."
        columns={[
          {
            key: "object",
            header: "Object",
            sortable: true,
            searchableTextExtractor: (_, row) => `${row.object} ${row.status}`,
            render: (value, row) => (
              <div className="flex items-center gap-2">
                <MetadataVisualizer name={String(value)} meta={row.meta} />
              </div>
            ),
          },
          {
            key: "action",
            header: "Action",
            sortable: true,
            render: (value) => {
              const val = value as string;
              if (val === "*") return <Badge variant="secondary">{val}</Badge>;
              return <Badge variant="outline">{val}</Badge>
            },
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
            onClick: permissionActions.openEdit,
            icon: Edit,
            variant: "ghost-primary",
          },
          {
            label: "Delete",
            icon: Trash2,
            onClick: permissionActions.openDelete,
            variant: "destructive",
          }
        ]}
        tableActions={[
          {
            label: "Create Permission",
            icon: Shield,
            onClick: permissionActions.openCreate,
            variant: "solid"
          }
        ]}
      />
      <PermissionDialog project_id={project_id}/>
    </>
  )
}
