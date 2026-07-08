import { useEffect, useRef, useCallback } from 'react'
import type { Canvas, Rect as FabricRect, Textbox, FabricImage } from 'fabric'
import type {
  CanvasElement,
  TextCanvasElement,
  ImageCanvasElement,
} from './types'
import {
  parseRichTextMarkup,
  plainIndexToMarkupIndex,
} from './rich-text'

function clamp(value: number, min: number, max: number): number {
  return Math.min(Math.max(value, min), max)
}

const DESIGN_WIDTH = 1000
const DESIGN_HEIGHT = 707

function getViewportTransform(canvasWidth: number, canvasHeight: number): [number, number, number, number, number, number] {
  const scale = Math.min(canvasWidth / DESIGN_WIDTH, canvasHeight / DESIGN_HEIGHT)
  const tx = (canvasWidth - DESIGN_WIDTH * scale) / 2
  const ty = (canvasHeight - DESIGN_HEIGHT * scale) / 2
  return [scale, 0, 0, scale, tx, ty]
}

interface CanvasElementObjects {
  field: FabricRect | null
  label: Textbox
  preview: FabricImage | null
}

interface FabricModule {
  Rect: typeof FabricRect
  Textbox: typeof Textbox
}

interface UseCertificateCanvasOptions {
  canvasHostRef: React.RefObject<HTMLDivElement | null>
  canvasElRef: React.RefObject<HTMLCanvasElement | null>
  backgroundUrl: string | null
  elements: CanvasElement[]
  selectedElementId: string | null
  signatureUrl: string | null
  onElementsChange: (elements: CanvasElement[]) => void
  onElementSelect: (id: string | null) => void
}

