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

export interface StorageModerateRequest {
  key: string;
}

export interface StorageModerateResponse {
  approved: boolean;
}

export interface StorageErrorResponse {
  error: string;
}
