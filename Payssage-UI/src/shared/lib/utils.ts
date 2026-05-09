import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Converts a percentage value to Basis Points (BPS).
 * 1% = 100 BPS
 * 
 * @param percentage - The percentage value (e.g., 1.5 for 1.5%)
 * @returns The value in Basis Points (e.g., 150)
 */
export const percentageToBps = (percentage: number): number => {
  return Math.round(percentage * 100);
};

/**
 * Converts Basis Points (BPS) to a percentage value.
 * 100 BPS = 1%
 * 
 * @param bps - The value in Basis Points (e.g., 150)
 * @returns The percentage value (e.g., 1.5)
 */
export const bpsToPercentage = (bps: number): number => {
  return bps / 100;
};

/**
 * Constrains a number between a minimum and maximum value.
 * 
 * @param value - The value to constrain
 * @param min - The minimum allowed value
 * @param max - The maximum allowed value
 * @returns The constrained value
 */
export const clamp = (value: number, min: number, max: number): number => {
  return Math.max(min, Math.min(max, value));
};
