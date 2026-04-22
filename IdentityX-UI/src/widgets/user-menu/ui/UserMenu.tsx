import { useEffect, useRef, useState } from 'react'
import { Settings, Settings2 } from 'lucide-react'
import { useNavigate, useRouter } from '@tanstack/react-router'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'

import { BasicLogoutButton } from '@soramux/identityx-sdk-ts/react'
import { cn } from '@/shared/lib/utils'
import { toast } from 'sonner'

export default function UserMenu() {
  const router = useRouter()
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement | null>(null)
  const navigate = useNavigate()

  const handleLogoutSuccess = async () => {
    const auth = router.options.context.auth
    if(auth) {
      router.update({ 
        context: { 
          ...router.options.context, 
          auth: {...auth, isAuthenticated: false },
        },
      })

      router.options.context.queryClient.clear()
      await navigate({ to: '/', replace: true })
      toast.success("Logout successful!")
      router.options.context.queryClient.invalidateQueries();
    } else toast.error("Auth Initialization Failed")
  }

  const handleFailure = async (message: string, trace?: string[]) => {
    const traceMsg = trace?.join("\n").replaceAll("trace: ", "")
    toast.warning(`Auth Failed: ${message}`, {description: traceMsg})
  }

  useEffect(() => {
    if (!open) return

    function handleClickOutside(e: MouseEvent) {
      if (!ref.current) return
      if (!ref.current.contains(e.target as Node)) setOpen(false)
    }

    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') setOpen(false)
    }

    document.addEventListener('mousedown', handleClickOutside)
    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [open])

  return (
    <div ref={ref} className="relative">
      <ShadowButton
        leftIcon={<Settings2 size={18} />}
        onClick={() => setOpen(v => !v)}
      />

      {open && (
        <div className={cn(
          "absolute right-0 mt-2 w-44 rounded-md border border-border bg-popover shadow-lg",
          "flex flex-col justify-center items-center py-3 gap-2 px-1"
        )}>
          <button
            type="button"
            className={cn(
              "w-full flex items-end text-popover-foreground px-3 ",
              "gap-1 cursor-pointer text-sm font-medium",
              "hover:scale-105 active:scale-[0.98] duration-200"
            )}
            onClick={() => navigate({ to: '/' })}
          >
            <Settings size={24} />
            <span>Settings</span>
          </button>
          <hr className='w-full border-t border-border'/>
          <div className='w-full px-3'>
            <BasicLogoutButton onSuccess={handleLogoutSuccess} onFailed={handleFailure}/>
          </div>
        </div>
      )}
    </div>
  )
}
