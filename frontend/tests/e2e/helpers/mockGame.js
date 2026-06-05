import http from 'node:http'

/** Minimal game API stub so Lobby can provision matches during E2E. */
export function startMockGameServer() {
  const server = http.createServer((req, res) => {
    if (req.method === 'POST' && req.url === '/api/v1/matches') {
      req.resume()
      res.writeHead(201, { 'Content-Type': 'application/json' })
      res.end('{}')
      return
    }
    res.writeHead(404).end()
  })

  return new Promise((resolve, reject) => {
    server.once('error', reject)
    server.listen(0, '127.0.0.1', () => {
      const address = server.address()
      const port = typeof address === 'object' && address ? address.port : 0
      resolve({
        server,
        baseUrl: `http://127.0.0.1:${port}`,
      })
    })
  })
}

export function stopMockGameServer(server) {
  return new Promise((resolve, reject) => {
    server.close((err) => {
      if (err) reject(err)
      else resolve()
    })
  })
}
