import { test, expect } from 'vitest'
import { readFileSync, existsSync } from 'fs'
import { join } from 'path'

test('build artifacts exist', () => {
  const distPath = join(process.cwd(), 'dist')
  
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
  const distPath = join(process.cwd(), 'dist')
  
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
