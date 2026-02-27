import { useState } from 'react'
import type { Node, NodeCustomName } from '../model/types'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'

const CustomNameLabelBuilder = (custom: NodeCustomName) => {
  return (
    <>
      <span>Assign to</span>
      <span className="text-sm font-medium px-3 py-1.5 rounded-lg border shadow-sm">
        {custom.receiverName}
      </span>
      <span>on</span>
      <span className="text-sm font-medium px-3 py-1.5 rounded-lg border shadow-sm">
        {custom.applicationName}
      </span>
    </>
  )
}

interface NodePropsI {
  node: Node
  level: number
  isLast: boolean
}

// Each tree item height (py-2 = 8px + 8px, h-6 = 24px, gap-2 = 8px ≈ 48px total)
const ITEM_HEIGHT = 48

function TreeNode({ node, level, isLast }: NodePropsI) {
  const [isExpanded, setIsExpanded] = useState(true)
  const hasChildren = node.children && node.children.length > 0

  return (
    <div className="relative">
      {level > 0 && (
        <div 
          className="absolute w-px bg-slate-300"
          style={{ 
            top: 0, 
            bottom: isLast ? `calc(100% - ${ITEM_HEIGHT / 2}px)` : 0,
            left: `${(level - 1) * 24 + 12}px`
          }}
        />
      )}
      
      <div className="flex items-center relative" style={{ height: ITEM_HEIGHT }}>
        {level > 0 && (
          <div 
            className="absolute h-px bg-slate-300"
            style={{ 
              width: '20px',
              left: `${(level - 1) * 24 + 12}px`,
              top: '50%'
            }}
          />
        )}
        
        <div 
          className="relative z-10 flex items-center gap-2"
          style={{ marginLeft: `${level * 24}px` }}
        >
          <button
            onClick={() => hasChildren && setIsExpanded(!isExpanded)}
            className={`
              w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold
              transition-all duration-200
              ${hasChildren 
                ? 'bg-indigo-500 text-white hover:bg-indigo-600 cursor-pointer shadow-md' 
                : 'bg-slate-200 text-slate-500'
              }
            `}
          >
            {hasChildren ? (isExpanded ? '−' : '+') : '•'}
          </button>
          
          {typeof node.name === "string" ?
            <span className="text-sm font-medium px-3 py-1.5 rounded-lg border shadow-sm">
              {node.name}
            </span>
          : <CustomNameLabelBuilder {...node.name}/>
          }
        </div>
      </div>

      {hasChildren && isExpanded && (
        <div className="relative">
          {node.children!.map((child, index) => (
            <TreeNode
              key={child.id}
              node={child}
              level={level + 1}
              isLast={index === node.children!.length - 1}
            />
          ))}
        </div>
      )}
    </div>
  )
}

interface TreePropsI {
  node: Node;
  goBack: () => void;
  onSubmit: () => void;
}

export default function UserPermTree({ node, goBack, onSubmit }: TreePropsI) {
  return (
    <>
      <div className="min-w-100 overflow-x-auto">
        <TreeNode
          node={node}
          level={0}
          isLast={true}
        />
      </div>
      <hr className="w-full"/>
      <div className="w-full flex justify-between items-center mt-3">
        <ShadowButton value="Add More" variant="ghost" onClick={goBack}/>
        <ShadowButton
          value="Grant Access" 
          variant="solid"
          onClick={onSubmit}
        />
      </div>
    </>
  )
}