import { toast } from "@/shared/ui/toast";

export async function handleShare(title: string, url = window.location.href) {
  try {
    if (typeof navigator.share === "function") {
      await navigator.share({ title, url });
      return;
    }
    await navigator.clipboard.writeText(url);
    toast.success("Link copiado!");
  } catch {
    toast.error("Erro ao compartilhar");
  }
}
