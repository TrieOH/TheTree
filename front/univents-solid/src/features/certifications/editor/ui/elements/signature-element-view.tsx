import type { JSX } from "@solidjs/web";
import type { SignatureCertificateElement } from "../../types";

interface SignatureElementViewProps {
  element: SignatureCertificateElement;
}

export function SignatureElementView(props: SignatureElementViewProps): JSX.Element {
  return (
    <img
      src={props.element.src}
      alt={`Assinatura de ${props.element.name}`}
      title={props.element.name}
      draggable={false}
      class="h-full w-full"
      style={{
        opacity: props.element.opacity,
        "object-fit": props.element.fit,
        "border-radius": `${props.element.radius}px`,
      }}
    />
  );
}
