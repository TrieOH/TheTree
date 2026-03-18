import { createFileRoute, Link } from '@tanstack/react-router'
import { SignIn, SignUp, useAuth } from '@soramux/node-auth-sdk/react'
import React, { useState } from 'react'
import { useServerFn } from '@tanstack/react-start'
import type { EventCreateI, EventI } from '@/features/events/model'
import type { EditionCreateI, EditionI } from '@/features/editions/model'
import type { ActivityCreateI, ActivityI } from '@/features/activities/model'
import type { TicketCreateI, TicketI } from '@/features/tickets/model'
import type { ProductCreateI, ProductI } from '@/features/products/model'
import type { CheckpointCreateI, CheckpointI } from '@/features/checkpoints/model'
import { createEventFn, getOwnEventsFn, publishEventFn } from '@/features/events/api'
import { createEditionFn, getAllAdminEditionsFn, publishEditionFn, disconnectPaymentAccountToEditionFn } from '@/features/editions/api'
import { createActivityFn, getAllAdminActivitiesFn } from '@/features/activities/api'
import { createTicketFn, getAllTicketsFn } from '@/features/tickets/api'
import { createProductFn, getAllAdminProductsFn, publishProductFn } from '@/features/products/api'
import { createCheckpointFn, getAllCheckpointsFn } from '@/features/checkpoints/api'
import {
  formatDateForDatetimeLocal,
  parseDatetimeLocal,
} from '@/shared/lib/date'
import { connectEditionSellerToWorkspaceFn, disconnectEditionSellerToWorkspaceFn } from '@/features/payments/api'
import { env } from '@/env'

