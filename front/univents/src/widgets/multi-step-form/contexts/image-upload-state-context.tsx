import type { ReactNode } from "react";
import { createContext, useContext, useMemo, useState } from "react";
import type { ImageItem } from "../model/types";

type ImageUploadStateMap = Record<string, ImageItem[]>;

interface ImageUploadStateContextValue {
  getItems: (key: string) => ImageItem[] | undefined;
  setItems: (key: string, items: ImageItem[]) => void;
}

const ImageUploadStateContext =
  createContext<ImageUploadStateContextValue | null>(null);

export function ImageUploadStateProvider({
  children,
}: {
  children: ReactNode;
}) {
  const [state, setState] = useState<ImageUploadStateMap>({});

  const value = useMemo<ImageUploadStateContextValue>(
    () => ({
      getItems: (key) => state[key],
      setItems: (key, items) => {
        setState((current) =>
          current[key] === items ? current : { ...current, [key]: items },
        );
      },
    }),
    [state],
  );

  return (
    <ImageUploadStateContext.Provider value={value}>
      {children}
    </ImageUploadStateContext.Provider>
  );
}

export function useImageUploadState() {
  const context = useContext(ImageUploadStateContext);
  if (!context)
    throw new Error(
      "useImageUploadState must be used within ImageUploadStateProvider",
    );
  return context;
}
