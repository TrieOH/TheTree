import { useMutation, useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'
import { useState, useMemo, useEffect } from 'react'
import { schemaQueryOptions } from '#/features/schema/api'
import { checkRelationship } from '../api'
import { parseSpiceDBSchema } from '../lib/schema-parser'
import type { RelationshipFormState } from '../model'
import { FieldLabel } from './field-label'
import {
  ArrowRight,
  FileText,
  User,
  ShieldCheck,
  ShieldX,
  ShieldQuestion,
  ChevronRight
} from 'lucide-react'
import CustomSelect from '#/shared/ui/custom-select'

const inputCls =
  'h-9 w-full rounded-md border border-input bg-background px-3 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring transition-shadow touch-manipulation'

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
  const [isMobileResultOpen, setIsMobileResultOpen] = useState(false)

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
      const firstRelation = parsedSchema.relationsByDefinition[firstDefinition]?.[0]?.name ?? ''
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
      setIsMobileResultOpen(true)
      if (allowed) toast.success('Permission granted')
      else toast.error('Permission denied')
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Error checking permission')
    }
  })

  function set(key: keyof RelationshipFormState, value: string) {
    setForm((prev) => ({ ...prev, [key]: value }))
    setResult(null)
  }

  function handleSubmit() {
    if (!form.resourceId.trim() || !form.relation || !form.subjectId.trim())
      return
    checkMutation.mutate(form)
  }

  const availableRelations = parsedSchema.relationsByDefinition[form.resource] || []
  const selectedRelation = availableRelations.find((r) => r.name === form.relation)
  const allowedSubjects = selectedRelation?.allowedSubjectTypes || []

  useEffect(() => {
    if (form.relation && allowedSubjects.length > 0) {
      if (!allowedSubjects.includes(form.subject)) {
        set('subject', allowedSubjects[0])
      }
    }
  }, [form.relation, allowedSubjects, form.subject])

  const mobilePreview = `${form.resource || 'resource'}:${form.resourceId || 'id'}#${form.relation || 'relation'}@${form.subject || 'subject'}:${form.subjectId || 'id'}`

  return (
    <main className="h-(--content-height) flex flex-col bg-background border-l">
      <div className="flex items-center border-b px-4 py-3 bg-background shrink-0 min-w-0 gap-3">
        <div className="flex items-center gap-2 shrink-0">
          <ShieldCheck size={18} className="text-primary" />
          <span className="text-sm font-bold whitespace-nowrap hidden sm:inline">Check Permission:</span>
        </div>
        <div className="flex-1 min-w-0 font-mono text-xs sm:text-sm overflow-hidden">
          <div className="truncate text-muted-foreground sm:text-foreground">
            {mobilePreview}
          </div>
        </div>
      </div>

      <div className="flex flex-1 overflow-hidden relative">
        {/* Form Section */}
        <div className={`flex flex-col w-full md:w-96 shrink-0 bg-background transition-transform duration-300 ease-out absolute md:relative inset-0 md:translate-x-0 z-10 ${isMobileResultOpen ? '-translate-x-full md:translate-x-0' : 'translate-x-0'}`}>
          <div className="flex flex-1 flex-col gap-4 overflow-y-auto p-4 sm:p-5">
            {/* Resource type */}
            <div>
              <FieldLabel icon={<FileText size={12} />}>Resource type</FieldLabel>
              <CustomSelect
                value={form.resource}
                onChange={(val) => {
                  setForm(prev => ({
                    ...prev,
                    resource: val ?? "",
                    relation: ''
                  }))
                  setResult(null)
                }}
                options={parsedSchema.definitions}
                placeholder="Select resource..."
              />
            </div>

            {/* Resource ID */}
            <div>
              <FieldLabel icon={<FileText size={12} />}>Resource ID</FieldLabel>
              <input
                type="text"
                inputMode="text"
                value={form.resourceId}
                onChange={(e) => set('resourceId', e.target.value)}
                placeholder="e.g. roadmap"
                className={inputCls}
              />
            </div>

            {/* Permission/Relation */}
            <div>
              <FieldLabel icon={<ArrowRight size={12} />}>Permission</FieldLabel>
              <CustomSelect
                value={form.relation}
                onChange={val => set('relation', val ?? "")}
                options={availableRelations.map(item => item.name)}
                placeholder="Select relation..."
              />
            </div>

            {/* Subject type */}
            <div>
              <FieldLabel icon={<User size={12} />}>Subject type</FieldLabel>
              <CustomSelect
                value={form.subject}
                onChange={val => set('subject', val ?? "")}
                options={allowedSubjects.length > 0 ? allowedSubjects : parsedSchema.definitions}
                placeholder="Select subject..."
              />
            </div>

            {/* Subject ID */}
            <div>
              <FieldLabel icon={<User size={12} />}>Subject ID</FieldLabel>
              <input
                type="text"
                inputMode="text"
                value={form.subjectId}
                onChange={(e) => set('subjectId', e.target.value)}
                placeholder="e.g. alice"
                className={inputCls}
              />
            </div>
          </div>

          <div className="border-t p-4 bg-background safe-area-pb">
            <button
              onClick={handleSubmit}
              disabled={checkMutation.isPending}
              className="inline-flex h-9 w-full items-center justify-center gap-2 rounded-lg bg-primary px-4 text-sm font-medium text-primary-foreground shadow-sm transition-all hover:bg-primary/90 active:scale-[0.98] disabled:opacity-50 touch-manipulation"
            >
              {checkMutation.isPending ? (
                <>
                  <span className="w-4 h-4 border-2 border-primary-foreground/30 border-t-primary-foreground rounded-full animate-spin" />
                  Checking...
                </>
              ) : (
                <>
                  <ShieldCheck size={13} />
                  Check Permission
                </>
              )}
            </button>
          </div>
        </div>

        {/* Result Section */}
        <div className={`absolute md:relative inset-0 md:flex-1 flex flex-col items-center justify-center p-4 sm:p-6 bg-muted/30 transition-transform duration-300 ease-out md:translate-x-0 ${isMobileResultOpen ? 'translate-x-0' : 'translate-x-full md:translate-x-0'}`}>
          {!result && !checkMutation.isPending && (
            <div className="flex flex-col items-center text-center max-w-sm animate-in fade-in zoom-in duration-300 px-4">
              <div className="w-16 h-16 rounded-full bg-background border flex items-center justify-center mb-4 shadow-sm">
                <ShieldQuestion size={32} className="text-muted-foreground" />
              </div>
              <h2 className="text-lg sm:text-xl font-semibold mb-2">Ready to Check</h2>
              <p className="text-muted-foreground text-sm sm:text-base">
                Fill out the form to verify if a subject has permission on a resource.
              </p>
            </div>
          )}

          {checkMutation.isPending && (
            <div className="flex flex-col items-center text-center">
              <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-4 animate-pulse">
                <ShieldCheck size={32} className="text-primary" />
              </div>
              <h2 className="text-lg sm:text-xl font-semibold mb-2">Checking...</h2>
            </div>
          )}

          {result && (
            <div className={`flex flex-col items-center text-center p-6 sm:p-8 rounded-xl border-2 max-w-sm w-full bg-background shadow-lg animate-in fade-in slide-in-from-bottom-4 duration-300 ${result.allowed ? 'border-green-500/20' : 'border-red-500/20'
              }`}>
              <div className={`w-16 h-16 sm:w-20 sm:h-20 rounded-full flex items-center justify-center mb-4 sm:mb-6 shadow-sm ${result.allowed ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'
                }`}>
                {result.allowed ? <ShieldCheck size={32} className="sm:w-10 sm:h-10" /> : <ShieldX size={32} className="sm:w-10 sm:h-10" />}
              </div>

              <h2 className={`text-xl sm:text-2xl font-bold mb-2 ${result.allowed ? 'text-green-600' : 'text-red-600'
                }`}>
                {result.allowed ? 'ALLOWED' : 'DENIED'}
              </h2>

              <div className="space-y-3 mt-4 sm:mt-6 text-sm text-left w-full border-t pt-4 sm:pt-6">
                <div className="flex justify-between items-center gap-2">
                  <span className="text-muted-foreground shrink-0">Resource:</span>
                  <span className="font-mono font-medium bg-muted/50 px-2 py-0.5 rounded truncate">{form.resource}:{form.resourceId}</span>
                </div>
                <div className="flex justify-between items-center gap-2">
                  <span className="text-muted-foreground shrink-0">Permission:</span>
                  <span className="font-mono font-medium bg-muted/50 px-2 py-0.5 rounded truncate">{form.relation}</span>
                </div>
                <div className="flex justify-between items-center gap-2">
                  <span className="text-muted-foreground shrink-0">Subject:</span>
                  <span className="font-mono font-medium bg-muted/50 px-2 py-0.5 rounded truncate">{form.subject}:{form.subjectId}</span>
                </div>
              </div>

              {/* Mobile action button */}
              <button
                onClick={() => setIsMobileResultOpen(false)}
                className="md:hidden mt-6 w-full h-11 flex items-center justify-center gap-2 rounded-lg border border-input bg-background text-sm font-medium hover:bg-muted transition-colors"
              >
                <ChevronRight size={16} className="rotate-180" />
                Back to Form
              </button>
            </div>
          )}
        </div>
      </div>
    </main>
  )
}