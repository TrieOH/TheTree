import type { UploadTask } from "../model/types";

/**
 * Resolves the corresponding file input ID for an upload task, if any.
 */
export function getAssociatedFileInputId(task: UploadTask): string | null {
  const handlerKey = task.association?.handlerKey;
  const field = task.association?.input?.field;

  if (handlerKey === "event-image") {
    if (field === "logo_url") return "event-logo-upload";
    if (field === "banner_url") return "event-banner-upload";
  }

  if (handlerKey === "edition-image") {
    const editionId = task.owner?.id;
    if (editionId && field === "logo_url") return `edition-${editionId}-logo-upload`;
    if (editionId && field === "banner_url") return `edition-${editionId}-banner-upload`;
  }

  return null;
}

/**
 * Prompts the user to pick a replacement image for a failed/invalid task.
 * Attempts to trigger the on-page file input, falling back to a dynamic file picker.
 */
export function promptImageReplacement(
  task: UploadTask,
  onFileChosen: (file: File) => void,
): void {
  if (typeof document === "undefined") return;

  const inputId = getAssociatedFileInputId(task);
  const existingInput = inputId ? document.getElementById(inputId) : null;

  if (existingInput instanceof HTMLInputElement) {
    existingInput.click();
    return;
  }

  const picker = document.createElement("input");
  picker.type = "file";
  picker.accept = "image/*";
  picker.style.display = "none";

  picker.onchange = () => {
    const file = picker.files?.[0];
    if (file) {
      onFileChosen(file);
    }
    picker.remove();
  };

  document.body.appendChild(picker);
  picker.click();
}
