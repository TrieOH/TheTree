import { SearchInput } from '@/shared/ui/form/SearchInput';
import React, { useState, useMemo, useCallback } from 'react';
import { 
  ChevronDown, 
  ChevronUp, 
  ChevronLeft, 
  ChevronRight,
  MoreHorizontal,
  X,
  Search,
} from 'lucide-react';
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuTrigger 
} from '@/shared/ui/shadcn/dropdown-menu';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/shadcn/select';
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { cn } from '@/shared/lib/utils';
import { AnimatePresence, motion } from 'motion/react';

// --- Type Definitions ---

export type SortDirection = 'asc' | 'desc';

export type FilterType = 'text' | 'select' | 'date' | 'boolean';

export type FilterOption = {
  label: string;
  value: string;
};

export type ColumnDef<T> = {
  key: keyof T;
  header: string;
  responsive?: 'default' | 'icon' | 'hidden';
  sortable?: boolean;
  /** Disables sorting for this column entirely */
  disabled?: boolean;
  primary?: boolean;
  render?: (value: T[keyof T], row: T) => React.ReactNode;
  searchableTextExtractor?: (value: T[keyof T], row: T) => string;
};

export type TableFilter<T> = {
  key: keyof T;
  type: FilterType;
  placeholder?: string;
  options?: FilterOption[];
  label?: string;
};

export type RowAction<T> = {
  label: string;
  onClick: (row: T) => void;
  icon?: React.ElementType;
  variant?: 'default' | 'destructive' | 'ghost' | 'ghost-primary';
  hideLabel?: boolean;
};

export type TableAction = {
  label: string;
  onClick: () => void;
  icon?: React.ElementType;
  variant?: 'default' | 'solid' | 'secondary-solid' | 'outline' | 'ghost' | 'ghost-primary';
};

export type DataTableProps<T> = {
  data: T[];
  columns: ColumnDef<T>[];
  rowActions?: RowAction<T>[];
  tableActions?: TableAction[];
  filters?: TableFilter<T>[];
  idKey?: keyof T;
  itemsPerPage?: number;
  searchPlaceholder?: string;
  initialSort?: { key: keyof T; direction: SortDirection };
  renderExpandedRow?: (row: T) => React.ReactNode;
  /** Force the mobile card view even on desktop */
  forceMobileView?: boolean;
  mobileColumnCount?: number;
};


// --- Sort Indicator Component ---

type SortIndicatorProps = {
  columnKey: string;
  sortConfig: { key: string; direction: SortDirection } | null;
};

const SortIndicator = React.memo(({ columnKey, sortConfig }: Pick<SortIndicatorProps, 'columnKey' | 'sortConfig'>) => {
  const isSorted = sortConfig?.key === columnKey;
  const isAsc = isSorted && sortConfig.direction === 'asc';
  const isDesc = isSorted && sortConfig.direction === 'desc';

  return (
    <div className="flex flex-col items-center justify-center ml-2 pointer-events-none">
      <ChevronUp 
        size={14} 
        className={cn(
          "transition-all duration-200",
          isAsc 
            ? 'text-primary opacity-100' 
            : isDesc
              ? 'text-muted-foreground opacity-30'
              : 'text-muted-foreground opacity-0 group-hover:opacity-40'
        )}
      />
      <ChevronDown 
        size={14} 
        className={cn(
          "-mt-1 transition-all duration-200",
          isDesc 
            ? 'text-primary opacity-100' 
            : isAsc
              ? 'text-muted-foreground opacity-30'
              : 'text-muted-foreground opacity-0 group-hover:opacity-40'
        )}
      />
    </div>
  );
});

SortIndicator.displayName = 'SortIndicator';

const NoResultsState = React.memo(() => (
  <div className="flex flex-col items-center justify-center gap-2 py-8 text-muted-foreground">
    <Search size={32} className="opacity-20" />
    <p>No results found.</p>
  </div>
));

NoResultsState.displayName = 'NoResultsState';


// --- Desktop Row Component ---

interface DesktopRowProps<T> {
  row: T;
  columns: ColumnDef<T>[];
  rowActions?: RowAction<T>[];
  isExpanded: boolean;
  onToggle: (id: string | number) => void;
  renderExpandedRow?: (row: T) => React.ReactNode;
  id: string | number;
}

