import type { JsonValue, ProfileData } from "@trieoh/identityx-sdk-ts";
import { ArrowLeft, CircleAlert, Globe, Save } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { preprocessImageUpload } from "@/features/storage/api";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";
import { Skeleton } from "@/shared/ui/shadcn/skeleton";
import blueskyIcon from "@/shared/ui/social-icons/assets/bluesky.svg";
import discordIcon from "@/shared/ui/social-icons/assets/discord.svg";
import githubIcon from "@/shared/ui/social-icons/assets/github.svg";
import instagramIcon from "@/shared/ui/social-icons/assets/instagram.svg";
import linkedinIcon from "@/shared/ui/social-icons/assets/linkedin.svg";
import twitchIcon from "@/shared/ui/social-icons/assets/twitch.svg";
import xIcon from "@/shared/ui/social-icons/assets/x.svg";
import youtubeIcon from "@/shared/ui/social-icons/assets/youtube.svg";
import type { ProfileSchemaNode } from "../model/profile-data";
import {
  applyProfileSchemaDefaults,
  asUniventsProfile,
  SYSTEM_PROFILE_FIELDS,
  socialHref,
  withProfileTimestamps,
} from "../model/profile-data";
import { ProfileImageInput } from "./profile-image-input";
import { TimezoneCombobox } from "./timezone-combobox";

type ProfileImageField = "pfpUrl" | "bannerUrl";
type PendingImages = Partial<Record<ProfileImageField, File>>;

export interface ProfileEditorProps {
  load: () => Promise<{
    schema: {
      success: boolean;
      data?: { schema: ProfileSchemaNode };
      message?: string;
    };
    profile: {
      success: boolean;
      data?: { handle?: string | null; profile: ProfileData };
      message?: string;
    };
  }>;
  save: (
    profile: ProfileData,
    handle?: string,
  ) => Promise<{
    success: boolean;
    data?: { profile: ProfileData };
    message?: string;
  }>;
  onCancel: () => void;
  onSaved: (profile: ProfileData) => void;
}

