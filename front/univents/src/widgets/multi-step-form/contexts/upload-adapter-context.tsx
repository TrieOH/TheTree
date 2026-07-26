import type { ReactNode } from "react";
import { createContext, useContext, useMemo } from "react";
import { preprocessImageUpload } from "@/features/storage/api";
import type { ImageUploadAdapter } from "../model/types";

const ImageUploadAdapterContext = createContext<ImageUploadAdapter | null>(
  null,
);

export const realImageUploadAdapter: ImageUploadAdapter = {
  preprocess: async (file: File) => {
    const url = await preprocessImageUpload(file);
    return { url };
  },
};

export function ImageUploadProvider({
  adapter,
  children,
}: {
  /** Omit to use the default adapter. */
  adapter?: ImageUploadAdapter;
  children: ReactNode;
}) {
  const value = useMemo(() => adapter ?? realImageUploadAdapter, [adapter]);
  return (
    <ImageUploadAdapterContext.Provider value={value}>
      {children}
    </ImageUploadAdapterContext.Provider>
  );
}

export function useImageUploadAdapter(): ImageUploadAdapter {
  return useContext(ImageUploadAdapterContext) ?? realImageUploadAdapter;
}
