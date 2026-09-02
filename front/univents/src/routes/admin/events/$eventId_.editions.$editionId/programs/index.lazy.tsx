import { useQueries, useQuery } from "@tanstack/react-query";
import { createLazyFileRoute, useNavigate } from "@tanstack/react-router";
import type { SortState } from "@trieoh/ui-base";
import { EmptyState, PaginatedContainer } from "@trieoh/ui-base";
import { CalendarClock, CalendarDays, CalendarRange, Plus } from "lucide-react";
import { useMemo, useState } from "react";
import {
  allCertificationTemplatesQueryOptions,
  certificationTemplateLinksQueryOptions,
} from "@/features/certifications/api";
import {
  useLinkCertificationTemplateMutation,
  useUnlinkCertificationTemplateMutation,
} from "@/features/certifications/api/mutations";
import { ToolbarCombobox } from "@/features/certifications/editor/ui/toolbar-combobox";
import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import {
  useDeleteProgramMutation,
  useProgramMutation,
} from "@/features/programs/api/mutations";
import type { ProgramI } from "@/features/programs/model";
import { ManageProgramModal } from "@/features/programs/ui/ManageProgramModal";
import { ProgramAdminCard } from "@/features/programs/ui/ProgramAdminCard";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/shadcn/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/shared/ui/shadcn/dropdown-menu";
import { AlertModal } from "@/widgets/ui/alert-modal";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/programs/",
)({ component: ProgramsRoute });

