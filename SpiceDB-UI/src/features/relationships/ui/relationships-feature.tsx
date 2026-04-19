import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { useState, useMemo, useEffect } from 'react'
import { schemaQueryOptions } from '#/features/schema/api'
import {
  createRelationship,
  deleteRelationship,
  relationshipQueryOptions,
  updateRelationship,
} from '../api'
import { parseSpiceDBSchema } from '../lib/schema-parser'
import type { RelationshipFormState } from '../model'
import { RelationshipForm } from './relationship-form'
import { RelationshipsTable } from './relationships-table'
import type { Relationship } from '@soramux/node-perm-sdk'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '#/shared/ui/shadcn/alert-dialog'
import { Database } from 'lucide-react'

const INITIAL_FORM: RelationshipFormState = {
  resource: '',
  resourceId: '',
  relation: '',
  subject: '',
  subjectId: '',
}

export function RelationshipsFeature({ envId }: { envId: string }) {
  const queryClient = useQueryClient()
  const [isEditing, setIsEditing] = useState(false)
  const [formData, setFormData] = useState<RelationshipFormState>(INITIAL_FORM)
  const [mobileView, setMobileView] = useState<'table' | 'form'>('table')
  const [relToDelete, setRelToDelete] = useState<Relationship.SpiceDBRelationshipI | null>(null)

  // Fetch schema to get dynamic definitions and relations
  const { data: schemaData } = useQuery(schemaQueryOptions(envId))

  const parsedSchema = useMemo(() => {
    if (!schemaData?.schemaText) return { definitions: [], relationsByDefinition: {} }
    const parsed = parseSpiceDBSchema(schemaData.schemaText)
    return {
      ...parsed,
      definitions: [...parsed.definitions].sort(),
    }
  }, [schemaData?.schemaText])

  useEffect(() => {
    if (parsedSchema.definitions.length > 0 && !isEditing && formData.resource === '') {
      const firstDefinition = parsedSchema.definitions[0]
      const firstRelation = parsedSchema.relationsByDefinition[firstDefinition]?.find(r => r.type === 'relation')?.name ?? ''
      setFormData({
        ...INITIAL_FORM,
        resource: firstDefinition,
        subject: firstDefinition,
        relation: firstRelation,
      })
    }
  }, [parsedSchema.definitions, parsedSchema.relationsByDefinition, formData.resource, isEditing])

  // Fetch relationships for all resource types
  const { data: relationships = [], isLoading: isQueryLoading } = useQuery({
    ...relationshipQueryOptions(envId, parsedSchema.definitions),
    enabled: parsedSchema.definitions.length > 0,
  })

  const createMutation = useMutation({
    mutationFn: (data: RelationshipFormState) =>
      createRelationship({ data: { ...data, envId } }),
    onSuccess: (res) => {
      if (res.success) {
        toast.success('Relationship created')
        queryClient.invalidateQueries({ queryKey: ['relationship', envId] })
        handleCancel()
      } else {
        toast.error(res.message || 'Error creating relationship')
      }
    },
  })

  const updateMutation = useMutation({
    mutationFn: (data: RelationshipFormState) =>
      updateRelationship({ data: { ...data, envId } }),
    onSuccess: (res) => {
      if (res.success) {
        toast.success('Relationship updated')
        queryClient.invalidateQueries({ queryKey: ['relationship', envId] })
        handleCancel()
      } else {
        toast.error(res.message || 'Error updating relationship')
      }
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (rel: Relationship.SpiceDBRelationshipI) =>
      deleteRelationship({
        data: {
          envId,
          resource: rel.resource.objectType,
          resourceId: rel.resource.objectId,
          relation: rel.relation,
          subject: rel.subject.object.objectType,
          subjectId: rel.subject.object.objectId,
        },
      }),
    onSuccess: (res) => {
      if (res.success) {
        toast.success('Relationship deleted')
        queryClient.invalidateQueries({ queryKey: ['relationship', envId] })
      } else {
        toast.error(res.message || 'Error deleting relationship')
      }
      setRelToDelete(null)
    },
  })

  // When opening "New", ensure the first available definition is selected if form is empty
  function handleNew() {
    setIsEditing(false)
    const initialForm = { ...INITIAL_FORM }
    if (parsedSchema.definitions.length > 0) {
      initialForm.resource = parsedSchema.definitions[0]
      initialForm.subject = parsedSchema.definitions[0]
    }
    setFormData(initialForm)
    setMobileView('form')
  }

  function handleEdit(rel: Relationship.SpiceDBRelationshipI) {
    setIsEditing(true)
    setFormData({
      resource: rel.resource.objectType,
      resourceId: rel.resource.objectId,
      relation: rel.relation,
      subject: rel.subject.object.objectType,
      subjectId: rel.subject.object.objectId,
    })
    setMobileView('form')
  }

  function handleCancel() {
    setIsEditing(false)
    setFormData(INITIAL_FORM)
    setMobileView('table')
  }

  function handleSubmit(data: RelationshipFormState) {
    if (isEditing) {
      updateMutation.mutate(data)
    } else {
      createMutation.mutate(data)
    }
  }

  function handleDelete(rel: Relationship.SpiceDBRelationshipI) {
    setRelToDelete(rel)
  }

  const isLoading = createMutation.isPending || updateMutation.isPending || isQueryLoading

  const mobilePreview = `${formData.resource || 'resource'}:${formData.resourceId || 'id'}#${formData.relation || 'relation'}@${formData.subject || 'subject'}:${formData.subjectId || 'id'}`

  return (
    <main className="h-(--content-height) flex flex-col bg-background border-l">
      <div className="flex h-14 items-center border-b px-4 bg-background shrink-0 min-w-0 gap-3">
        <div className="flex items-center gap-2 shrink-0">
          <Database size={18} className="text-primary" />
          <span className="text-sm font-bold whitespace-nowrap hidden sm:inline">Relationships:</span>
        </div>
        <div className="flex-1 min-w-0 font-mono text-xs sm:text-sm overflow-hidden">
          <div className="truncate text-muted-foreground sm:text-foreground">
            {mobilePreview}
          </div>
        </div>
      </div>

      <div className="flex flex-1 overflow-hidden relative">
        {/* Form sidebar */}
        <div
          className={`flex flex-col border-r w-full md:w-85 md:shrink-0 bg-background transition-transform duration-300 ease-out absolute md:relative inset-0 md:translate-x-0 z-10 ${mobileView === 'table' ? '-translate-x-full md:translate-x-0' : 'translate-x-0'
            }`}
        >
          <RelationshipForm
            key={isEditing ? 'edit' : 'new'}
            initial={formData}
            isEditing={isEditing}
            isLoading={isLoading}
            definitions={parsedSchema.definitions}
            relationsByDefinition={parsedSchema.relationsByDefinition}
            onSubmit={handleSubmit}
            onCancel={handleCancel}
          />
        </div>

        {/* Table pane */}
        <div
          className={`flex-1 overflow-y-auto w-full transition-transform duration-300 ease-out absolute md:relative inset-0 md:translate-x-0 ${mobileView === 'form' ? 'translate-x-full md:translate-x-0' : 'translate-x-0'
            }`}
        >
          <RelationshipsTable
            rows={relationships}
            onEdit={handleEdit}
            onDelete={handleDelete}
            onNew={handleNew}
          />
        </div>
      </div>


      <AlertDialog open={!!relToDelete} onOpenChange={(open) => !open && setRelToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the relationship.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              onClick={() => relToDelete && deleteMutation.mutate(relToDelete)}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? 'Deleting...' : 'Delete'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </main>
  )
}
