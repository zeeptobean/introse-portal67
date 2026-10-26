import { test, expect, type Page } from "@playwright/test";
async function login(page: Page, role = "student") {
  await page.goto("/sign-in");
  await page.getByLabel("Email", { exact: true }).fill(`${role}1@portal67.test`);
  await page.getByLabel("Password", { exact: true }).fill("portal67");
  await page.getByRole("button", { name: "Sign in", exact: true }).click();
  await expect(page).toHaveURL(new RegExp(`/${role}(?:/|$)`));
}
async function logout(page: Page) {
  await page
    .locator(".account-block")
    .getByRole("button", { name: "Sign out", exact: true })
    .click();
  await expect(page).toHaveURL(/\/sign-in$/);
}

test("SSR works without JavaScript; session cookies rotate and unauthorized actions fail closed", async ({
  browser,
  baseURL,
}) => {
  const context = await browser.newContext({ javaScriptEnabled: false });
  const page = await context.newPage();
  await login(page);
  await expect(page.getByRole("heading", { name: "Today", exact: true })).toBeVisible();
  const session = (await context.cookies()).find((c) => c.name === "session");
  expect(session).toMatchObject({ httpOnly: true, secure: true, sameSite: "Strict" });
  const forbidden = await context.request.post("/staff/documents/4", {
    form: { status: "approved", comment: "Unauthorized test" },
    headers: { Origin: baseURL! },
  });
  expect(forbidden.status()).toBe(403);
  await page.goto("/student/results");
  await expect(page.getByRole("heading", { name: "Course results" })).toBeVisible();
  await logout(page);
  expect((await context.cookies()).find((c) => c.name === "session")).toBeUndefined();
  await login(page);
  expect((await context.cookies()).find((c) => c.name === "session")?.value).not.toBe(
    session?.value,
  );
  await context.close();
});

