import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { PenLine, Plus, Trash2 } from "lucide-react";
import { useMemo, useState } from "react";
import { allSignaturesQueryOptions } from "@/features/signatures/api";
import { useRemoveSignatureMutation } from "@/features/signatures/api/mutations";
import type { SignatureI } from "@/features/signatures/model";
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
    allSignaturesQueryOptions(editionId),
  );
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
    <>
      <PaginatedContainer<SignatureI>
        items={filteredSignatures}
        layout="grid"
        minItemWidth="16rem"
        pageSize={6}
        gap="6"
        sortFields={[
          { key: "signatory_name", label: "Signatário" },
          { key: "created_at", label: "Data de criação" },
        ]}
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
              <div className="relative aspect-video overflow-hidden bg-white">
                <img
                  src={signature.image_url}
                  alt={signature.signatory_name}
                  className={cn(
                    "h-full w-full object-contain bg-white transition-transform duration-700 ease-out",
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
            editionId,
            signatureId: removingSignature.id,
          });
          setRemovingSignature(null);
        }}
      />
    </>
  );
}
