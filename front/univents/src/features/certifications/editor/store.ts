import { Store, useStore } from '@tanstack/react-store'
import type { CertificationTemplateElement } from '../model'
import { DEFAULT_CERTIFICATE_CANVAS } from './constants'
import { createHashElement } from './factories'
import type { CertificateCanvasSize, CertificationTemplateDraft } from './types'

export interface AvailableCertificateSignature {
  id: string
  url: string
  name: string
}

export interface CertificateRichTextController {
  elementId: string
  commit: () => void
  toggleBold: () => void
  toggleItalic: () => void
  toggleUnderline: () => void
  setAlign: (align: 'left' | 'center' | 'right' | 'justify') => void
  setLineHeight: (lineHeight: number) => void
  setColor: (color: string) => void
  setFontSize: (fontSize: number) => void
  setFontFamily: (fontFamily: string) => void
  insertText: (text: string) => void
}

export interface CertificateTextSelectionStyles {
  bold: boolean | null
  italic: boolean | null
  underline: boolean | null
  align: 'left' | 'center' | 'right' | 'justify' | null
  lineHeight: number | null
  color: string | null
  fontSize: number | null
  fontFamily: string | null
}

export interface CertificateEditorState {
  draft: CertificationTemplateDraft
  canvas: CertificateCanvasSize
  availableSignatures: AvailableCertificateSignature[]
  selectedElementId: string | null
  editingElementId: string | null
  richTextController: CertificateRichTextController | null
  textSelectionStyles: CertificateTextSelectionStyles | null
}

function createInitialDraft(): CertificationTemplateDraft {
  return {
    title: 'Certificado sem título',
    url: null,
    data: {
      canvas: { ...DEFAULT_CERTIFICATE_CANVAS },
      background: null,
      elements: [createHashElement(DEFAULT_CERTIFICATE_CANVAS)],
    },
  }
}

function createInitialState(): CertificateEditorState {
  return {
    draft: createInitialDraft(),
    canvas: { ...DEFAULT_CERTIFICATE_CANVAS },
    availableSignatures: [],
    selectedElementId: null,
    editingElementId: null,
    richTextController: null,
    textSelectionStyles: null,
  }
}

export const certificateEditorStore = new Store<CertificateEditorState>(
  createInitialState(),
)

function setElements(
  state: CertificateEditorState,
  elements: CertificationTemplateElement[],
): CertificateEditorState {
  return {
    ...state,
    draft: {
      ...state.draft,
      data: { ...state.draft.data, elements },
    },
  }
}

function updateState(
  updater: (state: CertificateEditorState) => CertificateEditorState,
): void {
  certificateEditorStore.setState(updater)
}

function normalizeLoadedElement(
  element: CertificationTemplateElement,
): CertificationTemplateElement {
  if (element.type === 'hash') {
    return { ...element, hash: '{{cert_hash}}', url: '{{verify_url}}' }
  }
  if (element.type === 'text') {
    return {
      ...element,
      paragraphs: element.paragraphs.map((paragraph) => ({
        ...paragraph,
        lineHeight: (paragraph as { lineHeight?: number }).lineHeight ?? 1.25,
      })),
    }
  }
  return { ...element }
}

