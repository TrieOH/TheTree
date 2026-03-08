import { createFileRoute } from '@tanstack/react-router'
import { SignIn, SignUp, useAuth } from '@trieoh/node-auth-sdk/react'
import React, { useState } from 'react'
import type { EventCreateI, EventI } from '@/features/events/model'
import type { EditionCreateI, EditionI } from '@/features/editions/model'
import { createEventFn, getOwnEventsFn } from '@/features/events/api'
import { createEditionFn, getAllEditionsFn } from '@/features/editions/api'
import {
  formatDateForDatetimeLocal,
  parseDatetimeLocal,
} from '@/shared/lib/date'

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

  React.useEffect(() => {
    const fetchEditions = async () => {
      if (selectedEventId) {
        const eds = await getAllEditionsFn(selectedEventId);
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
            Contact Email
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
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  )
}
