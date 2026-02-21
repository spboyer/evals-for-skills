import { test, expect } from "@playwright/test";
import { mockAllAPIs } from "./helpers/api-mock";

test.describe("Confidence Intervals & Significance", () => {
  test("shows significance badges in run detail task table", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    const badges = page.locator('[data-testid="significance-badge"]');
    // 3 tasks have isSignificant set in mock data (fibonacci, quicksort, binary-search)
    await expect(badges).toHaveCount(3);

    // First two tasks are significant
    await expect(badges.nth(0)).toContainText("significant");
    await expect(badges.nth(1)).toContainText("significant");

    // Third task is not significant
    await expect(badges.nth(2)).toContainText("not significant");
  });

  test("shows CI range inline for tasks with bootstrap CI", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    const ciRanges = page.locator('[data-testid="ci-range"]');
    // 3 tasks have bootstrapCI set
    await expect(ciRanges).toHaveCount(3);

    // explain-fibonacci: lower=0.82, upper=0.98 → [82.0%, 98.0%]
    await expect(ciRanges.nth(0)).toContainText("[82.0%, 98.0%]");

    // explain-binary-search: lower=-0.05, upper=0.15 → [-5.0%, 15.0%]
    await expect(ciRanges.nth(2)).toContainText("[-5.0%, 15.0%]");
  });

  test("does not show badges for tasks without CI data", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    // explain-merge-sort has no CI data — should not render badge or range
    // There are 4 task rows total but only 3 badges/ranges
    const badges = page.locator('[data-testid="significance-badge"]');
    const ciRanges = page.locator('[data-testid="ci-range"]');
    await expect(badges).toHaveCount(3);
    await expect(ciRanges).toHaveCount(3);
  });

  test("significant badge has green styling", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    const significantBadge = page.locator('[data-testid="significance-badge"]').first();
    await expect(significantBadge).toContainText("✓ significant");
    await expect(significantBadge).toHaveClass(/text-green-400/);
  });

  test("not-significant badge has yellow styling", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    const badges = page.locator('[data-testid="significance-badge"]');
    const notSigBadge = badges.nth(2);
    await expect(notSigBadge).toContainText("⚠ not significant");
    await expect(notSigBadge).toHaveClass(/text-yellow-400/);
  });

  test("CI range has tooltip with full interval", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    const ciRange = page.locator('[data-testid="ci-range"]').first();
    await expect(ciRange).toHaveAttribute("title", /95% CI/);
  });
});
