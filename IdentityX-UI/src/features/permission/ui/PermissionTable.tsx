import { permissionsQueryOptions } from "../api";
import { useQuery } from "@tanstack/react-query";
import { formatDate } from "@/shared/lib/date-utils";
import { Edit, Plus, Shield, Trash2, Filter } from "lucide-react";
import { permissionActions } from "../store";
import PermissionDialog from "./PermissionDialog";
import { Badge } from "@/shared/ui/shadcn/badge";
import TruncatedId from "@/shared/ui/TruncatedId";
import { MetadataVisualizer, type VisualMetadata } from "@/shared/ui/MetadataVisualizer";
import type { Permission } from "../model/types";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/shared/ui/shadcn/accordion";
import { useState, useMemo } from "react";
import { SearchInput } from "@/shared/ui/form/SearchInput";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/shadcn/select";

interface PropsI {
  project_id: string;
}

interface FlattenedPermission extends Permission {
  status: string;
}

interface GroupedPermission {
  object: string;
  permissions: FlattenedPermission[];
  meta?: VisualMetadata;
}

export default function PermissionTable({ project_id }: PropsI) {
  const { data = [] } = useQuery(permissionsQueryOptions(project_id))
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<string>("all");

  const groupedData = useMemo(() => {
    return data.reduce((acc, p) => {
      const status = (p.meta?.status as string) || "active";
      const flattened: FlattenedPermission = {
        ...p,
        status,
      };
      
      if (!acc[p.object]) {
        acc[p.object] = {
          object: p.object,
          permissions: [],
          meta: p.meta,
        };
      }
      acc[p.object].permissions.push(flattened);
      return acc;
    }, {} as Record<string, GroupedPermission>);
  }, [data]);

  const filteredData = useMemo(() => {
    const searchLower = search.toLowerCase();
    
    return Object.values(groupedData)
      .map(group => {
        // Filter permissions within the group by status and search
        const filteredPermissions = group.permissions.filter(p => {
          const matchesStatus = statusFilter === "all" || p.status === statusFilter;
          const matchesSearch = p.action.toLowerCase().includes(searchLower) || 
                               group.object.toLowerCase().includes(searchLower);
          return matchesStatus && matchesSearch;
        });

        if (filteredPermissions.length === 0) return null;

        return {
          ...group,
          permissions: filteredPermissions
        };
      })
      .filter((group): group is GroupedPermission => group !== null);
  }, [groupedData, search, statusFilter]);

  const getBadgeStyle = (meta?: VisualMetadata) => {
    if (!meta?.color) return {};
    
    const isGradient = meta.color.includes("linear-gradient");
    return isGradient 
      ? { backgroundImage: meta.color, color: 'white', border: 'none' } 
      : { backgroundColor: meta.color, color: 'white', border: 'none' };
  };

  return (
    <div className="space-y-6 pb-2">
      <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
        <div className="flex flex-col sm:flex-row flex-1 items-center gap-2 w-full max-w-xl">
          <div className="w-full sm:flex-1">
            <SearchInput
              placeholder="Search by object or action..."
              value={search}
              onChange={setSearch}
            />
          </div>
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-full sm:w-35 h-full!">
              <div className="flex items-center gap-2">
                <Filter size={14} className="text-muted-foreground" />
                <SelectValue placeholder="Status" />
              </div>
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Status</SelectItem>
              <SelectItem value="active">Active</SelectItem>
              <SelectItem value="restricted">Restricted</SelectItem>
              <SelectItem value="beta">Beta</SelectItem>
              <SelectItem value="deprecated">Deprecated</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <ShadowButton
          onClick={() => permissionActions.openCreate()}
          leftIcon={<Plus size={16} />}
          variant="solid"
          className="w-full md:w-auto"
          value="Create Permission"
        />
      </div>

      {filteredData.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-12 text-center border-2 border-dashed rounded-lg bg-muted/20 border-border">
          <Shield className="h-12 w-12 text-muted-foreground mb-4 opacity-20" />
          <h3 className="text-lg font-medium">No permissions found</h3>
          <p className="text-sm text-muted-foreground">
            {search || statusFilter !== "all" ? "Try adjusting your filters" : "Start by creating your first permission"}
          </p>
        </div>
      ) : (
        <Accordion type="multiple" className="w-full space-y-4">
          {filteredData.map((group) => (
            <AccordionItem 
              key={group.object} 
              value={group.object}
              className="border border-border rounded-md bg-card"
            >
              <AccordionTrigger className="hover:no-underline p-4">
                <div className="flex items-center justify-between w-full overflow-hidden pr-2">
                  <div className="flex items-center gap-3 overflow-hidden">
                    <span className="font-semibold text-sm text-foreground tracking-tight shrink-0">
                      {group.object}
                    </span>
                    <div className="flex items-center gap-1.5 overflow-hidden">
                      {/* Show 1 tag on mobile, up to 4 on desktop */}
                      <div className="flex items-center gap-1.5 flex-nowrap">
                        <div className="flex sm:hidden items-center gap-1.5">
                          {group.permissions.slice(0, 1).map((p) => (
                            <Badge 
                              key={p.id} 
                              variant="secondary" 
                              style={getBadgeStyle(p.meta)}
                              className="text-[10px] py-0 px-2 h-5 font-mono shrink-0"
                            >
                              {p.action}
                            </Badge>
                          ))}
                          {group.permissions.length > 1 && (
                            <Badge variant="outline" className="text-[10px] py-0 px-2 h-5 shrink-0">
                              +{group.permissions.length - 1}
                            </Badge>
                          )}
                        </div>
                        <div className="hidden sm:flex items-center gap-1.5">
                          {group.permissions.slice(0, 4).map((p) => (
                            <Badge 
                              key={p.id} 
                              variant="secondary" 
                              style={getBadgeStyle(p.meta)}
                              className="text-[10px] py-0 px-2 h-5 font-mono shrink-0"
                            >
                              {p.action}
                            </Badge>
                          ))}
                          {group.permissions.length > 4 && (
                            <Badge variant="outline" className="text-[10px] py-0 px-2 h-5 shrink-0">
                              +{group.permissions.length - 4}
                            </Badge>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </AccordionTrigger>
              <AccordionContent className="pb-0">
                <div className="overflow-hidden rounded-md border border-border bg-background">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border bg-muted/40 text-muted-foreground">
                        <th className="px-4 py-2 text-left text-[10px] font-bold uppercase tracking-wider">Action</th>
                        <th className="px-4 py-2 text-left text-[10px] font-bold uppercase tracking-wider hidden sm:table-cell">ID</th>
                        <th className="px-4 py-2 text-left text-[10px] font-bold uppercase tracking-wider hidden md:table-cell">Created At</th>
                        <th className="px-4 py-2 text-right"></th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border">
                      {group.permissions.map((p) => (
                        <tr key={p.id} className="hover:bg-muted/50 transition-colors">
                          <td className="px-4 py-2">
                            <div className="flex items-center gap-2">
                              {p.meta ? (
                                <MetadataVisualizer name={p.action} meta={p.meta} />
                              ) : (
                                <Badge variant={p.action === "*" ? "secondary" : "outline"} className="font-mono text-[10px]">
                                  {p.action}
                                </Badge>
                              )}
                            </div>
                          </td>
                          <td className="px-4 py-2 text-muted-foreground hidden sm:table-cell">
                            <TruncatedId id={p.id} />
                          </td>
                          <td className="px-4 py-2 text-muted-foreground text-xs font-mono hidden md:table-cell">
                            {formatDate(p.created_at)}
                          </td>
                          <td className="px-4 py-2 text-right">
                            <div className="flex items-center justify-end gap-1">
                              <ShadowButton
                                variant="ghost-primary"
                                className="h-8 w-8 p-0 flex items-center justify-center rounded-md hover:bg-primary/10"
                                onClick={() => permissionActions.openEdit(p)}
                                leftIcon={<Edit size={14} />}
                              />
                              <ShadowButton
                                variant="ghost-primary"
                                className="h-8 w-8 p-0 flex items-center justify-center rounded-md text-destructive hover:bg-destructive/10 hover:text-destructive"
                                onClick={() => permissionActions.openDelete(p)}
                                leftIcon={<Trash2 size={14} />}
                              />
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </AccordionContent>
            </AccordionItem>
          ))}
        </Accordion>
      )}
      
      <PermissionDialog project_id={project_id}/>
    </div>
  )
}
