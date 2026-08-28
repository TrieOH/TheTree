import { useLocation, useNavigate } from "@tanstack/react-router";
import { useAuthActions } from "@trieoh/front-core";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import type { LucideIcon } from "lucide-react";
import { Calendar, Home, LayoutGrid, LogIn, LogOut, User } from "lucide-react";
import { motion, useMotionValue, useSpring, useTransform } from "motion/react";
import { memo, useMemo, useRef, useState } from "react";
import {
  clearAuthReturnTo,
  isAuthOnlyPath,
} from "@/features/auths/lib/auth-path";
import { cn } from "@/shared/lib/utils";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/shared/ui/shadcn/tooltip";

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

const getNavItems = (
  actions: { logout: () => Promise<void> },
  isAuthenticated: boolean,
): NavItemType[] =>
  [
    { id: "home", label: "Home", icon: Home, href: "/" },
    { id: "events", label: "Evento", icon: Calendar, href: "/events" },
    {
      id: "admin",
      label: "Admin",
      icon: LayoutGrid,
      href: "/admin/events",
      authRequired: true,
    },
    {
      id: "profile",
      label: "Perfil",
      icon: User,
      href: "/profile",
      authRequired: true,
    },
    {
      id: "logout",
      label: "Logout",
      icon: LogOut,
      onClick: actions.logout,
      authRequired: true,
    },
    {
      id: "login",
      label: "Entrar",
      icon: LogIn,
      href: "/auth",
      hideIfAuthenticated: true,
    },
  ].filter((item) => {
    if (item.authRequired && !isAuthenticated) return false;
    if (item.hideIfAuthenticated && isAuthenticated) return false;
    return true;
  });

const DesktopNavItem = ({
  item,
  isActive,
  isAdmin,
  onClick,
  mouseX,
  isDockHovered,
}: {
  item: NavItemType;
  isActive: boolean;
  isAdmin: boolean;
  onClick: () => void;
  mouseX: ReturnType<typeof useMotionValue<number>>;
  isDockHovered: boolean;
}) => {
  const ref = useRef<HTMLButtonElement>(null);
  const Icon = item.icon;

  const distance = useTransform(mouseX, (val) => {
    if (!isDockHovered) return -1000;
    const b = ref.current?.getBoundingClientRect() ?? { x: 0, width: 0 };
    return val - (b.x + b.width / 2);
  });

  const sizeRaw = useTransform(distance, [-130, 0, 130], [40, 56, 40]);
  const size = useSpring(sizeRaw, { mass: 0.08, stiffness: 200, damping: 18 });

  const iconSizeRaw = useTransform(distance, [-130, 0, 130], [16, 22, 16]);
  const iconSize = useSpring(iconSizeRaw, {
    mass: 0.08,
    stiffness: 200,
    damping: 18,
  });

  return (
    <Tooltip>
      <TooltipTrigger
        render={
          <motion.button
            ref={ref}
            onClick={onClick}
            style={{ width: size, height: size }}
            className={cn(
              "relative flex items-center justify-center rounded-full outline-none transition-colors duration-200",
              isActive
                ? isAdmin
                  ? "bg-accent text-accent-foreground shadow-lg shadow-accent/30"
                  : "bg-primary text-primary-foreground shadow-lg shadow-primary/30"
                : "bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground",
            )}
            aria-label={item.label}
            aria-current={isActive ? "page" : undefined}
            whileTap={{ scale: 0.88 }}
          >
            {isActive && (
              <motion.div
                className={cn(
                  "absolute inset-0 rounded-full ring-2 ring-offset-2 ring-offset-background",
                  isAdmin ? "ring-accent" : "ring-primary",
                )}
                initial={{ opacity: 0, scale: 0.5 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.5 }}
                transition={{ type: "spring", stiffness: 400, damping: 30 }}
              />
            )}
            <motion.div
              style={{ width: iconSize, height: iconSize }}
              className="flex items-center justify-center"
            >
              <Icon
                style={{ width: "100%", height: "100%" }}
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
        "relative flex flex-col items-center justify-center flex-1 py-3 gap-1.5 outline-none",
        isActive
          ? isAdmin
            ? "text-accent"
            : "text-primary"
          : "text-muted-foreground hover:text-foreground",
      )}
      aria-label={item.label}
      aria-current={isActive ? "page" : undefined}
      whileTap={{ scale: 0.95 }}
    >
      <div className="absolute top-0 left-1/2 -translate-x-1/2">
        <motion.div
          className={cn(
            "h-1 rounded-b-full",
            isAdmin ? "bg-accent" : "bg-primary",
          )}
          initial={false}
          animate={{
            width: isActive ? 32 : 0,
            opacity: isActive ? 1 : 0,
          }}
          transition={{ type: "spring", stiffness: 500, damping: 35 }}
        />
      </div>

      <motion.div
        animate={isActive ? { scale: 1.1, y: -1 } : { scale: 1, y: 0 }}
        transition={{ type: "spring", stiffness: 400, damping: 25 }}
      >
        <Icon size={22} strokeWidth={isActive ? 2.4 : 2} />
      </motion.div>

      <span
        className={cn(
          "text-[10px] font-medium tracking-tight transition-colors duration-200",
          isActive
            ? isAdmin
              ? "text-accent"
              : "text-primary"
            : "text-muted-foreground",
        )}
      >
        {item.label}
      </span>
    </motion.button>
  );
};

