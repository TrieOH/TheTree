import { type MouseEvent, useEffect, useState } from "react";
import SessionCard from "./SessionCard";
import { useAuth } from "../../AuthProvider";
import type { SessionI } from "../../../types/sessions-types";

export interface SessionsProps {
  /** If true will revoke even the current session */
  revokeAll?: boolean;
  /** What will happen when sessions are revoked */
  onSuccess?: (message?: string) => Promise<void>;
}

export function Sessions({
  revokeAll = false,
  onSuccess,
}: SessionsProps) {
  const { auth } = useAuth();
  const [loading, setLoading] = useState(false);
  const [sessions, setSessions] = useState<SessionI[]>([]);

  const fetchSessions = async () => {
    const res = await auth.sessions();
    if (res.success) {
      const currentSessionId = auth.profile()?.session_id;

      const sessions = (res.data).sort((a, b) => {
        if (a.session_id === currentSessionId) return -1;
        if (b.session_id === currentSessionId) return 1;
        return 0;
      });
      setSessions(sessions);
    }
  }

  useEffect(() => {
    fetchSessions();
  }, []);

  const handleRevokeASession = async (e: MouseEvent<SVGElement>, id: string) => {
    e.preventDefault();
    const sessionToRemove = sessions.find(s => s.session_id === id);
    if (!sessionToRemove) return;
    setSessions(sessions.filter(s => s.session_id !== id))
    try {
      const res = await auth.revokeASession(id);
      if (!res.success) setSessions(prev => [...prev, sessionToRemove]);
    } catch {
      setSessions(prev => [...prev, sessionToRemove]);
    }
  }

  const handleRevokeSessions = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    if (loading) return;
    setLoading(true);

    const res = await auth.revokeSessions(revokeAll);
    const id = auth.profile()?.session_id;
    if (res.success) {
      setSessions(revokeAll ? [] : sessions.filter(s => s.session_id === id));
      if (onSuccess) onSuccess(res.message);
    }
    setLoading(false);
  }

  return (
    <div className="font-sans w-full min-w-[20rem] m-2 @container bg-background text-foreground p-[1.5rem_0.5rem] rounded-lg">
      <div className="w-full flex flex-wrap items-center justify-center text-center gap-4 box-border px-3 @[640px]:justify-between @[640px]:text-left @[640px]:gap-8">
        <div className="flex-[0_1_auto] max-w-100">
          <h3 className="font-semibold text-2xl m-0 mb-1">Navegadores e Dispositivos</h3>
          <p className="font-extralight text-base m-0">Esses navegadores e dispositivos estão atualmente conectados à sua conta. Remova quaisquer dispositivos não autorizados.</p>
        </div>
        <button
          type="submit"
          onClick={handleRevokeSessions}
          disabled={loading}
          className={`font-sans w-full max-w-56 p-[1rem_0] h-auto text-base font-semibold outline-none bg-transparent relative overflow-hidden shrink-0 border-2 border-foreground text-foreground cursor-pointer transition-transform duration-500 rounded-lg hover:scale-[1.02] active:scale-[0.99] disabled:opacity-60 disabled:cursor-not-allowed disabled:transform-none! ${loading ? "button-loading" : ""
            }`}
        >
          Revogar todas as sessões
        </button>
      </div>
      <div className="mt-4">
        {sessions.length > 0 ? sessions.map(s => (
          <SessionCard
            key={s.session_id}
            {...s}
            is_current={auth.profile()?.session_id === s.session_id}
            onClick={handleRevokeASession}
          />
        )) : <span className="block border-t border-border p-[1.25rem_0.75rem] text-center font-semibold">Nenhuma Sessão Disponível</span>}
      </div>
    </div>
  )
}