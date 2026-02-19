import { test, expect } from "@playwright/test";
import { mockAllAPIs } from "./helpers/api-mock";

test.describe("API Health", () => {
  test("/api/health returns 200 with status ok", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/");

    // Use page.evaluate so the fetch goes through the browser (route interception)
    const result = await page.evaluate(async () => {
      const res = await fetch("/api/health");
      return { status: res.status, body: await res.json() };
    });

    expect(result.status).toBe(200);
    expect(result.body.status).toBe("ok");
  });
});

test.describe("Theme", () => {
  test("dark theme applied — page has dark background", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/");

    // The root div with min-h-screen bg-zinc-900 defines the dark theme
    const rootDiv = page.locator("div.min-h-screen").first();
    await expect(rootDiv).toBeVisible();

    const bg = await rootDiv.evaluate((el) => getComputedStyle(el).backgroundColor);
    // Tailwind v4 uses oklch color space — zinc-900 is oklch(0.21 0.006 285.885)
    expect(bg).toMatch(/oklch|rgb/);
    // Verify it's a dark color (oklch lightness < 0.3, or rgb values < 50)
    if (bg.startsWith("oklch")) {
      const lightness = parseFloat(bg.match(/oklch\(([0-9.]+)/)?.[1] ?? "1");
      expect(lightness).toBeLessThan(0.3);
    } else {
      const [r, g, b] = bg.match(/\d+/g)!.map(Number);
      expect(Math.max(r, g, b)).toBeLessThan(50);
    }
  });

  test("nav header has dark border styling", async ({ page }) => {
    await mockAllAPIs(page);
    await page.goto("/");

    const header = page.locator("header");
    await expect(header).toBeVisible();

    const borderColor = await header.evaluate((el) =>
      getComputedStyle(el).borderBottomColor,
    );
    // border-zinc-800 — should be a dark gray, not transparent
    expect(borderColor).not.toBe("rgba(0, 0, 0, 0)");
  });
});
