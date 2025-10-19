import { test, expect } from '@playwright/test'

test.describe('Minimal App Tests', () => {
  test('app loads without errors', async ({ page }) => {
    // Just check that the page loads without throwing errors
    await page.goto('/', { timeout: 30000 })
    
    // Wait for basic page load
    await page.waitForLoadState('domcontentloaded', { timeout: 10000 })
    
    // Check that we get some content (not an error page)
    const content = await page.content()
    expect(content.length).toBeGreaterThan(100)
    
    // Check that it's not an error page
    expect(content).not.toContain('404')
    expect(content).not.toContain('500')
    expect(content).not.toContain('Error')
  })
})
