import { Copy, QrCode } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import { buttonVariants } from "@/shared/ui/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/shared/ui/shadcn/dialog";
import { ProfileQrCode } from "./profile-qr-code";

interface ProfileShareDialogProps {
  profileUrl: string;
  className?: string;
}

export function ProfileShareDialog({
  profileUrl,
  className,
}: ProfileShareDialogProps) {
  return (
    <Dialog>
      <DialogTrigger
        className={cn(
          "inline-flex size-11 items-center justify-center rounded-full",
          "border border-white/40 bg-background/90 text-foreground shadow-lg",
          "backdrop-blur-sm transition-transform active:scale-95",
          className,
        )}
        aria-label="Abrir QR Code do perfil"
      >
        <QrCode className="size-5" />
      </DialogTrigger>
      <DialogContent className="max-w-[calc(100%-2rem)] shadow-2xl sm:max-w-sm">
        <DialogHeader className="pr-8 text-center">
          <DialogTitle>Compartilhar perfil</DialogTitle>
          <DialogDescription>
            Escaneie o QR Code para abrir este perfil.
          </DialogDescription>
        </DialogHeader>
        <div className="flex justify-center rounded-xl bg-white p-4 shadow-inner">
          <ProfileQrCode value={profileUrl} size={220} />
        </div>
        <button
          type="button"
          onClick={() => navigator.clipboard.writeText(profileUrl)}
          className={buttonVariants({
            variant: "outline",
            className: "h-10 w-full rounded-md shadow-sm",
          })}
        >
          <Copy className="mr-2 size-4" />
          Copiar link do perfil
        </button>
      </DialogContent>
    </Dialog>
  );
}
