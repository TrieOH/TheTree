import { GenericForm } from "./generic-form"
import type { ZodType } from "zod"
import type { DefaultValues, FieldValues } from "react-hook-form"
import type { FormFieldI } from "@/shared/model/field"
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerDescription,
} from '@/shared/ui/shadcn/drawer'

interface FormDrawerProps<T extends FieldValues> {
  idPrefix?: string
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description?: string
  fields: readonly FormFieldI<T>[]
  schema: ZodType<T>
  onSubmit: (data: T) => void | Promise<void>
  defaultValues?: DefaultValues<T>
  submitLabel?: string
  loading?: boolean
}

export function FormDrawer<T extends FieldValues>({
  idPrefix,
  open,
  onOpenChange,
  title,
  description,
  fields,
  schema,
  onSubmit,
  defaultValues,
  submitLabel,
  loading,
}: FormDrawerProps<T>) {
  const handleFormSubmit = async (data: T) => {
    await onSubmit(data)
    onOpenChange(false)
  }

  const handleCancel = () => {
    onOpenChange(false)
  }

  return (
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent className="z-60 rounded-t-2xl border-t border-border bg-card max-h-[90vh]">
        <DrawerHeader className="pb-4 border-b border-border">
          <DrawerTitle className="text-base font-semibold text-left">
            {title}
          </DrawerTitle>
          <DrawerDescription className="text-sm text-left">
            {description ?? 'Preencha os campos abaixo'}
          </DrawerDescription>
        </DrawerHeader>


        <GenericForm
          idPrefix={idPrefix}
          fields={fields}
          schema={schema}
          onSubmit={handleFormSubmit}
          onCancel={handleCancel}
          defaultValues={defaultValues}
          submitLabel={submitLabel}
          loading={loading}
        />
      </DrawerContent>
    </Drawer>
  )
}
