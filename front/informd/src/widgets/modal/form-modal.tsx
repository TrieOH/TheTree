import { useForm, Controller } from "react-hook-form";
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
  FieldValues,
  Path,
  SubmitHandler,
} from "react-hook-form";

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
}: PropsI<T>) {
  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = useForm<T>({
    resolver: zodResolver(schema),
    defaultValues: defaultValues,
  });

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
    const error = errors[fieldName];

    if (field.type === "select") {
      return (
        <Controller
          name={fieldName}
          control={control}
          render={({ field: { onChange, value } }) => (
            <Select onValueChange={onChange} defaultValue={value}>
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
      <form
        id={formId}
        onSubmit={handleSubmit(handleFormSubmit)}
        className="space-y-6"
      >
        {fields.map((field) => {
          const fieldName = field.name as Path<T>;
          const error = errors[fieldName];
          return (
            <div className="space-y-2" key={"t_" + field.name.toString()}>
              <Label
                htmlFor={fieldName}
                className="text-xs font-semibold text-muted-foreground"
              >
                {field.label}
              </Label>
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
    </Modal>
  );
}