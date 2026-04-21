import { AuthProvider, Sessions } from "../../react";

interface SessionsWithProviderProps {
  isProjectMode?: boolean;
}

export default function SessionsWithProvider({ isProjectMode = true }: SessionsWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080" isProjectMode={isProjectMode}>
      <Sessions />
    </AuthProvider>
  )
}