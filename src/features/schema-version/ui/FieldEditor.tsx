import {
  DndContext, 
  PointerSensor, 
  useSensor, 
  useSensors, 
  type DragEndEvent, 
  type DragStartEvent 
} from "@dnd-kit/core";
import { arrayMove, SortableContext } from "@dnd-kit/sortable";
import FieldCard from "./FieldCard";
import { useState } from "react";
import { createPortal } from "react-dom";
import { DragOverlay } from "@dnd-kit/core";
import type { VersionField, VersionFieldList } from "../model/types";
import { defaultEmailVersionField, defaultPasswordVersionField, defaultVersionFieldList } from "../model/default";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Plus } from "lucide-react"; 

export default function FieldEditor() {
  const [items, setItems] = useState<VersionFieldList>(defaultVersionFieldList);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [nextId, setNextId] = useState(defaultVersionFieldList.length);

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

  const handleAddField = () => {
    const newField: VersionField = {
      key: `custom_field_${nextId}`,
      title: `Custom Field ${nextId}`,
      type: "string",
      owner: "user",
      mutable: true,
      required: false,
      position: items.length,
      default_value: null,
      options: [],
      required_rules: [],
      visibility_rules: []
    };
    setItems(currentItems => [...currentItems, newField]);
    setNextId(prevId => prevId + 1);
  };

  const handleDeleteField = (fieldKey: string) => {
    setItems(currentItems => currentItems.filter(item => item.key !== fieldKey));
  };

  const handleEditField = (fieldKey: string) => {
    console.log(`Edit field with key: ${fieldKey}`);
  };

  const activeItem = activeId ? items.find(item => item.key === activeId) : null;

  return (
    <main className="flex w-full min-h-(--screen--minus-header)">
      <DndContext
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onDragCancel={handleDragCancel}
        sensors={sensors}
      >
        <div className="flex-1 max-w-79 border-r border-r-border py-4 px-2 space-y-2">
          <FieldCard key={defaultEmailVersionField.key} field={defaultEmailVersionField} isFixed={true} />
          <SortableContext items={items.map(item => item.key)}>
            {items.map((item) => (
              <FieldCard 
                key={item.key} 
                field={item} 
                onEdit={handleEditField} 
                onDelete={handleDeleteField} 
              />
            ))}
          </SortableContext>
          <ShadowButton
            onClick={handleAddField} 
            className="w-full justify-center"
            value="Add Field" 
            variant="solid"
            leftIcon={<Plus className="w-4 h-4" />} 
          />
          <FieldCard 
            key={defaultPasswordVersionField.key} 
            field={{...defaultPasswordVersionField}} 
            overwriteType="password"
            isFixed={true}
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