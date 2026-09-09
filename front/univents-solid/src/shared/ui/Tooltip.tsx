import type { JSX } from '@solidjs/web';

interface TooltipProps {
  label: string;
  children: JSX.Element;
}

export function Tooltip(props: TooltipProps) {
  return (
    <div class="group relative">
      {props.children}
      <span
        role="tooltip"
        class="pointer-events-none absolute bottom-full left-1/2 z-50 mb-1.5 -translate-x-1/2 translate-y-1 rounded-md border border-border/50 bg-popover px-3 py-1.5 text-xs leading-none font-medium whitespace-nowrap text-popover-foreground opacity-0 shadow-lg shadow-black/10 transition-[opacity,transform] duration-150 group-hover:translate-y-0 group-hover:opacity-100 group-focus-within:translate-y-0 group-focus-within:opacity-100 max-md:hidden"
      >
        {props.label}
        <span
          aria-hidden="true"
          class="absolute top-[calc(100%-4px)] left-1/2 size-2 -translate-x-1/2 rotate-45 rounded-[1px] border-r border-b border-border/50 bg-popover"
        />
      </span>
    </div>
  );
}
