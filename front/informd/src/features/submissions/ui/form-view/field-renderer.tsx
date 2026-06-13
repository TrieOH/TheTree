import type { OptionI } from "#/features/fields/model";
import { cn } from "#/shared/lib/utils";
import { formatPhoneMask } from "#/shared/lib/helpers/mask";
import type { FieldAnswerable } from "@trieoh/informd-models";
import {
  Type,
  Mail,
  Hash,
  CheckSquare,
  Calendar as CalendarIcon,
  Clock,
  CalendarClock,
  List,
  Paperclip,
  Phone,
  Link,
  ChevronDown
} from "lucide-react";
import { motion } from "motion/react";

interface FieldRendererProps {
  field: FieldAnswerable;
  value: unknown;
  error?: string;
  onChange: (fieldId: string, value: unknown) => void;
}

function getPlaceholder(field: FieldAnswerable): string {
  const p = field.field.placeholder;
  if (typeof p === 'object' && p !== null && 'value' in p) return (p as { value: string }).value;
  return typeof p === 'string' ? p : "";
}

function getInputType(fieldType: string): string {
  switch (fieldType) {
    case "email":
      return "email";
    case "phone":
      return "tel";
    case "url":
      return "url";
    case "date":
      return "date";
    case "time":
      return "time";
    case "datetime":
      return "datetime-local";
    case "int":
    case "float":
      return "number";
    default:
      return "text";
  }
}

function getStep(fieldType: string): string | undefined {
  if (fieldType === "int") return "1";
  if (fieldType === "float") return "0.01";
  return undefined;
}

function getIcon(fieldType: string) {
  switch (fieldType) {
    case "email": return <Mail className="size-3.5" />;
    case "phone": return <Phone className="size-3.5" />;
    case "url": return <Link className="size-3.5" />;
    case "date": return <CalendarIcon className="size-3.5" />;
    case "time": return <Clock className="size-3.5" />;
    case "datetime": return <CalendarClock className="size-3.5" />;
    case "int":
    case "float": return <Hash className="size-3.5" />;
    case "select": return <List className="size-3.5" />;
    case "file": return <Paperclip className="size-3.5" />;
    case "bool": return <CheckSquare className="size-3.5" />;
    default: return <Type className="size-3.5" />;
  }
}



