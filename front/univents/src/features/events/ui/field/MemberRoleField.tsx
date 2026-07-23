import { ShieldCheck, UserCog, Users } from 'lucide-react'
import type { FieldFormApi } from '@/widgets/multi-step-form/model/types'
import type {
  EventMemberCreateInput,
  EventMemberRole,
} from '../../model/member'
import { getFieldError } from '@/widgets/multi-step-form/utils/get-field-error'
import { cn } from '@/shared/lib/utils'

const roleOptions: Array<{
  value: EventMemberRole
  label: string
  description: string
  icon: typeof Users
}> = [
  {
    value: 'staff',
    label: 'Equipe',
    description: 'Acesso às operações do dia a dia.',
    icon: Users,
  },
  {
    value: 'admin',
    label: 'Administrador',
    description: 'Pode gerenciar o evento e sua equipe.',
    icon: UserCog,
  },
  {
    value: 'owner',
    label: 'Proprietário',
    description: 'Nível máximo de acesso ao evento.',
    icon: ShieldCheck,
  },
]

export function MemberRoleField({
  form,
}: {
  form: FieldFormApi<EventMemberCreateInput>
}) {
  const selectedRole = form.watch('role')
  const error = getFieldError(form.formState.errors, 'role')

  return (
    <fieldset className="space-y-2">
      <legend className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
        Função
      </legend>

      <div className="grid gap-2 sm:grid-cols-3">
        {roleOptions.map((option) => {
          const Icon = option.icon
          const selected = selectedRole === option.value

          return (
            <button
              key={option.value}
              type="button"
              aria-pressed={selected}
              onClick={() =>
                form.setValue('role', option.value, {
                  shouldDirty: true,
                  shouldTouch: true,
                  shouldValidate: true,
                })
              }
              className={cn(
                'flex min-h-28 flex-col items-start gap-2 rounded-xl border p-3 text-left transition-colors',
                selected
                  ? 'border-primary bg-primary/10 text-primary ring-2 ring-primary/15'
                  : 'border-border bg-background hover:border-primary/40 hover:bg-muted/40',
              )}
            >
              <Icon className="size-5" />
              <span className="text-sm font-semibold">{option.label}</span>
              <span
                className={cn(
                  'text-xs leading-relaxed',
                  selected ? 'text-primary/80' : 'text-muted-foreground',
                )}
              >
                {option.description}
              </span>
            </button>
          )
        })}
      </div>

      {error ? <p className="text-xs text-destructive">{error}</p> : null}
    </fieldset>
  )
}
