export type HeaderVariant = | 'landing' | 'projects' | 'auth' | 'none'
export type VisibleOn = 'mobile' | 'desktop' | 'fixed'
export type TitlePosition = 'left' | 'none'

export type HeaderAction =
  | { type: 'link'; label: string; to: string; visibleOn?: VisibleOn, collapseToMenu?: boolean }
  | { type: 'button'; label?: string; icon?: React.ReactNode; onClick: () => void; visibleOn?: VisibleOn, collapseToMenu?: boolean }
  | { type: 'authButton'; visibleOn?: VisibleOn, collapseToMenu?: boolean }
  | { type: 'back'; visibleOn?: VisibleOn, collapseToMenu?: boolean }

export interface HeaderConfigI {
  variant: HeaderVariant
  title?: string
  titlePosition?: TitlePosition
  leftActions?: HeaderAction[]
  centerActions?: HeaderAction[]
  rightActions?: HeaderAction[]
  disableMobileMenu?: boolean // if true never show the menu button
}