import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Upload } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { publishSchemaVersionFn } from "../api";
import { toast } from "sonner";
import { useStore as useNavigationStore } from "@tanstack/react-store";
import { navigationStore } from "@/features/navigation";
import { useParams } from "@tanstack/react-router";

interface PublishSchemaVersionButtonProps {
  hasChanges: boolean;
  isMobile: boolean;
}

export default function PublishSchemaVersionButton({ 
  hasChanges, 
  isMobile,
}: PublishSchemaVersionButtonProps) {
  const queryClient = useQueryClient();
  const { projectId: currentProjectId, schemaId: currentSchemaId } = useParams({ strict: false });
  const { currentSchemaVersion } = useNavigationStore(navigationStore);

  const publishSchemaVersionMutation = useMutation({
    mutationFn: publishSchemaVersionFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["latestSchemaVersion", currentProjectId, currentSchemaId] });
        queryClient.invalidateQueries({ queryKey: ["currentSchemaVersion", currentProjectId, currentSchemaId] });
        queryClient.invalidateQueries({ queryKey: ["schemaVersionById", currentProjectId, currentSchemaId, currentSchemaVersion] });
      }
    },
    onError: (error) => {
      toast.error(`Failed to publish schema version: ${error.message}`);
    }
  });

  const handleSubmit = () => {
    publishSchemaVersionMutation.mutate({
      project_id: currentProjectId || "",
      schema_id: currentSchemaId || "",
    })
  };

  return (
    <ShadowButton
      onClick={handleSubmit}
      disabled={!hasChanges}
      variant="solid"
      value={isMobile ? '' : 'Publish Version'}
      leftIcon={<Upload className="w-4 h-4" />}
    />
  );
}
