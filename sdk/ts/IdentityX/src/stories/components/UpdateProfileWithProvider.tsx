import { AuthProvider, UpdateProfile } from "../../react";
import type { UpdateProfileProps } from "../../react/components/Profile/UpdateProfile";

export default function UpdateProfileWithProvider(props: UpdateProfileProps) {
  return (
    <AuthProvider baseURL="http://localhost:8080">
      <UpdateProfile {...props} />
    </AuthProvider>
  )
}
