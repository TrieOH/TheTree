import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { AlertTriangle, Award } from "lucide-react";
import { useMemo, useState } from "react";
import { allAdminEditionsQueryOptions } from "@/features/editions/api";
import { getActorEmailsServerFn } from "@/features/events/server";
import { programsQueryOptions } from "@/features/programs/api";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/shared/ui/shadcn/alert-dialog";
import {
  certificationsByEditionQueryOptions,
  emissionErrorsByEditionQueryOptions,
} from "../api";
import { useInvalidateCertificationMutation } from "../api/mutations";
import type { CertificationEmissionErrorI, CertificationI } from "../model";
import {
  CertificationCard,
  CertificationEmissionErrorCard,
} from "./CertificationCards";

interface CertificationListProps {
  eventId: string;
  editionId: string;
}

export function CertificationList({
  eventId,
  editionId,
}: CertificationListProps) {
  const navigate = useNavigate();
  const { data: editions = [] } = useQuery(
    allAdminEditionsQueryOptions(eventId),
  );
  const { data: programs = [] } = useQuery(programsQueryOptions(editionId));
  const { data: certifications = [] } = useQuery(
    certificationsByEditionQueryOptions(editionId),
  );
  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState<SortState<CertificationI>>({
    field: "issued_at",
    direction: "desc",
  });
  const [certToInvalidate, setCertToInvalidate] =
    useState<CertificationI | null>(null);
  const invalidateMutation = useInvalidateCertificationMutation();
  const { data: actorEmails = {} } = useQuery({
    queryKey: [
      "certification-actor-emails",
      editionId,
      certifications.map((cert) => cert.user_id),
    ],
    queryFn: () =>
      getActorEmailsServerFn({
        data: { actorIds: certifications.map((cert) => cert.user_id) },
      }),
    enabled: certifications.length > 0,
  });
  const programNames = new Map(
    programs.map((program) => [program.id, program.name]),
  );
  const editionName =
    editions.find((edition) => edition.id === editionId)?.name ?? "Edição";
  const visibleCertifications = useMemo(() => {
    const search = filter.trim().toLowerCase();
    return certifications.filter((cert) =>
      [
        actorEmails[cert.user_id],
        cert.verification_hash,
        cert.program_id ? programNames.get(cert.program_id) : editionName,
      ]
        .filter(Boolean)
        .some((value) => value?.toLowerCase().includes(search)),
    );
  }, [actorEmails, certifications, editionName, filter, programNames]);

  return (
    <>
      <PaginatedContainer<CertificationI>
        items={visibleCertifications}
        minItemWidth="16rem"
        gap="4"
        pageSize={8}
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          {
            key: "issued_at",
            label: "Emitido em",
            comparator: (a, b) =>
              new Date(a.issued_at).getTime() - new Date(b.issued_at).getTime(),
          },
          { key: "verification_hash", label: "Hash" },
          { key: "valid", label: "Status" },
          {
            key: "user_id",
            label: "E-mail",
            comparator: (a, b) =>
              (actorEmails[a.user_id] ?? "").localeCompare(
                actorEmails[b.user_id] ?? "",
              ),
          },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        itemLabel="certificados"
        filterPlaceholder="Buscar certificados..."
        emptyState={
          <EmptyState
            icon={Award}
            title="Nenhum certificado emitido"
            description="Os certificados emitidos aparecerão aqui."
          />
        }
        renderItems={(items) =>
          items.map((cert, index) => (
            <CertificationCard
              key={cert.id}
              index={index}
              certification={cert}
              email={actorEmails[cert.user_id]}
              scopeName={
                cert.program_id
                  ? (programNames.get(cert.program_id) ?? "Programação")
                  : editionName
              }
              onView={() =>
                void navigate({
                  to: "/verify/$hash",
                  params: { hash: cert.verification_hash },
                })
              }
              onInvalidate={() => setCertToInvalidate(cert)}
            />
          ))
        }
      />
      <AlertDialog
        open={certToInvalidate !== null}
        onOpenChange={(open) => !open && setCertToInvalidate(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Invalidar certificado?</AlertDialogTitle>
            <AlertDialogDescription>
              Essa ação tornará o certificado inválido.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              onClick={() => {
                if (!certToInvalidate) return;
                invalidateMutation.mutate(
                  {
                    certificationId: certToInvalidate.id,
                    reason: "Invalidado pelo administrador",
                  },
                  { onSuccess: () => setCertToInvalidate(null) },
                );
              }}
            >
              Invalidar
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

export function CertificationEmissionErrorsList({
  eventId,
  editionId,
}: CertificationListProps) {
  const { data: editions = [] } = useQuery(
    allAdminEditionsQueryOptions(eventId),
  );
  const { data: emissionErrors = [] } = useQuery(
    emissionErrorsByEditionQueryOptions(editionId),
  );
  const { data: programs = [] } = useQuery(programsQueryOptions(editionId));
  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState<SortState<CertificationEmissionErrorI>>({
    field: "created_at",
    direction: "desc",
  });
  const { data: actorEmails = {} } = useQuery({
    queryKey: [
      "certification-error-actor-emails",
      editionId,
      emissionErrors.map((error) => error.user_id),
    ],
    queryFn: () =>
      getActorEmailsServerFn({
        data: { actorIds: emissionErrors.map((error) => error.user_id) },
      }),
    enabled: emissionErrors.length > 0,
  });
  const programNames = new Map(
    programs.map((program) => [program.id, program.name]),
  );
  const editionName =
    editions.find((edition) => edition.id === editionId)?.name ?? "Edição";
  const visibleErrors = useMemo(() => {
    const search = filter.trim().toLowerCase();
    return emissionErrors.filter((error) =>
      [
        actorEmails[error.user_id],
        error.error_message,
        error.program_id ? programNames.get(error.program_id) : editionName,
      ]
        .filter(Boolean)
        .some((value) => value?.toLowerCase().includes(search)),
    );
  }, [actorEmails, editionName, emissionErrors, filter, programNames]);

  return (
    <PaginatedContainer<CertificationEmissionErrorI>
      items={visibleErrors}
      gap="4"
      pageSize={8}
      sort={sort}
      onSortChange={setSort}
      sortFields={[
        {
          key: "created_at",
          label: "Data",
          comparator: (a, b) =>
            new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
        },
        { key: "error_message", label: "Mensagem" },
        {
          key: "user_id",
          label: "E-mail",
          comparator: (a, b) =>
            (actorEmails[a.user_id] ?? "").localeCompare(
              actorEmails[b.user_id] ?? "",
            ),
        },
      ]}
      filterValue={filter}
      onFilterChange={setFilter}
      itemLabel="erros"
      filterPlaceholder="Buscar erros..."
      emptyState={
        <EmptyState
          icon={AlertTriangle}
          title="Nenhum erro de emissão"
          description="As falhas de emissão aparecerão aqui."
        />
      }
      renderItems={(items) =>
        items.map((error) => (
          <CertificationEmissionErrorCard
            key={error.id}
            error={error}
            email={actorEmails[error.user_id]}
            scopeName={
              error.program_id
                ? (programNames.get(error.program_id) ?? "Programação")
                : editionName
            }
          />
        ))
      }
    />
  );
}
