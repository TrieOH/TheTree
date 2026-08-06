import type { JsonValue } from "@trieoh/identityx-sdk-ts";
import { useEffect, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { Input } from "@/shared/ui/shadcn/input";
import type { ProfileSchemaNode } from "../model/profile-data";
import { socialHref } from "../model/profile-data";
import {
  fieldLabel,
  fieldPlaceholder,
  inputValue,
  nameFromId,
  normalizeSocialProfile,
  normalizeWebsite,
  socialIconSources,
  TIMEZONES,
  timezoneLabel,
} from "./profile-editor-utils";
import { TimezoneCombobox } from "./timezone-combobox";

export function SchemaInput({
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

export const PRONOUNS = [
  "ela/dela",
  "ele/dele",
  "elu/delu",
  "ela/ele",
  "prefiro não informar",
];

export function PronounsInput({
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
        placeholder="Selecione seus pronomes…"
        searchPlaceholder="Buscar pronome…"
        emptyMessage="Nenhum pronome encontrado."
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

export function ArrayDraftInput({
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
