import type { ParentProps } from "solid-js";
export interface AuthLayoutProps extends ParentProps {
  class?: string;
  backLink?: ParentProps["children"];
}
export function AuthLayout(props: AuthLayoutProps) {
  return (
    <main
      class={`bg-background h-full text-foreground min-h-screen relative overflow-hidden flex flex-col px-4 antialiased ${props.class ?? ""}`}
    >
      <div class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-3xl pointer-events-none" />
      <div class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-3xl pointer-events-none" />
      {props.backLink && <div class="relative z-10 pt-4">{props.backLink}</div>}
      <div class="relative z-10 flex-1 flex justify-center items-center">
        {props.children}
      </div>
    </main>
  );
}
