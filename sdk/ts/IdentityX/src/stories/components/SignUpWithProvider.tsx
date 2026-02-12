import { AuthProvider, SignUp } from "../../react";

export default function SignUpWithProvider({ flow_id }: { flow_id: string }) {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <SignUp flow_id={flow_id} />
    </AuthProvider>
  )
}