export const Route = createFileRoute('/temp/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { isAuthenticated } = useAuth()

  const timezones = Intl.supportedValuesOf('timeZone')
  const defaultTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone

  // events list and selection
  const [events, setEvents] = useState<EventI[]>([])
  const [selectedEventId, setSelectedEventId] = useState<string | null>(null)
  const [result, setResult] = useState<EventI | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [eventCreated, setEventCreated] = useState<boolean>(false);


  React.useEffect(() => {
    const fetchEvents = async () => {
      const evs = await getOwnEventsFn();
      setEvents(evs);
      setEventCreated(false); // Reset after fetching
    };
    fetchEvents();
  }, [eventCreated]);

  // event creation form
  const [form, setForm] = useState<EventCreateI>({
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
  })

  // editions-related state
  const [editionForm, setEditionForm] = useState<EditionCreateI>({
    type: 'year',
    edition_name: '',
    tagline: undefined,
    description: undefined,
    registration_opens_at: undefined,
    registration_closes_at: undefined,
    starts_at: formatDateForDatetimeLocal(new Date()),
    ends_at: formatDateForDatetimeLocal(
      new Date(new Date().getTime() + 24 * 60 * 60 * 1000),
    ),
    timezone: defaultTimezone,
    location_name: '',
    location_address: '',
    logo_url: undefined,
    banner_url: undefined,
    contact_email: undefined,
    contact_phone: undefined,
    organizer_name: undefined,
  })
  const [editions, setEditions] = useState<EditionI[]>([])
  const [editionError, setEditionError] = useState<string | null>(null)
  const [editionCreated, setEditionCreated] = useState<boolean>(false);
  const [selectedEditionId, setSelectedEditionId] = useState<string | null>(null);

  // activities-related state
  const [activityForm, setActivityForm] = useState<ActivityCreateI>({
    title: '',
    description: undefined,
    location: '',
    starts_at: formatDateForDatetimeLocal(new Date()),
    ends_at: formatDateForDatetimeLocal(
      new Date(new Date().getTime() + 60 * 60 * 1000), // 1 hour after start
    ),
    presenter_name: undefined,
    token_cost: 0,
    has_capacity: false,
    capacity: 0,
    difficulty: 'no_prerequisites',
  })
  const [activities, setActivities] = useState<ActivityI[]>([])
  const [activityError, setActivityError] = useState<string | null>(null)
  const [activityCreated, setActivityCreated] = useState<boolean>(false);

  // tickets-related state
  const [ticketForm, setTicketForm] = useState<TicketCreateI>({
    name: '',
    description: undefined,
  })
  const [tickets, setTickets] = useState<TicketI[]>([])
  const [ticketError, setTicketError] = useState<string | null>(null)
  const [ticketCreated, setTicketCreated] = useState<boolean>(false);

  // products-related state
  const [productForm, setProductForm] = useState<ProductCreateI>({
    edition_scope_id: '', // Will be set from selectedEditionId
    name: '',
    description: undefined,
    type: 'merchandise',
    ticket_id: undefined, // Will be set from selectedTicketId if applicable
    price_cents: 0,
    available_from: undefined,
    available_until: undefined,
    has_inventory: false,
    inventory_quantity: 0,
  })
  const [products, setProducts] = useState<ProductI[]>([])
  const [productError, setProductError] = useState<string | null>(null)
  const [productCreated, setProductCreated] = useState<boolean>(false);

  // checkpoints-related state
  const [checkpointForm, setCheckpointForm] = useState<CheckpointCreateI>({
    name: '',
    access_mode: 'open',
    type: 'entry',
    starts_at: undefined,
    ends_at: undefined,
  })
  const [checkpoints, setCheckpoints] = useState<CheckpointI[]>([])
  const [checkpointError, setCheckpointError] = useState<string | null>(null)
  const [checkpointCreated, setCheckpointCreated] = useState<boolean>(false);

  React.useEffect(() => {
    const fetchCheckpoints = async () => {
      if (selectedEventId && selectedEditionId) {
        const checkpts = await getAllCheckpointsFn(selectedEventId, selectedEditionId);
        setCheckpoints(checkpts);
        setCheckpointCreated(false); // Reset after fetching
      } else {
        setCheckpoints([]);
      }
    };
    fetchCheckpoints();
  }, [selectedEventId, selectedEditionId, checkpointCreated]);

  React.useEffect(() => {
    const fetchProducts = async () => {
      if (selectedEventId && selectedEditionId) {
        const prods = await getAllAdminProductsFn(selectedEventId, selectedEditionId);
        setProducts(prods);
        setProductCreated(false); // Reset after fetching
      } else {
        setProducts([]);
      }
    };
    fetchProducts();
  }, [selectedEventId, selectedEditionId, productCreated]);

  React.useEffect(() => {
    const fetchTickets = async () => {
      if (selectedEventId && selectedEditionId) {
        const tix = await getAllTicketsFn(selectedEventId, selectedEditionId);
        setTickets(tix);
        setTicketCreated(false); // Reset after fetching
      } else {
        setTickets([]);
      }
    };
    fetchTickets();
  }, [selectedEventId, selectedEditionId, ticketCreated]);

  React.useEffect(() => {
    const fetchActivities = async () => {
      if (selectedEventId && selectedEditionId) {
        const acts = await getAllAdminActivitiesFn(selectedEventId, selectedEditionId);
        setActivities(acts);
        setActivityCreated(false); // Reset after fetching
      } else {
        setActivities([]);
      }
    };
    fetchActivities();
  }, [selectedEventId, selectedEditionId, activityCreated]);

  React.useEffect(() => {
    const fetchEditions = async () => {
      if (selectedEventId) {
        const eds = await getAllAdminEditionsFn(selectedEventId);
        setEditions(eds);
        setEditionCreated(false); // Reset after fetching
      } else {
        setEditions([]);
      }
    };
    fetchEditions();
  }, [selectedEventId, editionCreated]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const target = e.target as HTMLInputElement
    const { name, value, type } = target
    const checked = type === 'checkbox' ? target.checked : undefined

    setForm((f) => ({
      ...f,
      [name]: type === 'checkbox' ? checked : value,
    }))
  }

  const handleEditionChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => {
    const target = e.target as HTMLInputElement | HTMLSelectElement
    const { name, value } = target
    setEditionForm((f) => ({
      ...f,
      [name]: value,
    } as unknown as EditionCreateI))
  }

  const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError(null)
    try {
      const res = await createEventFn(form)
      if (res.success) {
        setResult(res.data);
        setEventCreated(true); // Trigger re-fetch via useEffect
      } else {
        throw new Error(res.message)
      }
    } catch (err: unknown) {
      let errorMessage = 'request failed';
      if (err instanceof Error) {
        errorMessage = err.message;
      } else if (typeof err === 'object' && err !== null && 'message' in err && typeof (err as { message: unknown }).message === 'string') {
        errorMessage = (err as { message: string }).message;
      }
      setError(errorMessage);
    }

  }

  const handleConnectFn = useServerFn(connectEditionSellerToWorkspaceFn)

  const handleConnect = async (eventId: string, editionId: string) => {
    const res = await handleConnectFn({
      data: {
        provider: "mercadopago",
        workspace_name: "Univents",
        final_redirect_url: window.location.origin + `/admin/events/${eventId}/editions/${editionId}/callback/payment`,
        provider_redirect_url: env.VITE_MERCADO_PAGO_CALLBACK_URL
      }
    })
    if (res.success) {
      window.location.href = res.data.redirect_url
    } else console.error("DEU RUIM: ", res.message)
  }

  const handleDisconnect = async (
    eventId: string,
    editionId: string,
    credential_id: string | null
  ) => {
    if (!credential_id) return;
    const workspaceRes = await disconnectEditionSellerToWorkspaceFn({
      data: {
        workspace_name: "Univents",
        credential_id
      }
    });
    if (!workspaceRes.success) {
      console.error("Erro ao desconectar do workspace:", workspaceRes.message);
      return;
    }
    const res = await disconnectPaymentAccountToEditionFn(eventId, editionId);
    if (res.success) {
      setEditionCreated(true);
    } else {
      console.error("Erro ao desconectar:", res.message);
    }
  }

  const selectedEvent = events.find((ev) => ev.id === selectedEventId) ?? null

  const handleEditionSubmit = async (
    e: React.SyntheticEvent<HTMLFormElement>,
  ) => {
    e.preventDefault()
    setEditionError(null)
    if (!selectedEventId) {
      setEditionError('pick an event first')
      return
    }

    try {
      const data = {
        ...editionForm,
        starts_at: parseDatetimeLocal(editionForm.starts_at).toISOString(),
        ends_at: parseDatetimeLocal(editionForm.ends_at).toISOString(),
        registration_opens_at: editionForm.registration_opens_at
          ? parseDatetimeLocal(editionForm.registration_opens_at).toISOString()
          : undefined,
        registration_closes_at: editionForm.registration_closes_at
          ? parseDatetimeLocal(
            editionForm.registration_closes_at,
          ).toISOString()
          : undefined,
      }
      const res = await createEditionFn(data, selectedEventId)
      if (res.success) {
        setEditionCreated(true); // Trigger re-fetch via useEffect
      } else {
        throw new Error(res.message)
      }
    } catch (err: unknown) {
      let errorMessage = 'request failed';
      if (err instanceof Error) {
        errorMessage = err.message;
      } else if (typeof err === 'object' && err !== null && 'message' in err && typeof (err as { message: unknown }).message === 'string') {
        errorMessage = (err as { message: string }).message;
      }
      setEditionError(errorMessage);
    }
  }

  const handleChangeActivity = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => {
    const target = e.target as HTMLInputElement | HTMLSelectElement
    const { name, value, type } = target
    const checked = type === 'checkbox' ? (target).checked : undefined

    setActivityForm((f) => ({
      ...f,
      [name]: type === 'checkbox' ? checked : value,
    } as unknown as ActivityCreateI))
  }

  const handleActivitySubmit = async (
    e: React.SyntheticEvent<HTMLFormElement>,
  ) => {
    e.preventDefault()
    setActivityError(null)
    if (!selectedEventId) {
      setActivityError('pick an event first')
      return
    }
    if (!selectedEditionId) {
      setActivityError('pick an edition first')
      return
    }

    try {
      const data = {
        ...activityForm,
        starts_at: parseDatetimeLocal(activityForm.starts_at).toISOString(),
        ends_at: parseDatetimeLocal(activityForm.ends_at).toISOString(),
        token_cost: Number(activityForm.token_cost),
        capacity: Number(activityForm.capacity),
      }
      const res = await createActivityFn(data, selectedEventId, selectedEditionId)
      if (res.success) {
        setActivityCreated(true); // Trigger re-fetch via useEffect
      } else {
        throw new Error(res.message)
      }
    } catch (err: unknown) {
      let errorMessage = 'request failed';
      if (err instanceof Error) {
        errorMessage = err.message;
      } else if (typeof err === 'object' && err !== null && 'message' in err && typeof (err as { message: unknown }).message === 'string') {
        errorMessage = (err as { message: string }).message;
      }
      setActivityError(errorMessage);
    }
  }

  const handleChangeTicket = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const target = e.target as HTMLInputElement
    const { name, value } = target

    setTicketForm((f) => ({
      ...f,
      [name]: value,
    }))
  }

  const handleTicketSubmit = async (
    e: React.SyntheticEvent<HTMLFormElement>,
  ) => {
    e.preventDefault()
    setTicketError(null)
    if (!selectedEventId) {
      setTicketError('pick an event first')
      return
    }
    if (!selectedEditionId) {
      setTicketError('pick an edition first')
      return
    }

    try {
      const res = await createTicketFn(ticketForm, selectedEventId, selectedEditionId)
      if (res.success) {
        setTicketCreated(true); // Trigger re-fetch via useEffect
      } else {
        throw new Error(res.message)
      }
    } catch (err: unknown) {
      let errorMessage = 'request failed';
      if (err instanceof Error) {
        errorMessage = err.message;
      } else if (typeof err === 'object' && err !== null && 'message' in err && typeof (err as { message: unknown }).message === 'string') {
        errorMessage = (err as { message: string }).message;
      }
      setTicketError(errorMessage);
    }
  }

  const handleChangeProduct = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => {
    const target = e.target as HTMLInputElement | HTMLSelectElement
    const { name, value, type } = target
    const checked = type === 'checkbox' ? (target).checked : undefined

    setProductForm((f) => ({
      ...f,
      [name]: type === 'checkbox' ? checked : value,
    } as unknown as ProductCreateI))
  }

  const handleProductSubmit = async (
    e: React.SyntheticEvent<HTMLFormElement>,
  ) => {
    e.preventDefault()
    setProductError(null)
    if (!selectedEventId) {
      setProductError('pick an event first')
      return
    }
    if (!selectedEditionId) {
      setProductError('pick an edition first')
      return
    }

    try {
      const data = {
        ...productForm,
        edition_scope_id: selectedEditionId,
        price_cents: Number(productForm.price_cents),
        inventory_quantity: Number(productForm.inventory_quantity),
        available_from: productForm.available_from ? parseDatetimeLocal(productForm.available_from).toISOString() : undefined,
        available_until: productForm.available_until ? parseDatetimeLocal(productForm.available_until).toISOString() : undefined,
      }
      const res = await createProductFn(data, selectedEventId, selectedEditionId)
      if (res.success) {
        setProductCreated(true); // Trigger re-fetch via useEffect
      } else {
        throw new Error(res.message)
      }
    } catch (err: unknown) {
      let errorMessage = 'request failed';
      if (err instanceof Error) {
        errorMessage = err.message;
      } else if (typeof err === 'object' && err !== null && 'message' in err && typeof (err as { message: unknown }).message === 'string') {
        errorMessage = (err as { message: string }).message;
      }
      setProductError(errorMessage);
    }
  }

  const handleChangeCheckpoint = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const target = e.target as HTMLInputElement | HTMLSelectElement
    const { name, value } = target

    setCheckpointForm((f) => ({
      ...f,
      [name]: value,
    } as unknown as CheckpointCreateI))
  }

  const handleCheckpointSubmit = async (
    e: React.SyntheticEvent<HTMLFormElement>,
  ) => {
    e.preventDefault()
    setCheckpointError(null)
    if (!selectedEventId) {
      setCheckpointError('pick an event first')
      return
    }
    if (!selectedEditionId) {
      setCheckpointError('pick an edition first')
      return
    }

    try {
      const data = {
        ...checkpointForm,
        starts_at: checkpointForm.starts_at ? parseDatetimeLocal(checkpointForm.starts_at).toISOString() : undefined,
        ends_at: checkpointForm.ends_at ? parseDatetimeLocal(checkpointForm.ends_at).toISOString() : undefined,
      }
      const res = await createCheckpointFn(data, selectedEventId, selectedEditionId)
      if (res.success) {
        setCheckpointCreated(true); // Trigger re-fetch via useEffect
      } else {
        throw new Error(res.message)
      }
    } catch (err: unknown) {
      let errorMessage = 'request failed';
      if (err instanceof Error) {
        errorMessage = err.message;
      } else if (typeof err === 'object' && err !== null && 'message' in err && typeof (err as { message: unknown }).message === 'string') {
        errorMessage = (err as { message: string }).message;
      }
      setCheckpointError(errorMessage);
    }
  }

  return (
    <div className='flex flex-col items-center my-4 gap-2'>
      <h3>Você {isAuthenticated ? "já" : "não"} está autenticado</h3>
      <SignUp />
      <SignIn />

      {/* event list */}
      <div className="w-full max-w-md mt-8">
        <h4 className="text-lg font-semibold mb-2">Your events</h4>

        <ul className="list-disc pl-5">
          {events.map((ev) => (
            <li key={ev.id} className="mb-1 p-2 border rounded-md">
              <button
                className="text-blue-600 underline text-left"
                onClick={() => {
                  setSelectedEventId(ev.id)
                }}
              >
                <strong>ID:</strong> {ev.id}<br />
                <strong>Name:</strong> {ev.name} ({ev.slug}) {ev.acronym ? `[${ev.acronym}]` : ''}<br />
                <strong>Tagline:</strong> {ev.tagline ?? 'N/A'}<br />
                <strong>Description:</strong> {ev.description ?? 'N/A'}<br />
                <strong>Is Series:</strong> {ev.is_series ? 'Yes' : 'No'}<br />
                <strong>Editions Count:</strong> {ev.editions_count}<br />
                <strong>Status:</strong> {ev.status}<br />
                <strong>Contact Email:</strong> {ev.contact_email ?? 'N/A'}<br />
                <strong>Created By:</strong> {ev.created_by}<br />
                <strong>Created At:</strong> {new Date(ev.created_at).toLocaleDateString()} {new Date(ev.created_at).toLocaleTimeString()}<br />
                <strong>Updated At:</strong> {new Date(ev.updated_at).toLocaleDateString()} {new Date(ev.updated_at).toLocaleTimeString()}<br />
                <strong>Deleted At:</strong> {ev.deleted_at ? `${new Date(ev.deleted_at).toLocaleDateString()} ${new Date(ev.deleted_at).toLocaleTimeString()}` : 'N/A'}
              </button>
              <button
                className="mt-2 block bg-blue-600 text-white py-1 px-3 rounded disabled:bg-gray-400"
                disabled={ev.status !== 'draft'}
                onClick={async (e) => {
                  e.stopPropagation();
                  await publishEventFn(ev.id);
                  setEventCreated(true);
                }}
              >
                {ev.status === 'draft' ? 'Publish Event' : 'Published'}
              </button>
            </li>
          ))}
        </ul>
      </div>

      {/* event creation form */}
      <div className="w-full max-w-md mt-8">
        <h4 className="text-lg font-semibold mb-2">Create an event</h4>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <label>
            OrganizationID
            <input
              type="text"
              name="organization_id"
              value={form.organization_id ?? ''}
              onChange={handleChange}
              className="border p-1"
            />
          </label>

          <label>
            Name*
            <input
              type="text"
              name="name"
              value={form.name}
              onChange={handleChange}
              className="border p-1"
            />
          </label>

          <label>
            Acronym
            <input
              type="text"
              name="acronym"
              value={form.acronym ?? ''}
              onChange={handleChange}
              className="border p-1"
            />
          </label>

          <label>
            Slug*
            <input
              type="text"
              name="slug"
              value={form.slug}
              onChange={handleChange}
              className="border p-1"
            />
          </label>

          <label>
            Tagline
            <input
              type="text"
              name="tagline"
              value={form.tagline ?? ''}
              onChange={handleChange}
              className="border p-1"
            />
          </label>

          <label>
            Description
            <textarea
              name="description"
              value={form.description ?? ''}
              onChange={handleChange}
              className="border p-1"
            />
          </label>

          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              name="is_series"
              checked={form.is_series}
              onChange={handleChange}
            />
            Is series?
          </label>

          <label>
            Logo URL
            <input
              type="text"
              name="logo_url"
              value={form.logo_url ?? ''}
              onChange={handleChange}
              className="border p-1"
            />
          </label>

          <label>
            Banner URL
            <input
              type="text"
              name="banner_url"
              value={form.banner_url ?? ''}
              onChange={handleChange}
              className="border p-1"
            />
          </label>

          <label>
            Contact Email*
            <input
              type="email"
              name="contact_email"
              value={form.contact_email ?? ''}
              onChange={handleChange}
              className="border p-1"
            />
          </label>

          <button
            type="submit"
            className="bg-blue-600 text-white py-1 px-3 rounded"
          >
            Submit
          </button>
        </form>
        {result && (
          <pre className="mt-4 p-2 border bg-gray-100">
            {JSON.stringify(result, null, 2)}
          </pre>
        )}
        {error && <div className="text-red-600 mt-2">{error}</div>}
      </div>

      {/* edition creation/listing */}
      <div className="w-full max-w-md mt-8">
        <h4 className="text-lg font-semibold mb-2">
          Create an edition {selectedEvent ? `for ${selectedEvent.name}` : ''}
        </h4>
        {selectedEventId ? (
          <form onSubmit={handleEditionSubmit} className="flex flex-col gap-4">
            <label>
              Type*
              <select
                name="type"
                value={editionForm.type}
                onChange={handleEditionChange}
                className="border p-1"
              >
                <option value="year">year</option>
                <option value="season">season</option>
                <option value="number">number</option>
                <option value="ordinal">ordinal</option>
                <option value="custom">custom</option>
              </select>
            </label>
            <label>
              Edition Name* (at least 3)
              <input
                type="text"
                name="edition_name"
                value={editionForm.edition_name}
                onChange={handleEditionChange}
                className="border p-1"
              />
            </label>
            <label>
              Registration opens At
              <input
                type="datetime-local"
                name="registration_opens_at"
                value={editionForm.registration_opens_at ?? ''}
                onChange={handleEditionChange}
                className="border p-1"
              />
            </label>
            <label>
              Registration closes At
              <input
                type="datetime-local"
                name="registration_closes_at"
                value={editionForm.registration_closes_at ?? ''}
                onChange={handleEditionChange}
                className="border p-1"
              />
            </label>
            <label>
              Starts At*
              <input
                type="datetime-local"
                name="starts_at"
                value={editionForm.starts_at}
                onChange={handleEditionChange}
                className="border p-1"
              />
            </label>
            <label>
              Ends At*
              <input
                type="datetime-local"
                name="ends_at"
                value={editionForm.ends_at}
                onChange={handleEditionChange}
                className="border p-1"
              />
            </label>
            <label>
              Timezone*
              <select
                name="timezone"
                value={editionForm.timezone}
                onChange={handleEditionChange}
                className="border p-1"
              >
                {timezones.map((tz) => (
                  <option key={tz} value={tz}>
                    {tz}
                  </option>
                ))}
              </select>
            </label>
            <label>
              Location Name*
              <input
                type="text"
                name="location_name"
                value={editionForm.location_name}
                onChange={handleEditionChange}
                className="border p-1"
              />
            </label>
            <label>
              Location Address*
              <input
                type="text"
                name="location_address"
                value={editionForm.location_address}
                onChange={handleEditionChange}
                className="border p-1"
              />
            </label>
            <button
              type="submit"
              className="bg-green-600 text-white py-1 px-3 rounded"
            >
              Create edition
            </button>
          </form>
        ) : (
          <p className="text-sm text-gray-600">Select an event above to create editions</p>
        )}
        {editionError && <div className="text-red-600 mt-2">{editionError}</div>}

        {selectedEventId && (
          <div className="mt-4">
            <h5 className="font-semibold">Editions for selected event</h5>
            <ul className="list-disc pl-5">
              {editions.map((ed) => (
                <li key={ed.id} className="mb-2 p-2 border rounded-md">
                  <button
                    className="text-blue-600 underline text-left"
                    onClick={() => {
                      setSelectedEditionId(ed.id)
                    }}
                  >
                    <strong>ID:</strong> {ed.id}<br />
                    <strong>Edition Name:</strong> {ed.edition_name} ({ed.type})<br />
                    <strong>Tagline:</strong> {ed.tagline || 'N/A'}<br />
                    <strong>Description:</strong> {ed.description || 'N/A'}<br />
                    <strong>Status:</strong> {ed.status}<br />
                    <strong>Monetary Type:</strong> {ed.monetary_type}<br />
                    <strong>Registration Opens:</strong> {ed.registration_opens_at ? `${new Date(ed.registration_opens_at).toLocaleDateString()} ${new Date(ed.registration_opens_at).toLocaleTimeString()}` : 'N/A'}<br />
                    <strong>Registration Closes:</strong> {ed.registration_closes_at ? `${new Date(ed.registration_closes_at).toLocaleDateString()} ${new Date(ed.registration_closes_at).toLocaleTimeString()}` : 'N/A'}<br />
                    <strong>Starts:</strong> {new Date(ed.starts_at).toLocaleDateString()} {new Date(ed.starts_at).toLocaleTimeString()}<br />
                    <strong>Ends:</strong> {new Date(ed.ends_at).toLocaleDateString()} {new Date(ed.ends_at).toLocaleTimeString()}<br />
                    <strong>Timezone:</strong> {ed.timezone}<br />
                    <strong>Location Name:</strong> {ed.location_name}<br />
                    <strong>Location Address:</strong> {ed.location_address}<br />
                    <strong>Contact Email:</strong> {ed.contact_email || 'N/A'}<br />
                    <strong>Contact Phone:</strong> {ed.contact_phone || 'N/A'}<br />
                    <strong>Organizer Name:</strong> {ed.organizer_name || 'N/A'}<br />
                    <strong>Created By:</strong> {ed.created_by}<br />
                    <strong>Created At:</strong> {new Date(ed.created_at).toLocaleDateString()} {new Date(ed.created_at).toLocaleTimeString()}<br />
                    <strong>Updated At:</strong> {new Date(ed.updated_at).toLocaleDateString()} {new Date(ed.updated_at).toLocaleTimeString()}<br />
                    <strong>Deleted At:</strong> {ed.deleted_at ? `${new Date(ed.deleted_at).toLocaleDateString()} ${new Date(ed.deleted_at).toLocaleTimeString()}` : 'N/A'}
                  </button>
                  {ed.trie_payments_credential_id ? (
                    <button
                      onClick={() => handleDisconnect(
                        ed.event_id,
                        ed.id,
                        ed.trie_payments_credential_id
                      )}
                    >
                      Desconectar
                    </button>
                  ) : (
                    <button onClick={() => handleConnect(ed.event_id, ed.id)}>Conectar</button>
                  )}
                  <Link
                    to="/events/$eventId/editions/$editionId/products"
                    params={{ eventId: ed.event_id, editionId: ed.id }}
                  >
                    Ver Produtos
                  </Link>
                  <button
                    className="mt-2 block bg-green-600 text-white py-1 px-3 rounded disabled:bg-gray-400"
                    disabled={ed.status !== 'draft'}
                    onClick={async (e) => {
                      e.stopPropagation();
                      if (selectedEventId) {
                        await publishEditionFn(selectedEventId, ed.id);
                        setEditionCreated(true);
                      }
                    }}
                  >
                    {ed.status === 'draft' ? 'Publish Edition' : 'Published'}
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>

      {/* activity creation/listing */}
      <div className="w-full max-w-md mt-8">
        <h4 className="text-lg font-semibold mb-2">
          Create an activity {selectedEvent && selectedEditionId ? `for ${selectedEvent.name} edition ${selectedEditionId}` : ''}
        </h4>
        {selectedEventId && selectedEditionId ? (
          <form onSubmit={handleActivitySubmit} className="flex flex-col gap-4">
            <label>
              Title*
              <input
                type="text"
                name="title"
                value={activityForm.title}
                onChange={handleChangeActivity}
                className="border p-1"
              />
            </label>
            <label>
              Description
              <textarea
                name="description"
                value={activityForm.description ?? ''}
                onChange={handleChangeActivity}
                className="border p-1"
              />
            </label>
            <label>
              Location*
              <input
                type="text"
                name="location"
                value={activityForm.location}
                onChange={handleChangeActivity}
                className="border p-1"
              />
            </label>
            <label>
              Starts At*
              <input
                type="datetime-local"
                name="starts_at"
                value={activityForm.starts_at}
                onChange={handleChangeActivity}
                className="border p-1"
              />
            </label>
            <label>
              Ends At*
              <input
                type="datetime-local"
                name="ends_at"
                value={activityForm.ends_at}
                onChange={handleChangeActivity}
                className="border p-1"
              />
            </label>
            <label>
              Presenter Name
              <input
                type="text"
                name="presenter_name"
                value={activityForm.presenter_name ?? ''}
                onChange={handleChangeActivity}
                className="border p-1"
              />
            </label>
            <label>
              Token Cost*
              <input
                type="number"
                name="token_cost"
                value={activityForm.token_cost}
                onChange={handleChangeActivity}
                className="border p-1"
              />
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                name="has_capacity"
                checked={activityForm.has_capacity}
                onChange={handleChangeActivity}
              />
              Has Capacity?
            </label>
            {activityForm.has_capacity && (
              <label>
                Capacity*
                <input
                  type="number"
                  name="capacity"
                  value={activityForm.capacity}
                  onChange={handleChangeActivity}
                  className="border p-1"
                />
              </label>
            )}
            <label>
              Difficulty*
              <select
                name="difficulty"
                value={activityForm.difficulty}
                onChange={handleChangeActivity}
                className="border p-1"
              >
                <option value="no_prerequisites">No Prerequisites</option>
                <option value="beginner">Beginner</option>
                <option value="intermediate">Intermediate</option>
                <option value="advanced">Advanced</option>
                <option value="expert">Expert</option>
              </select>
            </label>
            <button
              type="submit"
              className="bg-purple-600 text-white py-1 px-3 rounded"
            >
              Create Activity
            </button>
          </form>
        ) : (
          <p className="text-sm text-gray-600">Select an event and an edition above to create activities</p>
        )}
        {activityError && <div className="text-red-600 mt-2">{activityError}</div>}

        {selectedEventId && selectedEditionId && (
          <div className="mt-4">
            <h5 className="font-semibold">Activities for selected edition</h5>
            <ul className="list-disc pl-5">
              {activities.map((act) => (
                <li key={act.id} className="mb-2 p-2 border rounded-md">
                  <strong>ID:</strong> {act.id}<br />
                  <strong>Title:</strong> {act.title}<br />
                  <strong>Description:</strong> {act.description || 'N/A'}<br />
                  <strong>Location:</strong> {act.location}<br />
                  <strong>Starts:</strong> {new Date(act.starts_at).toLocaleDateString()} {new Date(act.starts_at).toLocaleTimeString()}<br />
                  <strong>Ends:</strong> {new Date(act.ends_at).toLocaleDateString()} {new Date(act.ends_at).toLocaleTimeString()}<br />
                  <strong>Presenter:</strong> {act.presenter_name || 'N/A'}<br />
                  <strong>Token Cost:</strong> {act.token_cost}<br />
                  <strong>Has Capacity:</strong> {act.has_capacity ? 'Yes' : 'No'}<br />
                  <strong>Capacity:</strong> {act.capacity}<br />
                  <strong>Remaining Capacity:</strong> {act.remaining_capacity}<br />
                  <strong>Difficulty:</strong> {act.difficulty}<br />
                  <strong>Status:</strong> {act.status}<br />
                  <strong>Created By:</strong> {act.created_by}<br />
                  <strong>Created At:</strong> {new Date(act.created_at).toLocaleDateString()} {new Date(act.created_at).toLocaleTimeString()}<br />
                  <strong>Updated At:</strong> {new Date(act.updated_at).toLocaleDateString()} {new Date(act.updated_at).toLocaleTimeString()}<br />
                  <strong>Deleted At:</strong> {act.deleted_at ? `${new Date(act.deleted_at).toLocaleDateString()} ${new Date(act.deleted_at).toLocaleTimeString()}` : 'N/A'}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>

      {/* ticket creation/listing */}
      <div className="w-full max-w-md mt-8">
        <h4 className="text-lg font-semibold mb-2">
          Create a ticket {selectedEvent && selectedEditionId ? `for ${selectedEvent.name} edition ${selectedEditionId}` : ''}
        </h4>
        {selectedEventId && selectedEditionId ? (
          <form onSubmit={handleTicketSubmit} className="flex flex-col gap-4">
            <label>
              Name*
              <input
                type="text"
                name="name"
                value={ticketForm.name}
                onChange={handleChangeTicket}
                className="border p-1"
              />
            </label>
            <label>
              Description
              <textarea
                name="description"
                value={ticketForm.description ?? ''}
                onChange={handleChangeTicket}
                className="border p-1"
              />
            </label>
            <button
              type="submit"
              className="bg-blue-600 text-white py-1 px-3 rounded"
            >
              Create Ticket
            </button>
          </form>
        ) : (
          <p className="text-sm text-gray-600">Select an event and an edition above to create tickets</p>
        )}
        {ticketError && <div className="text-red-600 mt-2">{ticketError}</div>}

        {selectedEventId && selectedEditionId && (
          <div className="mt-4">
            <h5 className="font-semibold">Tickets for selected edition</h5>
            <ul className="list-disc pl-5">
              {tickets.map((tix) => (
                <li key={tix.id} className="mb-2 p-2 border rounded-md">
                  <strong>ID:</strong> {tix.id}<br />
                  <strong>Name:</strong> {tix.name}<br />
                  <strong>Description:</strong> {tix.description || 'N/A'}<br />
                  <strong>Created By:</strong> {tix.created_by}<br />
                  <strong>Created At:</strong> {new Date(tix.created_at).toLocaleDateString()} {new Date(tix.created_at).toLocaleTimeString()}<br />
                  <strong>Updated At:</strong> {new Date(tix.updated_at).toLocaleDateString()} {new Date(tix.updated_at).toLocaleTimeString()}<br />
                  <strong>Deleted At:</strong> {tix.deleted_at ? `${new Date(tix.deleted_at).toLocaleDateString()} ${new Date(tix.deleted_at).toLocaleTimeString()}` : 'N/A'}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>

      {/* product creation/listing */}
      <div className="w-full max-w-md mt-8">
        <h4 className="text-lg font-semibold mb-2">
          Create a product {selectedEvent && selectedEditionId ? `for ${selectedEvent.name} edition ${selectedEditionId}` : ''}
        </h4>
        {selectedEventId && selectedEditionId ? (
          <form onSubmit={handleProductSubmit} className="flex flex-col gap-4">
            <label>
              Name*
              <input
                type="text"
                name="name"
                value={productForm.name}
                onChange={handleChangeProduct}
                className="border p-1"
              />
            </label>
            <label>
              Description
              <textarea
                name="description"
                value={productForm.description ?? ''}
                onChange={handleChangeProduct}
                className="border p-1"
              />
            </label>
            <label>
              Type*
              <select
                name="type"
                value={productForm.type}
                onChange={handleChangeProduct}
                className="border p-1"
              >
                <option value="merchandise">Merchandise</option>
                <option value="ticket">Ticket</option>
                <option value="token">Token</option>
                <option value="bundle">Bundle</option>
              </select>
            </label>
            <label>
              Ticket ID (Optional - if product is a ticket)
              <input
                type="text"
                name="ticket_id"
                value={productForm.ticket_id ?? ''}
                onChange={handleChangeProduct}
                className="border p-1"
              />
            </label>
            <label>
              Price (cents)*
              <input
                type="number"
                name="price_cents"
                value={productForm.price_cents}
                onChange={handleChangeProduct}
                className="border p-1"
              />
            </label>
            <label>
              Available From
              <input
                type="datetime-local"
                name="available_from"
                value={productForm.available_from ?? ''}
                onChange={handleChangeProduct}
                className="border p-1"
              />
            </label>
            <label>
              Available Until
              <input
                type="datetime-local"
                name="available_until"
                value={productForm.available_until ?? ''}
                onChange={handleChangeProduct}
                className="border p-1"
              />
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                name="has_inventory"
                checked={productForm.has_inventory}
                onChange={handleChangeProduct}
              />
              Has Inventory?
            </label>
            {productForm.has_inventory && (
              <label>
                Inventory Quantity*
                <input
                  type="number"
                  name="inventory_quantity"
                  value={productForm.inventory_quantity}
                  onChange={handleChangeProduct}
                  className="border p-1"
                />
              </label>
            )}
            <button
              type="submit"
              className="bg-orange-600 text-white py-1 px-3 rounded"
            >
              Create Product
            </button>
          </form>
        ) : (
          <p className="text-sm text-gray-600">Select an event and an edition above to create products</p>
        )}
        {productError && <div className="text-red-600 mt-2">{productError}</div>}

        {selectedEventId && selectedEditionId && (
          <div className="mt-4">
            <h5 className="font-semibold">Products for selected edition</h5>
            <ul className="list-disc pl-5">
              {products.map((prod) => (
                <li key={prod.id} className="mb-2 p-2 border rounded-md">
                  <strong>ID:</strong> {prod.id}<br />
                  <strong>Name:</strong> {prod.name}<br />
                  <strong>Description:</strong> {prod.description || 'N/A'}<br />
                  <strong>Type:</strong> {prod.type}<br />
                  <strong>Ticket ID:</strong> {prod.ticket_id || 'N/A'}<br />
                  <strong>Price (cents):</strong> {prod.price_cents}<br />
                  <strong>Status:</strong> {prod.status}<br />
                  <strong>Available From:</strong> {prod.available_from ? `${new Date(prod.available_from).toLocaleDateString()} ${new Date(prod.available_from).toLocaleTimeString()}` : 'N/A'}<br />
                  <strong>Available Until:</strong> {prod.available_until ? `${new Date(prod.available_until).toLocaleDateString()} ${new Date(prod.available_until).toLocaleTimeString()}` : 'N/A'}<br />
                  <strong>Has Inventory:</strong> {prod.has_inventory ? 'Yes' : 'No'}<br />
                  <strong>Inventory Quantity:</strong> {prod.inventory_quantity}<br />
                  <strong>Inventory Remaining:</strong> {prod.inventory_remaining}<br />
                  <strong>Created By:</strong> {prod.created_by}<br />
                  <strong>Created At:</strong> {new Date(prod.created_at).toLocaleDateString()} {new Date(prod.created_at).toLocaleTimeString()}<br />
                  <strong>Updated At:</strong> {new Date(prod.updated_at).toLocaleDateString()} {new Date(prod.updated_at).toLocaleTimeString()}<br />
                  <strong>Deleted At:</strong> {prod.deleted_at ? `${new Date(prod.deleted_at).toLocaleDateString()} ${new Date(prod.deleted_at).toLocaleTimeString()}` : 'N/A'}
                  <button
                    className="mt-2 block bg-orange-600 text-white py-1 px-3 rounded disabled:bg-gray-400"
                    disabled={prod.status !== 'draft'}
                    onClick={async (e) => {
                      e.stopPropagation();
                      if (selectedEventId && selectedEditionId) {
                        await publishProductFn(selectedEventId, selectedEditionId, prod.id);
                        setProductCreated(true);
                      }
                    }}
                  >
                    {prod.status === 'draft' ? 'Publish Product' : 'Published'}
                  </button>
                </li>))}
            </ul>
          </div>
        )}
      </div>

      {/* checkpoint creation/listing */}
      <div className="w-full max-w-md mt-8">
        <h4 className="text-lg font-semibold mb-2">
          Create a checkpoint {selectedEvent && selectedEditionId ? `for ${selectedEvent.name} edition ${selectedEditionId}` : ''}
        </h4>
        {selectedEventId && selectedEditionId ? (
          <form onSubmit={handleCheckpointSubmit} className="flex flex-col gap-4">
            <label>
              Name*
              <input
                type="text"
                name="name"
                value={checkpointForm.name}
                onChange={handleChangeCheckpoint}
                className="border p-1"
              />
            </label>
            <label>
              Access Mode*
              <select
                name="access_mode"
                value={checkpointForm.access_mode}
                onChange={handleChangeCheckpoint}
                className="border p-1"
              >
                <option value="open">Open</option>
                <option value="ticket">Ticket</option>
                <option value="staff_only">Staff Only</option>
              </select>
            </label>
            <label>
              Type*
              <select
                name="type"
                value={checkpointForm.type}
                onChange={handleChangeCheckpoint}
                className="border p-1"
              >
                <option value="entry">Entry</option>
                <option value="zone">Zone</option>
                <option value="amenity">Amenity</option>
                <option value="session">Session</option>
                <option value="exit">Exit</option>
              </select>
            </label>
            <label>
              Starts At
              <input
                type="datetime-local"
                name="starts_at"
                value={checkpointForm.starts_at ?? ''}
                onChange={handleChangeCheckpoint}
                className="border p-1"
              />
            </label>
            <label>
              Ends At
              <input
                type="datetime-local"
                name="ends_at"
                value={checkpointForm.ends_at ?? ''}
                onChange={handleChangeCheckpoint}
                className="border p-1"
              />
            </label>
            <button
              type="submit"
              className="bg-red-600 text-white py-1 px-3 rounded"
            >
              Create Checkpoint
            </button>
          </form>
        ) : (
          <p className="text-sm text-gray-600">Select an event and an edition above to create checkpoints</p>
        )}
        {checkpointError && <div className="text-red-600 mt-2">{checkpointError}</div>}

        {selectedEventId && selectedEditionId && (
          <div className="mt-4">
            <h5 className="font-semibold">Checkpoints for selected edition</h5>
            <ul className="list-disc pl-5">
              {checkpoints.map((chkpt) => (
                <li key={chkpt.id} className="mb-2 p-2 border rounded-md">
                  <strong>ID:</strong> {chkpt.id}<br />
                  <strong>Name:</strong> {chkpt.name}<br />
                  <strong>Access Mode:</strong> {chkpt.access_mode}<br />
                  <strong>Type:</strong> {chkpt.type}<br />
                  <strong>Starts:</strong> {chkpt.starts_at ? `${new Date(chkpt.starts_at).toLocaleDateString()} ${new Date(chkpt.starts_at).toLocaleTimeString()}` : 'N/A'}<br />
                  <strong>Ends:</strong> {chkpt.ends_at ? `${new Date(chkpt.ends_at).toLocaleDateString()} ${new Date(chkpt.ends_at).toLocaleTimeString()}` : 'N/A'}<br />
                  <strong>Created By:</strong> {chkpt.created_by}<br />
                  <strong>Created At:</strong> {new Date(chkpt.created_at).toLocaleDateString()} {new Date(chkpt.created_at).toLocaleTimeString()}<br />
                  <strong>Updated At:</strong> {new Date(chkpt.updated_at).toLocaleDateString()} {new Date(chkpt.updated_at).toLocaleTimeString()}<br />
                  <strong>Deleted At:</strong> {chkpt.deleted_at ? `${new Date(chkpt.deleted_at).toLocaleDateString()} ${new Date(chkpt.deleted_at).toLocaleTimeString()}` : 'N/A'}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  )
}
