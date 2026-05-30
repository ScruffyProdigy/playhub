import { test, expect } from '@playwright/test'
import { setUserDisplayName, completeMagicSignIn } from './helpers/auth.js'

test.describe('Auth flow', () => {
  test('signs in with a magic link and logs out', async ({ page }) => {
    const email = `e2e-${Date.now()}@example.com`

    await page.goto('/')
    await expect(page.getByRole('heading', { name: 'Sign in' })).toBeVisible()

    await completeMagicSignIn(page, email)

    await expect(page.getByText(email)).toBeVisible()
    await expect(page.getByText(`${email.split('@')[0]} (new)`)).toBeVisible()

    await page.getByRole('button', { name: 'Log out' }).click()
    await expect(page.getByRole('heading', { name: 'Sign in' })).toBeVisible()
    await expect(page.getByText(email)).not.toBeVisible()
  })

  test('keeps an existing display name when a returning user signs in again', async ({ page }) => {
    const email = `returning-${Date.now()}@example.com`
    const customDisplayName = 'Returning Player'

    await page.goto('/')
    await completeMagicSignIn(page, email)

    setUserDisplayName(email, customDisplayName)

    await page.getByRole('button', { name: 'Log out' }).click()
    await expect(page.getByRole('heading', { name: 'Sign in' })).toBeVisible()

    await completeMagicSignIn(page, email)

    await expect(page.getByText(customDisplayName)).toBeVisible()
    await expect(page.getByText(`${email.split('@')[0]} (new)`)).not.toBeVisible()
  })
})
