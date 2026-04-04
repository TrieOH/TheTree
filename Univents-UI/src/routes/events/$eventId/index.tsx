import { createFileRoute } from '@tanstack/react-router'
import z from 'zod'

const searchSchema = z.object({
  edit: z.boolean().default(false),
})

export const Route = createFileRoute('/events/$eventId/')({
  validateSearch: searchSchema,
})
