import type { ServiceId } from "@/features/catalog/services";
import { ProblemList, type ProblemListProps } from "@/features/problems/ui/ProblemList";

/**
 * `<Problems service="payssage" />` in a doc; see ProblemListProps for the
 * options (`shared`, `only`, `group`, `open`, `description`).
 */
export function Problems(props: Omit<ProblemListProps, "set"> & { service: ServiceId | "shared" }) {
  return (
    <ProblemList
      set={props.service}
      shared={props.shared}
      only={props.only}
      group={props.group}
      open={props.open}
      description={props.description}
    />
  );
}
