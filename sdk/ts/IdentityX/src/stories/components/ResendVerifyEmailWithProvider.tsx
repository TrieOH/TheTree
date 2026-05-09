import { AuthProvider, ResendVerifyEmail } from "../../react";

interface ResendVerifyEmailWithProviderProps {
  isProjectMode?: boolean;
}

export default function ResendVerifyEmailWithProvider({ 
  isProjectMode = true
}: ResendVerifyEmailWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080" isProjectMode={isProjectMode}>
      <ResendVerifyEmail />
    </AuthProvider>
  )
}
