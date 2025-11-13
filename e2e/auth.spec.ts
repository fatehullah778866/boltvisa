import { test, expect } from '@playwright/test'

test('signup shows friendly error if email exists', async ({ page }) => {
  await page.goto('http://localhost:3000/signup')
  
  // Fill in form with an email that might already exist
  await page.getByLabel(/first name/i).fill('Test')
  await page.getByLabel(/last name/i).fill('User')
  await page.getByLabel(/email/i).fill('already@exists.com')
  await page.getByLabel(/password/i).fill('Password123!')
  
  await page.getByRole('button', { name: /sign up/i }).click()
  
  // Should still be on signup page (not redirected)
  await expect(page).toHaveURL(/signup/)
  
  // Assert error message contains relevant text
  await expect(page.locator('body')).toContainText(/already registered|exists/i)
});

test.describe('Authentication Flow', () => {
  test('should register a new user', async ({ page }) => {
    await page.goto('/signup');
    
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[name="firstName"]', 'Test');
    await page.fill('input[name="lastName"]', 'User');
    await page.fill('input[type="password"]', 'password123');
    
    await page.click('button[type="submit"]');
    
    // Should redirect to dashboard
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should login with existing user', async ({ page }) => {
    await page.goto('/login');
    
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', 'password123');
    
    await page.click('button[type="submit"]');
    
    // Should redirect to dashboard
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should show error for invalid credentials', async ({ page }) => {
    await page.goto('/login');
    
    await page.fill('input[type="email"]', 'wrong@example.com');
    await page.fill('input[type="password"]', 'wrongpassword');
    
    await page.click('button[type="submit"]');
    
    // Should show error message
    await expect(page.locator('text=/error|invalid/i')).toBeVisible();
  });
});

