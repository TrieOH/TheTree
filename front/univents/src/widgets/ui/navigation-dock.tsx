import { useRef, memo, useMemo } from 'react';
import { motion, useMotionValue, useSpring, useTransform } from 'motion/react';
import { Home, User, Calendar, LogIn, LogOut, Store, Activity } from 'lucide-react';
import { useLocation, useNavigate } from '@tanstack/react-router';
import type { LucideIcon } from 'lucide-react';
import { Tooltip, TooltipTrigger, TooltipContent } from '@/shared/ui/shadcn/tooltip';
import { cn } from '@/shared/lib/utils';
import { useAuthActions } from '@/features/auths/hooks/use-auth-actions';

interface NavItemType {
  id: string;
  label: string;
  icon: LucideIcon | React.ComponentType;
  href?: string;
  authRequired?: boolean;
  hideIfAuthenticated?: boolean;
  onClick?: () => void | Promise<void>;
}

interface NavigationDockProps {
  className?: string;
}

const UVIcon = () => (
  <span
    className={cn(
      "font-heading font-semibold text-lg flex",
      "items-center justify-center w-full h-full"
    )}
  >
    UV
  </span>
);

/**
 * Ordered by specificity: most specific matches first, default at the end.
 */
const navConfigs = (actions: { logout: () => Promise<void> }, isAdmin: boolean) => [
  {
    id: 'edition-context',
    match: (parts: string[]) => {
      // /events/:eventId/editions/:editionId
      if (parts[0] === 'events' && parts[2] === 'editions' && parts[3]) return true;
      // /admin/events/:eventId/editions/:editionId
      if (parts[0] === 'admin' && parts[1] === 'events' && parts[3] === 'editions' && parts[4])
        return true;
      return false;
    },
    getItems: (parts: string[]): NavItemType[] => {
      const isAdminRoute = parts[0] === 'admin';
      const eventId = isAdminRoute ? parts[2] : parts[1];
      const editionId = isAdminRoute ? parts[4] : parts[3];

      const eventBase = `/events/${eventId}`;
      const editionBase = `${eventBase}/editions/${editionId}`;
      const adminEditionBase = `/admin/events/${eventId}/editions/${editionId}`;

      return [
        { id: 'back-home', label: 'Univents', icon: UVIcon, href: '/' },
        { id: 'edition-home', label: 'Edição', icon: Home, href: editionBase },
        {
          id: 'edition-activities',
          label: 'Atividades',
          icon: Activity,
          href: `${isAdminRoute ? adminEditionBase : editionBase}/activities`
        },
        {
          id: 'edition-products',
          label: 'Produtos',
          icon: Store,
          href: `${isAdminRoute ? adminEditionBase : editionBase}/products`
        },
        { id: 'edition-profile', label: 'Perfil', icon: User, href: `${editionBase}/profile`, authRequired: true },
        { id: 'edition-login', label: 'Entrar', icon: LogIn, href: '/auth', hideIfAuthenticated: true },
      ];
    }
  },
  {
    id: 'event-context',
    // Matches /events/$eventId/... or /admin/events/$eventId/...
    match: (parts: string[]) => {
      if (parts[0] === 'events' && parts[1] && parts[1] !== 'index') return true;
      if (parts[0] === 'admin' && parts[1] === 'events' && parts[2]) return true;
      return false;
    },
    getItems: (parts: string[]): NavItemType[] => {
      const isAdminRoute = parts[0] === 'admin';
      const eventId = isAdminRoute ? parts[2] : parts[1];
      const eventBase = `/events/${eventId}`;
      const adminBase = `/admin/events/${eventId}`;

      return [
        { id: 'back-home', label: 'Univents', icon: UVIcon, href: '/' },
        { id: 'event-home', label: 'Evento', icon: Home, href: eventBase },
        {
          id: 'event-editions',
          label: 'Edições',
          icon: Calendar,
          href: `${isAdminRoute ? adminBase : eventBase}/editions`
        },
        { id: 'event-profile', label: 'Perfil', icon: User, href: `${eventBase}/profile`, authRequired: true },
        { id: 'event-login', label: 'Entrar', icon: LogIn, href: '/auth', hideIfAuthenticated: true },
      ];
    }
  },
  {
    id: 'default',
    match: () => true,
    getItems: (): NavItemType[] => [
      { id: 'home', label: 'Início', icon: Home, href: '/' },
      {
        id: 'events',
        label: 'Eventos',
        icon: Calendar,
        href: isAdmin ? '/admin/events' : '/events'
      },
      { id: 'profile', label: 'Perfil', icon: User, href: '/profile', authRequired: true },
      { id: 'logout', label: 'Sair', icon: LogOut, onClick: actions.logout, authRequired: true },
      { id: 'login', label: 'Entrar', icon: LogIn, href: '/auth', hideIfAuthenticated: true },
    ]
  }
];

