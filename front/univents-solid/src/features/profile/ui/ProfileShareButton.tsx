import type { JSX } from "@solidjs/web";
import Share2Icon from "~icons/lucide/share-2";
import { handleShare } from "@/shared/lib/share";

const Share2 = Share2Icon as unknown as () => JSX.Element;

export function ProfileShareButton(props: {
	name: string;
	url: string;
	class?: string;
}) {
	return (
		<button
			type="button"
			class={`inline-flex size-10 items-center justify-center rounded-md border border-border bg-background shadow-sm ${props.class ?? ""}`}
			onClick={() => void handleShare(props.name, props.url)}
			aria-label="Compartilhar perfil"
		>
			<Share2 />
		</button>
	);
}
