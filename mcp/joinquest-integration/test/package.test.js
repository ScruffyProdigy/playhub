import assert from 'node:assert/strict'
import { test } from 'node:test'
import pkg from '../package.json' with { type: 'json' }

test('package exposes default npx bin matching scoped name', () => {
  assert.equal(pkg.bin['mcp-integration'], 'src/index.js')
})
