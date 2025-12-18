import { AuthProvider, SignUp } from "../../react";

export default function SignUpWithProvider() {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <SignUp />
    </AuthProvider>
  )
}