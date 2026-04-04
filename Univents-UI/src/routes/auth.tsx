import { createFileRoute } from '@tanstack/react-router'
import { SignIn, SignUp } from '@soramux/node-auth-sdk/react'
import { useState } from 'react';
import { motion } from "motion/react";
import { toast } from 'sonner';
import z from 'zod';
import { requireGuest } from '@/features/auths/lib/route-guard';
import { useAuthActions } from '@/features/auths/hooks/use-auth-actions';
import { cn } from '@/shared/lib/utils';

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(''),
})

export const Route = createFileRoute('/auth')({
  validateSearch: (search) => authSearchSchema.parse(search),
  beforeLoad: requireGuest,
  component: App,
})

function App() {
  const [isLogin, setIsLogin] = useState(true);

  const search = Route.useSearch();
  const { handleLoginSuccess } = useAuthActions();

  const onLoginSuccess = async () => {
    await handleLoginSuccess(search.redirect)
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  const handleSignUpSuccess = async () => {
    setIsLogin(true);
    toast.success("Account successfully created!")
  }
  // eslint-disable-next-line @typescript-eslint/require-await
  const handleFailure = async (message: string, trace?: string[]) => {
    const traceMsg = trace?.join("\n").replaceAll("trace: ", "")
    toast.warning(`Auth Failed: ${message}`, { description: traceMsg })
  }

  return (
    <main className={cn(
      "bg-background h-full text-foreground",
      "antialiased selection:bg-muted selection:text-foreground"
    )}>
      <motion.div
        key={isLogin ? 'signin' : 'signup'}
        initial={{ opacity: 0, scale: 0.8, y: 5 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        transition={{ duration: 0.4, ease: "easeOut" }}
        className="flex justify-center items-center py-2 min-h-screen"
      >
        {isLogin ? (
          <SignIn
            signUpRedirect={() => { setIsLogin(false); }}
            onSuccess={onLoginSuccess}
            onFailed={handleFailure}
          />
        ) : (
          <SignUp
            loginRedirect={() => { setIsLogin(true); }}
            onSuccess={handleSignUpSuccess}
            onFailed={handleFailure}
          />
        )}
      </motion.div>
    </main>
  )
}
