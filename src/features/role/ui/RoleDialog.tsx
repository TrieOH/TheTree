import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useStore } from "@tanstack/react-store";
import { roleStore } from "../store";
import { createRoleFn, deleteRoleFn, patchRoleFn } from "../api";
import { toast } from "sonner";
import { useCrudOperations } from "@/shared/lib/hooks/useCrudStore";
import type { FieldConfig } from "@/shared/ui/form/types";
import { getFieldError } from "@/shared/lib/utils";
import { type RoleCRUD, roleCRUDSchema } from "../model/types";
import { formOptions } from "@tanstack/react-form";
import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";
import CrudForm from "@/shared/ui/form/CrudForm";

interface PropsI {
  project_id: string;
}

export default function RoleDialog({ project_id }: PropsI) {
  const queryClient = useQueryClient();
  const { formData, mode } = useStore(roleStore, (state) => state);

  const createRoleMutation = useMutation({
    mutationFn: createRoleFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["roles", project_id] });
        queryClient.setQueryData(["roles", project_id, response.data.id], response.data);
      } else toast.error(`Failed to create role: ${response.message}`);
    },
  });

  const patchRoleMutation = useMutation({
    mutationFn: patchRoleFn,
    onSuccess: (response, data) => {
      if (response.success) {
        toast.success(response.message || "Updated role");
        queryClient.setQueryData(["roles", project_id, data.id], data);
        queryClient.invalidateQueries({ queryKey: ["roles", project_id] });
      } else toast.error(`Failed to update roles: ${response.message}`);
    },
  });

  const deleteRoleMutation = useMutation({
    mutationFn: deleteRoleFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message || "Role deleted sucessfuly!");
        queryClient.invalidateQueries({ queryKey: ["roles"] });
      } else toast.error(`Failed to delete role: ${response.message}`);
    },
  });

  const { handleSubmit } = useCrudOperations({
    store: roleStore,
    autoClose: true,
    onCreate: async (data) => {
      createRoleMutation.mutate({ ...data, project_id });
    },
    onUpdate: async (id, data) => {
      patchRoleMutation.mutate({ id, ...data } as RoleCRUD);
    },
    onDelete: async (id) => {
      deleteRoleMutation.mutate({project_id, id});
    }
  });

  const fields: FieldConfig[] = [
    {
      name: "description", 
      label: "Description", 
      placeholder: "My Awesome Description", 
      autoComplete: "description",
      errors: getFieldError(roleCRUDSchema.shape.description, "a")
    }
  ];
  if(mode !== "edit") {
    fields.push({
      name: "name", 
      label: "Name", 
      placeholder: "My Awesome Role", 
      autoComplete: "name",
      errors: getFieldError(roleCRUDSchema.shape.name, "a")
    })
    fields.reverse();
  }

  const roleOpts = formOptions({
    defaultValues: (mode === 'create' ? { id: "", project_id: "" } : formData) as RoleCRUD,
    validators: {
      onChange: roleCRUDSchema,
      onMount: roleCRUDSchema,
    }
  });

  return (
    <CrudDialog
      formId="role-form"
      store={roleStore}
      title="Role"
      onSubmit={() => handleSubmit(formData as RoleCRUD)}
    >
      <CrudForm
        formId="role-form"
        fields={fields}
        options={{
          defaultValues: roleOpts.defaultValues,
          validators: roleOpts.validators,
          onSubmit: async ({ value }) => handleSubmit(value)
        }}
      />
    </CrudDialog>
  );
}