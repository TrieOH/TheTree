import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";
import { scopeCRUDSchema, type ScopeCRUD } from "../model/types";
import CrudForm from "@/shared/ui/form/CrudForm";
import { scopeStore } from "../store";
import { useStore } from "@tanstack/react-store";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createScopeFn } from "../api";
import { toast } from "sonner";
import { useCrudOperations } from "@/shared/lib/hooks/useCrudStore";
import type { FieldConfig } from "@/shared/ui/form/types";
import { formOptions } from "@tanstack/react-form";

interface PropsI {
  project_id: string;
}

export default function ScopeDialog({ project_id }: PropsI) {
  const queryClient = useQueryClient();
  const { formData, mode } = useStore(scopeStore, (state) => state);

  const createScopeMutation = useMutation({
    mutationFn: createScopeFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["scopes"] });
      } else toast.error(`Failed to create scope: ${response.message}`);
    },
    onError: (error) => {
      toast.error(`Failed to create scope: ${error.message}`);
    },
  });

  const { handleSubmit } = useCrudOperations({
    store: scopeStore,
    autoClose: true,
    onCreate: async (data) => {
      createScopeMutation.mutate({ ...data, project_id });
    },
  });
  const fields: FieldConfig[] = [
    {
      name: "name", 
      label: "Name", 
      placeholder: "My New Scope", 
      autoComplete: "name"
    },
    {
      name: "external_id", 
      label: "External ID", 
      placeholder: "Event", 
      autoComplete: "external_id"
    }
  ]
  const scopeOpts = formOptions({
    defaultValues: (mode === 'create' ? { id: "", external_id: "", name: "", project_id: "" } : formData) as ScopeCRUD,
    validators: {
      onChange: scopeCRUDSchema,
      onMount: scopeCRUDSchema,
    }
  });

  return (
    <CrudDialog
      formId="scope-form"
      store={scopeStore}
      title="Scope"
      onSubmit={() => handleSubmit(formData as ScopeCRUD)}
    >
      <CrudForm 
        formId="scope-form"
        fields={fields}
        options={{
          defaultValues: scopeOpts.defaultValues,
          validators: scopeOpts.validators,
          onSubmit: async ({ value }) => handleSubmit(value)
        }}
      />
    </CrudDialog>
  )
}