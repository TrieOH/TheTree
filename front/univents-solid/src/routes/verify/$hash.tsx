import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { cn } from "@trieoh/ui-solid";
import { Show, createMemo, createSignal } from "solid-js";

import BadgeCheckIcon from "~icons/lucide/badge-check";
import FileX2Icon from "~icons/lucide/file-x-2";
import Loader2Icon from "~icons/lucide/loader-2";

import {
  certificationTemplateQueryOptions,
  certificationVerificationQueryOptions,
} from "@/features/certifications/api";
import { getCertificationTemplateOrDefault } from "@/features/certifications/default-template";
import { DEFAULT_CERTIFICATE_CANVAS } from "@/features/certifications/editor/constants";
import {
  buildCertificateVariables,
  formatCertifiedAt,
  formatWorkloadLabel,
  type CertificateVariableValues,
  type OccurrenceTimeSpan,
} from "@/features/certifications/lib/certificate-variables";
import {
  CertificateDownloadButtons,
  CertificateTemplateStaticView,
} from "@/features/certifications/ui/CertViewer";
import { editionLocationQueryOptions } from "@/features/editions/api";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import { ownProfileQueryOptions } from "@/features/profile/api";
import { asUniventsProfile, profileDisplayName } from "@/features/profile/model/profile-data";
import {
  myParticipationsQueryOptions,
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";

const BadgeCheck = BadgeCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const FileX2 = FileX2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute("/verify/$hash")({
  head: () => ({
    meta: [{ title: "Verificação de Autenticidade de Certificado - Univents" }],
  }),
  component: VerifyCertificationPage,
});

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

  const [canvasEl, setCanvasEl] = createSignal<HTMLDivElement | undefined>(undefined);

  // 1. Verify query
  const verificationQuery = useQuery(() => certificationVerificationQueryOptions(hash()));

  const verified = () => verificationQuery().data?.valid === true;
  const cert = () => verificationQuery().data?.cert ?? null;
  const isLoading = () => verificationQuery().isLoading;
  const isError = () => verificationQuery().isError;

  const isNotFound = () => {
    const err = verificationQuery().error as { envelope?: { code?: number }; code?: number } | null;
    return err?.envelope?.code === 404 || err?.code === 404;
  };

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

  // 3. The participant's name, through the profile cache.
  const participantQuery = useQuery(() =>
    ownProfileQueryOptions(cert()?.user_id ?? "", auth),
  );

  const participantName = createMemo(() => {
    const participant = participantQuery().data;
    if (!participant) return participantQuery().isError ? "Participante" : "";
    const profile = asUniventsProfile(participant.profile ?? {});
    return profile.legalName || profileDisplayName(profile) || "Participante";
  });

  // 4. Events & Editions
  const eventsQuery = useQuery(() => allPublicEventsQueryOptions());
  const editionQuery = useQuery(() =>
    editionLocationQueryOptions(
      cert()?.edition_id,
      (eventsQuery().data ?? []) as Array<{ id: string }>,
    ),
  );
  const editionName = () => editionQuery().data?.name ?? "";
  const locationName = () => editionQuery().data?.location_name ?? "";
  const eventName = createMemo(() => {
    const edition = editionQuery().data;
    if (!edition) return "";
    return (eventsQuery().data ?? []).find((item) => item.id === edition.event_id)?.full_name ?? "";
  });

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

  const isCertificateOwner = () => {
    const profileId = auth.profile?.()?.id;
    const certUserId = cert()?.user_id;
    return Boolean(profileId && certUserId && profileId === certUserId);
  };

  const myParticipationsQuery = useQuery(() => ({
    ...myParticipationsQueryOptions(cert()?.edition_id ?? ""),
    enabled: Boolean(cert()?.edition_id && isCertificateOwner()),
  }));

  // The certificate is drawn from nine pieces of data (template, participant,
  // event, edition, programs, occurrences, attendance...). Painting it as each
  // arrives would redraw the artwork and reflow the details several times, so
  // it, and the details beside it, wait until every request has settled. A
  // failed one just leaves its field out.
  const detailsReady = createMemo(
    () =>
      !eventsQuery().isLoading &&
      ![
        templateQuery(),
        participantQuery(),
        editionQuery(),
        programsQuery(),
        occurrencesQuery(),
        myParticipationsQuery(),
      ].some((query) => query.isLoading),
    { name: "VerifyCertificationPage.detailsReady" },
  );

  const attendedOccurrences = createMemo(() => {
    const currentCert = cert();
    const list = (myParticipationsQuery().data ?? []) as Array<{
      status: string;
      occurrence: OccurrenceTimeSpan;
      program: { id: string };
    }>;
    return list
      .filter(
        (p) =>
          p.status === "attended" &&
          (!currentCert?.program_id || p.program.id === currentCert.program_id),
      )
      .map((p) => p.occurrence);
  });

  const workloadIsApproximate = createMemo(() => {
    if (cert()?.program_id) return false;
    const hasOccurrences =
      (occurrencesQuery().data ?? []).length > 0 ||
      Boolean(editionQuery().data?.starts_at);
    return hasOccurrences && attendedOccurrences().length === 0;
  });

  // 6. Certificate variables
  const variables = createMemo<CertificateVariableValues>(() => {
    const currentCert = cert();
    const ed = editionQuery().data;
    return buildCertificateVariables({
      participantName: participantName() || "Participante",
      eventName: eventName() || "Evento",
      editionName: editionName() || "",
      editionStartsAt: ed?.starts_at,
      editionEndsAt: ed?.ends_at,
      programName: programName(),
      programId: currentCert?.program_id,
      locationName: locationName() || "",
      occurrences: (occurrencesQuery().data ?? []) as OccurrenceTimeSpan[],
      attendedOccurrences: attendedOccurrences(),
      issuedAt: currentCert?.issued_at,
      certHash: hash(),
    });
  });

  const statusTitle = createMemo(() => {
    if (isLoading()) return "Verificando certificado";
    if (isNotFound()) return "Certificado não encontrado";
    if (isError()) return "Falha na verificação";
    if (verified()) return "Certificado autêntico";
    return "Certificado inválido";
  });

  const statusDescription = createMemo(() => {
    if (isLoading()) return "Consultando os dados de emissão e integridade.";
    if (isNotFound()) return "Nenhum certificado foi encontrado com o código informado.";
    if (isError()) return "Não foi possível consultar este certificado agora.";
    if (verified()) return "Emissão e integridade confirmadas.";
    return "Este documento não possui uma emissão válida.";
  });

  const issuedAtFormatted = createMemo(() => {
    const c = cert();
    return c ? formatCertifiedAt(c.issued_at) : null;
  });

  return (
    <div class="min-h-screen bg-background pb-24 text-foreground">
      <div class="mx-auto max-w-7xl px-4 py-6 md:px-6 md:py-8">
        <div class="grid gap-6 lg:grid-cols-2 lg:items-stretch">
          {/* Left Column: Certificate Preview Canvas */}
          <div class={showCertificateArea() ? "min-w-0" : "hidden"}>
            <Show
              when={!isLoading() && (detailsReady() || !verified())}
              fallback={
                <div class="grid h-80 place-items-center rounded-xl border border-border/60 bg-muted/20">
                  <div class="flex flex-col items-center gap-2 text-muted-foreground">
                    <Loader2 class="size-7 animate-spin text-primary" />
                    <span class="text-xs">Carregando visualização...</span>
                  </div>
                </div>
              }
            >
              <Show when={verified() && cert()}>
                <div
                  ref={setCanvasEl}
                  class="relative w-full rounded-xl overflow-hidden shadow-md border border-border/60 bg-white"
                  style={{ "aspect-ratio": `${canvas().width} / ${canvas().height}` }}
                >
                  <CertificateTemplateStaticView
                    template={template()}
                    variables={variables()}
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
                    : isNotFound() || (!verified() && !isLoading())
                      ? "bg-destructive/10 text-destructive"
                      : "bg-primary text-primary-foreground",
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
                <h2 class="text-sm leading-none font-semibold text-foreground">
                  {statusTitle()}
                </h2>
                <p class="mt-0.5 text-xs leading-tight text-muted-foreground">
                  {statusDescription()}
                </p>
              </div>
            </div>

            {/* Content Details */}
            <div class="flex flex-col gap-3 px-4 pb-0">
              <Show when={isLoading()}>
                <div class="h-1 overflow-hidden bg-muted">
                  <div class="h-full w-1/2 animate-pulse bg-primary" />
                </div>
              </Show>

              <Show when={!isLoading() && verified() && !detailsReady()}>
                <div class="h-40 animate-pulse border-y border-border/40 bg-muted/30" />
              </Show>

              <Show when={!isLoading() && verified() && detailsReady()}>
                <dl class="grid grid-cols-2 gap-x-6 border-y border-border/40">
                  <DetailItem label="Evento" value={eventName()} />
                  <DetailItem label="Edição" value={editionName()} />
                  <DetailItem label="Atividade" value={programName()} />
                  <DetailItem
                    label="Emitido em"
                    value={issuedAtFormatted()}
                  />
                  <DetailItem
                    label="Carga horária"
                    value={
                      formatWorkloadLabel(variables().workload_hours, {
                        approximate: workloadIsApproximate(),
                      }) || null
                    }
                  />
                  <DetailItem
                    label="Data de realização"
                    value={variables().participation_date || null}
                  />
                </dl>
              </Show>

              <Show when={!isLoading() && !verified()}>
                <div class="p-3 text-xs border-l-2 border-destructive bg-destructive/5 space-y-1">
                  <p class="font-medium text-destructive">
                    {isNotFound() ? "Código não localizado" : "Falha na verificação"}
                  </p>
                  <p class="text-muted-foreground">
                    {isNotFound()
                      ? "Confira se o código informado está completo ou solicite uma nova emissão ao organizador."
                      : isError()
                        ? "Tente novamente em instantes. Se o problema continuar, confirme se o link está completo."
                        : cert()?.invalid_reason ??
                        "Confira o código informado ou solicite uma nova emissão ao organizador."}
                  </p>
                </div>
              </Show>

              <div class="space-y-1">
                <p class="text-[10px] uppercase tracking-wider text-muted-foreground">
                  Código de verificação
                </p>
                <p class="break-all font-mono text-xs text-muted-foreground">
                  {hash()}
                </p>
              </div>

              <Show when={verified() && cert()}>
                <div class="border-t border-border/40 pt-3">
                  <CertificateDownloadButtons
                    canvasElement={canvasEl}
                    templateName="Certificado"
                    prominent
                  />
                </div>
              </Show>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
