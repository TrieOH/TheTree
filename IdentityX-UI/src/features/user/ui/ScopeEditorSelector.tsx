import type { Scope } from "@/features/scope/model/types";
import { buildScopeHierarchyToNodeTree } from "../lib/node-tree-utils";
import UserPermTree from "./UserPermTree";

interface PropsI {
  setCurrentScopeID: (value: string | null) => void;
  setCurrentType: (value: null) => void;
  allScopes: Scope[];
  currentType: string;
}

export default function ScopeEditorSelector({ setCurrentScopeID, setCurrentType, allScopes, currentType }: PropsI) {
  const scopeTree = buildScopeHierarchyToNodeTree(allScopes);

  return (
    <div className="flex flex-col items-center gap-3 text-foreground">
      <div className="text-center w-full">
        <span className="text-primary font-bold">
          SELECT SCOPE FOR {currentType.toUpperCase()}
        </span>
        <p className="text-xs text-muted-foreground">What type of access do you want to grant?</p>
      </div>
      <div className="w-full">
        <UserPermTree
          node={scopeTree}
          goBack={() => setCurrentType(null)}
          onNodeClick={(node) => setCurrentScopeID(node.id === 'null' ? null : node.id)}
          defaultExpanded={false}
        />
      </div>
    </div>
  )
}