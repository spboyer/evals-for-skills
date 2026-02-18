import { test, expect } from '@playwright/test';
import { mockAPI } from './helpers/api-mock';
import { mockRunDetail, mockHealth } from './fixtures/mock-data';

test.describe('Run Detail', () => {
  test('navigates to run detail from dashboard', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/');
    await page.getByRole('link', { name: 'code-explainer' }).click();
    await expect(page).toHaveURL(/\/runs\/run-001/);
    await expect(page.getByRole('heading', { name: 'code-explainer' })).toBeVisible();
  });

  test('displays summary cards with run info', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/runs/run-001');
    await expect(page.getByText('gpt-4o')).toBeVisible();
    await expect(page.getByText('passed').first()).toBeVisible();
    await expect(page.getByText('4/4')).toBeVisible();
  });

  test('lists all tasks', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/runs/run-001');
    await expect(page.getByText('explain-function')).toBeVisible();
    await expect(page.getByText('explain-class')).toBeVisible();
    await expect(page.getByText('explain-module')).toBeVisible();
    await expect(page.getByText('explain-api')).toBeVisible();
  });

  test('task badges show correct variant', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/runs/run-001');
    const taskList = page.locator('ul');
    const items = taskList.locator('li');
    await expect(items).toHaveCount(4);
    // explain-module is "failed"
    const failedItem = items.filter({ hasText: 'explain-module' });
    await expect(failedItem.locator('span').last()).toHaveText('failed');
  });

  test('back link returns to dashboard', async ({ page }) => {
    await mockAPI(page);
    await page.goto('/runs/run-001');
    await page.getByRole('link', { name: /back to dashboard/i }).click();
    await expect(page).toHaveURL('/');
  });

  test('shows not found for missing run', async ({ page }) => {
    // Mock health and return null for the run detail to simulate not-found
    await page.route('**/api/health', route =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(mockHealth) })
    );
    await page.route('**/api/runs/nonexistent', route =>
      route.fulfill({ status: 200, contentType: 'application/json', body: 'null' })
    );
    await page.goto('/runs/nonexistent');
    await expect(page.getByText('Run not found.')).toBeVisible();
  });
});
