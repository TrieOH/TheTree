import { createFileRoute, Link } from '@tanstack/react-router'
import React, { useState } from 'react'
import type { EventCreateI, EventI } from '@/features/events/model'
import { createEventFn, getOwnEventsFn, publishEventFn } from '@/features/events/api'
import {
  AdminShell, PageHeader, FormField, ErrorMsg, StatusBadge,
  EmptyState, cardClass, inputClass, btnPrimary, btnSecondary
} from '@/shared/ui'

export const Route = createFileRoute('/admin/events/')({
  component: RouteComponent,
})

const INITIAL_FORM: EventCreateI = {
  organization_id: undefined,
  name: '',
  acronym: undefined,
  slug: '',
  tagline: undefined,
  description: undefined,
  is_series: false,
  logo_url: undefined,
  banner_url: undefined,
  contact_email: '',
}

function RouteComponent() {
  const [events, setEvents] = useState<EventI[]>([])
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState<EventCreateI>(INITIAL_FORM)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [refresh, setRefresh] = useState(0)

  React.useEffect(() => {
    getOwnEventsFn().then(setEvents)
  }, [refresh])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const t = e.target as HTMLInputElement
    setForm(f => ({ ...f, [t.name]: t.type === 'checkbox' ? t.checked : t.value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const res = await createEventFn(form)
      if (res.success) {
        setForm(INITIAL_FORM)
        setShowForm(false)
        setRefresh(r => r + 1)
      } else throw new Error(res.message)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Erro ao criar evento')
    } finally {
      setLoading(false)
    }
  }

  return (
    <AdminShell>
      <PageHeader
        title="Eventos"
        subtitle={`${events.length} evento${events.length !== 1 ? 's' : ''}`}
        action={
          <button className={btnPrimary} onClick={() => setShowForm(s => !s)}>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" />
            </svg>
            Novo evento
          </button>
        }
      />

      {/* Create form */}
      {showForm && (
        <div className="bg-white border border-gray-100 rounded-xl p-6 mb-6 shadow-sm">
          <h2 className="text-sm font-semibold text-gray-900 mb-5">Criar evento</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <FormField label="Nome" required>
              <input className={inputClass} name="name" value={form.name} onChange={handleChange} placeholder="Ex: TechConf" />
            </FormField>
            <FormField label="Slug" required>
              <input className={inputClass} name="slug" value={form.slug} onChange={handleChange} placeholder="ex: tech-conf" />
            </FormField>
            <FormField label="Sigla">
              <input className={inputClass} name="acronym" value={form.acronym ?? ''} onChange={handleChange} placeholder="TC" />
            </FormField>
            <FormField label="E-mail de contato" required>
              <input className={inputClass} type="email" name="contact_email" value={form.contact_email ?? ''} onChange={handleChange} placeholder="contato@email.com" />
            </FormField>
            <FormField label="Tagline">
              <input className={inputClass} name="tagline" value={form.tagline ?? ''} onChange={handleChange} placeholder="Uma frase curta sobre o evento" />
            </FormField>
            <FormField label="Organization ID">
              <input className={inputClass} name="organization_id" value={form.organization_id ?? ''} onChange={handleChange} />
            </FormField>
            <div className="sm:col-span-2">
              <FormField label="Descrição">
                <textarea className={inputClass + ' resize-none'} rows={3} name="description" value={form.description ?? ''} onChange={handleChange} placeholder="Descreva o evento..." />
              </FormField>
            </div>
            <div className="sm:col-span-2">
              <FormField label="Logo URL">
                <input className={inputClass} name="logo_url" value={form.logo_url ?? ''} onChange={handleChange} />
              </FormField>
            </div>
            <div className="sm:col-span-2">
              <FormField label="Banner URL">
                <input className={inputClass} name="banner_url" value={form.banner_url ?? ''} onChange={handleChange} />
              </FormField>
            </div>
            <div className="sm:col-span-2 flex items-center gap-2">
              <input
                type="checkbox" id="is_series" name="is_series"
                checked={form.is_series}
                onChange={handleChange}
                className="w-4 h-4 rounded border-gray-300 accent-gray-900"
              />
              <label htmlFor="is_series" className="text-sm text-gray-600">É uma série de eventos</label>
            </div>

            <ErrorMsg msg={error} />

            <div className="sm:col-span-2 flex gap-2 pt-2">
              <button type="submit" className={btnPrimary} disabled={loading}>
                {loading ? 'Criando...' : 'Criar evento'}
              </button>
              <button type="button" className={btnSecondary} onClick={() => setShowForm(false)}>
                Cancelar
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Events list */}
      {events.length === 0 ? (
        <EmptyState
          icon={<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M8 2v4M16 2v4M3 10h18M5 4h14a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2z" /></svg>}
          title="Nenhum evento ainda"
          description="Crie seu primeiro evento para começar."
        />
      ) : (
        <div className="grid gap-3">
          {events.map(ev => (
            <div key={ev.id} className={cardClass + ' flex items-center justify-between'}>
              <div className="flex items-center gap-4 min-w-0">
                <div className="w-9 h-9 rounded-lg bg-gray-100 flex items-center justify-center text-gray-500 shrink-0 text-xs font-bold">
                  {ev.acronym ? ev.acronym.slice(0, 2).toUpperCase() : ev.name.slice(0, 2).toUpperCase()}
                </div>
                <div className="min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-medium text-gray-900 truncate">{ev.name}</span>
                    <StatusBadge status={ev.status} />
                    {ev.is_series && (
                      <span className="text-xs text-gray-400 bg-gray-100 px-2 py-0.5 rounded-md">série</span>
                    )}
                  </div>
                  <p className="text-xs text-gray-400 truncate">{ev.slug} · {ev.contact_email} · {ev.editions_count} edição{ev.editions_count !== 1 ? 'ões' : ''}</p>
                </div>
              </div>

              <div className="flex items-center gap-2 shrink-0 ml-4">
                {ev.status === 'draft' && (
                  <button
                    className={btnSecondary + ' text-xs py-1.5'}
                    onClick={async () => { await publishEventFn(ev.id); setRefresh(r => r + 1) }}
                  >
                    Publicar
                  </button>
                )}
                <Link
                  to="/admin/events/$eventId/editions"
                  params={{ eventId: ev.id }}
                  className={btnPrimary + ' text-xs py-1.5'}
                >
                  Ver edições
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M5 12h14M12 5l7 7-7 7" />
                  </svg>
                </Link>
              </div>
            </div>
          ))}
        </div>
      )}
    </AdminShell>
  )
}