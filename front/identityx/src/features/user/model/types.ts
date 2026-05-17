export interface User {
  id: string;
  user_type: "client" | "project";
  project_id: string | null;
  email: string;
  last_login_at: string | null;
  is_verified: boolean;
  verified_at: string | null;
  created_at: string;
  updated_at: string;
}
