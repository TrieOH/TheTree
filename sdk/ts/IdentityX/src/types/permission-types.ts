export interface CheckPermissionRequest {
  project_id?: string;
  scope_id?: string;
  entity_id: string;
  object: string;
  action: string;
  resource?: Record<string, unknown>;
}

export interface CheckPermissionResponse {
  allowed: boolean;
}
