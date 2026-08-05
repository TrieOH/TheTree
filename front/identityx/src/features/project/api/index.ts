import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  addOrganizationProjectMember,
  addProjectMember,
  createOrganizationProject,
  createProject,
  listOrganizationProjectMembers,
  listOrganizationProjects,
  listProjectMembers,
  listProjects,
  removeOrganizationProjectMember,
  removeProjectMember,
} from "@trieoh/identityx-api";
import type {
  MemberAddToProjectI,
  ProjectCreateI,
  ProjectI,
  ProjectMemberI,
} from "../model";

/**
 * Creates a new project on the server.
 * @param projectData - The data for the new project.
 * @param orgId - The organization ID to which the project belongs (optional).
 * @returns A promise that resolves to the API response containing the newly created project.
 */
export const createProjectFn = createClientOnlyFn(
  (projectData: ProjectCreateI, orgId?: string) => {
    if (orgId)
      return createOrganizationProject(orgId, projectData).then(
        orvalData<ProjectI>,
      );
    return createProject(projectData).then(orvalData<ProjectI>);
  },
);

/**
 * Fetches all projects from the server.
 * @param orgId - The organization ID to filter projects by (optional).
 * @returns A promise that resolves to an array of ProjectI objects.
 */
export const getProjectsFn = createClientOnlyFn(async (orgId?: string) => {
  if (orgId) return listOrganizationProjects(orgId).then(orvalData<ProjectI[]>);
  return listProjects().then(orvalData<ProjectI[]>);
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
export const addMemberToProjectFn = createClientOnlyFn(
  (
    project_id: string,
    memberData: MemberAddToProjectI,
    organization_id?: string,
  ) => {
    if (organization_id)
      return addOrganizationProjectMember(
        organization_id,
        project_id,
        memberData,
      ).then(orvalData<void>);
    return addProjectMember(project_id, memberData).then(orvalData<void>);
  },
);

/**
 * Removes a member from a organization on the server.
 * @param project_id - The ID of the project to remove the member from.
 * @param actor_email - The email of the user to remove from the project.
 * @param organization_id - The ID of the organization to remove the member from (optional).
 * @returns A promise that resolves to the API response confirming the removal of the member.
 */
export const removeMemberFromProjectFn = createClientOnlyFn(
  (project_id: string, actor_email: string, organization_id?: string) => {
    if (organization_id)
      return removeOrganizationProjectMember(organization_id, project_id, {
        actor_email,
      }).then(orvalData<void>);
    return removeProjectMember(project_id, { actor_email }).then(
      orvalData<void>,
    );
  },
);

/**
 * Fetches all project members from the server.
 * @param project_id - The ID of the project to fetch members for.
 * @param organization_id - The ID of the organization to fetch members for (optional).
 * @returns A promise that resolves to an array of project members objects.
 */
export const getAllProjectMembersFn = createClientOnlyFn(
  (project_id: string, organization_id?: string) => {
    if (organization_id)
      return listOrganizationProjectMembers(organization_id, project_id).then(
        orvalData<ProjectMemberI[]>,
      );
    return listProjectMembers(project_id).then(orvalData<ProjectMemberI[]>);
  },
);

/**
 * Query options for fetching all Project Members, using TanStack Query.
 * @param project_id - The ID of the project to fetch members for.
 * @param organization_id - The ID of the organization to fetch members for (optional).
 * @returns An object containing the query key and query function for fetching all Project Members.
 */
export const allProjectMembersQueryOptions = (
  project_id: string,
  organization_id?: string,
) => {
  return queryOptions({
    queryKey: ["projects", project_id, "members"],
    queryFn: () => getAllProjectMembersFn(project_id, organization_id),
  });
};
