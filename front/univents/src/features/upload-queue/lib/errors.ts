import { StorageImageError } from '@/features/storage/api'
import type { UploadTaskError } from '../model/types'

export class UploadAssociationError extends Error {
  constructor(
    message: string,
    readonly options: {
      code?: string
      status?: number
      retryable?: boolean
    } = {},
  ) {
    super(message)
    this.name = 'UploadAssociationError'
  }
}

export function classifyUploadError(error: unknown): UploadTaskError {
  const occurredAt = Date.now()

  if (error instanceof StorageImageError) {
    if (error.code === 'MODERATION_REJECTED') {
      return {
        kind: 'moderation',
        code: error.code,
        message: 'A imagem não foi aprovada pela moderação.',
        retryable: false,
        requiresReplacement: true,
        occurredAt,
      }
    }

    if (error.status === 400 || error.status === 413 || error.status === 415) {
      return {
        kind: 'validation',
        code: error.code,
        message:
          error.message || 'O arquivo não atende aos requisitos do upload.',
        retryable: false,
        requiresReplacement: true,
        occurredAt,
      }
    }

    if (error.status === 401 || error.status === 403) {
      return {
        kind: 'authentication',
        code: error.code,
        message: 'A sessão precisa ser renovada para continuar o upload.',
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      }
    }

    if (
      error.status &&
      error.status < 500 &&
      error.status !== 408 &&
      error.status !== 429
    ) {
      return {
        kind: 'unknown',
        code: error.code,
        message: error.message || 'O servidor recusou o upload.',
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      }
    }

    return {
      kind: error.status && error.status >= 500 ? 'server' : 'network',
      code: error.code,
      message: 'Não foi possível enviar a imagem.',
      retryable: true,
      requiresReplacement: false,
      occurredAt,
    }
  }

  if (error instanceof UploadAssociationError) {
    const status = error.options.status
    if (status === 401 || status === 403) {
      return {
        kind: 'authentication',
        code: error.options.code ?? 'ASSOCIATION_AUTHENTICATION_FAILED',
        message: 'A sessão precisa ser renovada para continuar o upload.',
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      }
    }

    if (status === 404) {
      return {
        kind: 'not_found',
        code: error.options.code ?? 'UPLOAD_OWNER_NOT_FOUND',
        message: 'O registro relacionado a esta imagem não existe mais.',
        retryable: false,
        requiresReplacement: false,
        occurredAt,
      }
    }

    return {
      kind: status && status >= 500 ? 'server' : 'unknown',
      code: error.options.code ?? 'ASSOCIATION_FAILED',
      message:
        error.message || 'A imagem foi enviada, mas não pôde ser associada.',
      retryable:
        error.options.retryable ??
        (status === undefined ||
          status === 408 ||
          status === 429 ||
          status >= 500),
      requiresReplacement: false,
      occurredAt,
    }
  }

  if (error instanceof TypeError) {
    return {
      kind: 'network',
      code: 'NETWORK_ERROR',
      message: 'Não foi possível conectar ao servidor.',
      retryable: true,
      requiresReplacement: false,
      occurredAt,
    }
  }

  return {
    kind: 'unknown',
    code: 'UNKNOWN_UPLOAD_ERROR',
    message:
      error instanceof Error && error.message
        ? error.message
        : 'O upload falhou por um motivo inesperado.',
    retryable: false,
    requiresReplacement: false,
    occurredAt,
  }
}
