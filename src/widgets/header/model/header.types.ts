export type HeaderVariant = | 'landing' | 'projects' | 'auth' | 'none'
export type VisibleOn = 'mobile' | 'desktop' | 'always'

export type HeaderAction =
  | { type: 'link'; label: string; to: string; visibleOn?: VisibleOn }
  | { type: 'button'; label?: string; icon?: React.ReactNode; onClick: () => void; visibleOn?: VisibleOn }
  | { type: 'authButton'; visibleOn?: VisibleOn }
  | { type: 'back'; visibleOn?: VisibleOn }

export interface HeaderConfigI {
  variant: HeaderVariant
  title?: string
  leftActions?: HeaderAction[]
  centerActions?: HeaderAction[]
  rightActions?: HeaderAction[]
}