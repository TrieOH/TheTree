// routes/admin/events/$eventId/editions/index.tsx — Editions list + create
import { createFileRoute, Link, useParams } from '@tanstack/react-router'
import React, { useState } from 'react'
import type { EditionCreateI, EditionI } from '@/features/editions/model'
import type { EventI } from '@/features/events/model'
import { getOwnEventsFn } from '@/features/events/api'
import {
  createEditionFn, getAllAdminEditionsFn, publishEditionFn,
  disconnectPaymentAccountToEditionFn,
} from '@/features/editions/api'
import { connectEditionSellerToWorkspaceFn, disconnectEditionSellerToWorkspaceFn } from '@/features/payments/api'
import { formatDateForDatetimeLocal, parseDatetimeLocal } from '@/shared/lib/date'
import { env } from '@/env'
import {
  AdminShell, PageHeader, FormField, ErrorMsg, StatusBadge,
  EmptyState, cardClass, inputClass, btnPrimary, btnSecondary, btnDanger
} from '@/shared/ui'

export const Route = createFileRoute('/admin/events/$eventId/editions/')({
  component: RouteComponent,
})

const defaultTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone
const timezones = Intl.supportedValuesOf('timeZone')

const INITIAL_FORM = (): EditionCreateI => ({
  type: 'year',
  edition_name: '',
  tagline: undefined,
  description: undefined,
  registration_opens_at: undefined,
  registration_closes_at: undefined,
  starts_at: formatDateForDatetimeLocal(new Date()),
  ends_at: formatDateForDatetimeLocal(new Date(Date.now() + 86400000)),
  timezone: defaultTimezone,
  location_name: '',
  location_address: '',
  logo_url: undefined,
  banner_url: undefined,
  contact_email: undefined,
  contact_phone: undefined,
  organizer_name: undefined,
})

