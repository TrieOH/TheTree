import handler from "@tanstack/react-start/server-entry";
import {
  S3Client,
  PutObjectCommand,
  DeleteObjectCommand,
  GetObjectCommand,
} from "@aws-sdk/client-s3";
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";

export interface Env {
  MINIO_ENDPOINT: string;
  MINIO_ACCESS_KEY: string;
  MINIO_SECRET_KEY: string;
  BUCKET_NAME: string;
  AI: Ai;
}

const ALLOWED_TYPES = ["image/png", "image/jpeg", "image/webp"];
const MAX_SIZE = 10 * 1024 * 1024; // 10MB

export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    const url = new URL(request.url);

    if (url.pathname === "/storage/upload" && request.method === "POST") {
      const { filename, contentType, size } = await request.json();

      if (!ALLOWED_TYPES.includes(contentType)) {
        return Response.json({ error: "Only PNG, JPEG, and WebP are allowed" }, { status: 400 });
      }

      if (size > MAX_SIZE) {
        return Response.json({ error: "File exceeds 10MB limit" }, { status: 400 });
      }

      const uploadUrl = await getUploadUrl(filename, contentType, env);
      return Response.json({ uploadUrl, key: filename });
    }

    if (url.pathname === "/storage/moderate" && request.method === "POST") {
      const { key } = await request.json();
      const isSafe = await moderateFile(key, env);
      if (!isSafe) {
        await deleteFile(key, env);
        return Response.json({ approved: false });
      }
      return Response.json({ approved: true });
    }

    return handler.fetch(request, env, ctx);
  },
};

const getS3Client = (env: Env) =>
  new S3Client({
    region: "auto",
    endpoint: env.MINIO_ENDPOINT,
    credentials: {
      accessKeyId: env.MINIO_ACCESS_KEY,
      secretAccessKey: env.MINIO_SECRET_KEY,
    },
    forcePathStyle: true,
  });

async function getUploadUrl(key: string, contentType: string, env: Env): Promise<string> {
  const client = getS3Client(env);
  return getSignedUrl(
    client,
    new PutObjectCommand({ Bucket: env.BUCKET_NAME, Key: key, ContentType: contentType }),
    { expiresIn: 300 }
  );
}

async function moderateFile(key: string, env: Env): Promise<boolean> {
  const client = getS3Client(env);
  const url = await getSignedUrl(
    client,
    new GetObjectCommand({ Bucket: env.BUCKET_NAME, Key: key }),
    { expiresIn: 60 }
  );

  const fileRes = await fetch(url);
  const buffer = await fileRes.arrayBuffer();

  const response = await env.AI.run("@cf/llava-hf/llava-1.5-7b-hf", {
    prompt:
      "Does this image contain any explicit, violent, or inappropriate content? Reply with only 'safe' or 'unsafe'.",
    image: [...new Uint8Array(buffer)],
  });

  const result = (response as { description: string }).description.trim().toLowerCase();
  return result.startsWith("safe");
}

async function deleteFile(key: string, env: Env): Promise<void> {
  const client = getS3Client(env);
  await client.send(new DeleteObjectCommand({ Bucket: env.BUCKET_NAME, Key: key }));
}