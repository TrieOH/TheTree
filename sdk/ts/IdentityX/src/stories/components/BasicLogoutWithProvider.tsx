import { AuthProvider, BasicLogoutButton, SignIn } from "../../next";

export default function BasicLogoutWithProvider() {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <BasicLogoutButton />
    </AuthProvider>
  )
}