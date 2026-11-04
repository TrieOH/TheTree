import problems from "virtual:problems";
import { findService } from "@/features/catalog/services";
import type { ProblemDoc, ProblemMember, ProblemSet } from "../../../mdx/problems-plugin";

export type { ProblemDoc, ProblemMember, ProblemSet };

export const PROBLEM_SETS: readonly ProblemSet[] = ["shared", "identityx", "payssage", "univents", "informd"];

export function isProblemSet(value: string): value is ProblemSet {
  return (PROBLEM_SETS as readonly string[]).includes(value);
}

export function problemsOf(set: ProblemSet): ProblemDoc[] {
  return problems[set];
}

export function findProblem(set: string, slug: string): ProblemDoc | undefined {
  return isProblemSet(set) ? problems[set].find((problem) => problem.slug === slug) : undefined;
}

export function setName(set: ProblemSet): string {
  return set === "shared" ? "Shared" : (findService(set)?.name ?? set);
}

const STATUS_TEXT: Record<number, string> = {
  400: "Bad Request",
  401: "Unauthorized",
  402: "Payment Required",
  403: "Forbidden",
  404: "Not Found",
  405: "Method Not Allowed",
  409: "Conflict",
  410: "Gone",
  412: "Precondition Failed",
  413: "Content Too Large",
  415: "Unsupported Media Type",
  422: "Unprocessable Content",
  429: "Too Many Requests",
  500: "Internal Server Error",
  501: "Not Implemented",
  502: "Bad Gateway",
  503: "Service Unavailable",
  504: "Gateway Timeout",
};

export function statusText(status: number): string {
  return STATUS_TEXT[status] ?? "";
}

const words = (text: string) => text.toLowerCase().replace(/[^\p{L}\p{N}]+/gu, " ").trim();

/** The description, unless it only restates the title ("Invalid access token."). */
export function extraDescription(problem: ProblemDoc): string | null {
  return problem.description && words(problem.description) !== words(problem.title) ? problem.description : null;
}

function exampleValue(member: ProblemMember): unknown {
  if (member.example !== undefined) return member.example;
  switch (member.type) {
    case "integer":
    case "number":
      return 0;
    case "boolean":
      return false;
    case "array":
      return [];
    case "object":
      return {};
    default:
      return "…";
  }
}

/** What a client receives: the standard members, then the problem's own. */
export function exampleBody(problem: ProblemDoc, set: ProblemSet): string {
  const api = set === "shared" ? "{API_URL}" : (findService(set)?.apiUrl ?? "{API_URL}");
  const body = {
    type: `${api}/problems/${problem.slug}`,
    status: problem.status,
    title: problem.title,
    detail: problem.description || problem.title,
    instance: "/the/request/path",
    ...Object.fromEntries(problem.members.map((member) => [member.name, exampleValue(member)])),
  };
  return JSON.stringify(body, null, 2);
}

/** Problems grouped by status code, lowest first. */
export function byStatus(list: readonly ProblemDoc[]): { status: number; problems: ProblemDoc[] }[] {
  const groups = new Map<number, ProblemDoc[]>();
  for (const problem of list) groups.set(problem.status, [...(groups.get(problem.status) ?? []), problem]);
  return [...groups.entries()].sort(([a], [b]) => a - b).map(([status, problems]) => ({ status, problems }));
}
