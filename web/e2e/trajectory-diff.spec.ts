import { test, expect } from "@playwright/test";
import { mockAllAPIs } from "./helpers/api-mock";

test.describe("Trajectory Diff", () => {
  test.beforeEach(async ({ page }) => {
    await mockAllAPIs(page);
  });

  test("compare view loads with run selectors", async ({ page }) => {
    await page.goto("/#/compare");

    await expect(page.getByRole("heading", { name: "Compare Runs" })).toBeVisible();

    // Two run selectors should be visible
    await expect(page.getByText("Run A")).toBeVisible();
    await expect(page.getByText("Run B")).toBeVisible();
  });

  test("selecting two runs shows comparison table", async ({ page }) => {
    await page.goto("/#/compare");

    // Select Run A
    const selects = page.locator("select");
    await selects.nth(0).selectOption("run-001");
    await selects.nth(1).selectOption("run-002");

    // Wait for comparison table to load
    await expect(page.getByText("Per-Task Comparison")).toBeVisible();
    await expect(page.getByText("explain-fibonacci")).toBeVisible();
  });

  test("task row is clickable and opens trajectory compare", async ({ page }) => {
    await page.goto("/#/compare");

    const selects = page.locator("select");
    await selects.nth(0).selectOption("run-001");
    await selects.nth(1).selectOption("run-002");

    // Wait for the comparison table
    await expect(page.getByText("Per-Task Comparison")).toBeVisible();

    // Click on explain-fibonacci row in the table
    const taskRow = page.locator("tbody tr").filter({ hasText: "explain-fibonacci" });
    await taskRow.click();

    // Should open TaskTrajectoryCompare
    await expect(page.getByText("Trajectory:")).toBeVisible();
    await expect(page.getByText("explain-fibonacci").last()).toBeVisible();
  });

  test("digest comparison shows token and turn deltas", async ({ page }) => {
    await page.goto("/#/compare");

    const selects = page.locator("select");
    await selects.nth(0).selectOption("run-001");
    await selects.nth(1).selectOption("run-002");

    await expect(page.getByText("Per-Task Comparison")).toBeVisible();

    // Click explain-fibonacci to open trajectory compare
    const taskRow = page.locator("tbody tr").filter({ hasText: "explain-fibonacci" });
    await taskRow.click();

    // Session Digest Comparison section should appear
    await expect(page.getByText("Session Digest Comparison")).toBeVisible();

    // Should show delta arrows (Run A has 3 turns, Run B has 4)
    await expect(page.getByText("Turns").last()).toBeVisible();
    await expect(page.getByText("Tool Calls").last()).toBeVisible();
  });

  test("diff entries render with correct labels", async ({ page }) => {
    await page.goto("/#/compare");

    const selects = page.locator("select");
    await selects.nth(0).selectOption("run-001");
    await selects.nth(1).selectOption("run-002");

    await expect(page.getByText("Per-Task Comparison")).toBeVisible();

    const taskRow = page.locator("tbody tr").filter({ hasText: "explain-fibonacci" });
    await taskRow.click();

    // Diff entries should show tool names and diff kinds
    // Both runs have read_file (changed — different result), write_file (changed — different args)
    // Run A has no run_tests, Run B has run_tests → Only in B
    await expect(page.getByText("read_file").first()).toBeVisible();
    await expect(page.getByText("write_file").first()).toBeVisible();
  });

  test("legend shows matched, changed, and missing counts", async ({ page }) => {
    await page.goto("/#/compare");

    const selects = page.locator("select");
    await selects.nth(0).selectOption("run-001");
    await selects.nth(1).selectOption("run-002");

    await expect(page.getByText("Per-Task Comparison")).toBeVisible();

    const taskRow = page.locator("tbody tr").filter({ hasText: "explain-fibonacci" });
    await taskRow.click();

    // Legend should show counts for matched, changed, missing
    await expect(page.getByText(/\d+ matched/)).toBeVisible();
    await expect(page.getByText(/\d+ changed/)).toBeVisible();
    await expect(page.getByText(/\d+ missing/)).toBeVisible();
  });

  test("clicking hint text is visible in comparison table", async ({ page }) => {
    await page.goto("/#/compare");

    const selects = page.locator("select");
    await selects.nth(0).selectOption("run-001");
    await selects.nth(1).selectOption("run-002");

    // Should show hint about clicking rows
    await expect(page.getByText("Click a row to view trajectory diff")).toBeVisible();
  });
});
