export interface StorageUploadRequest {
  filename: string;
  contentType: string;
  size: number;
}

export interface StorageUploadResponse {
  uploadUrl: string;
  key: string;
}

export interface StorageErrorResponse {
  error: string;
}

export interface StoragePreprocessResponse {
  approved: boolean;
  key?: string;
}
