import CustomDataTable from "@/widgets/table/ui/CustomDataTable";
import { Eye, Upload, Database } from "lucide-react";
import { formatDate } from "../../../shared/lib/date-utils";
import { schemaActions } from "../store";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { publishSchemaFn, schemasQueryOptions } from "../api";
import { SchemaDialog } from "./SchemaDialog";
import { Badge } from "@/shared/ui/shadcn/badge";
import { useState } from "react";
import { PublishConfirmDialog } from "./PublishConfirmDialog";
import type { Schema } from "../model/types";
import { toast } from "sonner";
import { useNavigate } from "@tanstack/react-router";
import { navigationActions } from "@/features/navigation";

interface PropsI {
  project_id: string;
}

export default function SchemaTable({ project_id }: PropsI) {
  const navigate = useNavigate({from: "/projects/config"});
  const queryClient = useQueryClient();
  const { data = [] } = useQuery(schemasQueryOptions(project_id))
  const [isPublishConfirmOpen, setIsPublishConfirmOpen] = useState(false);
  const [selectedSchemaToPublish, setSelectedSchemaToPublish] = useState<Schema | null>(null);

  const publishSchemaMutation = useMutation({
    mutationFn: publishSchemaFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["schemas"] });
      } else toast.error(`Failed to publish schema: ${response.message}`);
    },
    onError: (error) => {
      toast.error(`Failed to publish schema: ${error.message}`);
    },
  });

  const handlePublish = (schema: Schema) => {
    setSelectedSchemaToPublish(schema);
    setIsPublishConfirmOpen(true);
  };

  const confirmPublish = () => {
    if (selectedSchemaToPublish) publishSchemaMutation.mutate(selectedSchemaToPublish)
    setIsPublishConfirmOpen(false);
    setSelectedSchemaToPublish(null);
  };

  return (
    <>
      <CustomDataTable
        data={data}
        columns={[
          {
            key: "title",
            header: "Title",
            sortable: true,
          },
          {
            key: "flow_id",
            header: "Flow ID",
            sortable: true,
            render: (value) => <Badge variant="outline">{value}</Badge>
          },
          {
            key: "status",
            header: "Status",
            sortable: true,
            render: (value) => {
              let variant: "default" | "secondary" | "outline" = "default";
              switch (value) {
                case "draft":
                  variant = "outline";
                  break;
                case "archived":
                  variant = "secondary";
                  break;
                default:
                  variant = "default";
                  break;
              }
              return <Badge variant={variant}>{value}</Badge>
            }
          },
          {
            key: "type",
            header: "Type",
            sortable: true,
            render: (value) => <Badge variant="default">{value}</Badge>
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
        rowActions={[
          {
            label: "Inspect",
            onClick: (row) => {
              navigationActions.setCurrentSchemaId(row.id);
              navigationActions.setCurrentSchemaVersion(null);
              navigate({to: '/schemas/editor'})
            },
            icon: Eye,
            variant: "ghost-primary",
          },
          {
            label: "Publish",
            onClick: handlePublish,
            icon: Upload,
            variant: "ghost-primary",
          },
        ]}
        tableActions={[
          {
            label: "Create Schema",
            icon: Database,
            onClick: schemaActions.openCreate,
            variant: "solid"
          }
        ]}
      />
      <SchemaDialog project_id={project_id}/>
      {selectedSchemaToPublish && (
        <PublishConfirmDialog
          isOpen={isPublishConfirmOpen}
          onClose={() => setIsPublishConfirmOpen(false)}
          onConfirm={confirmPublish}
          schemaTitle={selectedSchemaToPublish.title}
        />
      )}
    </>
  )
}