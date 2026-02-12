import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { Globe } from "lucide-react";
import { formatDate } from "../../../shared/lib/date-utils";
import ScopeDialog from "./ScopeDialog";
import { scopeActions } from "../store";
import { useQuery } from "@tanstack/react-query";
import { scopesQueryOptions } from "../api";

interface PropsI {
  project_id: string;
}

export default function ScopeTable({ project_id }: PropsI) {
  const { data = [] } = useQuery(scopesQueryOptions(project_id))
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