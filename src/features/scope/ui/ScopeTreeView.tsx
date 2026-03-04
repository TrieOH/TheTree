import { useMemo, useState } from 'react'
import { Edit, Trash2, ChevronRight, ChevronDown, PlusCircle, ChevronUp, MoreHorizontal } from 'lucide-react'
import { Badge } from "@/shared/ui/shadcn/badge"
import { MetadataVisualizer } from "@/shared/ui/MetadataVisualizer"
import TruncatedId from "@/shared/ui/TruncatedId"
import { formatDate } from "@/shared/lib/date-utils"
import { scopeActions } from "../store"
import type { Scope } from "../model/types"
import { SearchInput } from "@/shared/ui/form/SearchInput"
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton"
import { cn } from "@/shared/lib/utils"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/ui/shadcn/dropdown-menu';

interface ScopeNode extends Scope {
  children: ScopeNode[]
  isMatch?: boolean
}

interface ScopeTreeViewProps {
  scopes: Scope[]
}

type SortDirection = 'asc' | 'desc';
type SortConfig = { key: keyof Scope; direction: SortDirection } | null;

const GRID_COLS_CLASS = "grid grid-cols-[minmax(400px,1fr)_140px_180px_180px_160px_80px] items-center"

const SortIndicator = ({ isAsc, isDesc }: { isAsc: boolean; isDesc: boolean }) => (
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

export default function ScopeTreeView({ scopes }: ScopeTreeViewProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set())
  const [sortConfig, setSortConfig] = useState<SortConfig>({ key: 'name', direction: 'asc' })

  const toggleExpand = (id: string, e?: React.MouseEvent) => {
    e?.stopPropagation()
    const newExpanded = new Set(expandedIds)
    if (newExpanded.has(id)) {
      newExpanded.delete(id)
    } else {
      newExpanded.add(id)
    }
    setExpandedIds(newExpanded)
  }

  const handleSort = (key: keyof Scope) => {
    setSortConfig(prev => {
      if (prev?.key === key) {
        return { key, direction: prev.direction === 'asc' ? 'desc' : 'asc' };
      }
      return { key, direction: 'asc' };
    });
  }

  const processedData = useMemo(() => {
    const scopeMap: Record<string, ScopeNode> = {}
    const sortedScopes = [...scopes].sort((a, b) => {
      if (!sortConfig) return 0;
      const { key, direction } = sortConfig;
      const aVal = String(a[key] || '').toLowerCase();
      const bVal = String(b[key] || '').toLowerCase();
      return direction === 'asc' ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal);
    });

    sortedScopes.forEach(scope => {
      scopeMap[scope.id] = { ...scope, children: [], isMatch: false }
    })

    const roots: ScopeNode[] = []
    const normalizedSearch = searchTerm.toLowerCase().trim()

    sortedScopes.forEach(scope => {
      const node = scopeMap[scope.id]
      if (normalizedSearch) {
        const matchesName = node.name.toLowerCase().includes(normalizedSearch)
        const matchesId = node.id.toLowerCase().includes(normalizedSearch)
        const matchesExternalId = node.external_id?.toLowerCase().includes(normalizedSearch)
        const matchesStatus = node.meta?.status?.toLowerCase().includes(normalizedSearch)
        if (matchesName || matchesId || matchesExternalId || matchesStatus) {
          node.isMatch = true
        }
      }

      if (scope.parent_id && scopeMap[scope.parent_id]) {
        scopeMap[scope.parent_id].children.push(node)
      } else {
        roots.push(node)
      }
    })

    const filterTree = (nodes: ScopeNode[]): ScopeNode[] => {
      return nodes.reduce<ScopeNode[]>((acc, node) => {
        const filteredChildren = filterTree(node.children)
        const hasMatchInDescendants = filteredChildren.length > 0
        if (node.isMatch || hasMatchInDescendants) {
          acc.push({ ...node, children: filteredChildren })
          if (hasMatchInDescendants && normalizedSearch) {
            setExpandedIds(prev => new Set(prev).add(node.id))
          }
        }
        return acc
      }, [])
    }

    return normalizedSearch ? filterTree(roots) : roots
  }, [scopes, searchTerm, sortConfig])

  const renderNode = (node: ScopeNode, level: number = 0) => {
    const isExpanded = expandedIds.has(node.id) || (searchTerm && node.children.length > 0)
    const hasChildren = node.children.length > 0

    return (
      <div key={node.id} className="w-full">
        <div 
          onClick={(e) => hasChildren && toggleExpand(node.id, e)}
          className={cn(
            "group hover:bg-muted/70 border-b border-border transition-colors cursor-pointer",
            GRID_COLS_CLASS,
            node.isMatch && "bg-indigo-50/50"
          )}
        >
          <div className="flex items-center gap-3 min-w-0 p-4" style={{ paddingLeft: `${level * 16 + 16}px` }}>
            <button
              type='button'
              onClick={(e) => toggleExpand(node.id, e)}
              className={cn(
                "w-4 h-4 flex items-center justify-center rounded transition-all shrink-0",
                !hasChildren ? "invisible" : isExpanded ? "text-muted-foreground" : "text-primary hover:bg-primary/10"
              )}
            >
              {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
            </button>
            <div className="flex-1 min-w-0">
              <MetadataVisualizer name={node.name} meta={node.meta} />
            </div>
          </div>

          <div className="p-4 text-sm">
            <ScopeBadge type={node.type} />
          </div>

          {node.external_id ? <TruncatedId id={node.external_id} /> : <span className="text-muted-foreground/50">-</span>}

          <TruncatedId id={node.id} />

          <div className="p-4 text-sm text-muted-foreground whitespace-nowrap">
            {formatDate(node.created_at)}
          </div>

          <div className="p-4 flex justify-end">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <ShadowButton
                  leftIcon={<MoreHorizontal size={16} />}
                  label='More Actions'
                  variant="ghost-primary"
                  onClick={(e) => e.stopPropagation()}
                />
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-40">
                <DropdownMenuItem onClick={(e) => { e.stopPropagation(); scopeActions.openCreate({ parent_id: node.id }); }} className="cursor-pointer">
                  <PlusCircle size={16} className="mr-2" />
                  <span>Add Child</span>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={(e) => { e.stopPropagation(); scopeActions.openEdit(node); }} className="cursor-pointer">
                  <Edit size={16} className="mr-2 text-primary" />
                  <span>Update</span>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={(e) => { e.stopPropagation(); scopeActions.openDelete(node); }} className="cursor-pointer text-destructive focus:text-destructive">
                  <Trash2 size={16} className="mr-2" />
                  <span>Delete</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
        
        {hasChildren && isExpanded && (
          <div className="w-full">
            {node.children.map(child => renderNode(child, level + 1))}
          </div>
        )}
      </div>
    )
  }

  return (
    <div className="w-full space-y-4 text-sm">
      <div className="flex items-center justify-between gap-4">
        <div className="flex-1 max-w-sm">
          <SearchInput
            placeholder="Search scopes..."
            value={searchTerm}
            onChange={(val) => setSearchTerm(val)}
          />
        </div>
        <ShadowButton 
          variant="solid"
          value="Create Scope"
          onClick={() => scopeActions.openCreate()}
          leftIcon={<PlusCircle className="mr-2 h-4 w-4" />}
        />
      </div>

      <div className="rounded-md border border-border bg-card shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <div className="min-w-max">
            <div className={cn(
              "sticky top-0 z-10 h-12 border-b border-border bg-muted/60 backdrop-blur-sm text-xs font-medium text-muted-foreground whitespace-nowrap select-none",
              GRID_COLS_CLASS
            )}>
              <HeaderItem label="Scope Identity" sortKey="name" currentSort={sortConfig} onSort={handleSort} />
              <HeaderItem label="Category" sortKey="type" currentSort={sortConfig} onSort={handleSort} />
              <HeaderItem label="External ID" sortKey="external_id" currentSort={sortConfig} onSort={handleSort} />
              <HeaderItem label="ID" sortKey="id" currentSort={sortConfig} onSort={handleSort} />
              <HeaderItem label="Created At" sortKey="created_at" currentSort={sortConfig} onSort={handleSort} />
              <div className="px-4" /> {/* Actions header spacer */}
            </div>

            {/* Body */}
            <div className="divide-y divide-border max-h-175 overflow-y-auto">
              {processedData.length > 0 ? (
                processedData.map(node => renderNode(node))
              ) : (
                <div className="py-20 text-center text-muted-foreground min-w-250">
                  <p>No results found.</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function HeaderItem({ label, sortKey, currentSort, onSort }: { 
  label: string; 
  sortKey: keyof Scope; 
  currentSort: SortConfig; 
  onSort: (key: keyof Scope) => void 
}) {
  const isSorted = currentSort?.key === sortKey;
  const isAsc = isSorted && currentSort?.direction === 'asc';
  const isDesc = isSorted && currentSort?.direction === 'desc';

  return (
    <button 
      type='button'
      className={cn(
        "flex items-center justify-between h-full px-4",
        "cursor-pointer hover:text-foreground transition-colors group"
      )}
      onClick={() => onSort(sortKey)}
    >
      <span>{label}</span>
      <SortIndicator isAsc={isAsc} isDesc={isDesc} />
    </button>
  )
}

function ScopeBadge({ type }: { type: string }) {
  let variant: "default" | "secondary" | "destructive" | "outline" = "default"
  let displayType = type
  
  switch (type) {
    case "global": variant = "outline"; break;
    case "project_root": variant = "secondary"; displayType = "Root"; break;
    case "project_scope": variant = "default"; displayType = "Scope"; break;
  }
  
  return <Badge variant={variant} className="capitalize text-[10px] px-2 h-5 font-bold">{displayType}</Badge>
}
