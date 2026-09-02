import { useQueries, useQuery } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { FileCheck2 } from "lucide-react";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import type { EditionI } from "@/features/editions/model";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import { programsQueryOptions } from "@/features/programs/api";
import type { ProgramI } from "@/features/programs/model";
import { Skeleton } from "@/shared/ui/shadcn/skeleton";
import {
  certificationsByUserQueryOptions,
  certificationTemplateQueryOptions,
} from "../api";
import { getCertificationTemplateOrDefault } from "../default-template";
import { DEFAULT_CERTIFICATE_CANVAS } from "../editor/constants";
import type { CertificationI } from "../model";
import { CertificateTemplateStaticView } from "./CertViewer";

export function UserCertificationsSection({
  participantName,
}: {
  participantName: string;
}) {
  const certifications = useQuery(certificationsByUserQueryOptions());
  const { data: events = [] } = useQuery(allPublicEventsQueryOptions());
  const editionQueries = useQueries({
    queries: events.map((event) => allPublicEditionsQueryOptions(event.id)),
  });
  const editions = editionQueries.flatMap((query) => query.data ?? []);
  const programQueries = useQueries({
    queries: editions.map((edition) => programsQueryOptions(edition.id)),
  });
  const programs = programQueries.flatMap((query) => query.data ?? []);

  if (certifications.isLoading) {
    return (
      <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
        {[0, 1, 2].map((item) => (
          <Skeleton key={item} className="aspect-[1.35] rounded-xl" />
        ))}
      </div>
    );
  }

  if (!certifications.data?.length) {
    return (
      <div className="grid min-h-64 place-items-center rounded-xl border border-dashed bg-card/40 p-8 text-center">
        <div>
          <FileCheck2 className="mx-auto size-9 text-muted-foreground" />
          <h2 className="mt-3 font-semibold">Nenhum certificado emitido</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            Seus certificados aparecerão aqui quando forem liberados.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-wrap items-start justify-center gap-3 sm:justify-start!">
      {certifications.data.map((certification) => {
        const edition = editions.find(
          (item) => item.id === certification.edition_id,
        );
        return (
          <CertificateProfileCard
            key={certification.id}
            certification={certification}
            participantName={participantName}
            eventName={
              events.find((event) => event.id === edition?.event_id)
                ?.full_name ?? "Evento"
            }
            edition={edition}
            program={programs.find(
              (program) => program.id === certification.program_id,
            )}
          />
        );
      })}
    </div>
  );
}

function CertificateProfileCard({
  certification,
  participantName,
  eventName,
  edition,
  program,
}: {
  certification: CertificationI;
  participantName: string;
  eventName: string;
  edition?: EditionI;
  program?: ProgramI;
}) {
  const { data: linkedTemplate } = useQuery({
    ...certificationTemplateQueryOptions(certification.template_id ?? ""),
    enabled: Boolean(certification.template_id),
  });
  const template = getCertificationTemplateOrDefault(linkedTemplate);
  const scopeName = program?.name ?? edition?.name ?? "Participação";
  const canvas = template.design_data.canvas ?? DEFAULT_CERTIFICATE_CANVAS;
  const width = 320;
  const height = (width * canvas.height) / canvas.width;
  const navigate = useNavigate();

  const openCertificate = () => {
    void navigate({
      to: "/verify/$hash",
      params: { hash: certification.verification_hash },
    });
  };

  const handleKeyDown = (event: React.KeyboardEvent<HTMLDivElement>) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      openCertificate();
    }
  };

  return (
    <div
      role="link"
      tabIndex={0}
      onClick={openCertificate}
      onKeyDown={handleKeyDown}
      className="block max-w-full cursor-pointer transition hover:-translate-y-0.5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
      style={{ width, height, maxWidth: "100%" }}
      aria-label={`Abrir certificado de ${scopeName}`}
    >
      <CertificateTemplateStaticView
        template={template}
        variables={{
          participant_name: participantName,
          event_name: eventName,
          edition_name: edition?.name ?? "",
          activity_name: scopeName,
          participation_type: program ? "atividade" : "edição",
          location: edition?.location_name ?? "",
          certified_at: new Date(certification.issued_at).toLocaleDateString(
            "pt-BR",
          ),
          cert_hash: certification.verification_hash,
          verify_url: `${window.location.origin}/verify/${certification.verification_hash}`,
        }}
      />
    </div>
  );
}
