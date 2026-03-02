import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";
import { scopeCRUDSchema, type ScopeCRUD } from "../model/types";
import CrudForm from "@/shared/ui/form/CrudForm";
import { scopeStore } from "../store";
import { useStore } from "@tanstack/react-store";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createScopeFn, patchScopeMetaFn } from "../api";
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
  description?: string;
  status?: "active" | "restricted" | "beta" | "deprecated";
}

export default function ScopeDialog({ project_id }: PropsI) {
  const queryClient = useQueryClient();
  const { formData, mode } = useStore(scopeStore, (state) => state);

  const createScopeMutation = useMutation({
    mutationFn: createScopeFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["scopes", project_id] });
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
        queryClient.invalidateQueries({ queryKey: ["scopes", project_id] });
      } else toast.error(`Failed to update scopes: ${response.message}`);
    },
  });

  const { handleSubmit } = useCrudOperations({
    store: scopeStore,
    autoClose: true,
    onCreate: async (data) => {
      const { icon, color, description, status, ...rest } = data as ScopeFormValues;
      const finalData: ScopeCRUD = {
        ...rest,
        meta: { icon, color, description, status }
      };
      createScopeMutation.mutate({ ...finalData, project_id });
    },
    onUpdate: async (id, data) => {
      const { icon, color, description, status, ...rest } = data as ScopeFormValues;
      const finalData: Partial<ScopeCRUD> = {
        ...rest,
        id,
        meta: { icon, color, description, status }
      };
      patchScopeMetaMutation.mutate(finalData);
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
      errors: getFieldError(scopeCRUDSchema.shape.external_id, "d")
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
    ? allFields.filter(f => ["description", "status", "icon", "color"].includes(f.name))
    : allFields;

  const scopeOpts = formOptions({
    defaultValues: (mode === 'create' 
      ? { id: "", external_id: "", name: "", project_id: "", icon: "Shield", color: "#6366f1", status: "active", description: "" } 
      : { ...formData, ...formData?.meta }) as ScopeFormValues,
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
          onSubmit: async ({ value }) => handleSubmit(value)
        }}
      />
    </CrudDialog>
  )
}
