import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { permissionsQueryOptions } from "../api";
import { useQuery } from "@tanstack/react-query";
import { formatDate } from "@/shared/lib/date-utils";
import { Shield } from "lucide-react";
import { permissionActions } from "../store";
import PermissionDialog from "./PermissionDialog";
import { Badge } from "@/shared/ui/shadcn/badge";
import TruncatedId from "@/shared/ui/TruncatedId";

interface PropsI {
  project_id: string;
}

export default function PermissionTable({ project_id }: PropsI) {
  const { data = [] } = useQuery(permissionsQueryOptions(project_id))
  return (
    <>
      <CustomDataTable
        data={data}
        columns={[
          {
            key: "object",
            header: "Object",
            sortable: true,
            render: (value) => {
              const val = value as string;
              if (val === "*") return <Badge variant="secondary">{val}</Badge>
              return <Badge variant="outline">{val}</Badge>
            },
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
            searchableTextExtractor: (value) => formatDate(value as string),
          },
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