import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";
import { projectStore } from "../store";
import z from "zod";
import CrudForm from "@/shared/ui/form/CrudForm";
import { formOptions } from "@tanstack/react-form";
import type { FieldConfig } from "@/shared/ui/form/types";

export function ProjectDialog() {
  const fields: FieldConfig[] = [
    {name: "project_name", label: "Project Name", placeholder: "My Awesome Project", autoComplete: "project_name"}
  ]
  const projectOpts = formOptions({
    defaultValues: {
      project_name: "",
    },
    validators: {
      onChange: z.object({
        project_name: z.string().min(3, "Project name must be at least 3 characters long"),
      }),
      onMount: z.object({
        project_name: z.string().min(3, "Project name must be at least 3 characters long"),
      })
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
          onSubmit: async ({ value, formApi }) => {
            formApi.reset()
            console.log(value);
          },
        }}
      />
    </CrudDialog>
  );
}