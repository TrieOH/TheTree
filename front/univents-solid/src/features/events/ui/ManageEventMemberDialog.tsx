import type { JSX } from "@solidjs/web";
import { For, createEffect, createSignal } from "solid-js";
import {
  MultiStepDialog,
  type MultiStepField,
  type MultiStepItem,
  cn,
} from "@trieoh/ui-solid";
import ShieldCheckIcon from "~icons/lucide/shield-check";
import UserCogIcon from "~icons/lucide/user-cog";
import UsersIcon from "~icons/lucide/users";
import { toast } from "@/shared/ui/toast";
import type { EventMemberRole } from "../model/member";

const ShieldCheck = ShieldCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const UserCog = UserCogIcon as unknown as (props: { class?: string }) => JSX.Element;
const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface ManageEventMemberValues {
  email: string;
  role: EventMemberRole;
}

export interface ManageEventMemberDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (values: ManageEventMemberValues) => Promise<boolean>;
}

const ROLE_OPTIONS: Array<{
  value: EventMemberRole;
  label: string;
  description: string;
  icon: (props: { class?: string }) => JSX.Element;
}> = [
  {
    value: "staff",
    label: "Equipe",
    description: "Acesso operacional ao evento (check-in, controle de entrada e apoio geral).",
    icon: Users,
  },
  {
    value: "admin",
    label: "Administrador",
    description: "Pode gerenciar edições, ingressos, certificados, produtos e membros.",
    icon: UserCog,
  },
  {
    value: "owner",
    label: "Proprietário",
    description: "Controle integral do evento, incluindo exclusão e transferências.",
    icon: ShieldCheck,
  },
];

const ROLE_LABELS: Record<EventMemberRole, string> = {
  staff: "Equipe",
  admin: "Administrador",
  owner: "Proprietário",
};

const emptyValues: ManageEventMemberValues = {
  email: "",
  role: "staff",
};

const MEMBER_EMAIL_FIELDS: MultiStepField<ManageEventMemberValues>[] = [
  {
    kind: "email",
    name: "email",
    label: "E-mail do colaborador",
    placeholder: "membro@exemplo.com",
    hint: "Será enviado um convite com acesso ao evento para este e-mail.",
    autocomplete: "email",
    required: true,
  },
];

