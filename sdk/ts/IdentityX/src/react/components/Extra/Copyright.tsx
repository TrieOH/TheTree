export interface CopyrightProps {
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'
}

export function Copyright({ size }: CopyrightProps) {
  const sizeClass = size ? `trieoh-copyright--${size}` : ''

  return (
    <span className={`trieoh trieoh-copyright ${sizeClass}`}>
      © {new Date().getFullYear()} TrieOH
    </span>
  )
}