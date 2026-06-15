import { cn } from "../../../../utils/cn";

interface AuthLayoutProps {
  children: React.ReactNode;
  className?: string;
}

export function AuthLayout({ children, className }: AuthLayoutProps) {
  return (
    <main className={cn(
      "bg-background h-full text-foreground min-h-screen relative overflow-hidden",
      "flex justify-center items-center px-4",
      "antialiased selection:bg-primary/10 selection:text-primary",
      className
    )}>
      {/* Decorative background elements */}
      <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-3xl pointer-events-none" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-3xl pointer-events-none" />

      {children}
    </main>
  );
}