import z from 'zod'

export const eventMemberRoles = ['owner', 'admin', 'staff'] as const

export type EventMemberRole = (typeof eventMemberRoles)[number]

export const eventMemberCreateSchema = z.object({
  email: z.email('Informe um e-mail válido'),
  role: z.enum(eventMemberRoles),
})

export type EventMemberCreateInput = z.input<typeof eventMemberCreateSchema>
export type EventMemberCreateOutput = z.output<typeof eventMemberCreateSchema>
