import type { MouseEvent } from "react";
import { AuthProvider, SignUp } from "../../react";

export interface SignUpWithProviderProps {
  loginRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
  flow_id?: string;
}

export default function SignUpWithProvider({ flow_id, loginRedirect }: SignUpWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <SignUp flow_id={flow_id} loginRedirect={loginRedirect} />
    </AuthProvider>
  )
}