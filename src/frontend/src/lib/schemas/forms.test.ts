import { describe, expect, it } from "vite-plus/test";
import { gradeBatchSchema, scoreInput } from "#lib/schemas/academics.ts";
import { documentInputSchema } from "#lib/schemas/documents.ts";
import { amountInput, paymentFormSchema } from "#lib/schemas/tuition.ts";
import { idSchema } from "#lib/schemas/shared.ts";
describe("form contracts", () => {
  it("keeps blank scores distinct from a genuine zero", () => {
    expect(scoreInput.safeParse("").success).toBe(false);
    expect(scoreInput.parse("0")).toBe("0");
    expect(scoreInput.safeParse("10.01").success).toBe(false);
    expect(scoreInput.safeParse("8.125").success).toBe(false);
  });
  it("rejects duplicate grade identities and an empty batch", () => {
    expect(gradeBatchSchema.safeParse({ grades: [] }).success).toBe(false);
    expect(
      gradeBatchSchema.safeParse({
        grades: [
          { enrollment_id: "8", score: "5" },
          { enrollment_id: "8", score: "6" },
        ],
      }).success,
    ).toBe(false);
  });
  it("matches the backend Unicode code-point limit and trims outer whitespace", () => {
    expect(
      documentInputSchema.parse({ request_type: "transcript", purpose: "  cần hồ sơ  " }).purpose,
    ).toBe("cần hồ sơ");
    expect(
      documentInputSchema.safeParse({ request_type: "other", purpose: "😀".repeat(1000) }).success,
    ).toBe(true);
    expect(
      documentInputSchema.safeParse({ request_type: "other", purpose: "😀".repeat(1001) }).success,
    ).toBe(false);
  });
  it("accepts only whole positive bounded VND amounts and bigint string IDs", () => {
    for (const invalid of ["", "0", "1.5", "1000000001", "1e3"])
      expect(amountInput.safeParse(invalid).success).toBe(false);
    expect(amountInput.parse("400000")).toBe("400000");
    expect(idSchema.parse("9007199254740993")).toBe("9007199254740993");
    expect(idSchema.safeParse("9223372036854775808").success).toBe(false);
    for (const malformed of ["bad", "0", "-1", "1".repeat(1000)])
      expect(idSchema.safeParse(malformed).success).toBe(false);
  });
  it("validates campus-time payment input before the server converts it to UTC", () => {
    const input = {
      amount: "400000",
      paid_at: "2026-10-26T15:30",
      note: "",
      voiding: false,
      balance_vnd: 5750000,
      revision: "",
    };
    expect(paymentFormSchema.safeParse(input).success).toBe(true);
    expect(paymentFormSchema.safeParse({ ...input, paid_at: "bad-date" }).success).toBe(false);
    expect(paymentFormSchema.safeParse({ ...input, paid_at: "2026-02-30T15:30" }).success).toBe(
      false,
    );
  });
});
