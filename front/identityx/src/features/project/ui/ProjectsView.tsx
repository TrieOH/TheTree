import { Plus } from 'lucide-react'
import { useState } from 'react'
import { projectCreateSchema, type ProjectCreateI, type ProjectI } from '../model';
import { PaginatedContainer } from '@trieoh/ui-base';
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { FormModal } from '@/widgets/modal/FormModal';
import ProjectCard from './project-card';

interface ProjectsViewProps {
  projects: ProjectI[];
  onCreate: (data: ProjectCreateI) => void;
  isCreating: boolean;
  title: string;
  description: string;
}

export function ProjectsView({
  projects,
  onCreate,
  isCreating,
  title,
  description,
}: ProjectsViewProps) {
  const [filter, setFilter] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)

  const count = projects.length

  const filteredProjects = projects.filter((project) => {
    const search = filter.toLowerCase().trim()

    if (!search) return true

    return (
      project.name.toLowerCase().includes(search) ||
      project.domain?.toLowerCase().includes(search)
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
                ? `No projects yet ${description}`
                : `${count} project${count !== 1 ? 's' : ''} ${description}`}
            </p>
          )}
        </div>
      )}

      <PaginatedContainer<ProjectI>
        items={filteredProjects}
        layout='grid'
        minItemWidth='14rem'
        gap='6'
        pageSize={10}
        sortFields={[
          { key: "name", label: "Name" },
          { key: "domain", label: "Domain" },
          { key: "created_at", label: "Created At" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by name..."
        itemLabel="projects"
        headerActions={
          <ShadowButton
            onClick={() => setIsCreateOpen(true)}
            variant="outline"
            className="h-9 sm:w-auto px-3 rounded-sm"
            leftIcon={<Plus size={16} />}
            value='Add Project'
          />
        }
        renderItems={(slice) => slice.map(item => {
          return (
            <ProjectCard key={item.id} data={item} />
          )
        })}
      />
      <FormModal<ProjectCreateI>
        title="Create Project"
        description="Give your project a name to identify it."
        schema={projectCreateSchema}
        formId="create-project-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={(data) => {
          onCreate(data);
          setIsCreateOpen(false);
        }}
        defaultValues={{ name: '', domain: '' }}
        isLoading={isCreating}
        submitLabel='Create Project'
        fields={[
          {
            name: 'name',
            label: 'Project Name',
            placeholder: 'e.g. My Own Project',
            required: true,
          },
          {
            name: 'domain',
            label: 'Project Domain',
            placeholder: 'e.g. my.domain.com',
            required: true,
          },
        ]}
      />
    </div>
  )
}
