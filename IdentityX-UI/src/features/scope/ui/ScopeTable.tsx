import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { Globe } from "lucide-react";
import { formatDate } from "../../../shared/lib/date-utils";
import type { Scope } from "@/features/scope/model/types";
import ScopeDialog from "./ScopeDialog";
import { scopeActions } from "../store";

interface PropsI {
  data: Scope[]
  project_id: string;
}

export default function ScopeTable({ data, project_id }: PropsI) {
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
            key: "type",
            header: "Type",
            sortable: true,
          },
          {
            key: "external_id",
            header: "External ID",
            sortable: true,
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