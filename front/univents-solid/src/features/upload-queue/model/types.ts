export type UploadTaskStage = "upload" | "association";

export type UploadTaskStatus =
  | "queued"
  | "uploading"
  | "associating"
  | "waiting_retry"
  | "completed"
  | "failed"
  | "rejected"
  | "paused";

export type UploadErrorKind =
  | "network"
  | "moderation"
  | "validation"
  | "authentication"
  | "not_found"
  | "conflict"
  | "server"
  | "configuration"
  | "unknown";

export interface UploadTaskError {
  kind: UploadErrorKind;
  code: string;
  message: string;
  retryable: boolean;
  requiresReplacement: boolean;
  occurredAt: number;
}

export interface UploadTaskOwner {
  type: string;
  id: string;
  label?: string;
}

export interface UploadTaskAssociation {
  handlerKey: string;
  input?: Record<string, unknown>;
}

export interface UploadTask {
  id: string;
  accountId: string;
  owner: UploadTaskOwner;
  mediaType: string;
  label: string;
  file: Blob;
  fileName: string;
  contentType: string;
  size: number;
  storagePath?: string;
  correctionPath?: string;
  association?: UploadTaskAssociation;
  uploadedUrl?: string;
  stage: UploadTaskStage;
  status: UploadTaskStatus;
  retryCount: number;
  nextAttemptAt?: number;
  error?: UploadTaskError;
  createdAt: number;
  updatedAt: number;
  completedAt?: number;
}

export interface EnqueueUploadInput {
  file: File;
  owner: UploadTaskOwner;
  mediaType: string;
  label?: string;
  storagePath?: string;
  correctionPath?: string;
  association?: UploadTaskAssociation;
}

export interface UploadQueueSnapshot {
  tasks: UploadTask[];
  initialized: boolean;
}
