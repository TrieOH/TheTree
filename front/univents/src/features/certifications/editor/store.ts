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

export interface CertificateEditorState {
  draft: CertificationTemplateDraft
  canvas: CertificateCanvasSize
  availableSignatures: AvailableCertificateSignature[]
  selectedElementId: string | null
  editingElementId: string | null
}

function createInitialDraft(): CertificationTemplateDraft {
  return {
    title: 'Certificado sem título',
    url: null,
    data: {
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

export const certificateEditorActions = {
  loadDraft(draft: CertificationTemplateDraft): void {
    updateState((state) => ({
      ...state,
      draft: {
        ...draft,
        data: {
          ...draft.data,
          elements: draft.data.elements.map((element) =>
            element.type === 'hash'
              ? { ...element, hash: '{{cert_hash}}', url: '{{verify_url}}' }
              : { ...element },
          ),
        },
      },
      selectedElementId: null,
      editingElementId: null,
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
      canvas.height < 320
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
        ...setElements(state, elements),
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
    updateState((state) => ({ ...state, editingElementId: null }))
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
