import { createFileRoute } from "@tanstack/solid-router";
import { createEffect } from "solid-js";
import type { JSX } from "@solidjs/web";
import Loader2Icon from "~icons/lucide/loader-2";
import { useCompleteEventSellerMutation } from "@/features/payments/api/mutations";
import { toast } from "@/shared/ui/toast";
import z from "zod";

const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;

const searchSchema = z.object({
  credential_id: z.uuid(),
  public_key: z.string().min(1),
});

export const Route = createFileRoute("/events/$eventId/payssage/oauth/callback")({
  validateSearch: searchSchema,
  component: SellerCallback,
});

function SellerCallback() {
  const params = Route.useParams();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const mutation = useCompleteEventSellerMutation();
  let started = false;

  createEffect(
    () => ({ params: params(), search: search() }),
    ({ params: currentParams, search: currentSearch }) => {
      const { eventId } = currentParams;
      const { credential_id, public_key } = currentSearch;
      if (started) return;
      started = true;

      void mutation.mutateAsync({ eventId, sellerId: credential_id, publicKey: public_key })
        .then(() => {
          toast.success("Mercado Pago conectado");
          void navigate({ to: "/admin/events/$eventId", params: { eventId } });
        })
        .catch(() => {
          toast.error("Não foi possível concluir a conexão do Mercado Pago");
          void navigate({ to: "/admin/events/$eventId", params: { eventId } });
        });
    },
  );

  return (
    <div class="flex min-h-100 flex-col items-center justify-center gap-4">
      <Loader2 class="size-10 animate-spin text-primary" />
      <p class="text-sm text-muted-foreground">Conectando Mercado Pago…</p>
    </div>
  );
}
