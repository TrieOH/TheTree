import type { FormFieldI } from "@/shared/model/field";
import type { EventCreateI } from ".";
import { uploadAndModerateFile } from "@/features/storage/api";

export const getEventFields = (id: string = 'temp'): FormFieldI<EventCreateI>[] => [
  { name: 'name' as const, label: 'Nome', type: 'text' as const, required: true, placeholder: 'Ex: TechConf' },
  { name: 'slug' as const, label: 'Slug', type: 'text' as const, required: true, placeholder: 'tech-conf' },
  { name: 'acronym' as const, label: 'Sigla', type: 'text' as const, placeholder: 'TC' },
  { name: 'contact_email' as const, label: 'E-mail de contato', type: 'email' as const, required: true, placeholder: 'contato@email.com' },
  { name: 'tagline' as const, label: 'Tagline', type: 'text' as const, placeholder: 'Uma frase curta sobre o evento', span: 'full' as const },
  { name: 'description' as const, label: 'Descrição', type: 'textarea' as const, placeholder: 'Descreva o evento...', span: 'full' as const, rows: 3 },
  {
    name: 'logo_url' as const,
    label: 'Logo URL',
    type: 'image-upload' as const,
    span: 'full',
    accept: 'image/png,image/jpeg',
    maxSize: 2 * 1024 * 1024,
    uploadFn: (file) => uploadAndModerateFile(file, `events/${id}`),
  },
  {
    name: 'banner_url' as const,
    label: 'Banner URL',
    type: 'image-upload' as const,
    span: 'full',
    accept: 'image/png,image/jpeg',
    maxSize: 2 * 1024 * 1024,
    uploadFn: (file) => uploadAndModerateFile(file, `events/${id}`),
  },
  { name: 'is_series' as const, label: 'É série', type: 'checkbox' as const, placeholder: 'É uma série de eventos', span: 'full' as const },
] as const