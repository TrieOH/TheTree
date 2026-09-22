import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { cn } from "@trieoh/ui-solid";
import { Show, createEffect, createMemo, createSignal } from "solid-js";

import BadgeCheckIcon from "~icons/lucide/badge-check";
import FileX2Icon from "~icons/lucide/file-x-2";
import Loader2Icon from "~icons/lucide/loader-2";

import {
  certificationTemplateQueryOptions,
  certificationVerificationQueryOptions,
} from "@/features/certifications/api";
import { getCertificationTemplateOrDefault } from "@/features/certifications/default-template";
import { DEFAULT_CERTIFICATE_CANVAS } from "@/features/certifications/editor/constants";
import type {
  CertificateVariableValues,
} from "@/features/certifications/editor/variables";
import {
  CertificateDownloadButtons,
  CertificateTemplateStaticView,
} from "@/features/certifications/ui/CertViewer";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import { editionLocationQueryOptions } from "@/features/editions/api";
import { asUniventsProfile, profileDisplayName } from "@/features/profile/model/profile-data";
import { occurrencesQueryOptions, programsQueryOptions } from "@/features/programs/api";

const BadgeCheck = BadgeCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const FileX2 = FileX2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute("/verify/$hash")({
  head: () => ({
    meta: [{ title: "Verificação de Autenticidade de Certificado - Univents" }],
  }),
  component: VerifyCertificationPage,
});

