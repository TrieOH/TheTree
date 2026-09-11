import { useLocation, useNavigate } from '@tanstack/solid-router';
import { useAuth } from '@trieoh/identityx-sdk-ts-solid';
import { Dynamic } from '@solidjs/web';
import {
  animate,
  mapValue,
  motionValue,
  springValue,
  styleEffect,
  transformValue,
} from 'motion';
import { For, createEffect, createMemo, onCleanup, onSettled, untrack } from 'solid-js';

import CalendarIcon from '~icons/lucide/calendar';
import HomeIcon from '~icons/lucide/home';
import LayoutGridIcon from '~icons/lucide/layout-grid';
import LogInIcon from '~icons/lucide/log-in';
import LogOutIcon from '~icons/lucide/log-out';
import UserIcon from '~icons/lucide/user';

import { useSessionActions } from '@/features/auths/hooks/use-session-actions';
import { isAuthOnlyPath } from '@/features/auths/lib/auth-path';
import { ownProfilePath } from '@/features/profile/lib/own-profile-path';
import { Tooltip } from '@/shared/ui/Tooltip';

const cn = (...classes: Array<string | false | null | undefined>) =>
  classes.filter(Boolean).join(' ');

type IconComponent = () => ReturnType<typeof HomeIcon>;
type NumericMotionValue = ReturnType<typeof motionValue<number>>;

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
  actions: { logout: () => Promise<void> },
  isAuthenticated: boolean,
): NavItemType[] =>
  [
    { id: 'home', label: 'Home', icon: Home, href: '/' },
    { id: 'events', label: 'Evento', icon: Calendar, href: '/events' },
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
    if (item.authRequired && !isAuthenticated) return false;
    if (item.hideIfAuthenticated && isAuthenticated) return false;
    return true;
  });

const playPressIn = (element: HTMLElement, pressedScale: number) =>
  animate(
    element,
    { scale: pressedScale },
    {
      duration: 0.14,
      ease: 'easeOut',
    },
  );

const playPressOut = (element: HTMLElement) =>
  animate(
    element,
    { scale: 1 },
    {
      duration: 0.28,
      ease: 'easeOut',
    },
  );


function DesktopNavItem(props: {
  item: NavItemType;
  isActive: boolean;
  isAdmin: boolean;
  onClick: () => void;
  mouseX: NumericMotionValue;
}) {
  let buttonRef: HTMLButtonElement | undefined;
  let iconRef: HTMLDivElement | undefined;
  let ringRef: HTMLDivElement | undefined;
  let tapAnimation: ReturnType<typeof animate> | undefined;
  const icon = untrack(() => props.item.icon);
  const mouseX = untrack(() => props.mouseX);

  const distance = transformValue(() => {
    const pointerX = mouseX.get();

    if (pointerX <= -999) return -1000;

    const bounds = buttonRef?.getBoundingClientRect();
    if (!bounds) return -1000;

    return pointerX - (bounds.x + bounds.width / 2);
  });

  const sizeRaw = mapValue(distance, [-130, 0, 130], [40, 56, 40]);
  const size = springValue(sizeRaw, {
    mass: 0.08,
    stiffness: 200,
    damping: 18,
  });

  const iconSizeRaw = mapValue(distance, [-130, 0, 130], [16, 22, 16]);
  const iconSize = springValue(iconSizeRaw, {
    mass: 0.08,
    stiffness: 200,
    damping: 18,
  });

  const activeTarget = motionValue(0);
  const activeSpring = springValue(activeTarget, {
    stiffness: 400,
    damping: 30,
  });
  const ringScale = mapValue(activeSpring, [0, 1], [0.5, 1], {
    clamp: false,
  });

  createEffect(
    () => props.isActive,
    (isActive) => {
      activeTarget.set(isActive ? 1 : 0);
    },
  );

  onSettled(() => {
    if (!buttonRef || !iconRef || !ringRef) return;

    const cancelButton = styleEffect(buttonRef, {
      width: size,
      height: size,
    });

    const cancelIcon = styleEffect(iconRef, {
      width: iconSize,
      height: iconSize,
    });

    const cancelRing = styleEffect(ringRef, {
      opacity: activeSpring,
      scale: ringScale,
    });

    return () => {
      tapAnimation?.stop();
      cancelButton();
      cancelIcon();
      cancelRing();
    };
  });

  return (
    <Tooltip label={props.item.label}>
      <button
        ref={(element) => { buttonRef = element; }}
        type="button"
        onPointerDown={() => {
          if (!buttonRef) return;

          tapAnimation?.stop();
          tapAnimation = playPressIn(buttonRef, 0.84);
        }}
        onPointerLeave={() => {
          if (!buttonRef) return;
          tapAnimation?.stop();
          tapAnimation = playPressOut(buttonRef);
        }}
        onPointerUp={() => {
          if (!buttonRef) return;

          tapAnimation?.stop();
          tapAnimation = playPressOut(buttonRef);

        }}
        onPointerCancel={() => {
          if (!buttonRef) return;

          tapAnimation?.stop();
          tapAnimation = playPressOut(buttonRef);

        }}
        onClick={() => props.onClick()}
        aria-label={props.item.label}
        aria-current={props.isActive ? 'page' : undefined}
        style={{
          width: '40px',
          height: '40px',
        }}
        class={cn(
          'relative flex shrink-0 items-center justify-center rounded-full outline-none',
          'transition-[background-color,color] duration-200 ease-out',
          props.isActive
            ? props.isAdmin
              ? 'bg-accent text-accent-foreground shadow-lg shadow-accent/30'
              : 'bg-primary text-primary-foreground shadow-lg shadow-primary/30'
            : 'bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground',
        )}
      >
        <div
          ref={(element) => { ringRef = element; }}
          class={cn(
            'pointer-events-none absolute inset-0 rounded-full',
            'ring-2 ring-offset-2 ring-offset-background',
            props.isAdmin ? 'ring-accent' : 'ring-primary',
          )}
          style={{ opacity: '0', transform: 'scale(0.5)' }}
        />

        <div
          ref={(element) => { iconRef = element; }}
          style={{
            width: '16px',
            height: '16px',
            'stroke-width': props.isActive ? '2.5' : '2',
          }}
          class="flex items-center justify-center [&>svg]:h-full [&>svg]:w-full"
        >
          <Dynamic component={icon} />
        </div>
      </button>
    </Tooltip>
  );
}


