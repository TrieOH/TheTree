import { useCallback, useEffect, useMemo, useRef } from "react";
import { useImageUploadState } from "../contexts/image-upload-state-context";
import { useImageUploadAdapter } from "../contexts/upload-adapter-context";
import type {
  ImageFieldChange,
  ImageItem,
  UploadedImage,
} from "../model/types";
import { registerImageUploadTask } from "./use-image-upload-queue";

let localIdCounter = 0;
function createLocalId(): string {
  localIdCounter += 1;
  return `img_${Date.now()}_${localIdCounter}`;
}

export interface UseImageUploadFieldOptions {
  fieldKey: string;
  /** URLs already present when the field mounted (edit mode). Captured
   * once by the caller via lazy `useState` — passing a fresh array
   * every render would keep resetting the tracked state. */
  initialUrls: string[];
  maxItems: number;
  accept?: string;
  maxSizeMB?: number;
  /** Called with the current list of *final* URLs (existing + newly
   * uploaded, excluding anything removed/still-processing) — this is
   * what you wire into `form.setValue`. */
  onValueChange: (urls: string[]) => void;
  /** Called with the running "added this session" / "removed this
   * session" sets — this is what you use at submit time to call your
   * own save/remove endpoints. */
  onTrackingChange?: (change: ImageFieldChange) => void;
}

export interface UseImageUploadFieldResult {
  items: ImageItem[];
  addFiles: (files: FileList | File[]) => void;
  removeItem: (id: string) => void;
  canAddMore: boolean;
}

export function useImageUploadField({
  fieldKey,
  initialUrls,
  maxItems,
  maxSizeMB,
  onValueChange,
  onTrackingChange,
}: UseImageUploadFieldOptions): UseImageUploadFieldResult {
  const adapter = useImageUploadAdapter();
  const { getItems, setItems: setStoredItems } = useImageUploadState();
  const itemsRef = useRef<ImageItem[]>([]);

  const initialItems = useMemo(
    (): ImageItem[] =>
      initialUrls.map((url) => ({
        id: createLocalId(),
        url,
        status: "existing",
        isExisting: true,
      })),
    [initialUrls],
  );
  const items = getItems(fieldKey) ?? initialItems;

  const removedRef = useRef<UploadedImage[]>([]);

  useEffect(() => {
    itemsRef.current = items;
  }, [items]);

  useEffect(() => {
    if (getItems(fieldKey)) return;
    setStoredItems(fieldKey, initialItems);
  }, [fieldKey, getItems, initialItems, setStoredItems]);

  const emitChanges = useCallback(
    (nextItems: ImageItem[]) => {
      const finalUrls = nextItems
        .filter(
          (item) =>
            item.status === "existing" ||
            (item.status === "ready" && !item.file),
        )
        .map((item) => item.url);
      onValueChange(finalUrls);

      const added = nextItems
        .filter(
          (item) => !item.isExisting && item.status === "ready" && !item.file,
        )
        .map((item) => ({ id: item.id, url: item.url }));
      onTrackingChange?.({ added, removed: removedRef.current });
    },
    [onValueChange, onTrackingChange],
  );

  const updateItem = useCallback(
    (id: string, patch: Partial<ImageItem>) => {
      const next = itemsRef.current.map((item) =>
        item.id === id ? { ...item, ...patch } : item,
      );
      setStoredItems(fieldKey, next);
      emitChanges(next);
    },
    [emitChanges, fieldKey, setStoredItems],
  );

  const flushPendingUploads = useCallback(async () => {
    const snapshot = itemsRef.current;
    const itemsToProcess = snapshot.filter((item) => item.file);
    let hadFailure = false;
    let failureMessage: string | null = null;

    for (const item of itemsToProcess) {
      if (!item.file) continue;

      if (maxSizeMB && item.file.size > maxSizeMB * 1024 * 1024) {
        const message = `Arquivo maior que ${maxSizeMB}MB`;
        updateItem(item.id, { status: "error", errorMessage: message });
        failureMessage ??= message;
        hadFailure = true;
        continue;
      }

      try {
        updateItem(item.id, { status: "processing" });
        const { url } = await adapter.preprocess(item.file);
        updateItem(item.id, { status: "ready", url, file: undefined });
      } catch (error) {
        const message =
          error instanceof Error ? error.message : "Falha no upload";
        updateItem(item.id, { status: "error", errorMessage: message });
        failureMessage ??= message;
        hadFailure = true;
      }
    }

    if (hadFailure)
      throw new Error(failureMessage ?? "One or more image uploads failed");
  }, [adapter, maxSizeMB, updateItem]);

  useEffect(() => {
    registerImageUploadTask(fieldKey, flushPendingUploads);
  }, [fieldKey, flushPendingUploads]);

  const addFiles = useCallback(
    (files: FileList | File[]) => {
      const incoming = Array.from(files);

      const current = itemsRef.current;
      const remainingSlots = Math.max(maxItems - current.length, 0);
      const accepted = incoming.slice(0, remainingSlots);

      const newItems: ImageItem[] = accepted.map((file) => ({
        id: createLocalId(),
        url: URL.createObjectURL(file),
        status: "ready" as const,
        file,
        isExisting: false,
      }));

      const next = [...current, ...newItems];
      setStoredItems(fieldKey, next);
      emitChanges(next);
    },
    [maxItems, emitChanges, fieldKey, setStoredItems],
  );

  const removeItem = useCallback(
    (id: string) => {
      const current = itemsRef.current;
      const target = current.find((item) => item.id === id);
      if (target?.isExisting)
        removedRef.current = [
          ...removedRef.current,
          { id: target.id, url: target.url },
        ];

      const next = current.filter((item) => item.id !== id);
      setStoredItems(fieldKey, next);
      emitChanges(next);
    },
    [emitChanges, fieldKey, setStoredItems],
  );

  return { items, addFiles, removeItem, canAddMore: items.length < maxItems };
}
