import { test, expect } from '@playwright/test';
import { mockAPI } from './helpers/api-mock';

test.describe('Theme', () => {
  test('layout structure renders correctly', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    const shell = page.locator('.min-h-screen');
    await expect(shell).toBeVisible();

    const nav = page.locator('header');
    await expect(nav).toBeVisible();

    const footer = page.locator('footer');
    await expect(footer).toBeVisible();
    await expect(footer).toContainText('waza');
  });

  test('dashboard page visual snapshot', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    await expect(page.getByText('Total Runs')).toBeVisible();
    await expect(page.getByText('code-explainer')).toBeVisible();

    await expect(page).toHaveScreenshot('dashboard.png', {
      maxDiffPixelRatio: 0.05,
    });
  });

  test('run detail page visual snapshot', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/runs/run-001');

    await expect(page.getByRole('heading', { name: 'code-explainer' })).toBeVisible();

    await expect(page).toHaveScreenshot('run-detail.png', {
      maxDiffPixelRatio: 0.05,
    });
  });
});
