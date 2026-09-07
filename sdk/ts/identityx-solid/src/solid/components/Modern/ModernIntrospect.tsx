import { createSignal } from "solid-js";
import { useAuth } from "../../AuthProvider";
import { Button, FormInput } from "./Shared";
export interface ModernIntrospectProps {
  data?: unknown;
}
export function ModernIntrospect(props: ModernIntrospectProps) {
  const { auth } = useAuth();
  const [apiKey, setApiKey] = createSignal("");
  const [data, setData] = createSignal(props.data);
  const load = async () => {
    const result = await auth.introspect(apiKey() || undefined);
    if (result.success) setData(result.data);
  };
  return (
    <div class="w-full flex flex-col gap-4">
      <FormInput
        label="API key (opcional)"
        value={apiKey()}
        onInput={(e) => setApiKey(e.currentTarget.value)}
      />
      <Button onClick={() => void load()}>Consultar</Button>
      {data() && (
        <pre class="max-h-96 overflow-auto rounded bg-muted p-4 text-xs">
          {JSON.stringify(data(), null, 2)}
        </pre>
      )}
    </div>
  );
}
