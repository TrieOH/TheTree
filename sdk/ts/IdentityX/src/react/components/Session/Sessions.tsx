import { type MouseEvent, useEffect, useState } from "react";
import SessionCard from "./SessionCard";
import { useAuth } from "../../AuthProvider";
import type { SessionI } from "../../../types/sessions-types";


export function Sessions() {
  const { auth } = useAuth();
  const [loading, setLoading] = useState(false);
  const [sessions, setSessions] = useState<SessionI[]>([]);

  const fetchSessions = async () => {
    const res = await auth.sessions();
    setSessions(res.data || []);
  }

  useEffect(() => {
    fetchSessions();
    // auth.sessions();
    // console.log(auth.profile())
  }, [])
  

  // const handleRevokeAllSessions = async (e: MouseEvent<HTMLButtonElement>) => {
  //   e.preventDefault();
  //   if (loading) return;

  //   setLoading(true);

  //   // const res = await auth.logout();
  //   // if(res.code === 200 && onSuccess) await onSuccess();
  //   // else if(onFailed) await onFailed(res.message);
  //   setLoading(false);
  // }

  return (
    <div className="trieoh trieoh-sessions">
      <div className="trieoh-sessions__header">
        <div>
          <h3>Navegadores e Dispositivos</h3>
          <p>Esses navegadores e dispositivos estão atualmente conectados à sua conta. Remova quaisquer dispositivos não autorizados.</p>
        </div>
        <button 
          type="submit"
          onClick={() => null}
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
          />
        ))}
      </div>
    </div>
  )
}