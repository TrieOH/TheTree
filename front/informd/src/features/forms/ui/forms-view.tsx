import { formCreateSchema } from '#/features/forms/model'
import type { FormCreateI, FormI } from '#/features/forms/model';
import { FormCard } from '#/features/forms/ui/form-card'
import { Button } from '#/shared/ui/shadcn/button'
import FormModal from '#/widgets/modal/form-modal'
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid'
import { Plus } from 'lucide-react'
import { useState } from 'react'

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
  const [isCreateOpen, setIsCreateOpen] = useState(false)

  const count = forms.length

  const filteredForms = forms.filter((form) => {
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
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by title…"
        itemLabel="forms"
        headerActions={
          <Button
            onClick={() => setIsCreateOpen(true)}
            size="icon"
            variant="outline"
            className="sm:w-auto px-3 rounded-sm"
          >
            <Plus size={16} />
            <span className="hidden sm:inline ml-2">Add Form</span>
          </Button>
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
