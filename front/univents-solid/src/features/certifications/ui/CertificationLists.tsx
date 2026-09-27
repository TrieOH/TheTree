import type { JSX } from "@solidjs/web";
import { useNavigate } from "@tanstack/solid-router";
import { For, Show, createEffect, createMemo, createSignal } from "solid-js";

import AlertTriangleIcon from "~icons/lucide/alert-triangle";
import AwardIcon from "~icons/lucide/award";
import BanIcon from "~icons/lucide/ban";
import ClockIcon from "~icons/lucide/clock";
import EyeIcon from "~icons/lucide/eye";
import ShieldCheckIcon from "~icons/lucide/shield-check";
import UserIcon from "~icons/lucide/user";

import { useQuery } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import type { SortState } from "@trieoh/ui-solid";
import { EmptyState, PaginatedContainer, cn } from "@trieoh/ui-solid";
import { asUniventsProfile, profileDisplayName } from "@/features/profile/model/profile-data";
import { programsQueryOptions } from "@/features/programs/api";
import { Reveal } from "@/shared/ui/Reveal";
import { AlertModal } from "@/widgets/ui/AlertModal";
import {
  allCertificationTemplatesQueryOptions,
  certificationsByEditionQueryOptions,
  emissionErrorsByEditionQueryOptions,
} from "../api";
import { useInvalidateCertificationMutation } from "../api/mutations";
import { DEFAULT_CERTIFICATION_TEMPLATE } from "../default-template";
import type {
  CertificationEmissionErrorI,
  CertificationI,
  CertificationTemplateI,
} from "../model";
import { CertViewer } from "./CertViewer";

const AlertTriangle = AlertTriangleIcon as unknown as (props: { class?: string }) => JSX.Element;
const Award = AwardIcon as unknown as (props: { class?: string }) => JSX.Element;
const Ban = BanIcon as unknown as (props: { class?: string }) => JSX.Element;
const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Eye = EyeIcon as unknown as (props: { class?: string }) => JSX.Element;
const ShieldCheck = ShieldCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const User = UserIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface CertificationListProps {
  eventId: string;
  editionId: string;
  editionName?: string;
}

