import { useEffect, useMemo, useState } from 'react'
import { Loader2, Trash2 } from 'lucide-react'
import type { EventMemberI } from '../api/members'
import type { EventMemberRole } from '../model/member'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/shared/ui/shadcn/alert-dialog'
import { Button } from '@/shared/ui/shadcn/button'
import { Input } from '@/shared/ui/shadcn/input'
import { Label } from '@/shared/ui/shadcn/label'

const roleLabels: Record<EventMemberRole, string> = {
  owner: 'Proprietário',
  admin: 'Administrador',
  staff: 'Equipe',
}

interface RemoveEventMemberModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  members: EventMemberI[]
  onRemove: (userId: string, email: string) => Promise<boolean>
}

export function RemoveEventMemberModal({
  open,
  onOpenChange,
  members,
  onRemove,
}: RemoveEventMemberModalProps) {
  const [selectedUserId, setSelectedUserId] = useState('')
  const [email, setEmail] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const selectedMember = useMemo(
    () => members.find((member) => member.user_id === selectedUserId),
    [members, selectedUserId],
  )

  useEffect(() => {
    if (!open) return

    const initialMember = members[0]
    setSelectedUserId(initialMember?.user_id ?? '')
    setEmail(initialMember?.email ?? '')
  }, [members, open])

  const handleMemberChange = (userId: string) => {
    const member = members.find((candidate) => candidate.user_id === userId)
    setSelectedUserId(userId)
    setEmail(member?.email ?? '')
  }

  const handleRemove = async () => {
    const normalizedEmail = email.trim().toLowerCase()
    if (!selectedMember || !normalizedEmail) return

    setIsSubmitting(true)
    const didRemove = await onRemove(selectedMember.user_id, normalizedEmail)
    setIsSubmitting(false)

    if (didRemove) onOpenChange(false)
  }

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Remover membro?</AlertDialogTitle>
          <AlertDialogDescription>
            Selecione o membro e confirme o e-mail para removê-lo deste evento.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <div className="space-y-2">
          <Label htmlFor="remove-member">Membro</Label>
          <select
            id="remove-member"
            value={selectedUserId}
            onChange={(event) => handleMemberChange(event.target.value)}
            className="h-10 w-full rounded-lg border border-input bg-background px-3 text-sm outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/50"
          >
            {members.map((member) => (
              <option key={member.id} value={member.user_id}>
                {member.email ?? `Usuário ${member.user_id.slice(0, 8)}`} —{' '}
                {roleLabels[member.role]}
              </option>
            ))}
          </select>
        </div>

        <div className="space-y-2">
          <Label htmlFor="remove-member-email">E-mail do membro</Label>
          <Input
            id="remove-member-email"
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            placeholder="membro@exemplo.com"
            autoComplete="email"
            required
          />
        </div>

        <AlertDialogFooter>
          <AlertDialogCancel disabled={isSubmitting}>
            Cancelar
          </AlertDialogCancel>
          <Button
            type="button"
            variant="destructive"
            onClick={handleRemove}
            disabled={!selectedMember || !email.trim() || isSubmitting}
          >
            {isSubmitting ? (
              <Loader2 className="size-4 animate-spin" />
            ) : (
              <Trash2 className="size-4" />
            )}
            Remover
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
