import { useState, type MouseEvent } from "react";
import { useAuth } from "../../AuthProvider";
import { ImExit } from "react-icons/im";

export interface LogoutProps {
  onSuccess?: () => Promise<void>;
  onFailed?: (message: string) => Promise<void>;
}

export function BasicLogoutButton({
  onSuccess,
  onFailed
}: LogoutProps) {
  const { auth } = useAuth();
  const [loading, setLoading] = useState(false);

  const handleLogout = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    if (loading) return;

    setLoading(true);

    const res = await auth.logout();
    if(res.code === 200 && onSuccess) await onSuccess();
    else if(onFailed) await onFailed(res.message);
    setLoading(false);
  }
  return (
    <button
      onClick={handleLogout}
      type="button"
      disabled={loading}
      className={"trieoh trieoh-button--logout"}
    >
      <ImExit size={32}/> <span>Log out</span>
    </button>
  )
}