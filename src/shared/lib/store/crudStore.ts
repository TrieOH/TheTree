import { Store } from '@tanstack/react-store';
export type CrudMode = 'create' | 'edit' | 'delete' | null;

export interface CrudState<T extends { id: string }> {
  mode: CrudMode;
  selectedItem: T | null;
  isLoading: boolean;
  isOpen: boolean;
  formData: Partial<T> | null;
}

export interface CrudActions<T extends { id: string }> {
  openCreate: () => void;
  openEdit: (item: T) => void;
  openDelete: (item: T) => void;
  close: () => void;
  setLoading: (loading: boolean) => void;
  setFormData: (data: Partial<T> | null) => void;
}

export type CrudStore<T extends { id: string }> = Store<CrudState<T>>;

export function createCrudStore<T extends { id: string }>() {
  return new Store<CrudState<T>>({
    mode: null,
    selectedItem: null,
    isLoading: false,
    isOpen: false,
    formData: null,
  });
}

export function createCrudActions<
  T extends { id: string },
  S extends CrudState<T>
>(
  store: Store<S>
): CrudActions<T> {
  return {
    openCreate: () => {
      store.setState((state) => ({
        ...state,
        mode: 'create',
        selectedItem: null,
        isOpen: true,
        formData: {},
      }));
    },

    openEdit: (item: T) => {
      store.setState((state) => ({
        ...state,
        mode: 'edit',
        selectedItem: item,
        isOpen: true,
        formData: item,
      }));
    },

    openDelete: (item: T) => {
      store.setState((state) => ({
        ...state,
        mode: 'delete',
        selectedItem: item,
        isOpen: true,
        formData: null,
      }));
    },

    close: () => {
      store.setState((state) => ({
        ...state,
        mode: null,
        selectedItem: null,
        isOpen: false,
        isLoading: false,
        formData: null,
      }));
    },

    setLoading: (loading: boolean) => {
      store.setState((state) => ({
        ...state,
        isLoading: loading,
      }));
    },

    setFormData: (data: Partial<T> | null) => {
      store.setState((state) => ({
        ...state,
        formData: data,
      }));
    },
  };
}