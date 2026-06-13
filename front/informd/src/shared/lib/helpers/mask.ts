/**
 * Formats a phone number string as (dd) dddd-dddd or (dd) ddddd-dddd.
 * Strips all non-digit characters first, then applies the mask.
 */
export function formatPhoneMask(value: string): string {
    const digits = value.replace(/\D/g, "").slice(0, 11);
    if (digits.length <= 2) {
        return digits.length ? `(${digits}` : "";
    }
    if (digits.length <= 7) {
        return `(${digits.slice(0, 2)}) ${digits.slice(2)}`;
    }
    return `(${digits.slice(0, 2)}) ${digits.slice(2, 7)}-${digits.slice(7)}`;
}
