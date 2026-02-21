import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'
import z from 'zod';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function truncateString(str: string, start: number, end: number) {
  if (str.length <= start + end) return str;
  return `${str.slice(0, start)}...${str.slice(str.length - end)}`;
}

export function maskStringMiddle(
  str: string,
  start: number,
  end: number,
  maskChar = "•••"
) {
  if (str.length <= start + end) return str;

  const hiddenLength = str.length - (start + end);

  return str.slice(0, start) +
    maskChar.repeat(Math.max(3, Math.min(hiddenLength, 8))) +
    str.slice(str.length - end)
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
