import { createFileRoute, useNavigate, useRouter, useSearch } from '@tanstack/react-router'
import { SignIn, SignUp } from '@trieoh/node-auth-sdk/react'
import { motion } from "motion/react";
import { useState } from 'react'
import z from 'zod';
import { requireGuest } from '@/features/auth/lib/route-guard';

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

  const handleSuccess = async () => {
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
    }
    console.error("Auth Failed:")
  }

  const handleFailure = async (message: string) => {
    console.error("Auth Failed:", message)
  }

  return (
    <main className='flex justify-center items-center py-2'>
      <motion.div
        key={isLogin ? 'signin' : 'signup'}
        initial={{ opacity: 0, scale: 0.98, y: 5 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        transition={{ duration: 0.3, ease: "easeOut" }}
        className="w-full max-w-md"
      >
        {isLogin ? (
          <SignIn 
            signUpRedirect={() => setIsLogin(false)} 
            onSuccess={handleSuccess}
            onFailed={handleFailure}
          />
        ) : (
          <SignUp 
            loginRedirect={() => setIsLogin(true)} 
            onSuccess={async () => setIsLogin(true)}
            onFailed={handleFailure}
          />
        )}
      </motion.div>
    </main>
  )
}
