import type { MouseEvent } from "react";
import { AuthProvider, TabbedSignUp } from "../../react";

export interface TabbedSignUpWithProviderProps {
  loginRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
  flowIds: { label: string; value: string; }[];
}

export default function TabbedSignUpWithProvider({ flowIds, loginRedirect }: TabbedSignUpWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <TabbedSignUp flowIds={flowIds} loginRedirect={loginRedirect} />
    </AuthProvider>
  )
}