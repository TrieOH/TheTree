import RolePermissionsEditor from "./RolePermissionsEditor";
import { useQuery } from "@tanstack/react-query";
import { roleQueryOptions } from "../api";
import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { formatDate } from "@/shared/lib/date-utils";
import { roleActions } from "../store";
import { Edit, ShieldCheck, Trash2 } from "lucide-react";
import RoleDialog from "./RoleDialog";
import TruncatedId from "@/shared/ui/TruncatedId";
import { MetadataVisualizer } from "@/shared/ui/MetadataVisualizer";
import type { Role } from "../model/types";


interface PropsI {
  project_id: string;
}

/**
 * Extended role type to include flattened meta fields for table filtering.
 */
interface FlattenedRole extends Role {
  status: string;
}

export default function RoleTable({ project_id }: PropsI) {
  const { data = [] } = useQuery(roleQueryOptions(project_id));

  const tableData: FlattenedRole[] = data.map((r) => {
    return {
      ...r,
      status: r.meta?.status || "active",
    };
  });

  return (
    <>
      <CustomDataTable<FlattenedRole>
        data={tableData}
        searchPlaceholder="Search roles by name or status..."
        columns={[
          {
            key: "name",
            header: "Name",
            sortable: true,
            searchableTextExtractor: (_, row) => `${row.name} ${row.status}`,
            render: (value, row) => (
              <div className="flex items-center gap-2">
                <MetadataVisualizer name={String(value)} meta={row.meta} />
              </div>
            ),
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
            searchableTextExtractor: (value) => formatDate(value as string),
          },
          {
            key: "updated_at",
            header: "Updated At",
            sortable: true,
            render: (value) => formatDate(value as string),
            searchableTextExtractor: (value) => formatDate(value as string),
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
        renderExpandedRow={(row) => <RolePermissionsEditor project_id={project_id} role={row} />}
        rowActions={[
          {
            label: "Update",
            onClick: roleActions.openEdit,
            icon: Edit,
            variant: "ghost-primary",
          },
          {
            label: "Delete",
            icon: Trash2,
            onClick: roleActions.openDelete,
            variant: "destructive",
          }
        ]}
        tableActions={[
          {
            label: "Create Role",
            icon: ShieldCheck,
            onClick: roleActions.openCreate,
            variant: "solid"
          }
        ]}
      />
      <RoleDialog project_id={project_id}/>
    </>
  )
}