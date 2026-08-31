import { useQuery } from "@tanstack/react-query";
import { createLazyFileRoute, useRouter } from "@tanstack/react-router";
import { programsQueryOptions } from "@/features/programs/api";
import { OccurrenceDrawPage } from "@/features/programs/ui/OccurrenceDrawPage";

export const Route = createLazyFileRoute(
  "/admin/events/$eventId_/editions/$editionId/programs/$programId/occurrences/$occurrenceId/draw",
)({ component: DrawRoute });

function DrawRoute() {
  const { editionId, programId, occurrenceId } = Route.useParams();
  const router = useRouter();
  const { data: programs = [] } = useQuery(programsQueryOptions(editionId));
  const program = programs.find((item) => item.id === programId);

  return (
    <OccurrenceDrawPage
      occurrenceId={occurrenceId}
      programName={program?.name ?? "Atividade"}
      onBack={() => router.history.back()}
    />
  );
}
