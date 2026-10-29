import { afterEach, describe, expect, it, vi } from "vitest";
import { preprocessImageUpload, uploadFile } from "../../src/features/storage/api/index";

afterEach(() => vi.restoreAllMocks());

describe("storage object references", () => {
  it("returns only the key after a signed upload", async () => {
    const fetchSpy = vi.spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(new Response(JSON.stringify({
        uploadUrl: "https://s3.test/bucket/signatures/image.png?signature=abc",
        key: "signatures/image.png",
      }), { status: 200 }))
      .mockResolvedValueOnce(new Response(null, { status: 200 }));

    const key = await uploadFile(new File(["image"], "image.png", { type: "image/png" }), "signatures");
    expect(key).toBe("signatures/image.png");
    expect(fetchSpy).toHaveBeenLastCalledWith(
      "https://s3.test/bucket/signatures/image.png?signature=abc",
      expect.objectContaining({ method: "PUT" }),
    );
  });

  it("rejects a signed-upload response without a key before sending the file", async () => {
    const fetchSpy = vi.spyOn(globalThis, "fetch").mockResolvedValue(new Response(JSON.stringify({
      uploadUrl: "https://s3.test/upload",
      publicUrl: "https://cdn.test/image.png",
    }), { status: 200 }));
    await expect(uploadFile(new File(["image"], "image.png", { type: "image/png" })))
      .rejects.toThrow("missing the object key");
    expect(fetchSpy).toHaveBeenCalledOnce();
  });

  it("returns the moderated upload key", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(new Response(JSON.stringify({
      approved: true,
      key: "profiles/images/avatar.png",
    }), { status: 200 }));
    expect(await preprocessImageUpload(new File(["image"], "image.png", { type: "image/png" })))
      .toBe("profiles/images/avatar.png");
  });
});