export function CertificationList(props: CertificationListProps): JSX.Element {
  const { auth } = useAuth();
  const navigate = useNavigate();
  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<CertificationI>>({
    field: "issued_at",
    direction: "desc",
  });
  const [certToInvalidate, setCertToInvalidate] =
    createSignal<CertificationI | null>(null);
  const [certToView, setCertToView] = createSignal<CertificationI | null>(null);
  const [invalidateReason, setInvalidateReason] = createSignal("");

  const certificationsQuery = useQuery(() =>
    certificationsByEditionQueryOptions(props.editionId),
  );
  const programsQuery = useQuery(() => programsQueryOptions(props.editionId));
  const templatesQuery = useQuery(() =>
    allCertificationTemplatesQueryOptions(props.editionId),
  );

  const invalidateMutation = useInvalidateCertificationMutation();

  const certifications = () =>
    (certificationsQuery().data ?? []) as CertificationI[];
  const programs = () =>
    (programsQuery().data ?? []) as Array<{ id: string; name: string }>;
  const templates = () =>
    (templatesQuery().data ?? []) as CertificationTemplateI[];

  const programNames = createMemo(() => {
    const map = new Map<string, string>();
    for (const prog of programs()) {
      map.set(prog.id, prog.name);
    }
    return map;
  });

  const templatesMap = createMemo(() => {
    const map = new Map<string, CertificationTemplateI>();
    for (const tmpl of templates()) {
      map.set(tmpl.id, tmpl);
    }
    return map;
  });

  // User display names cache
  const [userDisplayNames, setUserDisplayNames] = createSignal<
    Record<string, string>
  >({});
  const actorIds = createMemo(() => [
    ...new Set(
      certifications()
        .map((c: CertificationI) => c.user_id)
        .filter(Boolean),
    ),
  ]);

  createEffect(
    () => actorIds(),
    (currentIds) => {
      if (!currentIds || currentIds.length === 0) return;
      void (async () => {
        const result: Record<string, string> = {};
        await Promise.all(
          currentIds.map(async (id) => {
            try {
              const res = await auth.getActorProfile(id);
              if (res.success && res.data) {
                const p = asUniventsProfile(res.data.profile ?? {});
                result[id] = profileDisplayName(p);
              }
            } catch {
              // ignore
            }
          }),
        );
        setUserDisplayNames((prev) => ({ ...prev, ...result }));
      })();
    },
  );

  const filteredCertifications = createMemo(() => {
    const list = certifications();
    const staticFilter = filter().trim().toLowerCase();
    if (!staticFilter) return list;

    return list.filter((cert) => {
      const staticName = userDisplayNames()[cert.user_id] ?? "";
      const staticProg = cert.program_id
        ? programNames().get(cert.program_id) ?? ""
        : (props.editionName ?? "");
      return (
        staticName.toLowerCase().includes(staticFilter) ||
        (cert.verification_hash ?? "").toLowerCase().includes(staticFilter) ||
        staticProg.toLowerCase().includes(staticFilter)
      );
    });
  });

  const handleConfirmInvalidate = async () => {
    const cert = certToInvalidate();
    if (!cert) return;
    try {
      await invalidateMutation.mutateAsync({
        editionId: props.editionId,
        certificationId: cert.id,
        reason: invalidateReason() || "Cancelado pelo organizador",
      });
      setCertToInvalidate(null);
      setInvalidateReason("");
    } catch {
      // handled
    }
  };

  const activeTemplateForView = createMemo<CertificationTemplateI>(() => {
    const cert = certToView();
    if (!cert || !cert.template_id) return DEFAULT_CERTIFICATION_TEMPLATE;
    return (
      templatesMap().get(cert.template_id) ?? DEFAULT_CERTIFICATION_TEMPLATE
    );
  });

  const activeVariablesForView = createMemo(() => {
    const cert = certToView();
    if (!cert) return {};
    const name = userDisplayNames()[cert.user_id] ?? "Participante";
    const activity = cert.program_id
      ? programNames().get(cert.program_id)
      : props.editionName;
    return {
      participant_name: name,
      event_name: "Evento",
      edition_name: props.editionName ?? "Edição",
      activity_name: activity ?? props.editionName ?? "Edição",
      cert_hash: cert.verification_hash,
      certified_at: new Date(cert.issued_at).toLocaleDateString("pt-BR"),
      verify_url:
        typeof window !== "undefined"
          ? `${window.location.origin}/verify/${cert.verification_hash}`
          : "",
    };
  });

  return (
    <div class="space-y-4">
      <PaginatedContainer<CertificationI>
        items={filteredCertifications()}
        layout="grid"
        minItemWidth="16rem"
        maxRows={(columns) => (columns === 1 ? 8 : 4)}
        gap="2"
        sort={sort()}
        onSortChange={setSort}
        sortFields={[
          {
            key: "issued_at",
            label: "Data de emissão",
            ascLabel: "Mais antigos primeiro",
            descLabel: "Mais recentes primeiro",
            comparator: (a, b) =>
              new Date(a.issued_at).getTime() - new Date(b.issued_at).getTime(),
          },
          {
            key: "user_id",
            label: "Participante",
            ascLabel: "A → Z",
            descLabel: "Z → A",
            comparator: (a, b) => {
              const nameA = userDisplayNames()[a.user_id] ?? a.user_id ?? "";
              const nameB = userDisplayNames()[b.user_id] ?? b.user_id ?? "";
              return nameA.localeCompare(nameB);
            },
          },
          {
            key: "verification_hash",
            label: "Código / Hash",
            ascLabel: "A → Z",
            descLabel: "Z → A",
            comparator: (a, b) =>
              (a.verification_hash ?? "").localeCompare(
                b.verification_hash ?? "",
              ),
          },
          {
            key: "valid",
            label: "Status",
            ascLabel: "Inválidos primeiro",
            descLabel: "Válidos primeiro",
            comparator: (a, b) => (a.valid === b.valid ? 0 : a.valid ? -1 : 1),
          },
        ]}
        filterValue={filter()}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por participante, hash ou atividade..."
        itemLabel="certificados"
        emptyState={
          <EmptyState
            class="border-0 bg-transparent px-0 py-8 shadow-none"
            icon={<Award class="size-6 text-foreground/70" />}
            eyebrow="Certificados emitidos"
            title="Nenhum certificado emitido"
            description={
              filter()
                ? "Nenhum certificado corresponde à busca."
                : "Os certificados emitidos para esta edição aparecerão aqui."
            }
          />
        }
        renderItems={(slice, options) => (
          <For each={slice}>
            {(cert, index) => {
              const participant = () =>
                userDisplayNames()[cert.user_id] ??
                (cert.user_id ? cert.user_id.slice(0, 8) : "Participante");
              const activity = () =>
                cert.program_id
                  ? programNames().get(cert.program_id) ?? "Atividade"
                  : props.editionName ?? "Participação na edição";

              return (
                <Reveal
                  delay={options.animate ? index() * 0.04 : 0}
                  animate={options.animate}
                  class="h-full"
                >
                  <article
                    class={cn(
                      "group relative flex h-full flex-col justify-between overflow-hidden rounded-lg border border-border/60 bg-card p-3.5 text-left transition-colors duration-150 hover:border-border",
                    )}
                  >
                    <div class="space-y-2">
                      <div class="flex items-center justify-between gap-2">
                        <span
                          class={cn(
                            "inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-[10px] font-medium border",
                            cert.valid
                              ? "border-emerald-500/20 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400"
                              : "border-destructive/20 bg-destructive/10 text-destructive",
                          )}
                        >
                          <Show
                            when={cert.valid}
                            fallback={<Ban class="size-3" />}
                          >
                            <ShieldCheck class="size-3" />
                          </Show>
                          <span>{cert.valid ? "Válido" : "Invalidado"}</span>
                        </span>

                        <span
                          class="font-mono text-[10px] text-muted-foreground truncate max-w-30"
                          title={cert.verification_hash}
                        >
                          {cert.verification_hash}
                        </span>
                      </div>

                      <div>
                        <h4
                          class="truncate text-xs font-medium text-foreground"
                          title={participant()}
                        >
                          {participant()}
                        </h4>
                        <p
                          class="truncate text-[11px] text-muted-foreground"
                          title={activity()}
                        >
                          {activity()}
                        </p>
                      </div>
                    </div>

                    <div class="mt-3 flex items-center justify-between border-t border-border/40 pt-2 text-[10px] text-muted-foreground">
                      <span>
                        {new Date(cert.issued_at).toLocaleDateString("pt-BR")}
                      </span>
                      <div class="flex items-center gap-2">
                        <button
                          type="button"
                          onClick={() => (() => { try { void navigate?.({ to: "/verify/$hash", params: { hash: cert.verification_hash } }); } catch {} })()}
                          class="inline-flex items-center gap-1 text-[11px] font-medium text-foreground hover:text-primary transition-colors cursor-pointer"
                        >
                          <Eye class="size-3" />
                          <span>Ver</span>
                        </button>
                        <Show when={cert.valid}>
                          <button
                            type="button"
                            onClick={() => setCertToInvalidate(cert)}
                            class="inline-flex items-center gap-1 text-[11px] font-medium text-destructive hover:underline cursor-pointer"
                          >
                            <Ban class="size-3" />
                            <span>Invalidar</span>
                          </button>
                        </Show>
                      </div>
                    </div>
                  </article>
                </Reveal>
              );
            }}
          </For>
        )}
      />

      {/* Invalidate Confirmation Modal */}
      <AlertModal
        open={certToInvalidate() !== null}
        onOpenChange={(open) => {
          if (!open) {
            setCertToInvalidate(null);
            setInvalidateReason("");
          }
        }}
        title="Invalidar certificado"
        description={`Tem certeza que deseja invalidar o certificado de ${userDisplayNames()[certToInvalidate()?.user_id ?? ""] ??
          "participante"
          }? Esta ação registrará o certificado como inválido.`}
        confirmLabel="Invalidar"
        variant="destructive"
        onConfirm={handleConfirmInvalidate}
      />

      {/* Certificate Viewer Modal */}
      <Show when={certToView()}>
        <CertViewer
          template={activeTemplateForView()}
          variables={activeVariablesForView()}
          open={certToView() !== null}
          onOpenChange={(open) => {
            if (!open) setCertToView(null);
          }}
        />
      </Show>
    </div>
  );
}

