import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'
import z from 'zod';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

function getDefaultInvalidValue(schema: z.ZodTypeAny): unknown {
  if (schema instanceof z.ZodString) return "";
  if (schema instanceof z.ZodNumber) return NaN;
  if (schema instanceof z.ZodBoolean) return null;
  if (schema instanceof z.ZodArray) return [];
  if (schema instanceof z.ZodObject) return {};
  return undefined;
}


export function getFieldError(schema: z.ZodTypeAny, value?: string): string[] {
  const result = schema.safeParse(value || getDefaultInvalidValue(schema));
  if (result.success) return [];

  return result.error.issues.map(i => i.message);
}
