import { cn } from "../../../../utils/cn";

interface AuthLayoutProps {
  children: React.ReactNode;
  className?: string;
  backLink?: React.ReactNode;
}

export function AuthLayout({ children, className, backLink }: AuthLayoutProps) {
  return (
    <main className={cn(
      "bg-background h-full text-foreground min-h-screen relative overflow-hidden",
      "flex flex-col px-4",
      "antialiased selection:bg-primary/10 selection:text-primary",
      className
    )}>
      {/* Decorative background elements */}
      <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-3xl pointer-events-none" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-3xl pointer-events-none" />

      {/* Back link at the top, before centered content */}
      {backLink && (
        <div className="relative z-10 pt-4">
          {backLink}
        </div>
      )}

      {/* Centered content — takes remaining space */}
      <div className="relative z-10 flex-1 flex justify-center items-center">
        {children}
      </div>
    </main>
  );
}