import { useLocation, useNavigate } from '@tanstack/solid-router';
import { useAuth } from '@trieoh/identityx-sdk-ts-solid';
import { For, createMemo, createSignal, onSettled } from 'solid-js';

import CalendarIcon from '~icons/lucide/calendar';
import HomeIcon from '~icons/lucide/home';
import LayoutGridIcon from '~icons/lucide/layout-grid';
import LogInIcon from '~icons/lucide/log-in';
import LogOutIcon from '~icons/lucide/log-out';
import UserIcon from '~icons/lucide/user';

import { useSessionActions } from '@/features/auths/hooks/use-session-actions';
import { isAuthOnlyPath } from '@/features/auths/lib/auth-path';

const cn = (...classes: Array<string | false | null | undefined>) =>
  classes.filter(Boolean).join(' ');

type IconComponent = () => ReturnType<typeof HomeIcon>;

const asIcon = (icon: typeof HomeIcon): IconComponent =>
  icon as unknown as IconComponent;

const Home = asIcon(HomeIcon);
const Calendar = asIcon(CalendarIcon);
const LayoutGrid = asIcon(LayoutGridIcon);
const LogIn = asIcon(LogInIcon);
const LogOut = asIcon(LogOutIcon);
const User = asIcon(UserIcon);

interface NavItemType {
  id: string;
  label: string;
  icon: IconComponent;
  href?: string;
  authRequired?: boolean;
  hideIfAuthenticated?: boolean;
  onClick?: () => void | Promise<void>;
}

interface NavigationDockProps {
  className?: string;
}

const getNavItems = (
  actions: {
    logout: () => Promise<void>;
  },
  isAuthenticated: boolean,
): NavItemType[] =>
  [
    {
      id: 'home',
      label: 'Home',
      icon: Home,
      href: '/',
    },
    {
      id: 'events',
      label: 'Evento',
      icon: Calendar,
      href: '/events',
    },
    {
      id: 'admin',
      label: 'Admin',
      icon: LayoutGrid,
      href: '/admin/events',
      authRequired: true,
    },
    {
      id: 'profile',
      label: 'Perfil',
      icon: User,
      href: '/profile',
      authRequired: true,
    },
    {
      id: 'logout',
      label: 'Logout',
      icon: LogOut,
      onClick: actions.logout,
      authRequired: true,
    },
    {
      id: 'login',
      label: 'Entrar',
      icon: LogIn,
      href: '/auth',
      hideIfAuthenticated: true,
    },
  ].filter((item) => {
    if (item.authRequired && !isAuthenticated) {
      return false;
    }

    if (item.hideIfAuthenticated && isAuthenticated) {
      return false;
    }

    return true;
  });

/* ------------------------------------------------------------------ */
/* Desktop item                                                        */
/* ------------------------------------------------------------------ */

function DesktopNavItem(props: {
  item: NavItemType;
  isActive: boolean;
  isAdmin: boolean;
  onClick: () => void;

  /*
   * Passamos accessors simples em vez de depender
   * de tipos internos do Solid.
   */
  mouseX: () => number;
  isDockHovered: () => boolean;
}) {
  let ref: HTMLButtonElement | undefined;

  const Icon = props.item.icon;

  /*
   * Mede a distância entre o cursor e o centro do botão.
   *
   * 0 px      -> influência máxima
   * >= 130 px -> sem influência
   */
  const influence = createMemo(() => {
    if (!props.isDockHovered() || !ref) {
      return 0;
    }

    const bounds = ref.getBoundingClientRect();
    const center = bounds.left + bounds.width / 2;

    const distance = Math.abs(props.mouseX() - center);

    return Math.max(0, 1 - distance / 130);
  });

  const size = createMemo(() => 40 + influence() * 16);

  const iconSize = createMemo(() => 16 + influence() * 6);

  return (
    <button
      ref={ref}
      type="button"
      onClick={() => props.onClick()}
      title={props.item.label}
      aria-label={props.item.label}
      aria-current={props.isActive ? 'page' : undefined}
      style={{
        width: `${size()}px`,
        height: `${size()}px`,
      }}
      class={cn(
        'relative flex shrink-0 items-center justify-center rounded-full outline-none',

        /*
         * Substituição do spring do Motion.
         */
        'transition-[width,height,background-color,color,transform,box-shadow] duration-150 ease-out',

        'active:scale-[0.88]',

        props.isActive
          ? props.isAdmin
            ? 'bg-accent text-accent-foreground shadow-lg shadow-accent/30'
            : 'bg-primary text-primary-foreground shadow-lg shadow-primary/30'
          : 'bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground',
      )}
    >
      {/* Ring do item ativo */}
      <div
        class={cn(
          'pointer-events-none absolute inset-0 rounded-full',
          'ring-2 ring-offset-2 ring-offset-background',
          'transition-[opacity,transform] duration-200 ease-out',

          props.isAdmin ? 'ring-accent' : 'ring-primary',

          props.isActive
            ? 'scale-100 opacity-100'
            : 'scale-50 opacity-0',
        )}
      />

      {/*
       * Não passamos width/style diretamente para Icon.
       *
       * Assim não dependemos da tipagem do unplugin-icons.
       * O SVG ocupa 100% do wrapper.
       */}
      <div
        style={{
          width: `${iconSize()}px`,
          height: `${iconSize()}px`,
          'stroke-width': props.isActive ? '2.5' : '2',
        }}
        class={cn(
          'flex items-center justify-center',
          'transition-[width,height] duration-150 ease-out',
          '[&>svg]:h-full [&>svg]:w-full',
        )}
      >
        <Icon />
      </div>
    </button>
  );
}

