import {DndContext, PointerSensor, useSensor, useSensors, type DragEndEvent, DragStartEvent } from "@dnd-kit/core";
import { arrayMove, SortableContext } from "@dnd-kit/sortable";
import FieldCard from "./FieldCard";
import { useState } from "react";
import { createPortal } from "react-dom";
import { DragOverlay } from "@dnd-kit/core";
import type { VersionFieldList } from "../model/types";
import { defaultEmailVersionField, defaultPasswordVersionField, defaultVersionFieldList } from "../model/default";

export default function FieldEditor() {
  const [items, setItems] = useState<VersionFieldList>(defaultVersionFieldList);
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
        const oldIndex = currentItems.findIndex(item => item.key === active.id);
        const newIndex = currentItems.findIndex(item => item.key === over.id);
        return arrayMove(currentItems, oldIndex, newIndex);
      });
    }
    setActiveId(null);
  }

  const handleDragCancel = () => {
    setActiveId(null);
  }

  const activeItem = activeId ? items.find(item => item.key === activeId) : null;

  return (
    <main className="flex w-full h-(--screen--minus-header)">
      <DndContext
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onDragCancel={handleDragCancel}
        sensors={sensors}
      >
        <div className="flex-1 max-w-79 border-r border-r-border p-2 space-y-2">
          <FieldCard key={defaultEmailVersionField.key} field={defaultEmailVersionField} />
          <SortableContext items={items.map(item => item.key)}>
            {items.map((item) => (
              <FieldCard key={item.key} field={item}/>
            ))}
          </SortableContext>
          <FieldCard 
            key={defaultPasswordVersionField.key} 
            field={{...defaultPasswordVersionField}} 
            overwrite_type="password"
          />
        </div>
        {createPortal(
          <DragOverlay>
            {activeItem ? (
              <FieldCard field={activeItem} className="shadow-2xl scale-105 ring-2 ring-primary" />
            ) : null}
          </DragOverlay>,
          document.body
        )}
      </DndContext>
      <div className="flex-1"></div>
    </main>
  )
}