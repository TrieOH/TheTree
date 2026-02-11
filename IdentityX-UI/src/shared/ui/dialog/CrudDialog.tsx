import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/shadcn/dialog';
import { AlertTriangle } from 'lucide-react';
import { ShadowButton } from '../buttons/ShadowButton';
import { createCrudActions, type CrudStore } from '@/shared/lib/store/crudStore';
import { useStore } from '@tanstack/react-store';
import { cn } from '@/shared/lib/utils';

interface CrudDialogProps<T extends { id: string }> {
  store: CrudStore<T>;
  title: string;
  description?: string;
  children?: React.ReactNode;
  onSubmit?: () => void; 
  formId: string;
}

export function CrudDialog<T extends { id: string }>({
  store,
  title,
  description,
  children,
  onSubmit,
  formId
}: CrudDialogProps<T>) {

  const state = useStore(store);
  const actions = createCrudActions(store);

  const config = {
    create: {
      title: `Create New ${title}`,
      description: description || `Enter the details to create a new ${title.toLowerCase()}.`,
      submitLabel: `Create ${title}`,
    },
    edit: {
      title: `Edit ${title}`,
      description: description || `Change ${title.toLowerCase()} data.`,
      submitLabel: `Update ${title}`,
    },
    delete: {
      title: `Delete ${title}`,
      description: description || `Are you sure you want to delete this ${title.toLowerCase()}?`,
      submitLabel: `Delete ${title}`,
    },
  };

  if (!state.mode) return null;
  const currentConfig = config[state.mode];

  return (
    <Dialog open={state.isOpen} onOpenChange={(open) => !open && actions.close()}>
      <DialogContent 
        className={cn(
          state.mode === 'delete' ? 'max-w-md' : 'max-w-lg',
          "min-w-[320px] w-11/12"
        )}
      >
        <DialogHeader>
          <DialogTitle>{currentConfig.title}</DialogTitle>
          <DialogDescription>{currentConfig.description}</DialogDescription>
        </DialogHeader>

        {state.mode === 'delete' ? (
          <div className="py-4">
            <div className="flex items-center gap-3 p-4 bg-destructive/10 rounded-lg text-destructive">
              <AlertTriangle className="h-5 w-5" />
              <p className="text-sm">This action cannot be undone.</p>
            </div>
          </div>
        ) : (
          children
        )}

        <DialogFooter showCloseButton closeButtonText='Cancel' isPerformingSubmit={state.isLoading}>
          <ShadowButton 
            type={state.mode === 'delete' ? "button" : "submit"}
            variant={state.mode === 'delete' ? 'destructive' : "accent-solid" }
            onClick={state.mode === 'delete' ? onSubmit : undefined}
            formId={state.mode === 'delete' ? undefined : formId}
            disabled={state.isLoading}
            className='justify-center px-4 font-normal sm:font-light text-sm'
            value={state.isLoading ? 'Submitting...' : currentConfig.submitLabel}
          /> 
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}