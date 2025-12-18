import { AuthProvider, SignIn } from "../../react";

export default function SignInWithProvider() {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <SignIn />
    </AuthProvider>
  )
}