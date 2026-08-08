import { createFileRoute } from "@tanstack/react-router";
import { Loader2 } from "lucide-react";
import { useEffect, useRef } from "react";
import { toast } from "sonner";
import z from "zod";
import { useCompleteEventSellerMutation } from "@/features/payments/api/mutations";

const searchSchema = z.object({
  credential_id: z.string().uuid(),
  public_key: z.string().min(1),
});

export const Route = createFileRoute(
  "/events/$eventId/payssage/oauth/callback",
)({
  validateSearch: searchSchema,
  component: SellerCallback,
});

function SellerCallback() {
  const { eventId } = Route.useParams();
  const { credential_id, public_key } = Route.useSearch();
  const navigate = Route.useNavigate();
  const mutation = useCompleteEventSellerMutation();
  const started = useRef(false);

  useEffect(() => {
    if (started.current) return;
    started.current = true;

    mutation.mutate(
      { eventId, sellerId: credential_id, publicKey: public_key },
      {
        onSuccess: () => {
          toast.success("Mercado Pago conectado");
          void navigate({ to: "/admin/events/$eventId", params: { eventId } });
        },
        onError: () => {
          toast.error("Não foi possível concluir a conexão do Mercado Pago");
          void navigate({ to: "/admin/events/$eventId", params: { eventId } });
        },
      },
    );
  }, [credential_id, eventId, mutation, navigate, public_key]);

  return (
    <div className="flex min-h-100 flex-col items-center justify-center gap-4">
      <Loader2 className="size-10 animate-spin text-primary" />
      <p className="text-sm text-muted-foreground">Conectando Mercado Pago…</p>
    </div>
  );
}
