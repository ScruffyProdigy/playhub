import { mkdirSync, writeFileSync } from 'fs'
import { dirname, join } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const publicDir = join(__dirname, '..', 'public')
const envPath = join(publicDir, 'env.js')

const env = {
  REACT_APP_ENV: process.env.REACT_APP_ENV || 'local',
  REACT_APP_API_BASE_URL: process.env.REACT_APP_API_BASE_URL || 'http://localhost:8081',
}

mkdirSync(publicDir, { recursive: true })

const content = `window.env = {
  REACT_APP_ENV: "${env.REACT_APP_ENV}",
  REACT_APP_API_BASE_URL: "${env.REACT_APP_API_BASE_URL}"
};
`

writeFileSync(envPath, content)
console.log(`Generated ${envPath}`)
