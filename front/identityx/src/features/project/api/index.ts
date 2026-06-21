import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import type { MemberAddToProjectI, ProjectCreateI, ProjectI, ProjectMemberI } from "../model";
import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";

/**
 * Creates a new project on the server.
 * @param projectData - The data for the new project.
 * @param orgId - The organization ID to which the project belongs (optional).
 * @returns A promise that resolves to the API response containing the newly created project.
 */
export const createProjectFn = createClientOnlyFn((projectData: ProjectCreateI, orgId?: string) => {
  if (orgId)
    return authFetcher.post<ProjectI>(`/organizations/${orgId}/projects`, projectData);
  return authFetcher.post<ProjectI>("/projects", projectData);
});

/**
 * Fetches all projects from the server.
 * @param orgId - The organization ID to filter projects by (optional).
 * @returns A promise that resolves to an array of ProjectI objects.
 */
export const getProjectsFn = createClientOnlyFn(async (orgId?: string) => {
  if (orgId)
    return await tanstackQueryFetcher<ProjectI[]>(`/organizations/${orgId}/projects`);
  return await tanstackQueryFetcher<ProjectI[]>("/projects");
});

/**
 * Query options for fetching projects, compatible with React Query's useQuery hook.
 * @param orgId - The organization ID to filter projects by (optional).
 * @returns An object containing the query key and query function for fetching projects.
 */
export const allProjectsQueryOptions = (orgId?: string) => {
  return queryOptions({
    queryKey: ["organizations", orgId, "projects"],
    queryFn: () => getProjectsFn(orgId),
  });
};

// Members

/**
 * Adds a new member to a organization on the server.
 * @param project_id - The ID of the project to add the member to.
 * @param memberData - The data for the new member.
 * @param organization_id - The ID of the organization to add the member to (optional).
 * @returns A promise that resolves to the API response containing the newly created member.
 */
export const addMemberToProjectFn = createClientOnlyFn((project_id: string, memberData: MemberAddToProjectI, organization_id?: string) => {
  if (organization_id)
    return authFetcher.post(`/organizations/${organization_id}/projects/${project_id}/members`, memberData);
  return authFetcher.post(`/projects/${project_id}/members`, memberData);
});

/**
 * Removes a member from a organization on the server.
 * @param project_id - The ID of the project to remove the member from.
 * @param actor_email - The email of the user to remove from the project.
 * @param organization_id - The ID of the organization to remove the member from (optional).
 * @returns A promise that resolves to the API response confirming the removal of the member.
 */
export const removeMemberFromProjectFn = createClientOnlyFn((project_id: string, actor_email: string, organization_id?: string) => {
  if (organization_id)
    return authFetcher.delete(`/organizations/${organization_id}/projects/${project_id}/members`, { actor_email });
  return authFetcher.delete(`/projects/${project_id}/members`, { actor_email });
});

/**
 * Fetches all project members from the server.
 * @param project_id - The ID of the project to fetch members for.
 * @param organization_id - The ID of the organization to fetch members for (optional).
 * @returns A promise that resolves to an array of project members objects.
 */
export const getAllProjectMembersFn = createClientOnlyFn((
  project_id: string,
  organization_id?: string
) => {
  if (organization_id)
    return tanstackQueryFetcher<ProjectMemberI[]>(`/organizations/${organization_id}/projects/${project_id}/members`);
  return tanstackQueryFetcher<ProjectMemberI[]>(`/projects/${project_id}/members`);
});

/**
 * Query options for fetching all Project Members, using TanStack Query.
 * @param project_id - The ID of the project to fetch members for.
 * @param organization_id - The ID of the organization to fetch members for (optional).
 * @returns An object containing the query key and query function for fetching all Project Members.
 */
export const allProjectMembersQueryOptions = (project_id: string, organization_id?: string) => {
  return queryOptions({
    queryKey: ['projects', project_id, 'members'],
    queryFn: () => getAllProjectMembersFn(project_id, organization_id),
  })
}
