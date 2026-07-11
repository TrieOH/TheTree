import { useCallback, useRef } from "react";
import type { ImageFieldChange } from "../model/types";

const emptyChange: ImageFieldChange = { added: [], removed: [] };

/**
 * Collects the `onTrackingChange` callbacks from every image/gallery
 * field in a form, keyed by whatever string you pass to `track(key)`
 * (typically the field name, e.g. "logo_url", "gallery_urls").
 *
 * Read `getChanges(key)` inside your `onSubmit` to know which images
 * were added (the form will flush the server-side preprocess step
 * right before submit) and which existing images were removed (need a
 * delete call via your own endpoint).
 */
export function useImageFieldTracking() {
  const changesRef = useRef<Record<string, ImageFieldChange>>({});

  const track = useCallback(
    (key: string) => (change: ImageFieldChange) => {
      changesRef.current[key] = change;
    },
    [],
  );

  const getChanges = useCallback((key: string): ImageFieldChange => {
    return changesRef.current[key] ?? emptyChange;
  }, []);

  const reset = useCallback(() => {
    changesRef.current = {};
  }, []);

  return { track, getChanges, reset };
}
