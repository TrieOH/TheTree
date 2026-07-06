export function OAuthDivider() {
  return (
    <div className="flex items-center gap-3 py-1">
      <div className="flex-1 h-px bg-border" />
      <span className="text-xs font-medium text-muted-foreground uppercase tracking-widest select-none">
        ou
      </span>
      <div className="flex-1 h-px bg-border" />
    </div>
  );
}
