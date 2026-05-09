import { ShadowButton } from "@/shared/ui/buttons/ShadowButton"
import { useNavigate, useRouteContext } from "@tanstack/react-router"
import { LogIn, User } from "lucide-react"

export default function AuthButton() {
  const navigate = useNavigate()
  const { auth } = useRouteContext({ from: '__root__' })

  return (
    <>
      {auth?.isAuthenticated === false ?
        <>
          <ShadowButton 
            value="Authenticate" 
            leftIcon={ <LogIn size={20}/> }
            className="xs:flex hidden"
            onClick={() => navigate({to: "/auth"})}
          />
          <ShadowButton
            leftIcon={ <LogIn size={16}/> }
            className="xs:hidden flex"
            onClick={() => navigate({to: "/auth"})}
          />
        </> : 
        <>
          <ShadowButton 
            value="Dashboard" 
            leftIcon={ <User size={20}/> }
            className="xs:flex hidden"
            onClick={() => navigate({to: "/projects"})}
          />
          <ShadowButton
            leftIcon={ <User size={16}/> }
            className="xs:hidden flex"
            onClick={() => navigate({to: "/projects"})}
          />
        </>
      }
    </>
  )
}