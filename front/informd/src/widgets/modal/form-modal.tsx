import { useForm, Controller, useWatch, FormProvider } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Modal } from "./modal";
import { cn } from "#/shared/lib/utils";
import { Input } from "#/shared/ui/shadcn/input";
import { Textarea } from "#/shared/ui/shadcn/textarea";
import { Label } from "#/shared/ui/shadcn/label";
import { Button } from "#/shared/ui/shadcn/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#/shared/ui/shadcn/select";
import { AlertCircle } from "lucide-react";
import type { FieldDefinition } from "#/shared/model/form-types";
import type { ZodType } from "zod";
import { useEffect } from "react";
import type {
  DefaultValues,
  FieldError,
  FieldValues,
  Path,
  SubmitHandler,
} from "react-hook-form";
import { formatPhoneMask } from "#/shared/lib/helpers/mask";

/** Safely access nested error objects for dot-notation paths like "select_config.behaviour". */
function getNestedError(errors: Record<string, unknown>, path: string) {
  return path.split(".").reduce<Record<string, unknown> | undefined>((acc, key) => {
    if (acc && typeof acc === "object" && key in acc) {
      return (acc)[key] as Record<string, unknown>;
    }
    return undefined;
  }, errors);
}

export interface PropsI<T> {
  isOpen: boolean;
  onSubmit: (data: T) => void;
  onClose: () => void;
  title: string;
  description: string;
  buttonTitle: string;
  formId: string;
  defaultValues?: DefaultValues<T>;
  fields: FieldDefinition<T>[];
  schema: ZodType<T, T>;
  disabled?: boolean;
  children?: React.ReactNode;
}

