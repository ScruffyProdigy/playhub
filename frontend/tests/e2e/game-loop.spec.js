import { test, expect } from '@playwright/test'
import { signInWithEmailCode } from './helpers/auth.js'
import {
  clearDemoMatchmakingState,
  restorePrimaryGameHandoffUrls,
  setPrimaryGameHandoffUrls,
} from './helpers/db.js'
import { startMockGameServer, stopMockGameServer } from './helpers/mockGame.js'
import {
  expectMatchedBanner,
  expectNoActiveQueueBanner,
  expectWaitingBanner,
  joinRockPaperQueue,
  readLaunchMatchId,
  returnFromMatch,
} from './helpers/queue.js'

test.describe('Game loop', () => {
  /** @type {import('node:http').Server | undefined} */
  let mockServer

  test.beforeAll(async () => {
    const mock = await startMockGameServer()
    mockServer = mock.server
    setPrimaryGameHandoffUrls(mock.baseUrl, mock.baseUrl)
  })

  test.afterAll(async () => {
    restorePrimaryGameHandoffUrls()
    if (mockServer) {
      await stopMockGameServer(mockServer)
    }
  })

  test.beforeEach(() => {
    clearDemoMatchmakingState()
  })

  test('match, return home, and re-queue with a fresh session', async ({ browser }) => {
    const emailA = `loop-a-${Date.now()}@example.com`
    const emailB = `loop-b-${Date.now()}@example.com`

    const contextA = await browser.newContext()
    const contextB = await browser.newContext()
    const pageA = await contextA.newPage()
    const pageB = await contextB.newPage()

    try {
      await pageA.goto('/')
      await signInWithEmailCode(pageA, emailA)
      await pageB.goto('/')
      await signInWithEmailCode(pageB, emailB)

      await joinRockPaperQueue(pageA)
      await expectWaitingBanner(pageA)

      await joinRockPaperQueue(pageB)
      await expectMatchedBanner(pageA)
      await expectMatchedBanner(pageB)

      const firstMatchId = await readLaunchMatchId(pageA)
      expect(await readLaunchMatchId(pageB)).toBe(firstMatchId)

      await returnFromMatch(pageA, firstMatchId)
      await expectNoActiveQueueBanner(pageA)

      await returnFromMatch(pageB, firstMatchId)
      await expectNoActiveQueueBanner(pageB)

      await joinRockPaperQueue(pageA)
      await expectWaitingBanner(pageA)
      await joinRockPaperQueue(pageB)
      await expectMatchedBanner(pageA)
      await expectMatchedBanner(pageB)

      const secondMatchId = await readLaunchMatchId(pageA)
      expect(secondMatchId).not.toBe(firstMatchId)
      expect(await readLaunchMatchId(pageB)).toBe(secondMatchId)
    } finally {
      await contextA.close()
      await contextB.close()
    }
  })
})
