import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useStore } from "@tanstack/react-store";
import { permissionStore } from "../store";
import { createPermissionFn, deletePermissionFn, patchPermissionMetaFn } from "../api";
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

interface PermissionFormValues extends Omit<PermissionCRUD, 'meta'> {
  icon?: string;
  color?: string;
  description?: string;
  status?: "active" | "restricted" | "beta" | "deprecated";
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
  });

  const patchPermissionMetaMutation = useMutation({
    mutationFn: patchPermissionMetaFn,
    onSuccess: (response, data) => {
      if (response.success) {
        toast.success(response.message || "Updated permission");
        queryClient.setQueryData(["permissions", project_id, data.id], data);
        queryClient.invalidateQueries({ queryKey: ["permissions", project_id] });
      } else toast.error(`Failed to update permissions: ${response.message}`);
    },
  });

  const deletePermissionMutation = useMutation({
    mutationFn: deletePermissionFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message || "Permission deleted sucessfuly!");
        queryClient.invalidateQueries({ queryKey: ["permissions"] });
      } else toast.error(`Failed to delete permission: ${response.message}`);
    },
  });

  const { handleSubmit } = useCrudOperations({
    store: permissionStore,
    autoClose: true,
    onCreate: async (data) => {
      const { icon, color, description, status, ...rest } = data as PermissionFormValues;
      const finalData: PermissionCRUD = {
        ...rest,
        meta: { icon, color, description, status }
      };
      createPermissionMutation.mutate({ ...finalData, project_id });
    },
    onUpdate: async (id, data) => {
      const { icon, color, description, status, ...rest } = data as PermissionFormValues;
      const finalData: PermissionCRUD = {
        ...rest,
        id,
        meta: { icon, color, description, status }
      };
      patchPermissionMetaMutation.mutate(finalData);
    },
    onDelete: async (id) => {
      deletePermissionMutation.mutate({ project_id, id });
    }
  });

  const allFields: FieldConfig[] = [
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
    },
    {
      name: "description",
      label: "Description",
      placeholder: "Brief explanation",
      type: "text"
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

  const permissionOpts = formOptions({
    defaultValues: (mode === 'create' 
      ? { id: "", object: "", action: "", project_id: "", icon: "Shield", color: "#6366f1", status: "active", description: "" } 
      : { ...formData, ...formData?.meta }) as PermissionFormValues,
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
          onSubmit: async ({ value }) => handleSubmit(value)
        }}
      />
    </CrudDialog>
  );
}