/* ------------------------------------------------------------------ */
/* Mobile item                                                         */
/* ------------------------------------------------------------------ */

function MobileNavItem(props: {
  item: NavItemType;
  isActive: boolean;
  isAdmin: boolean;
  onClick: () => void;
}) {
  const Icon = props.item.icon;

  return (
    <button
      type="button"
      onClick={() => props.onClick()}
      title={props.item.label}
      aria-label={props.item.label}
      aria-current={props.isActive ? 'page' : undefined}
      class={cn(
        'relative flex flex-1 flex-col items-center justify-center gap-1.5 py-3 outline-none',
        'transition-[color,transform] duration-200',
        'active:scale-95',

        props.isActive
          ? props.isAdmin
            ? 'text-accent'
            : 'text-primary'
          : 'text-muted-foreground hover:text-foreground',
      )}
    >
      {/* Indicador superior */}
      <div class="absolute top-0 left-1/2 -translate-x-1/2">
        <div
          class={cn(
            'h-1 rounded-b-full',
            'transition-[width,opacity] duration-200 ease-out',

            props.isAdmin ? 'bg-accent' : 'bg-primary',

            props.isActive
              ? 'w-8 opacity-100'
              : 'w-0 opacity-0',
          )}
        />
      </div>

      {/* Ícone */}
      <div
        style={{
          'stroke-width': props.isActive ? '2.4' : '2',
        }}
        class={cn(
          'h-5.5 w-5.5',
          'transition-transform duration-200 ease-out',
          '[&>svg]:h-full [&>svg]:w-full',

          props.isActive
            ? '-translate-y-px scale-110'
            : 'translate-y-0 scale-100',
        )}
      >
        <Icon />
      </div>

      <span
        class={cn(
          'text-[10px] font-medium tracking-tight',
          'transition-colors duration-200',

          props.isActive
            ? props.isAdmin
              ? 'text-accent'
              : 'text-primary'
            : 'text-muted-foreground',
        )}
      >
        {props.item.label}
      </span>
    </button>
  );
}

/* ------------------------------------------------------------------ */
/* Navigation Dock                                                     */
/* ------------------------------------------------------------------ */

