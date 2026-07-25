import { createFileRoute, redirect } from "@tanstack/react-router";
import { env } from "#/env";
import { CreateTestModeIntentForm } from "#/features/payment-intents/ui/create-test-mode-intent-form";

export const Route = createFileRoute("/admin/test-mode-intents")({
  beforeLoad: () => {
    if (env.VITE_INTENT_TEST_MODE !== "true") throw redirect({ to: "/admin" });
  },
  component: TestModeIntentsPage,
});

function TestModeIntentsPage() {
  return (
    <div className="space-y-6 p-6">
      <div>
        <p className="text-xs font-medium uppercase tracking-widest text-amber-600">
          Test mode
        </p>
        <h1 className="mt-1 text-2xl font-semibold tracking-tight">
          Create test intent
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Insert an intent directly without running the provider checkout.
        </p>
      </div>
      <CreateTestModeIntentForm />
    </div>
  );
}