function RouteComponent() {
  const { eventId } = useParams({ from: '/admin/events/$eventId/editions/' })
  const [event, setEvent] = useState<EventI | null>(null)
  const [editions, setEditions] = useState<EditionI[]>([])
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState<EditionCreateI>(INITIAL_FORM())
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [refresh, setRefresh] = useState(0)

  React.useEffect(() => {
    getOwnEventsFn().then(evs => { setEvent(evs.find(e => e.id === eventId) ?? null); })
    getAllAdminEditionsFn(eventId).then(setEditions)
  }, [eventId, refresh])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setForm(f => ({ ...f, [name]: value } as unknown as EditionCreateI))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const data = {
        ...form,
        starts_at: parseDatetimeLocal(form.starts_at).toISOString(),
        ends_at: parseDatetimeLocal(form.ends_at).toISOString(),
        registration_opens_at: form.registration_opens_at ? parseDatetimeLocal(form.registration_opens_at).toISOString() : undefined,
        registration_closes_at: form.registration_closes_at ? parseDatetimeLocal(form.registration_closes_at).toISOString() : undefined,
      }
      const res = await createEditionFn(data, eventId)
      if (res.success) {
        setForm(INITIAL_FORM())
        setShowForm(false)
        setRefresh(r => r + 1)
      } else throw new Error(res.message)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Erro ao criar edição')
    } finally {
      setLoading(false)
    }
  }

  const handleConnect = async (editionId: string) => {
    const res = await connectEditionSellerToWorkspaceFn({
      data: {
        provider: 'mercadopago',
        workspace_name: 'Univents',
        final_redirect_url: `${window.location.origin}/admin/events/${eventId}/editions/${editionId}/callback/payment`,
        provider_redirect_url: env.VITE_MERCADO_PAGO_CALLBACK_URL,
      },
    })
    if (res.success) window.location.href = res.data.redirect_url
  }

  const handleDisconnect = async (editionId: string, credentialId: string) => {
    const ws = await disconnectEditionSellerToWorkspaceFn({ data: { workspace_name: 'Univents', credential_id: credentialId } })
    if (!ws.success) return
    const res = await disconnectPaymentAccountToEditionFn(eventId, editionId)
    if (res.success) setRefresh(r => r + 1)
  }

  const links = (editionId: string) => [
    // { label: 'Atividades', to: '/admin/events/$eventId/editions/$editionId/activities' as const },
    // { label: 'Ingressos', to: '/admin/events/$eventId/editions/$editionId/tickets' as const },
    { label: 'Produtos', to: '/events/$eventId/editions/$editionId/products' as const },
    { label: 'Checkpoints', to: '/admin/events/$eventId/editions/$editionId/checkpoints/' as const },
  ]

  return (
    <AdminShell
      breadcrumbs={
        <div className="flex items-center gap-2 text-sm text-gray-500">
          <Link to="/admin/events" className="hover:text-gray-900 transition-colors">Eventos</Link>
          <span className="text-gray-300">/</span>
          <span className="text-gray-700">{event?.name ?? '...'}</span>
          <span className="text-gray-300">/</span>
          <span className="text-gray-900">Edições</span>
        </div>
      }
    >
      <PageHeader
        title="Edições"
        subtitle={event ? `${event.name} · ${editions.length} edição${editions.length !== 1 ? 'ões' : ''}` : ''}
        action={
          <button className={btnPrimary} onClick={() => { setShowForm(s => !s); }}>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" />
            </svg>
            Nova edição
          </button>
        }
      />

      {showForm && (
        <div className="bg-white border border-gray-100 rounded-xl p-6 mb-6 shadow-sm">
          <h2 className="text-sm font-semibold text-gray-900 mb-5">Criar edição</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <FormField label="Tipo" required>
              <select className={inputClass} name="type" value={form.type} onChange={handleChange}>
                <option value="year">Ano</option>
                <option value="season">Temporada</option>
                <option value="number">Número</option>
                <option value="ordinal">Ordinal</option>
                <option value="custom">Personalizado</option>
              </select>
            </FormField>
            <FormField label="Nome da edição" required>
              <input className={inputClass} name="edition_name" value={form.edition_name} onChange={handleChange} placeholder="Ex: 2025" />
            </FormField>
            <FormField label="Início" required>
              <input className={inputClass} type="datetime-local" name="starts_at" value={form.starts_at} onChange={handleChange} />
            </FormField>
            <FormField label="Fim" required>
              <input className={inputClass} type="datetime-local" name="ends_at" value={form.ends_at} onChange={handleChange} />
            </FormField>
            <FormField label="Abertura de inscrições">
              <input className={inputClass} type="datetime-local" name="registration_opens_at" value={form.registration_opens_at ?? ''} onChange={handleChange} />
            </FormField>
            <FormField label="Fechamento de inscrições">
              <input className={inputClass} type="datetime-local" name="registration_closes_at" value={form.registration_closes_at ?? ''} onChange={handleChange} />
            </FormField>
            <FormField label="Nome do local" required>
              <input className={inputClass} name="location_name" value={form.location_name} onChange={handleChange} placeholder="Ex: Centro de Convenções" />
            </FormField>
            <FormField label="Endereço" required>
              <input className={inputClass} name="location_address" value={form.location_address} onChange={handleChange} placeholder="Rua, número, cidade" />
            </FormField>
            <div className="sm:col-span-2">
              <FormField label="Fuso horário" required>
                <select className={inputClass} name="timezone" value={form.timezone} onChange={handleChange}>
                  {timezones.map(tz => <option key={tz} value={tz}>{tz}</option>)}
                </select>
              </FormField>
            </div>
            <FormField label="Tagline">
              <input className={inputClass} name="tagline" value={form.tagline ?? ''} onChange={handleChange} />
            </FormField>
            <FormField label="Organizador">
              <input className={inputClass} name="organizer_name" value={form.organizer_name ?? ''} onChange={handleChange} />
            </FormField>

            <ErrorMsg msg={error} />
            <div className="sm:col-span-2 flex gap-2 pt-2">
              <button type="submit" className={btnPrimary} disabled={loading}>
                {loading ? 'Criando...' : 'Criar edição'}
              </button>
              <button type="button" className={btnSecondary} onClick={() => { setShowForm(false); }}>Cancelar</button>
            </div>
          </form>
        </div>
      )}

      {editions.length === 0 ? (
        <EmptyState
          icon={<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><rect x="3" y="4" width="18" height="18" rx="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>}
          title="Nenhuma edição ainda"
          description="Crie a primeira edição para este evento."
        />
      ) : (
        <div className="grid gap-3">
          {editions.map(ed => (
            <div key={ed.id} className={cardClass}>
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0">
                  <div className="flex items-center gap-2 flex-wrap mb-1">
                    <span className="text-sm font-medium text-gray-900">{ed.edition_name}</span>
                    <StatusBadge status={ed.status} />
                    <span className="text-xs text-gray-400">{ed.type}</span>
                    <span className="text-xs text-gray-400 bg-gray-50 px-2 py-0.5 rounded">{ed.monetary_type}</span>
                  </div>
                  <p className="text-xs text-gray-400">
                    {new Date(ed.starts_at).toLocaleDateString('pt-BR')} → {new Date(ed.ends_at).toLocaleDateString('pt-BR')} · {ed.location_name}
                  </p>
                </div>

                <div className="flex items-center gap-2 shrink-0">
                  {ed.status === 'draft' && (
                    <button
                      className={btnSecondary + ' text-xs py-1.5'}
                      onClick={async () => { await publishEditionFn(eventId, ed.id); setRefresh(r => r + 1) }}
                    >
                      Publicar
                    </button>
                  )}
                  {ed.trie_payments_credential_id ? (
                    <button
                      className={btnDanger + ' text-xs py-1.5'}
                      onClick={() => handleDisconnect(ed.id, ed.trie_payments_credential_id!)}
                    >
                      Desconectar pagamento
                    </button>
                  ) : (
                    <button className={btnSecondary + ' text-xs py-1.5'} onClick={() => handleConnect(ed.id)}>
                      Conectar pagamento
                    </button>
                  )}
                </div>
              </div>

              {/* Sub-nav links */}
              <div className="flex gap-2 mt-4 pt-4 border-t border-gray-50 flex-wrap">
                {links(ed.id).map(({ label, to }) => (
                  <Link
                    key={label}
                    to={to}
                    params={{ eventId: ed.event_id, editionId: ed.id }}
                    className="text-xs text-gray-500 bg-gray-50 hover:bg-gray-100 px-3 py-1.5 rounded-md transition-colors font-medium"
                  >
                    {label}
                  </Link>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </AdminShell>
  )
}