import type React from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogClose,
} from "@/shared/ui/shadcn/dialog";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";

interface PublishConfirmDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  schemaTitle: string;
}

export const PublishConfirmDialog: React.FC<PublishConfirmDialogProps> = ({
  isOpen,
  onClose,
  onConfirm,
  schemaTitle,
}) => {
  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Confirm Publication</DialogTitle>
          <DialogDescription>
            You are about to publish the schema "{schemaTitle}". This action will make the schema live and accessible.
            Are you sure you want to proceed? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <ShadowButton variant="outline" onClick={onClose} value="Cancel" />
          </DialogClose>
          <ShadowButton variant="solid" onClick={onConfirm} value="Publish" />
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
