export interface CopyrightProps {
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'
}

const sizeMap = {
  xs: 'text-xs',
  sm: 'text-sm',
  md: 'text-base',
  lg: 'text-lg',
  xl: 'text-xl',
  '2xl': 'text-2xl',
}

export function Copyright({ size }: CopyrightProps) {
  const sizeClass = size ? sizeMap[size] : 'text-base'

  return (
    <span className={`font-sans font-medium ${sizeClass}`}>
      © {new Date().getFullYear()} TrieOH
    </span>
  )
}