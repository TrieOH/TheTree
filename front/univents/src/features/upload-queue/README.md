# Persistent upload queue

This feature owns background image processing independently from any CRUD screen.
Tasks and their `Blob`s are stored in IndexedDB, retried with exponential backoff,
and resumed when the application starts again.

The root provider scopes every task to `auth.profile().id`. The active account is
applied inside the store, so callers cannot accidentally read or process another
account's queue after logout/login on the same browser. Legacy tasks without an
account identifier stay quarantined and are never shown or processed.

## Integrating a feature

Register an association handler once from that feature's application setup:

```ts
const unregister = registerUploadAssociationHandler(
  'event-gallery',
  async (task, uploadedUrl) => {
    const response = await addImage(task.owner.id, uploadedUrl)
    if (!response.success) {
      throw new UploadAssociationError(response.message, {
        status: response.status,
      })
    }
  },
)
```

After the owning record exists, enqueue each file independently:

```ts
await enqueue({
  file,
  owner: { type: 'event', id: event.id, label: event.name },
  mediaType: 'gallery',
  storagePath: `events/${event.id}/gallery`,
  correctionPath: `/admin/events/${event.id}/media`,
  association: {
    handlerKey: 'event-gallery',
    input: { position: 0 },
  },
})
```

Do not store callbacks in tasks: IndexedDB cannot persist them. Persist only the
handler key and serializable input, then register the matching handler at startup.

Moderation and validation errors require file replacement. Network, rate-limit,
timeout, and server errors use automatic retry. Association handlers should throw
`UploadAssociationError` with an HTTP status so the queue can make the correct
decision without parsing messages.
