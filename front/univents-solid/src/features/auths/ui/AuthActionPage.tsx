import type { JSX } from "@solidjs/web";
import type { ParentProps } from "solid-js";
import { AuthLayout } from "@trieoh/identityx-sdk-ts-solid";
import Logo from "@/shared/ui/Logo";

export function AuthActionPage(
  props: ParentProps<{ title?: string; description?: string }>,
): JSX.Element {
  return (
    <div class="relative [&>main]:py-16 [&>main]:pt-52 md:[&>main]:pt-56 [&>main]:pb-28">
      <div class="absolute left-1/2 top-16 md:top-20 z-20 w-32 md:w-40 -translate-x-1/2">
        <Logo
          variant="complete"
          priority
          imgClassName="h-auto max-h-16 md:max-h-20"
        />
      </div>

      <AuthLayout>
        <div class="w-full max-w-md z-10 flex flex-col">
          <div class="text-center mb-6">
            {props.title && (
              <h1 class="font-heading text-3xl font-bold tracking-tight mb-2">
                {props.title}
              </h1>
            )}
            {props.description && (
              <p class="text-muted-foreground text-sm mb-4">
                {props.description}
              </p>
            )}
          </div>
          {props.children}
        </div>
      </AuthLayout>
    </div>
  );
}
