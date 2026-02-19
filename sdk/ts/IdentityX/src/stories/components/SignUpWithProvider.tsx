import type { MouseEvent } from "react";
import { AuthProvider, SignUp } from "../../react";
import type { FieldDefinitionResultI } from "../../types/fields-types";

export interface SignUpWithProviderProps {
  loginRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
  flow_id?: string;
  fields?: FieldDefinitionResultI[];
}

export default function SignUpWithProvider({ flow_id, loginRedirect, fields }: SignUpWithProviderProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <SignUp flow_id={flow_id} loginRedirect={loginRedirect} fields={fields} />
    </AuthProvider>
  )
}