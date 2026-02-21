import { test, expect } from "@playwright/test";
import { mockAllAPIs } from "./helpers/api-mock";

test.describe("Judge Model", () => {
  test("run detail shows judge model badge when present", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-002");

    // Should display the judge model badge
    const badge = page.getByTestId("judge-model-badge");
    await expect(badge).toBeVisible();
    await expect(badge).toHaveText("Judge: claude-opus-4.6");
  });

  test("run detail hides judge model badge when absent", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    // run-001 has no judgeModel — badge should not exist
    await expect(page.getByRole("heading", { name: "code-explainer" })).toBeVisible();
    await expect(page.getByTestId("judge-model-badge")).not.toBeVisible();
  });

  test("runs table shows judge indicator when model differs", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/");

    // run-002 (skill-checker) has judgeModel different from model → ⚖ indicator
    const rows = page.locator("tbody tr");
    const run002Row = rows.nth(1);
    await expect(run002Row.locator("text=⚖")).toBeVisible();

    // run-001 (code-explainer) has no judgeModel → no indicator
    const run001Row = rows.nth(0);
    await expect(run001Row.locator("text=⚖")).not.toBeVisible();
  });
});
