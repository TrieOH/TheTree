import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useEffect, useMemo, useRef, useState } from 'react'
import {
  Check,
  ChevronsUpDown,
  FileText,
  Link2,
  Plus,
  Search,
} from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { EmptyState, PaginatedContainer } from '@trieoh/ui-base'
import type { CertificationTemplateI } from '@/features/certifications/model'
import { allCertificationTemplatesQueryOptions } from '@/features/certifications/api'
import {
  useSetActivityCertificationTemplateMutation,
  useSetEditionCertificationTemplateMutation,
} from '@/features/certifications/api/mutations'
import { allAdminActivitiesQueryOptions } from '@/features/activities/api'
import { allAdminEditionsQueryOptions } from '@/features/editions/api'
import { AdminCertificationTemplateCard } from '@/features/certifications/ui/AdminCertificationTemplateCard'
import { Button } from '@/shared/ui/shadcn/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/shadcn/card'
import { Label } from '@/shared/ui/shadcn/label'
import { cn } from '@/shared/lib/utils'
import type { ComboboxOption } from '@/widgets/multi-step-form/model/types'

export const Route = createFileRoute(
  '/admin/events/$eventId_/editions/$editionId/certifications/',
)({
  component: RouteComponent,
})

function SelectionCombobox({
  value,
  options,
  placeholder,
  searchPlaceholder = 'Buscar...',
  onChange,
  disabled = false,
}: {
  value: string
  options: ComboboxOption[]
  placeholder: string
  searchPlaceholder?: string
  onChange: (value: string) => void
  disabled?: boolean
}) {
  const containerRef = useRef<HTMLDivElement>(null)
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const [highlightedIndex, setHighlightedIndex] = useState(0)

  const normalizedQuery = query.trim().toLowerCase()
  const selectedOption = options.find((option) => option.value === value)
  const visibleOptions = useMemo(
    () =>
      options.filter((option) =>
        `${option.label} ${option.description ?? ''}`
          .toLowerCase()
          .includes(normalizedQuery),
      ),
    [normalizedQuery, options],
  )

  useEffect(() => {
    if (!open) return

    function closeOnOutsideClick(event: MouseEvent) {
      if (!containerRef.current?.contains(event.target as Node)) {
        setOpen(false)
      }
    }

    document.addEventListener('mousedown', closeOnOutsideClick)
    return () => document.removeEventListener('mousedown', closeOnOutsideClick)
  }, [open])

  useEffect(() => {
    setHighlightedIndex(0)
  }, [open, query])

  function select(option: ComboboxOption) {
    onChange(option.value)
    setOpen(false)
    setQuery('')
  }

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        disabled={disabled}
        aria-haspopup="listbox"
        aria-expanded={open}
        className={cn(
          'flex h-10 w-full items-center justify-between gap-2 rounded-xl border border-input bg-background px-3 text-left text-sm shadow-sm',
          'transition-colors hover:border-primary/30 hover:bg-muted/30',
          'focus:border-primary/40 focus:outline-none focus:ring-2 focus:ring-primary/10',
          disabled && 'cursor-not-allowed opacity-60',
        )}
        onClick={() => setOpen((current) => !current)}
      >
        <span
          className={cn(
            'min-w-0 truncate',
            !selectedOption && 'text-muted-foreground',
          )}
        >
          {selectedOption?.label ?? placeholder}
        </span>
        <ChevronsUpDown className="size-4 shrink-0 text-muted-foreground" />
      </button>

      {open ? (
        <div className="absolute left-0 top-full z-50 mt-2 w-full overflow-hidden rounded-xl border border-border/70 bg-popover text-popover-foreground shadow-xl">
          <div className="flex items-center gap-2 border-b border-border/60 px-3 py-2">
            <Search className="size-3.5 shrink-0 text-muted-foreground" />
            <input
              autoFocus
              value={query}
              placeholder={searchPlaceholder}
              className="h-5 w-full bg-transparent text-sm outline-none placeholder:text-muted-foreground"
              onChange={(event) => setQuery(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === 'ArrowDown') {
                  event.preventDefault()
                  setHighlightedIndex((index) =>
                    Math.min(index + 1, Math.max(visibleOptions.length - 1, 0)),
                  )
                } else if (event.key === 'ArrowUp') {
                  event.preventDefault()
                  setHighlightedIndex((index) => Math.max(index - 1, 0))
                } else if (event.key === 'Enter') {
                  event.preventDefault()
                  select(visibleOptions[highlightedIndex])
                } else if (event.key === 'Escape') {
                  setOpen(false)
                }
              }}
            />
          </div>

          <ul role="listbox" className="max-h-56 overflow-y-auto py-1">
            {visibleOptions.length === 0 ? (
              <li className="px-3 py-3 text-sm text-muted-foreground">
                Nenhum resultado encontrado
              </li>
            ) : (
              visibleOptions.map((option, index) => (
                <li key={option.value}>
                  <button
                    type="button"
                    role="option"
                    aria-selected={option.value === value}
                    className={cn(
                      'flex w-full items-center justify-between gap-2 px-3 py-2 text-left text-sm transition-colors',
                      index === highlightedIndex
                        ? 'bg-primary/10'
                        : 'hover:bg-muted/70',
                    )}
                    onMouseEnter={() => setHighlightedIndex(index)}
                    onClick={() => select(option)}
                  >
                    <span className="min-w-0">
                      <span className="block truncate">{option.label}</span>
                      {option.description ? (
                        <span className="block truncate text-xs text-muted-foreground">
                          {option.description}
                        </span>
                      ) : null}
                    </span>
                    {option.value === value ? (
                      <Check className="size-4 shrink-0" />
                    ) : null}
                  </button>
                </li>
              ))
            )}
          </ul>
        </div>
      ) : null}
    </div>
  )
}

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  const navigate = useNavigate()
  const [filter, setFilter] = useState('')
  const [selectedTemplateId, setSelectedTemplateId] = useState('')
  const [selectedActivityId, setSelectedActivityId] = useState('')

  const { data: editions = [] } = useQuery(
    allAdminEditionsQueryOptions(eventId),
  )
  const { data: templates = [] } = useQuery(
    allCertificationTemplatesQueryOptions(eventId, editionId),
  )
  const { data: activities = [] } = useQuery(
    allAdminActivitiesQueryOptions(eventId, editionId),
  )

  const edition = editions.find((item) => item.id === editionId) ?? null

  const filteredTemplates = useMemo(() => {
    const search = filter.trim().toLowerCase()
    if (!search) return templates

    return templates.filter((template) =>
      [template.title, template.url ?? ''].some((value) =>
        value.toLowerCase().includes(search),
      ),
    )
  }, [filter, templates])

  let selectedTemplate: CertificationTemplateI | null = null
  const matchedTemplate = filteredTemplates.find(
    (template) => template.id === selectedTemplateId,
  )
  if (matchedTemplate) {
    selectedTemplate = matchedTemplate
  } else if (filteredTemplates.length > 0) {
    selectedTemplate = filteredTemplates[0]
  } else if (templates.length > 0) {
    selectedTemplate = templates[0]
  }

  const editionTemplateMutation = useSetEditionCertificationTemplateMutation()
  const activityTemplateMutation = useSetActivityCertificationTemplateMutation()

  const activityOptions = activities.filter(
    (activity) => activity.status !== 'canceled',
  )
  const templateOptions = templates.map((template) => ({
    value: template.id,
    label: template.title,
    description: template.url
      ? 'Com fundo configurado'
      : 'Sem fundo configurado',
  }))
  const activitySelectOptions = activityOptions.map((activity) => ({
    value: activity.id,
    label: activity.title,
    description: activity.location,
  }))
  const isPending =
    editionTemplateMutation.isPending || activityTemplateMutation.isPending

  return (
    <div className="flex flex-wrap p-6 pb-28!">
      <PaginatedContainer<CertificationTemplateI>
        items={filteredTemplates}
        layout="grid"
        minItemWidth="16rem"
        pageSize={6}
        gap="6"
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Buscar por título ou URL..."
        itemLabel="templates"
        headerActions={
          <Link
            to="/admin/events/$eventId/editions/$editionId/certifications/editor"
            params={{ eventId, editionId }}
            className={cn(
              'inline-flex h-9 items-center justify-center gap-2 rounded-lg px-4 text-sm font-medium',
              'bg-primary text-primary-foreground shadow-sm transition-colors hover:bg-primary/90',
              'sm:min-w-40 sm:px-5',
            )}
          >
            <Plus className="size-4 shrink-0" />
            <span className="whitespace-nowrap">Novo template</span>
          </Link>
        }
        emptyState={
          <EmptyState
            icon={FileText}
            eyebrow="Certificações"
            title="Nenhum template encontrado"
            description="Crie o primeiro template para começar a emitir certificados nessa edição."
            className="border-0 bg-transparent px-0 py-4 shadow-none"
          />
        }
        renderItems={(slice) =>
          slice.map((template, index) => {
            const isSelected = Boolean(
              selectedTemplate && selectedTemplate.id === template.id,
            )

            return (
              <AdminCertificationTemplateCard
                key={template.id}
                template={template}
                selected={isSelected}
                index={index}
                onSelect={setSelectedTemplateId}
                onEdit={() => {
                  void navigate({
                    to: '/admin/events/$eventId/editions/$editionId/certifications/editor',
                    params: { eventId, editionId },
                  })
                }}
                verifyUrl={window.location.href}
                editionName={edition?.name ?? 'Nome da edição'}
              />
            )
          })
        }
      />

      <section className="mx-auto mt-6 w-full max-w-7xl px-6">
        <Card size="sm" className="border-border/60 bg-card/95 shadow-sm">
          <CardHeader className="border-b border-border/60 pb-4">
            <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <Link2 className="size-4 text-muted-foreground" />
                  <CardTitle className="text-sm font-semibold">
                    Vínculo do template
                  </CardTitle>
                </div>
                <CardDescription className="max-w-2xl text-xs">
                  Escolha um template na grade e aplique na edição inteira ou em
                  uma atividade específica.
                </CardDescription>
              </div>

              <span className="inline-flex w-fit items-center gap-1.5 rounded-full border border-border/60 bg-muted/40 px-3 py-1 text-[11px] font-medium text-muted-foreground">
                {selectedTemplate
                  ? 'Template pronto para aplicar'
                  : 'Selecione um template'}
              </span>
            </div>
          </CardHeader>

          <CardContent className="space-y-4 pt-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label className="text-xs text-muted-foreground">
                  Template
                </Label>
                <SelectionCombobox
                  value={selectedTemplate ? selectedTemplate.id : ''}
                  options={templateOptions}
                  placeholder="Selecione um template"
                  searchPlaceholder="Buscar template..."
                  onChange={setSelectedTemplateId}
                />
              </div>

              <div className="space-y-2">
                <Label className="text-xs text-muted-foreground">
                  Atividade
                </Label>
                <SelectionCombobox
                  value={selectedActivityId}
                  options={activitySelectOptions}
                  placeholder="Aplicar na edição inteira"
                  searchPlaceholder="Buscar atividade..."
                  onChange={setSelectedActivityId}
                />
              </div>
            </div>

            <div className="grid gap-3 rounded-2xl border border-border/60 bg-muted/20 p-3 md:grid-cols-2">
              <div className="min-w-0 space-y-1">
                <p className="text-[11px] font-medium uppercase tracking-[0.2em] text-muted-foreground">
                  Template ativo
                </p>
                <p className="truncate text-sm font-medium text-foreground">
                  {selectedTemplate?.title ?? 'Nenhum template selecionado'}
                </p>
              </div>
              <div className="min-w-0 space-y-1">
                <p className="text-[11px] font-medium uppercase tracking-[0.2em] text-muted-foreground">
                  Destino
                </p>
                <p className="truncate text-sm font-medium text-foreground">
                  {selectedActivityId
                    ? (activityOptions.find(
                        (activity) => activity.id === selectedActivityId,
                      )?.title ?? 'Atividade selecionada')
                    : 'Aplicação na edição inteira'}
                </p>
              </div>
            </div>

            <div className="flex flex-col gap-2 md:flex-row md:justify-end">
              <Button
                type="button"
                className="h-9 gap-2 md:min-w-44"
                disabled={selectedTemplate === null || isPending}
                onClick={() => {
                  if (selectedTemplate === null) return
                  editionTemplateMutation.mutate({
                    eventId,
                    editionId,
                    data: {
                      certification_template_id: selectedTemplate.id,
                    },
                  })
                }}
              >
                Definir para a edição
              </Button>
              <Button
                type="button"
                variant="outline"
                className="h-9 gap-2 md:min-w-44"
                disabled={
                  selectedTemplate === null || !selectedActivityId || isPending
                }
                onClick={() => {
                  if (selectedTemplate === null || !selectedActivityId) return
                  activityTemplateMutation.mutate({
                    eventId,
                    editionId,
                    activityId: selectedActivityId,
                    data: {
                      certification_template_id: selectedTemplate.id,
                    },
                  })
                }}
              >
                Definir para a atividade
              </Button>
            </div>
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