export function CertificationEmissionErrorsList(props: {
  editionId: string;
}): JSX.Element {
  const { auth } = useAuth();
  const [filter, setFilter] = createSignal("");
  const [sort, setSort] = createSignal<SortState<CertificationEmissionErrorI>>({
    field: "created_at",
    direction: "desc",
  });

  const errorsQuery = useQuery(() =>
    emissionErrorsByEditionQueryOptions(props.editionId),
  );
  const errors = () =>
    (errorsQuery().data ?? []) as CertificationEmissionErrorI[];

  // User display names cache
  const [userDisplayNames, setUserDisplayNames] = createSignal<
    Record<string, string>
  >({});
  const actorIds = createMemo(() => [
    ...new Set(
      errors()
        .map((e: CertificationEmissionErrorI) => e.user_id)
        .filter(Boolean),
    ),
  ]);

  createEffect(
    () => actorIds(),
    (currentIds) => {
      if (!currentIds || currentIds.length === 0) return;
      void (async () => {
        const result: Record<string, string> = {};
        await Promise.all(
          currentIds.map(async (id) => {
            try {
              const res = await auth.getActorProfile(id);
              if (res.success && res.data) {
                const p = asUniventsProfile(res.data.profile ?? {});
                result[id] = profileDisplayName(p);
              }
            } catch {
              // ignore
            }
          }),
        );
        setUserDisplayNames((prev) => ({ ...prev, ...result }));
      })();
    },
  );

  const filteredErrors = createMemo(() => {
    const list = errors();
    const staticFilter = filter().trim().toLowerCase();
    if (!staticFilter) return list;

    return list.filter((err) => {
      const staticName = userDisplayNames()[err.user_id] ?? "";
      return (
        staticName.toLowerCase().includes(staticFilter) ||
        (err.user_id ?? "").toLowerCase().includes(staticFilter) ||
        (err.error_message ?? "").toLowerCase().includes(staticFilter) ||
        (err.program_id ?? "").toLowerCase().includes(staticFilter)
      );
    });
  });

  return (
    <div class="space-y-4">
      <PaginatedContainer<CertificationEmissionErrorI>
        items={filteredErrors()}
        layout="grid"
        minItemWidth="16rem"
        maxRows={(columns) => (columns === 1 ? 8 : 4)}
        gap="2"
        sort={sort()}
        onSortChange={setSort}
        sortFields={[
          {
            key: "created_at",
            label: "Data do erro",
            ascLabel: "Mais antigos primeiro",
            descLabel: "Mais recentes primeiro",
            comparator: (a, b) =>
              new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
          },
          {
            key: "user_id",
            label: "Participante",
            ascLabel: "A → Z",
            descLabel: "Z → A",
            comparator: (a, b) => {
              const nameA = userDisplayNames()[a.user_id] ?? a.user_id ?? "";
              const nameB = userDisplayNames()[b.user_id] ?? b.user_id ?? "";
              return nameA.localeCompare(nameB);
            },
          },
          {
            key: "error_message",
            label: "Mensagem de erro",
            ascLabel: "A → Z",
            descLabel: "Z → A",
            comparator: (a, b) =>
              (a.error_message ?? "").localeCompare(b.error_message ?? ""),
          },
        ]}
        filterValue={filter()}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por participante, erro ou ID..."
        itemLabel="erros de emissão"
        emptyState={
          <EmptyState
            class="border-0 bg-transparent px-0 py-8 shadow-none"
            icon={<AlertTriangle class="size-6 text-foreground/70" />}
            eyebrow="Erros de emissão"
            title="Nenhum erro registrado"
            description={
              filter()
                ? "Nenhum erro corresponde à busca informada."
                : "Todas as emissões de certificados desta edição ocorreram normalmente."
            }
          />
        }
        renderItems={(slice, options) => (
          <For each={slice}>
            {(err, index) => {
              const participant = () =>
                userDisplayNames()[err.user_id] ??
                (err.user_id ? err.user_id.slice(0, 8) : "Participante");

              return (
                <Reveal
                  delay={options.animate ? index() * 0.04 : 0}
                  animate={options.animate}
                  class="h-full"
                >
                  <article
                    class={cn(
                      "group relative flex h-full flex-col justify-between overflow-hidden rounded-lg border border-destructive/25 bg-card p-3.5 text-left transition-colors duration-150 hover:border-destructive/40",
                    )}
                  >
                    <div class="space-y-2">
                      <div class="flex items-center justify-between gap-2">
                        <span class="inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-[10px] font-medium border border-destructive/20 bg-destructive/10 text-destructive">
                          <AlertTriangle class="size-3 shrink-0" />
                          <span>Falha de emissão</span>
                        </span>

                        <span class="inline-flex items-center gap-1 text-[10px] text-muted-foreground">
                          <Clock class="size-3" />
                          <span>
                            {new Date(err.created_at).toLocaleDateString(
                              "pt-BR",
                              {
                                day: "2-digit",
                                month: "2-digit",
                                hour: "2-digit",
                                minute: "2-digit",
                              },
                            )}
                          </span>
                        </span>
                      </div>

                      <div>
                        <div class="flex items-center gap-1.5 text-xs font-medium text-foreground">
                          <User class="size-3.5 shrink-0 text-muted-foreground" />
                          <h4 class="truncate" title={participant()}>
                            {participant()}
                          </h4>
                        </div>
                        <p class="mt-1.5 text-xs leading-relaxed text-muted-foreground wrap-break-word">
                          {err.error_message}
                        </p>
                      </div>
                    </div>

                    <div class="mt-3 flex items-center justify-between border-t border-border/40 pt-2 text-[10px] text-muted-foreground">
                      <span
                        class="font-mono text-muted-foreground/70 truncate max-w-40"
                        title={err.user_id}
                      >
                        ID: {err.user_id ? err.user_id.slice(0, 16) : "-"}
                      </span>
                    </div>
                  </article>
                </Reveal>
              );
            }}
          </For>
        )}
      />
    </div>
  );
}
