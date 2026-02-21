import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { Globe } from "lucide-react";
import { formatDate } from "../../../shared/lib/date-utils";
import ScopeDialog from "./ScopeDialog";
import { scopeActions } from "../store";
import { useQuery } from "@tanstack/react-query";
import { scopesQueryOptions } from "../api";
import { Badge } from "@/shared/ui/shadcn/badge";
import type { Scope } from "../model/types";
import TruncatedId from "@/shared/ui/TruncatedId";

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
            searchableTextExtractor: (value) => (value ? (value as string) : "N/A"),
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