import { expect } from "vitest";
import { z } from "zod";

// ============================================================================
// MATCHER INTERFACE
// ============================================================================

export interface Matcher {
    _isMatcher: true;
    match(val: unknown, path: string): unknown;
}

function isMatcher(v: unknown): v is Matcher {
    return typeof v === "object" && v !== null && (v as any)._isMatcher === true;
}

// ============================================================================
// PRIMITIVE MATCHERS
// ============================================================================

/** Matches any non-empty string */
export const AnyString: Matcher = {
    _isMatcher: true,
    match(val, path) {
        expect(typeof val, `${path}: expected string`).toBe("string");
        expect(val as string, `${path}: expected non-empty string`).not.toBe("");
        return val;
    },
};

/** Matches any number */
export const AnyNumber: Matcher = {
    _isMatcher: true,
    match(val, path) {
        expect(typeof val, `${path}: expected number`).toBe("number");
        return val;
    },
};

/** Validates UUID v4 format */
export const AnyUUID: Matcher = {
    _isMatcher: true,
    match(val, path) {
        expect(typeof val, `${path}: expected string UUID`).toBe("string");
        expect(val as string, `${path}: expected valid UUID`).toMatch(
            /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i
        );
        return val;
    },
};

/** Validates RFC3339 / ISO datetime string */
export const AnyDate: Matcher = {
    _isMatcher: true,
    match(val, path) {
        expect(typeof val, `${path}: expected string date`).toBe("string");
        const d = new Date(val as string);
        expect(isNaN(d.getTime()), `${path}: expected valid date, got "${val}"`).toBe(false);
        return val;
    },
};

/** Validates that value is not null/undefined */
export const NotEmpty: Matcher = {
    _isMatcher: true,
    match(val, path) {
        expect(val, `${path}: expected non-null value`).not.toBeNull();
        expect(val, `${path}: expected non-undefined value`).not.toBeUndefined();
        return val;
    },
};

/** Validates that value is null */
export const Null: Matcher = {
    _isMatcher: true,
    match(val, path) {
        expect(val, `${path}: expected null`).toBeNull();
        return val;
    },
};

// ============================================================================
// STORE MATCHERS - Capture values for later assertions
// ============================================================================

/** Captures the value into a ref object: { current: T } */
export function Store<T = unknown>(
    ref: { current: T },
    matcher?: Matcher
): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            const result = matcher ? matcher.match(val, path) : val;
            ref.current = result as T;
            return result;
        },
    };
}

/** Captures value directly into a holder array (index 0), avoids let reassign */
export function StoreString(
    into: string[],
    matcher?: Matcher
): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            const result = matcher ? matcher.match(val, path) : val;
            expect(typeof result, `${path}: StoreString expected string`).toBe("string");
            into[0] = result as string;
            return result;
        },
    };
}

// ============================================================================
// EQUALITY MATCHERS - Assert value AND validate format
// ============================================================================

/** Assert string equality, optionally also validate format */
export function AsString(value: string, matcher?: Matcher): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            expect(val, `${path}: string mismatch`).toBe(value);
            if (matcher) matcher.match(val, path);
            return val;
        },
    };
}

/** Assert number equality */
export function AsInt(value: number, matcher?: Matcher): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            expect(val, `${path}: number mismatch`).toBe(value);
            if (matcher) matcher.match(val, path);
            return val;
        },
    };
}

/** Assert boolean equality */
export function AsBool(value: boolean, matcher?: Matcher): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            expect(val, `${path}: boolean mismatch`).toBe(value);
            if (matcher) matcher.match(val, path);
            return val;
        },
    };
}

/** Assert deep equality against a reference value */
export function SameAs(ref: unknown): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            expect(val, `${path}: SameAs mismatch`).toEqual(ref);
            return val;
        },
    };
}

// ============================================================================
// ARRAY MATCHERS
// ============================================================================

/**
 * Index an array of objects by a key field, then validate each by key.
 * AllowExtra: if true, ignores keys not in spec.
 */
