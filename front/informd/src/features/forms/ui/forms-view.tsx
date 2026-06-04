import { formCreateSchema } from '#/features/forms/model'
import type { FormCreateI, FormI } from '#/features/forms/model';
import { FormCard } from '#/features/forms/ui/form-card'
import { Button } from '#/shared/ui/shadcn/button'
import FormModal from '#/widgets/modal/form-modal'
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid'
import { FormStatusArchived } from '@trieoh/informd-models'
import { Plus } from 'lucide-react'
import { useMemo, useState } from 'react'

interface FormsViewProps {
  forms: FormI[];
  onCreate: (data: FormCreateI) => void;
  isCreating: boolean;
  title: string;
  description: string;
}

export function FormsView({
  forms,
  onCreate,
  isCreating,
  title,
  description,
}: FormsViewProps) {
  const [filter, setFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'archived'>('active')
  const [isCreateOpen, setIsCreateOpen] = useState(false)

  const visibleForms = useMemo(() => {
    return forms.filter((form) => {
      if (statusFilter === 'archived') return form.status === FormStatusArchived
      if (statusFilter === 'active') return form.status !== FormStatusArchived
      return true
    })
  }, [forms, statusFilter])

  const count = visibleForms.length

  const filteredForms = visibleForms.filter((form) => {
    const search = filter.toLowerCase().trim()

    if (!search) return true

    return (
      form.title.toLowerCase().includes(search) ||
      form.created_at.toLowerCase().includes(search) ||
      form.updated_at.toLowerCase().includes(search) ||
      form.status.toLowerCase().includes(search)
    )
  })

  return (
    <div>
      {(title || description) && (
        <div className="mb-6">
          {title && <h1 className="text-lg font-semibold tracking-tight">{title}</h1>}
          {description && (
            <p className="text-sm text-muted-foreground">
              {count === 0
                ? `No forms yet ${description}`
                : `${count} form${count !== 1 ? 's' : ''} ${description}`}
            </p>
          )}
        </div>
      )}

      <PaginatedContainer<FormI>
        items={filteredForms}
        layout='grid'
        minItemWidth='14rem'
        gap='6'
        pageSize={10}
        sortFields={[
          { key: "title", label: "Title" },
          { key: "created_at", label: "Created At" },
          { key: "updated_at", label: "Updated At" },
          {
            key: "namespace_id",
            label: "Scope",
            comparator: (a, b) => Number(Boolean(b.namespace_id)) - Number(Boolean(a.namespace_id)),
          },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by title…"
        itemLabel="forms"
        headerActions={
          <div className="flex flex-wrap items-center gap-2">
            <div className="flex h-9 items-center gap-1 rounded-md border border-border bg-background p-1">
              {[
                { key: 'active', label: 'Active' },
                { key: 'archived', label: 'Archived' },
                { key: 'all', label: 'All' },
              ].map((option) => (
                <button
                  key={option.key}
                  type="button"
                  onClick={() => setStatusFilter(option.key as 'all' | 'active' | 'archived')}
                  className={[
                    'flex h-full items-center rounded-sm px-2.5 text-xs font-medium transition-colors',
                    statusFilter === option.key
                      ? 'bg-primary text-primary-foreground'
                      : 'text-muted-foreground hover:bg-muted hover:text-foreground',
                  ].join(' ')}
                >
                  {option.label}
                </button>
              ))}
            </div>

            <Button
              onClick={() => setIsCreateOpen(true)}
              size="icon"
              variant="outline"
              className="h-9 sm:w-auto px-3 rounded-sm"
            >
              <Plus size={16} />
              <span className="hidden sm:inline ml-2">Add Form</span>
            </Button>
          </div>
        }
        renderItems={(slice) => slice.map(item => {
          return (
            <FormCard key={item.id} data={item} />
          )
        })}
      />
      <FormModal<FormCreateI>
        title="Create Form"
        description="Give your form a title to identify it."
        buttonTitle="Create Form"
        schema={formCreateSchema}
        formId="create-form-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={(data) => {
          onCreate(data);
          setIsCreateOpen(false);
        }}
        fields={[
          {
            name: 'title',
            label: 'Form Title',
            type: 'text',
            placeholder: 'e.g. Contact Form',
          },
        ]}
        disabled={isCreating}
      />
    </div>
  )
}
