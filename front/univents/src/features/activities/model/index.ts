import z from "zod";

const difficultyTypeSchema = z
  .enum(
    ["no_prerequisites", "beginner", "intermediate", "advanced", "expert"],
    { error: "Invalid difficulty type" }
  ).default("no_prerequisites")
type DifficultyType = z.infer<typeof difficultyTypeSchema>

const requiredText = (label: string) =>
  z.string({ error: `${label} é obrigatório` }).trim().min(1, { message: `${label} é obrigatório` })

const requiredDateTime = (label: string) =>
  z.string({ error: `${label} é obrigatório` })
    .trim()
    .min(1, { message: `${label} é obrigatório` })
    .pipe(z.iso.datetime({ message: `${label} inválida` }))

const optionalText = () =>
  z.string().optional().nullable().or(z.literal('')).transform((value) => (value === '' ? null : value))

export const activityCreateSchema = z.object({
  title: requiredText('Título').min(3, { message: 'Título precisa ter pelo menos 3 caracteres' }),
  description: optionalText(),
  location: requiredText('Local'),
  starts_at: requiredDateTime('Data de início'),
  ends_at: requiredDateTime('Data de término'),
  presenter_name: optionalText(),
  token_cost: z.coerce.number({ error: 'Custo em tokens é obrigatório' }).int({ message: 'Custo em tokens precisa ser um número inteiro' }).nonnegative({ message: 'Custo em tokens não pode ser negativo' }),
  has_capacity: z.boolean().default(false),
  capacity: z.coerce.number({ error: 'Capacidade é obrigatória' }).int({ message: 'Capacidade precisa ser um número inteiro' }).nonnegative({ message: 'Capacidade não pode ser negativa' }).default(0),
  difficulty: difficultyTypeSchema,
})

export type ActivityCreateInputI = z.input<typeof activityCreateSchema>
export type ActivityCreateOutputI = z.output<typeof activityCreateSchema>


export interface ActivityI {
  id: string;
  scope_id: string;
  edition_id: string;
  title: string;
  description: string | null;
  status: "draft" | "published" | "ongoing" | "completed" | "canceled";
  location: string;
  starts_at: string;
  ends_at: string;
  presenter_name: string | null;
  token_cost: number;
  has_capacity: boolean;
  capacity: number;
  remaining_capacity: number;
  difficulty: DifficultyType;
  certification_template_id?: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}

export type AttendanceStatusType = "registered" | "waitlisted" | "promoted" |
  "checked_in" | "checked_out" | "completed" | "partial" | "no_show" | "cancelled";

export interface AttendanceRecordI {
  id: string;
  activity_id: string;
  user_id: string;
  status: AttendanceStatusType;
  checked_in_at: string | null;
  cancelled_at: string | null;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}
