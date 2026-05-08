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
      className="font-inter border-none bg-transparent cursor-pointer flex items-end gap-1 text-trieoh-base font-medium text-[oklch(0.628_0.2577_29.23)] transition-transform duration-200 hover:scale-[1.05] active:scale-[0.98] disabled:opacity-60 disabled:cursor-not-allowed disabled:!transform-none"
    >
      <ImExit size={24} /> <span>Log out</span>
    </button>
  )
  }