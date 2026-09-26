export function ProfileEmptyState(props: { message: string }) {
  return (
    <p class="rounded-md border border-dashed border-border p-6 text-center text-sm text-muted-foreground">
      {props.message}
    </p>
  );
}
