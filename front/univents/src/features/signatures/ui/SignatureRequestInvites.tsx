import { useQuery } from "@tanstack/react-query";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { Mail, Plus } from "lucide-react";
import { useMemo, useState } from "react";
import { allSignatureRequestsQueryOptions } from "@/features/signatures/api";
import { useCancelSignatureRequestMutation } from "@/features/signatures/api/mutations";
import type { SignatureRequestI } from "@/features/signatures/model";
import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Button } from "@/shared/ui/shadcn/button";
import { CreateSignatureRequestModal } from "./CreateSignatureRequestModal";

export function SignatureRequestInvites({ editionId }: { editionId: string }) {
  const requestsQuery = useQuery(allSignatureRequestsQueryOptions(editionId));
  const cancelRequest = useCancelSignatureRequestMutation();
  const [createOpen, setCreateOpen] = useState(false);
  const [filter, setFilter] = useState("");

  const filteredRequests = useMemo(() => {
    const search = filter.trim().toLowerCase();
    if (!search) return requestsQuery.data ?? [];
    return (requestsQuery.data ?? []).filter((request) =>
      [
        request.signatory_name,
        request.signatory_email ?? "",
        request.signatory_title ?? "",
        request.status,
      ].some((value) => value.toLowerCase().includes(search)),
    );
  }, [filter, requestsQuery.data]);

  return (
    <>
      <PaginatedContainer<SignatureRequestI>
        items={filteredRequests}
        layout="grid"
        minItemWidth="16rem"
        pageSize={6}
        gap="6"
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome, e-mail ou status..."
        sortFields={[
          { key: "signatory_name", label: "Signatário" },
          { key: "status", label: "Status" },
          { key: "expires_at", label: "Validade" },
          { key: "created_at", label: "Data de criação" },
        ]}
        itemLabel="convites"
        headerActions={
          <Button onClick={() => setCreateOpen(true)} className="h-9 shrink-0">
            <Plus className="size-4" /> Novo convite
          </Button>
        }
        emptyState={
          <EmptyState
            icon={Mail}
            eyebrow="Convites"
            title="Nenhum convite enviado"
            description="Envie um convite para alguém assinar esta edição."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(requests) =>
          requests.map((request) => (
            <article
              key={request.id}
              className="group relative flex min-h-48 w-full min-w-0 flex-col overflow-hidden rounded-2xl bg-card text-left ring-1 ring-foreground/10 shadow-xs transition-all duration-300 hover:-translate-y-0.5 hover:ring-foreground/20 hover:shadow-sm"
            >
              <div className="relative flex-1 overflow-hidden bg-linear-to-br from-primary/10 via-muted to-background p-5">
                <div className="absolute inset-0 bg-linear-to-t from-background/90 via-background/20 to-transparent" />
                <div className="relative z-10 space-y-2">
                  <div className="flex items-center justify-between gap-2">
                    <Badge variant="outline" className="text-muted-foreground">
                      Convite
                    </Badge>
                    <Badge
                      variant="outline"
                      className={cn(
                        request.status === "pending" &&
                          "border-primary/30 bg-primary/5 text-primary",
                        request.status === "cancelled" &&
                          "border-destructive/30 bg-destructive/5 text-destructive",
                      )}
                    >
                      {request.status === "pending"
                        ? "Pendente"
                        : request.status === "completed"
                          ? "Concluído"
                          : request.status === "cancelled"
                            ? "Cancelado"
                            : "Expirado"}
                    </Badge>
                  </div>
                  <h3 className="line-clamp-2 pt-2 text-lg font-semibold transition-colors group-hover:text-primary">
                    {request.signatory_name}
                  </h3>
                  <p className="break-all text-muted-foreground">
                    {request.signatory_email ?? "—"}
                  </p>
                  {request.signatory_title && (
                    <p className="text-xs text-muted-foreground">
                      {request.signatory_title}
                    </p>
                  )}
                  <div className="border-t border-border/70 pt-3 text-xs text-muted-foreground">
                    <p className="text-xs text-muted-foreground">
                      Expira em{" "}
                      {new Date(request.expires_at).toLocaleDateString("pt-BR")}
                    </p>
                  </div>
                </div>
              </div>
              <div className="flex min-h-16 items-center justify-between gap-3 border-t border-border bg-card px-5 py-3">
                {request.status === "pending" && (
                  <Button
                    size="sm"
                    variant="outline"
                    disabled={cancelRequest.isPending}
                    onClick={() =>
                      cancelRequest.mutate({
                        editionId,
                        requestId: request.id,
                      })
                    }
                  >
                    Cancelar
                  </Button>
                )}
                {request.status !== "pending" && <span aria-hidden="true" />}
              </div>
            </article>
          ))
        }
      />
      <CreateSignatureRequestModal
        editionId={editionId}
        open={createOpen}
        onOpenChange={setCreateOpen}
      />
    </>
  );
}
