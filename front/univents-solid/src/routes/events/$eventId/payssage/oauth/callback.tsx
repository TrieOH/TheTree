import { createFileRoute } from "@tanstack/solid-router";
import { createEffect } from "solid-js";
import type { JSX } from "@solidjs/web";
import Loader2Icon from "~icons/lucide/loader-2";
import { useCompleteEventSellerMutation } from "@/features/payments/api/mutations";
import { toast } from "@/shared/ui/toast";
import { Schema, createSchema, uuid } from "@/shared/lib/effect-schema";

const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;

const searchSchema = createSchema(
  Schema.Struct({
    credential_id: Schema.String.pipe(uuid()),
    public_key: Schema.NonEmptyString,
  }),
);

export const Route = createFileRoute("/events/$eventId/payssage/oauth/callback")({
  validateSearch: searchSchema,
  head: () => ({
    meta: [{ title: "Conectando Mercado Pago - Univents" }],
  }),
  component: SellerCallback,
});

function SellerCallback() {
  const params = Route.useParams();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const mutation = useCompleteEventSellerMutation();
  let started = false;

  createEffect(
    () => ({ currentParams: params(), currentSearch: search() }),
    ({ currentParams, currentSearch }) => {
      const { eventId } = currentParams;
      const { credential_id, public_key } = currentSearch;
      if (started) return;
      started = true;

      void mutation.mutateAsync({ eventId, sellerId: credential_id, publicKey: public_key })
        .then(() => {
          toast.success("Mercado Pago conectado");
          void navigate({ to: "/admin/events/$eventId", params: { eventId } });
        })
        .catch((error: unknown) => {
          const message = error instanceof Error ? error.message : "Falha ao vincular vendedor";
          toast.error(message);
          void navigate({ to: "/admin/events/$eventId", params: { eventId } });
        });
    },
  );

  return (
    <div class="min-h-screen bg-background text-foreground flex items-center justify-center p-4">
      <div class="flex items-center gap-3 text-muted-foreground text-sm">
        <Loader2 class="size-5 animate-spin text-primary" />
        <span>Concluindo conexão com Mercado Pago...</span>
      </div>
    </div>
  );
}
