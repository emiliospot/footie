import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display login page', async ({ page }) => {
    await expect(page).toHaveTitle(/Footie/);
    await expect(page.getByRole('heading', { name: /welcome to footie/i })).toBeVisible();
  });

  test('should login with valid credentials', async ({ page }) => {
    // Navigate to login
    await page.goto('/auth/login');

    // Fill in credentials
    await page.getByLabel(/email/i).fill('test@example.com');
    await page.getByLabel(/password/i).fill('password123');

    // Submit form
    await page.getByRole('button', { name: /login/i }).click();

    // Should redirect to dashboard
    await expect(page).toHaveURL(/\/dashboard/);
    await expect(page.getByRole('heading', { name: /dashboard/i })).toBeVisible();
  });

  test('should show error with invalid credentials', async ({ page }) => {
    await page.goto('/auth/login');

    await page.getByLabel(/email/i).fill('invalid@example.com');
    await page.getByLabel(/password/i).fill('wrongpassword');
    await page.getByRole('button', { name: /login/i }).click();

    // Should show error message
    await expect(page.getByText(/invalid credentials/i)).toBeVisible();
  });

  test('should navigate to registration page', async ({ page }) => {
    await page.goto('/auth/login');

    await page.getByRole('link', { name: /register/i }).click();

    await expect(page).toHaveURL(/\/auth\/register/);
  });

  test('should logout successfully', async ({ page }) => {
    // Login first
    await page.goto('/auth/login');
    await page.getByLabel(/email/i).fill('test@example.com');
    await page.getByLabel(/password/i).fill('password123');
    await page.getByRole('button', { name: /login/i }).click();
    await page.waitForURL(/\/dashboard/);

    // Logout
    await page.getByRole('button', { name: /logout/i }).click();

    // Should redirect to login
    await expect(page).toHaveURL(/\/auth\/login/);
  });
});

