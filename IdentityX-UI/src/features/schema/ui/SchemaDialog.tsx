import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";
import CrudForm from "@/shared/ui/form/CrudForm";
import { formOptions } from "@tanstack/react-form";
import type { FieldConfig } from "@/shared/ui/form/types";
import { useCrudOperations } from "@/shared/lib/hooks/useCrudStore";
import { toast } from "sonner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useStore } from "@tanstack/react-store";
import { schemaStore } from "../store";
import { createSchemaFn } from "../api";
import  { type SchemaCRUD, schemaCRUDSchema } from "../model/types";
import { getFieldError } from "@/shared/lib/utils";

interface PropsI {
  project_id: string;
}

export function SchemaDialog({ project_id }: PropsI) {
  const queryClient = useQueryClient();
  const { formData, mode } = useStore(schemaStore, (state) => state);

  const createSchemaMutation = useMutation({
    mutationFn: createSchemaFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["schemas"] });
      } else toast.error(`Failed to create schema: ${response.message}`);
    },
    onError: (error) => {
      toast.error(`Failed to create schema: ${error.message}`);
    },
  });

  const { handleSubmit } = useCrudOperations({
    store: schemaStore,
    autoClose: true,
    onCreate: async (data) => {
      createSchemaMutation.mutate({ ...data, project_id });
    },
  });
  const fields: FieldConfig[] = [
    {
      name: "title", 
      label: "Title", 
      placeholder: "My New Schema", 
      autoComplete: "title",
      errors: getFieldError(schemaCRUDSchema.shape["title"])
    },
    {
      name: "flow_id", 
      label: "Flow ID", 
      placeholder: "auth flow", 
      autoComplete: "flow_id",
      errors: getFieldError(schemaCRUDSchema.shape["flow_id"])
    }
  ]
  const projectOpts = formOptions({
    defaultValues: (mode === 'create' ? { id: "", flow_id: "", title: "", project_id: "" } : formData) as SchemaCRUD,
    validators: {
      onChange: schemaCRUDSchema,
      onMount: schemaCRUDSchema,
    }
  });
  
  return (
    <CrudDialog
      formId="schema-form"
      store={schemaStore}
      title="Schema"
      onSubmit={() => handleSubmit(formData as SchemaCRUD)}
    >
      <CrudForm 
        formId="schema-form"
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