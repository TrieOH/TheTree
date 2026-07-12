import type {
  FieldPath,
  FieldPathByValue,
  FieldValues,
  UseFormReturn
} from "react-hook-form";
import type { ReactNode } from "react";

/**
 * The subset of UseFormReturn that field renderers actually need.
 * Deliberately excludes `handleSubmit` / `control`, which are the only
 * members whose type depends on the zod schema's *output* (transformed)
 * type — this lets a single-generic form API be shared by every field
 * regardless of whether the schema applies `.transform()`.
 */
export type FieldFormApi<TFieldValues extends FieldValues> = Pick<
  UseFormReturn<TFieldValues>,
  "register" | "watch" | "getValues" | "setValue" | "trigger" | "resetField" | "unregister" | "formState"
>;

/**
 * Shared by every field kind. `layout` controls how much horizontal
 * space the field takes in the step's grid — "half" lets two fields
 * (e.g. Logo + Banner) sit side by side like in the design.
 */
export interface BaseFieldConfig<TFieldValues extends FieldValues> {
  visibleIf?: VisibilityRule<TFieldValues>;
  layout?: "full" | "half";
  optional?: boolean;
}


/**
 * Native input variants supported by the text field renderer.
 * Adding a new *string-based* variant is just a new entry here,
 * no extra component or branching logic needed.
 */
export type TextFieldInputType = "text" | "email" | "url" | "password" | "tel" | "number";

export interface TextFieldConfig<TFieldValues extends FieldValues> extends BaseFieldConfig<TFieldValues> {
  kind: "text";
  /** Dot-path into the form values, fully typed against the zod schema. */
  name: FieldPath<TFieldValues>;
  label: string;
  placeholder?: string;
  description?: string;
  inputType?: TextFieldInputType;
  disabled?: boolean;
}

export interface UrlFieldConfig<TFieldValues extends FieldValues> extends BaseFieldConfig<TFieldValues> {
  kind: "url";
  name: FieldPath<TFieldValues>;
  label: string;
  placeholder?: string;
  description?: string;
  disabled?: boolean;
  /** Prepend `https://` automatically when the user leaves the field
   * without a scheme. Defaults to `true`. */
  autoPrependScheme?: boolean;
}

export interface ComboboxOption {
  value: string;
  label: string;
  description?: string;
}

/**
 * Searchable single-select. `options` can be a static list or an async
 * loader (e.g. hitting your API as the user types) — same field kind
 * either way, the renderer just debounces calls to the loader.
 */
export interface ComboboxFieldConfig<TFieldValues extends FieldValues> extends BaseFieldConfig<TFieldValues> {
  kind: "combobox";
  /** Dot-path holding the selected option's `value` (a string). */
  name: FieldPath<TFieldValues>;
  label: string;
  placeholder?: string;
  description?: string;
  disabled?: boolean;
  emptyMessage?: string;
  options: ComboboxOption[] | ((query: string) => Promise<ComboboxOption[]>);
  /** Debounce (ms) applied when `options` is a loader function. Defaults to 250. */
  debounceMs?: number;
}

export interface CustomFieldRenderArgs<TFieldValues extends FieldValues> {
  form: FieldFormApi<TFieldValues>;
}

/**
 * Escape hatch for bespoke composite UI (e.g. the social-network picker)
 * that doesn't fit a single simple input, while still living inside the
 * same step/field-list/visibility machinery as every other field.
 */
export interface CustomFieldConfig<TFieldValues extends FieldValues> extends BaseFieldConfig<TFieldValues> {
  kind: "custom";
  /** Identifier only, used as React key. Not necessarily a form path. */
  name: string;
  render: (args: CustomFieldRenderArgs<TFieldValues>) => ReactNode;
}

/** Single image (logo, banner, avatar, ...). `name` must point to a
 * string-valued path (e.g. `logo_url`). */
export interface ImageFieldConfig<TFieldValues extends FieldValues> extends BaseFieldConfig<TFieldValues> {
  kind: "image";
  name: FieldPathByValue<TFieldValues, string | null | undefined>;
  label: string;
  hint?: string;
  accept?: string;
  maxSizeMB?: number;
  /** Fires with the running added/removed sets so you can call your
   * own save/remove endpoints at submit time (see README). */
  onTrackingChange?: (change: ImageFieldChange) => void;
}

/** "N" images (gallery). `name` must point to a string[]-valued path
 * (e.g. `gallery_urls`). */
export interface GalleryFieldConfig<TFieldValues extends FieldValues> extends BaseFieldConfig<TFieldValues> {
  kind: "gallery";
  name: FieldPathByValue<TFieldValues, string[] | null | undefined>;
  label: string;
  hint?: string;
  accept?: string;
  maxSizeMB?: number;
  maxItems?: number;
  onTrackingChange?: (change: ImageFieldChange) => void;
}

export type FieldConfig<TFieldValues extends FieldValues> =
  | TextFieldConfig<TFieldValues>
  | UrlFieldConfig<TFieldValues>
  | ComboboxFieldConfig<TFieldValues>
  | CustomFieldConfig<TFieldValues>
  | ImageFieldConfig<TFieldValues>
  | GalleryFieldConfig<TFieldValues>;

export type FieldKind = FieldConfig<FieldValues>["kind"];

export interface StepConfig<TFieldValues extends FieldValues> {
  id: string;
  label: string;
  description?: string;
  fields: FieldConfig<TFieldValues>[];
}

/**
 * Declarative visibility rules. Keeps "show field X only if field Y has a
 * value / equals something" out of ad-hoc JSX conditionals.
 */
export type VisibilityRule<TFieldValues extends FieldValues> =
  | { type: "equals"; field: FieldPath<TFieldValues>; value: unknown }
  | { type: "notEquals"; field: FieldPath<TFieldValues>; value: unknown }
  | { type: "hasValue"; field: FieldPath<TFieldValues> }
  | { type: "predicate"; predicate: (values: TFieldValues) => boolean };


/**
 * A minimal, backend-agnostic reference to an uploaded image — just
 * enough for a parent form to know what to persist/delete via its own
 * dedicated save/remove endpoints.
 */
export interface UploadedImage {
  /** Stable local id (React key + removal tracking). Not a backend id. */
  id: string;
  url: string;
}

export type ImageItemStatus = "existing" | "processing" | "ready" | "error";

export interface ImageItem {
  id: string;
  /** Object URL while processing; final signed URL once "ready", or the
   * original URL for "existing" items. */
  url: string;
  status: ImageItemStatus;
  file?: File;
  errorMessage?: string;
  /** True for items that were already present when the field mounted
   * (i.e. came from `values` in edit mode) — false for anything the
   * user added in this session. */
  isExisting: boolean;
}

/**
 * Emitted whenever the tracked set of "images added this session" /
 * "existing images the user removed this session" changes. This is
 * deliberately separate from the field's actual form value (which only
 * ever holds the final list of URLs to submit) — the parent uses this
 * to call its own dedicated save/remove endpoints.
 */
export interface ImageFieldChange {
  added: UploadedImage[];
  removed: UploadedImage[];
}

export interface ImageUploadAdapter {
  /** Runs the server-side preprocess step and returns the final public URL. */
  preprocess: (file: File) => Promise<{ url: string }>;
}