function ProgramsRoute() {
  const { eventId, editionId } = Route.useParams();
  const navigate = useNavigate();
  const { data: programs = [] } = useQuery(programsQueryOptions(editionId));
  const { data: occurrences = [] } = useQuery(
    occurrencesQueryOptions(editionId),
  );

  const { data: templates = [] } = useQuery(
    allCertificationTemplatesQueryOptions(editionId),
  );

  const programTemplates = templates.filter(
    (template) => template.kind === "program_attendance",
  );

  const linkQueries = useQueries({
    queries: programTemplates.map((template) =>
      certificationTemplateLinksQueryOptions(template.id),
    ),
  });

  const linkedTemplateByProgram = new Map(
    linkQueries
      .flatMap((query, index) =>
        (query.data ?? []).map(
          (link) => [link.program_id, programTemplates[index]?.id] as const,
        ),
      )
      .filter((entry): entry is [string, string] => Boolean(entry[1])),
  );

  const mutation = useProgramMutation(editionId);
  const deleteMutation = useDeleteProgramMutation();
  const [filter, setFilter] = useState("");
  const [startsAfter, setStartsAfter] = useState("");
  const [endsBefore, setEndsBefore] = useState("");
  const [sort, setSort] = useState<SortState<ProgramI>>({
    field: "name",
    direction: "asc",
  });
  const [editing, setEditing] = useState<ProgramI>();
  const [modalOpen, setModalOpen] = useState(false);
  const [programToDelete, setProgramToDelete] = useState<ProgramI | null>(null);
  const [certificateProgram, setCertificateProgram] = useState<ProgramI | null>(
    null,
  );
  const [certificateTemplateId, setCertificateTemplateId] = useState("");
  const linkCertificate = useLinkCertificationTemplateMutation();
  const unlinkCertificate = useUnlinkCertificationTemplateMutation();
  const [programToUnlink, setProgramToUnlink] = useState<ProgramI | null>(null);

  const firstOccurrenceByProgram = useMemo(
    () =>
      occurrences.reduce((starts, occurrence) => {
        const current = starts.get(occurrence.program_id);
        const timestamp = new Date(occurrence.starts_at).getTime();
        if (current === undefined || timestamp < current) {
          starts.set(occurrence.program_id, timestamp);
        }
        return starts;
      }, new Map<string, number>()),
    [occurrences],
  );
  const startsAfterTimestamp = startsAfter
    ? new Date(startsAfter).getTime()
    : Number.NEGATIVE_INFINITY;
  const endsBeforeTimestamp = endsBefore
    ? new Date(endsBefore).getTime()
    : Number.POSITIVE_INFINITY;
  const hasDateTimeFilter = Boolean(startsAfter || endsBefore);
  const programsInDateTimeRange = useMemo(
    () =>
      new Set(
        occurrences
          .filter(
            (item) =>
              new Date(item.starts_at).getTime() <= endsBeforeTimestamp &&
              new Date(item.ends_at).getTime() >= startsAfterTimestamp,
          )
          .map((item) => item.program_id),
      ),
    [endsBeforeTimestamp, occurrences, startsAfterTimestamp],
  );
  const filtered = programs.filter(
    (program) =>
      program.name.toLowerCase().includes(filter.trim().toLowerCase()) &&
      (!hasDateTimeFilter || programsInDateTimeRange.has(program.id)),
  );

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<ProgramI>
        items={filtered}
        layout="grid"
        minItemWidth="16rem"
        pageSize={9}
        gap="6"
        sort={sort}
        onSortChange={setSort}
        sortFields={[
          { key: "name", label: "Nome" },
          { key: "kind", label: "Tipo" },
          {
            key: "created_at",
            label: "Horário",
            comparator: (a, b) => {
              const aStart = firstOccurrenceByProgram.get(a.id);
              const bStart = firstOccurrenceByProgram.get(b.id);
              if (aStart === bStart) return a.name.localeCompare(b.name);
              if (aStart === undefined)
                return sort.direction === "asc" ? 1 : -1;
              if (bStart === undefined)
                return sort.direction === "asc" ? -1 : 1;
              return aStart - bStart;
            },
          },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar programa..."
        itemLabel="programas"
        headerActions={
          <>
            <DropdownMenu>
              <DropdownMenuTrigger
                render={
                  <Button
                    type="button"
                    variant="outline"
                    className="relative h-9 gap-2"
                  />
                }
              >
                <CalendarClock className="size-4" />
                Data e hora
                {hasDateTimeFilter ? (
                  <span className="size-2 rounded-full bg-primary" />
                ) : null}
              </DropdownMenuTrigger>
              <DropdownMenuContent
                align="end"
                className="w-[calc(100vw-2rem)]! max-w-sm! space-y-3 p-4 sm:w-80!"
                onClick={(event) => event.stopPropagation()}
                onKeyDown={(event) => event.stopPropagation()}
              >
                <div>
                  <p className="text-sm font-semibold">Período da ocorrência</p>
                  <p className="mt-0.5 text-xs text-muted-foreground">
                    Encontre programações dentro de um intervalo.
                  </p>
                </div>
                <label className="block space-y-1.5">
                  <span className="text-xs font-medium text-muted-foreground">
                    De
                  </span>
                  <input
                    type="datetime-local"
                    value={startsAfter}
                    onChange={(event) => setStartsAfter(event.target.value)}
                    className="h-10 w-full rounded-md border border-border bg-background px-3 text-sm outline-none focus:border-ring"
                  />
                </label>
                <label className="block space-y-1.5">
                  <span className="text-xs font-medium text-muted-foreground">
                    Até
                  </span>
                  <input
                    type="datetime-local"
                    value={endsBefore}
                    min={startsAfter || undefined}
                    onChange={(event) => setEndsBefore(event.target.value)}
                    className="h-10 w-full rounded-md border border-border bg-background px-3 text-sm outline-none focus:border-ring"
                  />
                </label>
                {hasDateTimeFilter ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="w-full"
                    onClick={() => {
                      setStartsAfter("");
                      setEndsBefore("");
                    }}
                  >
                    Limpar período
                  </Button>
                ) : null}
              </DropdownMenuContent>
            </DropdownMenu>
            <Button
              type="button"
              className="h-9 gap-2"
              onClick={() => {
                setEditing(undefined);
                setModalOpen(true);
              }}
            >
              <Plus className="size-4" />
              Novo programa
            </Button>
            <Button
              type="button"
              variant="outline"
              className="h-9 gap-2"
              onClick={() =>
                navigate({
                  to: "/admin/events/$eventId/editions/$editionId/programs/calendar",
                  params: { eventId, editionId },
                })
              }
            >
              <CalendarDays className="size-4" />
              Abrir calendário
            </Button>
          </>
        }
        emptyState={
          <EmptyState
            icon={CalendarRange}
            eyebrow="Programação"
            title="Nenhum programa encontrado"
            description="Crie o primeiro programa desta edição."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((program, index) => (
            <ProgramAdminCard
              key={program.id}
              program={program}
              index={index}
              occurrences={occurrences.filter(
                (item) => item.program_id === program.id,
              )}
              onEdit={() => {
                setEditing(program);
                setModalOpen(true);
              }}
              onManageOccurrences={() =>
                navigate({
                  to: "/admin/events/$eventId/editions/$editionId/programs/$programId/occurrences",
                  params: {
                    eventId,
                    editionId,
                    programId: program.id,
                  },
                })
              }
              onDelete={() => {
                setProgramToDelete(program);
              }}
              hasCertificate={linkedTemplateByProgram.has(program.id)}
              onManageCertificate={() => {
                setCertificateProgram(program);
                setCertificateTemplateId("");
              }}
              onUnlinkCertificate={() => setProgramToUnlink(program)}
            />
          ))
        }
      />
      <AlertModal
        open={Boolean(programToDelete)}
        onOpenChange={(open) => !open && setProgramToDelete(null)}
        title="Excluir programa?"
        description={`O programa "${programToDelete?.name ?? ""}" e suas ocorrências serão excluídos. Essa ação não pode ser desfeita.`}
        confirmLabel="Excluir programa"
        variant="destructive"
        loading={deleteMutation.isPending}
        onConfirm={async () => {
          if (!programToDelete) return;
          await deleteMutation.mutateAsync(programToDelete.id);
          setProgramToDelete(null);
        }}
      />
      <Dialog
        open={certificateProgram !== null}
        onOpenChange={(open) => !open && setCertificateProgram(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Certificado do programa</DialogTitle>
            <DialogDescription>
              Selecione o certificado que será usado por este programa.
            </DialogDescription>
          </DialogHeader>
          <ToolbarCombobox
            value={certificateTemplateId}
            options={programTemplates.map((template) => ({
              value: template.id,
              label: template.name,
              description: template.description ?? undefined,
            }))}
            placeholder="Selecionar certificado"
            searchPlaceholder="Buscar certificado..."
            onChange={setCertificateTemplateId}
            className="w-full"
            triggerClassName="h-10"
          />
          <DialogFooter>
            <Button
              disabled={
                !certificateProgram ||
                !certificateTemplateId ||
                linkCertificate.isPending
              }
              onClick={() => {
                if (!certificateProgram || !certificateTemplateId) return;
                linkCertificate.mutate(
                  {
                    templateId: certificateTemplateId,
                    programId: certificateProgram.id,
                  },
                  { onSuccess: () => setCertificateProgram(null) },
                );
              }}
            >
              Vincular certificado
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      <AlertModal
        open={Boolean(programToUnlink)}
        onOpenChange={(open) => !open && setProgramToUnlink(null)}
        title="Desvincular certificado?"
        description={`O programa "${programToUnlink?.name ?? ""}" ficará sem certificado específico.`}
        confirmLabel="Desvincular certificado"
        variant="destructive"
        loading={unlinkCertificate.isPending}
        onConfirm={async () => {
          if (!programToUnlink) return;
          const templateId = linkedTemplateByProgram.get(programToUnlink.id);
          if (templateId)
            await unlinkCertificate.mutateAsync({
              templateId,
              programId: programToUnlink.id,
            });
          setProgramToUnlink(null);
        }}
      />
      <ManageProgramModal
        open={modalOpen}
        program={editing}
        onOpenChange={setModalOpen}
        onSave={async (data) => {
          const program = await mutation.mutateAsync({
            id: editing?.id,
            data,
          });
          return Boolean(program);
        }}
      />
    </div>
  );
}
