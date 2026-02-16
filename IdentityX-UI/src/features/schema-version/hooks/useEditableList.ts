import { useCallback, useEffect, useRef, useState } from "react";

export type DiffResult<T> = {
  creates: T[];
  updates: { id: string; value: T }[];
  deletes: { id: string; value: T }[];
};

type EditableListConfig<T> = {
  initial: T[];
  getId: (item: T) => string | undefined;
  isEqual: (a: T, b: T) => boolean;
  historyLimit?: number;

  onCreate?: (items: T[]) => Promise<void>; // For now i will just ignore the results
  onUpdate?: (items: { id: string; value: T }[]) => Promise<void>; // For now i will just ignore the results
  onDelete?: (items: { id: string; value: T }[]) => Promise<void>; // For now i will just ignore the results
};

export function useEditableList<T>({
  initial,
  getId,
  isEqual,
  onCreate,
  onUpdate,
  onDelete,
  historyLimit = 100,
}: EditableListConfig<T>) {
  const [items, setItems] = useState<T[]>(initial);
  const [original, setOriginal] = useState<T[]>(initial);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // HISTORY
  const historyRef = useRef<T[][]>([]);
  const indexRef = useRef(-1);

  const pushHistory = useCallback((snapshot: T[]) => {
    const hist = historyRef.current.slice(0, indexRef.current + 1);
    hist.push(structuredClone(snapshot));
    if(hist.length > historyLimit) hist.shift();

    historyRef.current = hist;
    indexRef.current = hist.length - 1;
  }, [historyLimit]);

  const update = useCallback((updater: (prev: T[]) => T[]) => {
    setItems(prev => {
      const next = updater(prev);
      pushHistory(next);
      return next;
    });
  }, [pushHistory]);

  const undo = () => {
    if(indexRef.current <= 0) return;
    indexRef.current--;
    setItems(structuredClone(historyRef.current[indexRef.current]));
  }

  const redo = () => {
    if (indexRef.current >= historyRef.current.length - 1) return;
    indexRef.current++;
    setItems(structuredClone(historyRef.current[indexRef.current]));
  };

  const canUndo = indexRef.current > 0;
  const canRedo = indexRef.current < historyRef.current.length - 1;

  useEffect(() => {
    setItems(initial);
    setOriginal(initial);
    historyRef.current = [structuredClone(initial)];
    indexRef.current = 0;
  }, [initial]);

  // DIFF
  const computeDiff = useCallback((): DiffResult<T> => {
    const creates: T[] = [];
    const updates: { id: string; value: T }[] = [];
    const deletes: { id: string; value: T }[] = [];

    const originalMap = new Map<string, T>();
    for(const o of original) {
      const id = getId(o);
      if(id) originalMap.set(id, o);
    }

    const currentMap = new Map<string, T>();
    for (const c of items) {
      const id = getId(c);
      if (id) currentMap.set(id, c);
    }

    // Create + Update
    for(const c of items){
      const id = getId(c);
      if(!id) {creates.push(c); continue;}
      const old = originalMap.get(id);
      if(!old) {creates.push(c); continue;}
      if(!isEqual(old, c)) updates.push({id, value: c});
    }

    // deletes
    for(const o of original) {
      const id = getId(o);
      if(!id) continue;
      if(!currentMap.has(id)) deletes.push({id, value: o});
    }
    return { creates, updates, deletes};
  }, [items, original, getId, isEqual]);

  // SUBMIT
  const submit = useCallback(async () => {
    if(!onCreate && !onUpdate && !onDelete) return;
    setIsSubmitting(true);
    const diff = computeDiff();

    try {
      // I need to call delete options, required rules and visibility, if their key matchs with the 
      // toDelete(diff.deletes) depends_on_field_key, i need to call this before delete and update
      if (diff.deletes.length && onDelete) await onDelete(diff.deletes);
      if (diff.updates.length && onUpdate) await onUpdate(diff.updates);
      if (diff.creates.length && onCreate) await onCreate(diff.creates);
    } finally {
      setIsSubmitting(false);
    }
  }, [computeDiff, onCreate, onUpdate, onDelete, items]);

  const hasChanges = JSON.stringify(items) !== JSON.stringify(original);

  const syncWith = useCallback((newItems: T[]) => {
    setItems(structuredClone(newItems));
    setOriginal(structuredClone(newItems));
    historyRef.current = [structuredClone(newItems)];
    indexRef.current = 0;
  }, []);

  return {
    items,
    setItems: update,

    undo,
    redo,
    canUndo,
    canRedo,

    submit,
    isSubmitting,
    hasChanges,

    syncWith
  };
}