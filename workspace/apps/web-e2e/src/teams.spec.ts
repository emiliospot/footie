import { test, expect } from '@playwright/test';

test.describe('Teams', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/auth/login');
    await page.getByLabel(/email/i).fill('analyst@example.com');
    await page.getByLabel(/password/i).fill('password123');
    await page.getByRole('button', { name: /login/i }).click();
    await page.waitForURL(/\/dashboard/);
  });

  test('should display teams list', async ({ page }) => {
    await page.goto('/teams');

    await expect(page.getByRole('heading', { name: /teams/i })).toBeVisible();
    // Should show teams grid or list
    await expect(page.getByTestId('teams-list')).toBeVisible();
  });

  test('should filter teams by country', async ({ page }) => {
    await page.goto('/teams');

    // Select country filter
    await page.getByLabel(/country/i).selectOption('England');

    // Should show only English teams
    await expect(page.getByText(/manchester/i)).toBeVisible();
    await expect(page.getByText(/liverpool/i)).toBeVisible();
  });

  test('should navigate to team detail page', async ({ page }) => {
    await page.goto('/teams');

    // Click on first team
    await page.getByRole('link', { name: /manchester united/i }).first().click();

    // Should show team details
    await expect(page).toHaveURL(/\/teams\/\d+/);
    await expect(page.getByRole('heading', { name: /manchester united/i })).toBeVisible();
    await expect(page.getByTestId('team-players')).toBeVisible();
    await expect(page.getByTestId('team-statistics')).toBeVisible();
  });

  test('should create new team (analyst role)', async ({ page }) => {
    await page.goto('/teams');

    await page.getByRole('button', { name: /add team/i }).click();

    // Fill in team form
    await page.getByLabel(/team name/i).fill('Test FC');
    await page.getByLabel(/short name/i).fill('TEST');
    await page.getByLabel(/code/i).fill('TST');
    await page.getByLabel(/country/i).fill('England');

    await page.getByRole('button', { name: /save/i }).click();

    // Should show success message and redirect
    await expect(page.getByText(/team created successfully/i)).toBeVisible();
    await expect(page).toHaveURL(/\/teams\/\d+/);
  });

  test('should search for teams', async ({ page }) => {
    await page.goto('/teams');

    await page.getByPlaceholder(/search teams/i).fill('Manchester');

    // Should filter results
    await expect(page.getByText(/manchester united/i)).toBeVisible();
    await expect(page.getByText(/manchester city/i)).toBeVisible();
    await expect(page.getByText(/liverpool/i)).not.toBeVisible();
  });
});