export function ByKey(
    key: string,
    spec: Record<string, unknown>,
    allowExtra = false
): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            expect(Array.isArray(val), `${path}: ByKey expected array`).toBe(true);
            const arr = val as Record<string, unknown>[];

            const actual: Record<string, unknown> = {};
            for (const item of arr) {
                const k = String(item[key]);
                actual[k] = item;
            }

            for (const k of Object.keys(spec)) {
                expect(actual, `${path}: missing key "${k}"`).toHaveProperty(k);
            }

            if (!allowExtra) {
                for (const k of Object.keys(actual)) {
                    expect(spec, `${path}: unexpected key "${k}"`).toHaveProperty(k);
                }
            }

            const results: Record<string, unknown> = {};
            for (const [k, itemSpec] of Object.entries(spec)) {
                results[k] = validate(actual[k], itemSpec, `${path}[${k}]`);
            }

            return results;
        },
    };
}

/** Validate each array element against the same spec */
export function Each(spec: unknown): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            expect(Array.isArray(val), `${path}: Each expected array`).toBe(true);
            return (val as unknown[]).map((item, i) =>
                validate(item, spec, `${path}[${i}]`)
            );
        },
    };
}

/** Validate array elements in exact positional order */
export function InOrder(specs: unknown[]): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            expect(Array.isArray(val), `${path}: InOrder expected array`).toBe(true);
            const arr = val as unknown[];
            expect(arr.length, `${path}: InOrder length mismatch`).toBe(specs.length);
            return arr.map((item, i) => validate(item, specs[i], `${path}[${i}]`));
        },
    };
}

/** Validate a specific index of an array */
export function AtIndex(index: number, spec: unknown): Matcher {
    return {
        _isMatcher: true,
        match(val, path) {
            expect(Array.isArray(val), `${path}: AtIndex expected array`).toBe(true);
            const arr = val as unknown[];
            expect(arr.length, `${path}: array too short for index ${index}`).toBeGreaterThan(index);
            return validate(arr[index], spec, `${path}[${index}]`);
        },
    };
}

// ============================================================================
// CORE VALIDATE FUNCTION
// ============================================================================

type Spec = Matcher | Record<string, unknown> | unknown[] | unknown;

function validate(val: unknown, spec: Spec, path: string): unknown {
    // Matcher
    if (isMatcher(spec)) {
        return spec.match(val, path);
    }

    // Object spec
    if (
        typeof spec === "object" &&
        spec !== null &&
        !Array.isArray(spec)
    ) {
        const specMap = spec as Record<string, unknown>;
        expect(
            typeof val === "object" && val !== null && !Array.isArray(val),
            `${path}: expected object`
        ).toBe(true);
        const obj = val as Record<string, unknown>;
        const results: Record<string, unknown> = {};

        for (const [key, expectedVal] of Object.entries(specMap)) {
            results[key] = validate(obj[key], expectedVal, `${path}.${key}`);
        }

        return results;
    }

    // Array spec (positional)
    if (Array.isArray(spec)) {
        expect(Array.isArray(val), `${path}: expected array`).toBe(true);
        const arr = val as unknown[];
        expect(arr.length, `${path}: array length mismatch`).toBe(spec.length);
        return arr.map((item, i) => validate(item, spec[i], `${path}[${i}]`));
    }

    // Primitive equality
    expect(val, `${path}: value mismatch`).toBe(spec);
    return val;
}

/**
 * Validate data against a spec. Returns captured values.
 * Use for flexible partial matching.
 */
export function Validate(data: unknown, spec: Spec): unknown {
    return validate(data, spec, "$");
}

/**
 * Like Validate but fails if the actual object has extra fields not in spec.
 */
export function ValidateExact(data: unknown, spec: Record<string, unknown>): unknown {
    const obj = data as Record<string, unknown>;
    for (const key of Object.keys(obj)) {
        expect(spec, `extra field "${key}" not in spec`).toHaveProperty(key);
    }
    return validate(data, spec, "$");
}

// ============================================================================
// ZOD SCHEMA VALIDATION (for Univents-style schema checks)
// ============================================================================

export function validateSchema<T>(schema: z.ZodType<T>, data: unknown): T {
    const result = schema.safeParse(data);
    if (!result.success) {
        console.error("Schema errors:", JSON.stringify(result.error.format(), null, 2));
        expect(result.success, "Zod schema validation failed").toBe(true);
        throw new Error("unreachable");
    }
    return result.data;
}