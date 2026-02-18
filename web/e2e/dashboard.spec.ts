import { test, expect } from '@playwright/test';
import { mockAPI, mockEmptyAPI } from './helpers/api-mock';

test.describe('Dashboard', () => {
  test('page loads with waza branding', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    const brand = page.locator('header').getByText('waza');
    await expect(brand).toBeVisible();
  });

  test('displays 4 KPI cards with data', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    await expect(page.getByText('Total Runs')).toBeVisible();
    await expect(page.getByText('Pass Rate')).toBeVisible();
    await expect(page.getByText('Avg Cost')).toBeVisible();
    await expect(page.getByText('Avg Duration')).toBeVisible();
  });

  test('KPI values are properly formatted', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    await expect(page.getByText('42', { exact: true })).toBeVisible();
    await expect(page.getByText('85.7%')).toBeVisible();
    await expect(page.getByText('$0.0342')).toBeVisible();
    await expect(page.getByText('4.5s')).toBeVisible();
  });

  test('runs table renders with column headers', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    const thead = page.locator('thead');
    await expect(thead.getByText('Spec')).toBeVisible();
    await expect(thead.getByText('Model')).toBeVisible();
    await expect(thead.getByText('Outcome')).toBeVisible();
    await expect(thead.getByText('Tasks')).toBeVisible();
    await expect(thead.getByText('Duration')).toBeVisible();
    await expect(thead.getByText('Cost')).toBeVisible();
  });

  test('runs table shows run data', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    await expect(page.getByText('code-explainer')).toBeVisible();
    await expect(page.getByText('bug-fixer')).toBeVisible();
    await expect(page.getByText('test-writer')).toBeVisible();
    await expect(page.getByText('gpt-4o').first()).toBeVisible();
  });

  test('pass/fail badges display correctly', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    await expect(page.getByText('Passed').first()).toBeVisible();
    await expect(page.getByText('Failed').first()).toBeVisible();
  });

  test('table shows task counts', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    await expect(page.getByText('3/4')).toBeVisible();
    await expect(page.getByText('2/4')).toBeVisible();
  });

  test('clicking spec link navigates to detail', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');

    await page.getByText('code-explainer').click();
    await expect(page).toHaveURL(/\/runs\/run-001/);
  });

  test('empty state shows when no runs', async ({ page }) => {
    await mockEmptyAPI(page);
    await page.goto('/');

    await expect(page.getByText('No evaluation runs yet.')).toBeVisible();
  });
});
