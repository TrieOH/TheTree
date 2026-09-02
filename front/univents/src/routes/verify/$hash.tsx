import { useQueries, useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { BadgeCheck, FileX2, Hash, Loader2, ShieldCheck } from "lucide-react";
import { useMemo } from "react";
import {
  certificationTemplateQueryOptions,
  verifyCertificationHashFn,
} from "@/features/certifications/api";
import { certificationKeys } from "@/features/certifications/api/query-keys";
import { getCertificationTemplateOrDefault } from "@/features/certifications/default-template";
import type {
  CertificationI,
  CertificationTemplateI,
} from "@/features/certifications/model";
import {
  CertificateTemplateStaticView,
  CertViewer,
} from "@/features/certifications/ui/CertViewer";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import type { EditionI } from "@/features/editions/model";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import { profileKeys } from "@/features/profile/api/query-keys";
import { asUniventsProfile } from "@/features/profile/model/profile-data";
import {
  myParticipationsQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import type { ProgramI } from "@/features/programs/model";
import { Badge } from "@/shared/ui/shadcn/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/shadcn/card";
import { Separator } from "@/shared/ui/shadcn/separator";

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
};

function VerifiedTemplateSection({
  hash,
  templateQuery,
  payload,
  activityLookup,
  editionLookup,
  eventLookup,
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
  const { data: participantName = "" } = useQuery({
    queryKey: profileKeys.certificateName(payload.user_id),
    queryFn: async () => {
      const response = await auth.getActorProfile(payload.user_id);
      if (!response.success || !response.data) return "";
      return asUniventsProfile(response.data.profile ?? {}).legalName ?? "";
    },
  });
  const { data: participationData } = useQuery({
    ...myParticipationsQueryOptions(
      payload.program_id
        ? (activityLookup.get(payload.program_id)?.edition_id ?? "")
        : payload.edition_id,
    ),
    enabled: auth.profile()?.id === payload.user_id,
  });
  const myParticipations = participationData ?? [];

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
    const attended = myParticipations.filter(
      (participation) => participation.status === "attended",
    );
    const attendedForCertificate = payload.program_id
      ? attended.filter(
          (participation) => participation.program.id === payload.program_id,
        )
      : attended;
    const startDate = attendedForCertificate[0]?.occurrence.starts_at;
    const workload = attendedForCertificate.reduce(
      (total, participation) =>
        total +
        new Date(participation.occurrence.ends_at).getTime() -
        new Date(participation.occurrence.starts_at).getTime(),
      0,
    );

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
      workload_hours:
        attendedForCertificate.length > 0 ? formatHours(workload) : "",
      participation_date: startDate ? formatParticipationDate(startDate) : "",
      certified_at: formatCertifiedAt(payload.issued_at),
      cert_hash: hash,
      verify_url: `${getOrigin()}/verify/${hash}`,
    };
  }, [
    activityLookup,
    myParticipations,
    editionLookup,
    eventLookup,
    hash,
    participantName,
    payload,
  ]);

  return (
    <div className="overflow-hidden rounded-xl border bg-card shadow-sm">
      <div className="flex items-center justify-between gap-3 border-b px-4 py-3">
        <div className="min-w-0">
          <p className="truncate text-sm font-medium">{templateData.name}</p>
          <p className="text-xs text-muted-foreground">
            Certificado verificado
          </p>
        </div>
        <CertViewer
          template={templateData}
          variables={variables}
          triggerLabel="Ampliar e baixar"
        />
      </div>
      <div className="h-[min(32rem,60vw)] min-h-64 bg-muted/40 p-3 sm:p-5">
        <CertificateTemplateStaticView
          template={templateData}
          variables={variables}
        />
      </div>
    </div>
  );
}

