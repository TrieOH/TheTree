import type { SignatureCertificateElement } from "../../types";

interface SignatureElementViewProps {
  element: SignatureCertificateElement;
}

export function SignatureElementView({ element }: SignatureElementViewProps) {
  return (
    <img
      src={element.src}
      alt={`Assinatura de ${element.name}`}
      title={element.name}
      draggable={false}
      className="h-full w-full"
      style={{
        opacity: element.opacity,
        objectFit: element.fit,
        borderRadius: element.radius,
      }}
    />
  );
}