function formatCertifiedAt(value: string) {
  return new Date(value).toLocaleString("pt-BR", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function DetailItem(props: { label: string; value?: string | null }) {
  return (
    <Show when={props.value}>
      <div class="py-3">
        <dt class="text-[10px] uppercase tracking-wider text-muted-foreground">
          {props.label}
        </dt>
        <dd class="mt-1 text-xs font-medium text-foreground truncate" title={props.value ?? ""}>
          {props.value}
        </dd>
      </div>
    </Show>
  );
}

function VerifyCertificationPage(): JSX.Element {
  const params = Route.useParams();
  const hash = () => params().hash;
  const { auth } = useAuth();

  const [participantName, setParticipantName] = createSignal("");
  const [canvasEl, setCanvasEl] = createSignal<HTMLDivElement | undefined>(undefined);

  // 1. Verify query
  const verificationQuery = useQuery(() => certificationVerificationQueryOptions(hash()));

  const verified = () => verificationQuery().data?.valid === true;
  const cert = () => verificationQuery().data?.cert ?? null;
  const isLoading = () => verificationQuery().isLoading;
  const isError = () => verificationQuery().isError;
  const showCertificateArea = () => Boolean((verified() && cert()) || isLoading());

  const templateId = () => cert()?.template_id ?? verificationQuery().data?.template_id ?? "";

  // 2. Fetch template
  const templateQuery = useQuery(() => ({
    ...certificationTemplateQueryOptions(templateId()),
    enabled: Boolean(templateId()),
  }));

  const template = createMemo(() =>
    getCertificationTemplateOrDefault(templateQuery().data),
  );

  const canvas = createMemo(() => template().design_data?.canvas ?? DEFAULT_CERTIFICATE_CANVAS);

  // 3. Fetch actor profile for participant name
  createEffect(
    () => cert()?.user_id,
    (userId) => {
      if (!userId) {
        setParticipantName("");
        return;
      }
      void auth.getActorProfile(userId).then((res) => {
        if (res.success && res.data) {
          const p = asUniventsProfile(res.data.profile ?? {});
          setParticipantName(p.legalName || profileDisplayName(p) || "Participante");
        }
      }).catch(() => {
        setParticipantName("Participante");
      });
    },
  );

  // 4. Events & Editions
  const eventsQuery = useQuery(() => allPublicEventsQueryOptions());
  const editionQuery = useQuery(() =>
    editionLocationQueryOptions(
      cert()?.edition_id,
      (eventsQuery().data ?? []) as Array<{ id: string }>,
    ),
  );
  const [editionName, setEditionName] = createSignal<string>("");
  const [eventName, setEventName] = createSignal<string>("");
  const [locationName, setLocationName] = createSignal<string>("");

  createEffect(
    () => ({ edition: editionQuery().data, events: eventsQuery().data ?? [] }),
    ({ edition, events }) => {
      if (!edition) return;
      setEditionName(edition.name);
      setLocationName(edition.location_name ?? "");
      const event = events.find((item) => item.id === edition.event_id);
      setEventName(event?.full_name ?? "");
    },
  );

  // 5. Program query if certified for a program
  const programsQuery = useQuery(() => ({
    ...programsQueryOptions(cert()?.edition_id ?? ""),
    enabled: Boolean(cert()?.edition_id),
  }));

  const programName = createMemo(() => {
    const progId = cert()?.program_id;
    if (!progId) return null;
    const list = programsQuery().data ?? [];
    const matched = list.find((p) => p.id === progId);
    return matched?.name ?? null;
  });

  const occurrencesQuery = useQuery(() => ({
    ...occurrencesQueryOptions(cert()?.edition_id ?? ""),
    enabled: Boolean(cert()?.edition_id),
  }));
  const workloadHours = createMemo(() => {
    const milliseconds = ((occurrencesQuery().data ?? []) as Array<{ starts_at: string; ends_at: string }>).reduce(
      (total, occurrence) => total + new Date(occurrence.ends_at).getTime() - new Date(occurrence.starts_at).getTime(),
      0,
    );
    const hours = milliseconds / 3_600_000;
    return hours > 0
      ? Number.isInteger(hours)
        ? String(hours)
        : hours.toFixed(1).replace(".", ",")
      : "";
  });

  // 6. Certificate variables
  const variables = createMemo<CertificateVariableValues>(() => {
    const currentCert = cert();
    const origin = typeof window !== "undefined" ? window.location.origin : "";

    return {
      participant_name: participantName() || "Participante",
      event_name: eventName() || "Evento",
      edition_name: editionName() || "",
      activity_name: programName() ?? editionName() ?? "",
      participation_type: currentCert?.program_id ? "atividade" : "edição",
      location: locationName() || "",
      workload_hours: "",
      participation_date: "",
      certified_at: currentCert ? formatCertifiedAt(currentCert.issued_at) : "",
      cert_hash: hash(),
      verify_url: `${origin}/verify/${hash()}`,
    };
  });

  const statusTitle = createMemo(() => {
    if (isLoading()) return "Verificando certificado";
    if (isError()) return "Falha na verificação";
    if (verified()) return "Certificado autêntico";
    return "Certificado inválido";
  });

  const statusDescription = createMemo(() => {
    if (isLoading()) return "Consultando os dados de emissão e integridade.";
    if (isError()) return "Não foi possível consultar este certificado agora.";
    if (verified()) return "Emissão e integridade confirmadas.";
    return "Este documento não possui uma emissão válida.";
  });

  return (
    <main class="min-h-screen bg-background pb-24 text-foreground">

      <div class="mx-auto max-w-7xl px-4 py-6 md:px-6 md:py-8">
        <div class="grid gap-6 lg:grid-cols-2 lg:items-stretch">
          {/* Left Column: Certificate Preview Canvas */}
          <div class={showCertificateArea() ? "min-w-0" : "hidden"}>
            <Show
              when={!isLoading()}
              fallback={
                <div class="grid h-80 place-items-center rounded-xl border border-border/60 bg-muted/20">
                  <div class="flex flex-col items-center gap-2 text-muted-foreground">
                    <Loader2 class="size-7 animate-spin text-primary" />
                    <span class="text-xs">Carregando visualização...</span>
                  </div>
                </div>
              }
            >
              <Show
                when={verified() && cert()}
                fallback={
                  <div class="grid h-80 place-items-center rounded-xl border border-dashed border-destructive/30 bg-destructive/5 p-6 text-center">
                    <div class="flex flex-col items-center gap-3">
                      <div class="flex size-12 items-center justify-center rounded-full bg-destructive/10 text-destructive">
                        <FileX2 class="size-6" />
                      </div>
                      <div>
                        <h3 class="text-sm font-semibold text-foreground">
                          Visualização indisponível
                        </h3>
                        <p class="mt-1 max-w-sm text-xs text-muted-foreground">
                          Não foi possível exibir a prévia gráfica porque o certificado não é válido ou foi invalidado.
                        </p>
                      </div>
                    </div>
                  </div>
                }
              >
                <div
                  class="relative w-full overflow-hidden rounded-xl border border-border/60 bg-card shadow-md"
                  style={{ "aspect-ratio": `${canvas().width}/${canvas().height}` }}
                >
                  <CertificateTemplateStaticView
                    template={template()}
                    variables={variables()}
                    setCanvasRef={(el) => setCanvasEl(el)}
                    overlay={
                      <div class="absolute right-3 top-3 z-10 flex size-8 items-center justify-center rounded-full bg-background/90 p-1.5 shadow-sm ring-1 ring-foreground/10 backdrop-blur-sm">
                        <img
                          src="/logo-icon-default.svg"
                          alt="Univents"
                          class="size-full object-contain opacity-80"
                        />
                      </div>
                    }
                  />
                </div>

              </Show>
            </Show>
          </div>

          {/* Right Column: Verification Status & Metadata */}
          <div
            class={cn(
              "flex flex-col gap-2 overflow-hidden rounded-md border border-border/40 bg-card py-3 shadow-sm lg:self-start",
              !showCertificateArea() && "lg:col-span-2 lg:w-full lg:max-w-xl lg:justify-self-center",
            )}
          >
            {/* Status Header */}
            <div class="flex items-center gap-2 px-4">
              <div
                class={cn(
                  "flex size-8 shrink-0 items-center justify-center rounded-md",
                  isLoading()
                    ? "bg-muted text-muted-foreground"
                    : verified()
                      ? "bg-primary text-primary-foreground"
                      : "bg-destructive/10 text-destructive",
                )}
              >
                <Show when={isLoading()}>
                  <Loader2 class="size-4 animate-spin" />
                </Show>
                <Show when={!isLoading() && verified()}>
                  <BadgeCheck class="size-4" />
                </Show>
                <Show when={!isLoading() && !verified()}>
                  <FileX2 class="size-4" />
                </Show>
              </div>

              <div class="min-w-0 flex-1">
                <div>
                  <h2 class="text-sm leading-none font-semibold text-foreground">
                    {statusTitle()}
                  </h2>
                  <p class="mt-0.5 text-xs leading-tight text-muted-foreground">
                    {statusDescription()}
                  </p>
                </div>
              </div>
            </div>

            {/* Content Details */}
            <div class="flex flex-col gap-3 px-4 pb-0">
              <Show when={isLoading()}>
                <div class="h-1 overflow-hidden bg-muted">
                  <div class="h-full w-1/2 animate-pulse bg-primary" />
                </div>
              </Show>

              <Show when={!isLoading() && verified()}>
                <dl class="grid grid-cols-2 gap-x-6 border-y border-border/40">
                  <DetailItem label="Evento" value={eventName()} />
                  <DetailItem label="Edição" value={editionName()} />
                  <DetailItem label="Atividade" value={programName()} />
                  <DetailItem
                    label="Emitido em"
                    value={cert() ? formatCertifiedAt(cert()!.issued_at) : null}
                  />
                  <DetailItem
                    label="Carga horária"
                    value={workloadHours() ? `${workloadHours()} ${workloadHours() === "1" ? "hora" : "horas"}` : null}
                  />
                </dl>
                <div>
                  <p class="text-[10px] uppercase tracking-wider text-muted-foreground">
                    Código de verificação
                  </p>
                  <p class="mt-1 break-all font-mono text-[11px] text-muted-foreground">
                    {hash()}
                  </p>
                </div>
              </Show>

              <Show when={!isLoading() && !verified()}>
                <div class="space-y-3 border-l-2 border-destructive pl-3">
                  <p class="text-sm text-muted-foreground">
                    {isError()
                      ? "Tente novamente em instantes. Se o problema continuar, confirme se o link está completo."
                      : cert()?.invalid_reason ??
                      "Confira o código informado ou solicite uma nova emissão ao organizador."}
                  </p>
                  <p class="break-all font-mono text-[11px] text-muted-foreground/70">
                    {hash()}
                  </p>
                </div>
              </Show>

              <Show when={verified() && cert()}>
                <div class="border-t border-border/60 pt-4">
                  <CertificateDownloadButtons
                    canvasElement={() => canvasEl()}
                    templateName="Certificado"
                    prominent
                  />
                </div>
              </Show>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
