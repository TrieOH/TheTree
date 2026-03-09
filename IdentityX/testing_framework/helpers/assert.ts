import { expect } from "vitest";

// ============================================================================
// RESPONSE ASSERTION HELPERS
// Mirrors the Go Response methods: HasMessage, HasErrID, ValidationError
// These operate on the raw axios error response body.
// ============================================================================

export interface ErrorBody {
    module?: string;
    message?: string;
    error_id?: string;
    trace?: string[];
    [key: string]: unknown;
}

/**
 * Asserts that an axios error response has the expected HTTP status.
 * Returns the parsed error body for further assertions.
 */
export function assertStatus(e: unknown, expectedStatus: number): ErrorBody {
    const err = e as any;
    const actual = err?.response?.status;
    expect(
        actual,
        `Expected HTTP ${expectedStatus} but got ${actual}.\nBody: ${JSON.stringify(err?.response?.data, null, 2)}`
    ).toBe(expectedStatus);
    return err?.response?.data ?? {};
}

/** Assert response body contains expected message substring */
export function assertMessage(body: ErrorBody, expected: string): void {
    expect(
        body.message ?? "",
        `Expected message to contain "${expected}", got "${body.message}"`
    ).toContain(expected);
}

/** Assert response body has the expected error_id */
export function assertErrID(body: ErrorBody, expectedErrID: string): void {
    expect(
        body.error_id,
        `Expected error_id "${expectedErrID}", got "${body.error_id}"`
    ).toBe(expectedErrID);
}

/** Assert validation error trace contains all expected substrings */
export function assertValidationError(
    body: ErrorBody,
    ...expectedErrors: string[]
): void {
    expect(body.message, "Expected validation error message").toContain("Validation failed");
    expect(Array.isArray(body.trace), "Expected trace array").toBe(true);

    const trace = body.trace ?? [];
    expect(
        trace.length,
        `Expected ${expectedErrors.length} trace entries, got ${trace.length}.\nTrace: ${JSON.stringify(trace)}`
    ).toBe(expectedErrors.length);

    for (const expected of expectedErrors) {
        const found = trace.some((entry) => entry.includes(expected));
        expect(
            found,
            `Expected trace to contain "${expected}".\nActual trace: ${JSON.stringify(trace)}`
        ).toBe(true);
    }
}

/**
 * Wraps an async call that is expected to fail with a specific HTTP status.
 * Returns the parsed error body for further assertions.
 *
 * Usage:
 *   const body = await shouldFail(client.http.post(...), 400);
 *   assertMessage(body, "email already in use");
 */
export async function shouldFail(
    promise: Promise<unknown>,
    expectedStatus: number
): Promise<ErrorBody> {
    try {
        await promise;
        expect.fail(`Expected request to fail with ${expectedStatus} but it succeeded`);
    } catch (e: any) {
        return assertStatus(e, expectedStatus);
    }
    // unreachable but satisfies TS
    return {};
}