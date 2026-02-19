import { Api } from "./core/api";
import { createServerAuthService } from "./core/services";

/**
 * Creates a new server-side auth service instance with a custom base URL.
 */
export const createServerAuth = (baseURL?: string) => createServerAuthService(new Api(baseURL));