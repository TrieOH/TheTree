import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { StaticText } from "@/features/editor";
import type { CertificationTemplateElement } from "../../../model";
import { HashElementView } from "./hash-element-view";
import { SignatureElementView } from "./signature-element-view";
import type {
  HashCertificateElement,
  ImageCertificateElement,
  SignatureCertificateElement,
  TextCertificateElement,
} from "../../types";

interface CertificateElementViewProps {
  element: CertificationTemplateElement;
}

export function CertificateElementView(
  props: CertificateElementViewProps,
): JSX.Element {
  return (
    <Show
      when={props.element.type === "text"}
      fallback={
        <Show
          when={props.element.type === "image"}
          fallback={
            <Show
              when={props.element.type === "signature"}
              fallback={
                <HashElementView
                  element={props.element as HashCertificateElement}
                />
              }
            >
              <SignatureElementView
                element={props.element as SignatureCertificateElement}
              />
            </Show>
          }
        >
          {(() => {
            const img = () => props.element as ImageCertificateElement;
            return (
              <img
                src={img().src}
                alt=""
                class="h-full w-full pointer-events-none select-none"
                style={{
                  "object-fit": img().fit,
                  opacity: img().opacity,
                  "border-radius": `${img().radius}px`,
                }}
                draggable={false}
              />
            );
          })()}
        </Show>
      }
    >
      <StaticText
        paragraphs={(props.element as TextCertificateElement).paragraphs}
        showVariables={false}
      />
    </Show>
  );
}
