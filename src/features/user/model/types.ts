import type { JsonValue } from "@/shared/model/types";


export interface User {
  id: string;
  email: string;
  is_active: boolean;
  is_verified: boolean;
  user_type: string;
  project_id: string;
  metadata: Record<string, JsonValue>
  last_login_at: string;
  verified_at: string;
  created_at: string;
  updated_at: string;
}
