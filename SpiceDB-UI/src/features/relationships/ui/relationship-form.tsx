import {
  ArrowRight,
  FileText,
  Pencil,
  Plus,
  User,
} from 'lucide-react'
import { useState, useEffect } from 'react'
import type { RelationshipFormState } from '../model'
import { FieldLabel } from './field-label'
import type { RelationInfo } from '../lib/schema-parser'
import CustomSelect from '#/shared/ui/custom-select'

const inputCls =
  'h-9 w-full rounded-md border border-input bg-background px-3 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring transition-shadow'

interface RelationshipFormProps {
  initial: RelationshipFormState
  isEditing: boolean
  isLoading: boolean
  definitions: string[]
  relationsByDefinition: Record<string, RelationInfo[] | undefined>
  onSubmit: (data: RelationshipFormState) => void
  onCancel: () => void
}

export function RelationshipForm({
  initial,
  isEditing,
  isLoading,
  definitions,
  relationsByDefinition,
  onSubmit,
  onCancel,
}: RelationshipFormProps) {
  const [form, setForm] = useState<RelationshipFormState>(initial)

  useEffect(() => {
    setForm(initial)
  }, [initial])

  function set(key: keyof RelationshipFormState, value: string) {
    setForm((prev) => ({ ...prev, [key]: value }))
  }

  function handleSubmit() {
    if (!form.resourceId.trim() || !form.relation || !form.subjectId.trim())
      return
    onSubmit(form)
  }

  const availableRelations = (relationsByDefinition[form.resource] || []).filter(
    (r) => r.type === 'relation'
  )

  const selectedRelation = availableRelations.find((r) => r.name === form.relation)
  const allowedSubjects = selectedRelation?.allowedSubjectTypes || []

  // Update subject type if it's no longer allowed
  useEffect(() => {
    if (form.relation && allowedSubjects.length > 0) {
      if (!allowedSubjects.includes(form.subject)) {
        set('subject', allowedSubjects[0])
      }
    }
  }, [form.relation, allowedSubjects, form.subject])

  return (
    <div className="flex h-full flex-col">
      {/* Header */}
      <div className="flex items-center justify-between border-b px-5 py-3.5">
        <div className="flex items-center gap-2">
          <ArrowRight size={14} className="text-muted-foreground" />
          <span className="text-sm font-medium">
            {isEditing ? 'Edit relationship' : 'New relationship'}
          </span>
          {isEditing && (
            <span className="rounded-full bg-accent/20 px-2 py-0.5 text-xs font-medium text-accent-foreground">
              editing
            </span>
          )}
        </div>
        {isEditing && (
          <button
            onClick={onCancel}
            className="text-xs text-muted-foreground transition-colors hover:text-foreground"
          >
            Cancel
          </button>
        )}
      </div>

      {/* Body */}
      <div className="flex flex-1 flex-col gap-5 overflow-y-auto p-5">

        {/* Resource type */}
        <div>
          <FieldLabel icon={<FileText size={11} />}>Resource type</FieldLabel>
          <CustomSelect
            value={form.resource}
            onChange={(val) => {
              setForm(prev => ({
                ...prev,
                resource: val ?? "",
                relation: ''
              }))
            }}
            placeholder="Select a resource..."
            options={definitions}
            disabled={isEditing}
          />
        </div>

        {/* Resource ID */}
        <div>
          <FieldLabel icon={<FileText size={11} />}>Resource ID</FieldLabel>
          <input
            type="text"
            value={form.resourceId}
            onChange={(e) => set('resourceId', e.target.value)}
            placeholder="e.g. roadmap"
            className={inputCls}
            disabled={isEditing}
          />
        </div>

        {/* Relation */}
        <div>
          <FieldLabel icon={<ArrowRight size={11} />}>Relation</FieldLabel>
          <CustomSelect
            value={form.relation}
            onChange={(val) => set('relation', val ?? "")}
            placeholder="Select a relation..."
            options={availableRelations.map(item => item.name)}
            disabled={isEditing}
          />
        </div>

        {/* Subject type */}
        <div>
          <FieldLabel icon={<User size={11} />}>Subject type</FieldLabel>
          <CustomSelect
            value={form.subject}
            onChange={(val) => set('subject', val ?? "")}
            placeholder="Select a subject..."
            options={allowedSubjects.length > 0 ? allowedSubjects : definitions}
            disabled={isEditing}
          />
        </div>

        {/* Subject ID */}
        <div>
          <FieldLabel icon={<User size={11} />}>Subject ID</FieldLabel>
          <input
            type="text"
            value={form.subjectId}
            onChange={(e) => set('subjectId', e.target.value)}
            placeholder="e.g. alice"
            className={inputCls}
            disabled={isEditing}
          />
        </div>
      </div>


      {/* Footer */}
      <div className="border-t p-4">
        <button
          onClick={handleSubmit}
          disabled={isLoading}
          className="inline-flex h-9 w-full items-center justify-center gap-2 rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50"
        >
          {isLoading ? (
            'Processing...'
          ) : isEditing ? (
            <>
              <Pencil size={13} />
              Update relationship
            </>
          ) : (
            <>
              <Plus size={13} />
              Write relationship
            </>
          )}
        </button>
      </div>
    </div>
  )
}
