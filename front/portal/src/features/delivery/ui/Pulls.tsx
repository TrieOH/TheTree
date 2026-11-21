import { Badge, EmptyState, PageSection, Table } from "@trieoh/ui-solid";
import { pullUrl, type PullRequest } from "../model";
import { formatDuration, timeAgo } from "../stats";

const external = { target: "_blank", rel: "noopener" } as const;

type Tone = "success" | "destructive" | "warning" | "muted";

function reviewPill(pr: PullRequest): { tone: Tone; label: string } {
  if (pr.draft) return { tone: "muted", label: "Draft" };
  if (pr.reviewState === "approved") return { tone: "success", label: "Approved" };
  if (pr.reviewState === "changes_requested") return { tone: "destructive", label: "Changes requested" };
  return { tone: "warning", label: "Needs review" };
}

const title = (pr: PullRequest) => (
  <a class="hover:underline" href={pullUrl(pr.number)} {...external}>
    <span class="mr-2 font-mono text-[12px] text-muted-foreground">#{pr.number}</span>
    {pr.title}
  </a>
);

export function Pulls(props: { pulls: readonly PullRequest[] }) {
  const open = () => props.pulls.filter((pr) => pr.state === "open");
  const merged = () => props.pulls.filter((pr) => pr.state === "merged");
  return (
    <>
      <PageSection
        title="Open"
        aside={
          <a href="https://github.com/TrieOH/TheTree/pulls" {...external}>
            All on GitHub
          </a>
        }
      >
        <Table
          caption="Open pull requests"
          class="text-[13.5px]"
          items={open()}
          emptyState={<EmptyState variant="inline" title="No open pull requests." />}
          columns={[
            { id: "title", header: "Pull request", cell: title },
            { id: "author", header: "Author", class: "whitespace-nowrap text-muted-foreground", cell: (pr) => pr.author },
            {
              id: "open",
              header: "Open for",
              align: "right",
              class: "whitespace-nowrap tabular-nums text-muted-foreground",
              cell: (pr) => formatDuration(Date.now() - Date.parse(pr.openedAt)),
            },
            {
              id: "review",
              header: "Review",
              align: "right",
              cell: (pr) => (
                <Badge shape="pill" variant={reviewPill(pr).tone}>
                  {reviewPill(pr).label}
                </Badge>
              ),
            },
          ]}
        />
      </PageSection>
      <PageSection title="Merged" aside="last 30 days">
        <Table
          caption="Merged pull requests"
          class="text-[13.5px]"
          density="compact"
          items={merged()}
          emptyState={<EmptyState variant="inline" title="Nothing merged in the last 30 days." />}
          columns={[
            { id: "title", header: "Pull request", cell: title },
            { id: "author", header: "Author", class: "whitespace-nowrap text-muted-foreground", cell: (pr) => pr.author },
            {
              id: "lead",
              header: "Open to merge",
              align: "right",
              class: "whitespace-nowrap tabular-nums",
              cell: (pr) => formatDuration(Date.parse(pr.mergedAt!) - Date.parse(pr.openedAt)),
            },
            {
              id: "merged",
              header: "Merged",
              align: "right",
              class: "whitespace-nowrap text-muted-foreground",
              cell: (pr) => timeAgo(pr.mergedAt!),
            },
          ]}
        />
      </PageSection>
    </>
  );
}
