import type { FieldI, FieldSelectConfigI } from "#/features/fields/model";
import { cn } from "#/shared/lib/utils";
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
  field: FieldI;
  value: unknown;
  error?: string;
  onChange: (fieldId: string, value: unknown) => void;
}

function getPlaceholder(field: FieldI): string {
  return field.placeholder?.value ?? "";
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
  const inputType = getInputType(field.type);
  const placeholder = getPlaceholder(field);

  const handleTextChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(field.id, e.target.value);
  };

  const renderSelect = () => {
    const config = field.config as FieldSelectConfigI | undefined;
    const options = config?.options;
    const behaviour = config?.behaviour ?? "dropdown-radio";
    const isMultiple = behaviour === "checkbox" || behaviour === "dropdown-checkbox";

    if (behaviour === "radio" || behaviour === "checkbox") {
      const selectedValues: string[] = value
        ? Array.isArray(value)
          ? (value as string[])
          : [value as string]
        : [];

      return (
        <div className="flex flex-col gap-2 w-full">
          {options?.map((opt) => {
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
                  name={field.id}
                  value={opt.value}
                  checked={isSelected}
                  onChange={(e) => {
                    if (isMultiple) {
                      const current = Array.isArray(value) ? [...(value as string[])] : [];
                      if (e.target.checked) {
                        onChange(field.id, [...current, opt.value]);
                      } else {
                        onChange(
                          field.id,
                          current.filter((v) => v !== opt.value)
                        );
                      }
                    } else {
                      onChange(field.id, opt.value);
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
          value={(value as string) || ""}
          onChange={(e) => onChange(field.id, e.target.value)}
          className={cn(
            "h-10 w-full appearance-none rounded-md border bg-card px-4 text-sm text-foreground outline-none transition-all focus:border-primary focus:ring-2 focus:ring-primary/5 focus:shadow-sm",
            error ? "border-destructive bg-destructive/5" : "border-border group-hover:border-border/80"
          )}
        >
          <option value="">Select...</option>
          {options?.map((opt) => (
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
    const checked = value === true || value === "true";
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
          onChange={(e) => onChange(field.id, e.target.checked)}
          className="h-4 w-4 accent-primary"
        />
        <span className="text-sm text-foreground">
          {field.title}
          {field.required && <span className="text-destructive"> *</span>}
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
              {value && (value as File).name
                ? (value as File).name
                : "Click to upload or drag a file here"}
            </span>
          </div>
          <input
            type="file"
            className="hidden"
            onChange={(e) => {
              if (e.target.files?.[0]) {
                onChange(field.id, e.target.files[0]);
              }
            }}
          />
        </label>
      </div>
    );
  };

  const renderTextInput = () => {
    let inputValue = (value as string) || "";

    // Normalize date values for standard HTML inputs
    if (inputValue) {
      if (inputType === "date" && inputValue.includes("T")) {
        inputValue = inputValue.split("T")[0];
      } else if (inputType === "datetime-local") {
        const date = new Date(inputValue);
        if (!isNaN(date.getTime())) {
          // Format as YYYY-MM-DDTHH:mm
          inputValue = date.toISOString().slice(0, 16);
        }
      }
    }

    return (
      <div className="relative group w-full">
        <input
          type={inputType}
          step={getStep(field.type)}
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
    switch (field.type) {
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
      {field.type !== "bool" && (
        <label className="mb-1.5 flex items-center gap-2 text-sm font-medium text-foreground">
          <span className="text-primary/70">{getIcon(field.type)}</span>
          {field.title}
          {field.required && <span className="text-destructive"> *</span>}
        </label>
      )}
      {renderField()}
      {field.description && !error && (
        <p className="mt-1 text-xs text-muted-foreground/80">{field.description}</p>
      )}
      {error && <p className="mt-1 text-xs text-destructive font-medium">{error}</p>}
    </motion.div>
  );
}