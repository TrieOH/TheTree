import { useMemo } from "react";
import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import {
  type ActivityCreateInputI,
  type ActivityCreateOutputI,
  type ActivityI,
  activityCreateSchema,
} from "../model";
import { createActivityFormSteps } from "../model/activity-form-steps";

export interface ManageActivityModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  activity?: ActivityI;
  onCreate?: (
    values: ActivityCreateOutputI,
  ) => Promise<ActivityI | null | boolean> | ActivityI | null | boolean;
  onUpdate?: (
    id: string,
    values: ActivityCreateOutputI,
  ) => Promise<ActivityI | null | boolean> | ActivityI | null | boolean;
}

const emptyDefaultValues: ActivityCreateInputI = {
  title: "",
  description: "",
  location: "",
  starts_at: "",
  ends_at: "",
  presenter_name: "",
  token_cost: 0,
  has_capacity: false,
  capacity: 0,
  difficulty: "no_prerequisites",
};

function toFormValues(activity: ActivityI): ActivityCreateInputI {
  return {
    title: activity.title,
    description: activity.description ?? "",
    location: activity.location,
    starts_at: activity.starts_at,
    ends_at: activity.ends_at,
    presenter_name: activity.presenter_name ?? "",
    token_cost: activity.token_cost,
    has_capacity: activity.has_capacity,
    capacity: activity.capacity,
    difficulty: activity.difficulty,
  };
}

export function ManageActivityModal({
  open,
  onOpenChange,
  activity,
  onCreate,
  onUpdate,
}: ManageActivityModalProps) {
  const isEditing = Boolean(activity);
  const editValues = useMemo(
    () => (activity ? toFormValues(activity) : undefined),
    [activity],
  );
  const steps = useMemo(() => createActivityFormSteps(), []);

  const controller = useMultiStepForm({
    schema: activityCreateSchema,
    steps,
    defaultValues: emptyDefaultValues,
    values: editValues,
    requireDirtyToSubmit: isEditing,
    onSubmit: async (values): Promise<boolean> => {
      const result = activity
        ? await onUpdate?.(activity.id, values)
        : await onCreate?.(values);

      return Boolean(result);
    },
    onSubmitSuccess: () => onOpenChange(false),
  });

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEditing ? "Editar atividade" : "Criar atividade"}
      controller={controller}
      submitLabel={isEditing ? "Salvar alterações" : "Criar atividade"}
    />
  );
}
