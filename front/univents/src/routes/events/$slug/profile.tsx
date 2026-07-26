import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/events/$slug/profile")({
  component: EventProfilePage,
});

function EventProfilePage() {
  return <main className="min-h-screen bg-background"></main>;
}
