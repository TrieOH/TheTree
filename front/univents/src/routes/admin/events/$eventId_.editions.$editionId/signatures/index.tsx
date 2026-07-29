import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { Mail, PenLine, Plus, Trash2 } from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import {
  allSignatureRequestsQueryOptions,
  allSignaturesQueryOptions,
  cancelSignatureRequestFn,
  createSignatureRequestFn,
} from "@/features/signatures/api";
import { useRemoveSignatureMutation } from "@/features/signatures/api/mutations";
import type {
  SignatureI,
  SignatureRequestI,
} from "@/features/signatures/model";
import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Button } from "@/shared/ui/shadcn/button";
import { AlertModal } from "@/widgets/ui/alert-modal";

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/signatures/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId, editionId } = Route.useParams();
  const { data: signatures = [] } = useQuery(
    allSignaturesQueryOptions(eventId, editionId),
  );
  const requestsQuery = useQuery(allSignatureRequestsQueryOptions(editionId));
  const [requestName, setRequestName] = useState("");
  const [requestEmail, setRequestEmail] = useState("");
  const [requestBusy, setRequestBusy] = useState(false);
  const removeSignatureMutation = useRemoveSignatureMutation();
  const [filter, setFilter] = useState("");
  const [removingSignature, setRemovingSignature] = useState<SignatureI | null>(
    null,
  );

  const filteredSignatures = useMemo(() => {
    const search = filter.trim().toLowerCase();
    if (!search) return signatures;

    return signatures.filter((signature) =>
      [signature.signatory_name, signature.image_url].some((value) =>
        value.toLowerCase().includes(search),
      ),
    );
  }, [filter, signatures]);

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<SignatureI>
        items={filteredSignatures}
        layout="grid"
        minItemWidth="16rem"
        pageSize={6}
        gap="6"
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por título ou URL..."
        itemLabel="assinaturas"
        headerActions={
          <Link
            to="/admin/events/$eventId/editions/$editionId/signatures/editor"
            params={{ eventId, editionId }}
            className={cn(
              "inline-flex h-9 items-center justify-center gap-2 rounded-lg px-4 text-sm font-medium",
              "bg-primary text-primary-foreground shadow-sm transition-colors hover:bg-primary/90",
              "sm:min-w-44 sm:px-5",
            )}
          >
            <Plus className="size-4 shrink-0" />
            <span className="whitespace-nowrap">Nova assinatura</span>
          </Link>
        }
        emptyState={
          <EmptyState
            icon={PenLine}
            eyebrow="Assinaturas"
            title="Nenhuma assinatura encontrada"
            description="Crie a primeira assinatura para usar nas certificações dessa edição."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((signature, index) => (
            <article
              key={signature.id}
              className={cn(
                "group relative flex w-full min-w-0 flex-col overflow-hidden rounded-2xl bg-card text-left",
                "ring-1 ring-foreground/10 shadow-xs",
                "transform-gpu will-change-transform",
                "transition-all duration-300 ease-out",
                "hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm",
              )}
            >
              <div className="relative aspect-video overflow-hidden bg-muted">
                <img
                  src={signature.image_url}
                  alt={signature.signatory_name}
                  className={cn(
                    "h-full w-full object-contain bg-background transition-transform duration-700 ease-out",
                    "group-hover:scale-[1.03]",
                  )}
                  loading={index < 4 ? "eager" : "lazy"}
                />

                <div className="absolute inset-0 bg-linear-to-t from-background/85 via-background/20 to-transparent" />

                <div className="absolute left-3 top-3">
                  <Badge
                    variant="outline"
                    className="bg-background/80 backdrop-blur-sm"
                  >
                    Assinatura
                  </Badge>
                </div>

                <div className="absolute inset-x-0 bottom-0 flex items-end justify-between gap-3 p-4 sm:p-5">
                  <div className="min-w-0 space-y-1">
                    <h3 className="line-clamp-2 text-balance text-lg font-semibold leading-snug text-foreground transition-colors duration-300 group-hover:text-primary sm:text-xl">
                      {signature.signatory_name}
                    </h3>
                  </div>

                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="size-9 shrink-0 bg-background/85 backdrop-blur-sm"
                    onClick={() => setRemovingSignature(signature)}
                  >
                    <Trash2 className="size-4 text-destructive" />
                  </Button>
                </div>
              </div>
            </article>
          ))
        }
      />
      <section className="mt-8 w-full rounded-2xl border bg-card p-5 shadow-xs">
        <div className="mb-4 flex items-start justify-between gap-3">
          <div>
            <h2 className="flex items-center gap-2 text-lg font-semibold">
              <Mail className="size-5" /> Solicitações de assinatura
            </h2>
            <p className="text-sm text-muted-foreground">
              Envie um convite para alguém assinar esta edição.
            </p>
          </div>
        </div>
        <form
          className="grid gap-3 sm:grid-cols-[1fr_1fr_auto]"
          onSubmit={async (event) => {
            event.preventDefault();
            setRequestBusy(true);
            try {
              const response = await createSignatureRequestFn(editionId, {
                signatory_name: requestName,
                signatory_email: requestEmail,
              });
              if (!response.success)
                throw new Error(
                  response.message || "Não foi possível criar a solicitação",
                );
              setRequestName("");
              setRequestEmail("");
              await requestsQuery.refetch();
            } catch (error) {
              toast.error(
                error instanceof Error
                  ? error.message
                  : "Erro ao criar solicitação",
              );
            } finally {
              setRequestBusy(false);
            }
          }}
        >
          <input
            className="h-9 rounded-md border bg-background px-3 text-sm"
            placeholder="Nome do signatário"
            value={requestName}
            onChange={(event) => setRequestName(event.target.value)}
            required
          />
          <input
            className="h-9 rounded-md border bg-background px-3 text-sm"
            type="email"
            placeholder="E-mail"
            value={requestEmail}
            onChange={(event) => setRequestEmail(event.target.value)}
            required
          />
          <Button type="submit" disabled={requestBusy}>
            {requestBusy ? "Enviando..." : "Enviar convite"}
          </Button>
        </form>
        <div className="mt-5 divide-y rounded-lg border">
          {(requestsQuery.data ?? []).map((request: SignatureRequestI) => (
            <div
              key={request.id}
              className="flex flex-wrap items-center justify-between gap-3 px-3 py-3 text-sm"
            >
              <div>
                <p className="font-medium">{request.signatory_name}</p>
                <p className="text-muted-foreground">
                  {request.signatory_email}
                </p>
              </div>
              <div className="flex items-center gap-3">
                <Badge variant="outline">{request.status}</Badge>
                {request.status === "pending" && (
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={async () => {
                      await cancelSignatureRequestFn(request.id);
                      await requestsQuery.refetch();
                    }}
                  >
                    Cancelar
                  </Button>
                )}
              </div>
            </div>
          ))}
          {!requestsQuery.isLoading &&
            (requestsQuery.data ?? []).length === 0 && (
              <p className="px-3 py-4 text-sm text-muted-foreground">
                Nenhuma solicitação enviada.
              </p>
            )}
        </div>
      </section>

      <AlertModal
        open={Boolean(removingSignature)}
        onOpenChange={() => setRemovingSignature(null)}
        title="Remover assinatura?"
        description={
          removingSignature
            ? `A assinatura "${removingSignature.signatory_name}" será removida da biblioteca.`
            : undefined
        }
        confirmLabel="Remover assinatura"
        variant="destructive"
        loading={removeSignatureMutation.isPending}
        onConfirm={async () => {
          if (!removingSignature) return;
          await removeSignatureMutation.mutateAsync({
            eventId,
            editionId,
            signatureId: removingSignature.id,
          });
          setRemovingSignature(null);
        }}
      />
    </div>
  );
}
