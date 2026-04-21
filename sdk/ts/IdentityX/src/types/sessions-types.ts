export interface SessionI {
  session_id: string;
  project_id: string;
  user_id: string;
  user_type: 'client' | 'project';
  family_id: string;
  token_id: string;
  user_agent: string;
  user_ip: string;
  created_at: string;
  updated_at: string;
  expires_at: string;
  issued_at: string;
  revoked_at: string | null;
}

export interface CurrentSessionI {
  user_id: string;
  project_id: string | null;
  session_id: string | null;
}