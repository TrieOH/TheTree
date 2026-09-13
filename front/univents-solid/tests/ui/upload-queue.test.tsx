import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import {
  classifyUploadError,
  getRetryDelay,
  uploadQueueConfig,
  UploadTaskCard,
  type UploadTask,
} from "@/features/upload-queue";
import {
  FileSizeExceededError,
  InvalidFileTypeError,
  StorageModerationError,
} from "@/features/storage/api";

afterEach(cleanup);

const mockTask: UploadTask = {
  id: "task-1",
  accountId: "acc-1",
  owner: { type: "event", id: "ev-1", label: "DevFest 2026" },
  mediaType: "banner",
  label: "Banner Principal",
  file: new Blob(["test"], { type: "image/png" }),
  fileName: "banner.png",
  contentType: "image/png",
  size: 1024 * 1024,
  stage: "upload",
  status: "queued",
  retryCount: 0,
  createdAt: Date.now(),
  updatedAt: Date.now(),
};

describe("upload-queue", () => {
  describe("classifyUploadError", () => {
    it("classifies file validation errors", () => {
      const err = new InvalidFileTypeError({
        fileType: "image/gif",
        allowedTypes: ["image/png", "image/jpeg", "image/webp"],
      });
      const classified = classifyUploadError(err);
      expect(classified.kind).toBe("validation");
      expect(classified.requiresReplacement).toBe(true);
      expect(classified.retryable).toBe(false);
    });

    it("classifies file size exceeded errors", () => {
      const err = new FileSizeExceededError({
        size: 15 * 1024 * 1024,
        maxSize: 10 * 1024 * 1024,
      });
      const classified = classifyUploadError(err);
      expect(classified.kind).toBe("validation");
      expect(classified.requiresReplacement).toBe(true);
    });

    it("classifies moderation errors", () => {
      const err = new StorageModerationError({
        reason: "nsfw",
        message: "Conteúdo impróprio detectado.",
      });
      const classified = classifyUploadError(err);
      expect(classified.kind).toBe("moderation");
      expect(classified.requiresReplacement).toBe(true);
      expect(classified.retryable).toBe(false);
    });

    it("classifies unexpected errors as retryable unknown errors", () => {
      const err = new Error("Network timeout");
      const classified = classifyUploadError(err);
      expect(classified.kind).toBe("unknown");
      expect(classified.retryable).toBe(true);
    });
  });

  describe("getRetryDelay", () => {
    it("calculates exponential delay with cap", () => {
      const delay = getRetryDelay(1);
      expect(delay).toBeGreaterThanOrEqual(uploadQueueConfig.baseDelayMs);
      expect(delay).toBeLessThanOrEqual(uploadQueueConfig.maxDelayMs * 1.5);
    });
  });

  describe("UploadTaskCard", () => {
    it("renders task metadata and status badge", () => {
      render(() => (
        <UploadTaskCard
          task={mockTask}
          highlighted={false}
          onRetry={vi.fn()}
          onReplace={vi.fn()}
          onRemove={vi.fn()}
        />
      ));

      expect(screen.getByText("Banner Principal")).toBeInTheDocument();
      expect(screen.getByText("DevFest 2026")).toBeInTheDocument();
      expect(screen.getByText("banner")).toBeInTheDocument();
      expect(screen.getByText("Na fila")).toBeInTheDocument();
    });

    it("renders retry button when task is failed and retryable", () => {
      const onRetry = vi.fn();
      const failedTask: UploadTask = {
        ...mockTask,
        status: "failed",
        error: {
          kind: "server",
          code: "SERVER_ERROR",
          message: "Falha temporária",
          retryable: true,
          requiresReplacement: false,
          occurredAt: Date.now(),
        },
      };

      render(() => (
        <UploadTaskCard
          task={failedTask}
          highlighted={false}
          onRetry={onRetry}
          onReplace={vi.fn()}
          onRemove={vi.fn()}
        />
      ));

      const retryBtn = screen.getByRole("button", { name: /tentar novamente/i });
      expect(retryBtn).toBeInTheDocument();
      fireEvent.click(retryBtn);
      expect(onRetry).toHaveBeenCalledTimes(1);
    });

    it("renders replace button when task is rejected", () => {
      const rejectedTask: UploadTask = {
        ...mockTask,
        status: "rejected",
        error: {
          kind: "moderation",
          code: "MODERATION_REJECTED",
          message: "Imagem recusada",
          retryable: false,
          requiresReplacement: true,
          occurredAt: Date.now(),
        },
      };

      render(() => (
        <UploadTaskCard
          task={rejectedTask}
          highlighted={false}
          onRetry={vi.fn()}
          onReplace={vi.fn()}
          onRemove={vi.fn()}
        />
      ));

      expect(
        screen.getByRole("button", { name: /trocar imagem/i }),
      ).toBeInTheDocument();
    });
  });
});
