import { motion } from 'motion/react'
import type { Mode } from '@/routes/index'

interface Props {
  current: Mode
  onChange: (mode: Mode) => void
}

export function ModeSelector({ current, onChange }: Props) {
  return (
    <div className="flex flex-col items-center gap-6 md:gap-8">
      {/* Toggle */}
      <div className="inline-flex p-1 bg-neutral-100 rounded-full">
        <button
          onClick={() => onChange('guest')}
          className="relative px-4 py-2 md:px-6 md:py-2.5 rounded-full text-xs md:text-sm font-medium transition-colors z-10"
        >
          {current === 'guest' && (
            <motion.div
              layoutId="activeTab"
              className="absolute inset-0 bg-white rounded-full shadow-sm"
              transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
            />
          )}
          <span className={`relative z-10 ${current === 'guest' ? 'text-neutral-900' : 'text-neutral-500'}`}>
            Quero Participar
          </span>
        </button>
        <button
          onClick={() => onChange('host')}
          className="relative px-4 py-2 md:px-6 md:py-2.5 rounded-full text-xs md:text-sm font-medium transition-colors z-10"
        >
          {current === 'host' && (
            <motion.div
              layoutId="activeTab"
              className="absolute inset-0 bg-neutral-900 rounded-full shadow-sm"
              transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
            />
          )}
          <span className={`relative z-10 ${current === 'host' ? 'text-white' : 'text-neutral-500'}`}>
            Quero Organizar
          </span>
        </button>
      </div>

      {/* Headline - usando suas frases exatas */}
      <h1 className="text-center px-2">
        {current === 'guest' ? (
          <span className="block text-3xl sm:text-4xl md:text-6xl lg:text-7xl font-semibold tracking-tight text-neutral-900 leading-[1.1]">
            Descubra eventos,<br />
            <span className="text-neutral-400">viva experiências.</span>
          </span>
        ) : (
          <span className="block text-3xl sm:text-4xl md:text-6xl lg:text-7xl font-semibold tracking-tight text-neutral-900 leading-[1.1]">
            Seus eventos,<br />
            <span className="text-neutral-400">sob controle total.</span>
          </span>
        )}
      </h1>
    </div>
  )
}