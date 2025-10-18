import type { MouseEvent } from "react";
import { useAuth } from "../../AuthProvider";
import LogoutIcon from "../Icons/LogoutIcon";

export interface LogoutProps {
  onSuccess?: () => void;
  onFailed?: (message: string) => void;
}

export function BasicLogoutButton({
  onSuccess,
  onFailed
}: LogoutProps) {
  const { auth } = useAuth();
  const handleLogout = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();

    const res = await auth.logout();
    if(res.code === 200 && onSuccess) onSuccess();
    else if(onFailed) onFailed(res.message);
  }
  return (
    <button
      onClick={handleLogout}
      type="button"
      className="trieoh trieoh-button trieoh-button--all-rounded trieoh-button--logout"
    >
      <LogoutIcon /> <span>Logout</span>
    </button>
  )
}