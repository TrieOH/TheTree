import type { PermissionApi, PermissionDomain, PermissionObject } from "../../types/permission-types";

export class PermissionResult {
  constructor(private readonly domain: PermissionDomain) {}

  get value(): PermissionDomain {
    return this.domain;
  }

  toApi(): PermissionApi {
    return {
      object: serializeObject(this.domain.object),
      action: this.domain.action,
    };
  }
}

function serializeObject(obj: PermissionObject): string {
  if (obj.segments.length === 0) return "*";

  const base = obj.segments
    .map(s => `${s.namespace}:${s.specifier}`)
    .join("/");

  return obj.suffix ? `${base}/${obj.suffix}` : base;
}
