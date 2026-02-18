import type { Page } from '@playwright/test';
import {
  mockSummary,
  mockRuns,
  mockRunDetail,
  mockHealth,
  mockEmptyRuns,
  mockEmptySummary,
} from '../fixtures/mock-data';

/** Intercept all API routes with default mock data. */
export async function mockAPI(page: Page) {
  await page.route('**/api/health', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(mockHealth),
    }),
  );

  await page.route('**/api/summary', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(mockSummary),
    }),
  );

  await page.route('**/api/runs', (route) => {
    const url = new URL(route.request().url());
    if (url.pathname === '/api/runs') {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockRuns),
      });
    }
    return route.continue();
  });

  await page.route('**/api/runs/*', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(mockRunDetail),
    }),
  );
}

/** Mock API with empty data to test empty states. */
export async function mockEmptyAPI(page: Page) {
  await page.route('**/api/health', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(mockHealth),
    }),
  );

  await page.route('**/api/summary', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(mockEmptySummary),
    }),
  );

  await page.route('**/api/runs', (route) => {
    const url = new URL(route.request().url());
    if (url.pathname === '/api/runs') {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockEmptyRuns),
      });
    }
    return route.continue();
  });

  await page.route('**/api/runs/*', (route) =>
    route.fulfill({
      status: 404,
      contentType: 'application/json',
      body: JSON.stringify({ error: 'Run not found' }),
    }),
  );
}
