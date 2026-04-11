import { CheckPermissionRequest, PermissionBuilder } from "../types/permission-types";

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
