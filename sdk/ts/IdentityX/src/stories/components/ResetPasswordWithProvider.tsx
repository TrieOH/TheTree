import { AuthProvider, ResetPassword } from "../../react";
import { type MouseEvent } from "react";

interface ResetPasswordWithProviderProps {
  token: string;
  loginRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  isProjectMode?: boolean;
}

export default function ResetPasswordWithProvider({ 
  token,
  loginRedirect,
  isProjectMode = true
}: ResetPasswordWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080" isProjectMode={isProjectMode}>
      <ResetPassword
        token={token}
        loginRedirect={loginRedirect}
      />
    </AuthProvider>
  )
}
