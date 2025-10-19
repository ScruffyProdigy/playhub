import { test, expect } from '@playwright/test'

test.describe('Environment Debug Tests', () => {
  test('debug environment configuration', async ({ page }) => {
    try {
      console.log('Navigating to page...')
      await page.goto('/', { timeout: 60000 })
      
      console.log('Waiting for page to load...')
      await page.waitForLoadState('domcontentloaded', { timeout: 30000 })
      
      // Check if env.js script is loaded
      console.log('Checking for env.js script...')
      const envScript = await page.locator('script[src*="env.js"]').count()
      console.log('Number of env.js scripts found:', envScript)
      
      // Try to access env.js directly
      console.log('Trying to fetch env.js directly...')
      const envResponse = await page.request.get('/env.js')
      console.log('env.js response status:', envResponse.status())
      if (envResponse.ok()) {
        const envContent = await envResponse.text()
        console.log('env.js content:', envContent)
      } else {
        console.log('Failed to fetch env.js')
      }
      
      // Check window.env
      console.log('Checking window.env...')
      const envConfig = await page.evaluate(() => {
        return {
          hasWindowEnv: typeof window.env !== 'undefined',
          windowEnv: window.env,
          windowKeys: Object.keys(window).filter(k => k.includes('env'))
        }
      })
      console.log('Window.env check:', envConfig)
      
      // Check if the script executed
      console.log('Checking if env script executed...')
      const scriptExecuted = await page.evaluate(() => {
        return typeof window.env !== 'undefined' && window.env !== null
      })
      console.log('Script executed:', scriptExecuted)
      
      // Take a screenshot for debugging
      await page.screenshot({ path: 'env-debug.png' })
      
      // Basic assertions
      expect(envScript).toBeGreaterThan(0)
      expect(envResponse.status()).toBe(200)
      
    } catch (error) {
      console.error('Test failed:', error)
      await page.screenshot({ path: 'env-debug-failure.png' })
      throw error
    }
  })
})
