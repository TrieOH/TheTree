import { toast } from "@/shared/ui/toast";

/**
 * Attempts to share via the Web Share API (`navigator.share`).
 * Falls back to copying the URL to the clipboard.
 * @param title - Title to share (e.g. event name)
 * @param url - URL to share (defaults to `window.location.href`)
 */
export async function handleShare(title: string, url?: string) {
  const shareUrl = url ?? window.location.href;
  try {
    if (typeof navigator.share === "function") {
      await navigator.share({ title, url: shareUrl });
      return;
    }
    await navigator.clipboard.writeText(shareUrl);
    toast.success("Link copiado!");
  } catch {
    toast.error("Erro ao compartilhar");
  }
}
