export type ObjSegment = {
  namespace: string;
  specifier: string;
};

export type ObjSuffix = "*" | "**" | null;

export interface PermissionObject {
  segments: ObjSegment[];
  suffix: ObjSuffix;
}

export interface PermissionDomain {
  object: PermissionObject;
  action: string; // "field:update" | "*"
}

export interface PermissionApi {
  object: string; // "project:1/schema:2/**"
  action: string; // "field:update"
}
