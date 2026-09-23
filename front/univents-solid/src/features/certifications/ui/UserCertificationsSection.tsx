import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { For, Show, createMemo } from "solid-js";

import { allPublicEventsQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import {
  certificationTemplateQueryOptions,
  myCertificationsQueryOptions,
} from "../api";
import { getCertificationTemplateOrDefault } from "../default-template";
import { DEFAULT_CERTIFICATE_CANVAS } from "../editor/constants";
import type { CertificateVariableValues } from "../editor/variables";
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

  const eventName = createMemo(() => {
    const matched = props.events.find(
      (ev) => ev.id === props.certification.edition_id,
    );
    return matched?.full_name ?? "Univents";
  });

  const scopeName = createMemo(() => "Certificado de Participação");

  const canvas = createMemo(
    () => template().design_data?.canvas ?? DEFAULT_CERTIFICATE_CANVAS,
  );

  const width = 320;
  const height = createMemo(() => (width * canvas().height) / canvas().width);

  const variables = createMemo<CertificateVariableValues>(() => {
    const origin = typeof window !== "undefined" ? window.location.origin : "";
    return {
      participant_name: props.participantName || "Participante",
      event_name: eventName(),
      edition_name: "Edição",
      activity_name: scopeName(),
      participation_type: props.certification.program_id ? "atividade" : "edição",
      location: "",
      certified_at: new Date(props.certification.issued_at).toLocaleDateString("pt-BR"),
      cert_hash: props.certification.verification_hash,
      verify_url: `${origin}/verify/${props.certification.verification_hash}`,
    };
  });

  return (
    <Link
      to="/verify/$hash"
      params={{ hash: props.certification.verification_hash }}
      class="block max-w-full cursor-pointer transition hover:-translate-y-0.5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 rounded-lg overflow-hidden shadow-xs hover:shadow-md"
      style={{
        width: `${width}px`,
        height: `${height()}px`,
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
