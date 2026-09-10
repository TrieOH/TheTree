import { AuthProvider, SignUp } from "../../react";
import { type MouseEvent } from "react";

interface SignUpWithProviderProps {
  loginRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  isProjectMode?: boolean;
}

export default function SignUpWithProvider({ 
  loginRedirect,
  isProjectMode = true
}: SignUpWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080" isProjectMode={isProjectMode}>
      <SignUp
        loginRedirect={loginRedirect}
      />
    </AuthProvider>
  )
}