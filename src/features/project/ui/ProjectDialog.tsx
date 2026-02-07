import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";
import { projectStore } from "../store";
import CrudForm from "@/shared/ui/form/CrudForm";
import { formOptions } from "@tanstack/react-form";
import type { FieldConfig } from "@/shared/ui/form/types";
import { useCrudOperations } from "@/shared/lib/hooks/useCrudStore";
import { projectCRUDSchema, type ProjectCRUD } from "../model/types";
// import { useAuth } from "@trieoh/node-auth-sdk/react";

export function ProjectDialog() {
  // useAuth
  const { handleSubmit } = useCrudOperations({
    store: projectStore,
    autoClose: true,
    onCreate: async (data) => {
      console.log(data)
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