import { useTheme } from "next-themes";
import type { HTMLAttributes } from "react";
import { cn } from "@/shared/lib/utils";

export interface LogoProps extends HTMLAttributes<HTMLDivElement> {
  variant?: "complete" | "icon" | "responsive";
  theme?: "light" | "dark" | "default" | "auto";
  priority?: boolean;
  imgClassName?: string;
}

export function Logo({
  variant = "responsive",
  theme = "auto",
  priority = false,
  imgClassName,
  className,
  ...props
}: LogoProps) {
  const { resolvedTheme } = useTheme();

  let loading: "eager" | "lazy" = "lazy";
  if (priority) {
    loading = "eager";
  }

  let isDark = false;
  if (theme === "auto") {
    if (resolvedTheme === "dark") {
      isDark = true;
    }
  } else if (theme === "dark") {
    isDark = true;
  }

  const useDefault = theme === "default";

  const getSrc = (v: "complete" | "icon") => {
    if (useDefault) {
      return `/logo-${v}-default.svg`;
    }

    if (isDark) {
      return `/logo-${v}-dark.svg`;
    } else {
      return `/logo-${v}-light.svg`;
    }
  };

  const renderImage = (v: "complete" | "icon", variantClassName: string) => (
    <img
      src={getSrc(v)}
      alt={`Univents Logo ${v}`}
      loading={loading}
      className={cn(
        "w-full h-auto object-contain",
        variantClassName,
        imgClassName,
      )}
    />
  );

  return (
    <div
      className={cn(
        "relative flex items-center justify-center w-full h-full",
        className,
      )}
      {...props}
    >
      {variant === "responsive" && (
        <>
          {renderImage("complete", "hidden md:block")}
          {renderImage("icon", "block md:hidden")}
        </>
      )}
      {variant === "complete" && renderImage("complete", "block")}
      {variant === "icon" && renderImage("icon", "block")}
    </div>
  );
}
