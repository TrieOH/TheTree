import { Link } from '@tanstack/react-router'
import type { HeaderAction } from '../model/header.types'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import AuthButton from '@/features/auth/ui/AuthButton'
import BackButton from '@/features/navigation/ui/BackButton'

export default function HeaderActionRenderer({ action }: { action: HeaderAction }) {
  
  switch (action.type) {
    case 'link':
      return (
        <Link to={action.to} className="custom-underline">
          {action.label}
        </Link>
      )

    case 'button':
      return (
        <ShadowButton value={action.label} leftIcon={action.icon} onClick={action.onClick} />
      )

    case 'back':
      return (
        <BackButton value='Back' to={action.to} />
      )

    case 'authButton':
      return (
        <AuthButton />
      )

    default:
      return null
  }
}
