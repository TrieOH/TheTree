import { Link, useParams } from '@tanstack/react-router'
import type { HeaderAction } from '../model/header.types'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import AuthButton from '@/features/auth/ui/AuthButton'
import BackButton from '@/features/navigation/ui/BackButton'
import CreateProjectButton from '@/features/project/ui/CreateProjectButton'
import SchemaVersionSelector from '@/features/schema-version/ui/SchemaVersionSelector'

export default function HeaderActionRenderer({ action }: { action: HeaderAction }) {
  const params = useParams({ strict: false });

  const resolveTo = (to?: string) => {
    if (!to) return to
    let resolved = to
    if (params.projectId) resolved = resolved.replace('$projectId', params.projectId)
    if (params.schemaId) resolved = resolved.replace('$schemaId', params.schemaId)
    return resolved
  }
  
  switch (action.type) {
    case 'link':
      return (
        <Link to={resolveTo(action.to)} className="custom-underline">
          {action.label}
        </Link>
      )

    case 'button':
      return (
        <ShadowButton value={action.label} leftIcon={action.icon} onClick={action.onClick} />
      )

    case 'createProject': return  <CreateProjectButton />

    case 'schemaVersionSelector': return <SchemaVersionSelector />

    case 'back': return <BackButton value='Back' to={resolveTo(action.to)} />

    case 'authButton': return <AuthButton />

    default:
      return null
  }
}
