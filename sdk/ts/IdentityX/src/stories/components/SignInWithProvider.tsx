import { AuthProvider, SignIn } from "../../react";
import { type MouseEvent } from "react";

interface SignInWithProviderProps {
  signUpRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
  forgotPasswordRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
}

export default function SignInWithProvider({ forgotPasswordRedirect, signUpRedirect }: SignInWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <SignIn 
        forgotPasswordRedirect={forgotPasswordRedirect} 
        signUpRedirect={signUpRedirect}
      />
    </AuthProvider>
  )
}