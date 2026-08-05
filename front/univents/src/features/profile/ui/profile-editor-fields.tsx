import type { ProfileData } from "@trieoh/identityx-sdk-ts";
import { CircleAlert } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import { Label } from "@/shared/ui/shadcn/label";
import type { ProfileSchemaNode } from "../model/profile-data";
import { SYSTEM_PROFILE_FIELDS } from "../model/profile-data";
import { SchemaInput } from "./profile-editor-inputs";
import { fieldLabel, readValue, writeValue } from "./profile-editor-utils";

export function EditorCard({
  title,
  icon,
  compact = false,
  children,
}: {
  title: string;
  icon?: React.ReactNode;
  compact?: boolean;
  children: React.ReactNode;
}) {
  return (
    <section
      className={cn(
        "rounded-md border border-border bg-card shadow-md shadow-foreground/5",
        compact ? "p-3" : "p-5",
      )}
    >
      <h2
        className={cn(
          "flex items-center gap-2 font-semibold",
          compact ? "mb-2 text-sm" : "mb-4 text-base",
        )}
      >
        {icon}
        {title}
      </h2>
      <div className="space-y-4">{children}</div>
    </section>
  );
}

export function EditorObject({
  schema,
  values,
  onChange,
  name,
}: {
  schema: ProfileSchemaNode;
  values: ProfileData;
  onChange: (values: ProfileData) => void;
  name: string;
}) {
  const field = schema.properties?.[name];
  if (!field)
    return (
      <p className="text-sm text-muted-foreground">
        Campo indisponível no schema.
      </p>
    );
  return (
    <SchemaFields
      schema={field}
      values={values}
      onChange={onChange}
      path={[name]}
    />
  );
}

export function EditorField({
  schema,
  values,
  onChange,
  path,
  compact = false,
  hideLabel = false,
}: {
  schema: ProfileSchemaNode;
  values: ProfileData;
  onChange: (values: ProfileData) => void;
  path: string[];
  compact?: boolean;
  hideLabel?: boolean;
}) {
  const field = findSchemaField(schema, path);
  if (!field) return null;
  const id = `profile-${path.join("-")}`;
  const value = readValue(values, path);
  const type = Array.isArray(field.type)
    ? field.type.find((item) => item !== "null")
    : field.type;
  return (
    <div className={cn("space-y-2", compact && "space-y-1.5")}>
      {!hideLabel && (
        <Label htmlFor={id} className="text-sm font-medium">
          {fieldLabel(path.at(-1) ?? "", field)}
        </Label>
      )}
      <SchemaInput
        id={id}
        field={field}
        type={type}
        value={value}
        required={
          path.length === 1 && Boolean(schema.required?.includes(path[0] ?? ""))
        }
        onChange={(next) => onChange(writeValue(values, path, next))}
      />
    </div>
  );
}

export function findSchemaField(schema: ProfileSchemaNode, path: string[]) {
  let current: ProfileSchemaNode | undefined = schema;
  for (const segment of path) current = current?.properties?.[segment];
  return current;
}

export function ProfileWarning({ message }: { message: string }) {
  return (
    <div
      role="alert"
      className="mx-auto flex max-w-xl gap-3 rounded-2xl border border-amber-500/30 bg-amber-500/10 p-5 text-amber-950 dark:text-amber-100"
    >
      <CircleAlert className="mt-0.5 size-5 shrink-0" />
      <div>
        <p className="font-medium">Não foi possível abrir a edição</p>
        <p className="mt-1 text-sm opacity-80">{message}</p>
      </div>
    </div>
  );
}

export function SchemaFields({
  schema,
  values,
  onChange,
  path = [],
}: {
  schema: ProfileSchemaNode;
  values: ProfileData;
  onChange: (values: ProfileData) => void;
  path?: string[];
}) {
  const required = new Set(schema.required ?? []);
  return (
    <div className={path.length ? "grid gap-3" : "space-y-6"}>
      {Object.entries(schema.properties ?? {}).map(([name, field]) => {
        if (SYSTEM_PROFILE_FIELDS.has(name)) return null;
        const fieldPath = [...path, name];
        const value = readValue(values, fieldPath);
        const fieldType = Array.isArray(field.type)
          ? field.type.find((type) => type !== "null")
          : field.type;
        if (fieldType === "object") {
          return (
            <fieldset
              key={fieldPath.join(".")}
              className="space-y-4 rounded-xl border p-4 md:col-span-2"
            >
              <legend className="px-2 font-semibold">
                {fieldLabel(name, field)}
              </legend>
              <SchemaFields
                schema={field}
                values={values}
                onChange={onChange}
                path={fieldPath}
              />
            </fieldset>
          );
        }
        const id = fieldPath.join("-");
        const isRequired = required.has(name);
        return (
          <div
            key={id}
            className={cn(
              fieldType === "boolean"
                ? "flex min-w-0 items-center justify-between gap-4 rounded-md border border-border bg-background p-3"
                : path[0] === "socials"
                  ? "flex min-w-0 items-center gap-2"
                  : "space-y-2",
            )}
          >
            {path[0] !== "socials" && (
              <div>
                <Label htmlFor={id}>
                  {fieldLabel(name, field)}
                  {isRequired ? " *" : ""}
                </Label>
              </div>
            )}
            <SchemaInput
              id={id}
              field={field}
              type={fieldType}
              value={value}
              required={isRequired}
              onChange={(next) => onChange(writeValue(values, fieldPath, next))}
            />
          </div>
        );
      })}
    </div>
  );
}
