import { createFileRoute, redirect } from '@tanstack/react-router'
import { SignIn, SignUp } from '@trieoh/node-auth-sdk/react'
import { useState } from 'react'

export const Route = createFileRoute('/auth/')({
  beforeLoad: async ({context}) => {
    console.log(context.auth?.isAuthenticated);
    if(context.auth?.isAuthenticated) throw redirect({to: "/projects"})
  },
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
