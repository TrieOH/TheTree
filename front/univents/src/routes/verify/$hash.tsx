import { useQueries, useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { BadgeCheck, FileX2, Hash, Loader2, ShieldCheck } from "lucide-react";
import { useMemo } from "react";
import { allPublicActivitiesQueryOptions } from "@/features/activities/api";
import type { ActivityI } from "@/features/activities/model";
import {
  certificationTemplateQueryOptions,
  verifyCertificationHashFn,
} from "@/features/certifications/api";
import { certificationKeys } from "@/features/certifications/api/query-keys";
import {
  CertificateTemplateStaticView,
  CertViewer,
} from "@/features/certifications/ui/CertViewer";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import type { EditionI } from "@/features/editions/model";
import { allPublicEventsQueryOptions } from "@/features/events/api";
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

function getOrigin() {
  if (window?.location?.origin) return window.location.origin;

  return "http://localhost:3002";
}

type VerifiedTemplateSectionProps = {
  hash: string;
  templateQuery: {
    query: ReturnType<typeof certificationTemplateQueryOptions>;
    eventId: string;
    editionId: string;
  };
  payload: {
    target_type: "edition" | "activity";
    target_id: string;
    certified_at: string;
  };
  activityLookup: Map<string, ActivityI>;
  editionLookup: Map<string, EditionI>;
};

function VerifiedTemplateSection({
  hash,
  templateQuery,
  payload,
  activityLookup,
  editionLookup,
}: VerifiedTemplateSectionProps) {
  const { data: templateData } = useQuery(templateQuery.query);

  const variables = useMemo(() => {
    const activity =
      payload.target_type === "activity"
        ? (activityLookup.get(payload.target_id) ?? null)
        : null;
    const editionId =
      payload.target_type === "edition"
        ? payload.target_id
        : (activity?.edition_id ?? null);
    const edition = editionId ? (editionLookup.get(editionId) ?? null) : null;

    return {
      activity_name:
        payload.target_type === "edition"
          ? (edition?.name ?? "")
          : (activity?.title ?? edition?.name ?? ""),
      certified_at: formatCertifiedAt(payload.certified_at),
      cert_hash: hash,
      verify_url: `${getOrigin()}/verify/${hash}`,
    };
  }, [activityLookup, editionLookup, hash, payload]);

  if (!templateData) return null;

  return (
    <div className="mb-6 overflow-hidden rounded-xl border bg-card shadow-sm">
      <div className="flex items-center justify-between gap-3 border-b px-4 py-3">
        <div className="min-w-0">
          <p className="truncate text-sm font-medium">{templateData.title}</p>
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

  const verified = data?.is_verified === true;
  const payload = data ?? null;

  const { data: events = [] } = useQuery(allPublicEventsQueryOptions());
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
      ...allPublicActivitiesQueryOptions(edition.event_id, edition.id),
      enabled: !!edition.event_id,
    })),
  });

  const activityLookup = useMemo(() => {
    const activities = new Map<string, ActivityI>();
    activityQueries.forEach((query) => {
      for (const activity of (query.data ?? []) as ActivityI[]) {
        activities.set(activity.id, activity);
      }
    });
    return activities;
  }, [activityQueries]);

  const templateQuery = useMemo(() => {
    if (!payload) return null;

    const activity =
      payload.target_type === "activity"
        ? (activityLookup.get(payload.target_id) ?? null)
        : null;
    const editionId =
      payload.target_type === "edition"
        ? payload.target_id
        : (activity?.edition_id ?? null);
    const edition = editionId ? (editionLookup.get(editionId) ?? null) : null;
    const templateId = activity?.certification_template_id ?? null;

    if (!templateId || !edition?.event_id || !edition?.id) return null;

    return {
      query: certificationTemplateQueryOptions(
        edition.event_id,
        edition.id,
        templateId,
      ),
      eventId: edition.event_id,
      editionId: edition.id,
    };
  }, [activityLookup, editionLookup, payload]);

  return (
    <main className="min-h-screen bg-background">
      <section className="border-b border-border/60 bg-linear-to-b from-muted/40 via-background to-background">
        <div className="mx-auto max-w-3xl px-4 py-10 md:px-6 md:py-14">
          <div className="space-y-3">
            <div className="flex items-center gap-2 text-xs uppercase tracking-[0.24em] text-muted-foreground">
              <ShieldCheck className="size-4" />
              Verificação pública
            </div>
            <h1 className="text-3xl font-semibold tracking-tight">
              Certificado
            </h1>
            <p className="max-w-2xl text-sm text-muted-foreground">
              Validação pública do certificado usando o hash da URL.
            </p>
          </div>
        </div>
      </section>

      <div className="mx-auto max-w-3xl px-4 py-6 md:px-6 md:py-8">
        {templateQuery && payload && (
          <VerifiedTemplateSection
            hash={hash}
            templateQuery={templateQuery}
            payload={payload}
            activityLookup={activityLookup}
            editionLookup={editionLookup}
          />
        )}

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

                <div className="grid gap-3 sm:grid-cols-2">
                  <div className="rounded-2xl border bg-muted/20 p-4">
                    <p className="text-[10px] uppercase tracking-wider text-muted-foreground">
                      Usuário
                    </p>
                    <p className="mt-1 font-mono text-sm break-all">
                      {payload?.user_id}
                    </p>
                  </div>
                  <div className="rounded-2xl border bg-muted/20 p-4">
                    <p className="text-[10px] uppercase tracking-wider text-muted-foreground">
                      Tipo
                    </p>
                    <p className="mt-1 text-sm font-medium capitalize">
                      {payload?.target_type}
                    </p>
                  </div>
                </div>

                <div className="rounded-2xl border bg-muted/20 p-4">
                  <p className="text-[10px] uppercase tracking-wider text-muted-foreground">
                    Target
                  </p>
                  <p className="mt-1 font-mono text-sm break-all">
                    {payload?.target_id}
                  </p>
                </div>

                <Separator />

                <div className="text-sm text-muted-foreground">
                  Emitido em{" "}
                  <span className="font-medium text-foreground">
                    {payload ? formatCertifiedAt(payload.certified_at) : "-"}
                  </span>
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

            {payload?.target_type && (
              <div className="flex flex-wrap gap-2">
                <Badge variant="outline" className="gap-1.5">
                  <Hash className="size-3.5" />
                  {payload.target_type}
                </Badge>
                <Badge
                  variant="secondary"
                  className="font-mono text-[10px] uppercase tracking-wider"
                >
                  {payload.is_verified ? "verified" : "unverified"}
                </Badge>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </main>
  );
}
