import { test, expect } from '@playwright/test';
import { mockAPI } from './helpers/api-mock';

test.describe('Navigation', () => {
  test('root path loads dashboard', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
    await expect(page.getByText('Total Runs')).toBeVisible();
  });

  test('direct URL to /runs/:id works', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/runs/run-001');

    await expect(page.getByRole('heading', { name: 'code-explainer' })).toBeVisible();
    await expect(page.getByText('Back to dashboard')).toBeVisible();
  });

  test('nav logo links to dashboard', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/runs/run-001');

    await page.locator('header').getByText('waza').click();
    await expect(page).toHaveURL('/');
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
  });
});
