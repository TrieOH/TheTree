import { useForm } from "react-hook-form";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema"
import { Modal } from "./modal";
import { cn, percentageToBps, clamp } from "#/shared/lib/utils";
import { Input } from "#/shared/ui/shadcn/input";
import { Label } from "#/shared/ui/shadcn/label";
import { Button } from "#/shared/ui/shadcn/button";
import { AlertCircle, Percent, Info } from "lucide-react";
import type { FieldDefinition } from "#/shared/model/form-types";
import type { ZodType } from "zod";
import type { DefaultValues, FieldValues, Path, PathValue } from "react-hook-form";


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
  schema: ZodType<T>;
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
  disabled = false
}: PropsI<T>) {

  const { register, reset, handleSubmit, watch, setValue, formState: { errors } } = useForm<T>({
    resolver: standardSchemaResolver(schema),
    defaultValues: defaultValues,
  });

  const handleFormSubmit = (data: T) => {
    onSubmit(data);
    reset();
  };

  const renderField = (field: FieldDefinition<T>) => {
    const fieldName = field.name as Path<T>;
    const error = errors[fieldName];
    const value = watch(fieldName);

    if (field.type === 'percentage') {
      return (
        <div className="space-y-2">
          <div className="relative">
            <Input
              id={fieldName}
              type="number"
              step="0.01"
              placeholder={field.placeholder}
              className={cn(
                "rounded-none border-border focus-visible:ring-0 font-bold pr-10",
                "focus-visible:border-primary transition-colors",
                error && "border-destructive"
              )}
              {...register(fieldName, {
                valueAsNumber: true,
                onChange: (e) => {
                  const val = parseFloat(e.target.value);
                  if (!isNaN(val)) {
                    const constrained = clamp(val, 0, 100);
                    if (constrained !== val) {
                      setValue(fieldName, constrained.toString() as PathValue<T, Path<T>>);
                    }
                  }
                }
              })}
            />
            <div className="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-muted-foreground">
              <Percent className="h-4 w-4" />
            </div>
          </div>

          {value !== undefined && value !== null && !isNaN(Number(value)) && (
            <div className="flex items-center gap-1.5 px-1 py-0.5 bg-primary/5 border border-primary/10 animate-in fade-in duration-300">
              <Info className="w-3 h-3 text-primary/60" />
              <span className="text-[9px] mt-0.5 font-black uppercase tracking-wider text-primary/70">
                Equivalent to {percentageToBps(Number(value))} Basis Points (BPS)
              </span>
            </div>
          )}
        </div>
      );
    }

    return (
      <Input
        id={fieldName}
        type={field.type}
        placeholder={field.placeholder}
        className={cn(
          "rounded-none border-border focus-visible:ring-0 font-bold",
          "focus-visible:border-primary transition-colors",
          error && "border-destructive"
        )}
        {...register(fieldName)}
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
      <form id={formId} onSubmit={handleSubmit(handleFormSubmit)} className="space-y-6">
        {fields.map(field => {
          const fieldName = field.name as Path<T>;
          const error = errors[fieldName];
          return (
            <div className="space-y-2" key={"t_" + field.name.toString()}>
              <Label
                htmlFor={fieldName}
                className="text-[10px] font-black uppercase tracking-[0.2em]"
              >
                {field.label}
              </Label>
              {renderField(field)}
              {error && (
                <span className={cn(
                  "text-[10px] font-bold text-destructive uppercase",
                  "tracking-widest flex items-start gap-1"
                )}>
                  <AlertCircle className="w-3 h-3" />
                  <span className="-mt-px">{error.message?.toString()}</span>
                </span>
              )}
            </div>
          )
        })}
        <div className="flex justify-end pt-2">
          <Button
            type="submit"
            disabled={disabled}
            className="w-full rounded-none font-black uppercase tracking-widest transition-all h-12"
          >
            {buttonTitle}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
