import { FileText, Pencil, Plus, Trash2 } from 'lucide-react'
import { useState } from 'react'
import { RelationBadge } from './relation-badge'
import type { SpiceDBRelationshipI } from '@trieoh/node-perm-sdk'

interface RelationshipsTableProps {
  rows: SpiceDBRelationshipI[]
  onEdit: (rel: SpiceDBRelationshipI) => void
  onDelete: (rel: SpiceDBRelationshipI) => void
  onNew: () => void
}

export function RelationshipsTable({
  rows,
  onEdit,
  onDelete,
  onNew,
}: RelationshipsTableProps) {
  const [filter, setFilter] = useState('')

  const filtered = rows.filter((r) =>
    [
      r.resource.objectType,
      r.resource.objectId,
      r.relation,
      r.subject.object.objectType,
      r.subject.object.objectId,
    ]
      .join(' ')
      .toLowerCase()
      .includes(filter.toLowerCase()),
  )

  return (
    <div className="flex flex-col">
      {/* Header */}
      <div className="flex h-14 items-center justify-between border-b px-4 shrink-0 gap-4">
        <div className="flex items-center gap-2 min-w-0">
          <FileText size={14} className="text-muted-foreground shrink-0" />
          <span className="text-sm font-medium truncate">Relationships</span>
          <span className="rounded-full bg-muted px-2 py-0.5 text-[10px] text-muted-foreground hidden sm:inline-flex shrink-0">
            {filtered.length}
          </span>
        </div>

        <div className="flex items-center gap-2 flex-1 justify-end min-w-0">
          <input
            type="text"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            placeholder="Filter..."
            className="h-8 min-w-0 max-w-30 sm:max-w-none flex-1 rounded-md border border-input bg-background px-3 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring sm:w-40 sm:flex-none"
          />
          {/* Mobile-only new button */}
          <button
            onClick={onNew}
            className="inline-flex h-8 shrink-0 items-center gap-1.5 rounded-md bg-primary px-3 text-xs font-medium text-primary-foreground transition-colors hover:bg-primary/90 md:hidden"
          >
            <Plus size={12} />
            New
          </button>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        {filtered.length === 0 ? (
          <div className="flex flex-col items-center gap-3 py-16 text-center">
            <FileText size={28} className="text-muted-foreground/30" />
            <p className="text-sm text-muted-foreground">No relationships found</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b">
                <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground sm:px-4">
                  Relationship
                </th>
                <th className="hidden px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground sm:table-cell">
                  Relation
                </th>
                <th className="w-16 px-3 py-3 sm:px-4" />
              </tr>
            </thead>
            <tbody>
              {filtered.map((rel, idx) => (
                <tr
                  key={`${rel.resource.objectId}-${rel.relation}-${rel.subject.object.objectId}-${idx}`}
                  className="border-b last:border-0 hover:bg-muted/40"
                >
                  <td className="px-3 py-2.5 sm:px-4 sm:py-3">
                    <div className="flex flex-col gap-0.5">
                      <div className="flex items-baseline gap-0.5">
                        <span className="font-mono text-xs text-muted-foreground">
                          {rel.resource.objectType}:
                        </span>
                        <span className="font-medium">{rel.resource.objectId}</span>
                      </div>
                      <div className="flex items-baseline gap-0.5">
                        <span className="font-mono text-xs text-muted-foreground">
                          {rel.subject.object.objectType}:
                        </span>
                        <span className="text-xs text-muted-foreground">
                          {rel.subject.object.objectId}
                        </span>
                      </div>
                      {/* Badge inline only on mobile */}
                      <div className="mt-1 sm:hidden">
                        <RelationBadge value={rel.relation} />
                      </div>
                    </div>
                  </td>

                  {/* Relation col — sm+ */}
                  <td className="hidden px-4 py-3 sm:table-cell">
                    <RelationBadge value={rel.relation} />
                  </td>

                  <td className="px-3 py-2.5 sm:px-4 sm:py-3">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        onClick={() => onEdit(rel)}
                        title="Edit"
                        className="inline-flex h-7 w-7 items-center justify-center rounded-md border border-transparent text-muted-foreground transition-colors hover:border-border hover:bg-muted hover:text-foreground"
                      >
                        <Pencil size={13} />
                      </button>
                      <button
                        onClick={() => onDelete(rel)}
                        title="Delete"
                        className="inline-flex h-7 w-7 items-center justify-center rounded-md border border-transparent text-muted-foreground transition-colors hover:border-destructive/30 hover:bg-destructive/10 hover:text-destructive"
                      >
                        <Trash2 size={13} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
