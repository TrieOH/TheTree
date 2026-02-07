import { useStore } from '@tanstack/react-store';
import { createCrudActions, type CrudStore, type CrudActions, type CrudState } from '../store/crudStore';

interface UseCrudStoreReturn<T extends { id: string }> extends CrudState<T> {
  actions: CrudActions<T>;
}

export function useCrudStore<T extends { id: string }>(
  store: CrudStore<T>
): UseCrudStoreReturn<T> {
  const state = useStore(store);
  const actions = createCrudActions(store);

  return {
    ...state,
    actions,
  };
}

interface UseCrudOperationsOptions<T extends { id: string }> {
  store: CrudStore<T>;
  onCreate?: (data: Omit<T, 'id'>) => Promise<void>;
  onUpdate?: (id: string, data: Partial<T>) => Promise<void>;
  onDelete?: (id: string) => Promise<void>;
  onSuccess?: () => void;
  onError?: (error: Error) => void;
  autoClose: boolean;
}

export function useCrudOperations<T extends { id: string }>(
  options: UseCrudOperationsOptions<T>
) {
  const { store, onCreate, onUpdate, onDelete, onSuccess, onError, autoClose = false } = options;
  const actions = createCrudActions(store);

  const executeOperation = async (
    operation: () => Promise<void>
  ) => {
    const state = store.state;
    actions.setLoading(true);
    try {
      await operation();
      onSuccess?.();
      if(autoClose || state.mode === "delete") actions.close();
    } catch (error) {
      onError?.(error as Error);
      console.error('CRUD Operation failed:', error);
    } finally {
      actions.setLoading(false);
    }
  };

  const handleSubmit = async (data: T) => {
    const state = store.state;
    const item = state.selectedItem;
    if (state.mode === "create" && onCreate) await executeOperation(() => onCreate(data));

    if(!item) return;

    if (state.mode === "edit" && onUpdate) await executeOperation(() => onUpdate(item.id, data));
    else if(state.mode === "delete" && onDelete) executeOperation(() => onDelete(item.id));
  };

  return {
    handleSubmit,
    actions,
  };
}