export const NavigationDock = memo(({ className }: NavigationDockProps) => {
  const { handleLogoutTo } = useAuthActions();
  const { isAuthenticated } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const locked = location.pathname === "/profile/setup";
  const logoutDestination = isAuthOnlyPath(location.pathname)
    ? "/"
    : location.href;

  const isAdmin = useMemo(
    () => location.pathname.startsWith("/admin"),
    [location.pathname],
  );
  const navItems = useMemo(
    () =>
      getNavItems(
        {
          logout: () => {
            clearAuthReturnTo(localStorage);
            return handleLogoutTo(logoutDestination);
          },
        },
        isAuthenticated,
      ),
    [handleLogoutTo, isAuthenticated, logoutDestination],
  );

  const activeId = useMemo(() => {
    const activeItem = [...navItems]
      .reverse()
      .find((item) =>
        item.href === "/"
          ? location.pathname === "/"
          : item.href
            ? location.pathname.startsWith(item.href)
            : false,
      );
    return activeItem?.id ?? "";
  }, [location.pathname, navItems]);

  const handleNavigate = (item: NavItemType) => {
    if (locked && !item.onClick) return;
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
  const [isDockHovered, setIsDockHovered] = useState(false);
  const hidden =
    location.pathname.endsWith("/certifications/editor") ||
    location.pathname.endsWith("/badges/editor") ||
    location.pathname.includes("/badges/") ||
    location.pathname === "/profile/edit";

  if (hidden) return null;
  if (navItems.length === 0) return null;

  return (
    <>
      {/* Desktop */}
      <nav
        className={cn(
          "fixed bottom-8 left-1/2 -translate-x-1/2 z-50 hidden md:flex",
          className,
        )}
        onMouseEnter={() => {
          setIsDockHovered(true);
        }}
        onMouseMove={(e) => {
          mouseX.set(e.clientX);
        }}
        onMouseLeave={() => {
          setIsDockHovered(false);
        }}
      >
        <motion.div
          initial={{ y: 20, opacity: 0, filter: "blur(10px)" }}
          animate={{ y: 0, opacity: 1, filter: "blur(0px)" }}
          transition={{
            type: "spring",
            stiffness: 260,
            damping: 24,
            delay: 0.05,
          }}
          className={cn(
            "flex items-center gap-2 px-3 py-3 rounded-full bg-background/80 backdrop-blur-2xl border shadow-lg shadow-black/5",
            isAdmin ? "border-accent/20" : "border-border/60",
          )}
        >
          {navItems.map((item) => (
            <DesktopNavItem
              key={item.id}
              item={item}
              isActive={activeId === item.id}
              isAdmin={isAdmin}
              onClick={() => {
                handleNavigate(item);
              }}
              mouseX={mouseX}
              isDockHovered={isDockHovered}
            />
          ))}
        </motion.div>
      </nav>

      {/* Mobile */}
      <nav
        className={cn(
          "fixed bottom-0 left-0 right-0 z-50 md:hidden",
          className,
        )}
      >
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ type: "spring", stiffness: 260, damping: 24 }}
          className={cn(
            "flex items-stretch justify-around px-2 pb-safe bg-background/90 backdrop-blur-2xl border-t",
            isAdmin ? "border-accent/30" : "border-border/40",
          )}
        >
          {navItems.map((item) => (
            <MobileNavItem
              key={item.id}
              item={item}
              isActive={activeId === item.id}
              isAdmin={isAdmin}
              onClick={() => {
                handleNavigate(item);
              }}
            />
          ))}
        </motion.div>
        <div className="h-safe-area-inset-bottom bg-background/90 backdrop-blur-2xl" />
      </nav>
    </>
  );
});
