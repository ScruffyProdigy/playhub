import { test, expect } from 'vitest'
import { readFileSync, existsSync } from 'fs'
import { join, dirname } from 'path'
import { fileURLToPath } from 'url'

// Get the directory of this test file
const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

test('build artifacts exist', () => {
  // Look for dist directory relative to the frontend directory
  const frontendDir = join(__dirname, '..')
  const distPath = join(frontendDir, 'dist')
  
  // Debug information
  console.log('Current working directory:', process.cwd())
  console.log('Test file directory:', __dirname)
  console.log('Frontend directory:', frontendDir)
  console.log('Dist path:', distPath)
  console.log('Dist exists:', existsSync(distPath))
  
  // Check that dist directory exists
  expect(existsSync(distPath)).toBe(true)
  
  // Check that index.html exists
  const indexPath = join(distPath, 'index.html')
  expect(existsSync(indexPath)).toBe(true)
  
  // Check that index.html has basic content
  const indexContent = readFileSync(indexPath, 'utf-8')
  expect(indexContent).toContain('<title>PlayHub</title>')
  expect(indexContent).toContain('<div id="root">')
  expect(indexContent).toContain('env.js')
  
  // Check that env.js exists and has content
  const envPath = join(distPath, 'env.js')
  expect(existsSync(envPath)).toBe(true)
  
  const envContent = readFileSync(envPath, 'utf-8')
  expect(envContent).toContain('window.env')
  expect(envContent).toContain('REACT_APP_ENV')
  expect(envContent).toContain('REACT_APP_API_BASE_URL')
})

test('build artifacts have correct structure', () => {
  // Look for dist directory relative to the frontend directory
  const frontendDir = join(__dirname, '..')
  const distPath = join(frontendDir, 'dist')
  
  // Check for assets directory
  const assetsPath = join(distPath, 'assets')
  expect(existsSync(assetsPath)).toBe(true)
  
  // Check for CSS file
  const cssFiles = require('fs').readdirSync(assetsPath).filter(f => f.endsWith('.css'))
  expect(cssFiles.length).toBeGreaterThan(0)
  
  // Check for JS file
  const jsFiles = require('fs').readdirSync(assetsPath).filter(f => f.endsWith('.js'))
  expect(jsFiles.length).toBeGreaterThan(0)
})
