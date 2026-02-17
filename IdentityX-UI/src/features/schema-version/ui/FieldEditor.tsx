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
import { useEffect, useState, useMemo } from "react";
import { createPortal } from "react-dom";
import { DragOverlay } from "@dnd-kit/core";
import type { VersionFieldResult } from "../model/types";
import { DraggableFieldEditPanel } from './DraggableFieldEditPanel';
import { FieldEditForm } from './FieldEditForm';
import { defaultEmailVersionField, defaultPasswordVersionField } from "../model/default";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Plus, SaveAll } from "lucide-react"; 
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { createSchemaVersionFieldFn, deleteSchemaFieldOptionFn, deleteSchemaVersionFieldFn, schemaVersionByIdQueryOptions, setSchemaFieldOptionsFn, updateSchemaVersionFieldFn } from "../api";
import { useStore } from "@tanstack/react-store";
import { navigationStore } from "@/features/navigation";
import { useEditableList } from '../hooks/useEditableList';
import { mapFieldIdsToKeys, mapVersionFieldResultToVersionField } from '../lib/convert-field-utils';
import { areFieldsEqual } from '../lib/field-utils';
import PublishSchemaVersionButton from './PublishSchemaVersionButton';
import { optionsDiff } from '../lib/field-options-diff-utils';

