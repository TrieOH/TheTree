import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { FilePlus, Loader2 } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { navigationStore } from "@/features/navigation";
import { useStore } from "@tanstack/react-store";
import { toast } from "sonner";
import { createSchemaVersionDraftFn } from "../api";

export default function CreateSchemaVersionButton() {
  const { currentProjectId, currentSchemaId } = useStore(navigationStore);
  const queryClient = useQueryClient();

  const createSchemaVersionMutation = useMutation({
    mutationFn: createSchemaVersionDraftFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["latestSchemaVersion"] });
        queryClient.setQueryData(["latestSchemaVersion", response.data.id], response.data);
      } else toast.error(`Failed to create draft: ${response.message}`);
    },
  });

  const handleClick = () => {
    createSchemaVersionMutation.mutate({
      project_id: currentProjectId || "", 
      schema_id: currentSchemaId || ""
    });
  };

  const isLoading = createSchemaVersionMutation.isPending;

  return (
    <>
      <ShadowButton
        value={isLoading ? "" : "New Version"}
        leftIcon={isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <FilePlus size={20} />}
        className="xs:flex hidden"
        variant="accent-solid"
        onClick={handleClick}
        disabled={isLoading || !currentProjectId || !currentSchemaId}
      />
      <ShadowButton
        leftIcon={isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <FilePlus size={16} />}
        className="xs:hidden flex"
        variant="accent-solid"
        onClick={handleClick}
        disabled={isLoading || !currentProjectId || !currentSchemaId}
      />
    </>
  )
}
