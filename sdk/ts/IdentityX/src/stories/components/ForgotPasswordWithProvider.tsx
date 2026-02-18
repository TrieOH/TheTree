import type { MouseEvent } from "react";
import { AuthProvider } from "../../react/AuthProvider";
import ForgotPassword from "../../react/components/ForgotPassword/ForgotPassword";

interface ForgotPasswordWithProviderProps {
  loginRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
}

export default function ForgotPasswordWithProvider({loginRedirect}: ForgotPasswordWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <ForgotPassword loginRedirect={loginRedirect} />
    </AuthProvider>
  );
}
