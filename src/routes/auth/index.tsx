import { requireGuest } from '@/hooks/auth-guard';
import { createFileRoute } from '@tanstack/react-router'
import { SignIn, SignUp } from '@trieoh/node-auth-sdk/react'
import { useState } from 'react'

export const Route = createFileRoute('/auth/')({
  beforeLoad: requireGuest,
  staticData: {
    components: {
      header: {
        test: false,
      }
    }
  },
  component: RouteComponent,
})

function RouteComponent() {
  const [isLogin, setIsLogin] = useState(true);
  return (
    <main className='flex justify-center items-center py-2'>
      {isLogin && <SignIn signUpRedirect={() => setIsLogin(false)} />}
      {!isLogin && <SignUp loginRedirect={() => setIsLogin(true)} />}
    </main>
  )
}
