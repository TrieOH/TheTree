import { AuthProvider, BasicLogoutButton } from "../../react";

interface BasicLogoutWithProviderProps {
  isProjectMode?: boolean;
}

export default function BasicLogoutWithProvider({ isProjectMode = true }: BasicLogoutWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080" isProjectMode={isProjectMode}>
      <BasicLogoutButton />
    </AuthProvider>
  )
}