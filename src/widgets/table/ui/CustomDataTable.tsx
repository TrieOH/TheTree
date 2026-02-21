import { SearchInput } from '@/shared/ui/form/SearchInput';
import React, { useState, useMemo } from 'react';
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
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { cn } from '@/shared/lib/utils';

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
};


// --- Sort Indicator Component ---

type SortIndicatorProps = {
  columnKey: string;
  sortConfig: { key: string; direction: SortDirection } | null;
};

const SortIndicator = ({ columnKey, sortConfig }: Pick<SortIndicatorProps, 'columnKey' | 'sortConfig'>) => {
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
};

// --- No Results State Component ---

const NoResultsState = () => (
  <div className="flex flex-col items-center justify-center gap-2 py-8 text-muted-foreground">
    <Search size={32} className="opacity-20" />
    <p>No results found.</p>
  </div>
);


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
}: DataTableProps<T>) {
  
  const defaultSortColumn = columns.find(col => col.primary && col.sortable) || columns.find(col => col.sortable);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortConfig, setSortConfig] = useState<{ key: keyof T; direction: SortDirection } | null>(
    initialSort || (defaultSortColumn ? { key: defaultSortColumn.key, direction: 'asc' } : null)
  );
  const [currentPage, setCurrentPage] = useState(1);
  const [activeFilters, setActiveFilters] = useState<Partial<Record<keyof T, string>>>({});
  const [expandedRows, setExpandedRows] = useState<Set<string | number>>(new Set());

  const toggleRow = (id: string | number) => {
    const newExpandedRows = new Set(expandedRows);
    if (newExpandedRows.has(id)) newExpandedRows.delete(id);
    else newExpandedRows.add(id);
    setExpandedRows(newExpandedRows);
  };

  const handleSort = (key: keyof T, disabled?: boolean) => {
    if (disabled) return;
    
    let direction: SortDirection = 'asc';
    if (sortConfig?.key === key && sortConfig.direction === 'asc') direction = 'desc';
    setSortConfig({ key, direction });
    setCurrentPage(1);
  };

  const clearFilters = () => {
    setSearchTerm('');
    setActiveFilters({});
    setCurrentPage(1);
  };

  const hasActiveFilters = searchTerm || Object.values(activeFilters).some(v => v !== undefined && v !== '');

    const filteredData = useMemo(() => {
      return data.filter((item) => {
        const itemSearchableString = columns.reduce((acc, col) => {
          let cellValue: string = '';
          if (col.searchableTextExtractor) cellValue = col.searchableTextExtractor(item[col.key], item);
          else cellValue = String(item[col.key] ?? '');

          return `${acc} ${cellValue}`;
        }, "").toLowerCase();
  
        const matchesSearch = !searchTerm || itemSearchableString.includes(searchTerm.toLowerCase());
      const matchesCustomFilters = Object.entries(activeFilters).every(([key, filterValue]) => {
        if (!filterValue) return true;
        const itemValue = String(item[key as keyof T]).toLowerCase();
        return itemValue === String(filterValue).toLowerCase() || itemValue.includes(String(filterValue).toLowerCase());
      });

      return matchesSearch && matchesCustomFilters;
    });
  }, [data, searchTerm, activeFilters, columns]);

  const sortedData = useMemo(() => {
    if (!sortConfig) return filteredData;
    
    return [...filteredData].sort((a, b) => {
      const aVal = a[sortConfig.key];
      const bVal = b[sortConfig.key];

      if (aVal == null) return 1;
      if (bVal == null) return -1;
      if (aVal === bVal) return 0;

      if (typeof aVal === 'string' && typeof bVal === 'string') {
        return sortConfig.direction === 'asc' 
          ? aVal.localeCompare(bVal) 
          : bVal.localeCompare(aVal);
      }
      
      return sortConfig.direction === 'asc' ? (aVal > bVal ? 1 : -1) : (aVal > bVal ? -1 : 1);
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

          {/* Filters - TODO: I need to make a custom for select and input */}
          {filters?.map((filter) => (
            <div key={String(filter.key)}>
              {filter.type === 'select' ? (
                <select
                  className="h-9 rounded-md border border-input bg-background px-3 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  value={activeFilters[filter.key] || ''}
                  onChange={(e) => {
                    setActiveFilters(prev => ({ ...prev, [filter.key]: e.target.value }));
                    setCurrentPage(1);
                  }}
                >
                  <option value="">{filter.placeholder || `All ${String(filter.key)}`}</option>
                  {filter.options?.map(opt => (
                    <option key={opt.value} value={opt.value}>{opt.label}</option>
                  ))}
                </select>
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

          {/* Spacer to push actions right on desktop */}
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
      <div className="hidden rounded-md border border-border bg-card shadow-sm md:block">
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
                  const isExpanded = expandedRows.has(id);
                  
                  return (
                  <React.Fragment key={id}>
                    <tr 
                      onClick={() => renderExpandedRow && toggleRow(id)}
                      className={cn(
                        "transition-colors hover:bg-muted/70",
                        renderExpandedRow && "cursor-pointer"
                      )}
                    >
                      {columns.map((col, colIndex) => (
                        <td key={String(col.key)} className="p-4 align-middle whitespace-nowrap">
                          <div className="flex items-center gap-3">
                            {renderExpandedRow && colIndex === 0 && (
                              <button
                                type='button'
                                onClick={(e) => {
                                  e.stopPropagation();
                                  toggleRow(id);
                                }}
                                className="flex items-center justify-center"
                              >
                                {isExpanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                              </button>
                            )}
                            {col.render 
                              ? col.render(row[col.key], row) 
                              : String(row[col.key] ?? '-')
                            }
                          </div>
                        </td>
                      ))}
                      {rowActions && rowActions.length > 0 && (
                        <td className="p-4 align-middle text-right whitespace-nowrap">
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <ShadowButton
                                leftIcon={<MoreHorizontal size={16} />}
                                label='More Actions'
                                variant="ghost-primary"
                              />
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end" className="w-40">
                              {rowActions.map((action, idx) => {
                                const Icon = action.icon;
                                return (
                                  <DropdownMenuItem
                                    key={`${action.label}${idx}`}
                                    onClick={(e) => {
                                    e.stopPropagation();
                                    action.onClick(row);
                                  }}
                                    className="cursor-pointer"
                                  >
                                    {Icon && React.createElement(Icon, { size: 16, className: "mr-2" })}
                                    <span>{action.label}</span>
                                  </DropdownMenuItem>
                                );
                              })}
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </td>
                      )}
                    </tr>
                    {isExpanded && renderExpandedRow && (
                      <tr className="bg-muted/30">
                        <td colSpan={columns.length + (rowActions ? 1 : 0)}>
                          {renderExpandedRow(row)}
                        </td>
                      </tr>
                    )}
                  </React.Fragment>
                )})
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
      <div className="space-y-3 md:hidden">
        {paginatedData.length > 0 ? (
          paginatedData.map((row, rowIndex) => {
            const id = String(row[idKey] ?? rowIndex);
            const isExpanded = expandedRows.has(id);
            
            return (
            <button 
              key={id} 
              type='button'
              className={cn(
                "rounded-md border border-border bg-background p-4 shadow-sm text-left w-full",
                renderExpandedRow && "cursor-pointer"
              )}
              onClick={() => renderExpandedRow && toggleRow(id)}
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1 min-w-0">
                  <h3 className="font-semibold text-foreground truncate">
                    {primaryColumn.render 
                      ? primaryColumn.render(row[primaryColumn.key], row)
                      : String(row[primaryColumn.key])
                    }
                  </h3>
                  {mobileVisibleColumns[1] && (
                    <div className="text-sm text-muted-foreground mt-1 truncate">
                      {mobileVisibleColumns[1].render 
                        ? mobileVisibleColumns[1].render(row[mobileVisibleColumns[1].key], row)
                        : String(row[mobileVisibleColumns[1].key])
                      }
                    </div>
                  )}
                </div>
                
                <div className="flex shrink-0 items-center gap-1">
                    {renderExpandedRow && (
                      <button
                        type='button'
                        onClick={(e) => {
                          e.stopPropagation(); 
                          toggleRow(id);
                        }}
                        className="p-2 text-muted-foreground"
                      >
                        <ChevronDown size={18} className={cn("transition-transform", isExpanded && "rotate-180")} />
                      </button>
                    )}
                    {rowActions && (
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <button type="button" className="p-2 text-muted-foreground">
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
                              {action.icon && React.createElement(action.icon, { size: 16, className: "mr-2" })}
                              {action.label}
                            </DropdownMenuItem>
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    )}
                  </div>
              </div>

              <div className="mt-4 grid grid-cols-2 gap-3 border-t border-border pt-3">
                {mobileVisibleColumns
                  .filter(c => String(c.key) !== String(primaryColumn.key))
                  .slice(0, 3) // Show first 3 after primary
                  .map((col) => (
                    <div key={String(col.key)} className="flex flex-col">
                      <span className="text-xs text-muted-foreground uppercase tracking-wider truncate">{col.header}</span>
                      <span className="text-sm font-medium text-foreground truncate">
                        {col.render 
                          ? col.render(row[col.key], row) 
                          : String(row[col.key] ?? '-')
                        }
                      </span>
                    </div>
                  ))}
              </div>
              
              {isExpanded && (
                <button
                  type='button'
                  className="mt-4 border-t border-border pt-4 w-full"
                  onClick={(e) => e.stopPropagation()}
                >
                  {renderExpandedRow 
                    ? renderExpandedRow(row)
                    : (
                      <div className="space-y-3">
                        {mobileVisibleColumns
                          .slice(3) // Show the rest of the columns
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
                </button>
              )}
            </button>
          )})
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
