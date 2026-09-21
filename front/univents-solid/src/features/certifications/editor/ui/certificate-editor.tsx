import type { JSX } from "@solidjs/web";
import { Link, useNavigate } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { Show, createEffect, createSignal, onCleanup } from "solid-js";
import ArrowLeftIcon from "~icons/lucide/arrow-left";
import Loader2Icon from "~icons/lucide/loader-2";
import MonitorIcon from "~icons/lucide/monitor";
import SaveIcon from "~icons/lucide/save";

import { Button } from "@trieoh/ui-solid";
import { toast } from "@/shared/ui/toast";
import { allSignaturesQueryOptions } from "@/features/signatures/api";
import type { SignatureI } from "@/features/signatures/model";
import { DEFAULT_CERTIFICATION_TEMPLATE } from "../../default-template";
import { certificationTemplateQueryOptions } from "../../api";
import {
  useCreateCertificationTemplateMutation,
  useUpdateCertificationTemplateMutation,
} from "../../api/mutations";
import { certificationTemplateCreateSchema } from "../../model";
import { certificateEditorActions } from "../store";
import { uploadCertificateAssets } from "../upload-assets";
import { CertificateCanvas } from "./certificate-canvas";
import { CertificatePropertiesPanel } from "./certificate-properties-panel";
import { CertificateTextToolbar } from "./certificate-text-toolbar";
import { CertificateToolsSidebar } from "./certificate-tools-sidebar";

const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Monitor = MonitorIcon as unknown as (props: { class?: string }) => JSX.Element;
const Save = SaveIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface CertificateEditorProps {
  eventId: string;
  editionId: string;
  templateId?: string;
  duplicate?: boolean;
}

export function CertificateEditor(props: CertificateEditorProps): JSX.Element {
  const navigate = useNavigate();

  const signaturesQuery = useQuery(() =>
    allSignaturesQueryOptions(props.editionId),
  );

  const templateQuery = useQuery(() => ({
    ...certificationTemplateQueryOptions(props.templateId ?? ""),
    enabled: Boolean(props.templateId),
  }));

  createEffect(
    () => ({
      data: templateQuery().data,
      duplicate: props.duplicate,
      templateId: props.templateId,
    }),
    ({ data, duplicate, templateId }) => {
      if (data) {
        certificateEditorActions.loadDraft(data);
        if (duplicate) {
          certificateEditorActions.setName(`${data.name} (cópia)`);
        }
      } else if (!templateId) {
        certificateEditorActions.loadDraft(DEFAULT_CERTIFICATION_TEMPLATE);
      } else {
        certificateEditorActions.reset();
      }
    },
  );

  createEffect(
    () => signaturesQuery().data,
    (signatures) => {
      if (signatures) {
        certificateEditorActions.setAvailableSignatures(
          (signatures as SignatureI[]).map((signature) => ({
            id: signature.id,
            name: signature.signatory_name,
            url: signature.image_url,
          })),
        );
      }
    },
  );

  const createMutation = useCreateCertificationTemplateMutation();
  const updateMutation = useUpdateCertificationTemplateMutation();
  const [saving, setSaving] = createSignal(false);

  const isSaving = () =>
    saving() ||
    createMutation.result().status === "pending" ||
    updateMutation.result().status === "pending";

  const handleSave = async () => {
    if (isSaving()) return;

    try {
      setSaving(true);
      const payload = certificateEditorActions.getDraft();

      const validation = certificationTemplateCreateSchema.safeParse(payload);
      if (!validation.success) {
        const errorMsg =
          validation.error.issues[0]?.message ?? "Dados inválidos";
        toast.error(errorMsg);
        setSaving(false);
        return;
      }

      const processedDraft = await uploadCertificateAssets(
        payload,
        props.eventId,
        props.editionId,
      );

      if (props.templateId && !props.duplicate) {
        await updateMutation.mutateAsync({
          editionId: props.editionId,
          templateId: props.templateId,
          data: processedDraft,
        });
        toast.success("Certificado atualizado com sucesso!");
      } else {
        await createMutation.mutateAsync({
          editionId: props.editionId,
          data: processedDraft,
        });
        toast.success("Certificado criado com sucesso!");
      }

      void navigate({
        to: "/admin/events/$eventId/editions/$editionId/certifications",
        params: { eventId: props.eventId, editionId: props.editionId },
      });
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Erro ao salvar template";
      toast.error(message);
    } finally {
      setSaving(false);
    }
  };

  onCleanup(() => {
    certificateEditorActions.reset();
  });

  return (
    <div class="h-dvh min-h-0 bg-background text-foreground">
      <div class="flex h-full flex-col items-center justify-center gap-3 px-6 text-center lg:hidden!">
        <Monitor class="size-10 text-muted-foreground" />
        <p class="max-w-xs text-sm text-muted-foreground">
          O editor de certificados foi desenvolvido para telas maiores. Abra esta página em um computador para editar o template.
        </p>
      </div>

      <div class="hidden h-full min-h-0 flex-col lg:flex">
        <header class="flex h-16 shrink-0 items-center justify-between border-b border-muted bg-card px-4 shadow-xs">
          <div class="flex items-center gap-3">
            <Link
              to="/admin/events/$eventId/editions/$editionId/certifications"
              params={{ eventId: props.eventId, editionId: props.editionId }}
              class="rounded p-1.5 hover:bg-muted"
            >
              <ArrowLeft class="size-4" />
            </Link>
            <strong class="text-sm">Editor de certificados</strong>
          </div>

          <Button
            type="button"
            variant="default"
            disabled={isSaving()}
            onClick={handleSave}
            class="cursor-pointer"
          >
            <Show
              when={isSaving()}
              fallback={
                <>
                  <Save class="size-4" />
                  <span>Salvar</span>
                </>
              }
            >
              <Loader2 class="size-4 animate-spin" />
              <span>Salvando...</span>
            </Show>
          </Button>
        </header>

        <div class="flex min-h-0 flex-1 overflow-x-auto">
          <CertificateToolsSidebar />
          <div class="flex min-h-0 flex-1 flex-col overflow-hidden">
            <CertificateTextToolbar />
            <CertificateCanvas />
          </div>
          <CertificatePropertiesPanel />
        </div>
      </div>
    </div>
  );
}
