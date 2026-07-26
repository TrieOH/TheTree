import type { ImageCertificateElement } from "../../types";

interface ImageElementViewProps {
  element: ImageCertificateElement;
}

export function ImageElementView({ element }: ImageElementViewProps) {
  return (
    <div
      className="h-full w-full overflow-hidden"
      style={{ borderRadius: element.radius, opacity: element.opacity }}
    >
      <img
        src={element.src}
        alt=""
        draggable={false}
        className="h-full w-full"
        style={{ objectFit: element.fit }}
      />
    </div>
  );
}
