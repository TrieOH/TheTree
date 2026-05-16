import { motion } from 'motion/react'
import { Link } from '@tanstack/react-router'
import { buttonVariants } from '@/shared/ui/shadcn/button'
import { cn } from '@/shared/lib/utils'

export default function NotFound() {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center min-h-screen px-4 text-center",
        "bg-background text-foreground"
      )}
    >
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="max-w-md w-full min-h-[80vh] space-y-8"
      >
        <div className="relative group">
          <motion.div
            initial={{ scale: 0.95, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ delay: 0.2, duration: 0.8 }}
            className="aspect-video rounded-3xl overflow-hidden bg-muted relative"
          >
            <img
              src="/images/lagoon-5.svg"
              alt="404 Lagoon"
              className="w-full h-full object-cover opacity-80 dark:opacity-60 transition-all duration-700 group-hover:scale-105"
            />
            <div className="absolute inset-0 flex items-center justify-center bg-background/20 backdrop-blur-[2px]">
              <span className="text-8xl font-black tracking-tighter select-none opacity-20">
                404
              </span>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.5 }}
            className="absolute -bottom-4 -right-4 w-24 h-24 bg-accent/20 rounded-full blur-2xl -z-10"
          />
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.6 }}
            className="absolute -top-4 -left-4 w-32 h-32 bg-primary/10 rounded-full blur-3xl -z-10"
          />
        </div>

        <div className="space-y-3">
          <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground font-heading">
            Ops! Página não encontrada
          </h1>
          <p className="text-muted-foreground text-balance max-w-xs mx-auto">
            Parece que o evento que você está procurando ainda não começou ou o link está incorreto.
          </p>
        </div>

        <div className="flex flex-col sm:flex-row gap-3 justify-center pt-4">
          <Link
            to="/"
            className={cn(buttonVariants({ size: 'xl' }), "rounded-full px-8")}
          >
            Voltar ao Início
          </Link>
          <Link
            to="/events"
            className={cn(buttonVariants({ variant: 'outline', size: 'xl' }), "rounded-full px-8")}
          >
            Explorar Eventos
          </Link>
        </div>
      </motion.div>
    </div>
  )
}
