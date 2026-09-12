import type { SVGProps } from 'react'

const strokeProps = {
  fill: 'none' as const,
  stroke: 'currentColor',
  strokeWidth: 1.75,
  strokeLinecap: 'round' as const,
  strokeLinejoin: 'round' as const,
}

export function LineIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <path d="M2.5 14.5 7 9l3.5 3L17.5 5" />
    </svg>
  )
}

export function AreaIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <path d="M2.5 14.5 7 9l3.5 3L17.5 5V16H2.5Z" fill="currentColor" fillOpacity={0.18} stroke="none" />
      <path d="M2.5 14.5 7 9l3.5 3L17.5 5" />
    </svg>
  )
}

export function BarIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <path d="M4 16V10M10 16V4M16 16v-7" />
    </svg>
  )
}

export function ScatterIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <circle cx="5" cy="14" r="1.4" fill="currentColor" stroke="none" />
      <circle cx="10" cy="7" r="1.4" fill="currentColor" stroke="none" />
      <circle cx="14" cy="12" r="1.4" fill="currentColor" stroke="none" />
      <circle cx="16.5" cy="5.5" r="1.4" fill="currentColor" stroke="none" />
    </svg>
  )
}

export function SearchIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <circle cx="8.5" cy="8.5" r="5.5" />
      <path d="m17 17-4-4" />
    </svg>
  )
}

export function CloseIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <path d="m5 5 10 10M15 5 5 15" />
    </svg>
  )
}

export function CalendarIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <rect x="3" y="4.5" width="14" height="12" rx="2" />
      <path d="M3 8.5h14M7 2.5v3M13 2.5v3" />
    </svg>
  )
}

export function ChevronLeftIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <path d="M12.5 4.5 7 10l5.5 5.5" />
    </svg>
  )
}

export function ChevronRightIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <path d="m7.5 4.5 5.5 5.5-5.5 5.5" />
    </svg>
  )
}

export function ChevronDownIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 20 20" {...strokeProps} {...props}>
      <path d="m4.5 7.5 5.5 5.5 5.5-5.5" />
    </svg>
  )
}