function DesktopRowInner<T extends object>({ 
  row, 
  columns, 
  rowActions, 
  isExpanded, 
  onToggle, 
  renderExpandedRow,
  id 
}: DesktopRowProps<T>) {
  return (
    <React.Fragment>
      <tr 
        className={cn(
          "transition-colors hover:bg-muted/70",
          renderExpandedRow && "cursor-pointer"
        )}
        onClick={() => renderExpandedRow && onToggle(id)}
      >
        {columns.map((col, colIndex) => (
          <td key={String(col.key)} className="p-4 align-middle whitespace-nowrap">
            <div className="flex items-center gap-3">
              {renderExpandedRow && colIndex === 0 && (
                <div className="transition-transform duration-200" style={{ transform: isExpanded ? 'rotate(90deg)' : 'rotate(0deg)' }}>
                  <ChevronRight size={16} />
                </div>
              )}
              <span className="flex-1">
                {col.render 
                  ? col.render(row[col.key], row) 
                  : String(row[col.key] ?? '-')
                }
              </span>
            </div>
          </td>
        ))}
        {rowActions && rowActions.length > 0 && (
          <td className="p-4 flex items-center gap-2 text-right whitespace-nowrap">
            {rowActions.map((action, idx) => {
              const Icon = action.icon;
              if(!Icon) return null;
              return (
                <ShadowButton
                  key={`${action.label}${idx}`}
                  leftIcon={<Icon size={16} className="" />}
                  label={action.label}
                  onClick={(e) => {
                    e.stopPropagation();
                    action.onClick(row);
                  }}
                  variant={action.variant}
                />
              );
            })}
          </td>
        )}
      </tr>
      
      <tr>
        <td colSpan={columns.length + (rowActions ? 1 : 0)} className="p-0 border-none">
          <AnimatePresence initial={false}>
            {isExpanded && renderExpandedRow && (
              <motion.div
                initial={{ height: 0, opacity: 0 }}
                animate={{ height: "auto", opacity: 1 }}
                exit={{ height: 0, opacity: 0 }}
                transition={{ duration: 0.2, ease: "easeInOut" }}
                className="overflow-hidden bg-muted/30"
              >
                <div className="p-4">
                  {renderExpandedRow(row)}
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </td>
      </tr>
    </React.Fragment>
  );
}

const DesktopRow = React.memo(DesktopRowInner) as typeof DesktopRowInner;

// --- Mobile Card Component ---

interface MobileCardProps<T> {
  row: T;
  primaryColumn: ColumnDef<T>;
  mobileVisibleColumns: ColumnDef<T>[];
  mobileColumnCount: number;
  isExpanded: boolean;
  onToggle: (id: string | number) => void;
  renderExpandedRow?: (row: T) => React.ReactNode;
  rowActions?: RowAction<T>[];
  id: string | number;
}

function MobileCardInner<T extends object>({
  row,
  primaryColumn,
  mobileVisibleColumns,
  mobileColumnCount,
  isExpanded,
  onToggle,
  renderExpandedRow,
  rowActions,
  id
}: MobileCardProps<T>) {
  const interactiveProps = renderExpandedRow ? {
    onClick: () => onToggle(id),
    onKeyDown: (e: React.KeyboardEvent) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        onToggle(id);
      }
    },
    role: "button" as const,
    tabIndex: 0
  } : {};

  return (
    <div className='relative'>
      <div 
        className={cn(
          "rounded-md border border-border bg-background shadow-sm text-left",
          "w-full overflow-hidden"
        )}
      >
        <div 
          className={cn("p-4", renderExpandedRow && "cursor-pointer")} 
          {...interactiveProps}
        >
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1 min-w-0">
              <h3 className="font-semibold text-foreground truncate">
                {primaryColumn.render 
                  ? primaryColumn.render(row[primaryColumn.key], row)
                  : String(row[primaryColumn.key])
                }
              </h3>
            </div>
          </div>

          <div className="mt-4 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {mobileVisibleColumns
              .filter(c => String(c.key) !== String(primaryColumn.key))
              .slice(0, mobileColumnCount)
              .map((col) => (
                <div key={String(col.key)} className="flex flex-col justify-center">
                  <span className="text-xs text-muted-foreground uppercase tracking-wider truncate">{col.header}</span>
                  <span className="text-sm font-medium text-foreground truncate">
                    {col.render 
                      ? col.render(row[col.key], row) 
                      : String(row[col.key] ?? '-')
                    }
                  </span>
                </div>
              ))
            }
          </div>
        </div>
        
        <AnimatePresence initial={false}>
          {isExpanded && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              transition={{ duration: 0.2, ease: "easeInOut" }}
              className="overflow-hidden border-t border-border bg-muted/30"
            >
              <div className="p-4">
                {renderExpandedRow 
                  ? renderExpandedRow(row)
                  : (
                    <div className="space-y-3">
                      {mobileVisibleColumns
                        .slice(mobileColumnCount) 
                        .map((col) => (
                        <div key={String(col.key)} className="flex flex-col">
                          <span className="text-xs text-muted-foreground uppercase tracking-wider truncate">{col.header}</span>
                          <span className="text-sm font-medium text-foreground">
                            {col.render 
                              ? col.render(row[col.key], row) 
                              : String(row[col.key] ?? '-')
                            }
                          </span>
                        </div>
                      ))}
                    </div>
                  )
                }
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
      {renderExpandedRow && (
        <ChevronDown 
          size={18} 
          className={cn(
            "transition-transform", isExpanded && "rotate-180",
            "absolute top-6 right-12 pointer-events-none"
          )} 
        />
      )}
      {rowActions && (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button 
              type="button" 
              className="p-2 text-muted-foreground hover:bg-muted rounded absolute top-4 right-4"
              onClick={(e) => e.stopPropagation()}
            >
              <MoreHorizontal size={18} />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-40 bg-popover rounded-md shadow-lg z-50">
            {rowActions.map((action, idx) => (
              <DropdownMenuItem 
                key={`${action.label}${idx}`} 
                onClick={(e) => {
                  e.stopPropagation();
                  action.onClick(row);
                }}
                className={cn(
                  "flex cursor-pointer items-center px-2 py-1.5 text-sm transition-colors",
                  "hover:bg-muted hover:text-muted-foreground",
                  action.variant === 'destructive' && "text-destructive"
                )}
              >
                {action.icon && <action.icon size={16} className="mr-2" />}
                {action.label}
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      )}
    </div>
  );
}

const MobileCard = React.memo(MobileCardInner) as typeof MobileCardInner;


// --- Main Component ---

export default function CustomDataTable<T extends object>({
  data,
  columns,
  rowActions,
  tableActions,
  filters,
  idKey = 'id' as keyof T,
  itemsPerPage = 10,
  searchPlaceholder = "Search...",
  initialSort,
  renderExpandedRow,
  forceMobileView = false,
  mobileColumnCount = 3,
}: DataTableProps<T>) {
  
  const defaultSortColumn = columns.find(col => col.primary && col.sortable) || columns.find(col => col.sortable);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortConfig, setSortConfig] = useState<{ key: keyof T; direction: SortDirection } | null>(
    initialSort || (defaultSortColumn ? { key: defaultSortColumn.key, direction: 'asc' } : null)
  );
  const [currentPage, setCurrentPage] = useState(1);
  const [activeFilters, setActiveFilters] = useState<Partial<Record<keyof T, string>>>({});
  const [expandedRows, setExpandedRows] = useState<Set<string | number>>(new Set());

  const toggleRow = useCallback((id: string | number) => {
    setExpandedRows(prev => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }, []);

  const handleSort = useCallback((key: keyof T, disabled?: boolean) => {
    if (disabled) return;
    
    setSortConfig(prev => {
      let direction: SortDirection = 'asc';
      if (prev?.key === key && prev.direction === 'asc') direction = 'desc';
      return { key, direction };
    });
    setCurrentPage(1);
  }, []);

  const clearFilters = useCallback(() => {
    setSearchTerm('');
    setActiveFilters({});
    setCurrentPage(1);
  }, []);

  const hasActiveFilters = searchTerm || Object.values(activeFilters).some(v => v !== undefined && v !== '');

  const filteredData = useMemo(() => {
    const term = searchTerm.toLowerCase();
    const filterEntries = Object.entries(activeFilters).filter(([_, v]) => v);

    return data.filter((item) => {
      const matchesSearch = !term || columns.some(col => {
        const val = col.searchableTextExtractor 
          ? col.searchableTextExtractor(item[col.key], item)
          : String(item[col.key] ?? '');
        return val.toLowerCase().includes(term);
      });

      if (!matchesSearch) return false;

      return filterEntries.every(([key, filterValue]) => {
        const itemValue = String(item[key as keyof T]).toLowerCase();
        const fv = String(filterValue).toLowerCase();
        return itemValue === fv || itemValue.includes(fv);
      });
    });
  }, [data, searchTerm, activeFilters, columns]);

  const sortedData = useMemo(() => {
    if (!sortConfig) return filteredData;
    
    const { key, direction } = sortConfig;
    const isAsc = direction === 'asc';

    return [...filteredData].sort((a, b) => {
      const aVal = a[key];
      const bVal = b[key];

      if (aVal == null) return 1;
      if (bVal == null) return -1;
      if (aVal === bVal) return 0;

      if (typeof aVal === 'string' && typeof bVal === 'string') {
        return isAsc ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal);
      }
      
      return isAsc ? (aVal > bVal ? 1 : -1) : (aVal > bVal ? -1 : 1);
    });
  }, [filteredData, sortConfig]);

  const totalPages = Math.ceil(sortedData.length / itemsPerPage);
  const paginatedData = sortedData.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  if (currentPage > totalPages && totalPages > 0) {
    setCurrentPage(totalPages);
  }

  const primaryColumn = columns.find(c => c.primary) || columns[0];
  const mobileVisibleColumns = columns.filter(c => c.responsive !== 'hidden');

  return (
    <div className="w-full space-y-4">
      
      {/* Toolbar */}
      <div className="flex flex-col gap-4">
        
        {/* Top row: Search + Filters + Clear + Actions */}
        <div className="flex flex-wrap items-center gap-2">
          {/* Search */}
          <div className="relative w-full sm:w-64">
            <SearchInput
              value={searchTerm}
              onChange={(value) => {
                setSearchTerm(value);
                setCurrentPage(1);
              }}
              placeholder={searchPlaceholder}
              className="w-full"
            />
          </div>

          {/* Filters */}
          {filters?.map((filter) => (
            <div key={String(filter.key)}>
              {filter.type === 'select' ? (
                <Select
                  value={activeFilters[filter.key] || 'all'}
                  onValueChange={(value) => {
                    setActiveFilters(prev => ({ ...prev, [filter.key]: value === 'all' ? '' : value }));
                    setCurrentPage(1);
                  }}
                >
                  <SelectTrigger size="sm" className="w-45 h-full!">
                    <SelectValue placeholder={filter.placeholder || `All ${String(filter.key)}`} />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    <SelectItem value="all">{filter.placeholder || `All ${String(filter.key)}`}</SelectItem>
                    {filter.options?.map(opt => (
                      <SelectItem key={opt.value} value={opt.value}>{opt.label}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              ) : (
                <input
                  type={filter.type === 'date' ? 'date' : 'text'}
                  placeholder={filter.placeholder}
                  className="h-9 rounded-md border border-input bg-background px-3 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  value={activeFilters[filter.key] || ''}
                  onChange={(e) => {
                    setActiveFilters(prev => ({ ...prev, [filter.key]: e.target.value }));
                    setCurrentPage(1);
                  }}
                />
              )}
            </div>
          ))}

          {/* Clear button */}
          {hasActiveFilters && (
            <ShadowButton
              onClick={clearFilters}
              leftIcon={<X size={14} />}
              value="Clear"
              variant="outline"
            />
          )}

          <div className="hidden lg:block lg:flex-1" />

          {/* Table Actions */}
          {tableActions?.map((action, idx) => {
            const Icon = action.icon;
            return (
              <ShadowButton
                key={`${action.label}${idx}`}
                onClick={action.onClick}
                variant={action.variant}
                label={action.label}
                value={action.label}
                leftIcon={Icon ? <Icon size={16} /> : undefined}
              />
            );
          })}
        </div>
      </div>

      {/* Desktop Table */}
      <div className={cn(
        "rounded-md border border-border bg-card shadow-sm",
        forceMobileView ? "hidden" : "hidden md:block"
      )}>
        <div className="overflow-x-auto">
          <table className="w-full caption-bottom text-sm">
            <thead className="border-b border-border bg-muted/60">
              <tr>
                {columns.map((col) => {
                  const isSortable = col.sortable && !col.disabled;
                  const isDisabled = col.disabled;
                  
                  return (
                    <th
                      key={String(col.key)}
                      className={cn(
                        "h-12 px-4 text-left align-middle text-xs font-medium text-muted-foreground whitespace-nowrap",
                        isSortable && "cursor-pointer select-none hover:text-foreground",
                        isDisabled && "cursor-not-allowed opacity-60",
                        isSortable && "group"
                      )}
                      onClick={() => handleSort(col.key, col.disabled)}
                    >
                      <div className="flex items-center justify-between">
                        {/* Label - Left aligned */}
                        <span className={cn(isSortable && 'pr-2')}>
                          {col.responsive === 'icon' ? (
                            <>
                              <span className="hidden lg:inline">{col.header}</span>
                              <span className="lg:hidden" title={col.header}>
                                {col.header.charAt(0)}
                              </span>
                            </>
                          ) : (
                            col.header
                          )}
                        </span>
                        
                        {/* Sort Indicator - Right aligned */}
                        {col.sortable && (
                          <SortIndicator 
                            columnKey={String(col.key)}
                            sortConfig={sortConfig ? { key: String(sortConfig.key), direction: sortConfig.direction } : null}
                          />
                        )}
                      </div>
                    </th>
                  );
                })}
                {rowActions && rowActions.length > 0 && (
                  <th className="h-12 px-4 text-right align-middle font-medium text-muted-foreground whitespace-nowrap">
                  </th>
                )}
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {paginatedData.length > 0 ? (
                paginatedData.map((row, rowIndex) => {
                  const id = String(row[idKey] ?? rowIndex);
                  return (
                    <DesktopRow
                      key={id}
                      id={id}
                      row={row}
                      columns={columns}
                      rowActions={rowActions}
                      isExpanded={expandedRows.has(id)}
                      onToggle={toggleRow}
                      renderExpandedRow={renderExpandedRow}
                    />
                  );
                })
              ) : (
                <tr>
                  <td 
                    colSpan={columns.length + (rowActions ? 1 : 0)} 
                    className="h-32 text-center align-middle text-muted-foreground"
                  >
                    <NoResultsState />
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Mobile Card View */}
      <div className={cn(
        "space-y-3",
        forceMobileView ? "block" : "md:hidden"
      )}>
        {paginatedData.length > 0 ? (
          paginatedData.map((row, rowIndex) => {
            const id = String(row[idKey] ?? rowIndex);
            return (
              <MobileCard
                key={id}
                id={id}
                row={row}
                primaryColumn={primaryColumn}
                mobileVisibleColumns={mobileVisibleColumns}
                mobileColumnCount={mobileColumnCount}
                isExpanded={expandedRows.has(id)}
                onToggle={toggleRow}
                renderExpandedRow={renderExpandedRow}
                rowActions={rowActions}
              />
            );
          })
        ) : (
          <div className="rounded-md border border-border bg-background p-4 shadow-sm">
            <NoResultsState />
          </div>
        )}
      </div>

      {/* Pagination Footer */}
      <div className="flex flex-col items-center justify-between gap-4 border-t border-border bg-card px-4 py-3 sm:flex-row">
        <span className="text-sm text-muted-foreground">
          Results <strong>{Math.min((currentPage - 1) * itemsPerPage + 1, sortedData.length)}-{Math.min(currentPage * itemsPerPage, sortedData.length)}</strong> of <strong>{sortedData.length}</strong>
        </span>
        
        <div className="flex items-center gap-2">
          <ShadowButton
            onClick={() => setCurrentPage(p => Math.max(1, p - 1))} 
            disabled={currentPage === 1}
            leftIcon={<ChevronLeft size={18}/>}
            variant="outline"
            label="Previous Page"
          />
          <div className="hidden sm:flex items-center gap-1">
            {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
              const pageNum = i + 1;
              const isActive = pageNum === currentPage;
              return (
                <ShadowButton
                  key={pageNum}
                  onClick={() => setCurrentPage(pageNum)}
                  value={pageNum.toString()}
                  variant={isActive ? "default" : "outline"}
                  className="h-9 w-9 justify-center"
                />
              );
            })}
            {totalPages > 5 && <span className="px-2 text-muted-foreground">...</span>}
          </div>

          <span className="text-sm text-muted-foreground sm:hidden">
            Page {currentPage} / {totalPages}
          </span>
          <ShadowButton
            onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))} 
            leftIcon={<ChevronRight size={18}/>}
            disabled={currentPage === totalPages || totalPages === 0} 
            variant="outline"
            label="Next Page"
          />
        </div>
      </div>
    </div>
  );
}
