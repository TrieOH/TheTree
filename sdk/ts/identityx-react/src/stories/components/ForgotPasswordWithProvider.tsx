import type { MouseEvent } from "react";
import { AuthProvider, ForgotPassword } from "../../react";

interface ForgotPasswordWithProviderProps {
  loginRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  isProjectMode?: boolean;
}

export default function ForgotPasswordWithProvider({
  loginRedirect,
  isProjectMode = true
}: ForgotPasswordWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080" isProjectMode={isProjectMode}>
      <ForgotPassword loginRedirect={loginRedirect} />
    </AuthProvider>
  );
}
