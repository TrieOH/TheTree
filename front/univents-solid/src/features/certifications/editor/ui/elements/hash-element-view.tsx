import type { JSX } from "@solidjs/web";
import ExternalLinkIcon from "~icons/lucide/external-link";
import type { HashCertificateElement } from "../../types";

const ExternalLink = ExternalLinkIcon as unknown as (props: { class?: string }) => JSX.Element;

interface HashElementViewProps {
  element: HashCertificateElement;
}

export function HashElementView(props: HashElementViewProps): JSX.Element {
  const hashLabel = () => props.element.hashLabel.trim() || "Código de verificação";
  const hash = () => props.element.hash.trim() || "{{cert_hash}}";
  const url = () => props.element.url.trim() || "{{verify_url}}";
  const linkLabel = () => props.element.linkLabel.trim() || url();

  return (
    <div
      class="flex h-full w-full flex-col justify-center gap-1 overflow-hidden px-1"
      style={{
        "text-align": props.element.align,
        color: props.element.color,
        "font-size": `${props.element.fontSize}px`,
      }}
    >
      <div class="truncate">
        <span class="opacity-80">{hashLabel()}: </span>
        <span class="font-semibold tracking-wide">{hash()}</span>
      </div>
      <a
        href={url()}
        onClick={(event) => event.preventDefault()}
        title={`Este link abre: ${url()}`}
        class="inline-flex items-center gap-1 truncate font-medium underline decoration-current underline-offset-2"
        style={{
          "justify-content":
            props.element.align === "center"
              ? "center"
              : props.element.align === "right"
                ? "flex-end"
                : "flex-start",
        }}
      >
        <span>{linkLabel()}</span>
        <ExternalLink class="size-3 shrink-0" />
      </a>
    </div>
  );
}
