import { Loader2, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/shared/ui/shadcn/alert-dialog";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";
import type { EventMemberI } from "../api/members";

interface RemoveEventMemberModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  member: EventMemberI | null;
  onRemove: (userId: string, email: string) => Promise<boolean>;
}

export function RemoveEventMemberModal({
  open,
  onOpenChange,
  member,
  onRemove,
}: RemoveEventMemberModalProps) {
  const [email, setEmail] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (open) setEmail("");
  }, [open]);

  const handleRemove = async () => {
    const normalizedEmail = email.trim().toLowerCase();
    if (!member || !normalizedEmail) return;

    setIsSubmitting(true);
    const didRemove = await onRemove(member.user_id, normalizedEmail);
    setIsSubmitting(false);

    if (didRemove) onOpenChange(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Remover membro?</AlertDialogTitle>
          <AlertDialogDescription>
            Confirme o e-mail do usuário{" "}
            <span className="font-mono">{member?.user_id}</span> para removê-lo
            deste evento.
          </AlertDialogDescription>
        </AlertDialogHeader>

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
            disabled={!member || !email.trim() || isSubmitting}
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
  );
}
