import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useStore } from "@tanstack/react-store";
import { roleStore } from "../store";
import { createRoleFn, deleteRoleFn, patchRoleFn, patchRoleMetaFn, roleQueryOptions } from "../api";
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

interface RoleFormValues extends Omit<RoleCRUD, 'meta'> {
  icon?: string;
  color?: string;
  status?: "active" | "restricted" | "beta" | "deprecated";
}

export default function RoleDialog({ project_id }: PropsI) {
  const queryClient = useQueryClient();
  const { formData, mode } = useStore(roleStore, (state) => state);

  const createRoleMutation = useMutation({
    mutationFn: createRoleFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries(roleQueryOptions(project_id));
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
        queryClient.invalidateQueries(roleQueryOptions(project_id));
      } else toast.error(`Failed to update roles: ${response.message}`);
    },
  });

  const patchRoleMetaMutation = useMutation({
    mutationFn: patchRoleMetaFn,
    onSuccess: (response, data) => {
      if (response.success) {
        toast.success(response.message || "Updated role metadata");
        queryClient.setQueryData(["roles", project_id, data.id], data);
        queryClient.invalidateQueries(roleQueryOptions(project_id));
      } else toast.error(`Failed to update role metadata: ${response.message}`);
    },
  });

  const deleteRoleMutation = useMutation({
    mutationFn: deleteRoleFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message || "Role deleted sucessfuly!");
        queryClient.invalidateQueries(roleQueryOptions(project_id));
      } else toast.error(`Failed to delete role: ${response.message}`);
    },
  });

  const { handleSubmit } = useCrudOperations({
    store: roleStore,
    autoClose: true,
    onCreate: async (data) => {
      const { icon, color, status, ...rest } = data as RoleFormValues;
      const finalData: RoleCRUD = {
        ...rest,
        meta: { icon, color, status, description: rest.description }
      };
      createRoleMutation.mutate({ ...finalData, project_id });
    },
    onUpdate: async (id, data) => {
      const { icon, color, status, ...rest } = data as RoleFormValues;
      const finalData: RoleCRUD = {
        ...rest,
        id,
        meta: { icon, color, status, description: rest.description }
      };

      await patchRoleMutation.mutateAsync(finalData);
      patchRoleMetaMutation.mutate(finalData);
    },
    onDelete: async (id) => {
      deleteRoleMutation.mutate({project_id, id});
    }
  });

  const allFields: FieldConfig[] = [
    {
      name: "name", 
      label: "Name", 
      placeholder: "My Awesome Role", 
      autoComplete: "name",
      errors: getFieldError(roleCRUDSchema.shape.name, "a")
    },
    {
      name: "description", 
      label: "Description", 
      placeholder: "My Awesome Description", 
      autoComplete: "description",
      errors: getFieldError(roleCRUDSchema.shape.description, "a")
    },
    {
      name: "status",
      label: "Status",
      type: "select",
      options: [
        { label: "Active", value: "active" },
        { label: "Restricted", value: "restricted" },
        { label: "Beta", value: "beta" },
        { label: "Deprecated", value: "deprecated" }
      ]
    },
    {
      name: "icon",
      label: "Select Icon",
      type: "icon"
    },
    {
      name: "color",
      label: "Select Color / Gradient",
      type: "color"
    }
  ];

  const fields = mode === 'edit' 
    ? allFields.filter(f => ["description", "status", "icon", "color"].includes(f.name))
    : allFields;

  const roleOpts = formOptions({
    defaultValues: (mode === 'create' 
      ? { id: "", project_id: "", name: "", description: "", icon: "Shield", color: "#6366f1", status: "active" } 
      : { ...formData, ...formData?.meta }) as RoleFormValues,
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