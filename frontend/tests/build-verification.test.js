import { test, expect, beforeAll } from 'vitest'
import { readFileSync, existsSync, readdirSync } from 'fs'
import { join, dirname } from 'path'
import { fileURLToPath } from 'url'
import { execSync } from 'child_process'

const __dirname = dirname(fileURLToPath(import.meta.url))
const frontendDir = join(__dirname, '..')
const distPath = join(frontendDir, 'dist')

beforeAll(() => {
  execSync('npm run build', { cwd: frontendDir, stdio: 'pipe', env: process.env })
}, 120000)

test('build artifacts exist', () => {
  expect(existsSync(distPath)).toBe(true)

  const indexPath = join(distPath, 'index.html')
  expect(existsSync(indexPath)).toBe(true)

  const indexContent = readFileSync(indexPath, 'utf-8')
  expect(indexContent).toContain('<title>JoinQuest</title>')
  expect(indexContent).toContain('<div id="root">')
  expect(indexContent).toContain('env.js')
  expect(indexContent).toContain('<base href="/"')
  // Absolute paths so /auth/complete and other deep links load assets from site root.
  expect(indexContent).toMatch(/src="\/env.js"/)
  expect(indexContent).toMatch(/src="\/assets\//)

  const envPath = join(distPath, 'env.js')
  expect(existsSync(envPath)).toBe(true)

  const envContent = readFileSync(envPath, 'utf-8')
  expect(envContent).toContain('window.env')
  expect(envContent).toContain('REACT_APP_ENV')
  expect(envContent).toContain('REACT_APP_API_BASE_URL')
})

test('build artifacts have correct structure', () => {
  const assetsPath = join(distPath, 'assets')
  expect(existsSync(assetsPath)).toBe(true)

  const cssFiles = readdirSync(assetsPath).filter((f) => f.endsWith('.css'))
  expect(cssFiles.length).toBeGreaterThan(0)

  const jsFiles = readdirSync(assetsPath).filter((f) => f.endsWith('.js'))
  expect(jsFiles.length).toBeGreaterThan(0)
})
