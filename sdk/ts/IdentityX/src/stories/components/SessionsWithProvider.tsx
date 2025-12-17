import { AuthProvider, Sessions } from "../../react";

export default function SessionsWithProvider() {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <Sessions />
    </AuthProvider>
  )
}