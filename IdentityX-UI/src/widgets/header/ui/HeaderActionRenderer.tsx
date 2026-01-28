import { Link } from '@tanstack/react-router'
import type { HeaderAction } from '../model/header.types'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import { ArrowLeft } from 'lucide-react'
import AuthButton from '@/features/auth/ui/AuthButton'

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
        <ShadowButton leftIcon={<ArrowLeft size={18} />} onClick={() => history.back()} />
      )

    case 'authButton':
      return (
        <AuthButton />
      )

    default:
      return null
  }
}
