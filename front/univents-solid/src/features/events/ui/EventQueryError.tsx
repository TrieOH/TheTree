export function EventQueryError(props: {
  message: string;
  onRetry: () => void;
}) {
  return (
    <div class="flex flex-col items-center justify-center gap-3 py-12 text-center">
      <p class="text-sm text-muted-foreground">{props.message}</p>
      <button
        type="button"
        class="rounded-md border border-border px-4 py-2 text-sm font-medium hover:bg-muted"
        onClick={props.onRetry}
      >
        Tentar novamente
      </button>
    </div>
  );
}
