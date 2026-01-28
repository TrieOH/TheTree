import { Link } from '@tanstack/react-router'
import { HeaderAction as HeaderActionType } from '../model/header.types'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import { ArrowLeft } from 'lucide-react'
import AuthButton from '@/features/auth/ui/AuthButton'

export function actionVisibilityClass(visibleOn?: 'mobile'|'desktop'|'always') {
  if (!visibleOn || visibleOn === 'always') return ''
  if (visibleOn === 'desktop') return 'hidden md:flex'
  return 'flex md:hidden' // mobile
}

export default function HeaderAction({ action }: { action: HeaderActionType }) {

  const cls = actionVisibilityClass(action.visibleOn)
  
  switch (action.type) {
    case 'link':
      return (
        <Link to={action.to} className={`${cls} custom-underline`}>
          {action.label}
        </Link>
      )

    case 'button':
      return (
        <div className={cls}>
          <ShadowButton value={action.label} leftIcon={action.icon} onClick={action.onClick} />
        </div>
      )

    case 'back':
      return (
        <div className={cls}>
          <ShadowButton leftIcon={<ArrowLeft size={18} />} onClick={() => history.back()} />
        </div>
      )

    case 'authButton':
      return (
        <div className={cls}>
          <AuthButton />
        </div>
      )

    default:
      return null
  }
}
