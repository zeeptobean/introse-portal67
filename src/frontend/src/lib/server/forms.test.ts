import { describe, expect, it, vi } from "vite-plus/test";
import { ApiError } from "#lib/server/client.ts";
import { verifyUncertain } from "#lib/server/forms.ts";
describe("uncertain write verification", () => {
  it("reads current records only when the outcome is unknown", async () => {
    const check = vi.fn(async () => "saved");
    expect(await verifyUncertain(new ApiError("lost", 0, true), check)).toBe("saved");
    expect(await verifyUncertain(new ApiError("rejected", 400), check)).toBeUndefined();
    expect(await verifyUncertain(new Error("bug"), check)).toBeUndefined();
    expect(check).toHaveBeenCalledTimes(1);
  });
  it("checks a conflict only when asked to", async () => {
    const conflict = new ApiError("decided", 409);
    expect(await verifyUncertain(conflict, async () => true)).toBeUndefined();
    expect(await verifyUncertain(conflict, async () => true, { conflicts: true })).toBe(true);
  });
  it("keeps the write uncertain when the check read fails", async () => {
    const failing = async () => {
      throw new ApiError("unreachable", 0);
    };
    expect(await verifyUncertain(new ApiError("lost", 0, true), failing)).toBeUndefined();
  });
});
