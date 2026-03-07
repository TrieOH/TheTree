import React, { useEffect, useState, useMemo, useCallback } from 'react';
import { Tabs, TabsList, TabsTrigger } from '@/shared/ui/shadcn/tabs';
import { AnimatePresence, motion } from 'motion/react';
import { cn } from '@/shared/lib/utils';
import { useNavigate } from '@tanstack/react-router';
import { RefreshCcw } from 'lucide-react';

type TabItem = {
  value: string;
  label: string;
  icon?: React.ElementType;
  content: React.ReactNode;
  onRefresh?: () => void;
};

type Props = {
  items: TabItem[];
  initialValue?: string;
  deferTabSwitch?: boolean;
  className?: string;
};

const TabIcon = React.memo(({ icon: Icon, className }: { icon: React.ElementType; className?: string }) => {
  return <Icon className={cn("h-4 w-4 shrink-0", className)} />;
});

TabIcon.displayName = 'TabIcon';

const TabTriggerItem = React.memo(({ 
  tab, 
  isActive, 
  isMobile 
}: { 
  tab: TabItem; 
  isActive: boolean; 
  isMobile: boolean 
}) => {
  const [isRefreshing, setIsRefreshing] = useState(false);

  const handleRefresh = useCallback(async (e: React.MouseEvent) => {
    e.stopPropagation();
    if (tab.onRefresh) {
      setIsRefreshing(true);
      await tab.onRefresh();
      // Artificial delay to show the animation
      setTimeout(() => setIsRefreshing(false), 500);
    }
  }, [tab]);

  return (
    <TabsTrigger
      value={tab.value}
      className={cn(
        "relative h-12 md:h-14 flex items-center gap-2 px-4 py-3 text-sm font-medium transition-colors z-10",
        "justify-center md:justify-start text-muted-foreground hover:text-foreground whitespace-nowrap shrink-0",
        "flex-1 md:flex-none data-[state=active]:bg-transparent data-[state=active]:text-foreground",
        "border-b-2 md:border-b-0 md:border-l-2 border-transparent rounded-none",
        "md:rounded-r-lg md:w-full min-w-fit"
      )}
    >
      {isActive && (
        <motion.div
          layoutId="active-tab-highlight"
          className={cn(
            "absolute inset-0 bg-muted/50 rounded-none",
            "md:rounded-r-lg md:rounded-l-none"
          )}
          initial={false}
          transition={{ type: isMobile ? 'tween' : 'spring', stiffness: 400, damping: 30 }}
        />
      )}
      {isActive && (
        <motion.div
          layoutId="active-tab-border"
          className={cn(
            "absolute bottom-0 left-0 right-0 h-0.5 bg-primary z-20",
            "md:left-0 md:top-0 md:bottom-0 md:right-auto md:w-0.5 md:h-full"
          )}
          initial={false}
          transition={{ type: isMobile ? 'tween' : 'spring', stiffness: 400, damping: 30 }}
        />
      )}
      <span className="relative z-20 flex items-center justify-between w-full gap-2">
        <div className="flex items-center gap-2">
          {tab.icon && <TabIcon icon={tab.icon} />}
          <span className="hidden sm:inline">{tab.label}</span>
          <span className="sm:hidden text-xs">{tab.label.split(' ')[0]}</span>
        </div>
        
        {isActive && tab.onRefresh && (
          <motion.button
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={handleRefresh}
            className="p-1 hover:bg-muted-foreground/10 rounded-full transition-colors hidden md:block"
          >
            <RefreshCcw className={cn("h-3.5 w-3.5", isRefreshing && "animate-spin")} />
          </motion.button>
        )}
      </span>
    </TabsTrigger>
  );
});

TabTriggerItem.displayName = 'TabTriggerItem';

