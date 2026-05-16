export interface InlineEditProps {
  value: string | null;
  onChange: (value: string) => void;
  isEditEnabled: boolean; // If enabled, render only the componet (without the edit)
  isEditing: boolean;           // Controlled by Parent
  onStartEdit: () => void;      // On Start
  onFinishEdit: () => void;     // On Finish (save)
  onCancelEdit?: () => void;    // Optional: On Cancel
  multiline?: boolean;
  className?: string;
  placeholder?: string;
}
