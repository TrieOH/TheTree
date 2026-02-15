import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Upload } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { publishSchemaVersionFn } from "../api"; // Use the new publishSchemaVersionFn
import { toast } from "sonner";
import { useStore as useNavigationStore } from "@tanstack/react-store";
import { navigationStore } from "@/features/navigation";
import type { VersionFieldList } from "../model/types";

interface PublishSchemaVersionButtonProps {
  items: VersionFieldList;
  hasChanges: boolean;
  setOriginalItems: (items: VersionFieldList) => void;
}

export default function PublishSchemaVersionButton({ items, hasChanges, setOriginalItems }: PublishSchemaVersionButtonProps) {
  const queryClient = useQueryClient();
  const { currentProjectId, currentSchemaId } = useNavigationStore(navigationStore);

  const publishSchemaVersionMutation = useMutation({
    mutationFn: publishSchemaVersionFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["latestSchemaVersion", currentProjectId, currentSchemaId] });
        queryClient.invalidateQueries({ queryKey: ["currentSchemaVersion", currentProjectId, currentSchemaId] });
        queryClient.invalidateQueries({ queryKey: ["schemaVersionById"] });

        setOriginalItems(items);
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
      value="Publish Version"
      variant="solid"
      leftIcon={<Upload className="w-4 h-4" />}
    />
  );
}
