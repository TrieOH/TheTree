export type HeaderVariant = | 'landing' | 'projects' | 'projects/config' | 'auth' | 'none'
export type VisibleOn = 'mobile' | 'desktop' | 'fixed'
export type TitlePosition = 'left' | 'none'

export type HeaderAction =
  | { id: string; type: 'link'; label: string; to: string; visibleOn?: VisibleOn, collapseToMenu?: boolean }
  | { id: string; type: 'button'; label?: string; icon?: React.ReactNode; onClick: () => void; visibleOn?: VisibleOn, collapseToMenu?: boolean }
  | { id: string; type: 'authButton'; visibleOn?: VisibleOn, collapseToMenu?: boolean }
  | { id: string; type: 'createProject'; visibleOn?: VisibleOn, collapseToMenu?: boolean; }
  | { id: string; type: 'back'; visibleOn?: VisibleOn, collapseToMenu?: boolean, to?: string; }

export interface HeaderConfigI {
  variant: HeaderVariant
  title?: string
  titlePosition?: TitlePosition

  leftActions?: HeaderAction[]
  centerActions?: HeaderAction[]
  rightActions?: HeaderAction[]

  disableMobileMenu?: boolean // if true never show the menu button
  showUserMenu?: boolean
}