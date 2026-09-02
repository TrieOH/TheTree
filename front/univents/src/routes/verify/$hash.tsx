import { useQueries, useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { BadgeCheck, FileX2, Loader2 } from "lucide-react";
import type { RefObject } from "react";
import { useMemo, useRef } from "react";
import {
  certificationTemplateQueryOptions,
  verifyCertificationHashFn,
} from "@/features/certifications/api";
import { certificationKeys } from "@/features/certifications/api/query-keys";
import { getCertificationTemplateOrDefault } from "@/features/certifications/default-template";
import { DEFAULT_CERTIFICATE_CANVAS } from "@/features/certifications/editor/constants";
import type {
  CertificationI,
  CertificationTemplateI,
} from "@/features/certifications/model";
import {
  CertificateDownloadButtons,
  CertificateTemplateStaticView,
} from "@/features/certifications/ui/CertViewer";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import type { EditionI } from "@/features/editions/model";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import { profileKeys } from "@/features/profile/api/query-keys";
import { asUniventsProfile } from "@/features/profile/model/profile-data";
import {
  myParticipationsQueryOptions,
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import type { ProgramI } from "@/features/programs/model";
import { Logo } from "@/shared/ui/logo";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/shadcn/card";

export const Route = createFileRoute("/verify/$hash")({
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

function formatParticipationDate(value: string) {
  return new Date(value).toLocaleDateString("pt-BR", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  });
}

function formatHours(milliseconds: number) {
  const hours = milliseconds / 3_600_000;
  return Number.isInteger(hours)
    ? String(hours)
    : hours.toFixed(1).replace(".", ",");
}

function getOrigin() {
  if (window?.location?.origin) return window.location.origin;

  return "http://localhost:3002";
}

type VerifiedTemplateSectionProps = {
  hash: string;
  templateQuery: {
    query?: ReturnType<typeof certificationTemplateQueryOptions>;
    eventId: string;
    editionId: string;
  };
  payload: CertificationI;
  activityLookup: Map<string, ProgramI>;
  editionLookup: Map<string, EditionI>;
  eventLookup: Map<string, EventI>;
  canvasRef: RefObject<HTMLDivElement | null>;
  workloadHours: string;
  participationDate: string;
};

function VerifiedTemplateSection({
  hash,
  templateQuery,
  payload,
  activityLookup,
  editionLookup,
  eventLookup,
  canvasRef,
  workloadHours,
  participationDate,
}: VerifiedTemplateSectionProps) {
  const { auth } = useAuth();
  const { data: linkedTemplate } = useQuery<CertificationTemplateI | null>({
    queryKey: templateQuery.query?.queryKey ?? ["certifications", "default"],
    queryFn: async (context): Promise<CertificationTemplateI | null> => {
      const queryFn = templateQuery.query?.queryFn;
      if (!queryFn) return null;
      return (await queryFn(context as never)) as CertificationTemplateI;
    },
  });
  const templateData = getCertificationTemplateOrDefault(linkedTemplate);
  const canvas = templateData.design_data.canvas ?? DEFAULT_CERTIFICATE_CANVAS;
  const { data: participantName = "" } = useQuery({
    queryKey: profileKeys.certificateName(payload.user_id),
    queryFn: async () => {
      const response = await auth.getActorProfile(payload.user_id);
      if (!response.success || !response.data) return "";
      return asUniventsProfile(response.data.profile ?? {}).legalName ?? "";
    },
  });
  const variables = useMemo(() => {
    const activity =
      payload.program_id !== null
        ? (activityLookup.get(payload.program_id) ?? null)
        : null;
    const editionId =
      payload.program_id === null
        ? payload.edition_id
        : (activity?.edition_id ?? null);
    const edition = editionId ? (editionLookup.get(editionId) ?? null) : null;
    const event = edition ? (eventLookup.get(edition.event_id) ?? null) : null;
    return {
      participant_name: participantName,
      event_name: event?.full_name ?? "",
      edition_name: edition?.name ?? "",
      activity_name:
        payload.program_id === null
          ? (edition?.name ?? "")
          : (activity?.name ?? edition?.name ?? ""),
      participation_type: payload.program_id === null ? "edição" : "atividade",
      location: edition?.location_name ?? "",
      workload_hours: workloadHours,
      participation_date: participationDate,
      certified_at: formatCertifiedAt(payload.issued_at),
      cert_hash: hash,
      verify_url: `${getOrigin()}/verify/${hash}`,
    };
  }, [
    activityLookup,
    editionLookup,
    eventLookup,
    hash,
    participantName,
    participationDate,
    payload,
    workloadHours,
  ]);

  return (
    <div
      className="relative w-full"
      style={{ aspectRatio: `${canvas.width}/${canvas.height}` }}
    >
      <CertificateTemplateStaticView
        ref={canvasRef}
        template={templateData}
        variables={variables}
        overlay={
          <div className="absolute right-2 top-2 z-10 flex size-7 items-center justify-center rounded-full bg-background/90 p-1 shadow-sm ring-1 ring-foreground/10 backdrop-blur-sm">
            <Logo
              variant="icon"
              className="size-full opacity-75"
              imgClassName="object-contain"
            />
          </div>
        }
      />
    </div>
  );
}

function VerifyCertificationPage() {
  const { auth } = useAuth();
  const { hash } = Route.useParams();
  const { data, isLoading, isError } = useQuery({
    queryKey: certificationKeys.verification(hash),
    queryFn: () => verifyCertificationHashFn(hash),
    retry: false,
  });

  const verified = data?.valid === true;
  const payload = data?.cert ?? null;
  const canvasRef = useRef<HTMLDivElement>(null);

  const { data: events = [] } = useQuery(allPublicEventsQueryOptions());
  const eventLookup = useMemo(
    () => new Map(events.map((event) => [event.id, event])),
    [events],
  );
  const editionQueries = useQueries({
    queries: events.map((event) => ({
      ...allPublicEditionsQueryOptions(event.id),
      enabled: !!event.id,
    })),
  });

  const editionLookup = useMemo(() => {
    const editions = new Map<string, EditionI>();
    editionQueries.forEach((query) => {
      for (const edition of (query.data ?? []) as EditionI[]) {
        editions.set(edition.id, edition);
      }
    });
    return editions;
  }, [editionQueries]);

  const activityQueries = useQueries({
    queries: [...editionLookup.values()].map((edition) => ({
      ...programsQueryOptions(edition.id),
      enabled: !!edition.event_id,
    })),
  });

  const activityLookup = useMemo(() => {
    const programs = new Map<string, ProgramI>();
    activityQueries.forEach((query) => {
      for (const activity of query.data ?? []) {
        programs.set(activity.id, activity);
      }
    });
    return programs;
  }, [activityQueries]);

  const templateQuery = useMemo(() => {
    if (!payload) return null;

    const activity =
      payload.program_id !== null
        ? (activityLookup.get(payload.program_id) ?? null)
        : null;
    const editionId =
      payload.program_id === null
        ? payload.edition_id
        : (activity?.edition_id ?? null);
    const edition = editionId ? (editionLookup.get(editionId) ?? null) : null;
    const templateId = payload.template_id;

    if (!edition?.event_id || !edition?.id) return null;

    return {
      query: templateId
        ? certificationTemplateQueryOptions(templateId)
        : undefined,
      eventId: edition.event_id,
      editionId: edition.id,
    };
  }, [activityLookup, editionLookup, payload]);
  const activity = payload?.program_id
    ? activityLookup.get(payload.program_id)
    : undefined;
  const edition = payload
    ? editionLookup.get(activity?.edition_id ?? payload.edition_id)
    : undefined;
  const event = edition ? eventLookup.get(edition.event_id) : undefined;
  const { data: publicOccurrences = [] } = useQuery({
    ...occurrencesQueryOptions(edition?.id ?? ""),
    enabled: Boolean(edition?.id),
  });
  const isCertificateOwner = auth.profile()?.id === payload?.user_id;
  const { data: participationData } = useQuery({
    ...myParticipationsQueryOptions(edition?.id ?? ""),
    enabled: Boolean(edition?.id) && isCertificateOwner,
  });
  const attendedOccurrences = (participationData ?? [])
    .filter(
      (participation) =>
        participation.status === "attended" &&
        (!payload?.program_id ||
          participation.program.id === payload.program_id),
    )
    .map((participation) => participation.occurrence);
  const approximateOccurrences = publicOccurrences.filter(
    (occurrence) =>
      !payload?.program_id || occurrence.program_id === payload.program_id,
  );
  const workloadOccurrences =
    attendedOccurrences.length > 0
      ? attendedOccurrences
      : approximateOccurrences;
  const workload = workloadOccurrences.reduce(
    (total, occurrence) =>
      total +
      new Date(occurrence.ends_at).getTime() -
      new Date(occurrence.starts_at).getTime(),
    0,
  );
  const firstOccurrence = [...workloadOccurrences].sort(
    (a, b) => new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime(),
  )[0];
  const workloadHours = workload > 0 ? formatHours(workload) : "";
  const participationDate = firstOccurrence
    ? formatParticipationDate(firstOccurrence.starts_at)
    : "";
  const workloadIsApproximate =
    workloadOccurrences.length > 0 && attendedOccurrences.length === 0;
  const showCertificateArea = Boolean(templateQuery && payload) || isLoading;
  const status = isLoading
    ? {
        title: "Verificando certificado",
        description: "Consultando os dados de emissão e integridade.",
      }
    : isError
      ? {
          title: "Falha na verificação",
          description: "Não foi possível consultar este certificado agora.",
        }
      : verified
        ? {
            title: "Certificado autêntico",
            description: "Emissão e integridade confirmadas.",
          }
        : {
            title: "Certificado inválido",
            description: "Este documento não possui uma emissão válida.",
          };

  return (
    <main className="min-h-screen bg-background">
      <div className="mx-auto grid max-w-7xl gap-6 px-4 py-6 md:px-6 md:py-8 lg:grid-cols-2 lg:items-stretch">
        <div className={showCertificateArea ? "min-w-0" : "hidden"}>
          {templateQuery && payload ? (
            <VerifiedTemplateSection
              hash={hash}
              templateQuery={templateQuery}
              payload={payload}
              activityLookup={activityLookup}
              editionLookup={editionLookup}
              eventLookup={eventLookup}
              canvasRef={canvasRef}
              workloadHours={workloadHours}
              participationDate={participationDate}
            />
          ) : isLoading ? (
            <div className="grid h-full min-h-80 place-items-center rounded-md border bg-muted/30">
              <Loader2 className="size-6 animate-spin text-muted-foreground" />
            </div>
          ) : null}
        </div>
        <Card
          className={`gap-2! overflow-hidden rounded-md! border-border/40 bg-card py-3! shadow-sm lg:self-start ${showCertificateArea ? "" : "lg:col-span-2 lg:w-full lg:max-w-xl lg:justify-self-center"}`}
        >
          <CardHeader className="px-4 py-0">
            <div className="flex items-center gap-2">
              <div
                className={`flex size-8 shrink-0 items-center justify-center rounded-md ${isLoading ? "bg-muted text-muted-foreground" : verified ? "bg-primary text-primary-foreground" : "bg-destructive/10 text-destructive"}`}
              >
                {isLoading ? (
                  <Loader2 className="size-4 animate-spin" />
                ) : verified ? (
                  <BadgeCheck className="size-4" />
                ) : (
                  <FileX2 className="size-4" />
                )}
              </div>
              <div className="min-w-0">
                <CardTitle className="text-sm leading-none">
                  {status.title}
                </CardTitle>
                <CardDescription className="mt-0.5 text-xs leading-tight">
                  {status.description}
                </CardDescription>
              </div>
            </div>
          </CardHeader>

          <CardContent className="flex flex-col gap-3 px-4 pb-0">
            {isLoading ? (
              <div className="h-1 overflow-hidden bg-muted">
                <div className="h-full w-1/2 animate-pulse bg-primary" />
              </div>
            ) : verified ? (
              <>
                <dl className="grid grid-cols-2 gap-x-6 border-y border-border/40">
                  <CertificateDetail label="Evento" value={event?.full_name} />
                  <CertificateDetail label="Edição" value={edition?.name} />
                  {activity ? (
                    <CertificateDetail
                      label="Atividade"
                      value={activity.name}
                    />
                  ) : null}
                  <CertificateDetail
                    label="Emitido em"
                    value={
                      payload ? formatCertifiedAt(payload.issued_at) : undefined
                    }
                  />
                  <CertificateDetail
                    label="Carga horária"
                    value={
                      workloadHours
                        ? `${workloadHours} ${workloadHours === "1" ? "hora" : "horas"}${workloadIsApproximate ? " (aproximada)" : ""}`
                        : undefined
                    }
                  />
                </dl>
                <div>
                  <p className="text-[10px] font-medium uppercase tracking-wider text-muted-foreground">
                    Código de verificação
                  </p>
                  <p className="mt-1 break-all font-mono text-[11px] text-muted-foreground">
                    {hash}
                  </p>
                </div>
              </>
            ) : (
              <div className="space-y-3 border-l-2 border-destructive pl-3">
                <p className="text-sm text-muted-foreground">
                  {isError
                    ? "Tente novamente em instantes. Se o problema continuar, confirme se o link está completo."
                    : (payload?.invalid_reason ??
                      "Confira o código informado ou solicite uma nova emissão ao organizador.")}
                </p>
                <p className="break-all font-mono text-[11px] text-muted-foreground/70">
                  {hash}
                </p>
              </div>
            )}

            {verified && templateQuery && payload ? (
              <div className="flex justify-end border-t border-border/60 pt-4">
                <CertificateDownloadButtons
                  canvasRef={canvasRef}
                  templateName="Certificado"
                  prominent
                />
              </div>
            ) : null}
          </CardContent>
        </Card>
      </div>
    </main>
  );
}

function CertificateDetail({
  label,
  value,
}: {
  label: string;
  value?: string;
}) {
  if (!value) return null;
  return (
    <div className="py-3">
      <dt className="text-[10px] uppercase tracking-wider text-muted-foreground">
        {label}
      </dt>
      <dd className="mt-1 text-sm font-medium">{value}</dd>
    </div>
  );
}
