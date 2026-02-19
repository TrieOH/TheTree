import type { MouseEvent } from "react";
import { AuthProvider, TabbedSignUp } from "../../react";
import type { TabbedFlowI } from "../../react/components/SignUp/TabbedSignUp";

export interface TabbedSignUpWithProviderProps {
  loginRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
  flowIds: TabbedFlowI[];
}

export default function TabbedSignUpWithProvider({ flowIds, loginRedirect }: TabbedSignUpWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <TabbedSignUp flowIds={flowIds} loginRedirect={loginRedirect} />
    </AuthProvider>
  )
}