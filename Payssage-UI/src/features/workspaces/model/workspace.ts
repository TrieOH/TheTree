export interface Workspace {
  id: string
  name: string
  slug: string
  role: 'owner' | 'admin' | 'member'
  plan: 'free' | 'pro' | 'enterprise'
  status: 'active' | 'archived'
}

export const MOCK_WORKSPACES: Workspace[] = [
  {
    id: '1',
    name: 'Trie Payments',
    slug: 'trie-payments',
    role: 'owner',
    plan: 'enterprise',
    status: 'active',
  },
  {
    id: '2',
    name: 'Personal Project',
    slug: 'personal',
    role: 'owner',
    plan: 'free',
    status: 'active',
  },
  {
    id: '3',
    name: 'Design Studio',
    slug: 'studio',
    role: 'admin',
    plan: 'pro',
    status: 'active',
  },
]
