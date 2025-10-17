import { AuthProvider, SignUp } from "../next";

export default function SignUpWithProvider() {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <SignUp />
    </AuthProvider>
  )
}