import { useState, type MouseEvent } from "react";
import { useAuth } from "../../AuthProvider";
import { ImExit } from "react-icons/im";

export interface LogoutProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  forceLogout?: boolean;
}

export function BasicLogoutButton({
  onSuccess,
  onFailed,
  forceLogout
}: LogoutProps) {
  const { auth } = useAuth();
  const [loading, setLoading] = useState(false);

  const handleLogout = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    if (loading) return;

    setLoading(true);

    const res = await auth.logout({ forceLogout });
    if (res.success) {
      if (onSuccess) await onSuccess(res.message);
    } else if (onFailed) await onFailed(res.message, res.trace);
    setLoading(false);
  }
  return (
    <button
      onClick={handleLogout}
      type="button"
      disabled={loading}
      className={"trieoh trieoh-button--logout"}
    >
      <ImExit size={24} /> <span>Log out</span>
    </button>
  )
}