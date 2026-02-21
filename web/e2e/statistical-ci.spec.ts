import { test, expect } from "@playwright/test";
import { mockAllAPIs } from "./helpers/api-mock";

test.describe("Statistical Confidence Intervals", () => {
  test("run detail shows digest-level statistics bar", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    const statsBar = page.getByTestId("run-statistics");
    await expect(statsBar).toBeVisible();
    await expect(statsBar.getByText("✓ significant")).toBeVisible();
    await expect(statsBar.getByText(/\+8\.2%.*\+28\.8%/)).toBeVisible();
    await expect(statsBar.getByText(/Gain.*15\.0%/)).toBeVisible();
  });

  test("task row shows significance badge for significant task", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    await expect(page.getByText("explain-fibonacci")).toBeVisible();

    // The fibonacci task has isSignificant: true
    const fibRow = page.getByRole("row").filter({ hasText: "explain-fibonacci" });
    await expect(fibRow.getByText("✓ significant")).toBeVisible();
    await expect(fibRow.getByText(/\+12\.0%.*\+35\.0%/)).toBeVisible();
  });

  test("task row shows not-significant badge", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    await expect(page.getByText("explain-binary-search")).toBeVisible();

    // The binary-search task has isSignificant: false
    const bsRow = page.getByRole("row").filter({ hasText: "explain-binary-search" });
    await expect(bsRow.getByText("⚠ not significant")).toBeVisible();
    await expect(bsRow.getByText(/-10\.0%.*\+18\.0%/)).toBeVisible();
  });

  test("tasks without stats show no badge or CI", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/#/runs/run-001");

    await expect(page.getByText("explain-quicksort")).toBeVisible();

    // quicksort has no stats — should show no significance badge
    const qsRow = page.getByRole("row").filter({ hasText: "explain-quicksort" });
    await expect(qsRow.getByTestId("significance-badge")).not.toBeVisible();
  });
});
