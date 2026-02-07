import { env } from "@/env";
import { createAuthenticatedFetch } from "@trieoh/node-auth-sdk";
import { toast } from "sonner";

interface RawCommonResponse {
  module: string;
  message: string;
  timestamp: string;
  code: number;
}

interface RawErrorData {
  debug?: string[];
}

interface RawSuccessResponse<T> extends RawCommonResponse {
  data: T;
}

interface RawErrorResponse extends RawCommonResponse {
  data?: RawErrorData; 
  error_id: string;
  trace?: string[];
}

type RawApiResponse<T> = RawSuccessResponse<T> | RawErrorResponse;

// --- Standardized Client-Facing API Response Types ---
export type ApiSuccessResponse<T> = RawSuccessResponse<T> & { success: true };
export type ApiErrorResponse = RawErrorResponse & { success: false; details?: unknown; };
export type ApiResponse<T> = ApiSuccessResponse<T> | ApiErrorResponse;

const authenticatedFetch = createAuthenticatedFetch({
  baseURL: env.VITE_API_URL,
});

/**
 * A wrapper around the standard fetch API that provides a standardized response object.
 * It automatically includes authentication and handles errors.
 *
 * @param path The URL path (e.g., "/projects", "/users/123").
 * @param init Optional RequestInit object for fetch.
 * @returns A Promise that resolves to an ApiResponse object.
 */
export async function authFetcher<TData>(
  path: string,
  init?: RequestInit
): Promise<ApiResponse<TData>> {
  try {
    let baseUrlString = env.VITE_API_URL;
    if (!baseUrlString.startsWith('http://') && !baseUrlString.startsWith('https://')) {
      baseUrlString = `http://${baseUrlString}`; // Default to http if missing
    }
    const baseUrl = new URL(baseUrlString);
    const fullUrl = new URL(path, baseUrl).toString();

    const response = await authenticatedFetch(fullUrl, init);
    const rawResponse: RawApiResponse<TData> = await response.json().catch(() => ({
        module: "Client",
        message: response.statusText || "Unknown error",
        timestamp: new Date().toISOString(),
        code: response.status,
        error_id: "CLIENT_PARSE_ERROR",
        trace: ["Failed to parse API response as JSON"],
      })
    );

    if (!response.ok) {
      const errorResponse = rawResponse as RawErrorResponse;
      const errorMessage = errorResponse.message || response.statusText || "An unknown error occurred";
      
      toast.error(errorMessage, {description: errorResponse.trace?.join("\n")});

      return {
        success: false,
        module: errorResponse.module,
        message: errorMessage,
        timestamp: errorResponse.timestamp,
        code: errorResponse.code,
        error_id: errorResponse.error_id,
        trace: errorResponse.trace,
        data: errorResponse.data, // Include data (debug) if present
      };
    }

    // If we reach here, it's a successful response
    const successResponse = rawResponse as RawSuccessResponse<TData>;
    return {
      success: true,
      module: successResponse.module,
      message: successResponse.message,
      timestamp: successResponse.timestamp,
      code: successResponse.code,
      data: successResponse.data,
    };

  } catch (error) {
    const errorMessage = error instanceof Error 
      ? error.message 
      : "A network or unknown error occurred.";
    
    toast.error(errorMessage);

    return {
      success: false,
      module: "Client-module",
      message: errorMessage,
      timestamp: new Date().toISOString(),
      code: 0, // Indicate client-side network error
      error_id: "CLIENT_NETWORK_ERROR",
      trace: error instanceof Error ? [error.stack || errorMessage] : [errorMessage],
      details: error, // for raw error object from catch
    };
  }
}

/**
 * A fetcher for TanStack Query that uses authFetcher but throws errors,
 * 
 * @param path The URL path (e.g., "/projects", "/users/123").
 * @returns The data if the fetch is successful.
 * @throws An error if the fetch fails.
 */
export const tanstackQueryFetcher = async <TData>(path: string): Promise<TData> => {
  const response = await authFetcher<TData>(path);
  if (!response.success) throw new Error(response.message);
  return response.data;
};
