import type { PermissionApi, PermissionDomain, PermissionObject } from "../../types/permission-types";

/**
 * Represents the final output of a permission construction.
 * It encapsulates the permission domain data and provides methods to format it for different layers.
 */
export class PermissionResult {
  /**
   * @param domain The internal structured representation of the permission.
   */
  constructor(private readonly domain: PermissionDomain) {}

  /**
   * Returns the raw, structured domain object.
   * Useful for internal logic where you need to inspect namespaces or actions.
   */
  get value(): PermissionDomain {
    return this.domain;
  }

  /**
   * Converts the permission into an API-ready format.
   * Serializes the object hierarchy into a standard string representation.
   */
  toApi(): PermissionApi {
    return {
      object: serializeObject(this.domain.object),
      action: this.domain.action,
    };
  }
}

/**
 * Internal utility to transform a PermissionObject into a string.
 * * Logic:
 * 1. If no segments exist, it returns a global wildcard "*".
 * 2. Joins segments using "namespace:specifier" format, separated by "/".
 * 3. Appends a suffix (like "*" or "**") if one is present.
 * * @example
 * // Result: "drive:root/folder:work/*"
 */
function serializeObject(obj: PermissionObject): string {
  if (obj.segments.length === 0) return "*";

  const base = obj.segments
    .map(s => `${s.namespace}:${s.specifier}`)
    .join("/");

  return obj.suffix ? `${base}/${obj.suffix}` : base;
}
