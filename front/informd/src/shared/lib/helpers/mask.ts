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

/**
 * Validates a phone number (masked or raw). Must have at least 10 digits (DDD + number).
 */
export function isValidPhone(value: string): boolean {
    const digits = value.replace(/\D/g, "");
    return digits.length >= 10 && digits.length <= 11;
}

/**
 * Validates a public URL string. Accepts with or without protocol.
 * Requires a valid TLD (e.g. .com, .org) or localhost / IP.
 */
export function isValidUrl(value: string): boolean {
    if (!value) return false;
    try {
        const url = /^https?:\/\//i.test(value)
            ? new URL(value)
            : new URL(`https://${value}`);
        const hostname = url.hostname;

        // Allow localhost
        if (hostname === "localhost" || hostname === "127.0.0.1") return true;

        // Allow valid IPv4 addresses
        if (/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/.test(hostname)) return true;

        // Hostname must only contain valid domain characters (letters, digits, dots, hyphens)
        if (!/^[a-zA-Z0-9.-]+$/.test(hostname)) return false;

        // Require at least one dot for a public domain (e.g. example.com, sub.domain.co.uk)
        if (!hostname.includes(".")) return false;

        // TLD must have at least 2 characters after the last dot
        const tld = hostname.split(".").pop() ?? "";
        return tld.length >= 2;
    } catch {
        return false;
    }
}
