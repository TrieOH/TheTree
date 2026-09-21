import type { JSX } from "@solidjs/web";
import { useQuery } from "@trieoh/front-core-solid";
import { EmptyState } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal } from "solid-js";
import FileCheck2Icon from "~icons/lucide/file-check-2";

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
import { CertViewer, CertificateTemplateStaticView } from "./CertViewer";

const FileCheck2 = FileCheck2Icon as unknown as (props: { class?: string }) => JSX.Element;

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

  const [activeCert, setActiveCert] = createSignal<{
    cert: CertificationI;
    template: CertificationTemplateI;
    variables: CertificateVariableValues;
  } | null>(null);
  const [viewerOpen, setViewerOpen] = createSignal(false);

  const handleOpenCertificate = (
    cert: CertificationI,
    template: CertificationTemplateI,
    variables: CertificateVariableValues,
  ) => {
    setActiveCert({ cert, template, variables });
    setViewerOpen(true);
  };

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
          <EmptyState
            icon={<FileCheck2 class="size-6 text-foreground/70" />}
            title="Nenhum certificado emitido"
            description="Seus certificados aparecerão aqui quando forem liberados."
          />
        }
      >
        <div class="flex flex-wrap items-start justify-center gap-3 sm:justify-start">
          <For each={certifications()}>
            {(certification) => (
              <CertificateProfileCard
                certification={certification}
                participantName={props.participantName}
                events={events()}
                onOpen={handleOpenCertificate}
              />
            )}
          </For>
        </div>
      </Show>

      <Show when={activeCert()}>
        {(current) => (
          <CertViewer
            template={current().template}
            variables={current().variables}
            open={viewerOpen()}
            onOpenChange={setViewerOpen}
          />
        )}
      </Show>
    </Show>
  );
}

function CertificateProfileCard(props: {
  certification: CertificationI;
  participantName: string;
  events: EventI[];
  onOpen: (
    cert: CertificationI,
    template: CertificationTemplateI,
    variables: CertificateVariableValues,
  ) => void;
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

  const handleClick = () => {
    props.onOpen(props.certification, template(), variables());
  };

  const handleKeyDown = (event: KeyboardEvent) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      handleClick();
    }
  };

  return (
    <div
      role="link"
      tabindex={0}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
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
    </div>
  );
}

