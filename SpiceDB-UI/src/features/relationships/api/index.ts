import { createServerFn } from "@tanstack/react-start";
import { z } from "zod";
import { spicedb } from "@soramux/node-perm-sdk";
import { queryOptions } from "@tanstack/react-query";

export const readRelationship = createServerFn({
  method: "GET",
})
  .inputValidator(z.object({
    envId: z.string(),
    resources: z.array(z.string()),
  }))
  .handler(async ({ data }) => {
    try {
      const service = spicedb.relationship(data.envId);

      const allResults = await Promise.all(
        data.resources.map(async (resource) => {
          const relationships = [];
          const stream = service.readStream({ resourceType: resource });
          for await (const response of stream) relationships.push(response);
          return relationships;
        })
      );

      return allResults.flat();
    } catch (error) {
      console.error("Erro ao ler stream do SpiceDB:", error);
      throw new Error(error instanceof Error ? error.message : "Erro ao ler o relationship");
    }
  })

export const relationshipQueryOptions = (envId: string, resources: string[]) => queryOptions({
  queryKey: ["relationship", envId, resources],
  queryFn: () => readRelationship({ data: { envId, resources } }),
  staleTime: 0,
});


export const createRelationship = createServerFn({
  method: "POST",
})
  .inputValidator(z.object({
    envId: z.string(),
    resource: z.string().min(1),
    resourceId: z.string().min(1),
    subject: z.string().min(1),
    subjectId: z.string().min(1),
    relation: z.string().min(1),
  }))
  .handler(async ({ data }) => {
    const response = await spicedb.relationship(data.envId).create(
      {
        resourceType: data.resource,
        resourceId: data.resourceId,
        subjectType: data.subject,
        subjectId: data.subjectId,
        relation: data.relation,
      }
    )
    if (!response.success) {
      return {
        success: false,
        message: response.message,
      }
    }
    return response
  })

export const updateRelationship = createServerFn({
  method: "POST",
})
  .inputValidator(z.object({
    envId: z.string(),
    resource: z.string().min(1),
    resourceId: z.string().min(1),
    subject: z.string().min(1),
    subjectId: z.string().min(1),
    relation: z.string().min(1),
  }))
  .handler(async ({ data }) => {
    const response = await spicedb.relationship(data.envId).update(
      {
        resourceType: data.resource,
        resourceId: data.resourceId,
        subjectType: data.subject,
        subjectId: data.subjectId,
        relation: data.relation,
      }
    )
    if (!response.success) {
      return {
        success: false,
        message: response.message,
      }
    }
    return response
  })

export const deleteRelationship = createServerFn({
  method: "POST",
})
  .inputValidator(z.object({
    envId: z.string(),
    resource: z.string().min(1),
    resourceId: z.string().min(1),
    subject: z.string().min(1),
    subjectId: z.string().min(1),
    relation: z.string().min(1),
  }))
  .handler(async ({ data }) => {
    const response = await spicedb.relationship(data.envId).delete(
      {
        resourceType: data.resource,
        resourceId: data.resourceId,
        subjectType: data.subject,
        subjectId: data.subjectId,
        relation: data.relation,
      }
    )
    if (!response.success) {
      return {
        success: false,
        message: response.message,
      }
    }
    return response
  })

export const checkRelationship = createServerFn({
  method: "POST",
})
  .inputValidator(z.object({
    envId: z.string(),
    resource: z.string().min(1),
    resourceId: z.string().min(1),
    subject: z.string().min(1),
    subjectId: z.string().min(1),
    relation: z.string().min(1),
  }))
  .handler(async ({ data }) => {
    const response = await spicedb.relationship(data.envId).check(
      {
        resourceType: data.resource,
        resourceId: data.resourceId,
        subjectType: data.subject,
        subjectId: data.subjectId,
        permission: data.relation,
      }
    )
    if (!response.success) {
      return {
        success: false,
        message: response.message,
      }
    }
    return response
  })
