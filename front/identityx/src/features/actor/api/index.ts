import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  createActor,
  createOrganizationProjectActor,
  listActors,
  listOrganizationProjectActors,
} from "@trieoh/identityx-api";
import type { ActorCreateI, ActorI } from "../model";

/**
 * Creates a new actor in the specified project.
 * @param project_id - The ID of the project to create the actor in.
 * @param actor - The actor data to create.
 * @param organization_id - The ID of the organization (optional).
 * @returns A promise that resolves to the created ActorI object.
 */
export const createActorFn = createClientOnlyFn(
  async (project_id: string, actor: ActorCreateI, organization_id?: string) => {
    if (organization_id)
      return createOrganizationProjectActor(
        organization_id,
        project_id,
        actor,
      ).then(orvalData<ActorI>);
    return createActor(project_id, actor).then(orvalData<ActorI>);
  },
);

/**
 * Fetches all actors from the server.
 * @param project_id - The project ID to fetch actors for.
 * @param organization_id - The organization ID to filter actors by (optional).
 * @returns A promise that resolves to an array of ActorI objects.
 */
export const getActorsFn = createClientOnlyFn(
  async (project_id: string, organization_id?: string) => {
    if (organization_id)
      return listOrganizationProjectActors(organization_id, project_id).then(
        orvalData<ActorI[]>,
      );
    return listActors(project_id).then(orvalData<ActorI[]>);
  },
);

/**
 * Query options for fetching actors, compatible with React Query's useQuery hook.
 * @param project_id - The project ID to fetch actors for.
 * @param organization_id - The organization ID to filter actors by (optional).
 * @returns An object containing the query key and query function for fetching actors.
 */
export const allActorsQueryOptions = (
  project_id: string,
  organization_id?: string,
) => {
  return queryOptions({
    queryKey: [
      "organizations",
      organization_id,
      "projects",
      project_id,
      "actors",
    ],
    queryFn: () => getActorsFn(project_id, organization_id),
  });
};
