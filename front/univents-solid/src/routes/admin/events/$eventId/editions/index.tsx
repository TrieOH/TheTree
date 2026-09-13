import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import CalendarDaysIcon from "~icons/lucide/calendar-days";
import { EmptyState } from "@trieoh/ui-solid";

const CalendarDays = CalendarDaysIcon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute("/admin/events/$eventId/editions/")({
  head: () => ({
    meta: [{ title: "Edições - Admin Univents" }],
  }),
  component: AdminEventEditionsRoute,
});

function AdminEventEditionsRoute(): JSX.Element {
  return (
    <div class="space-y-6">
      <div class="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-foreground sm:text-3xl">
            Edições
          </h1>
          <p class="text-sm text-muted-foreground">
            Gerencie todas as edições, datas e locais vinculados a este evento.
          </p>
        </div>
      </div>

      <EmptyState
        icon={<CalendarDays class="h-6 w-6 text-foreground/70" />}
        eyebrow="Edições do Evento"
        title="Nenhuma edição cadastrada"
        description="As edições deste evento serão exibidas aqui para você acompanhar datas, locais e ingressos."
      />
    </div>
  );
}
