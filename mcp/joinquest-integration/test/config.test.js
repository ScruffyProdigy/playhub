import { describe, it } from 'node:test'
import assert from 'node:assert/strict'
import { DEFAULT_JOINQUEST_API_URL, loadConfig } from '../src/config.js'

describe('loadConfig', () => {
  it('defaults JOINQUEST_API_URL to production', () => {
    const original = process.env.JOINQUEST_API_URL
    delete process.env.JOINQUEST_API_URL
    process.env.JOINQUEST_API_KEY = 'lq_dev_test'

    try {
      const config = loadConfig()
      assert.equal(config.apiUrl, DEFAULT_JOINQUEST_API_URL)
      assert.equal(config.authHeader, 'Bearer lq_dev_test')
    } finally {
      delete process.env.JOINQUEST_API_KEY
      if (original === undefined) {
        delete process.env.JOINQUEST_API_URL
      } else {
        process.env.JOINQUEST_API_URL = original
      }
    }
  })

  it('honors JOINQUEST_API_URL override for local dev', () => {
    process.env.JOINQUEST_API_URL = 'http://localhost:8080/graphql'
    process.env.JOINQUEST_API_KEY = 'lq_dev_test'

    try {
      const config = loadConfig()
      assert.equal(config.apiUrl, 'http://localhost:8080/graphql')
    } finally {
      delete process.env.JOINQUEST_API_KEY
      delete process.env.JOINQUEST_API_URL
    }
  })
})
