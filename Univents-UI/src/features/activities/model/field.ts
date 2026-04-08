import type { FormFieldI } from "@/shared/model/field"
import type { ActivityCreateI } from "."

export const getActivityFields = (): FormFieldI<ActivityCreateI>[] => [
  { name: 'title' as const, label: 'Título', type: 'text' as const, placeholder: 'Nome da atividade', required: true, span: 'full' as const },
  { name: 'description' as const, label: 'Descrição', type: 'textarea' as const, placeholder: 'Descreva a atividade...', span: 'full' as const, rows: 4 },
  { name: 'presenter_name' as const, label: 'Apresentador/Palestrante', type: 'text' as const, placeholder: "Nome do responsável" },
  { name: 'location' as const, label: 'Local', type: 'text' as const, placeholder: 'Sala, auditório, etc.', required: true },
  { name: 'starts_at' as const, label: 'Início', type: 'datetime' as const, required: true },
  { name: 'ends_at' as const, label: 'Término', type: 'datetime' as const, required: true },
  { name: 'token_cost' as const, label: 'Custo em Tokens', type: 'number' as const, placeholder: '0 (gratuito)', },
  {
    name: 'difficulty' as const,
    label: 'Nível de dificuldade',
    type: 'select' as const,
    options: [
      { label: "Sem pré-requisitos", value: "no_prerequisites" },
      { label: "Iniciante", value: "beginner" },
      { label: "Intermediário", value: "intermediate" },
      { label: "Avançado", value: "advanced" },
      { label: "Especialista", value: "expert" },
    ],
    required: true
  },
  { name: 'has_capacity' as const, label: 'Limitar vagas', type: 'checkbox' as const, placeholder: 'Definir capacidade máxima de participantes', span: 'full' as const },
  { name: 'capacity' as const, label: 'Capacidade', type: 'number' as const, placeholder: 'Número de vagas', span: 'full' as const },
]