export const certificateEditorActions = {
  loadDraft(draft: CertificationTemplateDraft): void {
    const canvas = draft.data.canvas ?? DEFAULT_CERTIFICATE_CANVAS
    updateState((state) => ({
      ...state,
      draft: {
        ...draft,
        data: {
          ...draft.data,
          canvas: { ...canvas },
          elements: draft.data.elements.map(normalizeLoadedElement),
        },
      },
      canvas: { ...canvas },
      selectedElementId: null,
      editingElementId: null,
      richTextController: null,
      textSelectionStyles: null,
    }))
  },

  reset(): void {
    certificateEditorStore.setState(() => createInitialState())
  },

  setTitle(title: string): void {
    updateState((state) => ({
      ...state,
      draft: { ...state.draft, title },
    }))
  },

  setBackgroundUrl(url: string | null): void {
    updateState((state) => ({
      ...state,
      draft: {
        ...state.draft,
        url,
        data: { ...state.draft.data, background: url },
      },
    }))
  },

  setAvailableSignatures(
    availableSignatures: AvailableCertificateSignature[],
  ): void {
    updateState((state) => ({ ...state, availableSignatures }))
  },

  setCanvasSize(canvas: CertificateCanvasSize): void {
    if (
      !Number.isFinite(canvas.width) ||
      !Number.isFinite(canvas.height) ||
      canvas.width < 320 ||
      canvas.height < 320 ||
      canvas.width > 6000 ||
      canvas.height > 6000
    ) {
      return
    }

    updateState((state) => {
      const scaleX = canvas.width / state.canvas.width
      const scaleY = canvas.height / state.canvas.height
      const fontScale = Math.sqrt(scaleX * scaleY)
      const elements = state.draft.data.elements.map((element) => {
        const scaled = {
          ...element,
          x: element.x * scaleX,
          y: element.y * scaleY,
          width: element.width * scaleX,
          height: element.height * scaleY,
        }

        if (scaled.type === 'text') {
          return {
            ...scaled,
            paragraphs: scaled.paragraphs.map((paragraph) => ({
              ...paragraph,
              runs: paragraph.runs.map((run) => ({
                ...run,
                fontSize: Math.max(6, Math.round(run.fontSize * fontScale)),
              })),
            })),
          }
        }

        if (scaled.type === 'hash') {
          return {
            ...scaled,
            fontSize: Math.max(6, Math.round(scaled.fontSize * fontScale)),
          }
        }

        return scaled
      })

      return {
        ...setElements(
          {
            ...state,
            draft: {
              ...state.draft,
              data: { ...state.draft.data, canvas: { ...canvas } },
            },
          },
          elements,
        ),
        canvas: { ...canvas },
      }
    })
  },

  addElement(element: CertificationTemplateElement): void {
    updateState((state) => {
      if (
        element.type === 'hash' &&
        state.draft.data.elements.some((item) => item.type === 'hash')
      ) {
        return state
      }

      return {
        ...setElements(state, [...state.draft.data.elements, element]),
        selectedElementId: element.id,
      }
    })
  },

  updateElement(
    id: string,
    updater: (
      element: CertificationTemplateElement,
    ) => CertificationTemplateElement,
  ): void {
    updateState((state) =>
      setElements(
        state,
        state.draft.data.elements.map((element) =>
          element.id === id ? updater(element) : element,
        ),
      ),
    )
  },

  updateElementBounds(
    id: string,
    bounds: Partial<
      Pick<CertificationTemplateElement, 'x' | 'y' | 'width' | 'height'>
    >,
  ): void {
    this.updateElement(id, (element) => ({ ...element, ...bounds }))
  },

  removeElement(id: string): void {
    const controller = certificateEditorStore.state.richTextController
    if (controller?.elementId === id) controller.commit()
    updateState((state) => {
      const target = state.draft.data.elements.find(
        (element) => element.id === id,
      )
      if (!target || target.type === 'hash') return state

      return {
        ...setElements(
          state,
          state.draft.data.elements.filter((element) => element.id !== id),
        ),
        selectedElementId:
          state.selectedElementId === id ? null : state.selectedElementId,
        editingElementId:
          state.editingElementId === id ? null : state.editingElementId,
        richTextController:
          state.richTextController?.elementId === id
            ? null
            : state.richTextController,
        textSelectionStyles:
          state.richTextController?.elementId === id
            ? null
            : state.textSelectionStyles,
      }
    })
  },

  selectElement(id: string | null): void {
    updateState((state) => ({ ...state, selectedElementId: id }))
  },

  startEditing(id: string): void {
    updateState((state) => ({
      ...state,
      selectedElementId: id,
      editingElementId: id,
    }))
  },

  stopEditing(): void {
    certificateEditorStore.state.richTextController?.commit()
    updateState((state) => ({
      ...state,
      editingElementId: null,
      richTextController: null,
      textSelectionStyles: null,
    }))
  },

  setRichTextController(
    richTextController: CertificateRichTextController | null,
  ): void {
    updateState((state) => ({ ...state, richTextController }))
  },

  setTextSelectionStyles(
    textSelectionStyles: CertificateTextSelectionStyles | null,
  ): void {
    updateState((state) => ({ ...state, textSelectionStyles }))
  },

  moveElement(id: string, targetIndex: number): void {
    updateState((state) => {
      const elements = [...state.draft.data.elements]
      const currentIndex = elements.findIndex((element) => element.id === id)
      if (currentIndex < 0) return state

      const nextIndex = Math.min(Math.max(targetIndex, 0), elements.length - 1)
      if (currentIndex === nextIndex) return state

      const [element] = elements.splice(currentIndex, 1)
      elements.splice(nextIndex, 0, element)
      return setElements(state, elements)
    })
  },

  bringForward(id: string): void {
    const state = certificateEditorStore.state
    const index = state.draft.data.elements.findIndex(
      (element) => element.id === id,
    )
    this.moveElement(id, index + 1)
  },

  sendBackward(id: string): void {
    const state = certificateEditorStore.state
    const index = state.draft.data.elements.findIndex(
      (element) => element.id === id,
    )
    this.moveElement(id, index - 1)
  },

  bringToFront(id: string): void {
    this.moveElement(
      id,
      certificateEditorStore.state.draft.data.elements.length,
    )
  },

  sendToBack(id: string): void {
    this.moveElement(id, 0)
  },

  getDraft(): CertificationTemplateDraft {
    return certificateEditorStore.state.draft
  },
}

export function useCertificateEditorState<T>(
  selector: (state: CertificateEditorState) => T,
): T {
  return useStore(certificateEditorStore, selector)
}
