import { useEffect, useState } from 'react';
import { Tabs, TabsList, TabsTrigger } from '@/shared/ui/shadcn/tabs';
import { AnimatePresence, motion } from 'motion/react';
import { cn } from '@/shared/lib/utils';
import { useNavigate } from '@tanstack/react-router';

type TabItem = {
  value: string;
  label: string;
  icon?: React.ElementType;
  content: React.ReactNode;
};

type Props = {
  items: TabItem[];
  initialValue?: string;
  deferTabSwitch?: boolean;
  className?: string;
};

export default function CustomTabs({
  items,
  initialValue,
  deferTabSwitch = false,
  className = ''
}: Props) {
  const safeItems = items ?? [];

  const first = initialValue ?? safeItems[0]?.value ?? '';

  const [activeTab, setActiveTab] = useState<string>(first)
  const [displayTab, setDisplayTab] = useState<string | null>(first);
  const [pendingTab, setPendingTab] = useState<string | null>(null);
  const [direction, setDirection] = useState<number>(0);
  const [isMobile, setIsMobile] = useState<boolean>(false);
  const navigate = useNavigate();

  useEffect(() => {
    const check = () => setIsMobile(window.innerWidth < 768);
    check();
    window.addEventListener('resize', check);
    return () => window.removeEventListener('resize', check);
  }, []);

  if (safeItems.length === 0) return null;

  const idxOf = (v: string) => items.findIndex((t) => t.value === v);

  const handleTabChange = (newValue: string) => {
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
    })
  };

  const contentVariants = {
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
  };

  return (
    <Tabs 
      value={activeTab} 
      onValueChange={(v) => handleTabChange(v)} 
      orientation="vertical" 
      className={cn(
        "flex flex-col-reverse md:flex-row w-full max-w-7xl h-full",
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
          {items.map((tab) => {
            const isActive = activeTab === tab.value;
            const Icon = tab.icon ?? (() => null);
            return (
              <TabsTrigger
                key={tab.value}
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
                <span className="relative z-20 flex items-center gap-2">
                  <Icon className="h-4 w-4 shrink-0" />
                  <span className="hidden sm:inline">{tab.label}</span>
                  <span className="sm:hidden text-xs">{tab.label.split(' ')[0]}</span>
                </span>
              </TabsTrigger>
            );
          })}
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
                {(() => {
                  const item = items.find((t) => t.value === displayTab);
                  if (!item) return null;
                  return item.content;
                })()}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </Tabs>
  );
}
