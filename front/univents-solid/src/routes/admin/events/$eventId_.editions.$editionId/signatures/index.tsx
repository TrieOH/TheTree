import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { For, Show, createMemo, createSignal } from "solid-js";

import MailIcon from "~icons/lucide/mail";
import PenLineIcon from "~icons/lucide/pen-line";
import PlusIcon from "~icons/lucide/plus";

import { useQuery } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { Button, EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { toast } from "@/shared/ui/toast";

import {
  allSignatureRequestsQueryOptions,
  allSignaturesQueryOptions,
} from "@/features/signatures/api";
import {
  useCancelSignatureRequestMutation,
  useDeleteSignatureMutation,
} from "@/features/signatures/api/mutations";
import type {
  SignatureI,
  SignatureRequestI,
} from "@/features/signatures/model";
import {
  AdminCreateSignatureCard,
  AdminCreateSignatureRequestCard,
  AdminSignatureCard,
  AdminSignatureRequestCard,
  CreateSignatureModal,
  CreateSignatureRequestModal,
  type SignatureSection,
  SignatureSectionTabs,
} from "@/features/signatures/ui";

const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;
const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/signatures/",
)({
  head: () => ({
    meta: [{ title: "Assinaturas - Admin Univents" }],
  }),
  component: AdminEditionSignaturesRoute,
});

function AdminEditionSignaturesRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;
  const editionId = () => params().editionId;

  const [activeSection, setActiveSection] =
    createSignal<SignatureSection>("signatures");
  const [filter, setFilter] = createSignal("");
  const [signatureSort, setSignatureSort] = createSignal<SortState<SignatureI>>({
    field: "signatory_name",
    direction: "asc",
  });
  const [requestSort, setRequestSort] = createSignal<
    SortState<SignatureRequestI>
  >({
    field: "signatory_name",
    direction: "asc",
  });

  const [createSignatureOpen, setCreateSignatureOpen] = createSignal(false);
  const [createRequestOpen, setCreateRequestOpen] = createSignal(false);
  const [deletingId, setDeletingId] = createSignal<string | null>(null);
  const [cancellingId, setCancellingId] = createSignal<string | null>(null);

  // Queries
  const signaturesQuery = useQuery(() => allSignaturesQueryOptions(editionId()));
  const requestsQuery = useQuery(() =>
    allSignatureRequestsQueryOptions(editionId()),
  );

  const signatures = () => (signaturesQuery().data ?? []) as SignatureI[];
  const requests = () => (requestsQuery().data ?? []) as SignatureRequestI[];

  // Filtered Signatures
  const filteredSignatures = createMemo(() => {
    const search = filter().trim().toLowerCase();
    if (!search) return signatures();

    return signatures().filter((sig) =>
      [sig.signatory_name, sig.signatory_title ?? "", sig.signatory_email ?? ""].some(
        (val) => val.toLowerCase().includes(search),
      ),
    );
  });

  // Filtered Requests
  const filteredRequests = createMemo(() => {
    const search = filter().trim().toLowerCase();
    if (!search) return requests();

    return requests().filter((req) =>
      [
        req.signatory_name,
        req.signatory_title ?? "",
        req.signatory_email ?? "",
        req.status,
      ].some((val) => val.toLowerCase().includes(search)),
    );
  });

  // Mutations
  const deleteSignatureMutation = useDeleteSignatureMutation();
  const cancelRequestMutation = useCancelSignatureRequestMutation();

  const handleDeleteSignature = async (sig: SignatureI) => {
    setDeletingId(sig.id);
    try {
      await deleteSignatureMutation.mutateAsync({
        editionId: editionId(),
        signatureId: sig.id,
      });
      toast.success("Assinatura removida com sucesso.");
    } catch {
      toast.error("Erro ao remover assinatura.");
    } finally {
      setDeletingId(null);
    }
  };

  const handleCancelRequest = async (req: SignatureRequestI, reason?: string) => {
    setCancellingId(req.id);
    try {
      await cancelRequestMutation.mutateAsync({
        editionId: editionId(),
        requestId: req.id,
        reason,
      });
      toast.success("Convite cancelado com sucesso.");
    } catch {
      toast.error("Erro ao cancelar convite.");
    } finally {
      setCancellingId(null);
    }
  };

  return (
    <div class="mx-auto max-w-7xl space-y-6 min-w-0">
      {/* Header with Title & Description */}
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-foreground">
            Assinaturas da Edição
          </h1>
          <p class="text-xs sm:text-sm text-muted-foreground">
            Cadastre assinaturas digitais e envie convites para emissão e autenticação de certificados.
          </p>
        </div>
      </div>

      {/* Navigation Tabs */}
      <SignatureSectionTabs
        active={activeSection()}
        onChange={(sec) => {
          setFilter("");
          setActiveSection(sec);
        }}
        signaturesCount={signatures().length}
        invitesCount={requests().length}
      />

      {/* Content: Approved Signatures */}
      <Show when={activeSection() === "signatures"}>
        <PaginatedContainer<SignatureI>
          items={filteredSignatures()}
          layout="grid"
          minItemWidth="16rem"
          maxRows={(columns) => (columns === 1 ? 8 : 4)}
          gap="3"
          sort={signatureSort()}
          onSortChange={setSignatureSort}
          sortFields={[
            {
              key: "signatory_name",
              label: "Nome do signatário",
              ascLabel: "A → Z",
              descLabel: "Z → A",
            },
            {
              key: "created_at",
              label: "Data de cadastro",
              ascLabel: "Mais antigas primeiro",
              descLabel: "Mais recentes primeiro",
              comparator: (a, b) =>
                new Date(a.created_at).getTime() -
                new Date(b.created_at).getTime(),
            },
          ]}
          filterValue={filter()}
          onFilterChange={setFilter}
          filterPlaceholder="Buscar por nome, cargo ou e-mail..."
          itemLabel="assinaturas"
          emptyState={
            <EmptyState
              class="border-0 bg-transparent px-0 py-8 shadow-none"
              icon={<PenLine class="size-6 text-foreground/70" />}
              title="Nenhuma assinatura cadastrada"
              description={
                filter()
                  ? "Nenhuma assinatura corresponde à busca informada."
                  : "Cadastre a primeira assinatura ou importe uma imagem para autenticar certificados desta edição."
              }
              action={
                <Button
                  onClick={() => setCreateSignatureOpen(true)}
                  class="gap-1.5 text-xs cursor-pointer"
                >
                  <Plus class="size-4" />
                  <span>Cadastrar primeira assinatura</span>
                </Button>
              }
            />
          }
          renderItems={(slice, options) => (
            <>
              <AdminCreateSignatureCard
                index={0}
                animate={options.animate}
                onCreate={() => setCreateSignatureOpen(true)}
              />
              <For each={slice}>
                {(sig) => (
                  <AdminSignatureCard
                    signature={sig}
                    isDeleting={deletingId() === sig.id}
                    onDelete={() => handleDeleteSignature(sig)}
                  />
                )}
              </For>
            </>
          )}
        />
      </Show>

      {/* Content: Signature Request Invites */}
      <Show when={activeSection() === "invites"}>
        <PaginatedContainer<SignatureRequestI>
          items={filteredRequests()}
          layout="grid"
          minItemWidth="18rem"
          maxRows={(columns) => (columns === 1 ? 8 : 4)}
          gap="3"
          sort={requestSort()}
          onSortChange={setRequestSort}
          sortFields={[
            {
              key: "signatory_name",
              label: "Signatário",
              ascLabel: "A → Z",
              descLabel: "Z → A",
            },
            {
              key: "status",
              label: "Status",
              ascLabel: "Status (A → Z)",
              descLabel: "Status (Z → A)",
            },
            {
              key: "created_at",
              label: "Data de envio",
              ascLabel: "Mais antigos primeiro",
              descLabel: "Mais recentes primeiro",
              comparator: (a, b) =>
                new Date(a.created_at).getTime() -
                new Date(b.created_at).getTime(),
            },
          ]}
          filterValue={filter()}
          onFilterChange={setFilter}
          filterPlaceholder="Buscar por nome, e-mail ou status..."
          itemLabel="convites"
          emptyState={
            <EmptyState
              class="border-0 bg-transparent px-0 py-8 shadow-none"
              icon={<Mail class="size-6 text-foreground/70" />}
              title="Nenhum convite enviado"
              description={
                filter()
                  ? "Nenhum convite corresponde à busca informada."
                  : "Envie um convite por e-mail para que um palestrante, diretor ou coordenador assine digitalmente."
              }
              action={
                <Button
                  onClick={() => setCreateRequestOpen(true)}
                  class="gap-1.5 text-xs cursor-pointer"
                >
                  <Mail class="size-4" />
                  <span>Enviar primeiro convite</span>
                </Button>
              }
            />
          }
          renderItems={(slice, options) => (
            <>
              <AdminCreateSignatureRequestCard
                index={0}
                animate={options.animate}
                onCreate={() => setCreateRequestOpen(true)}
              />
              <For each={slice}>
                {(req) => (
                  <AdminSignatureRequestCard
                    request={req}
                    isCancelling={cancellingId() === req.id}
                    onCancel={(reason) => handleCancelRequest(req, reason)}
                  />
                )}
              </For>
            </>
          )}
        />
      </Show>

      {/* Modal: New Direct Signature */}
      <CreateSignatureModal
        open={createSignatureOpen()}
        onOpenChange={setCreateSignatureOpen}
        eventId={eventId()}
        editionId={editionId()}
      />

      {/* Modal: New Signature Request Invite */}
      <CreateSignatureRequestModal
        open={createRequestOpen()}
        onOpenChange={setCreateRequestOpen}
        editionId={editionId()}
      />
    </div>
  );
}