test("enhanced forms preserve the student/staff/lecturer journeys without browser Go requests", async ({
  page,
}) => {
  const directApi: string[] = [];
  page.on("request", (req) => {
    if (new URL(req.url()).pathname.startsWith("/api/")) directApi.push(req.url());
  });
  // Record unexpected draft prompts rather than letting a successful save hang.
  const dialogs: string[] = [];
  page.on("dialog", async (dialog) => {
    dialogs.push(`${dialog.type()} ${page.url()} ${dialog.message()}`);
    await dialog.accept();
  });
  await login(page);
  await page.goto("/student/registration?section=9#section-preview");
  await page
    .locator(".preview-panel")
    .getByRole("button", { name: "Register for section 01", exact: true })
    .click();
  await expect(page.locator(".registered-row").filter({ hasText: "CSC15001" })).toBeVisible();
  await page
    .locator(".registered-row")
    .filter({ hasText: "CSC15001" })
    .getByRole("button", { name: "Cancel", exact: true })
    .click();
  await page.getByRole("button", { name: "Cancel registration", exact: true }).click();
  await expect(page.getByText("Your seat has been released.", { exact: false })).toBeVisible();
  await expect(page.locator(".registered-row").filter({ hasText: "CSC15001" })).toHaveCount(0);
  await page.goto("/student/documents/new");
  await page.getByRole("button", { name: "Submit request", exact: true }).click();
  await expect(page.locator("[role=alert]")).toContainText("Enter a value");
  const purpose = `Server migration check ${Date.now()}`;
  await page.getByLabel("Purpose (required)").fill(purpose);
  await page.getByRole("button", { name: "Submit request", exact: true }).click();
  await expect(page).toHaveURL(/\/student\/documents\/\d+$/);
  const requestId = page.url().split("/").pop()!;
  await expect(page.locator(".status").first()).toHaveText(/pending/i);
  await logout(page);
  await login(page, "staff");
  await page.goto(`/staff/documents/${requestId}`);
  await page
    .getByLabel("Comment visible to the student (required)")
    .fill("Approved in the server migration smoke test.");
  await page.getByRole("button", { name: "Review decision", exact: true }).click();
  await page.getByRole("button", { name: "Confirm decision", exact: true }).click();
  await expect(page.locator(".status").first()).toHaveText(/approved/i);
  await page.goto("/staff/tuition");
  await page.getByRole("searchbox").fill("24120001");
  await page.getByRole("button", { name: "Find student", exact: true }).click();
  await page.getByRole("link", { name: "View tuition →" }).click();
  await page.getByRole("link", { name: "Record payment", exact: true }).click();
  await page.getByLabel("Amount (VND)").fill("400000");
  await page.getByLabel("Note (optional, 500 characters maximum)").fill(`Payment for ${requestId}`);
  await page.getByRole("button", { name: "Review payment", exact: true }).click();
  await page.getByRole("button", { name: "Confirm and save", exact: true }).click();
  await expect(page).toHaveURL(/saved=\d+/);
  const paymentId = new URL(page.url()).searchParams.get("saved")!;
  await page.locator(`a[href*='/payments/${paymentId}?']`).click();
  await page.getByLabel("Amount (VND)").fill("500000");
  await page.getByRole("button", { name: "Review payment", exact: true }).click();
  await page.getByRole("button", { name: "Confirm and save", exact: true }).click();
  await expect(page).toHaveURL(/saved=\d+/);
  await page.locator(`a[href*='/payments/${paymentId}?']`).click();
  await page.getByLabel("Void this payment").check();
  await page.getByRole("button", { name: "Review void", exact: true }).click();
  await page.getByRole("button", { name: "Confirm and save", exact: true }).click();
  await expect(
    page.locator("tr.voided").filter({ hasText: `Payment for ${requestId}` }),
  ).toContainText("Excluded from totals");
  await page.getByRole("link", { name: "Add charge", exact: true }).click();
  await page.getByLabel("Amount (VND)").fill("100000");
  await page.getByLabel("Charge description (required)").fill(`Charge for ${requestId}`);
  await page.getByRole("button", { name: "Review charge", exact: true }).click();
  await page.getByRole("button", { name: "Confirm and save", exact: true }).click();
  await expect(page).toHaveURL(/saved=\d+/);
  await expect(
    page.locator("tbody tr").filter({ hasText: `Charge for ${requestId}` }),
  ).toBeVisible();
  await logout(page);
  await login(page, "lecturer");
  await page.goto("/lecturer/sections/4/grades");
  const previousScore = await page.getByRole("spinbutton").inputValue();
  const nextScore = previousScore === "8.25" ? "8.50" : "8.25";
  await page.getByRole("spinbutton").fill("12");
  await page.getByRole("button", { name: "Save 1 grades", exact: true }).click();
  await expect(page.locator("[role=alert]")).toContainText("score from 0 to 10");
  await page.getByRole("spinbutton").fill(nextScore);
  await page.getByRole("button", { name: "Save 1 grades", exact: true }).click();
  await expect(page.getByText("1 grades saved.", { exact: false })).toBeVisible();
  await expect(page.getByRole("button", { name: "Save 0 grades", exact: true })).toBeDisabled();
  await logout(page);
  await login(page);
  await page.goto("/student/results");
  await expect(page.locator("tbody tr").filter({ hasText: "CSC13002" })).toContainText(
    Number(nextScore).toFixed(2),
  );
  await page.goto(`/student/documents/${requestId}`);
  await expect(page.locator(".status").first()).toHaveText(/approved/i);
  await page.setViewportSize({ width: 320, height: 800 });
  await page.goto("/student/registration?section=3");
  await expect(page.getByText("This section overlaps your timetable")).toBeVisible();
  expect(await page.evaluate(() => document.documentElement.scrollWidth > innerWidth)).toBe(false);
  expect(directApi).toEqual([]);
  expect(dialogs).toEqual([]);
  await page.screenshot({ path: "test-results/migration-registration-mobile.png", fullPage: true });
});

test("a lost form response preserves the draft and blocks a blind repeat", async ({ page }) => {
  await login(page);
  await page.goto("/student/documents/new");
  const purpose = "Uncertain action response test";
  await page.getByLabel("Purpose (required)").fill(purpose);
  await expect(page.locator(".review-box")).toContainText(purpose);
  let submissions = 0;
  await page.route("**/student/documents/new", async (route) => {
    if (route.request().method() === "POST") {
      submissions++;
      await route.abort("failed");
    } else await route.continue();
  });
  await page.getByRole("button", { name: "Submit request", exact: true }).click();
  await expect(page.getByText("Do not repeat the submission yet.", { exact: false })).toBeVisible();
  await expect(page.getByRole("button", { name: "Submit request", exact: true })).toBeDisabled();
  await expect(page.getByLabel("Purpose (required)")).toHaveValue(purpose);
  expect(submissions).toBe(1);
});
