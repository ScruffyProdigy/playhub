#!/usr/bin/env node

import { spawn } from 'node:child_process'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'

// Cursor injects ELECTRON_RUN_AS_NODE and a bundled node on PATH, which breaks stdio MCP.
delete process.env.ELECTRON_RUN_AS_NODE

const root = join(dirname(fileURLToPath(import.meta.url)), '..')
const entry = join(root, 'src/index.js')
const node = process.env.JOINQUEST_NODE || process.execPath

const child = spawn(node, [entry], { stdio: 'inherit', env: process.env })
child.on('exit', (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal)
    return
  }
  process.exit(code ?? 1)
})
