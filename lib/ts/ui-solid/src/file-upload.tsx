import { For, Show, createSignal, createUniqueId, onCleanup } from "solid-js";
import Upload from "~icons/lucide/upload";
import FileIcon from "~icons/lucide/file";
import X from "~icons/lucide/x";
import RotateCcw from "~icons/lucide/rotate-ccw";
import CircleCheck from "~icons/lucide/circle-check";
import { Button } from "./button";
import { Loading } from "./loading";
import { cn } from "./lib/cn";
import { shapeClass, type Shape } from "./lib/shape";
export interface UploadFile {
  id: string;
  file: File;
  status: "ready" | "uploading" | "success" | "error";
  progress?: number;
  error?: string;
  previewUrl?: string;
}
export interface FileRejection {
  file: File;
  reason: string;
}
export interface FileUploadProps {
  accept?: string;
  multiple?: boolean;
  maxFiles?: number;
  maxSize?: number;
  disabled?: boolean;
  label?: string;
  description?: string;
  showPreview?: boolean;
  onUpload?: (
    file: File,
    context: { signal: AbortSignal; onProgress: (progress: number) => void },
  ) => Promise<void>;
  onFilesChange?: (files: readonly UploadFile[]) => void;
  onReject?: (rejections: readonly FileRejection[]) => void;
  shape?: Shape;
  class?: string;
}
const formatBytes = (bytes: number) =>
  bytes < 1024
    ? `${bytes} B`
    : bytes < 1024 * 1024
      ? `${(bytes / 1024).toFixed(1)} KB`
      : `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
const accepts = (file: File, accept?: string) =>
  !accept ||
  accept
    .split(",")
    .map((value) => value.trim().toLowerCase())
    .filter(Boolean)
    .some((value) =>
      value.startsWith(".")
        ? file.name.toLowerCase().endsWith(value)
        : value.endsWith("/*")
          ? file.type.toLowerCase().startsWith(value.slice(0, -1))
          : file.type.toLowerCase() === value,
    );
/** Owns local selection and cancellation; transport and real progress belong to the caller. */
export function FileUpload(props: FileUploadProps) {
  const id = createUniqueId();
  const [files, setFiles] = createSignal<readonly UploadFile[]>([]),
    [dragging, setDragging] = createSignal(false),
    [errors, setErrors] = createSignal<readonly FileRejection[]>([]);
  const tasks = new Map<string, AbortController>();
  // Uploads can start within the same input event; keep a synchronous snapshot
  // while Solid batches signal writes for rendering.
  let current: readonly UploadFile[] = [];
  let input!: HTMLInputElement,
    sequence = 0,
    depth = 0,
    disposed = false;
  const publish = (next: readonly UploadFile[]) => {
    if (disposed) return;
    current = next;
    setFiles(next);
    props.onFilesChange?.(next);
  };
  const update = (fileId: string, patch: Partial<UploadFile>) =>
    publish(
      current.map((file) =>
        file.id === fileId ? { ...file, ...patch } : file,
      ),
    );
  const upload = async (item: UploadFile) => {
    const handler = props.onUpload;
    if (!handler || props.disabled || tasks.has(item.id)) return;
    const controller = new AbortController();
    tasks.set(item.id, controller);
    update(item.id, {
      status: "uploading",
      error: undefined,
      progress: undefined,
    });
    try {
      await handler(item.file, {
        signal: controller.signal,
        onProgress: (value) => {
          if (!controller.signal.aborted && Number.isFinite(value))
            update(item.id, { progress: Math.max(0, Math.min(100, value)) });
        },
      });
      if (!controller.signal.aborted)
        update(item.id, { status: "success", progress: 100 });
    } catch (error) {
      if (!controller.signal.aborted)
        update(item.id, {
          status: "error",
          error:
            error instanceof Error
              ? error.message
              : "Não foi possível enviar o arquivo.",
        });
    } finally {
      tasks.delete(item.id);
    }
  };
  const add = (incoming: readonly File[]) => {
    if (props.disabled) return;
    const rejected: FileRejection[] = [],
      added: UploadFile[] = [];
    const limit =
      props.multiple === false ? 1 : Math.max(1, props.maxFiles ?? 5);
    for (const file of incoming) {
      let reason = "";
      if (!accepts(file, props.accept)) reason = "Formato não permitido.";
      else if (props.maxSize !== undefined && file.size > props.maxSize)
        reason = `O arquivo excede ${formatBytes(props.maxSize)}.`;
      else if (current.length + added.length >= limit)
        reason = `Limite de ${limit} arquivo${limit === 1 ? "" : "s"}.`;
      else if (
        [...current, ...added].some(
          (item) =>
            item.file.name === file.name &&
            item.file.size === file.size &&
            item.file.lastModified === file.lastModified,
        )
      )
        reason = "Este arquivo já foi selecionado.";
      if (reason) {
        rejected.push({ file, reason });
        continue;
      }
      added.push({
        id: `${id}-${++sequence}`,
        file,
        status: "ready",
        previewUrl:
          props.showPreview !== false &&
          /^image\/(png|jpeg|gif|webp|avif|bmp)$/.test(file.type) &&
          typeof URL.createObjectURL === "function"
            ? URL.createObjectURL(file)
            : undefined,
      });
    }
    setErrors(rejected);
    if (rejected.length) props.onReject?.(rejected);
    if (added.length) {
      publish([...current, ...added]);
      for (const item of added) void upload(item);
    }
  };
  const remove = (item: UploadFile) => {
    tasks.get(item.id)?.abort();
    tasks.delete(item.id);
    if (item.previewUrl) URL.revokeObjectURL(item.previewUrl);
    publish(current.filter((file) => file.id !== item.id));
  };
  onCleanup(() => {
    disposed = true;
    for (const task of tasks.values()) task.abort();
    for (const item of current)
      if (item.previewUrl) URL.revokeObjectURL(item.previewUrl);
  });
  return (
    <div class={cn("w-full min-w-0", props.class)}>
      <input
        ref={(element) => {
          input = element;
        }}
        id={id}
        type="file"
        class="sr-only"
        tabindex={-1}
        accept={props.accept}
        multiple={props.multiple !== false}
        disabled={props.disabled}
        aria-label={props.label ?? "Selecionar arquivos"}
        onChange={(e) => {
          add(Array.from(e.currentTarget.files ?? []));
          e.currentTarget.value = "";
        }}
      />
      <div
        onDragEnter={(e) => {
          e.preventDefault();
          depth++;
          if (!props.disabled) setDragging(true);
        }}
        onDragLeave={(e) => {
          e.preventDefault();
          if (--depth <= 0) {
            depth = 0;
            setDragging(false);
          }
        }}
        onDragOver={(e) => {
          e.preventDefault();
          e.dataTransfer!.dropEffect = props.disabled ? "none" : "copy";
        }}
        onDrop={(e) => {
          e.preventDefault();
          depth = 0;
          setDragging(false);
          add(Array.from(e.dataTransfer?.files ?? []));
        }}
        class={cn(
          "flex flex-col items-center gap-3 rounded-xl border border-dashed border-border bg-muted/15 p-5 text-center transition-colors",
          shapeClass(props.shape, "surface"),
          dragging() && "border-primary bg-accent",
          props.disabled && "opacity-50",
        )}
      >
        <span class="flex size-10 items-center justify-center rounded-xl border border-border bg-background text-muted-foreground">
          <Upload class="size-5" />
        </span>
        <div class="min-w-0">
          <p class="text-sm font-medium wrap-anywhere">
            {props.label ?? "Arraste arquivos para cá"}
          </p>
          <p class="mt-1 text-xs text-muted-foreground wrap-anywhere">
            {props.description ??
              "Ou selecione os arquivos no seu dispositivo."}
          </p>
          <Show when={props.maxSize}>
            <p class="mt-1 text-[11px] text-muted-foreground">
              Até {formatBytes(props.maxSize!)} por arquivo
            </p>
          </Show>
        </div>
        <Button
          size="sm"
          variant="outline"
          disabled={props.disabled}
          onClick={() => input.click()}
          aria-controls={id}
        >
          Selecionar arquivos
        </Button>
      </div>
      <Show when={errors().length}>
        <ul role="alert" class="mt-2 space-y-1 text-xs text-destructive">
          <For each={errors()}>
            {(error) => (
              <li class="wrap-anywhere">
                {error.file.name}: {error.reason}
              </li>
            )}
          </For>
        </ul>
      </Show>
      <ul aria-label="Arquivos selecionados" class="mt-3 space-y-2">
        <For each={files()}>
          {(item) => (
            <li class="flex min-w-0 items-center gap-3 rounded-lg border border-border p-3">
              <Show
                when={item.previewUrl && props.showPreview !== false}
                fallback={
                  <span class="flex size-10 shrink-0 items-center justify-center rounded-md bg-muted">
                    <FileIcon class="size-5 text-muted-foreground" />
                  </span>
                }
              >
                <img
                  src={item.previewUrl}
                  alt={`Prévia de ${item.file.name}`}
                  class="size-10 shrink-0 rounded-md object-cover"
                />
              </Show>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm" title={item.file.name}>
                  {item.file.name}
                </p>
                <p class="text-xs text-muted-foreground">
                  {formatBytes(item.file.size)} ·{" "}
                  {
                    {
                      ready: "Selecionado",
                      uploading: "Enviando",
                      success: "Enviado",
                      error: "Falha no envio",
                    }[item.status]
                  }
                </p>
                <Show when={item.status === "uploading"}>
                  <Show
                    when={item.progress !== undefined}
                    fallback={
                      <Loading
                        size="sm"
                        tone="muted"
                        label="Enviando arquivo"
                        class="mt-1 justify-start"
                      />
                    }
                  >
                    <div
                      role="progressbar"
                      aria-label={`Envio de ${item.file.name}`}
                      aria-valuemin={0}
                      aria-valuemax={100}
                      aria-valuenow={item.progress}
                      class="mt-2 h-1 overflow-hidden rounded-full bg-muted"
                    >
                      <div
                        class="h-full bg-primary transition-[width]"
                        style={{ width: `${item.progress}%` }}
                      />
                    </div>
                  </Show>
                </Show>
                <Show when={item.error}>
                  <p
                    role="alert"
                    class="mt-1 text-xs text-destructive wrap-anywhere"
                  >
                    {item.error}
                  </p>
                </Show>
              </div>
              <Show when={item.status === "success"}>
                <CircleCheck
                  class="size-4 shrink-0 text-emerald-600"
                  aria-label="Enviado"
                />
              </Show>
              <Show when={item.status === "error"}>
                <Button
                  size="icon"
                  variant="ghost"
                  disabled={props.disabled}
                  aria-label={`Tentar novamente ${item.file.name}`}
                  onClick={() => void upload(item)}
                >
                  <RotateCcw class="size-4" />
                </Button>
              </Show>
              <Button
                size="icon"
                variant="ghost"
                disabled={props.disabled}
                aria-label={`${item.status === "uploading" ? "Cancelar envio de" : "Remover"} ${item.file.name}`}
                onClick={() => remove(item)}
              >
                <X class="size-4" />
              </Button>
            </li>
          )}
        </For>
      </ul>
    </div>
  );
}