export function ProfileEditor({
  load,
  save,
  onCancel,
  onSaved,
}: ProfileEditorProps) {
  const [schema, setSchema] = useState<ProfileSchemaNode>();
  const [original, setOriginal] = useState<ProfileData>({});
  const [values, setValues] = useState<ProfileData>({});
  const [handle, setHandle] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string>();
  const [pendingImages, setPendingImages] = useState<PendingImages>({});

  useEffect(() => {
    let active = true;
    setError(undefined);
    load()
      .then(({ schema: schemaResult, profile }) => {
        if (!active) return;
        if (!schemaResult.success || !schemaResult.data) {
          setError(
            schemaResult.message ||
              "Não foi possível carregar o schema do perfil.",
          );
        } else {
          setSchema(schemaResult.data.schema);
        }
        if (profile.success && profile.data) {
          setHandle(profile.data.handle ?? "");
          setOriginal(profile.data.profile);
          setValues(
            schemaResult.data
              ? applyProfileSchemaDefaults(
                  schemaResult.data.schema,
                  profile.data.profile,
                )
              : profile.data.profile,
          );
        } else if (schemaResult.data) {
          setValues(applyProfileSchemaDefaults(schemaResult.data.schema, {}));
        }
      })
      .catch((cause: unknown) => {
        if (!active) return;
        setError(
          cause instanceof Error
            ? cause.message
            : "Não foi possível carregar o perfil.",
        );
      })
      .finally(() => active && setLoading(false));
    return () => {
      active = false;
    };
  }, [load]);

  if (loading)
    return (
      <main className="min-h-dvh bg-background pb-32">
        <Skeleton className="h-40 w-full sm:h-48 md:h-56" />
        <div className="mx-auto max-w-7xl px-4">
          <Skeleton className="-mt-12 size-24 rounded-full border-4 border-background md:-mt-16 md:size-32" />
          <div className="mt-4 grid gap-4 md:grid-cols-2">
            <Skeleton className="h-28 rounded-md" />
            <Skeleton className="h-28 rounded-md" />
          </div>
        </div>
        <div className="mx-auto mt-5 grid max-w-7xl gap-5 px-4 md:grid-cols-[minmax(0,1fr)_360px]">
          <Skeleton className="h-96 rounded-md" />
          <Skeleton className="h-64 rounded-md" />
        </div>
      </main>
    );
  if (error || !schema)
    return (
      <main className="mx-auto max-w-3xl px-4 py-16 md:px-6">
        <ProfileWarning message={error || "Schema de perfil indisponível."} />
        <div className="mt-4 flex justify-center">
          <Button type="button" variant="outline" onClick={onCancel}>
            <ArrowLeft className="size-4" />
            Voltar ao perfil
          </Button>
        </div>
      </main>
    );

  return (
    <main className="min-h-dvh bg-background pb-32">
      <form
        onSubmit={async (event) => {
          event.preventDefault();
          setSaving(true);
          try {
            const imageEntries = Object.entries(pendingImages) as [
              ProfileImageField,
              File,
            ][];
            const uploadResults = await Promise.allSettled(
              imageEntries.map(async ([field, file]) => ({
                field,
                url: await preprocessImageUpload(
                  file,
                  "profiles/images",
                  crypto.randomUUID(),
                ),
              })),
            );
            const nextValues = normalizeProfileLinks(values);
            const failedImages: string[] = [];
            uploadResults.forEach((result, index) => {
              const [field] = imageEntries[index] ?? [];
              if (!field) return;
              if (result.status === "fulfilled") {
                nextValues[result.value.field] = result.value.url;
              } else {
                const label = field === "pfpUrl" ? "foto do perfil" : "banner";
                const reason =
                  result.reason instanceof Error
                    ? result.reason.message
                    : "falha no upload ou na moderação";
                failedImages.push(`${label} (${reason})`);
              }
            });
            const response = await save(
              projectProfileValues(
                schema,
                withProfileTimestamps(original, nextValues),
              ),
              handle.trim() || undefined,
            );
            if (!response.success)
              throw new Error(
                response.message || "Não foi possível salvar o perfil.",
              );
            toast.success("Perfil atualizado");
            if (failedImages.length > 0) {
              toast.warning(
                `As outras informações foram salvas, mas houve erro em: ${failedImages.join("; ")}.`,
              );
            }
            onSaved(nextValues);
          } catch (cause) {
            toast.error(
              cause instanceof Error
                ? cause.message
                : "Não foi possível salvar o perfil.",
            );
          } finally {
            setSaving(false);
          }
        }}
      >
        <InlineProfileEditor
          schema={schema}
          values={values}
          handle={handle}
          onHandleChange={setHandle}
          onChange={setValues}
          onCancel={onCancel}
          pendingImages={pendingImages}
          onImageSelect={(field, file) =>
            setPendingImages((current) => ({ ...current, [field]: file }))
          }
        />
        <div className="fixed inset-x-0 bottom-0 z-60 border-t border-border bg-background/95 p-3 shadow-[0_-8px_30px_rgb(0_0_0/0.08)] backdrop-blur">
          <div className="mx-auto flex max-w-7xl items-center justify-between gap-3 px-1 md:px-4">
            <Button
              type="button"
              variant="ghost"
              className="h-8 rounded-sm"
              disabled={saving}
              onClick={onCancel}
            >
              <ArrowLeft className="size-4" />
              Cancelar
            </Button>
            <p className="hidden text-sm text-muted-foreground sm:block">
              Revise a prévia antes de salvar.
            </p>
            <Button type="submit" disabled={saving} className="h-8 rounded-sm">
              <Save className="size-4" />
              {saving ? "Salvando…" : "Salvar alterações"}
            </Button>
          </div>
        </div>
      </form>
    </main>
  );
}

