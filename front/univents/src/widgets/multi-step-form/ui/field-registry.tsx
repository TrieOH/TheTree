import type { FieldValues } from "react-hook-form";
import type { ReactElement } from "react";
import type { FieldConfig, FieldFormApi, FieldKind } from "../model/types";
import { TextFieldRenderer } from "./fields/text-field";
import { CustomFieldRenderer } from "./fields/custom-field";
import { ImageFieldRenderer } from "./fields/image-field";
import { GalleryFieldRenderer } from "./fields/gallery-field";
import { UrlFieldRenderer } from "./fields/url-field";
import { ComboboxFieldRenderer } from "./fields/combobox-field";

export type FieldRenderer<TFieldValues extends FieldValues> = (props: {
  field: FieldConfig<TFieldValues>;
  form: FieldFormApi<TFieldValues>;
}) => ReactElement | null;

/**
 * Single source of truth mapping a field "kind" to the component that
 * knows how to render it. Adding a new field type (select, checkbox,
 * date, ...) means adding one entry here + one renderer file — never
 * touching a big if/switch somewhere else.
 */
export type FieldRegistry<TFieldValues extends FieldValues> = Record<FieldKind, FieldRenderer<TFieldValues>>;

export function createFieldRegistry<TFieldValues extends FieldValues>(): FieldRegistry<TFieldValues> {
  return {
    text: TextFieldRenderer,
    url: UrlFieldRenderer,
    combobox: ComboboxFieldRenderer,
    custom: CustomFieldRenderer,
    image: ImageFieldRenderer,
    gallery: GalleryFieldRenderer,
  };
}

export function renderField<TFieldValues extends FieldValues>(
  field: FieldConfig<TFieldValues>,
  form: FieldFormApi<TFieldValues>,
  registry: FieldRegistry<TFieldValues>,
): ReactElement | null {
  const Renderer = registry[field.kind];
  return <Renderer field={field} form={form} />;
}
