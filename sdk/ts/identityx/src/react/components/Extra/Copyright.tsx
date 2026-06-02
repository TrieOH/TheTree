export interface CopyrightProps {
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'
}

const sizeMap = {
  xs: 'text-trieoh-sm',
  sm: 'text-trieoh-base',
  md: 'text-trieoh-xl',
  lg: 'text-trieoh-2xl',
  xl: 'text-trieoh-3xl',
  '2xl': 'text-trieoh-6xl',
}

export function Copyright({ size }: CopyrightProps) {
  const sizeClass = size ? sizeMap[size] : 'text-trieoh-xl'

  return (
    <span className={`font-inter font-medium ${sizeClass}`}>
      © {new Date().getFullYear()} TrieOH
    </span>
  )
}