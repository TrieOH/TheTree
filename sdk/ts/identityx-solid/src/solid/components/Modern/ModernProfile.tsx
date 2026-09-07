import { createSignal, onSettled, Show } from "solid-js";
import { useAuth } from "../../AuthProvider";
export interface ModernProfileProps {
  actorId?: string;
  projectId?: string;
  schema?: unknown;
  onSuccess?: (message?: string) => Promise<void>;
}
export function ModernProfile(props: ModernProfileProps) {
  const { auth } = useAuth();
  const [profile, setProfile] = createSignal<{
    success: boolean;
    data?: unknown;
  } | null>(null);
  onSettled(() => {
    if (props.actorId)
      void auth
        .getProjectProfile(props.actorId, props.projectId)
        .then(setProfile);
  });
  return (
    <Show when={profile()?.success} fallback={<p>Perfil indisponível.</p>}>
      <pre class="max-h-96 overflow-auto rounded bg-muted p-4 text-xs">
        {JSON.stringify(profile()?.data, null, 2)}
      </pre>
    </Show>
  );
}
