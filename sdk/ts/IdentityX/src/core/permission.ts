import type { CheckPermissionRequest } from "../types/permission-types";

type ForbiddenChar =
  | ' ' | '.' | '-' | '/' | '@' | '!' | '#' | '$'
  | '%' | '&' | '(' | ')' | '+' | '=' | '[' | ']'
  | '{' | '}' | '|' | ':' | ';' | '"' | "'" | '<'
  | '>' | ',' | '?' | '`' | '~' | '\\';

type Alphabet =
  | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z'
  | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z';

/**
 * Validates identifiers at compile time to ensure they follow the required naming conventions.
 * 
 * Rules: 
 * 1. Cannot be empty.
 * 2. Must be a single asterisk ('*') OR start with a letter and contain only alphanumeric characters or underscores.
 */
type ValidIdent<T extends string> =
  T extends "" ? "Error: Identifier cannot be empty" :
  T extends "*" ? "*" :
  T extends `${infer First}${string}`
  ? First extends Alphabet
  ? (T extends `${string}${ForbiddenChar}${string}`
    ? `Error: Identifier '${T}' contains invalid characters`
    : T)
  : `Error: Identifier '${T}' must start with a letter`
  : never;

type AllMethods = 'project' | 'scope' | 'user' | 'object' | 'action' | 'resource';

export type BuilderMethods<Called extends AllMethods> = {
  /** 
   * Sets the Project ID for this permission check.
   */
  project: (id: string) => PermissionBuilder<Called | 'project'>;

  /** 
   * Sets the Scope ID (e.g., a specific organization or group) for this check.
   */
  scope: (id: string) => PermissionBuilder<Called | 'scope'>;

  /** 
   * Sets the User (Entity) ID that is performing the action.
   * This is a MANDATORY field.
   */
  user: (id: string) => PermissionBuilder<Called | 'user'>;

  /** 
   * Sets the target Object for the permission check.
   * This is a MANDATORY field.
   * 
   * Validation: Must be '*' or start with a letter and contain only alphanumeric/underscores.
   * @example .object('documents')
   * @example .object('*')
   */
  object: <T extends string>(obj: ValidIdent<T>) => PermissionBuilder<Called | 'object'>;

  /** 
   * Sets the Action to be checked against the object.
   * This is a MANDATORY field.
   * 
   * Validation: Must be '*' or start with a letter and contain only alphanumeric/underscores.
   * @example .action('read')
   * @example .action('*')
   */
  action: <T extends string>(act: ValidIdent<T>) => PermissionBuilder<Called | 'action'>;

  /** 
   * Attaches additional resource metadata (JSON) to the check.
   * Useful for attribute-based access control (ABAC).
   */
  resource: (res: Record<string, unknown>) => PermissionBuilder<Called | 'resource'>;
};

/**
 * Builder for permission checks.
 * 
 * This builder ensures that mandatory fields (user, object, action) are provided 
 * before the request can be finalized. It also prevents duplicate calls to the same method.
 * 
 * @template Called - Tracks which methods have already been called.
 */
export type PermissionBuilder<Called extends AllMethods = never> =
  Pick<BuilderMethods<Called>, Exclude<AllMethods, Called>> & (
    'user' | 'object' | 'action' extends Called
    ? {
      /** Finalizes the builder and returns the raw request object. */
      build(): CheckPermissionRequest;
      /** Returns the request object for JSON serialization. */
      toJSON(): CheckPermissionRequest
    }
    : unknown
  );

/**
 * Represents a builder that has fulfilled all mandatory requirements (user, object, and action).
 * This interface allows the service to consume the builder once it is in a valid state.
 */
export interface CompletePermissionBuilder {
  /** Finalizes the builder and returns the raw request object. */
  build(): CheckPermissionRequest;
  /** Returns the request object for JSON serialization. */
  toJSON(): CheckPermissionRequest;
}

const IDENT_REGEX = /^(\*|[a-zA-Z][a-zA-Z0-9_]*)$/;

/**
 * Internal implementation of the PermissionBuilder.
 */
class PermissionBuilderImpl {
  constructor(private readonly req: Partial<CheckPermissionRequest> = {}) { }

  private clone(patch: Partial<CheckPermissionRequest>): PermissionBuilderImpl {
    return new PermissionBuilderImpl({ ...this.req, ...patch });
  }

  /** @see BuilderMethods.project */
  project(id: string) { return this.clone({ project_id: id }); }

  /** @see BuilderMethods.scope */
  scope(id: string) { return this.clone({ scope_id: id }); }

  /** @see BuilderMethods.user */
  user(id: string) { return this.clone({ entity_id: id }); }

  /** @see BuilderMethods.resource */
  resource(res: Record<string, unknown>) { return this.clone({ resource: res }); }

  /** @see BuilderMethods.object */
  object(obj: string) {
    if (!obj) throw new Error("Permission object cannot be empty.");
    if (!IDENT_REGEX.test(obj))
      throw new Error(`Invalid object: "${obj}". Must start with a letter and contain only alphanumeric/underscores, or be '*'.`);
    return this.clone({ object: obj });
  }

  /** @see BuilderMethods.action */
  action(act: string) {
    if (!act) throw new Error("Permission action cannot be empty.");
    if (!IDENT_REGEX.test(act))
      throw new Error(`Invalid action: "${act}". Must start with a letter and contain only alphanumeric/underscores, or be '*'.`);
    return this.clone({ action: act });
  }

  /** @see PermissionBuilder.build */
  build(): CheckPermissionRequest {
    if (!this.req.entity_id || !this.req.object || !this.req.action) {
      throw new Error('PermissionBuilder: user, object, and action are required.');
    }
    return this.req as CheckPermissionRequest;
  }

  /** @see PermissionBuilder.toJSON */
  toJSON() { return this.build(); }
}

/**
 * Starts building a new permission check request.
 * 
 * Mandatory fields: `.user(id)`, `.object(name)`, and `.action(name)`.
 * Optional fields: `.project(id)`, `.scope(id)`, and `.resource(data)`.
 * 
 * @example
 * // Basic usage
 * const check = await auth.checkPermission(
 *   permission()
 *     .user('user-uuid')
 *     .object('documents')
 *     .action('read')
 * );
 */
export const permission = (): PermissionBuilder =>
  new PermissionBuilderImpl() as unknown as PermissionBuilder;
