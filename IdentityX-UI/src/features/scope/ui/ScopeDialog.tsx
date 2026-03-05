import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";
import { scopeCRUDSchema, type ScopeCRUD } from "../model/types";
import CrudForm from "@/shared/ui/form/CrudForm";
import { scopeStore } from "../store";
import { useStore } from "@tanstack/react-store";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createScopeFn, deleteScopeFn, patchScopeMetaFn, scopesQueryOptions } from "../api";
import { toast } from "sonner";
import { useCrudOperations } from "@/shared/lib/hooks/useCrudStore";
import type { FieldConfig } from "@/shared/ui/form/types";
import { formOptions } from "@tanstack/react-form";
import { getFieldError } from "@/shared/lib/utils";

interface PropsI {
  project_id: string;
}

// Flat type for form usage
interface ScopeFormValues extends Omit<ScopeCRUD, 'meta'> {
  icon?: string;
  color?: string;
  folder?: string;
  description?: string;
  status?: "active" | "restricted" | "beta" | "deprecated";
}

export default function ScopeDialog({ project_id }: PropsI) {
  const queryClient = useQueryClient();
  const { formData, mode } = useStore(scopeStore, (state) => state);
  const { data: scopes } = useQuery(scopesQueryOptions(project_id));

  const createScopeMutation = useMutation({
    mutationFn: createScopeFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries(scopesQueryOptions(project_id));
        queryClient.setQueryData(["scopes", project_id, response.data.id], response.data);
      } else toast.error(`Failed to create scope: ${response.message}`);
    },
    onError: (error) => {
      toast.error(`Failed to create scope: ${error.message}`);
    },
  });

  const patchScopeMetaMutation = useMutation({
    mutationFn: patchScopeMetaFn,
    onSuccess: (response, data) => {
      if (response.success) {
        toast.success(response.message || "Updated scope");
        queryClient.setQueryData(["scopes", project_id, data.id], data);
        queryClient.invalidateQueries(scopesQueryOptions(project_id));
      } else toast.error(`Failed to update scopes: ${response.message}`);
    },
  });

  const deleteScopeMutation = useMutation({
    mutationFn: deleteScopeFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message || "Scope deleted sucessfuly!");
        queryClient.invalidateQueries(scopesQueryOptions(project_id));
      } else toast.error(`Failed to delete scope: ${response.message}`);
    },
  });

  const { handleSubmit } = useCrudOperations({
    store: scopeStore,
    autoClose: true,
    onCreate: async (data) => {
      const { icon, color, folder, description, status, ...rest } = data as ScopeFormValues;
      const finalData: ScopeCRUD = {
        ...rest,
        parent_id: (rest.parent_id === "" || rest.parent_id === "none") ? null : rest.parent_id,
        meta: { icon, color, folder, description, status }
      };
      createScopeMutation.mutate({ ...finalData, project_id });
    },
    onUpdate: async (id, data) => {
      const { icon, color, folder, description, status, ...rest } = data as ScopeFormValues;
      const finalData: Partial<ScopeCRUD> = {
        ...rest,
        id,
        parent_id: (rest.parent_id === "" || rest.parent_id === "none") ? null : rest.parent_id,
        meta: { icon, color, folder, description, status }
      };
      patchScopeMetaMutation.mutate(finalData);
    },
    onDelete: async (id) => {
      deleteScopeMutation.mutate({ project_id, id });
    }
  });

  const allFields: FieldConfig[] = [
    {
      name: "name", 
      label: "Name", 
      placeholder: "My New Scope", 
      autoComplete: "name",
      errors: getFieldError(scopeCRUDSchema.shape.name)
    },
    {
      name: "external_id", 
      label: "External ID", 
      placeholder: "Event", 
      autoComplete: "external_id",
      errors: getFieldError(scopeCRUDSchema.shape.external_id)
    },
    {
      name: "folder",
      label: "Folder",
      placeholder: "e.g. Activity or Settings",
      type: "text"
    },
    {
      name: "parent_id",
      label: "Parent Scope",
      type: "select",
      options: [
        { label: "No Parent (Root)", value: "none" },
        ...(scopes?.filter(s => s.id !== formData?.id).map(s => ({ label: s.name, value: s.id })) || [])
      ]
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
  ]

  const fields = mode === 'edit' 
    ? allFields.filter(f => ["description", "status", "icon", "color", "folder"].includes(f.name))
    : allFields;

  const scopeOpts = formOptions({
    defaultValues: (mode === 'create' 
      ? { id: "", external_id: "", name: "", project_id, parent_id: "none", icon: "Shield", color: "#6366f1", status: "active", description: "", folder: "", ...formData } 
      : { ...formData, ...formData?.meta, parent_id: formData?.parent_id || "none" }) as ScopeFormValues,
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
