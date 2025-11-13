import { test, expect } from '@playwright/test'

test('signup shows clear error for existing email', async ({ page }) => {
  await page.goto('http://localhost:3000/signup')
  
  await page.getByPlaceholder(/email/i).fill('already@exists.com')
  await page.getByPlaceholder(/password/i).fill('Passw0rd!')
  await page.getByLabel(/first name/i).fill('Test')
  await page.getByLabel(/last name/i).fill('User')
  
  await page.getByRole('button', { name: /sign up/i }).click()
  
  await expect(page.locator('body')).toContainText(/already registered|exists/i)
})

test('login shows clear invalid credentials error', async ({ page }) => {
  await page.goto('http://localhost:3000/login')
  
  await page.getByPlaceholder(/email/i).fill('nope@example.com')
  await page.getByPlaceholder(/password/i).fill('wrongpass')
  
  await page.getByRole('button', { name: /log in/i }).click()
  
  await expect(page.locator('body')).toContainText(/invalid.*credentials|invalid.*email|invalid.*password/i)
})

