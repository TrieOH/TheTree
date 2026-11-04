import { createFileRoute, notFound } from "@tanstack/solid-router";
import { isProblemSet, problemsOf, setName } from "@/features/problems/model";
import { ProblemList } from "@/features/problems/ui/ProblemList";
import { PageHeader } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/problems/$set/")({
  loader: ({ params }) => {
    if (!isProblemSet(params.set)) throw notFound();
    return params.set;
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData ? setName(loaderData) : ""} problems · TrieOH dev` }] }),
  component: ProblemsPage,
});

function ProblemsPage() {
  const set = Route.useLoaderData();
  return (
    <>
      <PageHeader
        title={set() === "shared" ? "Shared problems" : `${setName(set())} problems`}
        description={
          set() === "shared"
            ? `The ${problemsOf("shared").length} RFC 9457 problems every TrieOH API can answer.`
            : `Every RFC 9457 problem ${setName(set())} can answer beyond the shared ones. Clients branch on the slug at the end of type.`
        }
      />
      <ProblemList set={set()} />
    </>
  );
}