export default function FormModal<T extends FieldValues>({
  isOpen,
  onClose,
  title,
  description,
  formId,
  onSubmit,
  fields,
  schema,
  defaultValues,
  buttonTitle,
  disabled = false,
  children,
}: PropsI<T>) {
  const methods = useForm<T>({
    resolver: zodResolver(schema),
    defaultValues: defaultValues,
  });

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = methods;

  const watchedValues = useWatch({ control });

  useEffect(() => {
    if (isOpen) {
      reset(defaultValues);
    }
  }, [isOpen, defaultValues, reset]);

  const handleFormSubmit: SubmitHandler<T> = (data) => {
    onSubmit(data);
  };

  const renderField = (field: FieldDefinition<T>) => {
    const fieldName = field.name as Path<T>;
    const error = getNestedError(errors, String(field.name)) as
      | FieldError
      | undefined;

    if (field.type === "boolean") {
      return (
        <Controller
          name={fieldName}
          control={control}
          render={({ field: { onChange, value } }) => (
            <div
              role="button"
              tabIndex={0}
              onClick={() => onChange(!value)}
              onKeyDown={(e) => {
                if (e.key === " " || e.key === "Enter") {
                  e.preventDefault();
                  onChange(!value);
                }
              }}
              className={cn(
                "flex items-center justify-between rounded-md border p-3.5 transition-all cursor-pointer group select-none",
                value
                  ? "border-primary/40 bg-primary/5 shadow-xs"
                  : "border-border bg-card/50 hover:border-border/80 hover:bg-muted/10"
              )}
            >
              <div className="space-y-1 pr-4">
                <div className="text-sm font-medium tracking-tight text-foreground">
                  {field.label}
                </div>
                {field.placeholder && (
                  <div className="text-[11px] leading-relaxed text-muted-foreground/80 font-medium">
                    {field.placeholder}
                  </div>
                )}
              </div>
              <div
                className={cn(
                  "relative inline-flex h-5 w-9 shrink-0 items-center rounded-full transition-colors",
                  value ? "bg-primary" : "bg-muted-foreground/30"
                )}
              >
                <span
                  className={cn(
                    "pointer-events-none block h-4 w-4 rounded-full bg-background shadow-md ring-0 transition-transform duration-200",
                    value ? "translate-x-4.5" : "translate-x-0.5"
                  )}
                />
              </div>
            </div>
          )}
        />
      );
    }

    if (field.type === "select") {
      return (
        <Controller
          name={fieldName}
          control={control}
          render={({ field: { onChange, value } }) => (
            <Select
              onValueChange={(val) => {
                const strVal = String(val);
                onChange(strVal === "true" ? true : strVal === "false" ? false : strVal);
              }}
              value={value ?? ""}
            >
              <SelectTrigger
                id={fieldName}
                className={cn(
                  "rounded-sm border-border w-full",
                  error && "border-destructive"
                )}
              >
                <SelectValue placeholder={field.placeholder} />
              </SelectTrigger>
              <SelectContent>
                {field.options?.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
      );
    }

    if (field.type === "textarea") {
      return (
        <Textarea
          id={fieldName}
          placeholder={field.placeholder}
          rows={field.rows ?? 3}
          className={cn(
            "rounded-md border-border min-h-20 resize-y",
            error && "border-destructive"
          )}
          {...register(fieldName)}
        />
      );
    }

    // Special case: default_value renders differently depending on the selected field type
    if (fieldName === "default_value" as unknown as Path<T>) {
      const currentType = (watchedValues as Record<string, unknown>)["type"];

      // "select" type → dropdown with options from select_config.options
      if (currentType === "select") {
        const selectConfig = (watchedValues as Record<string, unknown>)["select_config"] as Record<string, unknown> | undefined;
        const rawOptions = selectConfig?.["options"] as string | undefined;
        const parsedOptions = (rawOptions ?? "")
          .split("\n")
          .map(s => s.trim())
          .filter(Boolean);

        return (
          <Controller
            name={fieldName}
            control={control}
            render={({ field: { onChange, value } }) => (
              <Select
                onValueChange={(val) => onChange(val === "__none__" ? undefined : val)}
                value={value ?? "__none__"}
              >
                <SelectTrigger
                  id={fieldName}
                  className={cn(
                    "rounded-sm border-border w-full",
                    error && "border-destructive"
                  )}
                >
                  <SelectValue placeholder="Select a default value…" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__none__">None</SelectItem>
                  {parsedOptions.map((opt) => (
                    <SelectItem key={opt} value={opt}>
                      {opt}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          />
        );
      }

      // "date" type → native date picker
      if (currentType === "date") {
        return (
          <Controller
            name={fieldName}
            control={control}
            render={({ field: { onChange, value } }) => (
              <Input
                type="date"
                value={value ?? ""}
                onChange={(e) => onChange(e.target.value || undefined)}
                className={cn(
                  "rounded-md border-border",
                  error && "border-destructive"
                )}
              />
            )}
          />
        );
      }

      // "time" type → native time picker
      if (currentType === "time") {
        return (
          <Controller
            name={fieldName}
            control={control}
            render={({ field: { onChange, value } }) => (
              <Input
                type="time"
                value={value ?? ""}
                onChange={(e) => onChange(e.target.value || undefined)}
                className={cn(
                  "rounded-md border-border",
                  error && "border-destructive"
                )}
              />
            )}
          />
        );
      }

      // "datetime" type → native datetime-local picker
      if (currentType === "datetime") {
        return (
          <Controller
            name={fieldName}
            control={control}
            render={({ field: { onChange, value } }) => (
              <Input
                type="datetime-local"
                value={value ?? ""}
                onChange={(e) => onChange(e.target.value || undefined)}
                className={cn(
                  "rounded-md border-border",
                  error && "border-destructive"
                )}
              />
            )}
          />
        );
      }

      // "phone" type → masked input (dd) dddd-dddd
      if (currentType === "phone") {
        return (
          <Controller
            name={fieldName}
            control={control}
            render={({ field: { onChange, value } }) => (
              <Input
                type="text"
                placeholder="(dd) dddd-dddd"
                value={value ?? ""}
                onChange={(e) => onChange(formatPhoneMask(e.target.value))}
                className={cn(
                  "rounded-md border-border",
                  error && "border-destructive"
                )}
              />
            )}
          />
        );
      }

      // "url" type → input with a fixed https:// prefix badge
      if (currentType === "url") {
        return (
          <Controller
            name={fieldName}
            control={control}
            render={({ field: { onChange, value } }) => {
              const rawValue = value ?? "";
              const stored = rawValue.startsWith("https://") ? rawValue.slice(8) : rawValue;

              return (
                <div className="flex rounded-md border-border border has-focus-within:border-ring has-focus-within:ring-3 has-focus-within:ring-ring/50 transition-[color,box-shadow]">
                  <span className="inline-flex items-center px-2.5 text-xs font-medium text-muted-foreground bg-muted/50 border-r border-border rounded-l-md select-none whitespace-nowrap">
                    https://
                  </span>
                  <input
                    type="text"
                    placeholder="example.com"
                    value={stored}
                    onChange={(e) => {
                      const typed = e.target.value;
                      // Strip any https:// the user may have typed/pasted so we don't double up
                      const clean = typed.replace(/^https?:\/\//i, "");
                      onChange(clean ? `https://${clean}` : undefined);
                    }}
                    className="h-9 w-full min-w-0 bg-transparent px-2.5 py-1 text-base outline-none md:text-sm"
                  />
                </div>
              );
            }}
          />
        );
      }
    }

    return (
      <Input
        id={fieldName}
        type={field.type}
        placeholder={field.placeholder}
        min={field.min}
        max={field.max}
        disabled={field.disabled || disabled}
        className={cn(
          "rounded-md border-border",
          error && "border-destructive"
        )}
        {...register(fieldName, field.type === "number" ? { valueAsNumber: true } : undefined)}
      />
    );
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      description={description}
    >
      <FormProvider {...methods}>
        <form
          id={formId}
          onSubmit={handleSubmit(handleFormSubmit)}
          className="space-y-6"
        >
          {fields.map((field) => {
            // Skip field if its dependency is not met
            if (field.dependsOn) {
              const depValue = (watchedValues)[field.dependsOn.field as string];
              const accepted = Array.isArray(field.dependsOn.value)
                ? field.dependsOn.value
                : [field.dependsOn.value];
              const isMet = accepted.some(v => depValue === v || String(depValue) === String(v));
              if (!isMet) return null;
            }

            const fieldName = field.name as Path<T>;
            const error = getNestedError(errors, String(field.name)) as FieldError | undefined;
            const isBoolean = field.type === "boolean";

            return (
              <div className="space-y-2" key={"t_" + field.name.toString()}>
                {!isBoolean && (
                  <Label
                    htmlFor={fieldName}
                    className="text-xs font-semibold text-muted-foreground"
                  >
                    {field.label}
                    {field.required === false ? (
                      <span className="ml-1.5 font-normal text-muted-foreground/60">(optional)</span>
                    ) : (
                      <span className="ml-1 text-destructive">*</span>
                    )}
                  </Label>
                )}
                {renderField(field)}
                {error && (
                  <span
                    className={cn(
                      "text-[10px] font-bold text-destructive uppercase",
                      "tracking-widest flex items-start gap-1"
                    )}
                  >
                    <AlertCircle className="w-3 h-3" />
                    <span className="-mt-px">{error.message?.toString()}</span>
                  </span>
                )}
              </div>
            );
          })}
          {children}
          <div className="flex justify-end pt-2">
            <Button
              type="submit"
              disabled={disabled}
              className="w-full rounded-sm font-bold transition-all h-10"
            >
              {buttonTitle}
            </Button>
          </div>
        </form>
      </FormProvider>
    </Modal>
  );
}