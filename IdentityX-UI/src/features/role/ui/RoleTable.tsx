import RolePermissionsEditor from "./RolePermissionsEditor";
import { useQuery } from "@tanstack/react-query";
import { roleQueryOptions } from "../api";
import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { formatDate } from "@/shared/lib/date-utils";
import { roleActions } from "../store";
import { Edit, ShieldCheck, Trash2 } from "lucide-react";
import RoleDialog from "./RoleDialog";
import TruncatedId from "@/shared/ui/TruncatedId";


interface PropsI {
  project_id: string;
}

export default function RoleTable({ project_id }: PropsI) {
  const { data = [] } = useQuery(roleQueryOptions(project_id));

  return (
    <>
      <CustomDataTable
        data={data}
        columns={[
          {
            key: "name",
            header: "Name",
            sortable: true,
          },
          {
            key: "description",
            header: "Description",
            sortable: true,
            render: (value) => (
              <p 
                title={value}
                className="max-w-64 line-clamp-3 whitespace-normal"
              >
                {value}
              </p>
            )
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