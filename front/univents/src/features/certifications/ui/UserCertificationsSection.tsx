import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { Award, CalendarDays, ExternalLink, FileCheck2 } from "lucide-react";
import { programsQueryOptions } from "@/features/programs/api";
import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import { buttonVariants } from "@/shared/ui/shadcn/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/shadcn/card";
import { Skeleton } from "@/shared/ui/shadcn/skeleton";
import { certificationsByUserQueryOptions } from "../api";

interface UserCertificationsSectionProps {
  editionId: string;
  title?: string;
  subtitle?: string;
}

export function UserCertificationsSection({
  editionId,
  title = "Meus certificados",
  subtitle = "Certificados emitidos nesta edição.",
}: UserCertificationsSectionProps) {
  const certificationsQuery = useQuery({
    ...certificationsByUserQueryOptions(),
  });
  const programsQuery = useQuery(programsQueryOptions(editionId));
  const programs = programsQuery.data ?? [];
  const activityNames = new Map(
    programs.map((program) => [program.id, program.name]),
  );
  const activityIds = new Set(activityNames.keys());
  const certifications = (certificationsQuery.data ?? []).filter(
    (certification) =>
      (certification.program_id === null &&
        certification.edition_id === editionId) ||
      (certification.program_id !== null &&
        activityIds.has(certification.program_id)),
  );
  const isLoading = certificationsQuery.isLoading || programsQuery.isLoading;

  return (
    <section className="space-y-4">
      <div>
        <h2 className="flex items-center gap-2 text-lg font-semibold">
          <Award className="size-5 text-primary" />
          {title}
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p>
      </div>

      {isLoading ? (
        <div className="grid gap-4 md:grid-cols-2">
          <Skeleton className="h-44 rounded-xl" />
          <Skeleton className="h-44 rounded-xl" />
        </div>
      ) : certifications.length === 0 ? (
        <div className="flex min-h-44 flex-col items-center justify-center rounded-xl border border-dashed bg-muted/20 px-6 text-center">
          <FileCheck2 className="mb-3 size-8 text-muted-foreground" />
          <p className="text-sm font-medium">Nenhum certificado nesta edição</p>
          <p className="mt-1 max-w-md text-xs text-muted-foreground">
            Quando um certificado for emitido para a edição ou uma de suas
            atividades, ele aparecerá aqui.
          </p>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {certifications.map((certification) => {
            const targetName =
              certification.program_id === null
                ? "Certificado da edição"
                : (activityNames.get(certification.program_id ?? "") ??
                  "Certificado de atividade");

            return (
              <Card key={certification.id} className="overflow-hidden">
                <CardHeader className="border-b bg-muted/20 pb-4">
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0">
                      <CardTitle className="truncate text-base">
                        {targetName}
                      </CardTitle>
                      <CardDescription className="mt-1 flex items-center gap-1.5 text-xs">
                        <CalendarDays className="size-3.5" />
                        Emitido em{" "}
                        {formatCertificateDate(certification.issued_at)}
                      </CardDescription>
                    </div>
                    <Badge variant="secondary">
                      {certification.program_id === null
                        ? "Edição"
                        : "Atividade"}
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent className="space-y-4 p-4">
                  <div className="rounded-md bg-muted px-3 py-2">
                    <p className="text-[10px] font-medium tracking-wide text-muted-foreground uppercase">
                      Código de verificação
                    </p>
                    <p className="mt-1 truncate font-mono text-xs">
                      {certification.verification_hash || "Código indisponível"}
                    </p>
                  </div>
                  {certification.verification_hash ? (
                    <Link
                      to="/verify/$hash"
                      params={{ hash: certification.verification_hash }}
                      className={cn(
                        buttonVariants({ variant: "outline", size: "sm" }),
                        "w-full",
                      )}
                    >
                      Verificar certificado
                      <ExternalLink className="size-3.5" />
                    </Link>
                  ) : null}
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </section>
  );
}

function formatCertificateDate(value: string) {
  return new Intl.DateTimeFormat("pt-BR", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  }).format(new Date(value));
}