function InlineProfileEditor({
  schema,
  values,
  handle,
  onHandleChange,
  onChange,
  onCancel,
  pendingImages,
  onImageSelect,
}: {
  schema: ProfileSchemaNode;
  values: ProfileData;
  handle: string;
  onHandleChange: (handle: string) => void;
  onChange: (values: ProfileData) => void;
  onCancel: () => void;
  pendingImages: PendingImages;
  onImageSelect: (field: ProfileImageField, file: File) => void;
}) {
  const profile = asUniventsProfile(values);
  return (
    <>
      <section className="relative border-b border-border bg-card shadow-md">
        <ProfileImageInput
          label="Banner"
          currentUrl={profile.bannerUrl}
          file={pendingImages.bannerUrl}
          onSelect={(file) => onImageSelect("bannerUrl", file)}
          variant="banner"
          className={cn(
            "relative h-40 bg-background bg-linear-to-br sm:h-48 md:h-56",
            "from-primary/40 via-primary/15 to-muted",
          )}
        />
        <div className="pointer-events-none absolute inset-x-4 top-4 z-20 mx-auto max-w-7xl">
          <div className="pointer-events-auto w-fit">
            <Button
              type="button"
              variant="secondary"
              size="icon"
              className="rounded-full shadow-lg"
              onClick={onCancel}
              aria-label="Voltar ao perfil"
            >
              <ArrowLeft className="size-4" />
            </Button>
          </div>
        </div>
        <div className="mx-auto max-w-7xl px-4">
          <div className="relative pb-5 md:pb-6">
            <ProfileImageInput
              label="Foto do perfil"
              currentUrl={profile.pfpUrl}
              file={pendingImages.pfpUrl}
              onSelect={(file) => onImageSelect("pfpUrl", file)}
              variant="avatar"
              className="-mt-12 md:-mt-16"
            />
            <div className="mt-3">
              <div className="grid items-end gap-3 sm:grid-cols-2 lg:grid-cols-3">
                <div className="min-w-0 space-y-1.5">
                  <Label htmlFor="profile-handle">Nome de usuário</Label>
                  <Input
                    id="profile-handle"
                    value={handle}
                    placeholder="seu nickname"
                    pattern={"[^\\s/]+"}
                    title="Use um nickname sem espaços ou barras"
                    autoCapitalize="none"
                    autoComplete="username"
                    spellCheck={false}
                    className="h-10 bg-background"
                    onChange={(event) => onHandleChange(event.target.value)}
                  />
                </div>
                {[
                  "preferredName",
                  "legalName",
                  "role",
                  "organization",
                  "pronouns",
                  "tagline",
                ].map((name) => (
                  <EditorField
                    key={name}
                    schema={schema}
                    values={values}
                    onChange={onChange}
                    path={[name]}
                  />
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      <div className="mx-auto mt-4 grid max-w-7xl gap-4 px-4 md:mt-5 md:grid-cols-[minmax(0,1fr)_360px] md:gap-5">
        <div className="space-y-5">
          <EditorCard title="Sobre mim">
            <EditorField
              schema={schema}
              values={values}
              onChange={onChange}
              path={["aboutMe"]}
              compact
              hideLabel
            />
          </EditorCard>
          <EditorCard title="Idiomas">
            <EditorField
              schema={schema}
              values={values}
              onChange={onChange}
              path={[
                schema.properties?.languages ? "languages" : "specializations",
              ]}
              compact
              hideLabel
            />
          </EditorCard>
          <EditorCard title="Contato" icon={<Globe className="size-4" />}>
            <div className="grid gap-4 sm:grid-cols-2">
              <EditorField
                schema={schema}
                values={values}
                onChange={onChange}
                path={["website"]}
              />
              <EditorField
                schema={schema}
                values={values}
                onChange={onChange}
                path={["contactEmail"]}
              />
            </div>
            <EditorObject
              schema={schema}
              values={values}
              onChange={onChange}
              name="socials"
            />
          </EditorCard>
        </div>

        <aside className="space-y-5">
          <EditorCard title="Privacidade">
            <EditorObject
              schema={schema}
              values={values}
              onChange={onChange}
              name="visibility"
            />
          </EditorCard>
          <EditorCard title="Outros detalhes">
            <EditorField
              schema={schema}
              values={values}
              onChange={onChange}
              path={["timezone"]}
              compact
            />
          </EditorCard>
        </aside>
      </div>
    </>
  );
}

function EditorCard({
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

function EditorObject({
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

function EditorField({
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

function findSchemaField(schema: ProfileSchemaNode, path: string[]) {
  let current: ProfileSchemaNode | undefined = schema;
  for (const segment of path) current = current?.properties?.[segment];
  return current;
}

function ProfileWarning({ message }: { message: string }) {
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

function SchemaFields({
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

function SchemaInput({
  id,
  field,
  type,
  value,
  required,
  onChange,
}: {
  id: string;
  field: ProfileSchemaNode;
  type?: string;
  value: JsonValue | undefined;
  required: boolean;
  onChange: (value: JsonValue) => void;
}) {
  if (type === "boolean")
    return (
      <label className="relative inline-flex shrink-0 cursor-pointer items-center">
        <input
          id={id}
          type="checkbox"
          checked={value === true}
          onChange={(event) => onChange(event.target.checked)}
          className="peer sr-only"
        />
        <span className="h-6 w-11 rounded-full bg-input transition-colors peer-checked:bg-primary peer-focus-visible:ring-2 peer-focus-visible:ring-ring peer-focus-visible:ring-offset-2 after:absolute after:left-0.5 after:top-0.5 after:size-5 after:rounded-full after:bg-background after:shadow-sm after:transition-transform peer-checked:after:translate-x-5" />
      </label>
    );
  if (type === "array")
    return (
      <ArrayDraftInput
        id={id}
        required={required}
        value={Array.isArray(value) ? value : []}
        placeholder={fieldPlaceholder(nameFromId(id))}
        onChange={onChange}
      />
    );
  if (nameFromId(id) === "pronouns")
    return <PronounsInput id={id} value={value} onChange={onChange} />;
  if (nameFromId(id) === "timezone") {
    const selectedTimezone = typeof value === "string" ? value : "";
    const timezoneOptions = Array.from(
      new Set([selectedTimezone, ...TIMEZONES].filter(Boolean)),
    ).map((timezone) => ({
      value: timezone,
      label: timezoneLabel(timezone),
    }));
    return (
      <TimezoneCombobox
        id={id}
        value={selectedTimezone}
        options={timezoneOptions}
        onChange={onChange}
      />
    );
  }
  if (field.enum)
    return (
      <select
        id={id}
        required={required}
        value={
          typeof value === "string" || typeof value === "number"
            ? String(value)
            : ""
        }
        onChange={(event) => onChange(event.target.value)}
        className="h-10 w-full rounded-md border border-input bg-background px-3 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        <option value="">Selecione…</option>
        {field.enum.map((option) => (
          <option key={String(option)} value={String(option)}>
            {String(option)}
          </option>
        ))}
      </select>
    );
  if (nameFromId(id) === "aboutMe")
    return (
      <textarea
        id={id}
        required={required}
        value={typeof value === "string" ? value : ""}
        placeholder={fieldPlaceholder("aboutMe")}
        onChange={(event) => onChange(event.target.value)}
        className="min-h-32 w-full resize-y rounded-md border border-input bg-background px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-ring"
      />
    );
  const inputType =
    type === "number" || type === "integer"
      ? "number"
      : field.format === "email"
        ? "email"
        : field.format === "uri"
          ? "url"
          : "text";
  const fieldName = nameFromId(id);
  const isSocial = id.startsWith("socials-");
  const displayType = fieldName === "website" || isSocial ? "text" : inputType;
  const stringValue =
    typeof value === "string" || typeof value === "number" ? String(value) : "";
  const normalizedSocial = isSocial
    ? normalizeSocialProfile(fieldName, stringValue)
    : "";
  return (
    <div className="min-w-0 flex-1 space-y-1.5">
      <div className="relative">
        {isSocial && socialIconSources[fieldName] && (
          <span className="pointer-events-none border border-border absolute inset-y-0 left-0 flex w-10 items-center justify-center rounded-l-lg bg-primary">
            <img
              src={socialIconSources[fieldName]}
              alt={`${fieldLabel(fieldName, field)} — ícone`}
              className={cn(
                "size-4 object-contain",
                !["youtube", "bluesky", "linkedin"].includes(fieldName) &&
                  "dark:invert",
              )}
            />
          </span>
        )}
        <Input
          id={id}
          type={displayType}
          inputMode={fieldName === "website" || isSocial ? "url" : undefined}
          required={required}
          placeholder={fieldPlaceholder(fieldName)}
          className={cn("h-10 bg-background", isSocial && "pl-11!")}
          value={stringValue}
          aria-label={isSocial ? fieldLabel(fieldName, field) : undefined}
          onChange={(event) => {
            const next = event.target.value;
            onChange(
              isSocial && /^https?:\/\//i.test(next)
                ? normalizeSocialProfile(fieldName, next)
                : inputValue(field, inputType, next),
            );
          }}
          onBlur={() => {
            if (fieldName === "website")
              onChange(normalizeWebsite(stringValue));
            else if (isSocial)
              onChange(normalizeSocialProfile(fieldName, stringValue));
          }}
        />
      </div>
      {normalizedSocial && (
        <p
          className="truncate text-right text-xs text-muted-foreground"
          title={socialHref(fieldName, normalizedSocial)}
        >
          {socialHref(fieldName, normalizedSocial)}
        </p>
      )}
    </div>
  );
}

const PRONOUNS = [
  "ela/dela",
  "ele/dele",
  "elu/delu",
  "ela/ele",
  "prefiro não informar",
];

function PronounsInput({
  id,
  value,
  onChange,
}: {
  id: string;
  value: JsonValue | undefined;
  onChange: (value: JsonValue) => void;
}) {
  const current = typeof value === "string" ? value : "";
  const isOther = current !== "" && !PRONOUNS.includes(current);
  const [otherSelected, setOtherSelected] = useState(isOther);
  const options = [
    ...PRONOUNS.map((pronoun) => ({ value: pronoun, label: pronoun })),
    { value: "outro", label: "Outro" },
  ];
  return (
    <div className="space-y-2">
      <TimezoneCombobox
        id={id}
        value={isOther || otherSelected ? "outro" : current}
        options={options}
        onChange={(next) => {
          setOtherSelected(next === "outro");
          onChange(next === "outro" ? "" : next);
        }}
      />
      {(isOther || otherSelected) && (
        <Input
          aria-label="Outros pronomes"
          value={isOther ? current : ""}
          placeholder="Digite seus pronomes"
          className="h-10 bg-background"
          onChange={(event) => onChange(event.target.value)}
        />
      )}
    </div>
  );
}

function ArrayDraftInput({
  id,
  value,
  placeholder,
  required,
  onChange,
}: {
  id: string;
  value: JsonValue[];
  placeholder: string;
  required: boolean;
  onChange: (value: JsonValue) => void;
}) {
  const isLanguageField = ["languages", "specializations"].includes(
    nameFromId(id),
  );
  const canonical = value.map(String).join(", ");
  const [draft, setDraft] = useState(isLanguageField ? "" : canonical);

  useEffect(() => {
    if (!isLanguageField) setDraft(canonical);
  }, [canonical, isLanguageField]);

  const commit = (nextDraft: string) => {
    const items = nextDraft
      .split(/[,;\n]+/)
      .map((item) => item.trim())
      .filter(Boolean);
    onChange(items);
    setDraft(items.join(", "));
  };

  const addItem = (item: string) => {
    const next = [...value.map(String), item.trim()].filter(Boolean);
    onChange(next);
    setDraft("");
  };

  if (isLanguageField)
    return (
      <div className="space-y-2">
        {value.length > 0 && (
          <div className="flex flex-wrap gap-1.5">
            {value.map((item, index) => (
              <button
                key={`${String(item)}-${index}`}
                type="button"
                className="rounded-full border border-border bg-muted px-2.5 py-1 text-xs"
                onClick={() =>
                  onChange(
                    value.filter((_, entryIndex) => entryIndex !== index),
                  )
                }
                title="Remover idioma"
              >
                {String(item)} ×
              </button>
            ))}
          </div>
        )}
        <Input
          id={id}
          value={draft}
          placeholder={placeholder}
          className="h-10 bg-background"
          onChange={(event) => {
            const next = event.target.value;
            const parts = next.split(/[,;\n]+/);
            if (parts.length > 1) {
              const added = parts
                .slice(0, -1)
                .map((part) => part.trim())
                .filter(Boolean);
              onChange([...value.map(String), ...added]);
              setDraft(parts.at(-1) ?? "");
            } else setDraft(next);
          }}
          onBlur={() => draft.trim() && addItem(draft)}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              event.preventDefault();
              if (draft.trim()) addItem(draft);
            }
          }}
        />
      </div>
    );

  return (
    <Input
      id={id}
      required={required}
      value={draft}
      placeholder={placeholder}
      className="h-10 bg-background"
      onChange={(event) => setDraft(event.target.value)}
      onBlur={() => commit(draft)}
      onKeyDown={(event) => {
        if (event.key === "Enter") {
          event.preventDefault();
          commit(draft);
        }
      }}
    />
  );
}

function inputValue(
  field: ProfileSchemaNode,
  inputType: string,
  value: string,
): JsonValue {
  if (inputType === "number") return Number(value);
  if (!value && Array.isArray(field.type) && field.type.includes("null")) {
    return null;
  }
  return value;
}

function readValue(values: ProfileData, path: string[]): JsonValue | undefined {
  let current: JsonValue = values;
  for (const segment of path) {
    if (!current || Array.isArray(current) || typeof current !== "object")
      return undefined;
    current = current[segment] as JsonValue;
  }
  return current;
}

function writeValue(
  values: ProfileData,
  path: string[],
  value: JsonValue,
): ProfileData {
  const next = structuredClone(values);
  let current: Record<string, JsonValue> = next;
  path.forEach((segment, index) => {
    if (index === path.length - 1) current[segment] = value;
    else {
      const child = current[segment];
      if (!child || Array.isArray(child) || typeof child !== "object")
        current[segment] = {};
      current = current[segment] as Record<string, JsonValue>;
    }
  });
  return next;
}

const socialIconSources: Record<string, string> = {
  x: xIcon,
  github: githubIcon,
  twitch: twitchIcon,
  bluesky: blueskyIcon,
  discord: discordIcon,
  youtube: youtubeIcon,
  linkedin: linkedinIcon,
  instagram: instagramIcon,
};

function humanize(value: string): string {
  return value
    .replace(/([a-z])([A-Z])/g, "$1 $2")
    .replace(/^./, (letter) => letter.toUpperCase());
}

const FIELD_LABELS: Record<string, string> = {
  preferredName: "Nome de exibição",
  legalName: "Nome completo",
  pronouns: "Pronomes",
  tagline: "Frase de apresentação",
  aboutMe: "Sobre mim",
  role: "Função",
  organization: "Organização",
  contactEmail: "E-mail de contato",
  website: "Site",
  socials: "Redes sociais",
  city: "Cidade",
  region: "Estado ou região",
  country: "País",
  countryCode: "Código do país",
  languages: "Idiomas",
  specializations: "Idiomas",
  timezone: "Fuso horário",
  hideSocials: "Ocultar redes sociais",
  hideLocation: "Ocultar localização",
  hideLegalName: "Ocultar nome completo",
  hideContactEmail: "Ocultar e-mail de contato",
  hideOrganization: "Ocultar organização",
};

const FIELD_PLACEHOLDERS: Record<string, string> = {
  preferredName: "Como você quer ser chamado",
  legalName: "Digite seu nome completo",
  pronouns: "Ex.: ela/dela",
  tagline: "Uma frase curta sobre você",
  aboutMe: "Conte um pouco sobre você, sua experiência e seus interesses…",
  role: "Ex.: Desenvolvedor de software",
  organization: "Ex.: Univents",
  contactEmail: "voce@exemplo.com",
  website: "https://seusite.com",
  city: "Ex.: São Paulo",
  region: "Ex.: São Paulo",
  country: "Ex.: Brasil",
  countryCode: "Ex.: BR",
  languages: "Digite um idioma e separe por vírgula",
  specializations: "Digite um idioma e separe por vírgula",
  x: "Ex.: @univents ou x.com/univents",
  twitter: "Ex.: @univents ou twitter.com/univents",
  instagram: "Ex.: @univents ou instagram.com/univents",
  linkedin: "Ex.: trieoh ou linkedin.com/in/trieoh",
  github: "Ex.: trieoh ou github.com/trieoh",
  twitch: "Ex.: univents ou twitch.tv/univents",
  bluesky: "Ex.: univents.bsky.social",
  youtube: "Ex.: @univents ou youtube.com/@univents",
};

const TIMEZONES = [
  "America/Sao_Paulo",
  "America/Manaus",
  "America/Cuiaba",
  "America/Rio_Branco",
  "America/Noronha",
  "America/New_York",
  "America/Los_Angeles",
  "Europe/Lisbon",
  "Europe/London",
  "Europe/Paris",
  "UTC",
];

function fieldLabel(name: string, field: ProfileSchemaNode): string {
  return FIELD_LABELS[name] ?? field.title ?? humanize(name);
}

function fieldPlaceholder(name: string): string {
  return FIELD_PLACEHOLDERS[name] ?? "Ex.: seu perfil ou informação";
}

function normalizeWebsite(value: string): string {
  const trimmed = value.trim();
  if (!trimmed) return "";
  return `https://${trimmed.replace(/^https?:\/\//i, "").replace(/^\/+/, "")}`;
}

function normalizeSocialProfile(network: string, value: string): string {
  const trimmed = value.trim();
  if (!trimmed) return "";
  const withoutProtocol = trimmed.replace(/^https?:\/\//i, "");
  if (!withoutProtocol.includes("/")) return withoutProtocol.replace(/^@/, "");

  try {
    const url = new URL(`https://${withoutProtocol}`);
    const parts = url.pathname.split("/").filter(Boolean);
    const prefixIndex =
      network === "linkedin"
        ? parts.indexOf("in")
        : network === "bluesky"
          ? parts.indexOf("profile")
          : -1;
    const profile = prefixIndex >= 0 ? parts[prefixIndex + 1] : parts.at(0);
    return (profile ?? trimmed).replace(/^@/, "").replace(/\/$/, "");
  } catch {
    return trimmed.replace(/^@/, "").replace(/\/$/, "");
  }
}

function normalizeProfileLinks(values: ProfileData): ProfileData {
  const next = structuredClone(values);
  if (typeof next.website === "string") {
    next.website = normalizeWebsite(next.website);
  }
  const socials = next.socials;
  if (socials && !Array.isArray(socials) && typeof socials === "object") {
    for (const [network, value] of Object.entries(socials)) {
      if (typeof value === "string") {
        socials[network] = normalizeSocialProfile(network, value);
      }
    }
  }
  return next;
}

function projectProfileValues(
  schema: ProfileSchemaNode,
  values: ProfileData,
): ProfileData {
  const project = (node: ProfileSchemaNode, value: JsonValue): JsonValue => {
    if (
      value === null ||
      Array.isArray(value) ||
      !node.properties ||
      typeof value !== "object"
    )
      return value;
    const result: Record<string, JsonValue> = {};
    for (const [name, child] of Object.entries(node.properties)) {
      const childValue = value[name];
      if (childValue !== undefined) result[name] = project(child, childValue);
    }
    return result;
  };

  return project(schema, values) as ProfileData;
}

function timezoneLabel(timezone: string): string {
  const labels: Record<string, string> = {
    "America/Sao_Paulo": "Brasília — São Paulo (UTC−03:00)",
    "America/Manaus": "Manaus (UTC−04:00)",
    "America/Cuiaba": "Cuiabá (UTC−04:00)",
    "America/Rio_Branco": "Rio Branco (UTC−05:00)",
    "America/Noronha": "Fernando de Noronha (UTC−02:00)",
    UTC: "UTC (tempo universal)",
  };
  return labels[timezone] ?? timezone.replaceAll("_", " ");
}

function nameFromId(id: string): string {
  return id.split("-").at(-1) ?? id;
}