function VerifyCertificationPage() {
  const { hash } = Route.useParams();
  const { data, isLoading, isError } = useQuery({
    queryKey: certificationKeys.verification(hash),
    queryFn: () => verifyCertificationHashFn(hash),
    retry: false,
  });

  const verified = data?.valid === true;
  const payload = data?.cert ?? null;

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
  const certificateName = activity?.name ?? edition?.name ?? "Certificado";

  return (
    <main className="min-h-screen bg-background">
      <section className="border-b border-border/60 bg-linear-to-b from-muted/40 via-background to-background">
        <div className="mx-auto max-w-6xl px-4 py-10 md:px-6 md:py-14">
          <div className="space-y-3">
            <div className="flex items-center gap-2 text-xs uppercase tracking-[0.24em] text-muted-foreground">
              <ShieldCheck className="size-4" />
              Verificação pública
            </div>
            <h1 className="text-3xl font-semibold tracking-tight">
              {certificateName}
            </h1>
            <p className="max-w-2xl text-sm text-muted-foreground">
              {event?.full_name
                ? `Certificado emitido por ${event.full_name}.`
                : "Validação pública de autenticidade do certificado."}
            </p>
          </div>
        </div>
      </section>

      <div className="mx-auto grid max-w-6xl gap-6 px-4 py-6 md:px-6 md:py-8 lg:grid-cols-[minmax(0,1fr)_22rem] lg:items-start">
        <div>
          {templateQuery && payload ? (
            <VerifiedTemplateSection
              hash={hash}
              templateQuery={templateQuery}
              payload={payload}
              activityLookup={activityLookup}
              editionLookup={editionLookup}
              eventLookup={eventLookup}
            />
          ) : isLoading ? (
            <div className="grid min-h-80 place-items-center rounded-xl border bg-muted/30">
              <Loader2 className="size-6 animate-spin text-muted-foreground" />
            </div>
          ) : null}
        </div>
        <Card className="overflow-hidden border-border/60 bg-card shadow-sm">
          <CardHeader className="border-b border-border/60">
            <CardTitle className="flex items-center gap-2 text-base">
              <BadgeCheck className="size-4 text-primary" />
              Status da verificação
            </CardTitle>
            <CardDescription>
              Hash: <span className="font-mono">{hash}</span>
            </CardDescription>
          </CardHeader>

          <CardContent className="space-y-5 p-5">
            {isLoading ? (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Loader2 className="size-4 animate-spin" />
                Verificando certificado...
              </div>
            ) : verified ? (
              <>
                <div className="flex items-center gap-2 rounded-2xl border border-emerald-500/20 bg-emerald-500/10 px-4 py-3 text-emerald-700">
                  <BadgeCheck className="size-5" />
                  Certificado verificado com sucesso
                </div>

                <dl className="divide-y rounded-xl border bg-muted/20 px-4">
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
                </dl>

                <Separator />

                <div>
                  <p className="text-xs font-medium text-muted-foreground">
                    Código de verificação
                  </p>
                  <p className="mt-1 break-all font-mono text-xs">{hash}</p>
                </div>
              </>
            ) : (
              <div className="space-y-4">
                <div className="flex items-center gap-2 rounded-2xl border border-destructive/20 bg-destructive/10 px-4 py-3 text-destructive">
                  <FileX2 className="size-5" />
                  Não foi possível validar este certificado
                </div>
                <p className="text-sm text-muted-foreground">
                  {isError
                    ? "A verificação falhou ao consultar o certificado."
                    : "O hash informado não corresponde a um certificado válido ou já não está ativo."}
                </p>
              </div>
            )}

            {payload && (
              <div className="flex flex-wrap gap-2">
                <Badge variant="outline" className="gap-1.5">
                  <Hash className="size-3.5" />
                  {payload.program_id === null ? "Edição" : "Atividade"}
                </Badge>
                <Badge
                  variant="secondary"
                  className="font-mono text-[10px] uppercase tracking-wider"
                >
                  {data?.valid ? "Verificado" : "Não verificado"}
                </Badge>
              </div>
            )}
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
