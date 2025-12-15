import { AuthProvider, BasicLogoutButton } from "../../react";

export default function BasicLogoutWithProvider() {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <BasicLogoutButton />
    </AuthProvider>
  )
}