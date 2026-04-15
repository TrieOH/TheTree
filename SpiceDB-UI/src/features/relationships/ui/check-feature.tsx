import { useMutation, useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'
import { useState, useMemo, useEffect } from 'react'
import { schemaQueryOptions } from '#/features/schema/api'
import { checkRelationship } from '../api'
import { parseSpiceDBSchema } from '../lib/schema-parser'
import type { RelationshipFormState } from '../model'
import { FieldLabel } from './field-label'
import { ArrowRight, FileText, User, ShieldCheck, ShieldX, ShieldQuestion } from 'lucide-react'

const inputCls =
  'h-9 w-full rounded-md border border-input bg-background px-3 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring transition-shadow'

const selectCls =
  'h-9 w-full rounded-md border border-input bg-background px-2.5 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring transition-shadow'

const INITIAL_FORM: RelationshipFormState = {
  resource: '',
  resourceId: '',
  relation: '',
  subject: '',
  subjectId: '',
}

export function CheckFeature({ envId }: { envId: string }) {
  const [form, setForm] = useState<RelationshipFormState>(INITIAL_FORM)
  const [result, setResult] = useState<{ allowed: boolean } | null>(null)

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
    if (parsedSchema.definitions.length > 0 && form.resource === '') {
      const firstDefinition = parsedSchema.definitions[0]
      const firstRelation = parsedSchema.relationsByDefinition[firstDefinition]?.[0] ?? ''
      setForm({
        ...INITIAL_FORM,
        resource: firstDefinition,
        subject: firstDefinition,
        relation: firstRelation,
      })
    }
  }, [parsedSchema.definitions, parsedSchema.relationsByDefinition, form.resource])

  const checkMutation = useMutation({
    mutationFn: (data: RelationshipFormState) =>
      checkRelationship({ data: { ...data, envId } }),
    onSuccess: (res) => {
      const allowed = res.success && res.data.permissionship === "PERMISSIONSHIP_HAS_PERMISSION"
      setResult({ allowed })
      if (allowed) toast.success('Permission granted')
      else toast.error('Permission denied')
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Error checking permission')
    }
  })

  function set(key: keyof RelationshipFormState, value: string) {
    setForm((prev) => ({ ...prev, [key]: value }))
    setResult(null) // Reset result when form changes
  }

  function handleSubmit() {
    if (!form.resourceId.trim() || !form.relation || !form.subjectId.trim())
      return
    checkMutation.mutate(form)
  }

  const availableRelations = parsedSchema.relationsByDefinition[form.resource] || []

  return (
    <main className="h-(--content-height) flex flex-col md:flex-row border-l overflow-hidden">
      {/* Form Section */}
      <div className="flex flex-col border-r w-full md:w-96 shrink-0 bg-background">
        <div className="flex items-center gap-2 border-b px-5 py-3.5">
          <ShieldCheck size={16} className="text-primary" />
          <span className="text-sm font-medium">Check Permission</span>
        </div>

        <div className="flex flex-1 flex-col gap-5 overflow-y-auto p-5">
          {/* Resource type */}
          <div>
            <FieldLabel icon={<FileText size={11} />}>Resource type</FieldLabel>
            <select
              value={form.resource}
              onChange={(e) => {
                const newResource = e.target.value
                setForm(prev => ({
                  ...prev,
                  resource: newResource,
                  relation: ''
                }))
                setResult(null)
              }}
              className={selectCls}
            >
              {parsedSchema.definitions.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
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
            />
          </div>

          {/* Permission/Relation */}
          <div>
            <FieldLabel icon={<ArrowRight size={11} />}>Permission / Relation</FieldLabel>
            <select
              value={form.relation}
              onChange={(e) => set('relation', e.target.value)}
              className={selectCls}
            >
              <option value="">Select permission</option>
              {availableRelations.map((r) => (
                <option key={r} value={r}>{r}</option>
              ))}
            </select>
          </div>

          {/* Subject type */}
          <div>
            <FieldLabel icon={<User size={11} />}>Subject type</FieldLabel>
            <select
              value={form.subject}
              onChange={(e) => set('subject', e.target.value)}
              className={selectCls}
            >
              {parsedSchema.definitions.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
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
            />
          </div>
        </div>

        <div className="border-t p-4">
          <button
            onClick={handleSubmit}
            disabled={checkMutation.isPending}
            className="inline-flex h-9 w-full items-center justify-center gap-2 rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50"
          >
            {checkMutation.isPending ? 'Checking...' : 'Check Permission'}
          </button>
        </div>
      </div>

      {/* Result Section */}
      <div className="flex-1 flex flex-col items-center justify-center p-10 bg-muted/30">
        {!result && !checkMutation.isPending && (
          <div className="flex flex-col items-center text-center max-w-md">
            <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4">
              <ShieldQuestion size={32} className="text-muted-foreground" />
            </div>
            <h2 className="text-xl font-semibold mb-2">Ready to Check</h2>
            <p className="text-muted-foreground">
              Fill out the form on the left to verify if a subject has a specific permission on a resource.
            </p>
          </div>
        )}

        {checkMutation.isPending && (
          <div className="flex flex-col items-center text-center">
            <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-4 animate-pulse">
              <ShieldCheck size={32} className="text-primary" />
            </div>
            <h2 className="text-xl font-semibold mb-2">Checking...</h2>
          </div>
        )}

        {result && (
          <div className={`flex flex-col items-center text-center p-8 rounded-xl border-2 max-w-md w-full bg-background shadow-sm ${result.allowed ? 'border-green-500/20' : 'border-red-500/20'
            }`}>
            <div className={`w-20 h-20 rounded-full flex items-center justify-center mb-6 ${result.allowed ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'
              }`}>
              {result.allowed ? <ShieldCheck size={40} /> : <ShieldX size={40} />}
            </div>

            <h2 className={`text-2xl font-bold mb-2 ${result.allowed ? 'text-green-600' : 'text-red-600'
              }`}>
              {result.allowed ? 'ALLOWED' : 'DENIED'}
            </h2>

            <div className="space-y-3 mt-6 text-sm text-left w-full border-t pt-6">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Subject:</span>
                <span className="font-mono font-medium">{form.subject}:{form.subjectId}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Permission:</span>
                <span className="font-mono font-medium">{form.relation}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Resource:</span>
                <span className="font-mono font-medium">{form.resource}:{form.resourceId}</span>
              </div>
            </div>
          </div>
        )}
      </div>
    </main>
  )
}
