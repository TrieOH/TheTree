export interface StorageUploadRequest {
  filename: string;
  contentType: string;
  size: number;
}

export interface StorageUploadResponse {
  uploadUrl: string;
  key: string;
  publicUrl: string;
}

export interface StorageErrorResponse {
  error: string;
}

export interface StoragePreprocessResponse {
  approved: boolean;
  publicUrl?: string;
}
