import ScopeDialog from "./ScopeDialog";
import { useQuery } from "@tanstack/react-query";
import { scopesQueryOptions } from "../api";
import ScopeTreeView from "./ScopeTreeView";

interface PropsI {
  project_id: string;
}

export default function ScopeTable({ project_id }: PropsI) {
  const { data = [] } = useQuery(scopesQueryOptions(project_id))

  return (
    <>
      <ScopeTreeView scopes={data} />
      <ScopeDialog project_id={project_id}/>
    </>
  )
}
