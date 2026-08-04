import { cpSync, mkdirSync, rmSync, writeFileSync } from 'node:fs'
import { dirname } from 'node:path'

export function copyTree(src, dest, { dryRun = false } = {}) {
  if (dryRun) {
    return { src, dest, copied: false }
  }
  mkdirSync(dirname(dest), { recursive: true })
  rmSync(dest, { recursive: true, force: true })
  cpSync(src, dest, { recursive: true })
  return { src, dest, copied: true }
}

export function writeFile(filePath, content, { dryRun = false } = {}) {
  if (dryRun) {
    return { filePath, wrote: false }
  }
  mkdirSync(dirname(filePath), { recursive: true })
  writeFileSync(filePath, content)
  return { filePath, wrote: true }
}