export function useCertificateCanvas({
  canvasHostRef,
  canvasElRef,
  backgroundUrl,
  elements,
  selectedElementId,
  signatureUrl,
  onElementsChange,
  onElementSelect,
}: UseCertificateCanvasOptions) {
  const canvasRef = useRef<Canvas | null>(null)
  const isReadyRef = useRef(false)
  const objectsRef = useRef<Map<string, CanvasElementObjects>>(new Map())
  const fabricRef = useRef<FabricModule | null>(null)
  const elementIdsRef = useRef<string[]>([])
  const elementsRef = useRef(elements)
  const editingTextElementIdRef = useRef<string | null>(null)
  const backgroundUrlRef = useRef(backgroundUrl)
  const backgroundLoadTokenRef = useRef(0)
  elementsRef.current = elements
  backgroundUrlRef.current = backgroundUrl

  const syncTextElementFromTextbox = useCallback((elementId: string, tb: Textbox) => {
    const currentElement = elementsRef.current.find((el) => el.id === elementId)
    if (!currentElement || currentElement.type !== 'text') return

    const textEl = currentElement as TextCanvasElement
    const parsed = parseRichTextMarkup(textEl.content, {
      fontFamily: textEl.fontFamily,
      fontSize: textEl.fontSize,
      color: textEl.color,
    })
    const nextPlainText = tb.text ?? ''
    const currentPlainText = parsed.plainText
    if (nextPlainText === currentPlainText) return

    let prefix = 0
    const maxPrefix = Math.min(currentPlainText.length, nextPlainText.length)
    while (prefix < maxPrefix && currentPlainText[prefix] === nextPlainText[prefix]) {
      prefix += 1
    }

    let oldSuffix = currentPlainText.length
    let newSuffix = nextPlainText.length
    while (
      oldSuffix > prefix &&
      newSuffix > prefix &&
      currentPlainText[oldSuffix - 1] === nextPlainText[newSuffix - 1]
    ) {
      oldSuffix -= 1
      newSuffix -= 1
    }

    const rawStart = plainIndexToMarkupIndex(parsed, prefix)
    const rawEnd = plainIndexToMarkupIndex(parsed, oldSuffix)
    const replacement = nextPlainText.slice(prefix, newSuffix)
    const nextContent = `${textEl.content.slice(0, rawStart)}${replacement}${textEl.content.slice(rawEnd)}`

    if (nextContent !== textEl.content) {
      onElementsChange(elementsRef.current.map((el) =>
        el.id === elementId ? { ...el, content: nextContent } as CanvasElement : el
      ))
    }
  }, [onElementsChange])

  const updateLabelPosition = useCallback((id: string, el: CanvasElement) => {
    const objs = objectsRef.current.get(id)
    if (!objs) return

    if (el.type === 'text') {
      objs.label.setCoords()
      return
    }

    const field = objs.field
    if (!field) return
    const center = field.getCenterPoint()
    const fw = field.getScaledWidth()
    const fh = field.getScaledHeight()
    const content = objs.preview

    objs.label.set({
      left: center.x,
      top: center.y,
      width: Math.max(120, fw * 0.88),
      fontSize: Math.max(14, Math.round(fh * 0.18)),
      visible: !content,
    })
    objs.label.setCoords()

    if (content) {
      const s = Math.min((fw * 0.88) / (content.width || 1), (fh * 0.88) / (content.height || 1))
      content.set({ left: center.x, top: center.y, scaleX: s, scaleY: s })
      content.setCoords()
    }
  }, [])

  const syncBackgroundImage = useCallback(async (targetCanvas?: Canvas | null) => {
    const canvas = targetCanvas ?? canvasRef.current
    if (!canvas || !isReadyRef.current) return

    const token = ++backgroundLoadTokenRef.current
    const { FabricImage } = await import('fabric')
    if (token !== backgroundLoadTokenRef.current) return

    const url = backgroundUrlRef.current
    if (!url) {
      canvas.backgroundImage = undefined
      canvas.renderAll()
      return
    }

    try {
      const img = await FabricImage.fromURL(url)
      if (token !== backgroundLoadTokenRef.current || !canvasRef.current) return

      const currentCanvas = canvasRef.current
      const scale = Math.max(
        DESIGN_WIDTH / (img.width || 1),
        DESIGN_HEIGHT / (img.height || 1)
      )

      img.set({
        left: DESIGN_WIDTH / 2,
        top: DESIGN_HEIGHT / 2,
        originX: 'center',
        originY: 'center',
        selectable: false,
        evented: false,
        scaleX: scale,
        scaleY: scale,
      })

      currentCanvas.backgroundImage = img
      currentCanvas.renderAll()
    } catch {
      // ignore background load errors
    }
  }, [])

  const applyStateToObjects = useCallback(() => {
    const canvas = canvasRef.current
    if (!canvas || !isReadyRef.current) return

    const cw = DESIGN_WIDTH
    const ch = DESIGN_HEIGHT
    let needsRender = false

    for (const element of elements) {
      const objs = objectsRef.current.get(element.id)
      if (!objs) continue

      if (element.type === 'text') {
        const textEl = element as TextCanvasElement
        const tb = objs.label
        const w = clamp(cw * (element.widthPct / 100), 60, cw * 0.85)
        const margin = w / 2
        const maxLeft = cw - margin
        const x = clamp(cw * (element.xPct / 100), margin, maxLeft)
        const y = clamp(ch * (element.yPct / 100), 20, ch - 20)
        const parsed = parseRichTextMarkup(textEl.content, {
          fontFamily: textEl.fontFamily,
          fontSize: textEl.fontSize,
          color: textEl.color,
        })

        tb.set({
          left: x,
          top: y,
          width: w,
          fontSize: textEl.fontSize,
          fontFamily: textEl.fontFamily,
          fontWeight: textEl.fontWeight,
          fill: textEl.color,
          text: parsed.plainText || 'Texto',
          styles: parsed.styles,
          editable: true,
        })
        tb.initDimensions()
        tb.setCoords()
        needsRender = true

      } else {
        const field = objs.field
        if (!field) continue

        const w = clamp(cw * (element.widthPct / 100), 60, cw * 1.5)
        const h = clamp(ch * (element.heightPct / 100), 30, ch * 0.8)
        const x = clamp(cw * (element.xPct / 100), -w * 0.3, cw + w * 0.3)
        const y = clamp(ch * (element.yPct / 100), -h * 0.3, ch + h * 0.3)

        field.set({ left: x, top: y, width: w, height: h, scaleX: 1, scaleY: 1 })
        field.setCoords()
        updateLabelPosition(element.id, element)
        needsRender = true
      }
    }

    if (needsRender) canvas.renderAll()
  }, [elements, updateLabelPosition])

  const createElementOnCanvas = useCallback((
    element: CanvasElement,
    x: number,
    y: number,
    w: number,
    h: number,
    canvas: Canvas,
    fabric: FabricModule
  ) => {
    const { Rect, Textbox } = fabric

    if (element.type === 'text') {
      const textEl = element as TextCanvasElement
      const parsed = parseRichTextMarkup(textEl.content, {
        fontFamily: textEl.fontFamily,
        fontSize: textEl.fontSize,
        color: textEl.color,
      })
      const tb = new Textbox(parsed.plainText || 'Texto', {
        left: x, top: y,
        originX: 'center', originY: 'center',
        width: w,
        fontFamily: textEl.fontFamily,
        fontSize: textEl.fontSize,
        fontWeight: textEl.fontWeight,
        fill: textEl.color,
        text: parsed.plainText || 'Texto',
        textAlign: 'center',
        splitByGrapheme: false,
        editable: true,
        transparentCorners: false, cornerStyle: 'circle',
        cornerColor: '#0f172a', cornerStrokeColor: '#ffffff',
        cornerSize: 10, lockRotation: true, lockScalingFlip: true,
        lockScalingY: true,
        styles: parsed.styles,
      })
      ;(tb as any).data = { elementId: element.id, role: 'text' }
      tb.initDimensions()

      tb.on('mousedblclick', () => {
        editingTextElementIdRef.current = element.id
        canvas.setActiveObject(tb)
        tb.enterEditing()
      })
      tb.on('editing:entered', () => {
        editingTextElementIdRef.current = element.id
      })
      tb.on('changed', () => {
        if (editingTextElementIdRef.current === element.id) {
          syncTextElementFromTextbox(element.id, tb)
        }
      })
      tb.on('editing:exited', () => {
        syncTextElementFromTextbox(element.id, tb)
        if (editingTextElementIdRef.current === element.id) {
          editingTextElementIdRef.current = null
        }
      })

      tb.on('modified', () => {
        // Clamp text center so textbox stays within canvas
        const c = tb.getCenterPoint()
        const cw = DESIGN_WIDTH
        const ch = DESIGN_HEIGHT
        const hw = tb.getScaledWidth() / 2
        const hh = tb.getScaledHeight() / 2
        const xClamped = clamp(c.x, hw, cw - hw)
        const yClamped = clamp(c.y, hh, ch - hh)
        if (xClamped !== c.x || yClamped !== c.y) {
          tb.set({ left: xClamped, top: yClamped })
          tb.setCoords()
        }
        onElementsChange(elementsRef.current.map(el =>
          el.id === element.id ? {
            ...el,
            xPct: clamp((c.x / cw) * 100, 5, 95),
            yPct: clamp((c.y / ch) * 100, 3, 97),
            widthPct: clamp((tb.getScaledWidth() / cw) * 100, 10, 85),
          } : el
        ))
      })
      tb.on('selected', () => onElementSelect(element.id))
      tb.on('deselected', () => onElementSelect(null))

      objectsRef.current.set(element.id, { field: null, label: tb, preview: null })
      canvas.add(tb)
      canvas.setActiveObject(tb)
      return
    }

    const isSignature = element.type === 'signature'
    const isImage = element.type === 'image'

    const field = new Rect({
      left: x, top: y,
      originX: 'center', originY: 'center',
      width: w, height: h,
      rx: isSignature ? 22 : 10, ry: isSignature ? 22 : 10,
      fill: isImage ? 'rgba(0,0,0,0.05)' : 'rgba(255,255,255,0.72)',
      stroke: '#0f172a', strokeWidth: 2,
      strokeDashArray: [10, 6] as unknown as number[],
      transparentCorners: false, cornerStyle: 'circle',
      cornerColor: '#0f172a', cornerStrokeColor: '#ffffff',
      cornerSize: 14, lockRotation: true, lockScalingFlip: true,
    })
    ;(field as any).data = { elementId: element.id, role: 'field' }

    const label = new Textbox(isSignature ? 'Assinatura' : 'Imagem', {
      left: x, top: y,
      originX: 'center', originY: 'center',
      width: Math.max(120, w * 0.88),
      fontFamily: 'Inter, system-ui, sans-serif',
      fontSize: Math.max(14, Math.round(h * 0.18)),
      fontWeight: '700', textAlign: 'center',
      fill: '#334155', selectable: false, evented: false,
      splitByGrapheme: false,
    })
    ;(label as any).data = { elementId: element.id, role: 'label' }

    field.on('scaling', () => { updateLabelPosition(element.id, element) })
      field.on('modified', () => {
        const c = field.getCenterPoint()
        const cw = DESIGN_WIDTH
        const ch = DESIGN_HEIGHT
        const fw = field.getScaledWidth()
        const fh = field.getScaledHeight()
      const xClamped = clamp(c.x, -fw * 0.3, cw + fw * 0.3)
      const yClamped = clamp(c.y, -fh * 0.3, ch + fh * 0.3)
      if (xClamped !== c.x || yClamped !== c.y) {
        field.set({ left: xClamped, top: yClamped })
        field.setCoords()
      }
      onElementsChange(elementsRef.current.map(el =>
        el.id === element.id ? {
          ...el,
          xPct: clamp((c.x / cw) * 100, -15, 115),
          yPct: clamp((c.y / ch) * 100, -15, 115),
          widthPct: clamp((fw / cw) * 100, 10, 150),
          heightPct: clamp((fh / ch) * 100, 6, 80),
        } : el
      ))
      updateLabelPosition(element.id, element)
    })
    field.on('selected', () => onElementSelect(element.id))
    field.on('deselected', () => onElementSelect(null))

    objectsRef.current.set(element.id, { field, label, preview: null })
    canvas.add(field, label)
    canvas.setActiveObject(field)

    if (isImage) {
      const imgEl = element as ImageCanvasElement
      if (imgEl.src) void loadImageForElement(element.id, imgEl.src)
    }
  }, [onElementsChange, onElementSelect, updateLabelPosition])

  const loadImageForElement = useCallback(async (elementId: string, src: string) => {
    const objs = objectsRef.current.get(elementId)
    const canvas = canvasRef.current
    if (!objs || !canvas) return
    if (objs.preview) { canvas.remove(objs.preview); objs.preview = null }
    try {
      const { FabricImage } = await import('fabric')
      const img = await FabricImage.fromURL(src)
      objs.preview = img
      img.set({ selectable: false, evented: false, originX: 'center', originY: 'center' })
      ;(img as any).data = { elementId, role: 'preview' }
      canvas.add(img)
      const el = elementsRef.current.find(e => e.id === elementId)
      if (el) updateLabelPosition(elementId, el)
      canvas.requestRenderAll()
    } catch { /* ignore */ }
  }, [updateLabelPosition])

  const syncFull = useCallback(() => {
    const canvas = canvasRef.current
    const fabric = fabricRef.current
    if (!canvas || !isReadyRef.current || !fabric) return

    const currentIds = new Set(elements.map(e => e.id))
    const map = objectsRef.current
    for (const [id] of map) {
      if (!currentIds.has(id)) {
        const objs = map.get(id)
        if (objs) {
          canvas.remove(objs.label)
          if (objs.field) canvas.remove(objs.field)
          if (objs.preview) canvas.remove(objs.preview)
        }
        map.delete(id)
      }
    }

    const cw = DESIGN_WIDTH
    const ch = DESIGN_HEIGHT

    for (const element of elements) {
      if (map.has(element.id)) continue
      const x = clamp(cw * (element.xPct / 100), 0, cw)
      const y = clamp(ch * (element.yPct / 100), 0, ch)
      const w = element.type === 'text' ? clamp(cw * (element.widthPct / 100), 60, cw * 0.85) : clamp(cw * (element.widthPct / 100), 60, cw * 1.5)
      const h = element.type === 'text' ? 40 : clamp(ch * (element.heightPct / 100), 30, ch * 0.8)
      createElementOnCanvas(element, x, y, w, h, canvas, fabric)
    }

    applyStateToObjects()
    elementIdsRef.current = elements.map(e => e.id)
  }, [elements, createElementOnCanvas, applyStateToObjects])

  useEffect(() => {
    const canvasEl = canvasElRef.current
    const hostEl = canvasHostRef.current
    if (!canvasEl || !hostEl) return

    const canvasElV: HTMLCanvasElement = canvasEl
    const hostElV: HTMLDivElement = hostEl
    let cancelled = false
    let instance: Canvas | null = null
    let cleanupScroll = () => {}

    async function init() {
      const fabricModule = await import('fabric')
      const { Canvas } = fabricModule
      if (cancelled) return

      instance = new Canvas(canvasElV, {
        selection: false, preserveObjectStacking: true, renderOnAddRemove: true,
        allowTouchScrolling: true, uniformScaling: false, uniScaleKey: null,
        centeredScaling: false, stopContextMenu: false, fireRightClick: true,
        skipOffscreen: false,
        backgroundColor: '#ffffff',
      })

      canvasRef.current = instance
        ; (window as any).__fabricCanvas = instance
      fabricRef.current = { Rect: fabricModule.Rect, Textbox: fabricModule.Textbox }

      const size = hostElV.getBoundingClientRect()
      instance.setDimensions({ width: Math.max(320, Math.round(size.width)), height: Math.max(220, Math.round(size.height)) })
      instance.setViewportTransform(getViewportTransform(instance.getWidth(), instance.getHeight()))
      instance.renderAll()

      const observer = new ResizeObserver(() => {
        if (!canvasRef.current || !canvasHostRef.current) return
        const ns = canvasHostRef.current.getBoundingClientRect()
        const width = Math.max(320, Math.round(ns.width))
        const height = Math.max(220, Math.round(ns.height))
        if (width < 10 || height < 10) return
        canvasRef.current.setDimensions({ width, height })
        canvasRef.current.setViewportTransform(getViewportTransform(width, height))
        canvasRef.current.calcOffset()
        void syncBackgroundImage(canvasRef.current)
        canvasRef.current.renderAll()
      })
      observer.observe(hostElV)

      let scrollFrame: number | null = null
      const handleScroll = () => {
        if (!canvasRef.current) return
        if (scrollFrame !== null) return
        scrollFrame = window.requestAnimationFrame(() => {
          scrollFrame = null
          if (!canvasRef.current) return
          canvasRef.current.calcOffset()
          canvasRef.current.renderAll()
        })
      }

      document.addEventListener('scroll', handleScroll, true)
      cleanupScroll = () => {
      document.removeEventListener('scroll', handleScroll, true)
        if (scrollFrame !== null) {
          window.cancelAnimationFrame(scrollFrame)
        }
      }
      isReadyRef.current = true
      syncFull()
      void syncBackgroundImage(instance)
    }

    void init()
    return () => {
      cancelled = true
      isReadyRef.current = false
      cleanupScroll()
      instance?.dispose()
      canvasRef.current = null
      objectsRef.current.clear()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    void syncBackgroundImage()
  }, [backgroundUrl, syncBackgroundImage])

  useEffect(() => {
    const c = canvasRef.current
    if (!c || !isReadyRef.current) return
    let cancelled = false
    async function updatePreviews() {
      const { FabricImage } = await import('fabric')
      if (cancelled || !canvasRef.current) return
      const cur = canvasRef.current
      for (const [, objs] of objectsRef.current) { if (objs.preview) { cur.remove(objs.preview); objs.preview = null } }
      for (const [id, objs] of objectsRef.current) {
        const el = elements.find(e => e.id === id)
        if (!el) continue
        let imgSrc: string | null = null
        if (el.type === 'signature') imgSrc = signatureUrl
        else if (el.type === 'image') imgSrc = (el as ImageCanvasElement).src
        if (!imgSrc) continue
        try {
          const img = await FabricImage.fromURL(imgSrc)
          if (cancelled || !canvasRef.current) continue
          img.set({ selectable: false, evented: false, originX: 'center', originY: 'center' })
          ;(img as any).data = { elementId: id, role: 'preview' }
          objs.preview = img; cur.add(img)
          updateLabelPosition(id, el)
        } catch { /* skip */ }
      }
      cur.renderAll()
    }
    void updatePreviews()
    return () => { cancelled = true }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [signatureUrl, elements, updateLabelPosition])

  useEffect(() => {
    if (!isReadyRef.current) return
    const currentIds = elements.map(e => e.id).sort().join(',')
    const prevIds = elementIdsRef.current.sort().join(',')
    if (currentIds !== prevIds) { syncFull() } else { applyStateToObjects() }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [elements])

  // Sync selectedElementId from sidebar click to fabric canvas selection
  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas || !isReadyRef.current) return
    if (!selectedElementId) {
      canvas.discardActiveObject()
      canvas.requestRenderAll()
      return
    }
    const objs = objectsRef.current.get(selectedElementId)
    if (!objs) return
    const target = objs.field || objs.label
    if (target) {
      canvas.setActiveObject(target)
      canvas.requestRenderAll()
    }
  }, [selectedElementId])

  return { canvasRef, isReady: isReadyRef.current, syncFull, applyStateToObjects, loadImageForElement }
}
