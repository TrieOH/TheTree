import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";
import { projectStore } from "../store";
import CrudForm from "@/shared/ui/form/CrudForm";
import { formOptions } from "@tanstack/react-form";
import type { FieldConfig } from "@/shared/ui/form/types";
import { useCrudOperations } from "@/shared/lib/hooks/useCrudStore";
import { projectCRUDSchema, type ProjectCRUD } from "../model/types";
import { createProjectFn } from "../api";
import { toast } from "sonner";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export function ProjectDialog() {
  const queryClient = useQueryClient();

  const createProjectMutation = useMutation({
    mutationFn: createProjectFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["projects"] });
      } else toast.error(`Failed to create project: ${response.message}`);
    },
    onError: (error) => {
      toast.error(`Failed to create project: ${error.message}`);
    },
  });

  const { handleSubmit } = useCrudOperations({
    store: projectStore,
    autoClose: true,
    onCreate: async (data) => {
      createProjectMutation.mutate(data);
    },
  });
  const fields: FieldConfig[] = [
    {name: "project_name", label: "Project Name", placeholder: "My Awesome Project", autoComplete: "project_name"}
  ]
  const projectOpts = formOptions({
    defaultValues: { id: "", project_name: "" } as ProjectCRUD,
    validators: {
      onChange: projectCRUDSchema,
      onMount: projectCRUDSchema,
    }
  });
  
  return (
    <CrudDialog
      formId="project-form"
      store={projectStore}
      title="Project"
    >
      <CrudForm 
        formId="project-form"
        fields={fields}
        options={{
          defaultValues: projectOpts.defaultValues,
          validators: projectOpts.validators,
          onSubmit: async ({ value }) => handleSubmit(value)
        }}
      />
    </CrudDialog>
  );
}