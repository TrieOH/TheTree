import { AuthProvider, SignIn } from "../../react";
import { type MouseEvent } from "react";

interface SignInWithProviderProps {
  signUpRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  forgotPasswordRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  isProjectMode?: boolean;
}

export default function SignInWithProvider({ 
  forgotPasswordRedirect, 
  signUpRedirect,
  isProjectMode = true
}: SignInWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080" isProjectMode={isProjectMode}>
      <SignIn
        forgotPasswordRedirect={forgotPasswordRedirect}
        signUpRedirect={signUpRedirect}
      />
    </AuthProvider>
  )
}