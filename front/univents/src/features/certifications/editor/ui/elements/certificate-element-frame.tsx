import { Trash2 } from 'lucide-react'
import type { ReactNode } from 'react'
import {
  CERTIFICATE_ELEMENT_OVERFLOW,
  MIN_CERTIFICATE_ELEMENT_SIZE,
} from '../../constants'
import { useDragResize } from '../../hooks/use-drag-resize'
import type { ElementBounds, ResizeHandle } from '../../hooks/use-drag-resize'
import type { CertificateElementType } from '../../types'
import { cn } from '@/shared/lib/utils'

const RESIZE_HANDLES: ResizeHandle[] = ['nw', 'ne', 'sw', 'se']

const RESIZE_HANDLE_CLASS: Record<ResizeHandle, string> = {
  nw: '-left-1.5 -top-1.5 cursor-nwse-resize',
  ne: '-right-1.5 -top-1.5 cursor-nesw-resize',
  sw: '-bottom-1.5 -left-1.5 cursor-nesw-resize',
  se: '-right-1.5 -bottom-1.5 cursor-nwse-resize',
}

interface CertificateElementFrameProps {
  type: CertificateElementType
  bounds: ElementBounds
  scale: number
  canvas: { width: number; height: number }
  selected: boolean
  editing: boolean
  deletable: boolean
  onSelect: () => void
  onDoubleClick?: () => void
  onChangeBounds: (bounds: ElementBounds) => void
  onDelete?: () => void
  children: ReactNode
}

export function CertificateElementFrame({
  type,
  bounds,
  scale,
  canvas,
  selected,
  editing,
  deletable,
  onSelect,
  onDoubleClick,
  onChangeBounds,
  onDelete,
  children,
}: CertificateElementFrameProps) {
  const { startDrag, startResize } = useDragResize({
    bounds,
    scale,
    canvas,
    overflowAllowance: CERTIFICATE_ELEMENT_OVERFLOW[type],
    minWidth: MIN_CERTIFICATE_ELEMENT_SIZE.width,
    minHeight: MIN_CERTIFICATE_ELEMENT_SIZE.height,
    onChange: onChangeBounds,
  })

  return (
    <div
      data-certificate-element={type}
      className={cn(
        'absolute select-none',
        selected ? 'z-20' : 'z-10',
        !editing && (selected ? 'cursor-move' : 'cursor-pointer'),
      )}
      style={{
        left: bounds.x,
        top: bounds.y,
        width: bounds.width,
        height: bounds.height,
      }}
      onPointerDown={(event) => {
        onSelect()
        if (!editing) startDrag(event)
      }}
      onDoubleClick={onDoubleClick}
    >
      <div
        className={cn(
          'h-full w-full',
          selected && !editing && 'outline-2 outline-offset-2 outline-ring',
          editing && 'outline-2 outline-offset-2 outline-accent',
        )}
      >
        {children}
      </div>

      {selected && !editing ? (
        <>
          {RESIZE_HANDLES.map((handle) => (
            <div
              key={handle}
              aria-hidden="true"
              onPointerDown={startResize(handle)}
              className={cn(
                'absolute size-3 rounded-full border-2 border-ring bg-popover shadow-sm',
                RESIZE_HANDLE_CLASS[handle],
              )}
            />
          ))}

          {deletable && onDelete ? (
            <button
              type="button"
              onPointerDown={(event) => event.stopPropagation()}
              onClick={onDelete}
              title="Excluir elemento"
              aria-label="Excluir elemento"
              className="absolute -top-3.5 -right-3.5 flex size-7 items-center justify-center rounded-full bg-destructive text-destructive-foreground shadow-sm hover:opacity-90"
            >
              <Trash2 className="size-3.5" />
            </button>
          ) : null}
        </>
      ) : null}
    </div>
  )
}
