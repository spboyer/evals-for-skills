import { test } from '@playwright/test';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const outputDir = path.resolve(__dirname, '../../docs/images');

const pages = [
  { name: 'home', path: '/' },
  { name: 'getting-started', path: '/getting-started/' },
  { name: 'graders', path: '/guides/graders/' },
  { name: 'cli-reference', path: '/reference/cli/' },
];

for (const pg of pages) {
  test(`screenshot: ${pg.name}`, async ({ page }) => {
    await page.goto(pg.path);
    await page.waitForLoadState('networkidle');
    await page.screenshot({
      path: path.join(outputDir, `site-${pg.name}.png`),
      fullPage: true,
    });
  });
}
