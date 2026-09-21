import type { JSX } from "@solidjs/web";
import type { ImageCertificateElement } from "../../types";

interface ImageElementViewProps {
  element: ImageCertificateElement;
}

export function ImageElementView(props: ImageElementViewProps): JSX.Element {
  return (
    <div
      class="h-full w-full overflow-hidden"
      style={{
        "border-radius": `${props.element.radius}px`,
        opacity: props.element.opacity,
      }}
    >
      <img
        src={props.element.src}
        alt=""
        draggable={false}
        class="h-full w-full"
        style={{ "object-fit": props.element.fit }}
      />
    </div>
  );
}
