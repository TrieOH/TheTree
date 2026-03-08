import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { SignIn, SignUp } from '@trieoh/node-auth-sdk/react'
import { useState } from 'react';
import { motion } from "motion/react";
import { toast } from 'sonner';

export const Route = createFileRoute('/')({ component: App })

function App() {
  const [isLogin, setIsLogin] = useState(true);

  const navigate = useNavigate()

  const handleLoginSuccess = async () => {
    toast.success("Login successful!")
    // await navigate({ to: "" })
  }

  const handleSignUpSuccess = async () => {
    setIsLogin(true);
    toast.success("Account successfully created!")
  }

  const handleFailure = async (message: string, trace?: string[]) => {
    const traceMsg = trace?.join("\n").replaceAll("trace: ", "")
    toast.warning(`Auth Failed: ${message}`, { description: traceMsg })
  }

  return (
    <motion.main
      key={isLogin ? 'signin' : 'signup'}
      initial={{ opacity: 0, scale: 0.8, y: 5 }}
      animate={{ opacity: 1, scale: 1, y: 0 }}
      transition={{ duration: 0.4, ease: "easeOut" }}
      className='flex justify-center items-center py-2 h-screen'
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
