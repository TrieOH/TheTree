import { Share2 } from "lucide-react";
import { handleShare } from "@/shared/lib/share";
import { cn } from "@/shared/lib/utils";

interface ProfileShareDialogProps {
  profileUrl: string;
  className?: string;
}

export function ProfileShareDialog({
  profileUrl,
  className,
}: ProfileShareDialogProps) {
  return (
    <button
      type="button"
      onClick={() => handleShare("Perfil", profileUrl)}
      className={cn(
        "inline-flex size-11 items-center justify-center rounded-full",
        "border border-border/40 bg-background/90 text-foreground shadow-lg",
        "backdrop-blur-sm transition-transform active:scale-95",
        className,
      )}
      aria-label="Compartilhar perfil"
    >
      <Share2 className="size-5" />
    </button>
  );
}
