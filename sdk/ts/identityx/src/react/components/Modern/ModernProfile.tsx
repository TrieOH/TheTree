import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import type { JsonSchemaProperty, ProfileData } from "../../../types/auth-types";
import type { JsonValue } from "../../../types/token-types";
import { useAuth } from "../../AuthProvider";
import { Button } from "./Shared/Button";

export interface ModernProfileProps {
  actorId?: string;
  projectId?: string;
  schema?: JsonSchemaProperty;
  onSuccess?: (profile: ProfileData) => void | Promise<void>;
}

export function ModernProfile({ actorId, projectId, schema: schemaOverride, onSuccess }: ModernProfileProps) {
  const { auth } = useAuth();
  const resolvedActorId = actorId ?? auth.profile()?.id;
  const [schema, setSchema] = useState<JsonSchemaProperty | null>(schemaOverride ?? null);
  const [values, setValues] = useState<ProfileData>({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    let cancelled = false;
    if (!resolvedActorId) {
      setLoading(false);
      return;
    }
    Promise.all([
      schemaOverride ? Promise.resolve(null) : auth.getProfileSchema(projectId),
      auth.getActorProfile(resolvedActorId, projectId),
    ]).then(([schemaResponse, profileResponse]) => {
      if (cancelled) return;
      if (schemaResponse?.success) setSchema(schemaResponse.data.schema);
      if (profileResponse.success) setValues(profileResponse.data.profile);
    }).finally(() => {
      if (!cancelled) setLoading(false);
    });
    return () => { cancelled = true; };
  }, [auth, projectId, resolvedActorId, schemaOverride]);

  const properties = useMemo(() => Object.entries(schema?.properties ?? {}), [schema]);
  const required = new Set(schema?.required ?? []);

  if (loading) return <p className="text-sm text-muted-foreground">Loading profile…</p>;
  if (!resolvedActorId) return <p className="text-sm text-destructive">No authenticated actor was found.</p>;
  if (!schema || properties.length === 0) return <p className="text-sm text-muted-foreground">No active profile fields are configured.</p>;

  const setValue = (name: string, value: JsonValue) =>
    setValues((current) => ({ ...current, [name]: value }));

  return (
    <form
      className="space-y-4"
      onSubmit={async (event) => {
        event.preventDefault();
        setSaving(true);
        try {
          const response = await auth.upsertActorProfile(resolvedActorId, { profile: values }, projectId);
          if (!response.success) throw new Error(response.message || "Could not save profile");
          setValues(response.data.profile);
          await onSuccess?.(response.data.profile);
          toast.success("Profile saved");
        } catch (error) {
          toast.error(error instanceof Error ? error.message : "Could not save profile");
        } finally {
          setSaving(false);
        }
      }}
    >
      {properties.map(([name, field]) => (
        <label key={name} className="block space-y-1.5 text-sm">
          <span className="font-medium">{field.title ?? name}{required.has(name) ? " *" : ""}</span>
          {field.description && <span className="block text-xs text-muted-foreground">{field.description}</span>}
          {field.type === "boolean" ? (
            <input type="checkbox" checked={values[name] === true} onChange={(e) => setValue(name, e.target.checked)} />
          ) : field.enum ? (
            <select
              required={required.has(name)}
              value={String(values[name] ?? "")}
              onChange={(e) => setValue(name, e.target.value)}
              className="w-full rounded-md border border-border bg-background px-3 py-2"
            >
              <option value="">Select…</option>
              {field.enum.map((option) => <option key={String(option)} value={String(option)}>{String(option)}</option>)}
            </select>
          ) : (
            <input
              required={required.has(name)}
              type={field.type === "number" || field.type === "integer" ? "number" : field.format === "email" ? "email" : field.format === "uri" ? "url" : "text"}
              value={typeof values[name] === "string" || typeof values[name] === "number" ? values[name] : ""}
              onChange={(e) => setValue(name, field.type === "number" || field.type === "integer" ? Number(e.target.value) : e.target.value)}
              className="w-full rounded-md border border-border bg-background px-3 py-2"
            />
          )}
        </label>
      ))}
      <Button type="submit" disabled={saving}>{saving ? "Saving…" : "Save profile"}</Button>
    </form>
  );
}
