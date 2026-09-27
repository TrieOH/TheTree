import type { JSX } from "@solidjs/web";
import { render, screen, fireEvent } from "@solidjs/testing-library";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { CertificateEditor } from "@/features/certifications/editor/CertificateEditor";
import { CertificateCanvas } from "@/features/certifications/editor/ui/certificate-canvas";
import {
  certificateEditorActions,
  certificateEditorStore,
} from "@/features/certifications/editor/store";
import { DEFAULT_CERTIFICATION_TEMPLATE } from "@/features/certifications/default-template";

vi.mock("@tanstack/solid-router", () => ({
  Link: (props: { href: string; children: JSX.Element }) => (
    <a href={props.href}>{props.children}</a>
  ),
  useNavigate: () => vi.fn(),
}));

vi.mock("@trieoh/front-core-solid", () => ({
  useQuery: (optionsFn: (() => { enabled?: boolean }) | { enabled?: boolean }) => () => {
    const opts = typeof optionsFn === "function" ? optionsFn() : optionsFn;
    if (opts?.enabled === false) {
      return { data: undefined, isLoading: false, isError: false };
    }
    return {
      data: [],
      isLoading: false,
      isError: false,
    };
  },
}));

vi.mock("@/features/certifications/api/mutations", () => ({
  useCreateCertificationTemplateMutation: () => ({
    mutateAsync: vi.fn(),
    result: () => ({ status: "idle" }),
  }),
  useUpdateCertificationTemplateMutation: () => ({
    mutateAsync: vi.fn(),
    result: () => ({ status: "idle" }),
  }),
}));

vi.mock("@/shared/ui/toast", () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

describe("CertificateEditor", () => {
  beforeEach(() => {
    certificateEditorActions.loadDraft(DEFAULT_CERTIFICATION_TEMPLATE);
  });

  it("renders editor sidebar and canvas properly", () => {
    const { container } = render(() => (
      <CertificateEditor eventId="ev-123" editionId="ed-456" />
    ));

    expect(screen.getByText("Editor de certificados")).toBeDefined();
    expect(screen.getByPlaceholderText("Nome do certificado")).toBeDefined();
    expect(screen.getByText("Salvar")).toBeDefined();
    expect(container.querySelector("[data-editor-element]")).toBeDefined();
  });

  it("does not render dynamic variables dropdown in the top toolbar (matches Badge Editor)", () => {
    const { container } = render(() => (
      <CertificateEditor eventId="ev-123" editionId="ed-456" />
    ));

    const topbar = container.querySelector("[data-text-toolbar]");
    expect(topbar).not.toBeNull();
    expect(topbar?.querySelector('button[title*="variável"]')).toBeNull();
  });

  it("enters edit mode on double click and keeps editable element mounted", async () => {
    const { container } = render(() => <CertificateCanvas />);

    const textEl = container.querySelector(".cursor-text");
    expect(textEl).not.toBeNull();
    if (textEl) {
      fireEvent.dblClick(textEl);
      await new Promise((resolve) => setTimeout(resolve, 50));
      const editable = container.querySelector("[contenteditable]");
      expect(editable).not.toBeNull();
      expect(certificateEditorStore.editingElementId()).not.toBeNull();
      expect(certificateEditorStore.richTextController()).not.toBeNull();
    }
  });

  it("stays in edit mode when user inputs text and synchronizes paragraphs", async () => {
    const { container } = render(() => <CertificateCanvas />);

    const textEl = container.querySelector(".cursor-text");
    expect(textEl).not.toBeNull();
    if (textEl) {
      fireEvent.dblClick(textEl);
      await new Promise((resolve) => setTimeout(resolve, 50));
      const editable = container.querySelector("[contenteditable]");
      expect(editable).not.toBeNull();

      if (editable) {
        editable.innerHTML = "<p>Texto editado do certificado</p>";
        fireEvent.input(editable);
        await new Promise((resolve) => setTimeout(resolve, 50));

        // Must still be in edit mode!
        const editableAfterInput = container.querySelector("[contenteditable]");
        expect(editableAfterInput).not.toBeNull();

        // Check draft was updated
        const currentElements = certificateEditorStore.draft().design_data.elements;
        const textElement = currentElements.find((e) => e.type === "text");
        expect(textElement).toBeDefined();
        if (textElement && textElement.type === "text") {
          expect(textElement.paragraphs[0]?.runs[0]?.text).toContain("Texto editado do certificado");
        }
      }
    }
  });

  it("allows formatting bold via controller or fallback", async () => {
    const { container } = render(() => (
      <CertificateEditor eventId="ev-123" editionId="ed-456" />
    ));

    await new Promise((resolve) => setTimeout(resolve, 50));

    // Select text element
    const firstText = certificateEditorStore
      .draft()
      .design_data.elements.find((e) => e.type === "text");
    expect(firstText).toBeDefined();

    if (firstText) {
      certificateEditorActions.selectElement(firstText.id);
      await new Promise((resolve) => setTimeout(resolve, 50));

      // Topbar bold button should be enabled
      const boldButton = container.querySelector('button[title*="Negrito"]') as HTMLButtonElement | null;
      expect(boldButton).not.toBeNull();
      expect(boldButton?.disabled).toBe(false);

      if (boldButton) {
        fireEvent.click(boldButton);
        await new Promise((resolve) => setTimeout(resolve, 50));
        const updated = certificateEditorStore
          .draft()
          .design_data.elements.find((e) => e.id === firstText.id);
        if (updated && updated.type === "text") {
          expect(updated.paragraphs[0]?.runs.some((r) => r.bold)).toBe(true);
        }
      }
    }
  });
});
