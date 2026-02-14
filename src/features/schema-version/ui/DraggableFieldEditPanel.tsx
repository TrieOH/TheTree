import { useState, useEffect, useRef, useCallback } from 'react';
import { 
  useDraggable, 
  DndContext, 
  PointerSensor, 
  useSensor, 
  useSensors, 
  type DragEndEvent 
} from '@dnd-kit/core';
import { CSS } from '@dnd-kit/utilities';
import { cn } from '@/shared/lib/utils';
import { X } from 'lucide-react';
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';

interface DraggableFieldEditPanelProps {
  children: React.ReactNode;
  onClose: () => void;
  title: string;
}

interface DraggablePanelContentProps extends DraggableFieldEditPanelProps {
  isMobile: boolean;
  position: { x: number; y: number };
  setPosition: React.Dispatch<React.SetStateAction<{ x: number; y: number }>>;
}

const DraggablePanelContent: React.FC<DraggablePanelContentProps> = ({ 
  children, onClose, title, isMobile, position 
}) => {
  const panelRef = useRef<HTMLDivElement>(null);

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    isDragging,
  } = useDraggable({
    id: 'draggable-field-edit-panel',
    disabled: isMobile,
  });

  const currentTransform = isDragging && transform
    ? { x: position.x + transform.x, y: position.y + transform.y, scaleX: 1, scaleY: 1 }
    : { x: position.x, y: position.y, scaleX: 1, scaleY: 1 };

  const panelStyle = {
    transform: isMobile ? undefined : CSS.Transform.toString(currentTransform),
  };

  return (
    <div
      ref={(node) => { setNodeRef(node); panelRef.current = node; }}
      style={panelStyle}
      className={cn(
        "fixed bg-popover text-popover-foreground rounded-lg shadow-2xl border border-border z-50",
        "flex flex-col overflow-hidden w-[calc(100vw-20px)] md:w-96 max-w-full max-h-(--screen--minus-header)",
        isMobile && "top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2",
      )}
    >
      <div
        {...listeners}
        {...attributes}
        className={cn(
          "flex justify-between items-center p-2 border-b border-border",
          !isMobile && "cursor-grab active:cursor-grabbing"
        )}
      >
        <h4 className="font-semibold text-lg">{title}</h4>
        <ShadowButton variant="ghost" onClick={onClose} leftIcon={<X className="w-4 h-4"/>} className="p-1 h-auto" />
      </div>
      <div className="flex-1 overflow-y-auto">
        {children}
      </div>
    </div>
  );
};


export const DraggableFieldEditPanel: React.FC<DraggableFieldEditPanelProps> = ({ children, onClose, title }) => {
  const DRAGGABLE_PANEL_POSITION_KEY = 'draggable-field-edit-panel-position';

  const [position, setPosition] = useState(() => {
    try {
      const savedPosition = localStorage.getItem(DRAGGABLE_PANEL_POSITION_KEY);
      return savedPosition ? JSON.parse(savedPosition) : { x: 100, y: 100 };
    } catch (error) {
      console.error("Failed to parse saved position from localStorage", error);
      return { x: 100, y: 100 };
    }
  });
  const [isMobile, setIsMobile] = useState(false);
  const [dragStartPos, setDragStartPos] = useState({ x: 0, y: 0 });

  useEffect(() => {
    const handleResize = () => {
      const newIsMobile = window.innerWidth < 768;
      setIsMobile(newIsMobile);
    };
    window.addEventListener('resize', handleResize);
    handleResize();
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  useEffect(() => {
    return () => localStorage.setItem(DRAGGABLE_PANEL_POSITION_KEY, JSON.stringify(position));
  }, [position]);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 1 },
    })
  );

  const constrainPosition = useCallback((x: number, y: number, width: number, height: number) => {
    const margin = 10;
    const maxX = window.innerWidth - width - margin;
    const maxY = window.innerHeight - height - margin;
    
    return {
      x: Math.max(margin, Math.min(x, maxX)),
      y: Math.max(margin, Math.min(y, maxY)),
    };
  }, []);

  const handleDragStart = () => {
    setDragStartPos({ ...position });
  };

  const handleDragEnd = (event: DragEndEvent) => {
    if (isMobile) return;
    
    const panel = document.querySelector('[data-panel="field-edit"]') as HTMLElement;
    const width = panel?.offsetWidth || 384;
    const height = panel?.offsetHeight || 400;
    
    const newX = dragStartPos.x + event.delta.x;
    const newY = dragStartPos.y + event.delta.y;
    
    setPosition(constrainPosition(newX, newY, width, height));
  };

  return (
    <DndContext 
      sensors={sensors} 
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <DraggablePanelContent 
        title={title} 
        onClose={onClose} 
        isMobile={isMobile} 
        position={position}
        setPosition={setPosition}
      >
        {children}
      </DraggablePanelContent>
    </DndContext>
  );
};