export function FieldRenderer({ field, value, error, onChange }: FieldRendererProps) {
  const inputType = getInputType(field.field.type);
  const placeholder = getPlaceholder(field);

  const currentValue = value !== undefined ? value : "";

  const handleTextChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(field.field.id, e.target.value);
  };

  const renderSelect = () => {
    const config = field.field_select_config;
    const rawOptions = (config?.options ?? []) as (OptionI | string)[];
    const behaviour = config?.behaviour ?? "dropdown-radio";
    const isMultiple = behaviour === "checkbox" || behaviour === "dropdown-checkbox";

    // Normalize options: convert string[] to OptionI[]
    const options = rawOptions.map(opt => {
      if (typeof opt === 'string') return { label: opt, value: opt };
      return opt;
    });

    const selectedValues: string[] = Array.isArray(currentValue)
      ? currentValue
      : currentValue
        ? [String(currentValue)]
        : [];

    if (behaviour === "radio" || behaviour === "checkbox") {
      return (
        <div className="flex flex-col gap-2 w-full">
          {options.map((opt) => {
            const isSelected = selectedValues.includes(opt.value);
            return (
              <label
                key={opt.value}
                className={cn(
                  "flex h-10 w-full cursor-pointer items-center gap-3 rounded-md border px-4 transition-all",
                  isSelected
                    ? "border-primary bg-primary/5 shadow-xs"
                    : "border-border bg-card/50 hover:border-border/80 hover:bg-card hover:shadow-xs"
                )}
              >
                <input
                  type={isMultiple ? "checkbox" : "radio"}
                  name={field.field.id}
                  value={opt.value}
                  checked={isSelected}
                  onChange={(e) => {
                    if (isMultiple) {
                      if (e.target.checked) {
                        onChange(field.field.id, [...selectedValues, opt.value]);
                      } else {
                        onChange(
                          field.field.id,
                          selectedValues.filter((v) => v !== opt.value)
                        );
                      }
                    } else {
                      onChange(field.field.id, [opt.value]);
                    }
                  }}
                  className="h-4 w-4 accent-primary"
                />
                <span className="text-sm text-foreground">{opt.label}</span>
              </label>
            );
          })}
        </div>
      );
    }

    // Dropdown
    return (
      <div className="relative group w-full">
        <select
          value={selectedValues[0] || ""}
          onChange={(e) => onChange(field.field.id, e.target.value ? [e.target.value] : [])}
          className={cn(
            "h-10 w-full appearance-none rounded-md border bg-card px-4 text-sm text-foreground outline-none transition-all focus:border-primary focus:ring-2 focus:ring-primary/5 focus:shadow-sm",
            error ? "border-destructive bg-destructive/5" : "border-border group-hover:border-border/80"
          )}
        >
          <option value="">{placeholder || "Select..."}</option>
          {options.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
        <div className="pointer-events-none absolute right-4 top-1/2 -translate-y-1/2">
          <ChevronDown className="size-4 text-muted-foreground transition-colors group-hover:text-foreground" />
        </div>
      </div>
    );
  };

  const renderBool = () => {
    const checked = currentValue === true || currentValue === "true";
    return (
      <label
        className={cn(
          "flex h-10 w-full cursor-pointer items-center gap-3 rounded-md border px-4 transition-all",
          checked
            ? "border-primary bg-primary/5 shadow-xs"
            : "border-border bg-card/50 hover:border-border/80 hover:bg-card hover:shadow-xs"
        )}
      >
        <input
          type="checkbox"
          checked={checked}
          onChange={(e) => onChange(field.field.id, e.target.checked)}
          className="h-4 w-4 accent-primary"
        />
        <span className="text-sm text-foreground">
          {field.field.title}
          {field.field.required && <span className="text-destructive"> *</span>}
        </span>
      </label>
    );
  };

  const renderFile = () => {
    return (
      <div className="w-full">
        <label
          className={cn(
            "flex h-20 w-full cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed transition-all",
            error ? "border-destructive/30 bg-destructive/5" : "border-border bg-card hover:border-primary hover:bg-primary/5 hover:shadow-sm"
          )}
        >
          <div className="flex items-center gap-2">
            <Paperclip className="size-4 text-muted-foreground" />
            <span className="text-xs text-muted-foreground text-center line-clamp-1 px-4">
              {currentValue instanceof File
                ? currentValue.name
                : placeholder || "Click to upload or drag a file here"}
            </span>
          </div>
          <input
            type="file"
            className="hidden"
            onChange={(e) => {
              if (e.target.files?.[0]) onChange(field.field.id, e.target.files[0]);
            }}
          />
        </label>
      </div>
    );
  };

  const renderTextInput = () => {
    const fieldType = field.field.type;

    // --- Phone: masked input ---
    if (fieldType === "phone") {
      return (
        <div className="relative group w-full">
          <input
            type="tel"
            placeholder={placeholder || "(dd) dddd-dddd"}
            value={currentValue ? formatPhoneMask(String(currentValue)) : ""}
            onChange={(e) => onChange(field.field.id, formatPhoneMask(e.target.value))}
            className={cn(
              "h-10 w-full rounded-md border bg-card px-4 text-sm text-foreground outline-none transition-all placeholder:text-muted-foreground/50 focus:border-primary focus:ring-2 focus:ring-primary/5 focus:shadow-sm",
              error ? "border-destructive bg-destructive/5" : "border-border group-hover:border-border/80"
            )}
          />
        </div>
      );
    }

    // --- URL: https:// prefix badge ---
    if (fieldType === "url") {
      const rawValue = currentValue ? String(currentValue) : "";
      const stored = rawValue.startsWith("https://") ? rawValue.slice(8) : rawValue;

      return (
        <div className="flex w-full rounded-md border border-border has-focus-within:border-primary has-focus-within:ring-2 has-focus-within:ring-primary/5 transition-all group hover:border-border/80">
          <span className="inline-flex items-center px-3 text-xs font-medium text-muted-foreground bg-muted/50 border-r border-border rounded-l-md select-none whitespace-nowrap">
            https://
          </span>
          <input
            type="text"
            placeholder="example.com"
            value={stored}
            onChange={(e) => {
              const typed = e.target.value;
              const clean = typed.replace(/^https?:\/\//i, "");
              onChange(field.field.id, clean ? `https://${clean}` : undefined);
            }}
            className="h-10 w-full min-w-0 bg-card px-3 text-sm text-foreground outline-none rounded-r-md"
          />
        </div>
      );
    }

    // --- Standard text input ---
    let inputValue = String(currentValue);

    // Normalize date values for standard HTML inputs
    if (inputValue && inputValue !== "null" && inputValue !== "undefined") {
      if (inputType === "date" && inputValue.includes("T")) {
        inputValue = inputValue.split("T")[0];
      } else if (inputType === "datetime-local") {
        const date = new Date(inputValue);
        if (!isNaN(date.getTime())) inputValue = date.toISOString().slice(0, 16);
      }
    } else inputValue = "";
    return (
      <div className="relative group w-full">
        <input
          type={inputType}
          step={getStep(field.field.type)}
          value={inputValue}
          placeholder={placeholder}
          onChange={handleTextChange}
          className={cn(
            "h-10 w-full rounded-md border bg-card px-4 text-sm text-foreground outline-none transition-all placeholder:text-muted-foreground/50 focus:border-primary focus:ring-2 focus:ring-primary/5 focus:shadow-sm",
            error ? "border-destructive bg-destructive/5" : "border-border group-hover:border-border/80"
          )}
        />
      </div>
    );
  };

  const renderField = () => {
    switch (field.field.type) {
      case "select":
        return renderSelect();
      case "bool":
        return renderBool();
      case "file":
        return renderFile();
      default:
        return renderTextInput();
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="mb-5"
    >
      {field.field.type !== "bool" && (
        <label className="mb-1.5 flex items-center gap-2 text-sm font-medium text-foreground">
          <span className="text-primary/70">{getIcon(field.field.type)}</span>
          {field.field.title}
          {field.field.required && <span className="text-destructive"> *</span>}
        </label>
      )}
      {renderField()}
      {field.field.description && !error && (
        <p className="mt-1 text-xs text-muted-foreground/80">{field.field.description}</p>
      )}
      {error && <p className="mt-1 text-xs text-destructive font-medium">{error}</p>}
    </motion.div>
  );
}