const DesktopNavItem = ({
  item,
  isActive,
  isAdmin,
  onClick,
  mouseX,
}: {
  item: NavItemType;
  isActive: boolean;
  isAdmin: boolean;
  onClick: () => void;
  mouseX: ReturnType<typeof useMotionValue<number>>;
}) => {
  const ref = useRef<HTMLButtonElement>(null);
  const Icon = item.icon;

  const distance = useTransform(mouseX, (val) => {
    const b = ref.current?.getBoundingClientRect() ?? { x: 0, width: 0 };
    return val - (b.x + b.width / 2);
  });

  const sizeRaw = useTransform(distance, [-130, 0, 130], [40, 56, 40]);
  const size = useSpring(sizeRaw, { mass: 0.08, stiffness: 200, damping: 18 });

  const iconSizeRaw = useTransform(distance, [-130, 0, 130], [16, 22, 16]);
  const iconSize = useSpring(iconSizeRaw, { mass: 0.08, stiffness: 200, damping: 18 });

  return (
    <Tooltip>
      <TooltipTrigger
        render={
          <motion.button
            ref={ref}
            onClick={onClick}
            style={{ width: size, height: size }}
            className={cn(
              'relative flex items-center justify-center rounded-full outline-none transition-colors duration-200',
              isActive
                ? (isAdmin ? 'bg-amber-500 text-white shadow-lg shadow-amber-500/30' : 'bg-primary text-primary-foreground shadow-lg shadow-primary/30')
                : 'bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground',
            )}
            aria-label={item.label}
            aria-current={isActive ? 'page' : undefined}
            whileTap={{ scale: 0.88 }}
          >
            {isActive && (
              <motion.div
                className={cn(
                  "absolute inset-0 rounded-full ring-2 ring-offset-2 ring-offset-background",
                  isAdmin ? "ring-amber-500" : "ring-primary"
                )}
                initial={{ opacity: 0, scale: 0.5 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.5 }}
                transition={{ type: 'spring', stiffness: 400, damping: 30 }}
              />
            )}
            <motion.div style={{ width: iconSize, height: iconSize }} className="flex items-center justify-center">
              <Icon
                style={{ width: '100%', height: '100%' }}
                strokeWidth={isActive ? 2.5 : 2}
              />
            </motion.div>
          </motion.button>
        }
      />
      <TooltipContent side="top" sideOffset={8}>
        {item.label}
      </TooltipContent>
    </Tooltip>
  );
};

const MobileNavItem = ({
  item,
  isActive,
  isAdmin,
  onClick,
}: {
  item: NavItemType;
  isActive: boolean;
  isAdmin: boolean;
  onClick: () => void;
}) => {
  const Icon = item.icon;

  return (
    <motion.button
      onClick={onClick}
      className={cn(
        'relative flex flex-col items-center justify-center flex-1 py-3 gap-1.5 outline-none',
        isActive
          ? (isAdmin ? 'text-amber-600' : 'text-primary')
          : 'text-muted-foreground hover:text-foreground'
      )}
      aria-label={item.label}
      aria-current={isActive ? 'page' : undefined}
      whileTap={{ scale: 0.95 }}
    >
      <div className="absolute top-0 left-1/2 -translate-x-1/2">
        <motion.div
          className={cn(
            "h-1 rounded-b-full",
            isAdmin ? "bg-amber-500" : "bg-primary"
          )}
          initial={false}
          animate={{
            width: isActive ? 32 : 0,
            opacity: isActive ? 1 : 0,
          }}
          transition={{ type: 'spring', stiffness: 500, damping: 35 }}
        />
      </div>

      <motion.div
        animate={isActive ? { scale: 1.1, y: -1 } : { scale: 1, y: 0 }}
        transition={{ type: 'spring', stiffness: 400, damping: 25 }}
      >
        <Icon size={22} strokeWidth={isActive ? 2.4 : 2} />
      </motion.div>

      <span
        className={cn(
          'text-[10px] font-medium tracking-tight transition-colors duration-200',
          isActive ? (isAdmin ? 'text-amber-600' : 'text-primary') : 'text-muted-foreground'
        )}
      >
        {item.label}
      </span>
    </motion.button>
  );
};

