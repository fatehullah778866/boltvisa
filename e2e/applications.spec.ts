import { test, expect } from '@playwright/test';

test.describe('Visa Applications', () => {
  test.beforeEach(async ({ page }) => {
    // Login first
    await page.goto('/login');
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/dashboard/);
  });

  test('should create a new application', async ({ page }) => {
    await page.goto('/dashboard');
    
    // Click "New Application" button
    await page.click('text=New Application');
    
    // Fill application form
    await page.selectOption('select', { index: 1 }); // Select first category
    await page.fill('input[name="passport_number"]', 'AB123456');
    await page.fill('input[type="date"][name="date_of_birth"]', '1990-01-01');
    await page.fill('input[name="nationality"]', 'US');
    
    // Submit form
    await page.click('button[type="submit"]');
    
    // Should redirect back to dashboard
    await expect(page).toHaveURL(/\/dashboard/);
    
    // Should see the new application
    await expect(page.locator('text=AB123456')).toBeVisible();
  });

  test('should display applications list', async ({ page }) => {
    await page.goto('/dashboard');
    
    // Should see applications section
    await expect(page.locator('text=My Applications')).toBeVisible();
  });

  test('should navigate to notifications', async ({ page }) => {
    await page.goto('/dashboard');
    
    await page.click('text=Notifications');
    
    await expect(page).toHaveURL(/\/notifications/);
  });

  test('should navigate to payments', async ({ page }) => {
    await page.goto('/dashboard');
    
    await page.click('text=Payments');
    
    await expect(page).toHaveURL(/\/payments/);
  });
});