export function NavigationDock(props: NavigationDockProps) {
  const { logoutTo } = useSessionActions();

  /*
   * SDK IdentityX:
   *
   * isAuthenticated é Accessor.
   *
   * CORRETO:
   * isAuthenticated()
   */
  const { isAuthenticated } = useAuth();

  /*
   * TanStack Solid Router:
   *
   * useLocation() também retorna Accessor.
   *
   * CORRETO:
   * location().pathname
   */
  const location = useLocation();

  const navigate = useNavigate();

  /*
   * Centralizamos pathname para evitar repetir location().
   */
  const pathname = () => location().pathname;

  const locked = () => pathname() === '/profile/setup';

  const logoutDestination = () =>
    isAuthOnlyPath(pathname())
      ? '/'
      : location().href;

  const isAdmin = () =>
    pathname().startsWith('/admin');

  /*
   * Recalcula somente quando isAuthenticated() mudar.
   *
   * Login/logout atualizam authStore e o SDK invalida
   * esse memo automaticamente.
   */
  const navItems = createMemo(() =>
    getNavItems(
      {
        logout: async () => {
          await logoutTo(logoutDestination());
        },
      },
      isAuthenticated(),
    ),
  );

  /*
   * Item ativo também é derivado de forma reativa.
   */
  const activeId = createMemo(() => {
    const path = pathname();

    const activeItem = [...navItems()]
      .reverse()
      .find((item) => {
        if (!item.href) {
          return false;
        }

        if (item.href === '/') {
          return path === '/';
        }

        return path.startsWith(item.href);
      });

    return activeItem?.id ?? '';
  });

  const handleNavigate = (item: NavItemType) => {
    /*
     * Durante profile/setup impedimos navegação normal,
     * mas ainda permitimos logout.
     */
    if (locked() && !item.onClick) {
      return;
    }

    if (item.onClick) {
      void item.onClick();
      return;
    }

    if (!item.href) {
      return;
    }

    if (pathname() === item.href) {
      return;
    }

    void navigate({
      to: item.href as never,
    });
  };

  /* ---------------------------------------------------------------- */
  /* Dock magnético                                                     */
  /* ---------------------------------------------------------------- */

  const [mouseX, setMouseX] = createSignal(0);
  const [isDockHovered, setIsDockHovered] = createSignal(false);

  let mouseFrame: number | undefined;

  const handleMouseMove = (event: MouseEvent) => {
    const clientX = event.clientX;

    if (mouseFrame !== undefined) {
      cancelAnimationFrame(mouseFrame);
    }

    mouseFrame = requestAnimationFrame(() => {
      setMouseX(clientX);
      mouseFrame = undefined;
    });
  };

  /* ---------------------------------------------------------------- */
  /* Entrada do dock — Solid 2                                         */
  /* ---------------------------------------------------------------- */

  const [settled, setSettled] = createSignal(false);

  /*
   * Solid 2:
   *
   * onMount -> onSettled
   *
   * Além disso, onSettled pode retornar cleanup.
   */
  onSettled(() => {
    const frame = requestAnimationFrame(() => {
      setSettled(true);
    });

    return () => {
      cancelAnimationFrame(frame);

      if (mouseFrame !== undefined) {
        cancelAnimationFrame(mouseFrame);
      }
    };
  });

  /* ---------------------------------------------------------------- */
  /* Rotas em que o dock deve desaparecer                              */
  /* ---------------------------------------------------------------- */

  const hidden = createMemo(() => {
    const path = pathname();

    return (
      path.endsWith('/certifications/editor') ||
      path.endsWith('/badges/editor') ||
      path.includes('/badges/') ||
      (path.includes('/occurrences/') &&
        path.endsWith('/draw')) ||
      path === '/profile/edit'
    );
  });

  return (
    <>
      {!hidden() && navItems().length > 0 && (
        <>
          {/* ---------------------------------------------------------- */}
          {/* Desktop                                                    */}
          {/* ---------------------------------------------------------- */}

          <nav
            class={cn(
              'fixed bottom-8 left-1/2 z-50',
              'hidden -translate-x-1/2 md:flex',
              props.className,
            )}
            onMouseEnter={() => {
              setIsDockHovered(true);
            }}
            onMouseMove={handleMouseMove}
            onMouseLeave={() => {
              setIsDockHovered(false);
            }}
          >
            <div
              class={cn(
                'flex items-center gap-2',
                'rounded-full border',
                'bg-background/80 px-3 py-3',
                'shadow-lg shadow-black/5',
                'backdrop-blur-2xl',

                /*
                 * Substitui:
                 *
                 * initial
                 * animate
                 * transition
                 *
                 * do motion/solid.
                 */
                'transition-[transform,opacity,filter] duration-300 ease-out',

                settled()
                  ? 'translate-y-0 opacity-100 blur-0'
                  : 'translate-y-5 opacity-0 blur-[10px]',

                isAdmin()
                  ? 'border-accent/20'
                  : 'border-border/60',
              )}
            >
              <For each={navItems()}>
                {(item) => (
                  <DesktopNavItem
                    item={item}
                    isActive={activeId() === item.id}
                    isAdmin={isAdmin()}
                    onClick={() => {
                      handleNavigate(item);
                    }}
                    mouseX={mouseX}
                    isDockHovered={isDockHovered}
                  />
                )}
              </For>
            </div>
          </nav>

          {/* ---------------------------------------------------------- */}
          {/* Mobile                                                     */}
          {/* ---------------------------------------------------------- */}

          <nav
            class={cn(
              'fixed right-0 bottom-0 left-0 z-50 md:hidden',
              props.className,
            )}
          >
            <div
              class={cn(
                'flex items-stretch justify-around',
                'border-t bg-background/90',
                'px-2 pb-safe',
                'backdrop-blur-2xl',

                'transition-[transform,opacity] duration-300 ease-out',

                settled()
                  ? 'translate-y-0 opacity-100'
                  : 'translate-y-5 opacity-0',

                isAdmin()
                  ? 'border-accent/30'
                  : 'border-border/40',
              )}
            >
              <For each={navItems()}>
                {(item) => (
                  <MobileNavItem
                    item={item}
                    isActive={activeId() === item.id}
                    isAdmin={isAdmin()}
                    onClick={() => {
                      handleNavigate(item);
                    }}
                  />
                )}
              </For>
            </div>

            <div class="h-safe-area-inset-bottom bg-background/90 backdrop-blur-2xl" />
          </nav>
        </>
      )}
    </>
  );
}