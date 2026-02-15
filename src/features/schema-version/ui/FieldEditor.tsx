import CustomTabs from '@/widgets/tabs/ui/CustomTabs';
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
import { useEffect, useState, useCallback } from "react";
import { createPortal } from "react-dom";
import { DragOverlay } from "@dnd-kit/core";
import type { VersionField, VersionFieldList, VersionFieldResult } from "../model/types";
import { DraggableFieldEditPanel } from './DraggableFieldEditPanel';
import { FieldEditForm } from './FieldEditForm';
import { defaultEmailVersionField, defaultPasswordVersionField, defaultVersionFieldList } from "../model/default";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Plus, SaveAll } from "lucide-react"; 
import { useMutation, useQueryClient, useQuery } from "@tanstack/react-query";
import { publishSchemaVersionFieldFn, schemaVersionByIdQueryOptions } from "../api";
import { toast } from "sonner";
import { useStore } from "@tanstack/react-store";
import { navigationStore } from "@/features/navigation";
import PublishSchemaVersionButton from './PublishSchemaVersionButton';

export default function FieldEditor() {
  const queryClient = useQueryClient();
  const { currentProjectId, currentSchemaId, currentSchemaVersion } = useStore(navigationStore);
  const isVersionNull = currentSchemaVersion === null || currentSchemaVersion === undefined;

  const { data: schemaVersionData } = useQuery(schemaVersionByIdQueryOptions(currentProjectId || "", currentSchemaId || "", currentSchemaVersion || 1));

  const [items, setItems] = useState<VersionFieldList>([]);
  const [originalItems, setOriginalItems] = useState<VersionFieldList>([]);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [nextId, setNextId] = useState(defaultVersionFieldList.length);
  const [editingField, setEditingField] = useState<VersionField | null>(null);
  const [hasChanges, setHasChanges] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };

    checkMobile();
    window.addEventListener("resize", checkMobile);

    return () => window.removeEventListener("resize", checkMobile);
  }, []);


  const mapFieldIdsToKeys = useCallback((fields: VersionFieldResult[]) => {
    const fieldMap = new Map(
      fields.map(field => [field.object_id, field.key])
    );

    return fields.map(field => ({
      ...field,
      required_rules: field.required_rules.map(rule => ({
        operator: rule.operator,
        value: rule.value,
        depends_on_field_key: fieldMap.get(rule.depends_on_field_id) ?? ""
      })),
      visibility_rules: field.visibility_rules.map(rule => ({
        operator: rule.operator,
        value: rule.value,
        depends_on_field_key: fieldMap.get(rule.depends_on_field_id) ?? ""
      }))
    }));
  }, []);


  useEffect(() => {
    if (schemaVersionData) {
      const mappedFields = mapFieldIdsToKeys(schemaVersionData.fields);  
      setItems(mappedFields);
      setOriginalItems(mappedFields);
    }
  }, [schemaVersionData, mapFieldIdsToKeys]);

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

  const handleUpdateField = (originalField: VersionField, updatedField: VersionField) => {
    setItems(currentItems =>
      currentItems.map(item =>
        item.key === originalField.key ? updatedField : item
      )
    );
    if (editingField?.key === originalField.key) setEditingField(updatedField);
  };

  const publishVersionFieldSchemaMutation = useMutation({
    mutationFn: publishSchemaVersionFieldFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ queryKey: ["schemaVersionById"] });
      }
    },
    onError: (error) => {
      toast.error(`Failed to publish schema: ${error.message}`);
    }
  });

  const handleSubmit = () => {
    publishVersionFieldSchemaMutation.mutate({
      fields: items, 
      project_id: currentProjectId || "",
      schema_id: currentSchemaId || "",
      version: currentSchemaVersion || 1
    })
    setOriginalItems(items);
  };

  const activeItem = activeId ? items.find(item => item.key === activeId) : null;

  const allFieldKeys = [
    ...items.map(item => item.key),
    defaultEmailVersionField.key,
    defaultPasswordVersionField.key
  ];

  return (
    <main className="flex w-full md:min-h-(--screen--minus-header) md:h-auto h-(--screen--minus-header)">
      <DndContext
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onDragCancel={handleDragCancel}
        sensors={sensors}
      >
        {isMobile ? (
          <CustomTabs
            items={[
              {
                value: "field",
                label: "Field",
                content: (
                  <div className="flex-1 py-4 px-2 space-y-2">
                    <FieldCard
                      key={defaultEmailVersionField.key}
                      field={defaultEmailVersionField}
                      isFixed={true}
                    />
                    <SortableContext items={items.map((item) => item.key)}>
                      {items.map((item) => (
                        <FieldCard
                          key={item.key}
                          field={item}
                          onDelete={handleDeleteField}
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
                      disabled={isVersionNull}
                    />
                    <FieldCard
                      key={defaultPasswordVersionField.key}
                      field={{ ...defaultPasswordVersionField }}
                      overwriteType="password"
                      isFixed={true}
                    />
                  </div>
                ),
              },
              {
                value: "preview",
                label: "Preview",
                content: (
                  <div className="flex-1 p-4">
                    {editingField && (
                      <FieldEditForm
                        field={editingField}
                        onSave={(originalField, updatedField) => {
                          handleUpdateField(originalField, updatedField);
                        }}
                        onCancel={handleCloseEditPanel}
                        allFieldKeys={allFieldKeys}
                      />
                    )}
                  </div>
                ),
              },
            ]}
          />
        ) : (
            <div className="flex-1 max-w-79 border-r border-r-border py-4 px-2 space-y-2">
              <FieldCard
                key={defaultEmailVersionField.key}
                field={defaultEmailVersionField}
                isFixed={true}
              />
              <SortableContext items={items.map((item) => item.key)}>
                {items.map((item) => (
                  <FieldCard
                    key={item.key}
                    field={item}
                    onDelete={handleDeleteField}
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
                disabled={isVersionNull}
              />
              <FieldCard
                key={defaultPasswordVersionField.key}
                field={{ ...defaultPasswordVersionField }}
                overwriteType="password"
                isFixed={true}
              />
            </div>
        )}

        {createPortal(
          <DragOverlay>
            {activeItem ? (
              <FieldCard
                field={activeItem}
                className="shadow-2xl scale-105 ring-2 ring-primary"
              />
            ) : null}
          </DragOverlay>,
          document.body
        )}
      </DndContext>
      {editingField && (
        <DraggableFieldEditPanel onClose={handleCloseEditPanel} title="Edit Field">
          <FieldEditForm
            field={editingField}
            onSave={(originalField, updatedField) => {
              handleUpdateField(originalField, updatedField);
              handleCloseEditPanel();
            }}
            onCancel={handleCloseEditPanel}
            allFieldKeys={allFieldKeys}
          />
        </DraggableFieldEditPanel>
      )}
      <div className="fixed right-4 md:bottom-4 bottom-16 flex flex-col items-center gap-2">
        <ShadowButton
          onClick={handleSubmit}
          disabled={!hasChanges || isVersionNull}
          leftIcon={<SaveAll className="w-4 h-4" />}
          value={isMobile ? '' : 'Save Fields'}
          variant="solid"
        />
        <PublishSchemaVersionButton
          items={items}
          isMobile={isMobile}
          hasChanges={!isVersionNull}
          setOriginalItems={setOriginalItems}
        />
      </div>
    </main>
  );
}