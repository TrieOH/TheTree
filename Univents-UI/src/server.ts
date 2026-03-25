/// <reference path="../worker-configuration.d.ts" />
import handler from "@tanstack/react-start/server-entry";
import { AwsClient } from "aws4fetch";

const ALLOWED_TYPES = ["image/png", "image/jpeg", "image/webp"];
const MAX_SIZE = 10 * 1024 * 1024; // 10MB

export interface Env extends Cloudflare.Env {
  MINIO_ACCESS_KEY: string;
  MINIO_SECRET_KEY: string;
  MINIO_ENDPOINT: string;
  BUCKET_NAME: string;
  AI: Ai; // Cloudflare AI binding
}

// Initialize AWS client for MinIO/S3-compatible storage
const getAwsClient = (env: Env) =>
  new AwsClient({
    accessKeyId: env.MINIO_ACCESS_KEY,
    secretAccessKey: env.MINIO_SECRET_KEY,
    service: "s3",
    region: "us-east-1", // MinIO typically uses us-east-1 or auto
  });

export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    const url = new URL(request.url);

    if (url.pathname === "/storage/upload" && request.method === "POST") {
      const { filename, contentType, size } = (await request.json()) as {
        filename: string;
        contentType: string;
        size: number;
      };

      if (!ALLOWED_TYPES.includes(contentType)) {
        return Response.json({ error: "Only PNG, JPEG, and WebP are allowed" }, { status: 400 });
      }

      if (size > MAX_SIZE) {
        return Response.json({ error: "File exceeds 10MB limit" }, { status: 400 });
      }

      const uploadUrl = await getUploadUrl(filename, contentType, env);
      const publicUrl = `${env.MINIO_ENDPOINT}/${env.BUCKET_NAME}/${filename}`;
      return Response.json({ uploadUrl, key: filename, publicUrl });
    }

    if (url.pathname === "/storage/moderate" && request.method === "POST") {
      const { key } = (await request.json()) as { key: string };
      const isSafe = await moderateFile(key, env);
      if (!isSafe) {
        await deleteFile(key, env);
        return Response.json({ approved: false });
      }
      return Response.json({ approved: true });
    }

    // Fallback to TanStack handler
    // @ts-expect-error handler.fetch might only accept 2 arguments in some environments
    return handler.fetch(request, env, ctx);
  },
};

async function getUploadUrl(key: string, contentType: string, env: Env): Promise<string> {
  const aws = getAwsClient(env);
  const uploadUrl = new URL(`${env.MINIO_ENDPOINT}/${env.BUCKET_NAME}/${key}`);
  uploadUrl.searchParams.set("X-Amz-Expires", "300"); // 5 minutes

  const signed = await aws.sign(
    new Request(uploadUrl, {
      method: "PUT",
      headers: { "Content-Type": contentType },
    }),
    { aws: { signQuery: true } }
  );

  return signed.url.toString();
}

async function moderateFile(key: string, env: Env): Promise<boolean> {
  const aws = getAwsClient(env);
  const downloadUrl = new URL(`${env.MINIO_ENDPOINT}/${env.BUCKET_NAME}/${key}`);
  downloadUrl.searchParams.set("X-Amz-Expires", "60"); // 1 minute

  const signed = await aws.sign(
    new Request(downloadUrl, { method: "GET" }),
    { aws: { signQuery: true } }
  );

  const fileRes = await fetch(signed.url);
  if (!fileRes.ok) throw new Error("Failed to fetch file for moderation");

  const buffer = await fileRes.arrayBuffer();

  const response = await env.AI.run("@cf/llava-hf/llava-1.5-7b-hf", {
    prompt: "Does this image contain any explicit, violent, or inappropriate content? Reply with only 'safe' or 'unsafe'.",
    image: [...new Uint8Array(buffer)],
    max_tokens: 5,
  });

  const result = (response as { description: string }).description.trim().toLowerCase();
  return result.startsWith("safe");
}

async function deleteFile(key: string, env: Env): Promise<void> {
  const aws = getAwsClient(env);
  const deleteUrl = new URL(`${env.MINIO_ENDPOINT}/${env.BUCKET_NAME}/${key}`);

  const signed = await aws.sign(
    new Request(deleteUrl, { method: "DELETE" }),
    { aws: { signQuery: true } }
  );

  const res = await fetch(signed.url, { method: "DELETE" });
  if (!res.ok) throw new Error(`Failed to delete file: ${res.statusText}`);
}