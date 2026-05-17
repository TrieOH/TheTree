import { useEffect, useState, useMemo, useCallback, useRef, memo } from 'react';
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

const TabIcon = memo(({ icon: Icon, className }: { icon: React.ElementType; className?: string }) => {
  return <Icon className={cn("h-4 w-4 shrink-0", className)} />;
});

TabIcon.displayName = 'TabIcon';

const springTransition = {
  type: "spring",
  stiffness: 500,
  damping: 38,
  mass: 1
} as const;

const TabTriggerItem = memo(({ 
  tab, 
  isActive, 
}: { 
  tab: TabItem; 
  isActive: boolean; 
}) => {
  const [isRefreshing, setIsRefreshing] = useState(false);

  const handleRefresh = useCallback(async (e: React.MouseEvent | React.KeyboardEvent) => {
    e.stopPropagation();
    if (tab.onRefresh) {
      setIsRefreshing(true);
      tab.onRefresh();
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
          transition={springTransition}
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
          transition={springTransition}
        />
      )}
      <span className="relative z-20 flex items-center justify-between w-full gap-2">
        <div className="flex items-center gap-2">
          {tab.icon && <TabIcon icon={tab.icon} />}
          <span className="hidden sm:inline">{tab.label}</span>
          <span className="sm:hidden text-xs">{tab.label.split(' ')[0]}</span>
        </div>
        
        {isActive && tab.onRefresh && (
          <motion.span
            role="button"
            tabIndex={0}
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={handleRefresh}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                handleRefresh(e)
              }
            }}
            className="p-1 hover:bg-muted-foreground/10 rounded-full transition-colors hidden md:block cursor-pointer"
          >
            <RefreshCcw className={cn("h-3.5 w-3.5", isRefreshing && "animate-spin")} />
          </motion.span>
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
  const [isMobile, setIsMobile] = useState<boolean>(false);
  const navigate = useNavigate();

  const idxOf = useCallback((v: string) => safeItems.findIndex((t) => t.value === v), [safeItems]);

  const activeIndex = idxOf(activeTab);
  const prevIndexRef = useRef(activeIndex);
  const direction = activeIndex > prevIndexRef.current ? 1 : activeIndex < prevIndexRef.current ? -1 : 0;

  useEffect(() => {
    prevIndexRef.current = activeIndex;
  }, [activeIndex]);

  // Sync with initialValue
  useEffect(() => {
    if (initialValue && initialValue !== activeTab) {
      const newIndex = idxOf(initialValue);
      if (newIndex !== -1) {
        setActiveTab(initialValue);
        setDisplayTab(initialValue);
        
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
    const activeItem = safeItems.find(item => item.value === newValue);
    if (activeItem?.onRefresh) activeItem.onRefresh();

    if (newValue === activeTab && !deferTabSwitch) return;

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
  }, [activeTab, deferTabSwitch, navigate, pendingTab, safeItems]);

  const contentVariants = useMemo(() => ({
    enter: (dir: number) => ({
      opacity: 0,
      x: isMobile ? (dir > 0 ? 20 : -20) : 0,
      y: isMobile ? 0 : (dir > 0 ? 15 : -15),
      filter: "blur(4px)",
      scale: 0.99,
      pointerEvents: "auto" as const
    }),
    center: { 
      opacity: 1, 
      x: 0, 
      y: 0, 
      filter: "blur(0px)",
      scale: 1,
      pointerEvents: "auto" as const
    },
    exit: (dir: number) => ({
      opacity: 0,
      x: isMobile ? (dir > 0 ? -20 : 20) : 0,
      y: isMobile ? 0 : (dir > 0 ? -15 : 15),
      filter: "blur(4px)",
      scale: 0.99,
      pointerEvents: "none" as const, // Impede que abas saindo bloqueiem cliques
      transition: { duration: 0.2 }
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
      orientation={isMobile ? "horizontal" : "vertical"} 
      className={cn(
        "flex flex-col-reverse md:flex-row w-full h-full",
        "gap-4 md:gap-8 mx-auto items-start overflow-x-hidden",
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
            "absolute bottom-0 left-0 bg-border/50",
            "md:left-0 md:top-0 md:bottom-0 md:right-auto md:w-0.5 md:h-full"
          )}
        />
        <TabsList
          className={cn(
            "flex h-full flex-row justify-around md:flex-col bg-transparent p-0 gap-0",
            "overflow-hidden md:overflow-visible scrollbar-hide w-full relative",
            "transform-gpu md:pr-2 md:py-2",
          )}
        >
            {safeItems.map((tab) => (
              <TabTriggerItem 
                key={tab.value} 
                tab={tab} 
                isActive={activeTab === tab.value} 
              />
            ))}
        </TabsList>
      </div>

      {/* Content */}
      <div className="relative flex-1 w-full h-full overflow-hidden min-w-0">
        <AnimatePresence
          mode="popLayout"
          custom={direction}
          initial={false}
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
              transition={{ 
                duration: 0.3, 
                ease: [0.23, 1, 0.32, 1] 
              }} 
              className={cn("absolute inset-0 h-full w-full", isMobile && "pb-12")}
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
