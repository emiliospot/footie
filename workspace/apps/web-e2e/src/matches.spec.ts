import { test, expect } from '@playwright/test';

test.describe('Matches & Analytics', () => {
  test.beforeEach(async ({ page }) => {
    // Login as analyst
    await page.goto('/auth/login');
    await page.getByLabel(/email/i).fill('analyst@example.com');
    await page.getByLabel(/password/i).fill('password123');
    await page.getByRole('button', { name: /login/i }).click();
    await page.waitForURL(/\/dashboard/);
  });

  test('should display matches list', async ({ page }) => {
    await page.goto('/matches');

    await expect(page.getByRole('heading', { name: /matches/i })).toBeVisible();
    await expect(page.getByTestId('matches-list')).toBeVisible();
  });

  test('should filter matches by competition', async ({ page }) => {
    await page.goto('/matches');

    // Filter by Premier League
    await page.getByLabel(/competition/i).selectOption('Premier League');

    // Should show only Premier League matches
    const matches = page.getByTestId('match-card');
    await expect(matches.first()).toContainText(/premier league/i);
  });

  test('should filter matches by date range', async ({ page }) => {
    await page.goto('/matches');

    // Select date range
    await page.getByLabel(/from date/i).fill('2024-01-01');
    await page.getByLabel(/to date/i).fill('2024-03-31');
    await page.getByRole('button', { name: /apply/i }).click();

    // Should show matches in date range
    const matches = page.getByTestId('match-card');
    await expect(matches).toHaveCount(10, { timeout: 5000 });
  });

  test('should view match detail with analytics', async ({ page }) => {
    await page.goto('/matches');

    // Click on a match
    await page.getByTestId('match-card').first().click();

    // Should show match details
    await expect(page).toHaveURL(/\/matches\/\d+/);
    await expect(page.getByTestId('match-score')).toBeVisible();
    await expect(page.getByTestId('match-stats')).toBeVisible();
    await expect(page.getByTestId('match-events')).toBeVisible();

    // Should show analytics charts
    await expect(page.getByTestId('possession-chart')).toBeVisible();
    await expect(page.getByTestId('shots-chart')).toBeVisible();
  });

  test('should view player statistics', async ({ page }) => {
    await page.goto('/matches/1'); // Specific match

    // Navigate to player stats tab
    await page.getByRole('tab', { name: /player statistics/i }).click();

    // Should show player stats table
    await expect(page.getByTestId('player-stats-table')).toBeVisible();

    // Should show sortable columns
    await page.getByRole('columnheader', { name: /goals/i }).click();

    // Top scorer should be first
    const firstRow = page.getByTestId('player-row').first();
    await expect(firstRow).toBeVisible();
  });

  test('should compare two teams', async ({ page }) => {
    await page.goto('/matches/1');

    // Click compare button
    await page.getByRole('button', { name: /compare teams/i }).click();

    // Should show comparison view
    await expect(page.getByTestId('team-comparison')).toBeVisible();
    await expect(page.getByTestId('comparison-chart')).toBeVisible();

    // Should show head-to-head stats
    await expect(page.getByText(/head to head/i)).toBeVisible();
    await expect(page.getByText(/last 5 matches/i)).toBeVisible();
  });

  test('should export match data', async ({ page }) => {
    await page.goto('/matches/1');

    // Click export button
    const downloadPromise = page.waitForEvent('download');
    await page.getByRole('button', { name: /export/i }).click();
    await page.getByRole('menuitem', { name: /csv/i }).click();

    // Should download file
    const download = await downloadPromise;
    expect(download.suggestedFilename()).toContain('match');
    expect(download.suggestedFilename()).toContain('.csv');
  });
});

