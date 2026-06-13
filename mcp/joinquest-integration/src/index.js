#!/usr/bin/env node

import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js'
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js'
import { loadConfig } from './config.js'
import { registerJoinQuestIntegrationTools } from './tools.js'

async function main() {
  const config = loadConfig()
  const server = new McpServer({
    name: 'joinquest-integration',
    version: '0.1.0',
  })

  registerJoinQuestIntegrationTools(server, config)

  const transport = new StdioServerTransport()
  await server.connect(transport)
  console.error('JoinQuest integration MCP server running (stdio)')
}

main().catch((error) => {
  console.error('JoinQuest integration MCP fatal error:', error.message)
  process.exit(1)
})
