import { Schema, SchemaIssue, type SchemaAST } from "effect";

const standardFormatter = SchemaIssue.makeFormatterStandardSchemaV1();
const defaultFormatter = SchemaIssue.makeFormatterDefault();

export interface ValidationSuccess<T> {
  readonly success: true;
  readonly data: T;
}

export interface ValidationFailure {
  readonly success: false;
  readonly error: {
    readonly message: string;
    readonly issues: ReadonlyArray<{
      readonly message: string;
      readonly path?: ReadonlyArray<PropertyKey>;
    }>;
  };
}

export type ValidationResult<T> = ValidationSuccess<T> | ValidationFailure;

export interface StandardSchemaProperties<A> {
  readonly "~standard": {
    readonly version: 1;
    readonly vendor: string;
    readonly validate: (
      value: unknown,
    ) =>
      | { readonly value: A; readonly issues?: undefined }
      | { readonly issues: ReadonlyArray<{ readonly message: string; readonly path?: ReadonlyArray<PropertyKey> }> };
    readonly types?: {
      readonly input: unknown;
      readonly output: A;
    };
  };
}

export type Writable<T> = T extends (...args: readonly unknown[]) => unknown
  ? T
  : T extends readonly (infer U)[]
    ? Writable<U>[]
    : T extends object
      ? { -readonly [P in keyof T]: Writable<T[P]> }
      : T;

export type StandardCompatibleSchema<T, S = unknown> = S &
  StandardSchemaProperties<T> & {
    readonly Type: T;
    parse(input: unknown): T;
    safeParse(input: unknown): ValidationResult<T>;
  };

export type Infer<S> = S extends { readonly "Type": infer T }
  ? Writable<T>
  : S extends StandardCompatibleSchema<infer T, unknown>
    ? Writable<T>
    : never;

/**
 * Options for check combinators (e.g. custom error message).
 */
export interface CheckOptions {
  readonly message?: string;
}

const DECODE_OPTIONS: SchemaAST.ParseOptions = {
  errors: "all",
};

/**
 * Extracts a map of { [fieldName]: errorMessage } from a ValidationResult.
 */
export function getFieldErrors<T extends object>(
  result: {
    readonly success: boolean;
    readonly error?: {
      readonly message?: string;
      readonly issues?: ReadonlyArray<{
        readonly message: string;
        readonly path?: ReadonlyArray<PropertyKey>;
      }>;
    };
  },
): Partial<Record<keyof T, string>> {
  if (result.success || !result.error) return {};
  const map: Partial<Record<keyof T, string>> = {};
  for (const issue of result.error.issues ?? []) {
    const key = issue.path?.[0] as keyof T | undefined;
    if (key && !map[key]) {
      map[key] = issue.message;
    }
  }
  return map;
}

/**
 * Wraps an Effect Schema with convenient `.parse` and `.safeParse` methods,
 * preserving full compatibility with TanStack Router search validation and form logic.
 */
export function createSchema<S extends Schema.ConstraintDecoder<unknown>>(
  schema: S,
): StandardCompatibleSchema<Infer<S>, S> {
  const decodeSync = Schema.decodeUnknownSync(schema);
  const decodeResult = Schema.decodeUnknownResult(schema, DECODE_OPTIONS);
  const std = Schema.toStandardSchemaV1(schema, { parseOptions: DECODE_OPTIONS });

  const enhanced = Object.assign(schema, {
    parse(input: unknown): Infer<S> {
      return decodeSync(input) as Infer<S>;
    },
    safeParse(input: unknown): ValidationResult<Infer<S>> {
      const res = decodeResult(input);
      if (res._tag === "Success") {
        return { success: true, data: res.success as Infer<S> };
      }
      const formatted = standardFormatter(res.failure.issue);
      const issues = (formatted.issues ?? []) as ReadonlyArray<{
        readonly message: string;
        readonly path?: ReadonlyArray<PropertyKey>;
      }>;
      const firstMessage =
        issues[0]?.message ?? defaultFormatter(res.failure.issue);
      return {
        success: false,
        error: {
          message: firstMessage,
          issues,
        },
      };
    },
    "~standard": {
      version: 1 as const,
      vendor: "effect",
      validate: (input: unknown) => {
        const res = std["~standard"].validate(input);
        return res as ReturnType<StandardSchemaProperties<Infer<S>>["~standard"]["validate"]>;
      },
      types: {
        input: undefined as unknown,
        output: undefined as unknown as Infer<S>,
      },
    },
  });

  return enhanced as unknown as StandardCompatibleSchema<Infer<S>, S>;
}

// Convenient Effect check combinators preserving underlying schema type and supporting custom messages
export const minLength = (n: number, options?: CheckOptions) =>
  <S extends Schema.Top & { readonly Type: { readonly length: number } }>(self: S): S =>
    self.check(
      Schema.isMinLength(n, options?.message ? { message: options.message } : undefined) as SchemaAST.Check<S["Type"]>,
    ) as unknown as S;

export const maxLength = (n: number, options?: CheckOptions) =>
  <S extends Schema.Top & { readonly Type: { readonly length: number } }>(self: S): S =>
    self.check(
      Schema.isMaxLength(n, options?.message ? { message: options.message } : undefined) as SchemaAST.Check<S["Type"]>,
    ) as unknown as S;

export const pattern = (reg: RegExp, options?: CheckOptions) =>
  <S extends Schema.Top & { readonly Type: string }>(self: S): S =>
    self.check(
      Schema.isPattern(reg, options?.message ? { message: options.message } : undefined) as SchemaAST.Check<S["Type"]>,
    ) as unknown as S;

export const uuid = (options?: CheckOptions) =>
  <S extends Schema.Top & { readonly Type: string }>(self: S): S =>
    self.check(
      Schema.isUUID(undefined, options?.message ? { message: options.message } : undefined) as SchemaAST.Check<S["Type"]>,
    ) as unknown as S;

export const min = (n: number, options?: CheckOptions) =>
  <S extends Schema.Top & { readonly Type: number }>(self: S): S =>
    self.check(
      Schema.isGreaterThanOrEqualTo(n, options?.message ? { message: options.message } : undefined) as SchemaAST.Check<S["Type"]>,
    ) as unknown as S;

export const max = (n: number, options?: CheckOptions) =>
  <S extends Schema.Top & { readonly Type: number }>(self: S): S =>
    self.check(
      Schema.isLessThanOrEqualTo(n, options?.message ? { message: options.message } : undefined) as SchemaAST.Check<S["Type"]>,
    ) as unknown as S;

export const positive = (options?: CheckOptions) =>
  <S extends Schema.Top & { readonly Type: number }>(self: S): S =>
    self.check(
      Schema.isGreaterThan(0, options?.message ? { message: options.message } : undefined) as SchemaAST.Check<S["Type"]>,
    ) as unknown as S;

export const nonNegative = (options?: CheckOptions) =>
  <S extends Schema.Top & { readonly Type: number }>(self: S): S =>
    self.check(
      Schema.isGreaterThanOrEqualTo(0, options?.message ? { message: options.message } : undefined) as SchemaAST.Check<S["Type"]>,
    ) as unknown as S;

export function refine<A = unknown>(
  predicate: (value: A) => boolean,
  options?: CheckOptions,
) {
  return <S extends Schema.Top & { readonly Type: A }>(self: S): S =>
    (self as unknown as { check: (f: unknown) => S }).check(
      Schema.makeFilter(
        (input: unknown) => predicate(input as A),
        options?.message ? { message: options.message } : undefined,
      ),
    );
}

export { Schema };