export default function FieldEditor() {
  const queryClient = useQueryClient();
  const { currentProjectId, currentSchemaId, currentSchemaVersion } = useStore(navigationStore);
  const isVersionNull = currentSchemaVersion === null || currentSchemaVersion === undefined;

  const [nextId, setNextId] = useState(-1);

  const { data: schemaVData } = useQuery(schemaVersionByIdQueryOptions(
    currentProjectId || "", currentSchemaId || "", currentSchemaVersion || 1
  ));

  const mappedFields = useMemo(() => {
    if (!schemaVData) return [];
    return mapFieldIdsToKeys(schemaVData.fields);
  }, [schemaVData]);

  const [activeId, setActiveId] = useState<string | null>(null);
  useEffect(() => {
    const allKeys = [
      ...mappedFields.map(field => field.key),
      defaultEmailVersionField.key,
      defaultPasswordVersionField.key
    ];

    let maxSuffix = -1;
    allKeys.forEach(key => {
      const match = key.match(/custom_field_(\d+)/);
      if (match && match[1]) {
        const suffix = parseInt(match[1], 10);
        if (!isNaN(suffix) && suffix > maxSuffix) maxSuffix = suffix;
      }
    });
    setNextId(maxSuffix + 1);
  }, [mappedFields]);
  const [editingField, setEditingField] = useState<VersionFieldResult | null>(null);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };

    checkMobile();
    window.addEventListener("resize", checkMobile);

    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  const handleOpenEditPanel = (field: VersionFieldResult) => {
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

  const fields = useEditableList<VersionFieldResult>({
    initial: mappedFields,

    getId: (f) => f.object_id,

    isEqual: areFieldsEqual,

    onCreate: async (creates) => {
      console.log("Create:", creates);
      if(!currentProjectId || !currentSchemaId || !currentSchemaVersion) return;
      await createSchemaVersionFieldFn({
        fields: creates.map(item => mapVersionFieldResultToVersionField(item)),
        project_id: currentProjectId,
        schema_id: currentSchemaId,
        version: currentSchemaVersion,
      });
    },

    onUpdate: async (updates) => {
      console.log("Update:", updates);
      if(!currentProjectId || !currentSchemaId || !currentSchemaVersion) return;
      const ordered = [...updates].sort((a, b) => a.value.position - b.value.position);
      for (const u of ordered) {
        await updateSchemaVersionFieldFn({
          project_id: currentProjectId,
          schema_id: currentSchemaId,
          version: currentSchemaVersion,
          field_id: u.id,
          fieldData: mapVersionFieldResultToVersionField(u.value),
        });
      }
    },

    onDelete: async (deletes) => {
      console.log("DELETE:", deletes)
      if(!currentProjectId || !currentSchemaId || !currentSchemaVersion) return;
      await Promise.all(
        deletes.map(d =>
          deleteSchemaVersionFieldFn({
            project_id: currentProjectId,
            schema_id: currentSchemaId,
            version: currentSchemaVersion,
            field_id: d.id,
          })
        )
      );
    },
    customDiffs: [
      optionsDiff({
        deleteOptions: async (fieldId, optionsDiff) => {
          console.log("DELETE OPTIONS:", optionsDiff);
          if(!currentProjectId || !currentSchemaId || !currentSchemaVersion) return;
          await Promise.all(
            optionsDiff.map(id =>
              deleteSchemaFieldOptionFn({
                project_id: currentProjectId,
                schema_id: currentSchemaId,
                version: currentSchemaVersion,
                field_id: fieldId,
                option_id: id,
              })
            )
          );
        },
        putOptions: async (fieldId, options) => {
          console.log("PUT OPTIONS:", options);
          if(!currentProjectId || !currentSchemaId || !currentSchemaVersion) return;
          await setSchemaFieldOptionsFn({
            project_id: currentProjectId,
            schema_id: currentSchemaId,
            version: currentSchemaVersion,
            field_id: fieldId,
            options: options
          });
        }
      })
    ]
  });

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if(e.ctrlKey && e.key === "z") fields.undo();
      if(e.ctrlKey && e.key === "y") fields.redo();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [fields])

  
  const handleDragStart = (ev: DragStartEvent) => {
    setActiveId(ev.active.id as string);
  }

  const handleDragEnd = (ev: DragEndEvent) => {
    const { active, over } = ev;
    if(over && active.id !== over.id) {
      fields.setItems(list => {
        const oldIndex = list.findIndex(i => i.key === active.id);
        const newIndex = list.findIndex(i => i.key === over.id);

        const moved = arrayMove(list, oldIndex, newIndex);

        return moved.map((item, index) => {
          if (item.position === index) return item;
          return { ...item, position: index };
        });
      });
    }
    setActiveId(null);
  }

  const handleDragCancel = () => {
    setActiveId(null);
  }

  const handleAddField = () => {
    const newField: VersionFieldResult = {
      key: `custom_field_${nextId}`,
      title: `Custom Field ${nextId}`,
      description: '',
      placeholder: '',
      object_id: undefined,
      id: 'Null',
      type: "string",
      owner: "user",
      mutable: true,
      required: false,
      position: fields.items.length,
      default_value: null,
      options: [],
      required_rules: [],
      visibility_rules: [],

    };
    fields.setItems(list => [...list, newField]);
    setNextId(prevId => prevId + 1);
  };

  const handleDeleteField = (fieldKey: string) => {
    fields.setItems(list => {
      const filteredList = list.filter(item => item.key !== fieldKey);
      return filteredList.map((item, index) => ({ ...item, position: index }));
    });
  };

  const handleUpdateField = (originalField: VersionFieldResult, updatedField: VersionFieldResult) => {
    fields.setItems(list =>
      list.map(item =>
        item.key === originalField.key ? updatedField : item
      )
    );
    if (editingField?.key === originalField.key) setEditingField(updatedField);
  };

  const activeItem = activeId ? fields.items.find(item => item.key === activeId) : null;

  const allFieldKeys = fields.items.map(item => item.key);

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
                      isSelected={editingField?.key === defaultEmailVersionField.key}
                    />
                    <FieldCard
                      key={defaultPasswordVersionField.key}
                      field={{ ...defaultPasswordVersionField }}
                      overwriteType="password"
                      isFixed={true}
                      isSelected={editingField?.key === defaultPasswordVersionField.key}
                    />
                    <SortableContext items={fields.items.map((item) => item.key)}>
                      {fields.items.map((item) => (
                        <FieldCard
                          key={item.key}
                          field={item}
                          onDelete={handleDeleteField}
                          onOpenEditPanel={handleOpenEditPanel}
                          isSelected={editingField?.key === item.key}
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
                isSelected={editingField?.key === defaultEmailVersionField.key}
              />
              <FieldCard
                key={defaultPasswordVersionField.key}
                field={{ ...defaultPasswordVersionField }}
                overwriteType="password"
                isFixed={true}
                isSelected={editingField?.key === defaultPasswordVersionField.key}
              />
              <SortableContext items={fields.items.map((item) => item.key)}>
                {fields.items.map((item) => (
                  <FieldCard
                    key={item.key}
                    field={item}
                    onDelete={handleDeleteField}
                    onOpenEditPanel={handleOpenEditPanel}
                    isSelected={editingField?.key === item.key}
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
          onClick={async() => {
            if(!currentProjectId || !currentSchemaId || !currentSchemaVersion) return;
            await fields.submit()
            await queryClient.invalidateQueries({
              queryKey: schemaVersionByIdQueryOptions(currentProjectId, currentSchemaId, currentSchemaVersion).queryKey
            });
            fields.syncWith(mappedFields);
          }}
          disabled={!fields.hasChanges || isVersionNull || fields.isSubmitting}
          leftIcon={<SaveAll className="w-4 h-4" />}
          value={isMobile ? '' : 'Save Fields'}
          variant="solid"
        />
        <PublishSchemaVersionButton
          isMobile={isMobile}
          hasChanges={!isVersionNull && schemaVData?.status === 'draft'}
        />
      </div>
    </main>
  );
}