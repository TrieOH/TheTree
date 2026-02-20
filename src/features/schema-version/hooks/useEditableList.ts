import { useCallback, useEffect, useRef, useState } from "react";

export type DiffResult<T> = {
  creates: T[];
  updates: { id: string; value: T }[];
  deletes: { id: string; value: T }[];
  currentMap: Map<string, T>;
  originalMap: Map<string, T>
};

export type CustomDiff<T> = (ctx: {
  getOriginalById: (id: string) => T | undefined;
  getCurrentById: (id: string) => T | undefined;
  diff: DiffResult<T>;
}) => Promise<void>;

type EditableListConfig<T> = {
  initial: T[];
  getId: (item: T) => string | undefined;
  isEqual: (a: T, b: T) => boolean;
  historyLimit?: number;

  onCreate?: (items: T[]) => Promise<void>; // For now i will just ignore the results
  onUpdate?: (items: { id: string; value: T }[]) => Promise<void>; // For now i will just ignore the results
  onDelete?: (items: { id: string; value: T }[]) => Promise<void>; // For now i will just ignore the results

  customDiffs?: CustomDiff<T>[];
};

export function useEditableList<T>({
  initial,
  getId,
  isEqual,
  onCreate,
  onUpdate,
  onDelete,
  historyLimit = 100,
  customDiffs
}: EditableListConfig<T>) {
  const itemsRef = useRef<T[]>(initial);
  const [, forceRender] = useState(0);
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
    const prevItems = itemsRef.current;
    const nextItems = updater(prevItems);

    itemsRef.current = nextItems;
    pushHistory(nextItems);
    forceRender(prev => prev + 1);
  }, [pushHistory]);

  const undo = useCallback(() => {
    if(indexRef.current <= 0) return;
    indexRef.current--;
    itemsRef.current = structuredClone(historyRef.current[indexRef.current]);
    forceRender(prev => prev + 1);
  }, []);

  const redo = useCallback(() => {
    if (indexRef.current >= historyRef.current.length - 1) return;
    indexRef.current++;
    itemsRef.current = structuredClone(historyRef.current[indexRef.current]);
    forceRender(prev => prev + 1);
  }, []);

  const canUndo = indexRef.current > 0;
  const canRedo = indexRef.current < historyRef.current.length - 1;

  useEffect(() => {
    itemsRef.current = initial;
    setOriginal(initial);
    historyRef.current = [structuredClone(initial)];
    indexRef.current = 0;
    forceRender(prev => prev + 1); 
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
    for (const c of itemsRef.current) { 
      const id = getId(c);
      if (id) currentMap.set(id, c);
    }

    // Create + Update
    for(const c of itemsRef.current){
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
    return { creates, updates, deletes, originalMap, currentMap };
  }, [original, getId, isEqual]);

  // SUBMIT
  const submit = useCallback(async () => {
    if(!onCreate && !onUpdate && !onDelete) return;
    setIsSubmitting(true);
    const diff = computeDiff();

    try {
      if (diff.deletes.length && onDelete) await onDelete(diff.deletes);
      if (diff.updates.length && onUpdate) await onUpdate(diff.updates);
      if (diff.creates.length && onCreate) await onCreate(diff.creates);

      for(const d of (customDiffs || [])) 
        await d({
          diff, 
          getCurrentById: id => diff.currentMap.get(id),
          getOriginalById: id => diff.originalMap.get(id),
        });
    } finally { setIsSubmitting(false); }
  }, [computeDiff, onCreate, onUpdate, onDelete, customDiffs]);

  const hasChanges = JSON.stringify(itemsRef.current) !== JSON.stringify(original);
  const syncWith = useCallback((newItems: T[]) => {
    itemsRef.current = structuredClone(newItems);
    setOriginal(structuredClone(newItems));
    historyRef.current = [structuredClone(newItems)];
    indexRef.current = 0;
    forceRender(prev => prev + 1);
  }, []);

  return {
    items: itemsRef.current, 
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