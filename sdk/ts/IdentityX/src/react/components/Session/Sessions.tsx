import { type MouseEvent, useEffect, useState } from "react";
import SessionCard from "./SessionCard";
import { useAuth } from "../../AuthProvider";
import type { SessionI } from "../../../types/sessions-types";

export interface SessionsProps {
  /** If true will revoke even the current session */
  revokeAll?: boolean;
  /** What will happen when sessions are revoked */
  onSuccess?: () => Promise<void>;
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
    setSessions(res.data || []);
  }

  useEffect(() => {
    fetchSessions();
  }, []);

  const handleRevokeASession = async (e: MouseEvent<SVGElement>, id: string) => {
    e.preventDefault();
    const sessionToRemove = sessions.find(s => s.session_id === id);
    if(!sessionToRemove) return;
    setSessions(sessions.filter(s => s.session_id !== id))
    try {
      const res = await auth.revokeASession(id);
      if(res.code !== 200) setSessions(prev => [...prev, sessionToRemove]);
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
    if(res.code === 200) {
      setSessions(revokeAll ? [] : sessions.filter(s => s.session_id === id));
      if(onSuccess) onSuccess();
    }
    setLoading(false);
  }

  return (
    <div className="trieoh trieoh-sessions">
      <div className="trieoh-sessions__header">
        <div>
          <h3>Navegadores e Dispositivos</h3>
          <p>Esses navegadores e dispositivos estão atualmente conectados à sua conta. Remova quaisquer dispositivos não autorizados.</p>
        </div>
        <button 
          type="submit"
          onClick={handleRevokeSessions}
          disabled={loading}
          className={
            `trieoh trieoh-button trieoh-button--all-rounded 
            ${loading ? "trieoh-button--loading" : ""}`
          }
        >
          Revogar todas as sessões
        </button>
      </div>
      <div className="trieoh-sessions__content">
        {sessions.map(s => (
          <SessionCard 
            key={s.session_id} 
            {...s} 
            is_current={auth.profile()?.session_id === s.session_id}
            onClick={handleRevokeASession}
          />
        ))}
      </div>
    </div>
  )
}