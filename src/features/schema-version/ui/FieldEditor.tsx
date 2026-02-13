import {DndContext, PointerSensor, useSensor, useSensors, type DragEndEvent, DragStartEvent } from "@dnd-kit/core";
import { arrayMove, SortableContext } from "@dnd-kit/sortable";
import FieldCard from "./FieldCard";
import { useState } from "react";
import { createPortal } from "react-dom";
import { DragOverlay } from "@dnd-kit/core";

interface Item {
  id: string;
  isFixed?: boolean;
}

export default function FieldEditor() {
  const [items, setItems] = useState<Item[]>([
    { id: "fixed-start", isFixed: true },
    { id: "dwdw" },
    { id: "fx" },
    { id: "fixed-end", isFixed: true },
  ]);
  const [activeId, setActiveId] = useState<string | null>(null);

  const sensors = useSensors(
  useSensor(PointerSensor, {
    activationConstraint: {
      distance: 6,
    },
  })
);
  
  const handleDragStart = (ev: DragStartEvent) => {
    setActiveId(ev.active.id as string);
  }

  const handleDragEnd = (ev: DragEndEvent) => {
    const { active, over } = ev;
    if(over && active.id !== over.id) {
      setItems(currentItems => {
        const activeItem = currentItems.find(item => item.id === active.id);
        const overItem = currentItems.find(item => item.id === over.id);

        if (activeItem?.isFixed || overItem?.isFixed) {
          return currentItems; // Do not reorder if either is fixed
        }

        const oldIndex = currentItems.findIndex(item => item.id === active.id);
        const newIndex = currentItems.findIndex(item => item.id === over.id);
        return arrayMove(currentItems, oldIndex, newIndex);
      });
    }
    setActiveId(null);
  }

  const handleDragCancel = () => {
    setActiveId(null);
  }

  const sortableItems = items.filter(item => !item.isFixed);
  const fixedStartItem = items.find(item => item.id === "fixed-start");
  const fixedEndItem = items.find(item => item.id === "fixed-end");

  const activeItem = activeId ? items.find(item => item.id === activeId) : null;

  return (
    <main className="flex w-full h-(--screen--minus-header)">
      <DndContext
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onDragCancel={handleDragCancel}
        sensors={sensors}
      >
        <div className="flex-1 max-w-79 border-r border-r-border p-2 space-y-2">
            {fixedStartItem && <FieldCard key={fixedStartItem.id} id={fixedStartItem.id} isFixed={true} />}
            <SortableContext items={sortableItems}>
                {sortableItems.map((item) => (
                <FieldCard key={item.id} id={item.id}/>
                ))}
            </SortableContext>
            {fixedEndItem && <FieldCard key={fixedEndItem.id} id={fixedEndItem.id} isFixed={true} />}
          </div>
        {createPortal(
          <DragOverlay>
            {activeItem ? (
              <FieldCard id={activeItem.id} isFixed={activeItem.isFixed} className="shadow-2xl scale-105 ring-2 ring-primary" />
            ) : null}
          </DragOverlay>,
          document.body
        )}
      </DndContext>
      <div className="flex-1"></div>
    </main>
  )
}