import { AuthProvider, SignUp, TabbedSignUp } from "../../react";

interface FlowIdsI {
  label: string;
  value: string;
}

export default function TabbedSignUpWithProvider({ flowIds }: { flowIds: FlowIdsI[] }) {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <TabbedSignUp flowIds={flowIds} />
    </AuthProvider>
  )
}