export const NavigationDock = memo(function ({ className }: NavigationDockProps) {
  const { handleLogout, isAuthenticated } = useAuthActions();
  const location = useLocation();
  const navigate = useNavigate();

  const isAdmin = useMemo(() => location.pathname.includes('/admin'), [location.pathname]);

  const configs = useMemo(() => navConfigs({ logout: handleLogout }, isAdmin), [handleLogout, isAdmin]);

  const navItems = useMemo(() => {
    const pathParts = location.pathname.split('/').filter(Boolean);
    const config = configs.find(c => c.match(pathParts));
    const allItems = config?.getItems(pathParts) ?? [];

    return allItems.filter(item => {
      if (item.authRequired && !isAuthenticated) return false;
      if (item.hideIfAuthenticated && isAuthenticated) return false;
      return true;
    });
  }, [location.pathname, isAuthenticated, configs]);

  const activeId = useMemo(() => {
    const activeItem = [...navItems].reverse().find(item =>
      item.href === '/' ? location.pathname === '/' : (item.href ? location.pathname.startsWith(item.href) : false)
    );
    return activeItem?.id ?? '';
  }, [location.pathname, navItems]);

  const handleNavigate = (item: NavItemType) => {
    if (item.onClick) {
      void item.onClick();
      return;
    }

    if (item.href) {
      if (location.pathname === item.href) return;
      void navigate({ to: item.href });
    }
  };

  const mouseX = useMotionValue(0);

  if (navItems.length === 0) return null;

  return (
    <>
      {/* Desktop */}
      <nav
        role="navigation"
        className={cn('fixed bottom-8 left-1/2 -translate-x-1/2 z-50 hidden md:flex', className)}
        onMouseMove={(e) => { mouseX.set(e.clientX); }}
        onMouseLeave={() => { mouseX.set(0); }}
      >
        <motion.div
          initial={{ y: 20, opacity: 0, filter: 'blur(10px)' }}
          animate={{ y: 0, opacity: 1, filter: 'blur(0px)' }}
          transition={{ type: 'spring', stiffness: 260, damping: 24, delay: 0.05 }}
          className={cn(
            "flex items-center gap-2 px-3 py-3 rounded-full bg-background/80 backdrop-blur-2xl border shadow-lg shadow-black/5",
            isAdmin ? "border-amber-500/20" : "border-border/60"
          )}
        >
          {navItems.map((item) => (
            <DesktopNavItem
              key={item.id}
              item={item}
              isActive={activeId === item.id}
              isAdmin={isAdmin}
              onClick={() => { handleNavigate(item); }}
              mouseX={mouseX}
            />
          ))}
        </motion.div>
      </nav>

      {/* Mobile */}
      <nav
        role="navigation"
        className={cn('fixed bottom-0 left-0 right-0 z-50 md:hidden', className)}
      >
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ type: 'spring', stiffness: 260, damping: 24 }}
          className={cn(
            "flex items-stretch justify-around px-2 pb-safe bg-background/90 backdrop-blur-2xl border-t",
            isAdmin ? "border-amber-500/30" : "border-border/40"
          )}
        >
          {navItems.map((item) => (
            <MobileNavItem
              key={item.id}
              item={item}
              isActive={activeId === item.id}
              isAdmin={isAdmin}
              onClick={() => { handleNavigate(item); }}
            />
          ))}
        </motion.div>
        <div className="h-safe-area-inset-bottom bg-background/90 backdrop-blur-2xl" />
      </nav>
    </>
  );
});
