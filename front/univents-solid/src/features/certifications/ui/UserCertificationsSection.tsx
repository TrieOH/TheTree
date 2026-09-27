import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { For, Show, createMemo } from "solid-js";

import { editionLocationQueryOptions } from "@/features/editions/api";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import {
  myParticipationsQueryOptions,
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";

import {
  certificationTemplateQueryOptions,
  myCertificationsQueryOptions,
} from "../api";
import { getCertificationTemplateOrDefault } from "../default-template";
import { DEFAULT_CERTIFICATE_CANVAS } from "../editor/constants";
import {
  buildCertificateVariables,
  type CertificateVariableValues,
  type OccurrenceTimeSpan,
} from "../lib/certificate-variables";
import type { CertificationI, CertificationTemplateI } from "../model";
import { CertificateTemplateStaticView } from "./CertViewer";

export interface UserCertificationsSectionProps {
  participantName: string;
  certifications?: CertificationI[];
}

export function UserCertificationsSection(
  props: UserCertificationsSectionProps,
): JSX.Element {
  const certsQuery = useQuery(() => ({
    ...myCertificationsQueryOptions(),
    enabled: props.certifications === undefined,
  }));

  const certifications = createMemo<CertificationI[]>(() => {
    if (props.certifications !== undefined) {
      return props.certifications;
    }
    return (certsQuery().data ?? []) as CertificationI[];
  });

  const isLoading = createMemo(() => {
    if (props.certifications !== undefined) return false;
    return certsQuery().isLoading;
  });

  const eventsQuery = useQuery(() => allPublicEventsQueryOptions());
  const events = createMemo<EventI[]>(() => (eventsQuery().data ?? []) as EventI[]);

  return (
    <Show
      when={!isLoading()}
      fallback={
        <div class="flex flex-wrap items-start justify-center gap-3 sm:justify-start">
          <For each={[1, 2, 3]}>
            {() => (
              <div
                class="animate-pulse rounded-lg border border-border/40 bg-muted/40"
                style={{
                  width: "320px",
                  height: "226px",
                  "max-width": "100%",
                }}
              />
            )}
          </For>
        </div>
      }
    >
      <Show
        when={certifications().length > 0}
        fallback={
          <div class="rounded-md border border-dashed border-border p-10 text-center">
            <h2 class="font-semibold">Você ainda não possui certificados</h2>
            <p class="mt-2 text-sm text-muted-foreground">
              Seus certificados aparecerão aqui quando forem emitidos.
            </p>
          </div>
        }
      >
        <div class="flex flex-wrap items-start justify-center gap-3 sm:justify-start">
          <For each={certifications()}>
            {(certification) => (
              <CertificateProfileCard
                certification={certification}
                participantName={props.participantName}
                events={events()}
              />
            )}
          </For>
        </div>
      </Show>
    </Show>
  );
}

function CertificateProfileCard(props: {
  certification: CertificationI;
  participantName: string;
  events: EventI[];
}): JSX.Element {
  const templateQuery = useQuery(() => ({
    ...certificationTemplateQueryOptions(props.certification.template_id ?? ""),
    enabled: Boolean(props.certification.template_id),
  }));

  const template = createMemo(() =>
    getCertificationTemplateOrDefault(templateQuery().data as CertificationTemplateI | undefined),
  );

  const editionQuery = useQuery(() =>
    editionLocationQueryOptions(props.certification.edition_id, props.events),
  );

  const edition = () => editionQuery().data;
  const editionName = () => edition()?.name ?? "";
  const locationName = () => edition()?.location_name ?? "";

  const eventName = createMemo(() => {
    const ed = edition();
    if (!ed) return "Univents";
    const matched = props.events.find((ev) => ev.id === ed.event_id);
    return matched?.full_name ?? "Univents";
  });

  const programsQuery = useQuery(() => ({
    ...programsQueryOptions(props.certification.edition_id),
    enabled: Boolean(props.certification.program_id),
  }));

  const programName = createMemo(() => {
    const progId = props.certification.program_id;
    if (!progId) return null;
    const list = programsQuery().data ?? [];
    const matched = list.find((p) => p.id === progId);
    return matched?.name ?? null;
  });

  const occurrencesQuery = useQuery(() => ({
    ...occurrencesQueryOptions(props.certification.edition_id),
    enabled: Boolean(props.certification.edition_id),
  }));

  const myParticipationsQuery = useQuery(() => ({
    ...myParticipationsQueryOptions(props.certification.edition_id),
    enabled: Boolean(props.certification.edition_id),
  }));

  const attendedOccurrences = createMemo(() => {
    const list = (myParticipationsQuery().data ?? []) as Array<{
      status: string;
      occurrence: OccurrenceTimeSpan;
      program: { id: string };
    }>;
    const progId = props.certification.program_id;
    return list
      .filter((p) => p.status === "attended" && (!progId || p.program.id === progId))
      .map((p) => p.occurrence);
  });

  const scopeName = createMemo(
    () => programName() ?? editionName() ?? "Certificado de Participação",
  );

  const canvas = createMemo(
    () => template().design_data?.canvas ?? DEFAULT_CERTIFICATE_CANVAS,
  );

  const width = 320;
  const height = createMemo(() => (width * canvas().height) / canvas().width);

  const variables = createMemo<CertificateVariableValues>(() => {
    const ed = edition();
    return buildCertificateVariables({
      participantName: props.participantName || "Participante",
      eventName: eventName(),
      editionName: editionName(),
      editionStartsAt: ed?.starts_at,
      editionEndsAt: ed?.ends_at,
      programName: programName(),
      programId: props.certification.program_id,
      locationName: locationName(),
      occurrences: (occurrencesQuery().data ?? []) as OccurrenceTimeSpan[],
      attendedOccurrences: attendedOccurrences(),
      issuedAt: props.certification.issued_at,
      certHash: props.certification.verification_hash,
    });
  });

  return (
    <Link
      to="/verify/$hash"
      params={{ hash: props.certification.verification_hash }}
      class="block max-w-full cursor-pointer transition hover:-translate-y-0.5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 rounded-lg overflow-hidden shadow-xs hover:shadow-md"
      style={{
        width: `${width}px` as `${number}px`,
        height: `${height()}px` as `${number}px`,
        "max-width": "100%",
      }}
      aria-label={`Abrir certificado de ${scopeName()}`}
    >
      <CertificateTemplateStaticView
        template={template()}
        variables={variables()}
      />
    </Link>
  );
}
