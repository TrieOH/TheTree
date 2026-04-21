import { createFileRoute, useNavigate, useRouter, useSearch } from '@tanstack/react-router'
import { SignIn, SignUp } from '@soramux/node-auth-sdk/react'
import { motion } from "motion/react";
import { useState } from 'react'
import z from 'zod';
import { requireGuest } from '@/features/auth/lib/route-guard';
import { toast } from 'sonner';

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(''),
})

export const Route = createFileRoute('/auth/')({
  validateSearch: (search) =>authSearchSchema.parse(search),
  beforeLoad: requireGuest,
  staticData: {
    components: {
      header: "auth"
    }
  },
  component: RouteComponent,
})

function RouteComponent() {
  const [isLogin, setIsLogin] = useState(true);

  const navigate = useNavigate()
  const router = useRouter()
  const search = useSearch({ from: '/auth/' })

  const handleLoginSuccess = async (message?: string) => {
    const auth = router.options.context.auth
    if(auth) {
      router.update({ 
        context: { 
          ...router.options.context, 
          auth: {...auth, isAuthenticated: true },
        },
      })
      const destination = search.redirect || '/projects'
      await navigate({ to: destination, replace: true })
      toast.success(message ?? "Login successful!")
      router.options.context.queryClient.invalidateQueries();
    } else toast.error("Auth Initialization Failed")
  }

  const handleSignUpSuccess = async (message?: string) => {
    setIsLogin(true);
    toast.success(message ?? "Account successfully created!")
  }

  const handleFailure = async (message: string, trace?: string[]) => {
    const traceMsg = trace?.join("\n").replaceAll("trace: ", "")
    toast.warning(`Auth Failed: ${message}`, {description: traceMsg})
  }

  return (
    <motion.main
      key={isLogin ? 'signin' : 'signup'}
      initial={{ opacity: 0, scale: 0.8, y: 5 }}
      animate={{ opacity: 1, scale: 1, y: 0 }}
      transition={{ duration: 0.4, ease: "easeOut" }}
      className='flex justify-center items-center py-2 h-(--screen--minus-header)'
    >
      {isLogin ? (
        <SignIn 
          signUpRedirect={() => setIsLogin(false)} 
          onSuccess={handleLoginSuccess}
          onFailed={handleFailure}
        />
      ) : (
        <SignUp 
          loginRedirect={() => setIsLogin(true)} 
          onSuccess={handleSignUpSuccess}
          onFailed={handleFailure}
        />
      )}
    </motion.main>
  )
}
