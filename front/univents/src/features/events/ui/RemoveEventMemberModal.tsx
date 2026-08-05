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
import type { EventMemberWithEmailI } from "../api/members";

interface RemoveEventMemberModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  member: EventMemberWithEmailI | null;
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
    if (
      !member?.email ||
      !normalizedEmail ||
      normalizedEmail !== member.email.trim().toLowerCase()
    )
      return;

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
            Digite novamente o e-mail abaixo para confirmar a remoção deste
            membro.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <div className="space-y-2">
          <Label htmlFor="remove-member-email">E-mail do membro</Label>
          <Input
            id="remove-member-email"
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            placeholder={member?.email ?? "membro@exemplo.com"}
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
            disabled={
              !member?.email ||
              email.trim().toLowerCase() !==
                member.email.trim().toLowerCase() ||
              isSubmitting
            }
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