export default function CustomTabs({
  items,
  initialValue,
  deferTabSwitch = false,
  className = ''
}: Props) {
  const safeItems = useMemo(() => items ?? [], [items]);
  const first = useMemo(() => initialValue ?? safeItems[0]?.value ?? '', [initialValue, safeItems]);

  const [activeTab, setActiveTab] = useState<string>(first);
  const [displayTab, setDisplayTab] = useState<string | null>(first);
  const [pendingTab, setPendingTab] = useState<string | null>(null);
  const [direction, setDirection] = useState<number>(0);
  const [isMobile, setIsMobile] = useState<boolean>(false);
  const navigate = useNavigate();

  const idxOf = useCallback((v: string) => safeItems.findIndex((t) => t.value === v), [safeItems]);

  // Sync with initialValue (e.g. from search params)
  useEffect(() => {
    if (initialValue && initialValue !== activeTab) {
      const newIndex = idxOf(initialValue);
      const currentIndex = idxOf(activeTab);
      if (newIndex !== -1) {
        setDirection(newIndex > currentIndex ? 1 : -1);
        setActiveTab(initialValue);
        setDisplayTab(initialValue);
        
        // Refresh when navigating via initialValue (back/forward)
        const item = safeItems[newIndex];
        if (item?.onRefresh) item.onRefresh();
      }
    }
  }, [initialValue, activeTab, idxOf, safeItems]);

  useEffect(() => {
    const mql = window.matchMedia('(max-width: 767px)');
    const onChange = (e: MediaQueryListEvent | MediaQueryList) => setIsMobile(e.matches);
    
    onChange(mql);
    mql.addEventListener('change', onChange);
    return () => mql.removeEventListener('change', onChange);
  }, []);

  const handleTabChange = useCallback((newValue: string) => {
    // Call refresh on every click, even if it's the active tab
    const activeItem = safeItems.find(item => item.value === newValue);
    if (activeItem?.onRefresh) activeItem.onRefresh();

    if (newValue === activeTab && !deferTabSwitch) return;
    const newIndex = idxOf(newValue);
    const currentIndex = idxOf(activeTab);
    setDirection(newIndex > currentIndex ? 1 : -1);

    if (deferTabSwitch) {
      if (pendingTab === newValue) return;
      setPendingTab(newValue);
      setDisplayTab(null);
    } else {
      setActiveTab(newValue);
      setDisplayTab(newValue);
    }
    
    navigate({
      to: ".",
      search: (prev) => ({
        ...prev,
        tab: newValue,
      }),
    });
  }, [activeTab, deferTabSwitch, idxOf, navigate, pendingTab, safeItems]);

  const contentVariants = useMemo(() => ({
    enter: (dir: number) => ({
      opacity: 0,
      x: isMobile ? (dir > 0 ? 30 : -30) : 0,
      y: isMobile ? 0 : (dir > 0 ? 20 : -20)
    }),
    center: { opacity: 1, x: 0, y: 0 },
    exit: (dir: number) => ({
      opacity: 0,
      x: isMobile ? (dir > 0 ? -30 : 30) : 0,
      y: isMobile ? 0 : (dir > 0 ? -20 : 20)
    })
  }), [isMobile]);

  const activeContent = useMemo(() => {
    return safeItems.find((t) => t.value === displayTab)?.content;
  }, [safeItems, displayTab]);

  if (safeItems.length === 0) return null;

  return (
    <Tabs 
      value={activeTab} 
      onValueChange={handleTabChange} 
      orientation="vertical" 
      className={cn(
        "flex flex-col-reverse md:flex-row w-full h-full",
        "gap-4 md:gap-8 mx-auto items-start",
        className
      )}
    >
      {/* Sidebar */}
      <div 
        className={cn(
          "fixed bottom-0 left-0 right-0 z-50 md:relative",
          "md:w-64 shrink-0 overflow-hidden bg-background",
          "md:my-5 md:ml-2",
          isMobile && "border-t border-border"
        )}
      >
        <div
          className={cn(
            "absolute bottom-0 left-0 bg-border",
            "md:left-0 md:top-0 md:bottom-0 md:right-auto md:w-0.5 md:h-full"
          )}
        />
        <TabsList
          className={cn(
            "flex h-full flex-row justify-around md:flex-col bg-transparent p-0 gap-0",
            "overflow-x-auto overflow-y-hidden md:overflow-visible scrollbar-hide w-full relative",
            "md:pr-2 md:py-2"
          )}
        >
          {safeItems.map((tab) => (
            <TabTriggerItem 
              key={tab.value} 
              tab={tab} 
              isActive={activeTab === tab.value} 
              isMobile={isMobile}
            />
          ))}
        </TabsList>
      </div>

      {/* Content */}
      <div className="relative flex-1 w-full h-full overflow-hidden">
        <AnimatePresence
          mode="wait"
          custom={direction}
          onExitComplete={() => {
            if (pendingTab) {
              setActiveTab(pendingTab);
              setDisplayTab(pendingTab);
              setPendingTab(null);
            }
          }}
        >
          {displayTab && (
            <motion.div 
              key={displayTab} 
              custom={direction} 
              variants={contentVariants} 
              initial="enter" 
              animate="center" 
              exit="exit" 
              transition={{ duration: 0.25, ease: [0.4, 0, 0.2, 1] }} 
              className={cn("absolute inset-0 h-full", isMobile && "pb-12")}
            >
              <div className="flex flex-col h-full overflow-y-auto overflow-x-hidden p-4">
                {activeContent}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </Tabs>
  );
}