export function ManageEventMemberDialog(props: ManageEventMemberDialogProps): JSX.Element {
  let rawValues: ManageEventMemberValues = { ...emptyValues };
  const [values, setValues] = createSignal<ManageEventMemberValues>({ ...emptyValues });
  const [submitting, setSubmitting] = createSignal(false);
  const [errors, setErrors] = createSignal<Partial<Record<keyof ManageEventMemberValues, string>>>({});

  // Reset values when dialog opens
  createEffect(
    () => props.open,
    (isOpen) => {
      if (isOpen) {
        rawValues = { ...emptyValues };
        setValues({ ...emptyValues });
        setErrors({});
      }
    },
  );

  const update = <K extends keyof ManageEventMemberValues>(key: K, value: ManageEventMemberValues[K]) => {
    rawValues = { ...rawValues, [key]: value };
    setValues({ ...rawValues });
    setErrors((prev) => {
      const updated = { ...prev };
      delete updated[key];
      if (rawValues.email.trim().includes("@")) {
        delete updated.email;
      }
      return updated;
    });
  };

  const validateStep = (stepIndex: number): boolean => {
    const v = rawValues;
    const stepErrors: Partial<Record<keyof ManageEventMemberValues, string>> = {};

    if (stepIndex === 0) {
      const cleanEmail = v.email.trim().toLowerCase();
      if (!cleanEmail || !cleanEmail.includes("@")) {
        stepErrors.email = "Informe um e-mail válido.";
      }
    } else if (stepIndex === 1) {
      if (!v.role) {
        stepErrors.role = "Selecione uma função para o membro.";
      }
    }

    setErrors((prev) => ({ ...prev, ...stepErrors }));
    return Object.keys(stepErrors).length === 0;
  };

  const submitCurrent = async (): Promise<boolean> => {
    if (!validateStep(0)) {
      return false;
    }
    if (!validateStep(1)) {
      return false;
    }

    const current: ManageEventMemberValues = {
      email: rawValues.email.trim().toLowerCase(),
      role: rawValues.role,
    };

    setSubmitting(true);
    try {
      const ok = await props.onSubmit(current);
      if (ok) {
        props.onOpenChange(false);
        return true;
      } else {
        return false;
      }
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "Erro ao adicionar membro.";
      toast.error(message);
      return false;
    } finally {
      setSubmitting(false);
    }
  };

  const selectedRoleOption = () =>
    ROLE_OPTIONS.find((opt) => opt.value === values().role) ?? ROLE_OPTIONS[0];

  const steps: MultiStepItem<ManageEventMemberValues>[] = [
    {
      id: "identificacao",
      title: "Identificação",
      description: "E-mail do colaborador",
      fields: MEMBER_EMAIL_FIELDS,
    },
    {
      id: "permissao",
      title: "Permissão",
      description: "Função no evento",
      fields: [
        {
          kind: "custom",
          name: "role",
          label: "Selecione o papel do colaborador",
          required: true,
          render: () => (
            <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
              <For each={ROLE_OPTIONS}>
                {(option) => {
                  const isSelected = () => values().role === option.value;
                  return (
                    <button
                      type="button"
                      onClick={() => update("role", option.value)}
                      class={cn(
                        "group flex flex-col items-start gap-2 rounded-xl border p-3.5 text-left transition-all",
                        isSelected()
                          ? "border-primary bg-primary/10 text-primary ring-2 ring-primary/20 shadow-xs"
                          : "border-border bg-card text-foreground hover:border-primary/40 hover:bg-muted/40",
                      )}
                    >
                      <div class="flex items-center gap-2 font-medium text-sm">
                        <span
                          class={cn(
                            "flex size-7 shrink-0 items-center justify-center rounded-lg transition-colors",
                            isSelected()
                              ? "bg-primary text-primary-foreground"
                              : "bg-muted text-muted-foreground group-hover:text-foreground",
                          )}
                        >
                          {option.icon({ class: "size-4" })}
                        </span>
                        <span class="truncate font-semibold">{option.label}</span>
                      </div>
                      <span
                        class={cn(
                          "text-xs leading-relaxed line-clamp-3",
                          isSelected() ? "text-primary/90" : "text-muted-foreground",
                        )}
                      >
                        {option.description}
                      </span>
                    </button>
                  );
                }}
              </For>
            </div>
          ),
        },
      ],
    },
    {
      id: "resumo",
      title: "Resumo",
      description: "Confirmação do convite",
      summary: {
        badge: "Pronto para adicionar",
        items: [
          {
            label: "E-mail do colaborador",
            value: () => values().email.trim().toLowerCase(),
            mono: true,
          },
          {
            label: "Função atribuída",
            value: () => ROLE_LABELS[values().role],
          },
          {
            label: "Permissões concedidas",
            value: () => selectedRoleOption().description,
            fullWidth: true,
          },
        ],
      },
    },
  ];

  return (
    <MultiStepDialog<ManageEventMemberValues>
      open={props.open}
      onOpenChange={props.onOpenChange}
      title="Adicionar membro"
      description="Convide uma pessoa para colaborar na gestão deste evento em etapas organizadas."
      steps={steps}
      values={values()}
      errors={errors()}
      onChange={(key, val) => update(key, val as ManageEventMemberValues[typeof key])}
      formId="manage-member-form"
      loading={submitting()}
      submitLabel="Adicionar membro"
      nextLabel="Continuar"
      onBeforeNext={(step) => validateStep(step)}
      onFormSubmit={submitCurrent}
      onSubmit={submitCurrent}
    />
  );
}
