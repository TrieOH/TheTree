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
import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { DragOverlay } from "@dnd-kit/core";
import type { VersionField, VersionFieldList } from "../model/types";
import { DraggableFieldEditPanel } from './DraggableFieldEditPanel';
import { FieldEditForm } from './FieldEditForm';
import { defaultEmailVersionField, defaultPasswordVersionField, defaultVersionFieldList } from "../model/default";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Plus } from "lucide-react"; 

export default function FieldEditor() {
  const [items, setItems] = useState<VersionFieldList>(defaultVersionFieldList);
  const [originalItems, setOriginalItems] = useState<VersionFieldList>(defaultVersionFieldList);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [nextId, setNextId] = useState(defaultVersionFieldList.length);
  const [editingField, setEditingField] = useState<VersionField | null>(null);
  const [hasChanges, setHasChanges] = useState(false);

  useEffect(() => {
    setHasChanges(JSON.stringify(items) !== JSON.stringify(originalItems));
  }, [items, originalItems]);

  const handleOpenEditPanel = (field: VersionField) => {
    setEditingField(field);
  };

  const handleCloseEditPanel = () => {
    setEditingField(null);
  };

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
        const newItems = arrayMove(currentItems, oldIndex, newIndex);
        return newItems.map((item, index) => ({ ...item, position: index }));
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

  const handleUpdateField = (updatedField: VersionField) => {
    setItems(currentItems =>
      currentItems.map(item =>
        item.key === updatedField.key ? updatedField : item
      )
    );
    if (editingField?.key === updatedField.key) setEditingField(updatedField);
  };

  const handleSubmit = () => {
    console.log("Submitting changes:", items);
    setOriginalItems(items);
  };

  const activeItem = activeId ? items.find(item => item.key === activeId) : null;

  const allFieldKeys = [
    ...items.map(item => item.key),
    defaultEmailVersionField.key,
    defaultPasswordVersionField.key
  ];

  return (
    <main className="flex w-full min-h-(--screen--minus-header)">
      <DndContext
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onDragCancel={handleDragCancel}
        sensors={sensors}
      >
        <div className="flex-1 max-w-79 border-r border-r-border py-4 px-2 space-y-2">
          <FieldCard 
            key={defaultEmailVersionField.key} 
            field={defaultEmailVersionField} 
            isFixed={true} 
            onOpenEditPanel={handleOpenEditPanel}
          />
          <SortableContext items={items.map(item => item.key)}>
            {items.map((item) => (
              <FieldCard
                key={item.key}
                field={item}
                onDelete={handleDeleteField}
                onUpdateField={handleUpdateField}
                onOpenEditPanel={handleOpenEditPanel}
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
            onOpenEditPanel={handleOpenEditPanel}
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
      <div className="flex-1 p-4 flex flex-col justify-end items-end">
        <ShadowButton
          onClick={handleSubmit}
          disabled={!hasChanges}
          value="Submit Changes"
          variant="solid"
          className="w-fit"
        />
      </div>
      {editingField && (
        <DraggableFieldEditPanel onClose={handleCloseEditPanel} title="Edit Field">
          <FieldEditForm
            field={editingField}
            onSave={(updatedField) => {
              handleUpdateField(updatedField);
              handleCloseEditPanel();
            }}
            onCancel={handleCloseEditPanel}
            allFieldKeys={allFieldKeys}
          />
        </DraggableFieldEditPanel>
      )}
    </main>
  )
}