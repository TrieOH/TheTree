import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { Show } from "solid-js";
import { CopyButton } from "@/shared/ui/CopyButton";
import { extraDescription, findProblem, setName, statusText, type ProblemSet } from "@/features/problems/model";
import { ProblemDetails } from "@/features/problems/ui/ProblemDetails";
import { StatusPill } from "@/features/problems/ui/ProblemList";
import { PageHeader } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/problems/$set/$slug")({
  loader: ({ params }) => {
    const problem = findProblem(params.set, params.slug);
    if (!problem) throw notFound();
    return { problem, set: params.set as ProblemSet };
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.problem.slug ?? "Problem"} · TrieOH dev` }] }),
  component: ProblemPage,
});

function ProblemPage() {
  const data = Route.useLoaderData();
  const problem = () => data().problem;

  return (
    <>
      <PageHeader title={problem().title}>
        <div class="flex flex-wrap items-center gap-2 text-[13px] text-muted-foreground">
          <StatusPill status={problem().status} />
          <span>{statusText(problem().status)}</span>
          <span>·</span>
          <code class="font-mono text-[13px] text-foreground">{problem().slug}</code>
          <CopyButton text={() => problem().slug} label="Copy slug" />
          <span>·</span>
          <Link to="/problems/$set" params={{ set: data().set }} class="underline underline-offset-4">
            {data().set === "shared" ? "Shared problems" : `${setName(data().set)} problems`}
          </Link>
        </div>
      </PageHeader>
      <Show when={extraDescription(problem())}>
        <p class="mb-6">{extraDescription(problem())}</p>
      </Show>
      <ProblemDetails problem={problem()} />
    </>
  );
}
