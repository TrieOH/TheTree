import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useStore } from "@tanstack/react-store";
import { permissionStore } from "../store";
import { createPermissionFn } from "../api";
import { toast } from "sonner";
import { useCrudOperations } from "@/shared/lib/hooks/useCrudStore";
import type { FieldConfig } from "@/shared/ui/form/types";
import { getFieldError } from "@/shared/lib/utils";
import { type PermissionCRUD, permissionCRUDSchema } from "../model/types";
import { formOptions } from "@tanstack/react-form";
import CrudForm from "@/shared/ui/form/CrudForm";
import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";

interface PropsI {
  project_id: string;
}

export default function PermissionDialog({ project_id }: PropsI) {
  const queryClient = useQueryClient();
  const { formData, mode } = useStore(permissionStore, (state) => state);

  const createPermissionMutation = useMutation({
    mutationFn: createPermissionFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["permissions", project_id] });
        queryClient.setQueryData(["permissions", project_id, response.data.id], response.data);
      } else toast.error(`Failed to create permission: ${response.message}`);
    },
    onError: (error) => {
      toast.error(`Failed to create scope: ${error.message}`);
    },
  });

  const { handleSubmit } = useCrudOperations({
    store: permissionStore,
    autoClose: true,
    onCreate: async (data) => {
      createPermissionMutation.mutate({ ...data, project_id });
    },
  });

  const fields: FieldConfig[] = [
    {
      name: "object", 
      label: "Object", 
      placeholder: "Events", 
      autoComplete: "object",
      errors: getFieldError(permissionCRUDSchema.shape.object)
    },
    {
      name: "action", 
      label: "Action", 
      placeholder: "read", 
      autoComplete: "action",
      errors: getFieldError(permissionCRUDSchema.shape.action)
    }
  ];

  const permissionOpts = formOptions({
    defaultValues: (mode === 'create' ? { id: "", object: "", action: "", project_id: "" } : formData) as PermissionCRUD,
    validators: {
      onChange: permissionCRUDSchema,
      onMount: permissionCRUDSchema,
    }
  });

  return (
    <CrudDialog
      formId="permission-form"
      store={permissionStore}
      title="Permission"
      onSubmit={() => handleSubmit(formData as PermissionCRUD)}
    >
      <CrudForm 
        formId="permission-form"
        fields={fields}
        options={{
          defaultValues: permissionOpts.defaultValues,
          validators: permissionOpts.validators,
          onSubmit: async ({ value }) => handleSubmit(value)
        }}
      />
    </CrudDialog>
  );
}