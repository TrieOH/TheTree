import type { JSX } from "@solidjs/web";
import { cn } from "@trieoh/ui-solid";
import Share2Icon from "~icons/lucide/share-2";
import { handleShare } from "@/shared/lib/share";

const Share2 = Share2Icon as unknown as (p: { class?: string }) => JSX.Element;

export function ProfileShareButton(props: {
	name: string;
	url: string;
	class?: string;
}) {
	return (
		<button
			type="button"
			class={cn(
				"inline-flex items-center justify-center transition-all active:scale-95",
				props.class ??
					"size-10 rounded-md border border-border bg-background text-foreground shadow-sm hover:bg-accent",
			)}
			onClick={() => void handleShare(props.name, props.url)}
			aria-label="Compartilhar perfil"
		>
			<Share2 class="size-4" />
		</button>
	);
}
