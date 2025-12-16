import { type MouseEvent, useState } from "react";


export function Session() {
  // const { auth } = useAuth();
  const [loading, setLoading] = useState(false);

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
    <div className="trieoh trieoh-session">
      <div className="trieoh-session__header">
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
    </div>
  )
}