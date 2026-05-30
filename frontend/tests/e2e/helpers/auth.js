import { execSync } from 'node:child_process'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { expect } from '@playwright/test'

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../../..')

function psqlAvailable() {
  try {
    execSync('command -v psql', { stdio: 'ignore' })
    return true
  } catch {
    return false
  }
}

function queryMagicLinkToken(sql) {
  if (process.env.DATABASE_URL && psqlAvailable()) {
    return execSync(`psql "${process.env.DATABASE_URL}" -tAc "${sql}"`, {
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'pipe'],
    }).trim()
  }

  return execSync(
    `docker compose exec -T postgres psql -U app -d playhub -tAc "${sql}"`,
    {
      cwd: repoRoot,
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'pipe'],
    },
  ).trim()
}

export function latestMagicLinkToken(email) {
  const normalizedEmail = email.trim().toLowerCase().replace(/'/g, "''")
  const sql = `SELECT token FROM magic_links WHERE email='${normalizedEmail}' ORDER BY created_at DESC LIMIT 1`
  return queryMagicLinkToken(sql)
}

export function setUserDisplayName(email, displayName) {
  const normalizedEmail = email.trim().toLowerCase().replace(/'/g, "''")
  const escapedName = displayName.replace(/'/g, "''")
  const sql = `UPDATE users SET display_name='${escapedName}' WHERE email='${normalizedEmail}'`
  queryMagicLinkToken(sql)
}

export async function completeMagicSignIn(page, email) {
  await page.getByLabel('Email').fill(email)
  await page.getByRole('button', { name: 'Send magic link' }).click()
  await expect(page.getByText('Check your email for your sign-in link.')).toBeVisible()

  const token = latestMagicLinkToken(email)
  expect(token).toBeTruthy()

  await page.goto(`/auth/complete?token=${encodeURIComponent(token)}`)
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible({ timeout: 15000 })
}
