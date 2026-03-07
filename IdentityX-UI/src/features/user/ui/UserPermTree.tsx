import * as Icons from "lucide-react";
import { useState } from 'react'
import type { Node, NodeCustomName } from '../model/types'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import { cn } from '@/shared/lib/utils'
import { Folder, FolderOpen, type LucideIcon } from 'lucide-react'

const CustomNameLabelBuilder = (custom: NodeCustomName) => {
  return (
    <div className="flex items-center gap-2 flex-wrap">
      <span className="text-muted-foreground">Assign to</span>
      <span className="text-[11px] font-bold uppercase tracking-wider px-2 py-0.5 rounded border bg-muted/30">
        {custom.receiverName}
      </span>
      <span className="text-muted-foreground">on</span>
      <span className="text-[11px] font-bold uppercase tracking-wider px-2 py-0.5 rounded border bg-muted/30">
        {custom.applicationName}
      </span>
    </div>
  )
}

const NodeIcon = (type: Node["type"], icon?: string): LucideIcon | undefined => {
  switch (type) {
    case "perm-folder": return Icons.FileCheck;
    case "role-folder": return Icons.HardHat;
    case "object": return Icons.Box;
    case "scope": return Icons.Globe;
    case "folder": return Icons.Folder;
    case "inherited": return Icons.ArrowDownLeft;
    default: break;
  }
  
  if (icon && icon in Icons) {
    return (Icons as unknown as Record<string, LucideIcon>)[icon];
  }
  
  return undefined;
}

interface NodePropsI {
  node: Node
  level: number
  isLast: boolean
  onNodeClick?: (node: Node) => void
  renderExtra?: (node: Node) => React.ReactNode
  isExpandedByDefault?: boolean
  childrenExpandedByDefault?: boolean
}

// Each tree item height
const ITEM_HEIGHT = 44

function TreeNode({ 
  node, 
  level, 
  isLast, 
  onNodeClick, 
  renderExtra,
  isExpandedByDefault = true,
  childrenExpandedByDefault = true
}: NodePropsI) {
  const [isExpanded, setIsExpanded] = useState(isExpandedByDefault)
  
  if (!node) return null;

  const hasChildren = node.children && node.children.length > 0
  const isFolder = node.id.startsWith('folder-') || node.type === 'folder' || node.type === 'perm-folder' || node.type === 'role-folder';
  
  const Icon = NodeIcon(node.type, node.icon)

  return (
    <div className="relative">
      {level > 0 && (
        <div 
          className="absolute w-px bg-border/60"
          style={{ 
            top: 0, 
            bottom: isLast ? `calc(100% - ${ITEM_HEIGHT / 2}px)` : 0,
            left: `${(level - 1) * 24 + 12}px`
          }}
        />
      )}

      <div className="flex items-center relative group" style={{ height: ITEM_HEIGHT }}>
        {level > 0 && (
          <div 
            className="absolute h-px bg-border/60"
            style={{ 
              width: '20px',
              left: `${(level - 1) * 24 + 12}px`,
              top: '50%'
            }}
          />
        )}

        <div
          className="relative z-10 flex items-center gap-2 shrink-0 min-w-0 flex-1"
          style={{ marginLeft: `${level * 24}px` }}
        >
          <button
            type='button'
            onClick={() => hasChildren && setIsExpanded(!isExpanded)}
            className={cn(
              "w-6 h-6 rounded-full flex items-center justify-center text-xs",
              "font-bold transition-all duration-200 shadow-sm shrink-0",
              !node.color && !hasChildren && "bg-accent/90! hover:bg-accent text-accent-foreground",
              hasChildren 
                ? 'bg-primary/90 text-primary-foreground hover:bg-primary cursor-pointer' 
                : 'bg-muted text-white'
            )}
            style={{ backgroundColor: node.color }}
          >
            {Icon ? <Icon size={14} className={node.color ? "text-white" : undefined}/> : (
              isExpanded ? <FolderOpen size={14} /> : <Folder size={14}/>
            )}
          </button>

          <div className="flex items-center justify-between flex-1 min-w-0 pr-2">
            <button
              type='button'
              onClick={() => {
                if (!isFolder && onNodeClick) onNodeClick(node);
              }}
              disabled={isFolder || !onNodeClick}
              className={cn(
                "flex items-center gap-2 px-2 py-1.5 rounded-lg transition-all min-w-0 text-left",
                onNodeClick && !isFolder && "cursor-pointer hover:bg-muted group-hover:bg-muted border border-transparent hover:border-border",
                !onNodeClick && "cursor-default"
              )}
            >
              {typeof node.name === "string" ?
                <span className={cn(
                  "text-xs font-medium truncate",
                  node.type === 'inherited' && "text-accent italic font-bold",
                  node.type === 'scope' && level === 0 && "text-primary font-bold uppercase tracking-wider",
                  (node.type === 'scope' || isFolder) && level > 0 && "text-muted-foreground uppercase tracking-tight text-[10px]"
                )}>
                  {node.name}
                </span>
              : <CustomNameLabelBuilder {...node.name}/>
              }
            </button>

            {renderExtra && (
              <div className="shrink-0">
                {renderExtra(node)}
              </div>
            )}
          </div>
        </div>
      </div>

      {hasChildren && isExpanded && (
        <div className="relative">
          {node.children?.map((child, index) => (
            <TreeNode
              key={child.id}
              node={child}
              level={level + 1}
              isLast={index === (node.children?.length || 0) - 1}
              onNodeClick={onNodeClick}
              renderExtra={renderExtra}
              isExpandedByDefault={childrenExpandedByDefault}
              childrenExpandedByDefault={childrenExpandedByDefault}
            />
          ))}
        </div>
      )}
    </div>
  )
}

interface TreePropsI {
  node: Node;
  goBack?: () => void;
  onSubmit?: () => void;
  onNodeClick?: (node: Node) => void;
  renderExtra?: (node: Node) => React.ReactNode;
  submitLabel?: string;
  defaultExpanded?: boolean;
  rootDefaultExpanded?: boolean;
  showFooter?: boolean;
}

export default function UserPermTree({ 
  node, 
  goBack, 
  onSubmit, 
  onNodeClick, 
  renderExtra,
  submitLabel = "Grant Access", 
  defaultExpanded = true,
  rootDefaultExpanded = true,
  showFooter = true
}: TreePropsI) {
  return (
    <div className="flex flex-col gap-1">
      <div className="overflow-x-auto">
        <TreeNode
          node={node}
          level={0}
          isLast={true}
          onNodeClick={onNodeClick}
          renderExtra={renderExtra}
          isExpandedByDefault={rootDefaultExpanded}
          childrenExpandedByDefault={defaultExpanded}
        />
      </div>
      {showFooter && (
        <>
          <hr className="w-full border-muted mt-2"/>
          <div className="w-full flex justify-between items-center mt-3">
            {goBack && <ShadowButton value="Back" variant="ghost" onClick={goBack}/>}
            {onSubmit && (
              <ShadowButton
                value={submitLabel} 
                variant="solid"
                onClick={onSubmit}
              />
            )}
          </div>
        </>
      )}
    </div>
  )
}