import type { JSX } from "@solidjs/web";
import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useQuery, useQueryClient } from "@trieoh/front-core-solid";
import type { SortState } from "@trieoh/ui-solid";
import { Button, EmptyState, PaginatedContainer } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal, onCleanup } from "solid-js";

import CalendarDaysIcon from "~icons/lucide/calendar-days";
import CalendarPlusIcon from "~icons/lucide/calendar-plus";
import FilterIcon from "~icons/lucide/filter";
import AwardIcon from "~icons/lucide/award";
import { Dialog } from "@trieoh/ui-solid";
import {
  allCertificationTemplatesQueryOptions,
  certificationKeys,
  editionProgramTemplateLinksQueryOptions,
} from "@/features/certifications/api";
import {
  useEmitProgramCertificationsMutation,
  useLinkCertificationTemplateMutation,
  useUnlinkCertificationTemplateMutation,
} from "@/features/certifications/api/mutations";
import type { CertificationTemplateI } from "@/features/certifications/model";


import {
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import {
  useCreateProgramMutation,
  useDeleteProgramMutation,
  useUpdateProgramMutation,
} from "@/features/programs/api/mutations";
import type {
  OccurrenceI,
  ProgramCreateInput,
  ProgramI,
} from "@/features/programs/model";
import { AdminCreateProgramCard } from "@/features/programs/ui/AdminCreateProgramCard";
import { AdminProgramCard } from "@/features/programs/ui/AdminProgramCard";
import { ManageProgramDialog } from "@/features/programs/ui/ManageProgramDialog";
import { toast } from "@/shared/ui/toast";
import { AlertModal } from "@/widgets/ui/AlertModal";
import { Combobox, type ComboboxOption } from "@/shared/ui/Combobox";

const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;
const CalendarPlus = CalendarPlusIcon as unknown as (props: { class?: string }) => JSX.Element;
const Filter = FilterIcon as unknown as (props: { class?: string }) => JSX.Element;
const Award = AwardIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute(
  "/admin/events/$eventId_/editions/$editionId/programs/",
)({
  head: () => ({
    meta: [{ title: "Programação - Univents Admin" }],
  }),
  component: AdminProgramsRoute,
});

function AdminProgramsRoute(): JSX.Element {
  const params = Route.useParams();
  const eventId = () => params().eventId;
  const editionId = () => params().editionId;
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // Queries
  const programsQuery = useQuery(() => programsQueryOptions(editionId()));
  const occurrencesQuery = useQuery(() => occurrencesQueryOptions(editionId()));

  // Mutations
  const createMutation = useCreateProgramMutation(editionId);
  const updateMutation = useUpdateProgramMutation(editionId);
  const deleteMutation = useDeleteProgramMutation(editionId);

  // Certification Queries & Mutations
  const templatesQuery = useQuery(() => allCertificationTemplatesQueryOptions(editionId()));
  const programTemplates = createMemo<CertificationTemplateI[]>(() =>
    ((templatesQuery().data ?? []) as CertificationTemplateI[]).filter(
      (t) => t.kind === "program_attendance",
    ),
  );

  const [linkingLoading, setLinkingLoading] = createSignal(false);

  const linksQuery = useQuery(() =>
    editionProgramTemplateLinksQueryOptions(
      editionId(),
      programTemplates(),
      queryClient,
    ),
  );

  const linkedTemplateByProgram = createMemo(
    () => linksQuery().data ?? new Map<string, string>(),
  );

  const linkMutation = useLinkCertificationTemplateMutation();
  const unlinkMutation = useUnlinkCertificationTemplateMutation();
  const emitMutation = useEmitProgramCertificationsMutation();

  const [certificateProgram, setCertificateProgram] = createSignal<ProgramI | null>(null);
  const [certificateTemplateId, setCertificateTemplateId] = createSignal<string>("");
  const [programToUnlink, setProgramToUnlink] = createSignal<ProgramI | null>(null);
  const [emittingProgramId, setEmittingProgramId] = createSignal<string | null>(null);

  // Cooldown handling
  const [now, setNow] = createSignal(Date.now());
  const cooldownTimer = setInterval(() => setNow(Date.now()), 1000);
  onCleanup(() => clearInterval(cooldownTimer));

  const COOLDOWN_MS = 60_000;
  function getCooldownKey(programId: string) {
    return `univents:cert-cooldown:${programId}`;
  }
  function startCooldown(programId: string) {
    try {
      localStorage.setItem(getCooldownKey(programId), String(Date.now() + COOLDOWN_MS));
      setNow(Date.now());
    } catch {
      // ignore
    }
  }
  function cooldownRemaining(programId: string) {
    now();
    try {
      const raw = localStorage.getItem(getCooldownKey(programId));
      if (!raw) return 0;
      const expiresAt = Number(raw);
      if (!Number.isFinite(expiresAt)) return 0;
      const remaining = expiresAt - Date.now();
      return remaining > 0 ? remaining : 0;
    } catch {
      return 0;
    }
  }
  function formatCooldown(ms: number) {
    const seconds = Math.ceil(ms / 1000);
    return `${seconds}s`;
  }

  const handleLinkCertificate = async () => {
    const prog = certificateProgram();
    const templateId = certificateTemplateId();
    if (!prog || !templateId) return;

    setLinkingLoading(true);
    try {
      await linkMutation.mutateAsync({ templateId, programId: prog.id });
      toast.success("Certificado vinculado com sucesso!");
      setCertificateProgram(null);
      setCertificateTemplateId("");
      await queryClient.invalidateQueries({ queryKey: certificationKeys.editionProgramLinks(editionId()) });
    } catch {
      toast.error("Erro ao vincular certificado ao programa.");
    } finally {
      setLinkingLoading(false);
    }
  };

  const handleUnlinkCertificate = async () => {
    const prog = programToUnlink();
    if (!prog) return;

    const templateId = linkedTemplateByProgram().get(prog.id);
    if (!templateId) {
      setProgramToUnlink(null);
      return;
    }

    try {
      await unlinkMutation.mutateAsync({ templateId, programId: prog.id });
      toast.success("Certificado desvinculado!");
      setProgramToUnlink(null);
      await queryClient.invalidateQueries({ queryKey: certificationKeys.editionProgramLinks(editionId()) });
    } catch {
      toast.error("Erro ao desvincular certificado.");
    }
  };

  const handleEmitCertificates = async (prog: ProgramI) => {
    if (cooldownRemaining(prog.id) > 0 || emitMutation.result().isPending) return;

    startCooldown(prog.id);
    setEmittingProgramId(prog.id);
    try {
      await emitMutation.mutateAsync({ editionId: editionId(), programId: prog.id });
      toast.success(`Emissão de certificados iniciada para "${prog.name}"!`);
    } catch {
      toast.error("Erro ao solicitar emissão de certificados.");
    } finally {
      setEmittingProgramId(null);
    }
  };


  const programs = createMemo(() => (programsQuery().data ?? []) as ProgramI[]);
  const occurrences = createMemo(() => (occurrencesQuery().data ?? []) as OccurrenceI[]);

  // States & Filters
  const [filter, setFilter] = createSignal("");
  const [kindFilter, setKindFilter] = createSignal<"all" | "activity" | "checkpoint">("all");
  const [dateFilter, setDateFilter] = createSignal<string | null>(null);
  const [sort, setSort] = createSignal<SortState<ProgramI>>({
    field: "name",
    direction: "asc",
  });

  const distinctDates = createMemo(() => {
    const datesMap = new Map<string, number>();
    for (const occ of occurrences()) {
      const key = occ.starts_at.slice(0, 10);
      datesMap.set(key, (datesMap.get(key) ?? 0) + 1);
    }
    return Array.from(datesMap.entries())
      .sort((a, b) => a[0].localeCompare(b[0]))
      .map(([dateKey, count]) => {
        const [y, m, d] = dateKey.split("-").map(Number);
        const dateObj = new Date(y, m - 1, d);
        const weekday = dateObj.toLocaleDateString("pt-BR", { weekday: "short" });
        return {
          key: dateKey,
          label: `${String(d).padStart(2, "0")}/${String(m).padStart(2, "0")}`,
          weekday,
          count,
        };
      });
  });

  const dateOptions = createMemo((): ComboboxOption[] => {
    const list: ComboboxOption[] = [
      { value: "all", label: "Todos os dias" },
    ];
    for (const d of distinctDates()) {
      list.push({
        value: d.key,
        label: `${d.label} (${d.weekday})`,
        description: `${d.count} horário${d.count > 1 ? "s" : ""}`,
      });
    }
    return list;
  });

  const activityCount = createMemo(() => programs().filter((p) => p.kind === "activity").length);
  const checkpointCount = createMemo(() => programs().filter((p) => p.kind === "checkpoint").length);

  const filteredPrograms = createMemo(() => {
    const all = programs();
    const currentKind = kindFilter();
    const targetDate = dateFilter();
    const allOccurrences = occurrences();

    const result: ProgramI[] = [];
    for (const p of all) {
      if (currentKind !== "all" && p.kind !== currentKind) {
        continue;
      }
      if (targetDate) {
        let hasOccurrence = false;
        for (const occ of allOccurrences) {
          if (occ.program_id === p.id && occ.starts_at.slice(0, 10) === targetDate) {
            hasOccurrence = true;
            break;
          }
        }
        if (!hasOccurrence) continue;
      }
      result.push(p);
    }
    return result;
  });



  const [filterMenuOpen, setFilterMenuOpen] = createSignal(false);
  let filterMenuRef: HTMLDivElement | undefined;
  let filterPopoverEl: HTMLDivElement | undefined;

  const clampFilterMenu = (el?: HTMLElement | null) => {
    const target = el ?? filterPopoverEl;
    if (!target || typeof window === "undefined") return;
    requestAnimationFrame(() => {
      target.style.transform = "";
      const rect = target.getBoundingClientRect();
      const vw = window.innerWidth;
      if (rect.left < 8) {
        target.style.transform = `translateX(${8 - rect.left}px)`;
      } else if (rect.right > vw - 8) {
        target.style.transform = `translateX(${vw - 8 - rect.right}px)`;
      }
    });
  };

  const activeFilterCount = createMemo(() => {
    let count = 0;
    if (kindFilter() !== "all") count++;
    if (dateFilter() !== null) count++;
    return count;
  });

  const clearAllFilters = () => {
    setKindFilter("all");
    setDateFilter(null);
  };

  if (typeof document !== "undefined") {
    const handlePointerDown = (event: MouseEvent) => {
      const target = event.target as Node | null;
      if (filterMenuRef && !filterMenuRef.contains(target)) {
        const portalEl = (target as HTMLElement | null)?.closest?.("[role='listbox']") || (target as HTMLElement | null)?.closest?.(".bg-popover");
        if (!portalEl) {
          setFilterMenuOpen(false);
        }
      }
    };
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setFilterMenuOpen(false);
      }
    };
    const handleResize = () => {
      if (filterMenuOpen()) {
        clampFilterMenu();
      }
    };
    document.addEventListener("mousedown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);
    window.addEventListener("resize", handleResize);
    onCleanup(() => {
      document.removeEventListener("mousedown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
      window.removeEventListener("resize", handleResize);
    });
  }

  const [creating, setCreating] = createSignal(false);
  const [editing, setEditing] = createSignal<ProgramI | null>(null);
  const [deleting, setDeleting] = createSignal<ProgramI | null>(null);

  const handleCreate = async (data: ProgramCreateInput) => {
    try {
      await createMutation.mutateAsync(data);
      toast.success("Programa criado com sucesso!");
      setCreating(false);
      return true;
    } catch {
      toast.error("Erro ao criar programa.");
      return false;
    }
  };

  const handleUpdate = async (id: string, data: ProgramCreateInput) => {
    try {
      await updateMutation.mutateAsync({ id, data });
      toast.success("Programa atualizado com sucesso!");
      setEditing(null);
      return true;
    } catch {
      toast.error("Erro ao atualizar programa.");
      return false;
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteMutation.mutateAsync(id);
      toast.success("Programa excluído com sucesso!");
      setDeleting(null);
      return true;
    } catch {
      toast.error("Erro ao excluir programa.");
      return false;
    }
  };

  const emptyStateAction = (
    <Button
      size="sm"
      class="gap-2 rounded-sm py-4"
      onClick={() => setCreating(true)}
    >
      <CalendarPlus class="size-4" />
      Novo programa
    </Button>
  );

  return (
    <div class="space-y-6">
      {/* Header with Navigation */}
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-foreground">
            Programação do Evento
          </h1>
          <p class="text-xs sm:text-sm text-muted-foreground">
            Gerencie as palestras, workshops, checkpoints e horários da edição.
          </p>
        </div>

        <div class="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() =>
              navigate({
                to: "/admin/events/$eventId/editions/$editionId/programs/calendar",
                params: { eventId: eventId(), editionId: editionId() },
              })
            }
            class="gap-2 cursor-pointer"
          >
            <CalendarDays class="size-4 text-primary" />
            <span>Abrir calendário</span>
          </Button>
        </div>
      </div>

      <PaginatedContainer<ProgramI>
        items={filteredPrograms()}
        layout="grid"
        minItemWidth="16rem"
        maxRows={(columns) => (columns === 1 ? 8 : 4)}
        gap="2"
        sort={sort()}
        onSortChange={setSort}
        sortFields={[
          {
            key: "name",
            label: "Nome",
            ascLabel: "A → Z",
            descLabel: "Z → A",
          },
          {
            key: "kind",
            label: "Tipo",
            ascLabel: "Atividade primeiro",
            descLabel: "Checkpoint primeiro",
          },
        ]}
        filterValue={filter()}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por nome ou descrição..."
        filterFields={["name", "description"]}
        headerActions={
          <div class="relative" ref={(el) => (filterMenuRef = el)}>
            <button
              type="button"
              onClick={() => setFilterMenuOpen((prev) => !prev)}
              class={`relative flex items-center gap-1.5 h-9 px-3 text-xs font-medium rounded-md border transition-all cursor-pointer ${activeFilterCount() > 0
                ? "border-primary/50 bg-primary/10 text-primary hover:bg-primary/20"
                : "border-border bg-background hover:bg-muted text-foreground"
                }`}
              title="Filtrar programações"
            >
              <Filter class="size-3.5" />
              <span>Filtrar</span>
              <Show when={activeFilterCount() > 0}>
                <span class="absolute -top-1.5 -right-1.5 size-4 rounded-full bg-primary text-primary-foreground text-[10px] flex items-center justify-center font-bold shadow-xs">
                  {activeFilterCount()}
                </span>
              </Show>
            </button>

            <Show when={filterMenuOpen()}>
              <div
                ref={(el) => {
                  filterPopoverEl = el;
                  clampFilterMenu(el);
                }}
                class="absolute left-0 sm:left-auto sm:right-0 mt-2 w-72 max-w-[calc(100vw-1rem)] rounded-xl border border-border bg-popover p-3 shadow-xl z-50 text-xs space-y-3"
              >
                <div class="flex items-center justify-between pb-1.5 border-b border-border/60">
                  <div class="flex items-center gap-1.5 font-semibold text-foreground">
                    <Filter class="size-3.5 text-primary" />
                    <span>Filtros</span>
                  </div>
                  <Show when={activeFilterCount() > 0}>
                    <button
                      type="button"
                      onClick={clearAllFilters}
                      class="text-[11px] text-primary hover:underline cursor-pointer font-medium"
                    >
                      Limpar todos
                    </button>
                  </Show>
                </div>

                {/* Kind Filter Tabs */}
                <div class="space-y-1">
                  <label class="text-[11px] font-medium text-muted-foreground">Tipo de atividade</label>
                  <div class="flex rounded-md border border-border bg-muted/40 p-0.5 text-[11px]">
                    <button
                      type="button"
                      onClick={() => setKindFilter("all")}
                      class={`flex-1 py-1 rounded text-center font-medium transition-colors cursor-pointer ${kindFilter() === "all"
                        ? "bg-background text-foreground shadow-2xs font-semibold"
                        : "text-muted-foreground hover:text-foreground"
                        }`}
                    >
                      Todos ({programs().length})
                    </button>
                    <button
                      type="button"
                      onClick={() => setKindFilter("activity")}
                      class={`flex-1 py-1 rounded text-center font-medium transition-colors cursor-pointer ${kindFilter() === "activity"
                        ? "bg-background text-foreground shadow-2xs font-semibold"
                        : "text-muted-foreground hover:text-foreground"
                        }`}
                    >
                      Atividades ({activityCount()})
                    </button>
                    <button
                      type="button"
                      onClick={() => setKindFilter("checkpoint")}
                      class={`flex-1 py-1 rounded text-center font-medium transition-colors cursor-pointer ${kindFilter() === "checkpoint"
                        ? "bg-background text-foreground shadow-2xs font-semibold"
                        : "text-muted-foreground hover:text-foreground"
                        }`}
                    >
                      Checkpoints ({checkpointCount()})
                    </button>
                  </div>
                </div>

                {/* Event Day Combobox */}
                <Show when={distinctDates().length > 0}>
                  <div class="space-y-1">
                    <label class="text-[11px] font-medium text-muted-foreground">Dia específico</label>
                    <Combobox
                      value={dateFilter() ?? "all"}
                      options={dateOptions()}
                      placeholder="Todos os dias..."
                      searchPlaceholder="Buscar dia..."
                      onChange={(val) => setDateFilter(val === "all" ? null : val)}
                      class="w-full"
                      triggerClass="h-8 text-xs"
                    />
                  </div>
                </Show>
              </div>
            </Show>
          </div>
        }
        itemLabel="programas"
        emptyState={
          <EmptyState
            icon={<CalendarDays class="size-6 text-foreground/70" />}
            eyebrow="Programação"
            title="Nenhum programa cadastrado"
            description={
              filter()
                ? "Nenhum programa corresponde à busca informada."
                : "Crie o primeiro programa desta edição para definir a grade de horários."
            }
            class="border-0 bg-transparent px-0 py-4 shadow-none"
            action={emptyStateAction}
          />
        }
        renderItems={(slice, options) => (
          <>
            <AdminCreateProgramCard
              index={0}
              animate={options.animate}
              onCreate={() => setCreating(true)}
            />

            <For each={slice}>
              {(program, index) => (
                <AdminProgramCard
                  program={program}
                  eventId={eventId()}
                  index={index() + 1}
                  animate={options.animate}
                  occurrences={occurrences().filter((o) => o.program_id === program.id)}
                  occurrencesHref={`/admin/events/${eventId()}/editions/${editionId()}/programs/${program.id}/occurrences`}
                  hasCertificate={Boolean(linkedTemplateByProgram().get(program.id))}
                  isEmittingCertificates={emittingProgramId() === program.id}
                  emissionCooldownLabel={
                    cooldownRemaining(program.id) > 0
                      ? formatCooldown(cooldownRemaining(program.id))
                      : undefined
                  }
                  onManageCertificate={(p) => {
                    setCertificateProgram(p);
                    setCertificateTemplateId(linkedTemplateByProgram().get(p.id) ?? "");
                  }}
                  onUnlinkCertificate={(p) => setProgramToUnlink(p)}
                  onEmitCertificates={(p) => void handleEmitCertificates(p)}
                  onEdit={setEditing}
                  onDelete={(p) => setDeleting(p)}
                  onManageOccurrences={(p) =>
                    navigate({
                      to: "/admin/events/$eventId/editions/$editionId/programs/$programId/occurrences",
                      params: {
                        eventId: eventId(),
                        editionId: editionId(),
                        programId: p.id,
                      },
                    })
                  }
                  onOpenCalendar={(_p) =>
                    navigate({
                      to: "/admin/events/$eventId/editions/$editionId/programs/calendar",
                      params: { eventId: eventId(), editionId: editionId() },
                    })
                  }
                />
              )}
            </For>
          </>
        )}
      />

      {/* Dialog: Create or Edit */}
      <ManageProgramDialog
        open={creating() || editing() !== null}
        program={editing()}
        onOpenChange={(open) => {
          if (!open) {
            setCreating(false);
            setEditing(null);
          }
        }}
        onSubmit={(data: ProgramCreateInput) => {
          const current = editing();
          if (current) {
            return handleUpdate(current.id, data);
          }
          return handleCreate(data);
        }}
      />

      {/* Dialog: Link Certificate */}
      <Dialog
        open={certificateProgram() !== null}
        onOpenChange={(open) => {
          if (!open) {
            setCertificateProgram(null);
            setCertificateTemplateId("");
          }
        }}
        title="Vincular certificado ao programa"
        description={
          certificateProgram()
            ? `Selecione o modelo de certificado emitido para os participantes de "${certificateProgram()?.name}".`
            : undefined
        }
      >
        <div class="space-y-4 pt-1">
          <Show
            when={programTemplates().length > 0}
            fallback={
              <div class="rounded-lg border border-border/60 bg-muted/40 p-4 text-center">
                <Award class="mx-auto size-8 text-muted-foreground/60" />
                <p class="mt-2 text-xs font-medium text-foreground">
                  Nenhum modelo de atividade disponível
                </p>
                <p class="mt-1 text-[11px] text-muted-foreground">
                  Crie um modelo com a categoria &quot;Presença em atividade&quot; na aba de Certificações desta edição.
                </p>
              </div>
            }
          >
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-foreground">Modelo de Certificado</label>
              <Combobox
                value={certificateTemplateId()}
                options={programTemplates().map((t) => ({
                  value: t.id,
                  label: t.name,
                  description: t.description ?? undefined,
                }))}
                placeholder="Selecione um modelo..."
                searchPlaceholder="Buscar modelo..."
                onChange={(val) => setCertificateTemplateId(val)}
                class="w-full"
              />
            </div>
          </Show>

          <div class="flex items-center justify-end gap-2 pt-2 border-t border-border/40">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                setCertificateProgram(null);
                setCertificateTemplateId("");
              }}
            >
              Cancelar
            </Button>
            <Button
              type="button"
              disabled={!certificateTemplateId() || linkingLoading()}
              onClick={handleLinkCertificate}
            >
              <Show when={linkingLoading()} fallback={"Vincular certificado"}>
                <span>Vinculando...</span>
              </Show>
            </Button>
          </div>
        </div>
      </Dialog>

      {/* Modal: Unlink Certificate */}
      <AlertModal
        open={programToUnlink() !== null}
        onOpenChange={(open) => {
          if (!open) setProgramToUnlink(null);
        }}
        title="Desvincular certificado?"
        description={
          programToUnlink()
            ? `O programa "${programToUnlink()?.name}" não terá mais um certificado específico vinculado.`
            : undefined
        }
        confirmLabel="Desvincular certificado"
        variant="destructive"
        loading={unlinkMutation.result().isPending}
        onConfirm={handleUnlinkCertificate}
      />

      {/* Modal: Delete */}
      <AlertModal
        open={deleting() !== null}
        onOpenChange={(open) => {
          if (!open) setDeleting(null);
        }}
        title="Excluir programa"
        description={
          deleting()
            ? `Tem certeza que deseja excluir "${deleting()?.name}"? Esta ação removerá os horários associados.`
            : undefined
        }
        confirmLabel="Excluir programa"
        variant="destructive"
        loading={deleteMutation.result().isPending}
        onConfirm={() => {
          const p = deleting();
          if (p) void handleDelete(p.id);
        }}
      />
    </div>
  );
}
