import { AuthProvider, SignIn } from "../next";

export default function SignInWithProvider() {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <SignIn />
    </AuthProvider>
  )
}