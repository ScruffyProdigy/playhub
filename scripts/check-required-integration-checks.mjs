#!/usr/bin/env node
/**
 * Ensures backend requiredPassChecks and frontend REQUIRED_INTEGRATION_CHECKS stay aligned.
 */
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')

function parseGoChecks(source) {
  const match = source.match(/requiredPassChecks = \[\]string\{([\s\S]*?)\}/)
  if (!match) throw new Error('requiredPassChecks block not found in developer.go')
  return [...match[1].matchAll(/"([^"]+)"/g)].map((m) => m[1])
}

function parseJSChecks(source) {
  const match = source.match(/REQUIRED_INTEGRATION_CHECKS = \[([\s\S]*?)\]/)
  if (!match) throw new Error('REQUIRED_INTEGRATION_CHECKS block not found in developers.js')
  return [...match[1].matchAll(/'([^']+)'/g)].map((m) => m[1])
}

const goSource = fs.readFileSync(path.join(root, 'backend/internal/store/developer.go'), 'utf8')
const jsSource = fs.readFileSync(path.join(root, 'frontend/src/lib/developers.js'), 'utf8')

const goChecks = parseGoChecks(goSource)
const jsChecks = parseJSChecks(jsSource)

const onlyGo = goChecks.filter((id) => !jsChecks.includes(id))
const onlyJs = jsChecks.filter((id) => !goChecks.includes(id))

if (onlyGo.length || onlyJs.length) {
  console.error('Required integration check IDs are out of sync:')
  if (onlyGo.length) console.error('  backend only:', onlyGo.join(', '))
  if (onlyJs.length) console.error('  frontend only:', onlyJs.join(', '))
  process.exit(1)
}

if (goChecks.length !== jsChecks.length) {
  console.error(`Count mismatch: backend ${goChecks.length}, frontend ${jsChecks.length}`)
  process.exit(1)
}

console.log(`OK: ${goChecks.length} required integration checks aligned (backend ↔ frontend)`)
