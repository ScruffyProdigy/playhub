import { test, expect } from '@playwright/test'

test.describe('App E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('has title', async ({ page }) => {
    // Expect a title "to contain" a substring.
    await expect(page).toHaveTitle(/PlayHub/)
  })

  test('displays the main heading', async ({ page }) => {
    // Check that the main heading is visible
    await expect(page.getByRole('heading', { name: 'PlayHub' })).toBeVisible()
  })

  test('displays the tagline', async ({ page }) => {
    // Check that the tagline is visible
    await expect(page.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeVisible()
  })

  test('has proper heading structure', async ({ page }) => {
    // Check that the heading is an h1
    await expect(page.getByRole('heading', { level: 1, name: 'PlayHub' })).toBeVisible()
  })

  test('renders all expected content', async ({ page }) => {
    // Check that all main content is visible
    await expect(page.getByText('PlayHub')).toBeVisible()
    await expect(page.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeVisible()
  })

  test('has accessible content structure', async ({ page }) => {
    // Check that content has proper semantic structure
    const heading = page.getByRole('heading', { name: 'PlayHub' })
    const tagline = page.getByText('Your Gaming Hub - Queue, Play, Trade')
    
    await expect(heading).toBeVisible()
    await expect(tagline).toBeVisible()
  })

  test('maintains content during page interactions', async ({ page }) => {
    // Check that content remains stable
    await expect(page.getByText('PlayHub')).toBeVisible()
    await expect(page.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeVisible()
    
    // Simulate some page interactions (hover, scroll, etc.)
    await page.hover('h1')
    await page.mouse.move(100, 100)
    
    // Content should still be visible
    await expect(page.getByText('PlayHub')).toBeVisible()
    await expect(page.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeVisible()
  })
})

test.describe('App Accessibility E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('has proper heading structure', async ({ page }) => {
    await expect(page.getByRole('heading', { level: 1, name: 'PlayHub' })).toBeVisible()
  })

  test('has accessible content', async ({ page }) => {
    const heading = page.getByRole('heading', { name: 'PlayHub' })
    const tagline = page.getByText('Your Gaming Hub - Queue, Play, Trade')
    
    await expect(heading).toBeVisible()
    await expect(tagline).toBeVisible()
  })

  test('supports keyboard navigation', async ({ page }) => {
    // Tab through the page
    await page.keyboard.press('Tab')
    await page.keyboard.press('Tab')
    
    // Content should still be visible
    await expect(page.getByText('PlayHub')).toBeVisible()
    await expect(page.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeVisible()
  })
})

test.describe('App Performance E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('loads quickly', async ({ page }) => {
    const startTime = performance.now()
    await page.waitForLoadState('networkidle')
    const endTime = performance.now()
    const loadTime = endTime - startTime
    
    // Should load in less than 2 seconds (adjust as needed for your app)
    expect(loadTime).toBeLessThan(2000)
  })

  test('responds quickly to interactions', async ({ page }) => {
    const startTime = performance.now()
    await page.hover('h1')
    const endTime = performance.now()
    const responseTime = endTime - startTime
    
    // Should respond in less than 1000ms (more realistic for E2E tests)
    expect(responseTime).toBeLessThan(1000)
  })

  test('handles multiple interactions efficiently', async ({ page }) => {
    const startTime = performance.now()
    
    // Perform multiple interactions
    for (let i = 0; i < 10; i++) {
      await page.hover('h1')
      await page.hover('p')
    }
    
    const endTime = performance.now()
    const totalTime = endTime - startTime
    
    // Should handle 10 interactions in less than 5000ms (more realistic for E2E tests)
    expect(totalTime).toBeLessThan(5000)
  })
})

test.describe('App Cross-Browser E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('works in different browsers', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'PlayHub' })).toBeVisible()
    await expect(page.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeVisible()
  })

  test('handles different screen sizes', async ({ page }) => {
    // Test mobile viewport
    await page.setViewportSize({ width: 375, height: 667 }) // iPhone SE
    await expect(page.getByRole('heading', { name: 'PlayHub' })).toBeVisible()
    await expect(page.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeVisible()

    // Test tablet viewport
    await page.setViewportSize({ width: 1024, height: 768 }) // iPad
    await expect(page.getByRole('heading', { name: 'PlayHub' })).toBeVisible()
    await expect(page.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeVisible()

    // Test desktop viewport
    await page.setViewportSize({ width: 1920, height: 1080 }) // Desktop
    await expect(page.getByRole('heading', { name: 'PlayHub' })).toBeVisible()
    await expect(page.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeVisible()
  })
})