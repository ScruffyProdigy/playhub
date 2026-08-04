#!/usr/bin/env node
import { runCreateCli } from 'joinquest'

runCreateCli(process.argv.slice(2))
  .then((code) => process.exit(code))
  .catch((err) => {
    console.error(err.message || err)
    process.exit(1)
  })
