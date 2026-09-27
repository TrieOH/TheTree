import { StorageImageError } from "@/features/storage/api/index";
import {
  FileSizeExceededError,
  InvalidFileTypeError,
  StorageModerationError,
  StorageUploadNetworkError,
} from "@/features/storage/api";
import type { UploadTaskError } from "../model/types";

export class UploadAssociationError extends Error {
  constructor(
    message: string,
    readonly options: {
      code?: string;
      status?: number;
      retryable?: boolean;
    } = {},
  ) {
    super(message);
    this.name = "UploadAssociationError";
  }
}

type FailedApiResponse = {
  code?: number;
  error_id?: string;
  message?: string;
};

export function uploadAssociationErrorFromResponse(
  response: FailedApiResponse,
  fallbackMessage: string,
) {
  return new UploadAssociationError(response.message || fallbackMessage, {
    code: response.error_id ?? "ASSOCIATION_FAILED",
    status: response.code,
  });
}

export function classifyUploadError(error: unknown): UploadTaskError {
  const occurredAt = Date.now();

  if (error instanceof InvalidFileTypeError || error instanceof FileSizeExceededError) {
    return {
      kind: "validation",
      code: error._tag,
      message: error.message,
      retryable: false,
      requiresReplacement: true,
      occurredAt,
    };
  }

  if (error instanceof StorageModerationError) {
    return {
      kind: "moderation",
      code: error._tag,
      message: error.message,
      retryable: false,
      requiresReplacement: true,
      occurredAt,
    };
  }

  if (error instanceof StorageUploadNetworkError) {
    return {
      kind: "network",
      code: error._tag,
      message: error.message,
      retryable: true,
      requiresReplacement: false,
      occurredAt,
    };
  }

  if (error instanceof StorageImageError) {
    if (error.code === "MODERATION_REJECTED") {
      return {
        kind: "moderation",
        code: error.code,
        message: "A imagem não foi aprovada pela moderação.",
        retryable: false,
        requiresReplacement: true,
        occurredAt,
      };
    }

    if (error.status === 400 || error.status === 413 || error.status === 415) {
      return {
        kind: "validation",
        code: error.code,
        message:
          error.message || "O arquivo não atende aos requisitos do upload.",
        retryable: false,
        requiresReplacement: true,
        occurredAt,
      };
    }

    if (error.status === 401 || error.status === 403) {
      return {
        kind: "authentication",
        code: error.code,
        message: "A sessão precisa ser renovada para continuar o upload.",
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      };
    }

    if (
      error.status &&
      error.status < 500 &&
      error.status !== 408 &&
      error.status !== 429
    ) {
      return {
        kind: "unknown",
        code: error.code,
        message: error.message || "O servidor recusou o upload.",
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      };
    }

    return {
      kind: error.status && error.status >= 500 ? "server" : "network",
      code: error.code,
      message: "Não foi possível enviar a imagem.",
      retryable: true,
      requiresReplacement: false,
      occurredAt,
    };
  }

  if (error instanceof UploadAssociationError) {
    const status = error.options.status;
    if (status === 401 || status === 403) {
      return {
        kind: "authentication",
        code: error.options.code ?? "ASSOCIATION_AUTHENTICATION_FAILED",
        message: "A sessão precisa ser renovada para continuar o upload.",
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      };
    }

    if (status === 404) {
      return {
        kind: "not_found",
        code: error.options.code ?? "UPLOAD_OWNER_NOT_FOUND",
        message: "O registro relacionado a esta imagem não existe mais.",
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      };
    }

    return {
      kind: status && status >= 500 ? "server" : "unknown",
      code: error.options.code ?? "ASSOCIATION_FAILED",
      message:
        error.message || "A imagem foi enviada, mas não pôde ser associada.",
      retryable:
        error.options.retryable ??
        (status === undefined ||
          status === 408 ||
          status === 429 ||
          status >= 500),
      requiresReplacement: false,
      occurredAt,
    };
  }

  if (error instanceof TypeError) {
    return {
      kind: "network",
      code: "NETWORK_ERROR",
      message: "Não foi possível conectar ao servidor.",
      retryable: true,
      requiresReplacement: false,
      occurredAt,
    };
  }

  if (typeof error === "object" && error !== null) {
    const result = error as FailedApiResponse;
    if (result.code === 401 || result.code === 403) {
      return {
        kind: "authentication",
        code: result.error_id ?? "AUTHENTICATION_FAILED",
        message: "A sessão precisa ser renovada para continuar o upload.",
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      };
    }

    if (result.code === 404) {
      return {
        kind: "not_found",
        code: result.error_id ?? "RECORD_NOT_FOUND",
        message: "O recurso não foi encontrado.",
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      };
    }

    if (result.code === 409) {
      return {
        kind: "conflict",
        code: result.error_id ?? "CONFLICT",
        message: result.message || "Conflito de dados.",
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      };
    }

    if (result.code && result.code >= 500) {
      return {
        kind: "server",
        code: result.error_id ?? "SERVER_ERROR",
        message: result.message || "Ocorreu um erro no servidor.",
        retryable: true,
        requiresReplacement: false,
        occurredAt,
      };
    }
  }

  return {
    kind: "unknown",
    code: "UNEXPECTED_UPLOAD_ERROR",
    message:
      error instanceof Error ? error.message : "Ocorreu um erro inesperado.",
    retryable: true,
    requiresReplacement: false,
    occurredAt,
  };
}