function MobileNavItem(props: {
  item: NavItemType;
  isActive: boolean;
  isAdmin: boolean;
  onClick: () => void;
}) {
  let buttonRef: HTMLButtonElement | undefined;
  let iconRef: HTMLDivElement | undefined;
  let tapAnimation: ReturnType<typeof animate> | undefined;

  const icon = untrack(() => props.item.icon);

  const iconTarget = motionValue(
    untrack(() => (props.isActive ? 1 : 0)),
  );
  const iconActive = springValue(iconTarget, {
    stiffness: 400,
    damping: 25,
  });
  const iconScale = mapValue(iconActive, [0, 1], [1, 1.1], {
    clamp: false,
  });
  const iconY = mapValue(iconActive, [0, 1], [0, -1], {
    clamp: false,
  });

  createEffect(
    () => props.isActive,
    (isActive) => {
      const next = isActive ? 1 : 0;
      iconTarget.set(next);
    },
  );

  onSettled(() => {
    const cancelIcon = styleEffect(iconRef, {
      scale: iconScale,
      y: iconY,
    });

    return () => {
      tapAnimation?.stop();
      cancelIcon();
    };
  });

  return (
    <Tooltip label={props.item.label}>
      <button
        ref={(element) => { buttonRef = element; }}
        type="button"
        onPointerDown={() => {
          if (!buttonRef) return;

          tapAnimation?.stop();
          tapAnimation = playPressIn(buttonRef, 0.9);
        }}
        onPointerUp={() => {
          if (!buttonRef) return;

          tapAnimation?.stop();
          tapAnimation = playPressOut(buttonRef);

        }}
        onPointerCancel={() => {
          if (!buttonRef) return;

          tapAnimation?.stop();
          tapAnimation = playPressOut(buttonRef);

        }}
        onClick={() => props.onClick()}
        aria-label={props.item.label}
        aria-current={props.isActive ? 'page' : undefined}
        class={cn(
          'relative flex flex-1 flex-col items-center justify-center gap-1.5 py-3 outline-none',
          'transition-colors duration-200 ease-out',
          props.isActive
            ? props.isAdmin
              ? 'text-accent'
              : 'text-primary'
            : 'text-muted-foreground hover:text-foreground',
        )}
      >
        <div class="absolute top-0 left-1/2 -translate-x-1/2">
          <div
            class={cn(
              'h-1 rounded-b-full transition-[width,opacity] duration-200 ease-out',
              props.isAdmin ? 'bg-accent' : 'bg-primary',
            )}
            style={{
              width: props.isActive ? '32px' : '0px',
              opacity: props.isActive ? '1' : '0',
            }}
          />
        </div>

        <div
          ref={(element) => { iconRef = element; }}
          class="h-5.5 w-5.5 [&>svg]:h-full [&>svg]:w-full"
          style={{
            'stroke-width': props.isActive ? '2.4' : '2',
          }}
        >
          <Dynamic component={icon} />
        </div>

        <span
          class={cn(
            'text-[10px] font-medium tracking-tight transition-colors duration-200',
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
    </Tooltip>
  );
}


export function NavigationDock(props: NavigationDockProps) {
  let desktopDockRef: HTMLDivElement | undefined;
  let mobileDockRef: HTMLDivElement | undefined;

  const { logoutTo } = useSessionActions();
  const { isAuthenticated } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const pathname = () => location().pathname;
  const locked = () => pathname() === '/profile/setup';
  const logoutDestination = () =>
    isAuthOnlyPath(pathname()) ? '/' : location().href;
  const isAdmin = () => pathname().startsWith('/admin');

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

  const activeId = createMemo(() => {
    const path = pathname();

    const activeItem = [...navItems()].reverse().find((item) => {
      if (!item.href) return false;
      if (item.href === '/') return path === '/';
      return path.startsWith(item.href);
    });

    return activeItem?.id ?? '';
  });

  const handleNavigate = (item: NavItemType) => {
    if (locked() && !item.onClick) return;

    if (item.onClick) {
      void item.onClick();
      return;
    }

    if (!item.href || pathname() === item.href) return;

    if (item.id === 'profile' && pathname() === ownProfilePath()) return;

    void navigate({
      to: item.href as never,
    });
  };

  const mouseX = motionValue(-1000);

  const hidden = createMemo(() => {
    const path = pathname();

    return (
      path.endsWith('/certifications/editor') ||
      path.endsWith('/badges/editor') ||
      path.includes('/badges/') ||
      (path.includes('/occurrences/') && path.endsWith('/draw')) ||
      path === '/profile/edit'
    );
  });

  onCleanup(() => {
    mouseX.set(-1000);
  });

  createEffect(
    () => !hidden() && navItems().length > 0,
    (visible, wasVisible) => {
      if (!visible || wasVisible) return;

      queueMicrotask(() => {
        if (desktopDockRef) {
          animate(
            desktopDockRef,
            {
              y: [20, 0],
              opacity: [0, 1],
              filter: ['blur(10px)', 'blur(0px)'],
            },
            {
              type: 'spring',
              stiffness: 260,
              damping: 24,
              delay: 0.05,
            },
          );
        }

        if (mobileDockRef) {
          animate(
            mobileDockRef,
            {
              y: [20, 0],
              opacity: [0, 1],
            },
            {
              type: 'spring',
              stiffness: 260,
              damping: 24,
            },
          );
        }
      });
    },
  );

  return (
    <>
      {!hidden() && navItems().length > 0 && (
        <>
          <nav
            class={cn(
              'fixed bottom-8 left-1/2 z-50',
              'hidden -translate-x-1/2 md:flex',
              props.className,
            )}
            onPointerEnter={(event) => {
              mouseX.set(event.clientX);
            }}
            onPointerMove={(event) => {
              mouseX.set(event.clientX);
            }}
            onPointerLeave={() => {
              mouseX.set(-1000);
            }}
          >
            <div
              ref={(element) => { desktopDockRef = element; }}
              class={cn(
                'flex items-center gap-2 rounded-full border',
                'bg-background/80 px-3 py-3',
                'shadow-lg shadow-black/5 backdrop-blur-2xl',
                isAdmin() ? 'border-accent/20' : 'border-border/60',
              )}
              style={{
                transform: 'translateY(20px)',
                opacity: '0',
                filter: 'blur(10px)',
              }}
            >
              <For each={navItems()}>
                {(item) => (
                  <DesktopNavItem
                    item={item}
                    isActive={activeId() === item.id}
                    isAdmin={isAdmin()}
                    onClick={() => handleNavigate(item)}
                    mouseX={mouseX}
                  />
                )}
              </For>
            </div>
          </nav>

          <nav
            class={cn(
              'fixed right-0 bottom-0 left-0 z-50 md:hidden',
              props.className,
            )}
          >
            <div
              ref={(element) => { mobileDockRef = element; }}
              class={cn(
                'flex items-stretch justify-around',
                'border-t bg-background/90 px-2 pb-safe',
                'backdrop-blur-2xl',
                isAdmin() ? 'border-accent/30' : 'border-border/40',
              )}
              style={{
                transform: 'translateY(20px)',
                opacity: '0',
              }}
            >
              <For each={navItems()}>
                {(item) => (
                  <MobileNavItem
                    item={item}
                    isActive={activeId() === item.id}
                    isAdmin={isAdmin()}
                    onClick={() => handleNavigate(item)}
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
