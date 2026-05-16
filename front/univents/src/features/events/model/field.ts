import type { FormFieldI } from "@/shared/model/field";
import type { EventCreateI } from ".";
import { uploadAndModerateFile } from "@/features/storage/api";

export const getEventFields = (id: string = 'temp'): FormFieldI<EventCreateI>[] => [
  { name: 'name' as const, label: 'Nome', type: 'text' as const, required: true, placeholder: 'Ex: TechConf', autocomplete: 'organization', autoFocus: true },
  { name: 'slug' as const, label: 'Slug', type: 'text' as const, required: true, placeholder: 'tech-conf' },
  { name: 'acronym' as const, label: 'Sigla', type: 'text' as const, placeholder: 'TC' },
  { name: 'contact_email' as const, label: 'E-mail de contato', type: 'email' as const, required: true, placeholder: 'contato@email.com', autocomplete: 'email' },
  { name: 'tagline' as const, label: 'Tagline', type: 'text' as const, placeholder: 'Uma frase curta sobre o evento', span: 'full' as const },
  { name: 'description' as const, label: 'Descrição', type: 'textarea' as const, placeholder: 'Descreva o evento...', span: 'full' as const, rows: 3 },
  { name: 'social_links.twitter' as const, label: 'Twitter (X)', type: 'url' as const, placeholder: 'https://twitter.com/...' },
  { name: 'social_links.instagram' as const, label: 'Instagram', type: 'url' as const, placeholder: 'https://instagram.com/...' },
  { name: 'social_links.linkedin' as const, label: 'LinkedIn', type: 'url' as const, placeholder: 'https://linkedin.com/...' },
  { name: 'social_links.website' as const, label: 'Site oficial', type: 'url' as const, placeholder: 'https://www.meusite.com' },
  {
    name: 'logo_url' as const,
    label: 'Logo',
    type: 'image-upload' as const,
    span: 'full',
    accept: 'image/png,image/jpeg,image/webp',
    maxSize: 2 * 1024 * 1024,
    uploadFn: (file) => uploadAndModerateFile(file, `events/${id}`),
  },
  {
    name: 'banner_url' as const,
    label: 'Banner',
    type: 'image-upload' as const,
    span: 'full',
    accept: 'image/png,image/jpeg,image/webp',
    maxSize: 5 * 1024 * 1024,
    uploadFn: (file) => uploadAndModerateFile(file, `events/${id}`),
  },
  { name: 'is_series' as const, label: 'É série', type: 'checkbox' as const, placeholder: 'É uma série de eventos', span: 'full' as const },
  {
    name: 'gallery_urls' as const,
    label: 'Fotos da Galeria',
    type: 'gallery-upload' as const,
    span: 'full',
    accept: 'image/png,image/jpeg,image/webp',
    maxSize: 5 * 1024 * 1024,
    uploadFn: (file) => uploadAndModerateFile(file, `events/${id}`),
    itemActions: [
      {
        label: 'Definir como Logo',
        icon: 'star' as const,
        onClick: (url, setValue) => { setValue('logo_url', url); },
      },
      {
        label: 'Definir como Banner',
        icon: 'layout' as const,
        onClick: (url, setValue) => { setValue('banner_url', url); },
      },
    ]
  },
]
